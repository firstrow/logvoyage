package main

import (
	"bytes"
	"encoding/json"
	"github.com/firstrow/logvoyage/common"
	"github.com/firstrow/logvoyage/tcp_server"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

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
		// TODO: Find recover solution.
		log.Fatal("%s", err)
	}
	defer resp.Body.Close()
	// Read body to close connection
	// If dont read body golang will keep connection open
	ioutil.ReadAll(resp.Body)
	log.Print("Message sent to Elastic: " + json)
}

func main() {
	log.Print("Initializing server")
	server := tcp_server.NewServer(":9999")
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
