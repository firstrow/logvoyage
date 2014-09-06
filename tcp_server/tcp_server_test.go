package tcp_server

import (
	"net"
	"testing"
	"time"
)

func buildTestServer() *server {
	return NewServer("localhost:9999")
}

func Test_accepting_new_client_callback(t *testing.T) {
	server := buildTestServer()

	var messageReceived bool
	var newClient bool
	var connectinClosed bool

	server.OnNewClient(func(c *Client) {
		newClient = true
	})
	server.OnNewMessage(func(c *Client, message string) {
		messageReceived = true
	})
	server.OnClientConnectionClosed(func(c *Client, err error) {
		connectinClosed = true
	})
	go server.Listen()

	// Wait for server
	time.Sleep(10 * time.Millisecond)

	conn, err := net.Dial("tcp", "localhost:9999")
	if err != nil {
		t.Fatal("Failed to connect to test server")
	}
	conn.Write([]byte("Test message\n"))
	conn.Close()

	// Wait for server
	time.Sleep(10 * time.Millisecond)

	if newClient != true {
		t.Error("New-client callback not called")
	}
	if messageReceived != true {
		t.Error("Message not received by server")
	}
	if connectinClosed != true {
		t.Error("Connection-closed callback not called")
	}
}
