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

	host := flag.String("host", defaultHost, "Host to open server. Set to `localhost` to accept only local connections.")
	port := flag.String("port", defaultPort, "Port to accept new connections. Default value: "+defaultPort)
	flag.Parse()

	server := tcp_server.New(*host + ":" + *port)
	server.OnNewClient(func(c *tcp_server.Client) {
		log.Print("New client")
	})

	// Receives new message and send it to Elastic server
	server.OnNewMessage(func(c *tcp_server.Client, message string) {
		indexName, err := getIndexName(message)
		if err != nil {
			// TODO: Log error
		} else {
			message = common.RemoveApiKey(message)
			message = strings.TrimSpace(message)

			// Send data to elastic
			record := &common.LogRecord{
				Message:  message,
				Datetime: time.Now().UTC(),
			}
			toElastic(indexName, record)
		}
	})

	server.OnClientConnectionClosed(func(c *tcp_server.Client, err error) {
		log.Print("Client disconnected")
	})
	server.Listen()
}

// Get users index name by apiKey
func getIndexName(message string) (string, error) {
	key, err := common.ExtractApiKey(message)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	user := common.FindUserByApiKey(key)
	if user == nil {
		log.Println("User not found")
		return "", errors.New("Error. User not found")
	}

	return user.GetIndexName(), nil
}

func toElastic(indexName string, record *common.LogRecord) {
	j, err := json.Marshal(record)
	if err != nil {
		log.Print("Error encoding message to JSON")
	} else {
		common.SendToElastic(indexName+"/logs", "POST", j)
	}
}
