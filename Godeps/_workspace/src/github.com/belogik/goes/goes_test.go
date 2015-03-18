// Copyright 2013 Belogik. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goes

import (
	. "gopkg.in/check.v1"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"
)

var (
	ES_HOST = "localhost"
	ES_PORT = "9200"
)

// Hook up gocheck into the gotest runner.
func Test(t *testing.T) { TestingT(t) }

type GoesTestSuite struct{}

var _ = Suite(&GoesTestSuite{})

func (s *GoesTestSuite) SetUpTest(c *C) {
	h := os.Getenv("TEST_ELASTICSEARCH_HOST")
	if h != "" {
		ES_HOST = h
	}

	p := os.Getenv("TEST_ELASTICSEARCH_PORT")
	if p != "" {
		ES_PORT = p
	}
}

func (s *GoesTestSuite) TestNewConnection(c *C) {
	conn := NewConnection(ES_HOST, ES_PORT)
	c.Assert(conn, DeepEquals, &Connection{ES_HOST, ES_PORT, http.DefaultClient})
}

func (s *GoesTestSuite) TestWithClient(c *C) {
	tr := &http.Transport{
		DisableCompression:    true,
		ResponseHeaderTimeout: 1 * time.Second,
	}
	cl := &http.Client{
		Transport: tr,
	}
	conn := NewConnection(ES_HOST, ES_PORT).WithClient(cl)

	c.Assert(conn, DeepEquals, &Connection{ES_HOST, ES_PORT, cl})
	c.Assert(conn.Client.Transport.(*http.Transport).DisableCompression, Equals, true)
	c.Assert(conn.Client.Transport.(*http.Transport).ResponseHeaderTimeout, Equals, 1*time.Second)
}

func (s *GoesTestSuite) TestUrl(c *C) {
	conn := NewConnection(ES_HOST, ES_PORT)

	r := Request{
		Conn:      conn,
		Query:     "q",
		IndexList: []string{"i"},
		TypeList:  []string{},
		method:    "GET",
		api:       "_search",
	}

	c.Assert(r.Url(), Equals, "http://"+ES_HOST+":"+ES_PORT+"/i/_search")

	r.IndexList = []string{"a", "b"}
	c.Assert(r.Url(), Equals, "http://"+ES_HOST+":"+ES_PORT+"/a,b/_search")

	r.TypeList = []string{"c", "d"}
	c.Assert(r.Url(), Equals, "http://"+ES_HOST+":"+ES_PORT+"/a,b/c,d/_search")

	r.ExtraArgs = make(url.Values, 1)
	r.ExtraArgs.Set("version", "1")
	c.Assert(r.Url(), Equals, "http://"+ES_HOST+":"+ES_PORT+"/a,b/c,d/_search?version=1")

	r.id = "1234"
	r.api = ""
	c.Assert(r.Url(), Equals, "http://"+ES_HOST+":"+ES_PORT+"/a,b/c,d/1234/?version=1")
}

func (s *GoesTestSuite) TestEsDown(c *C) {
	conn := NewConnection("a.b.c.d", "1234")

	var query = map[string]interface{}{"query": "foo"}

	r := Request{
		Conn:      conn,
		Query:     query,
		IndexList: []string{"i"},
		method:    "GET",
		api:       "_search",
	}
	_, err := r.Run()

	c.Assert(err, ErrorMatches, "Get http://a.b.c.d:1234/i/_search:(.*)lookup a.b.c.d: no such host")
}

func (s *GoesTestSuite) TestRunMissingIndex(c *C) {
	conn := NewConnection(ES_HOST, ES_PORT)

	var query = map[string]interface{}{"query": "foo"}

	r := Request{
		Conn:      conn,
		Query:     query,
		IndexList: []string{"i"},
		method:    "GET",
		api:       "_search",
	}
	_, err := r.Run()

	c.Assert(err.Error(), Equals, "[404] IndexMissingException[[i] missing]")
}

func (s *GoesTestSuite) TestCreateIndex(c *C) {
	indexName := "testcreateindexgoes"

	conn := NewConnection(ES_HOST, ES_PORT)
	defer conn.DeleteIndex(indexName)

	mapping := map[string]interface{}{
		"settings": map[string]interface{}{
			"index.number_of_shards":   1,
			"index.number_of_replicas": 0,
		},
		"mappings": map[string]interface{}{
			"_default_": map[string]interface{}{
				"_source": map[string]interface{}{
					"enabled": false,
				},
				"_all": map[string]interface{}{
					"enabled": false,
				},
			},
		},
	}

	resp, err := conn.CreateIndex(indexName, mapping)

	c.Assert(err, IsNil)
	c.Assert(resp.Acknowledged, Equals, true)
}

func (s *GoesTestSuite) TestDeleteIndexInexistantIndex(c *C) {
	conn := NewConnection(ES_HOST, ES_PORT)
	resp, err := conn.DeleteIndex("foobar")

	c.Assert(err.Error(), Equals, "[404] IndexMissingException[[foobar] missing]")
	c.Assert(resp, DeepEquals, Response{})
}

