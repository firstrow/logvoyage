package common

import (
	"github.com/belogik/goes"
	"github.com/mitchellh/mapstructure"
	. "github.com/smartystreets/goconvey/convey"
	_ "log"
	"net/url"
	"testing"
	"time"
)

func TestElasticResponseToStruct(t *testing.T) {
	indexname := "testingindex"
	conn := GetConnection()
	defer conn.DeleteIndex(indexname)

	settings := `{
		"settings": {
			"index": {
				"number_of_shards": 5,
				"number_of_replicas": 1
			}
		},
		"mappings": {
			"user" : {
				"_source" : {"enabled" : true},
				"properties" : {
					"email" : {"type" : "string", "index": "not_analyzed" },
					"password" : {"type" : "string", "index": "not_analyzed" },
					"tokens" : {"type" : "string", "index": "not_analyzed" }
				}
			}
		}
	}`
	SendToElastic(indexname, "PUT", []byte(settings))

	doc := goes.Document{
		Index: indexname,
		Type:  "user",
		Fields: map[string]string{
			"email":    "test@localhost.loc",
			"password": "password",
			"apiKey":   "api_key_123",
		},
	}
	conn.Index(doc, url.Values{})

	time.Sleep(2 * time.Second)

	// Search user
	var query = map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": map[string]interface{}{
					"term": map[string]interface{}{
						"email": map[string]interface{}{
							"value": "test@localhost.loc",
						},
					},
				},
			},
		},
	}

	searchResults, err := conn.Search(query, []string{indexname}, []string{"user"}, url.Values{})

	if err != nil {
		t.Fatal(err.Error())
	}

	if searchResults.Hits.Total == 0 {
		t.Fatal("User not found. Probably insert error.")
	}

	user := &User{}

	mapstructure.Decode(searchResults.Hits.Hits[0].Source, user)

	Convey("It should populate user from goes search response", t, func() {
		So(user.Email, ShouldEqual, "test@localhost.loc")
		So(user.Password, ShouldEqual, "password")
		So(user.ApiKey, ShouldEqual, "api_key_123")
	})
}
