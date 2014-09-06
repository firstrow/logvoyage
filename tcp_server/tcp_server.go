package tcp_server

import (
	"bufio"
	"log"
	"net"
)

type Client struct {
	conn     net.Conn
	Server   *server
	incoming chan string // Channel for incoming data from client
}

type server struct {
	clients                  []*Client
	address                  string        // Address to open connection: localhost:9999
	joins                    chan net.Conn // Channel for new connections
	onNewClientCallback      func(c *Client)
	onClientConnectionClosed func(c *Client, err error)
	onNewMessage             func(c *Client, message string)
}

func (c *Client) listen() {
	reader := bufio.NewReader(c.conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Print("Closing connection")
			c.conn.Close()
			c.Server.onClientConnectionClosed(c, err)
			return
		}
		c.Server.onNewMessage(c, message)
	}
}

func (s *server) OnNewClient(callback func(c *Client)) {
	s.onNewClientCallback = callback
}

func (s *server) OnClientConnectionClosed(callback func(c *Client, err error)) {
	s.onClientConnectionClosed = callback
}

func (s *server) OnNewMessage(callback func(c *Client, message string)) {
	s.onNewMessage = callback
}

func (s *server) newClient(conn net.Conn) {
	log.Print("Adding new connection")
	client := &Client{
		conn:   conn,
		Server: s,
	}
	go client.listen()
	s.onNewClientCallback(client)
}

// Listens new connections channel and creating new client
func (s *server) listenChannels() {
	for {
		select {
		case conn := <-s.joins:
			s.newClient(conn)
		}
	}
}

// Start network server
func (s *server) Listen() {
	log.Print("Starting listen channels...")
	go s.listenChannels()

	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		log.Fatal("Error starting TCP server.")
	}

	for {
		conn, _ := listener.Accept()
		s.joins <- conn
	}
}

// Create new tcp server instance
func NewServer(address string) *server {
	log.Print("Creating server with address " + address)
	server := &server{
		address: address,
		joins:   make(chan net.Conn),
	}
	return server
}
