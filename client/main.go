// Clinet is linux daemon to collect logs and send it to LogVoyage service.
// If service is down Client will write to file and try to send them again
// when service is up.
//
// Client can accept messages in two ways: by tcp and http interface
package main

import (
	"errors"
	"flag"
	"github.com/firstrow/logvoyage/tcp_server"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
)

var (
	logVoyageConnection net.Conn
	httpDsn             string
	tcpDsn              string
	logVoyageDsn        string
	emptyStringError    = errors.New("Received empty string")
)

func main() {
	flag.StringVar(&httpDsn, "httpDsn", "localhost:27078", "Host and port to start local HTTP server.")
	flag.StringVar(&tcpDsn, "tcpDsn", "localhost:27079", "Host and port to start local TCP server.")
	flag.StringVar(&logVoyageDsn, "logVoyageDsn", "localhost:27077", "LogVoyage server host and port.")

	flag.Parse()

	// Start servers
	err := connectLogVoyage(logVoyageDsn)
	if err != nil {
		log.Fatal("Error connecting go LogVoyage server.")
	}

	startServers()

	// Where to defer connection close?
	defer logVoyageConnection.Close()
}

// Setup persistent tcp connection to LogVoyage server
func connectLogVoyage(dsn string) error {
	conn, err := net.Dial("tcp", dsn)
	if err != nil {
		return errors.New("Error connecting to LogVoyage server.")
	}
	logVoyageConnection = conn
	return nil
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("Received new http request")
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Print("Error reading http request body")
	} else {
		send(bytes)
	}
}

func startHttpServer(httpDsn string) {
	log.Printf("Starting http server at %s", httpDsn)
	http.HandleFunc("/", httpHandler)
	http.ListenAndServe(httpDsn, nil)
}

func startTcpServer(tcpDsn string) {
	log.Printf("Starting tcp server at %s", tcpDsn)
	server := tcp_server.New(tcpDsn)
	server.OnNewMessage(func(c *tcp_server.Client, message string) {
		send([]byte(message))
	})
	server.Listen()
}

// Sends message to LogVoyage server
func send(message []byte) {
	text, err := prepareMessage(string(message))
	if err == nil {
		_, err := logVoyageConnection.Write([]byte(text))
		if err != nil {
			log.Print("Error sending " + text)
			backupMessage(text)
			startTryingReconnect()
		}
	}
}

func backupMessage(text string) {
	// Write message to file
}

func startTryingReconnect() {
	log.Print("Reconnecting")
	err := connectLogVoyage(logVoyageDsn)
	if err == nil {
		log.Print("Connected")
	}
}

// Trims message and adds \n to end so it can be properly read by the server
func prepareMessage(message string) (string, error) {
	result := strings.TrimSpace(message)
	if len(result) > 0 {
		return result + "\n", nil
	}
	return "", emptyStringError
}

// Starts http and tcp servers
func startServers() {
	go startHttpServer(httpDsn)
	startTcpServer(tcpDsn)
}
