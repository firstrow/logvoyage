// Backend server - main part of LogVoyage service.
// It accepts connections from "Client", parses string and pushes it to ElasticSearch index
package main

import (
	"encoding/json"
	"flag"
	"github.com/firstrow/logvoyage/common"
	"github.com/firstrow/tcp_server"
	"log"
	"runtime"
	"strings"
	"time"
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
		message = strings.TrimSpace(message)
		// Send data to elastic
		record := &common.LogRecord{
			Datetime: time.Now().UTC(),
			Message:  message,
		}
		toElastic(record)
	})

	server.OnClientConnectionClosed(func(c *tcp_server.Client, err error) {
		log.Print("Client disconnected")
	})
	server.Listen()
}

func toElastic(record *common.LogRecord) {
	j, err := json.Marshal(record)
	if err != nil {
		log.Print("Error encoding message to JSON")
	} else {
		common.SendToElastic("firstrow/logs", "POST", j)
	}
}
