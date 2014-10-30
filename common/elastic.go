package common

import (
	"bytes"
	"github.com/belogik/goes"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	ES_HOST = "localhost"
	ES_PORT = "9200"
)

func GetConnection() *goes.Connection {
	return goes.NewConnection(ES_HOST, ES_PORT)
}

// Send raw bytes to elastic search server
func SendToElastic(url string, method string, b []byte) (string, error) {
	eurl := "http://localhost:9200/"
	eurl += url

	req, err := http.NewRequest(method, eurl, bytes.NewBuffer(b))
	if err != nil {
		log.Print("Error creating POST request to storage")
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read body to close connection
	// If dont read body golang will keep connection open
	r, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	return string(r), nil
}
