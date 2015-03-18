// Copyright 2013 Belogik. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goes_test

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/belogik/goes"
)

var (
	ES_HOST = "localhost"
	ES_PORT = "9200"
)

func getConnection() (conn *goes.Connection) {
	h := os.Getenv("TEST_ELASTICSEARCH_HOST")
	if h == "" {
		h = ES_HOST
	}

	p := os.Getenv("TEST_ELASTICSEARCH_PORT")
	if p == "" {
		p = ES_PORT
	}

	conn = goes.NewConnection(h, p)

	return
}

func ExampleConnection_CreateIndex() {
	conn := getConnection()

	mapping := map[string]interface{}{
		"settings": map[string]interface{}{
			"index.number_of_shards":   1,
			"index.number_of_replicas": 0,
		},
		"mappings": map[string]interface{}{
			"_default_": map[string]interface{}{
				"_source": map[string]interface{}{
					"enabled": true,
				},
				"_all": map[string]interface{}{
					"enabled": false,
				},
			},
		},
	}

	resp, err := conn.CreateIndex("test", mapping)

	if err != nil {
		panic(err)
	}

	fmt.Printf("%s", resp)
}

func ExampleConnection_DeleteIndex() {
	conn := getConnection()
	resp, err := conn.DeleteIndex("yourinde")

	if err != nil {
		panic(err)
	}

	fmt.Printf("%s", resp)
}

func ExampleConnection_RefreshIndex() {
	conn := getConnection()
	resp, err := conn.RefreshIndex("yourindex")

	if err != nil {
		panic(err)
	}

	fmt.Printf("%s", resp)
}

func ExampleConnection_Search() {
	conn := getConnection()

	var query = map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": map[string]interface{}{
					"match_all": map[string]interface{}{},
				},
			},
		},
		"from":   0,
		"size":   100,
		"fields": []string{"onefield"},
		"filter": map[string]interface{}{
			"range": map[string]interface{}{
				"somefield": map[string]interface{}{
					"from":          "some date",
					"to":            "some date",
					"include_lower": false,
					"include_upper": false,
				},
			},
		},
	}

	extraArgs := make(url.Values, 1)

	searchResults, err := conn.Search(query, []string{"someindex"}, []string{""}, extraArgs)

	if err != nil {
		panic(err)
	}

	fmt.Printf("%s", searchResults)
}

func ExampleConnection_Index() {
	conn := getConnection()

	d := goes.Document{
		Index: "twitter",
		Type:  "tweet",
		Fields: map[string]interface{}{
			"user":    "foo",
			"message": "bar",
		},
	}

	extraArgs := make(url.Values, 1)
	extraArgs.Set("ttl", "86400000")

	response, err := conn.Index(d, extraArgs)

	if err != nil {
		panic(err)
	}

	fmt.Printf("%s", response)
}

func ExampleConnection_Delete() {
	conn := getConnection()

	//[create index, index document ...]

	d := goes.Document{
		Index: "twitter",
		Type:  "tweet",
		Id:    "1",
		Fields: map[string]interface{}{
			"user": "foo",
		},
	}

	response, err := conn.Delete(d, url.Values{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s", response)
}

func ExampleConnectionOverrideHttpClient() {
	tr := &http.Transport{
		ResponseHeaderTimeout: 1 * time.Second,
	}
	cl := &http.Client{
		Transport: tr,
	}
	conn := getConnection()
	conn.WithClient(cl)

	fmt.Printf("%v\n", conn.Client)
}