func (s *GoesTestSuite) TestDeleteIndexExistingIndex(c *C) {
	conn := NewConnection(ES_HOST, ES_PORT)

	indexName := "testdeleteindexexistingindex"

	_, err := conn.CreateIndex(indexName, map[string]interface{}{})

	c.Assert(err, IsNil)

	resp, err := conn.DeleteIndex(indexName)
	c.Assert(err, IsNil)

	expectedResponse := Response{}
	expectedResponse.Acknowledged = true
	resp.Raw = nil
	c.Assert(resp, DeepEquals, expectedResponse)
}

func (s *GoesTestSuite) TestRefreshIndex(c *C) {
	conn := NewConnection(ES_HOST, ES_PORT)
	indexName := "testrefreshindex"

	_, err := conn.CreateIndex(indexName, map[string]interface{}{})
	c.Assert(err, IsNil)

	_, err = conn.RefreshIndex(indexName)
	c.Assert(err, IsNil)

	_, err = conn.DeleteIndex(indexName)
	c.Assert(err, IsNil)
}

func (s *GoesTestSuite) TestOptimize(c *C) {
	conn := NewConnection(ES_HOST, ES_PORT)
	indexName := "testoptimize"

	conn.DeleteIndex(indexName)
	_, err := conn.CreateIndex(indexName, map[string]interface{}{})
	c.Assert(err, IsNil)

	// we must wait for a bit otherwise ES crashes
	time.Sleep(1 * time.Second)

	response, err := conn.Optimize([]string{indexName}, url.Values{"max_num_segments": []string{"1"}})
	c.Assert(err, IsNil)

	c.Assert(response.All.Indices[indexName].Primaries["docs"].Count, Equals, 0)

	_, err = conn.DeleteIndex(indexName)
	c.Assert(err, IsNil)
}

func (s *GoesTestSuite) TestBulkSend(c *C) {
	indexName := "testbulkadd"
	docType := "tweet"

	tweets := []Document{
		{
			Id:          "123",
			Index:       indexName,
			Type:        docType,
			BulkCommand: BULK_COMMAND_INDEX,
			Fields: map[string]interface{}{
				"user":    "foo",
				"message": "some foo message",
			},
		},

		{
			Id:          nil,
			Index:       indexName,
			Type:        docType,
			BulkCommand: BULK_COMMAND_INDEX,
			Fields: map[string]interface{}{
				"user":    "bar",
				"message": "some bar message",
			},
		},
	}

	conn := NewConnection(ES_HOST, ES_PORT)

	conn.DeleteIndex(indexName)
	_, err := conn.CreateIndex(indexName, nil)
	c.Assert(err, IsNil)

	response, err := conn.BulkSend(tweets)
	i := Item{
		Id:      "123",
		Type:    docType,
		Version: 1,
		Index:   indexName,
		Status:  201, //201 for indexing ( https://issues.apache.org/jira/browse/CONNECTORS-634 )
	}
	c.Assert(response.Items[0][BULK_COMMAND_INDEX], Equals, i)
	c.Assert(err, IsNil)

	_, err = conn.RefreshIndex(indexName)
	c.Assert(err, IsNil)

	var query = map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
	}

	searchResults, err := conn.Search(query, []string{indexName}, []string{}, url.Values{})
	c.Assert(err, IsNil)

	var expectedTotal uint64 = 2
	c.Assert(searchResults.Hits.Total, Equals, expectedTotal)

	extraDocId := ""
	checked := 0
	for _, hit := range searchResults.Hits.Hits {
		if hit.Source["user"] == "foo" {
			c.Assert(hit.Id, Equals, "123")
			checked++
		}

		if hit.Source["user"] == "bar" {
			c.Assert(len(hit.Id) > 0, Equals, true)
			extraDocId = hit.Id
			checked++
		}
	}
	c.Assert(checked, Equals, 2)

	docToDelete := []Document{
		{
			Id:          "123",
			Index:       indexName,
			Type:        docType,
			BulkCommand: BULK_COMMAND_DELETE,
		},
		{
			Id:          extraDocId,
			Index:       indexName,
			Type:        docType,
			BulkCommand: BULK_COMMAND_DELETE,
		},
	}

	response, err = conn.BulkSend(docToDelete)
	i = Item{
		Id:      "123",
		Type:    docType,
		Version: 2,
		Index:   indexName,
		Status:  200, //200 for updates
	}
	c.Assert(response.Items[0][BULK_COMMAND_DELETE], Equals, i)

	c.Assert(err, IsNil)

	_, err = conn.RefreshIndex(indexName)
	c.Assert(err, IsNil)

	searchResults, err = conn.Search(query, []string{indexName}, []string{}, url.Values{})
	c.Assert(err, IsNil)

	expectedTotal = 0
	c.Assert(searchResults.Hits.Total, Equals, expectedTotal)

	_, err = conn.DeleteIndex(indexName)
	c.Assert(err, IsNil)
}

