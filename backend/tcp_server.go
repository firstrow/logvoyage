package main

import (
	"log"

	"github.com/firstrow/tcp_server"
)

func initTcpServer() {
	server := tcp_server.New(tcpDsn)
	server.OnNewClient(func(c *tcp_server.Client) {
		log.Print("New client")
	})

	// Receives new message and send it to Elastic server
	server.OnNewMessage(func(c *tcp_server.Client, message string) {
		processMessage(message)
	})

	server.OnClientConnectionClosed(func(c *tcp_server.Client, err error) {
		log.Print("Client disconnected")
	})
	server.Listen()
}
