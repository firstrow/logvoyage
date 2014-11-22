package common

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/belogik/goes"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	ES_HOST = "localhost"
	ES_PORT = "9200"
)

var (
	ErrSendingElasticSearchRequest = errors.New("Error sending request to ES.")
	ErrDecodingJson                = errors.New("Error decoding ES response")
)

func GetConnection() *goes.Connection {
	return goes.NewConnection(ES_HOST, ES_PORT)
}

type IndexMapping map[string]map[string]map[string]interface{}

// Retuns list of types available in search index
func GetTypes(index string) ([]string, error) {
	var mapping IndexMapping
	result, err := SendToElastic(index+"/_mapping", "GET", []byte{})
	if err != nil {
		return nil, ErrSendingElasticSearchRequest
	}
	err = json.Unmarshal([]byte(result), &mapping)
	if err != nil {
		return nil, ErrDecodingJson
	}

	keys := []string{}
	for k := range mapping[index]["mappings"] {
		keys = append(keys, k)
	}
	return keys, nil
}

// Send raw bytes to elastic search serve	r
func SendToElastic(url string, method string, b []byte) (string, error) {
	eurl := "http://" + ES_HOST + ":" + ES_PORT + "/"
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