func (s *GoesTestSuite) TestStats(c *C) {
	conn := NewConnection(ES_HOST, ES_PORT)
	indexName := "teststats"

	conn.DeleteIndex(indexName)
	_, err := conn.CreateIndex(indexName, map[string]interface{}{})
	c.Assert(err, IsNil)

	// we must wait for a bit otherwise ES crashes
	time.Sleep(1 * time.Second)

	response, err := conn.Stats([]string{indexName}, url.Values{})
	c.Assert(err, IsNil)

	c.Assert(response.All.Indices[indexName].Primaries["docs"].Count, Equals, 0)

	_, err = conn.DeleteIndex(indexName)
	c.Assert(err, IsNil)
}

func (s *GoesTestSuite) TestIndexWithFieldsInStruct(c *C) {
	indexName := "testindexwithfieldsinstruct"
	docType := "tweet"
	docId := "1234"

	conn := NewConnection(ES_HOST, ES_PORT)
	// just in case
	conn.DeleteIndex(indexName)

	_, err := conn.CreateIndex(indexName, map[string]interface{}{})
	c.Assert(err, IsNil)
	defer conn.DeleteIndex(indexName)

	d := Document{
		Index: indexName,
		Type:  docType,
		Id:    docId,
		Fields: struct {
			user    string
			message string
		}{
			"foo",
			"bar",
		},
	}

	extraArgs := make(url.Values, 1)
	extraArgs.Set("ttl", "86400000")
	response, err := conn.Index(d, extraArgs)
	c.Assert(err, IsNil)

	expectedResponse := Response{
		Index:   indexName,
		Id:      docId,
		Type:    docType,
		Version: 1,
	}

	response.Raw = nil
	c.Assert(response, DeepEquals, expectedResponse)
}

func (s *GoesTestSuite) TestIndexWithFieldsNotInMapOrStruct(c *C) {
	indexName := "testindexwithfieldsnotinmaporstruct"
	docType := "tweet"
	docId := "1234"

	conn := NewConnection(ES_HOST, ES_PORT)
	// just in case
	conn.DeleteIndex(indexName)

	_, err := conn.CreateIndex(indexName, map[string]interface{}{})
	c.Assert(err, IsNil)
	defer conn.DeleteIndex(indexName)

	d := Document{
		Index:  indexName,
		Type:   docType,
		Id:     docId,
		Fields: "test",
	}

	extraArgs := make(url.Values, 1)
	extraArgs.Set("ttl", "86400000")
	_, err = conn.Index(d, extraArgs)
	c.Assert(err, Not(IsNil))
}

func (s *GoesTestSuite) TestIndexIdDefined(c *C) {
	indexName := "testindexiddefined"
	docType := "tweet"
	docId := "1234"

	conn := NewConnection(ES_HOST, ES_PORT)
	// just in case
	conn.DeleteIndex(indexName)

	_, err := conn.CreateIndex(indexName, map[string]interface{}{})
	c.Assert(err, IsNil)
	defer conn.DeleteIndex(indexName)

	d := Document{
		Index: indexName,
		Type:  docType,
		Id:    docId,
		Fields: map[string]interface{}{
			"user":    "foo",
			"message": "bar",
		},
	}

	extraArgs := make(url.Values, 1)
	extraArgs.Set("ttl", "86400000")
	response, err := conn.Index(d, extraArgs)
	c.Assert(err, IsNil)

	expectedResponse := Response{
		Index:   indexName,
		Id:      docId,
		Type:    docType,
		Version: 1,
	}

	response.Raw = nil
	c.Assert(response, DeepEquals, expectedResponse)
}

func (s *GoesTestSuite) TestIndexIdNotDefined(c *C) {
	indexName := "testindexidnotdefined"
	docType := "tweet"

	conn := NewConnection(ES_HOST, ES_PORT)
	// just in case
	conn.DeleteIndex(indexName)

	_, err := conn.CreateIndex(indexName, map[string]interface{}{})
	c.Assert(err, IsNil)
	defer conn.DeleteIndex(indexName)

	d := Document{
		Index: indexName,
		Type:  docType,
		Fields: map[string]interface{}{
			"user":    "foo",
			"message": "bar",
		},
	}

	response, err := conn.Index(d, url.Values{})
	c.Assert(err, IsNil)

	c.Assert(response.Index, Equals, indexName)
	c.Assert(response.Type, Equals, docType)
	c.Assert(response.Version, Equals, 1)
	c.Assert(response.Id != "", Equals, true)
}

