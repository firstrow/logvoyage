package common

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/belogik/goes"
)

const (
	DEFAULT_ES_HOST = "127.0.0.1"
	DEFAULT_ES_PORT = "9200"
)

var (
	ErrSendingElasticSearchRequest = errors.New("Error sending request to ES.")
	ErrCreatingHttpRequest         = errors.New("Could not create http.NewRequest")
	ErrReadResponse                = errors.New("Could not read ES response")
	ErrDecodingJson                = errors.New("Error decoding ES response")
)

func getElasticHostPort() (string, string) {
	es_host := DEFAULT_ES_HOST
	es_port := DEFAULT_ES_PORT
	if es_env := os.Getenv("LOGVOYAGE_ES"); len(es_env) > 1 {
		es_params := strings.Split(es_env, ":")
		if len(es_params[0]) > 0 {
			es_host = es_params[0]
		}
		if len(es_params[1]) > 0 {
			es_port = es_params[1]
		}
	}
	return es_host, es_port
}

func GetConnection() *goes.Connection {
	es_host, es_port := getElasticHostPort()
	return goes.NewConnection(es_host, es_port)
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
	_, err := SendToElastic(fmt.Sprintf("%s/%s", index, logType), "DELETE", nil)
	if err != nil {
		log.Println(err.Error())
	}
}

// Send raw bytes to elastic search server
// TODO: Bulk processing
func SendToElastic(url string, method string, b []byte) (string, error) {
	es_host, es_port := getElasticHostPort()
	eurl := fmt.Sprintf("http://%s:%s/%s", es_host, es_port, url)

	req, err := http.NewRequest(method, eurl, bytes.NewBuffer(b))
	if err != nil {
		return "", err
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
