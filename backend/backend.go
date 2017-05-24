// Backend server - main part of LogVoyage service.
// It accepts connections from "Client", parses string and pushes it to ElasticSearch index
package backend

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	c "github.com/firstrow/logvoyage/configuration"
	"github.com/firstrow/logvoyage/common"
	"github.com/firstrow/logvoyage/web_socket"
	"github.com/garyburd/redigo/redis"
)

var (
	tcpDsn    = ":27077"
	httpDsn   = ":27078"
	redisConn redis.Conn

	errUserNotFound = errors.New("Error. User not found")
)

func Start(userTcpDsn, userHttpDsn string) {
	if userTcpDsn != "" {
		tcpDsn = userTcpDsn
	}
	if userHttpDsn != "" {
		httpDsn = userHttpDsn
	}

	log.Println("Initializing server")

	initRedis()
	go initTimers()
	go initBacklog()
	go initTcpServer()
	initHttpServer()
}

func initRedis() {
	r, err := redis.Dial("tcp", c.ReadConf().Redis.GetURI())
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
		log.Println("Error extracting index name and type", err.Error())
		switch err {
		case common.ErrSendingElasticSearchRequest:
			toBacklog(origMessage)
		case errUserNotFound:
			log.Println("Backend: user not found.")
		}
	} else {
		message = common.RemoveApiKey(message)

		log.Println("Sending message to elastic")

		err = toElastic(indexName, logType, buildMessageStruct(message))
		if err == common.ErrSendingElasticSearchRequest {
			toBacklog(origMessage)
		} else {
			increaseCounter(indexName)
		}
		toRedis(indexName, logType, message)
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
		log.Println("Error encoding message to JSON")
	} else {
		_, err := common.SendToElastic(fmt.Sprintf("%s/%s", indexName, logType), "POST", j)
		if err != nil {
			return err
		}
	}
	return nil
}

func toRedis(indexName string, logType string, msg string) {
	var message web_socket.RedisMessage
	log.Println("Sending message to redis")
	message = web_socket.RedisMessage{ApiKey: indexName, Data: map[string]string{
		"type":     "log_message",
		"log_type": logType,
		"message":  msg,
	}}
	message.Send(redisConn)
}