func (s *GoesTestSuite) TestDelete(c *C) {
	indexName := "testdelete"
	docType := "tweet"
	docId := "1234"

	conn := NewConnection(ES_HOST, ES_PORT)
	// just in case
	conn.DeleteIndex(indexName)

	_, err := conn.CreateIndex(indexName, map[string]interface{}{})
	c.Assert(err, IsNil)
	defer conn.DeleteIndex(indexName)

	d := Document{
		Index: indexName,
		Type:  docType,
		Id:    docId,
		Fields: map[string]interface{}{
			"user": "foo",
		},
	}

	_, err = conn.Index(d, url.Values{})
	c.Assert(err, IsNil)

	response, err := conn.Delete(d, url.Values{})
	c.Assert(err, IsNil)

	expectedResponse := Response{
		Found: true,
		Index: indexName,
		Type:  docType,
		Id:    docId,
		// XXX : even after a DELETE the version number seems to be incremented
		Version: 2,
	}
	response.Raw = nil
	c.Assert(response, DeepEquals, expectedResponse)

	response, err = conn.Delete(d, url.Values{})
	c.Assert(err, IsNil)

	expectedResponse = Response{
		Found: false,
		Index: indexName,
		Type:  docType,
		Id:    docId,
		// XXX : even after a DELETE the version number seems to be incremented
		Version: 3,
	}
	response.Raw = nil
	c.Assert(response, DeepEquals, expectedResponse)
}

func (s *GoesTestSuite) TestDeleteByQuery(c *C) {
	indexName := "testdeletebyquery"
	docType := "tweet"
	docId := "1234"

	conn := NewConnection(ES_HOST, ES_PORT)
	// just in case
	conn.DeleteIndex(indexName)

	_, err := conn.CreateIndex(indexName, map[string]interface{}{})
	c.Assert(err, IsNil)
	defer conn.DeleteIndex(indexName)

	d := Document{
		Index: indexName,
		Type:  docType,
		Id:    docId,
		Fields: map[string]interface{}{
			"user": "foo",
		},
	}

	_, err = conn.Index(d, url.Values{})
	c.Assert(err, IsNil)

	_, err = conn.RefreshIndex(indexName)
	c.Assert(err, IsNil)

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"match_all": map[string]interface{}{},
					},
				},
			},
		},
	}

	//should be 1 doc before delete by query
	response, err := conn.Search(query, []string{indexName}, []string{docType}, url.Values{})
	c.Assert(err, IsNil)
	c.Assert(response.Hits.Total, Equals, uint64(1))

	response, err = conn.Query(query, []string{indexName}, []string{docType}, "DELETE", url.Values{})

	c.Assert(err, IsNil)

	expectedResponse := Response{
		Found:   false,
		Index:   "",
		Type:    "",
		Id:      "",
		Version: 0,
	}
	response.Raw = nil
	c.Assert(response, DeepEquals, expectedResponse)

	//should be 0 docs after delete by query
	response, err = conn.Search(query, []string{indexName}, []string{docType}, url.Values{})
	c.Assert(err, IsNil)
	c.Assert(response.Hits.Total, Equals, uint64(0))
}

func (s *GoesTestSuite) TestGet(c *C) {
	indexName := "testget"
	docType := "tweet"
	docId := "111"
	source := map[string]interface{}{
		"f1": "foo",
		"f2": "foo",
	}

	conn := NewConnection(ES_HOST, ES_PORT)
	conn.DeleteIndex(indexName)

	_, err := conn.CreateIndex(indexName, map[string]interface{}{})
	c.Assert(err, IsNil)
	defer conn.DeleteIndex(indexName)

	d := Document{
		Index:  indexName,
		Type:   docType,
		Id:     docId,
		Fields: source,
	}

	_, err = conn.Index(d, url.Values{})
	c.Assert(err, IsNil)

	response, err := conn.Get(indexName, docType, docId, url.Values{})
	c.Assert(err, IsNil)

	expectedResponse := Response{
		Index:   indexName,
		Type:    docType,
		Id:      docId,
		Version: 1,
		Found:   true,
		Source:  source,
	}

	response.Raw = nil
	c.Assert(response, DeepEquals, expectedResponse)

	fields := make(url.Values, 1)
	fields.Set("fields", "f1")
	response, err = conn.Get(indexName, docType, docId, fields)
	c.Assert(err, IsNil)

	expectedResponse = Response{
		Index:   indexName,
		Type:    docType,
		Id:      docId,
		Version: 1,
		Found:   true,
		Fields: map[string]interface{}{
			"f1": []interface{}{"foo"},
		},
	}

	response.Raw = nil
	c.Assert(response, DeepEquals, expectedResponse)
}

func (s *GoesTestSuite) TestSearch(c *C) {
	indexName := "testsearch"
	docType := "tweet"
	docId := "1234"
	source := map[string]interface{}{
		"user":    "foo",
		"message": "bar",
	}

	conn := NewConnection(ES_HOST, ES_PORT)
	conn.DeleteIndex(indexName)

	_, err := conn.CreateIndex(indexName, map[string]interface{}{})
	c.Assert(err, IsNil)
	defer conn.DeleteIndex(indexName)

	d := Document{
		Index:  indexName,
		Type:   docType,
		Id:     docId,
		Fields: source,
	}

	_, err = conn.Index(d, url.Values{})
	c.Assert(err, IsNil)

	_, err = conn.RefreshIndex(indexName)
	c.Assert(err, IsNil)

	// I can feel my eyes bleeding
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"match_all": map[string]interface{}{},
					},
				},
			},
		},
	}
	response, err := conn.Search(query, []string{indexName}, []string{docType}, url.Values{})

	expectedHits := Hits{
		Total:    1,
		MaxScore: 1.0,
		Hits: []Hit{
			{
				Index:  indexName,
				Type:   docType,
				Id:     docId,
				Score:  1.0,
				Source: source,
			},
		},
	}

	c.Assert(response.Hits, DeepEquals, expectedHits)
}

