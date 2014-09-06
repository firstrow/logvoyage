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
	"time"
)

var (
	logVoyageConnection net.Conn
	emptyStringError    = errors.New("Received empty string")
)

func main() {
	httpDsn := flag.String("httpHost", "localhost:27078", "Host and port to start local HTTP server.")
	tcpDsn := flag.String("tcpDsn", "localhost:27079", "Host and port to start local TCP server.")
	logVoyageDsn := flag.String("logVoyage", "localhost:27077", "LogVoyage server host and port.")

	flag.Parse()
	initialize(*httpDsn, *tcpDsn, *logVoyageDsn)
}

func initialize(httpDsn string, tcpDsn string, logVoyageDsn string) {
	// Connect to to LogVoyage server
	conn, err := net.Dial("tcp", logVoyageDsn)
	if err != nil {
		log.Print("Error connecting to LogVoyage server. Will retry in 10 seconds.")
		time.Sleep(10 * time.Second)
		initialize(httpDsn, tcpDsn, logVoyageDsn)
		return
	}
	// Set package variable
	logVoyageConnection = conn
	defer conn.Close()

	go startHttpServer(httpDsn)
	startTcpServer(tcpDsn)
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
	server := tcp_server.NewServer(tcpDsn)
	server.OnNewClient(func(c *tcp_server.Client) {})

	server.OnNewMessage(func(c *tcp_server.Client, message string) {
		b := []byte(message)
		go send(b)
	})

	server.OnClientConnectionClosed(func(c *tcp_server.Client, err error) {})
	server.Listen()
}

// Sends message to LogVoyage server
func send(message []byte) {
	text, err := prepareMessage(string(message))
	if err == nil {
		// TODO: Handle write.
		// Create restore-log.
		_, err := logVoyageConnection.Write([]byte(text))
		if err != nil {
			log.Print("Connection with LogVoyage server lost. Will try again after 10 sec.")
		}
	}
}

// Trims message and adds \n to end
func prepareMessage(message string) (string, error) {
	result := strings.TrimSpace(message)
	if len(result) > 0 {
		return result + "\n", nil
	}
	return "", emptyStringError
}
