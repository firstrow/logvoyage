package common

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
)

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