func (s *GoesTestSuite) TestCount(c *C) {
	indexName := "testcount"
	docType := "tweet"
	docId := "1234"
	source := map[string]interface{}{
		"user":    "foo",
		"message": "bar",
	}

	conn := NewConnection(ES_HOST, ES_PORT)
	conn.DeleteIndex(indexName)

	_, err := conn.CreateIndex(indexName, map[string]interface{}{})
	c.Assert(err, IsNil)
	defer conn.DeleteIndex(indexName)

	d := Document{
		Index:  indexName,
		Type:   docType,
		Id:     docId,
		Fields: source,
	}

	_, err = conn.Index(d, url.Values{})
	c.Assert(err, IsNil)

	_, err = conn.RefreshIndex(indexName)
	c.Assert(err, IsNil)

	// I can feel my eyes bleeding
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"match_all": map[string]interface{}{},
					},
				},
			},
		},
	}
	response, err := conn.Count(query, []string{indexName}, []string{docType}, url.Values{})

	c.Assert(response.Count, Equals, 1)
}

func (s *GoesTestSuite) TestIndexStatus(c *C) {
	indexName := "testindexstatus"
	conn := NewConnection(ES_HOST, ES_PORT)
	conn.DeleteIndex(indexName)

	mapping := map[string]interface{}{
		"settings": map[string]interface{}{
			"index.number_of_shards":   1,
			"index.number_of_replicas": 1,
		},
	}

	_, err := conn.CreateIndex(indexName, mapping)
	c.Assert(err, IsNil)
	defer conn.DeleteIndex(indexName)

	// gives ES some time to do its job
	time.Sleep(1 * time.Second)

	response, err := conn.IndexStatus([]string{"testindexstatus"})
	c.Assert(err, IsNil)

	expectedShards := Shard{Total: 2, Successful: 1, Failed: 0}
	c.Assert(response.Shards, Equals, expectedShards)

	primarySizeInBytes := response.Indices[indexName].Index["primary_size_in_bytes"].(float64)
	sizeInBytes := response.Indices[indexName].Index["size_in_bytes"].(float64)

	c.Assert(primarySizeInBytes > 0, Equals, true)
	c.Assert(sizeInBytes > 0, Equals, true)

	expectedIndices := map[string]IndexStatus{
		indexName: {
			Index: map[string]interface{}{
				"primary_size_in_bytes": primarySizeInBytes,
				"size_in_bytes":         sizeInBytes,
			},
			Translog: map[string]uint64{
				"operations": 0,
			},
			Docs: map[string]uint64{
				"num_docs":     0,
				"max_doc":      0,
				"deleted_docs": 0,
			},
			Merges: map[string]interface{}{
				"current":               float64(0),
				"current_docs":          float64(0),
				"current_size_in_bytes": float64(0),
				"total":                 float64(0),
				"total_time_in_millis":  float64(0),
				"total_docs":            float64(0),
				"total_size_in_bytes":   float64(0),
			},
			Refresh: map[string]interface{}{
				"total":                float64(1),
				"total_time_in_millis": float64(0),
			},
			Flush: map[string]interface{}{
				"total":                float64(0),
				"total_time_in_millis": float64(0),
			},
		},
	}

	c.Assert(response.Indices, DeepEquals, expectedIndices)
}

