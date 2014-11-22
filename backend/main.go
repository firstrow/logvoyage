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
	defaultHost = ""
	defaultPort = "27077"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	log.Print("Initializing server")

	// Initalize counter timer
	go initTimers()

	host := flag.String("host", defaultHost, "Host to open server. Set to `localhost` to accept only local connections.")
	port := flag.String("port", defaultPort, "Port to accept new connections. Default value: "+defaultPort)
	flag.Parse()

	server := tcp_server.New(*host + ":" + *port)
	server.OnNewClient(func(c *tcp_server.Client) {
		log.Print("New client")
	})

	// Receives new message and send it to Elastic server
	// Message examples:
	// apiKey Some text
	// apiKey {message: "Some text", field:"value", ...}
	server.OnNewMessage(func(c *tcp_server.Client, message string) {
		indexName, logType, err := extractIndexAndType(message)
		if err != nil {
			// TODO: Log error
		} else {
			message = common.RemoveApiKey(message)
			message = strings.TrimSpace(message)

			var data map[string]interface{}
			err := json.Unmarshal([]byte(message), &data)

			if err == nil {
				// Save parsed json
				data["datetime"] = time.Now().UTC()
				toElastic(indexName, logType, data)
			} else {
				// Could not parse json, save entire message.
				record := &common.LogRecord{
					Message:  message,
					Datetime: time.Now().UTC(),
				}
				toElastic(indexName, logType, record)
			}

			increaseCounter(indexName)
		}
	})

	server.OnClientConnectionClosed(func(c *tcp_server.Client, err error) {
		log.Print("Client disconnected")
	})
	server.Listen()
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
	} else {
		user := common.FindUserByApiKey(key)
		if user == nil {
			return "", "", errors.New("Error. User not found")
		}
		userIndexNameCache[user.GetIndexName()] = user.GetIndexName()
		return user.GetIndexName(), logType, nil
	}
}

// Sends data to elastic index
func toElastic(indexName string, logType string, record interface{}) {
	j, err := json.Marshal(record)
	if err != nil {
		log.Print("Error encoding message to JSON")
	} else {
		common.SendToElastic(indexName+"/"+logType, "POST", j)
	}
}
