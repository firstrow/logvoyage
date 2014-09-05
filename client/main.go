// Clinet is linux daemon to collect logs and send it to LogVoyage service.
// If service is down Client will write to file and try to send them again
// when service is up.
//
// Client can accept messages in two ways: by tcp and http interface
package main

import (
	"errors"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
)

var (
	log_voyage_connection net.Conn
	emptyStringError      = errors.New("Received empty string")
	httpServerPort        = "9998"
	httpServerHost        = "localhost"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:9999")
	if err != nil {
		log.Printf("Error connection to server: %s", err)
	}
	// Set package variable
	log_voyage_connection = conn
	defer conn.Close()

	startHttpServer()
	startTcpServer()

	log.Print("Ready.")
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

func startHttpServer() {
	connectionString := httpServerHost + ":" + httpServerPort
	log.Printf("Starting server at %s", connectionString)
	http.HandleFunc("/", httpHandler)
	http.ListenAndServe(connectionString, nil)
}

func startTcpServer() {

}

// Sends message to LogVoyage server
func send(message []byte) {
	text, err := prepareMessage(string(message))
	if err == nil {
		log_voyage_connection.Write([]byte(text))
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