func (s *GoesTestSuite) TestScroll(c *C) {
	indexName := "testscroll"
	docType := "tweet"

	tweets := []Document{
		{
			Id:          nil,
			Index:       indexName,
			Type:        docType,
			BulkCommand: BULK_COMMAND_INDEX,
			Fields: map[string]interface{}{
				"user":    "foo",
				"message": "some foo message",
			},
		},

		{
			Id:          nil,
			Index:       indexName,
			Type:        docType,
			BulkCommand: BULK_COMMAND_INDEX,
			Fields: map[string]interface{}{
				"user":    "bar",
				"message": "some bar message",
			},
		},

		{
			Id:          nil,
			Index:       indexName,
			Type:        docType,
			BulkCommand: BULK_COMMAND_INDEX,
			Fields: map[string]interface{}{
				"user":    "foo",
				"message": "another foo message",
			},
		},
	}

	conn := NewConnection(ES_HOST, ES_PORT)

	mapping := map[string]interface{}{
		"settings": map[string]interface{}{
			"index.number_of_shards":   1,
			"index.number_of_replicas": 0,
		},
	}

	defer conn.DeleteIndex(indexName)

	_, err := conn.CreateIndex(indexName, mapping)
	c.Assert(err, IsNil)

	_, err = conn.BulkSend(tweets)
	c.Assert(err, IsNil)

	_, err = conn.RefreshIndex(indexName)
	c.Assert(err, IsNil)

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"filtered": map[string]interface{}{
				"filter": map[string]interface{}{
					"term": map[string]interface{}{
						"user": "foo",
					},
				},
			},
		},
	}

	scan, err := conn.Scan(query, []string{indexName}, []string{docType}, "1m", 1)
	c.Assert(err, IsNil)
	c.Assert(len(scan.ScrollId) > 0, Equals, true)

	searchResults, err := conn.Scroll(scan.ScrollId, "1m")
	c.Assert(err, IsNil)

	// some data in first chunk
	c.Assert(searchResults.Hits.Total, Equals, uint64(2))
	c.Assert(len(searchResults.ScrollId) > 0, Equals, true)
	c.Assert(len(searchResults.Hits.Hits), Equals, 1)

	searchResults, err = conn.Scroll(searchResults.ScrollId, "1m")
	c.Assert(err, IsNil)

	// more data in second chunk
	c.Assert(searchResults.Hits.Total, Equals, uint64(2))
	c.Assert(len(searchResults.ScrollId) > 0, Equals, true)
	c.Assert(len(searchResults.Hits.Hits), Equals, 1)

	searchResults, err = conn.Scroll(searchResults.ScrollId, "1m")
	c.Assert(err, IsNil)

	// nothing in third chunk
	c.Assert(searchResults.Hits.Total, Equals, uint64(2))
	c.Assert(len(searchResults.ScrollId) > 0, Equals, true)
	c.Assert(len(searchResults.Hits.Hits), Equals, 0)
}

func (s *GoesTestSuite) TestAggregations(c *C) {
	indexName := "testaggs"
	docType := "tweet"

	tweets := []Document{
		{
			Id:          nil,
			Index:       indexName,
			Type:        docType,
			BulkCommand: BULK_COMMAND_INDEX,
			Fields: map[string]interface{}{
				"user":    "foo",
				"message": "some foo message",
				"age":     25,
			},
		},

		{
			Id:          nil,
			Index:       indexName,
			Type:        docType,
			BulkCommand: BULK_COMMAND_INDEX,
			Fields: map[string]interface{}{
				"user":    "bar",
				"message": "some bar message",
				"age":     30,
			},
		},

		{
			Id:          nil,
			Index:       indexName,
			Type:        docType,
			BulkCommand: BULK_COMMAND_INDEX,
			Fields: map[string]interface{}{
				"user":    "foo",
				"message": "another foo message",
			},
		},
	}

	conn := NewConnection(ES_HOST, ES_PORT)

	mapping := map[string]interface{}{
		"settings": map[string]interface{}{
			"index.number_of_shards":   1,
			"index.number_of_replicas": 0,
		},
	}

	defer conn.DeleteIndex(indexName)

	_, err := conn.CreateIndex(indexName, mapping)
	c.Assert(err, IsNil)

	_, err = conn.BulkSend(tweets)
	c.Assert(err, IsNil)

	_, err = conn.RefreshIndex(indexName)
	c.Assert(err, IsNil)

	query := map[string]interface{}{
		"aggs": map[string]interface{}{
			"user": map[string]interface{}{
				"terms": map[string]interface{}{
					"field": "user",
					"order": map[string]interface{}{
						"_term": "asc",
					},
				},
				"aggs": map[string]interface{}{
					"age": map[string]interface{}{
						"stats": map[string]interface{}{
							"field": "age",
						},
					},
				},
			},
			"age": map[string]interface{}{
				"stats": map[string]interface{}{
					"field": "age",
				},
			},
		},
	}

	resp, err := conn.Search(query, []string{indexName}, []string{docType}, url.Values{})

	user, ok := resp.Aggregations["user"]
	c.Assert(ok, Equals, true)

	c.Assert(len(user.Buckets()), Equals, 2)
	c.Assert(user.Buckets()[0].Key(), Equals, "bar")
	c.Assert(user.Buckets()[1].Key(), Equals, "foo")

	barAge := user.Buckets()[0].Aggregation("age")
	c.Assert(barAge["count"], Equals, 1.0)
	c.Assert(barAge["sum"], Equals, 30.0)

	fooAge := user.Buckets()[1].Aggregation("age")
	c.Assert(fooAge["count"], Equals, 1.0)
	c.Assert(fooAge["sum"], Equals, 25.0)

	age, ok := resp.Aggregations["age"]
	c.Assert(ok, Equals, true)

	c.Assert(age["count"], Equals, 2.0)
	c.Assert(age["sum"], Equals, 25.0+30.0)
}

