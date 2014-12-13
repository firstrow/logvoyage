package common

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/belogik/goes"
)

const (
	ES_HOST = "localhost"
	ES_PORT = "9200"
)

var (
	ErrSendingElasticSearchRequest = errors.New("Error sending request to ES.")
	ErrCreatingHttpRequest         = errors.New("Could not create http.NewRequest")
	ErrReadResponse                = errors.New("Could not read ES response")
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

// Count documents in collection
func CountTypeDocs(index string, logType string) float64 {
	result, err := SendToElastic(fmt.Sprintf("%s/%s/_count", index, logType), "GET", nil)
	if err != nil {
		return 0
	}

	var m map[string]interface{}
	err = json.Unmarshal([]byte(result), &m)
	if err != nil {
		return 0
	}
	return m["count"].(float64)
}

func DeleteType(index string, logType string) {
	_, err := SendToElastic(fmt.Sprint("%s/%s", index, logType), "DELETE", nil)
	if err != nil {
		// TODO: Log error
	}
}

// Send raw bytes to elastic search server
func SendToElastic(url string, method string, b []byte) (string, error) {
	eurl := fmt.Sprintf("http://%s:%s/%s", ES_HOST, ES_PORT, url)

	req, err := http.NewRequest(method, eurl, bytes.NewBuffer(b))
	if err != nil {
		return "", ErrCreatingHttpRequest
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", ErrSendingElasticSearchRequest
	}
	defer resp.Body.Close()

	// Read body to close connection
	// If dont read body golang will keep connection open
	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", ErrReadResponse
	}

	return string(r), nil
}
