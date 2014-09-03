package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type Client struct {
	conn     net.Conn
	incoming chan string
}

type Server struct {
	clients []*Client
	joins   chan net.Conn // Channel for new connections
}

// Add new connection to server
func (server *Server) addConnection(conn net.Conn) {
	log.Print("Adding new connection")
	client := &Client{
		conn: conn,
	}

	go func() {
		reader := bufio.NewReader(conn)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				log.Print("Closing connection")
				conn.Close()
				return
			}
			log.Print(line)
		}
	}()

	server.clients = append(server.clients, client)
}

func NewServer() *Server {
	server := &Server{
		joins: make(chan net.Conn),
	}

	go func() {
		for {
			select {
			// On new connection
			case conn := <-server.joins:
				server.addConnection(conn)
			}
		}
	}()

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
