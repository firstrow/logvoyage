// Backend server - main part of LogVoyage service.
// It accepts connections from "Client", parses string and pushes it to ElasticSearch index
package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"runtime"
	"strings"
	"time"

	"github.com/firstrow/logvoyage/common"
	"github.com/garyburd/redigo/redis"
)

var (
	tcpDsn          string
	httpDsn         string
	errUserNotFound = errors.New("Error. User not found")
	redisConn       redis.Conn
)

func init() {
	flag.StringVar(&tcpDsn, "tcpDsn", ":27077", "Host and port to accept tcp connections.")
	flag.StringVar(&httpDsn, "httpDsn", ":27078", "Host and port to accept http messages.")
	flag.Parse()
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	log.Print("Initializing server")

	initRedis()
	go initTimers()
	go initBacklog()
	go initTcpServer()
	initHttpServer()
}

func initRedis() {
	r, err := redis.Dial("tcp", ":6379")
	if err != nil {
		log.Fatal("Cannot connect to redis")
	}
	r.Flush()
	redisConn = r
}

// Process text message from tcp or http client
// Extract user api key, check send message to search index.
// Message examples:
// apiKey@logType Some text
// apiKey@logType {message: "Some text", field:"value", ...}
func processMessage(message string) {
	origMessage := message
	indexName, logType, err := extractIndexAndType(message)
	if err != nil {
		switch err {
		case common.ErrSendingElasticSearchRequest:
			toBacklog(origMessage)
		case errUserNotFound:
			log.Println("Backend: user not found.")
		}
	} else {
		message = common.RemoveApiKey(message)
		message = strings.TrimSpace(message)

		err = toElastic(indexName, logType, buildMessageStruct(message))
		if err == common.ErrSendingElasticSearchRequest {
			toBacklog(origMessage)
		} else {
			increaseCounter(indexName)
		}
		toWebSocket(indexName, message)
	}
}

// Stores [apiKey]indexName
var userIndexNameCache = make(map[string]string)

// Get users index name by apiKey
func extractIndexAndType(message string) (string, string, error) {
	key, logType, err := common.ExtractApiKey(message)
	if err != nil {
		return "", "", err
	}

	if indexName, ok := userIndexNameCache[key]; ok {
		return indexName, logType, nil
	}

	user, err := common.FindUserByApiKey(key)
	if err != nil {
		return "", "", err
	}
	if user == nil {
		return "", "", errUserNotFound
	}
	userIndexNameCache[user.GetIndexName()] = user.GetIndexName()
	return user.GetIndexName(), logType, nil
}

// Prepares message to be inserted into ES.
// Builds struct based on message.
func buildMessageStruct(message string) interface{} {
	var data map[string]interface{}
	err := json.Unmarshal([]byte(message), &data)

	if err == nil {
		// Save parsed json
		data["datetime"] = time.Now().UTC()
		return data
	} else {
		// Could not parse json, save entire message.
		return common.LogRecord{
			Message:  message,
			Datetime: time.Now().UTC(),
		}
	}
}

// Sends data to elastic index
func toElastic(indexName string, logType string, record interface{}) error {
	j, err := json.Marshal(record)
	if err != nil {
		log.Print("Error encoding message to JSON")
	} else {
		_, err := common.SendToElastic(fmt.Sprintf("%s/%s", indexName, logType), "POST", j)
		if err != nil {
			return err
		}
	}
	return nil
}

func toWebSocket(indexName string, message string) {

}
