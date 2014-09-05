package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/firstrow/logvoyage/common"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	_ "time"
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
	// If dont read body
	ioutil.ReadAll(resp.Body)
	log.Print("Message sent")
}

type Client struct {
	conn     net.Conn
	incoming chan string // Channel for incoming data from client
}

type Server struct {
	clients []*Client
	joins   chan net.Conn // Channel for new connections
}

func (client *Client) listen() {
	reader := bufio.NewReader(client.conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Print("Closing connection")
			// TODO: Remove client from list
			client.conn.Close()
			return
		}
		// Send data to elastic
		// parse line(lines)
		record := &common.LogRecord{
			Message: line,
		}
		json, err := json.Marshal(record)
		sendToElastic(string(json))
	}
}

// Add new connection to server
// TODO: Write authentication by key
func (server *Server) addConnection(conn net.Conn) {
	log.Print("Adding new connection")
	client := &Client{
		conn: conn,
	}
	go client.listen()
	server.clients = append(server.clients, client)
}

// Listen channels
func (server *Server) listen() {
	for {
		select {
		case conn := <-server.joins:
			server.addConnection(conn)
		}
	}
}

func NewServer() *Server {
	server := &Server{
		joins: make(chan net.Conn),
	}
	go server.listen()
	return server
}

func main() {
	fmt.Print("Initializing server")
	listener, err := net.Listen("tcp", ":9999")
	if err != nil {
		log.Fatal("Error starting TCP server.")
	}

	server := NewServer()

	for {
		conn, _ := listener.Accept()
		server.joins <- conn
	}
}