func (s *GoesTestSuite) TestPutMapping(c *C) {
	indexName := "testputmapping"
	docType := "tweet"

	conn := NewConnection(ES_HOST, ES_PORT)
	// just in case
	conn.DeleteIndex(indexName)

	_, err := conn.CreateIndex(indexName, map[string]interface{}{})
	c.Assert(err, IsNil)
	defer conn.DeleteIndex(indexName)

	d := Document{
		Index: indexName,
		Type:  docType,
		Fields: map[string]interface{}{
			"user":    "foo",
			"message": "bar",
		},
	}

	response, err := conn.Index(d, url.Values{})
	c.Assert(err, IsNil)

	mapping := map[string]interface{}{
		"tweet": map[string]interface{}{
			"properties": map[string]interface{}{
				"count": map[string]interface{}{
					"type":  "integer",
					"index": "not_analyzed",
					"store": true,
				},
			},
		},
	}
	response, err = conn.PutMapping("tweet", mapping, []string{indexName})
	c.Assert(err, IsNil)

	c.Assert(response.Acknowledged, Equals, true)
	c.Assert(response.TimedOut, Equals, false)
}

func (s *GoesTestSuite) TestIndicesExist(c *C) {
	indices := []string{"testindicesexist"}

	conn := NewConnection(ES_HOST, ES_PORT)
	// just in case
	conn.DeleteIndex(indices[0])

	exists, err := conn.IndicesExist(indices)
	c.Assert(exists, Equals, false)

	_, err = conn.CreateIndex(indices[0], map[string]interface{}{})
	c.Assert(err, IsNil)
	defer conn.DeleteIndex(indices[0])
	time.Sleep(200 * time.Millisecond)

	exists, err = conn.IndicesExist(indices)
	c.Assert(exists, Equals, true)

	indices = append(indices, "nonexistent")
	exists, err = conn.IndicesExist(indices)
	c.Assert(exists, Equals, false)
}

func (s *GoesTestSuite) TestUpdate(c *C) {
	indexName := "testupdate"
	docType := "tweet"
	docId := "1234"

	conn := NewConnection(ES_HOST, ES_PORT)
	// just in case
	conn.DeleteIndex(indexName)

	_, err := conn.CreateIndex(indexName, map[string]interface{}{})
	c.Assert(err, IsNil)
	defer conn.DeleteIndex(indexName)

	d := Document{
		Index: indexName,
		Type:  docType,
		Id:    docId,
		Fields: map[string]interface{}{
			"user":    "foo",
			"message": "bar",
			"counter": 1,
		},
	}

	extraArgs := make(url.Values, 1)
	response, err := conn.Index(d, extraArgs)
	c.Assert(err, IsNil)
	time.Sleep(200 * time.Millisecond)

	expectedResponse := Response{
		Index:   indexName,
		Id:      docId,
		Type:    docType,
		Version: 1,
	}

	response.Raw = nil
	c.Assert(response, DeepEquals, expectedResponse)

	// Now that we have an ordinary document indexed, try updating it
	query := map[string]interface{}{
		"script": "ctx._source.counter += count",
		"params": map[string]interface{}{
			"count": 5,
		},
		"upsert": map[string]interface{}{
			"message": "candybar",
			"user":    "admin",
			"counter": 1,
		},
	}

	response, err = conn.Update(d, query, extraArgs)
	if err != nil && strings.Contains(err.(*SearchError).Msg, "dynamic scripting disabled") {
		c.Skip("Scripting is disabled on server, skipping this test")
		return
	}
	time.Sleep(200 * time.Millisecond)

	c.Assert(err, Equals, nil)

	response, err = conn.Get(indexName, docType, docId, url.Values{})
	c.Assert(err, Equals, nil)
	c.Assert(response.Source["counter"], Equals, float64(6))
	c.Assert(response.Source["user"], Equals, "foo")
	c.Assert(response.Source["message"], Equals, "bar")

	// Test another document, non-existent
	docId = "555"
	d.Id = docId
	response, err = conn.Update(d, query, extraArgs)
	c.Assert(err, Equals, nil)
	time.Sleep(200 * time.Millisecond)

	response, err = conn.Get(indexName, docType, docId, url.Values{})
	c.Assert(err, Equals, nil)
	c.Assert(response.Source["user"], Equals, "admin")
	c.Assert(response.Source["message"], Equals, "candybar")
}

func (s *GoesTestSuite) TestGetMapping(c *C) {
	indexName := "testmapping"
	docType := "tweet"

	conn := NewConnection(ES_HOST, ES_PORT)
	// just in case
	conn.DeleteIndex(indexName)

	_, err := conn.CreateIndex(indexName, map[string]interface{}{})
	c.Assert(err, IsNil)
	defer conn.DeleteIndex(indexName)

	time.Sleep(300 * time.Millisecond)

	response, err := conn.GetMapping([]string{docType}, []string{indexName})
	c.Assert(err, Equals, nil)
	c.Assert(len(response.Raw), Equals, 0)

	d := Document{
		Index: indexName,
		Type:  docType,
		Fields: map[string]interface{}{
			"user":    "foo",
			"message": "bar",
		},
	}

	response, err = conn.Index(d, url.Values{})
	c.Assert(err, IsNil)
	time.Sleep(200 * time.Millisecond)

	response, err = conn.GetMapping([]string{docType}, []string{indexName})
	c.Assert(err, Equals, nil)
	c.Assert(len(response.Raw), Not(Equals), 0)
}

