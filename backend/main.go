// Backend server - main part of LogVoyage service.
// It accepts connections from "Client", parses string and pushes it to ElasticSearch index
package main

import (
	"encoding/json"
	"errors"
	"flag"
	"log"
	"runtime"
	"strings"
	"time"

	"github.com/firstrow/logvoyage/common"
	"github.com/firstrow/tcp_server"
)

var (
	host            = ""
	port            = "27077"
	errUserNotFound = errors.New("Error. User not found")
)

func init() {
	flag.StringVar(&host, "host", "", "Host to open server. Set to `localhost` to accept only local connections.")
	flag.StringVar(&port, "port", "27077", "Port to accept new connections. Default value: ")
	flag.Parse()
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	log.Print("Initializing server")

	// Initalize counter timer
	go initTimers()
	go initBacklog()
	go initTcpServer()
}

func initTcpServer() {
	server := tcp_server.New(host + ":" + port)
	server.OnNewClient(func(c *tcp_server.Client) {
		log.Print("New client")
	})

	// Receives new message and send it to Elastic server
	// Message examples:
	// apiKey@logType Some text
	// apiKey@logType {message: "Some text", field:"value", ...}
	server.OnNewMessage(func(c *tcp_server.Client, message string) {
		processMessage(message)
	})

	server.OnClientConnectionClosed(func(c *tcp_server.Client, err error) {
		log.Print("Client disconnected")
	})
	server.Listen()
}

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

		err = toElastic(indexName, logType, buildMessage(message))
		if err == common.ErrSendingElasticSearchRequest {
			toBacklog(origMessage)
		} else {
			increaseCounter(indexName)
		}
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
// Builds object based on message.
func buildMessage(message string) interface{} {
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
		_, err := common.SendToElastic(indexName+"/"+logType, "POST", j)
		if err != nil {
			return err
		}
	}
	return nil
}
