// Backend server - main part of LogVoyage service.
// It accepts connections from "Client", parses string and pushes it to ElasticSearch index
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"github.com/firstrow/logvoyage/common"
	"github.com/firstrow/logvoyage/tcp_server"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var (
	defaultHost = ""
	defaultPort = "27077"
)

func main() {
	log.Print("Initializing server")

	host := flag.String("host", defaultHost, "Host to open server. Set to `localhost` to accept only local connections.")
	port := flag.String("port", defaultPort, "Port to accept new connections. Default value: "+defaultPort)
	flag.Parse()

	server := tcp_server.New(*host + ":" + *port)
	server.OnNewClient(func(c *tcp_server.Client) {
		log.Print("New client")
	})

	server.OnNewMessage(func(c *tcp_server.Client, message string) {
		message = strings.TrimSpace(message)
		// Send data to elastic
		record := &common.LogRecord{
			Message: message,
		}
		json, err := json.Marshal(record)
		if err != nil {
			log.Print("Error encoding message")
		}
		sendToElastic(string(json))
	})

	server.OnClientConnectionClosed(func(c *tcp_server.Client, err error) {
		log.Print("Client disconnected")
	})
	server.Listen()
}

func sendToElastic(json string) {
	url := "http://localhost:9200/firstrow/logs"

	var jsonStr = []byte(json)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Print("Error creating POST request to storage")
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// Here we can't send data to elastic.
		// Write to log. restore.
		log.Fatal("%s", err)
	}
	defer resp.Body.Close()
	// Read body to close connection
	// If dont read body golang will keep connection open
	ioutil.ReadAll(resp.Body)
	log.Print("Message sent to Elastic: " + json)
}