func (s *GoesTestSuite) TestDeleteMapping(c *C) {
	indexName := "testdeletemapping"
	docType := "tweet"

	conn := NewConnection(ES_HOST, ES_PORT)
	// just in case
	conn.DeleteIndex(indexName)

	_, err := conn.CreateIndex(indexName, map[string]interface{}{})
	c.Assert(err, IsNil)
	defer conn.DeleteIndex(indexName)

	d := Document{
		Index: indexName,
		Type:  docType,
		Fields: map[string]interface{}{
			"user":    "foo",
			"message": "bar",
		},
	}

	response, err := conn.Index(d, url.Values{})
	c.Assert(err, IsNil)

	mapping := map[string]interface{}{
		"tweet": map[string]interface{}{
			"properties": map[string]interface{}{
				"count": map[string]interface{}{
					"type":  "integer",
					"index": "not_analyzed",
					"store": true,
				},
			},
		},
	}
	response, err = conn.PutMapping("tweet", mapping, []string{indexName})
	c.Assert(err, IsNil)
	time.Sleep(200 * time.Millisecond)

	response, err = conn.DeleteMapping("tweet", []string{indexName})
	c.Assert(err, IsNil)

	c.Assert(response.Acknowledged, Equals, true)
	c.Assert(response.TimedOut, Equals, false)
}

func (s *GoesTestSuite) TestAddAlias(c *C) {
	aliasName := "testAlias"
	indexName := "testalias_1"
	docType := "testDoc"
	docId := "1234"
	source := map[string]interface{}{
		"user":    "foo",
		"message": "bar",
	}

	conn := NewConnection(ES_HOST, ES_PORT)
	defer conn.DeleteIndex(indexName)

	_, err := conn.CreateIndex(indexName, map[string]interface{}{})
	c.Assert(err, IsNil)
	defer conn.DeleteIndex(indexName)

	d := Document{
		Index:  indexName,
		Type:   docType,
		Id:     docId,
		Fields: source,
	}

	// Index data
	_, err = conn.Index(d, url.Values{})
	c.Assert(err, IsNil)

	// Add alias
	_, err = conn.AddAlias(aliasName, []string{indexName})
	c.Assert(err, IsNil)

	// Get document via alias
	response, err := conn.Get(aliasName, docType, docId, url.Values{})
	c.Assert(err, IsNil)

	expectedResponse := Response{
		Index:   indexName,
		Type:    docType,
		Id:      docId,
		Version: 1,
		Found:   true,
		Source:  source,
	}

	response.Raw = nil
	c.Assert(response, DeepEquals, expectedResponse)
}

func (s *GoesTestSuite) TestRemoveAlias(c *C) {
	aliasName := "testAlias"
	indexName := "testalias_1"
	docType := "testDoc"
	docId := "1234"
	source := map[string]interface{}{
		"user":    "foo",
		"message": "bar",
	}

	conn := NewConnection(ES_HOST, ES_PORT)
	defer conn.DeleteIndex(indexName)

	_, err := conn.CreateIndex(indexName, map[string]interface{}{})
	c.Assert(err, IsNil)
	defer conn.DeleteIndex(indexName)

	d := Document{
		Index:  indexName,
		Type:   docType,
		Id:     docId,
		Fields: source,
	}

	// Index data
	_, err = conn.Index(d, url.Values{})
	c.Assert(err, IsNil)

	// Add alias
	_, err = conn.AddAlias(aliasName, []string{indexName})
	c.Assert(err, IsNil)

	// Remove alias
	_, err = conn.RemoveAlias(aliasName, []string{indexName})
	c.Assert(err, IsNil)

	// Get document via alias
	_, err = conn.Get(aliasName, docType, docId, url.Values{})
	c.Assert(err.Error(), Equals, "[404] IndexMissingException[["+aliasName+"] missing]")
}

func (s *GoesTestSuite) TestAliasExists(c *C) {
	index := "testaliasexist_1"
	alias := "testaliasexists"

	conn := NewConnection(ES_HOST, ES_PORT)
	// just in case
	conn.DeleteIndex(index)

	exists, err := conn.AliasExists(alias)
	c.Assert(exists, Equals, false)

	_, err = conn.CreateIndex(index, map[string]interface{}{})
	c.Assert(err, IsNil)
	defer conn.DeleteIndex(index)
	time.Sleep(200 * time.Millisecond)

	_, err = conn.AddAlias(alias, []string{index})
	c.Assert(err, IsNil)
	time.Sleep(200 * time.Millisecond)
	defer conn.RemoveAlias(alias, []string{index})

	exists, err = conn.AliasExists(alias)
	c.Assert(exists, Equals, true)
}
