package commands

import (
	"github.com/codegangsta/cli"
	"log"

	"github.com/firstrow/logvoyage/backend"
	"github.com/firstrow/logvoyage/common"
	"github.com/firstrow/logvoyage/web"
)

var CreateUsersIndex = cli.Command{
	Name:        "create_users_index",
	Usage:       "Will create `user` index in ES",
	Description: "",
	Action:      createUsersIndexFunc,
	Flags:       []cli.Flag{},
}

var DeleteIndex = cli.Command{
	Name:        "delete_index",
	Usage:       "Will delete elastic search index",
	Description: "",
	Action:      deleteIndexFunc,
	Flags:       []cli.Flag{},
}

var CreateIndex = cli.Command{
	Name:        "create_index",
	Usage:       "Create search index",
	Description: "",
	Action:      createIndexFunc,
	Flags:       []cli.Flag{},
}

var StartBackendServer = cli.Command{
	Name:   "backend",
	Usage:  "Starts TCP server to accept logs from clients",
	Action: startBackendServer,
}

var StartWebServer = cli.Command{
	Name:        "web",
	Usage:       "Starts web-ui server",
	Description: "",
	Action:      startWebServer,
	Flags:       []cli.Flag{},
}

var StartAll = cli.Command{
	Name:        "start-all",
	Usage:       "Starts backend and web server",
	Description: "",
	Action:      startAll,
	Flags:       []cli.Flag{},
}

func startBackendServer(c *cli.Context) {
	backend.Start()
}

func startWebServer(c *cli.Context) {
	web.Start()
}

func startAll(c *cli.Context) {
	go backend.Start()
	web.Start()
}

func createUsersIndexFunc(c *cli.Context) {
	log.Println("Creating users index in ElasticSearch")
	settings := `{
		"settings": {
			"index": {
				"number_of_shards": 5,
				"number_of_replicas": 1,
				"refresh_interval" : "1s"
			}
		},
		"mappings": {
			"user" : {
				"_source" : {"enabled" : true},
				"properties" : {
					"email" : {"type" : "string", "index": "not_analyzed" },
					"password" : {"type" : "string", "index": "not_analyzed" },
					"apiKey" : {"type" : "string", "index": "not_analyzed" }
				}
			}
		}
	}`
	result, _ := common.SendToElastic("users", "PUT", []byte(settings))
	log.Println(result)
}

func createIndexFunc(c *cli.Context) {
	settings := `{
		"settings": {
			"index": {
				"number_of_shards": 5,
				"number_of_replicas": 1,
				"refresh_interval" : "2s"
			}
		}
	}`
	result, _ := common.SendToElastic(c.Args()[0], "PUT", []byte(settings))
	log.Println(result)
}

func deleteIndexFunc(c *cli.Context) {
	if len(c.Args()) > 0 {
		for _, name := range c.Args() {
			result, _ := common.SendToElastic(name, "DELETE", nil)
			log.Println(result)
		}
	} else {
		log.Println("Provide index name. e.g: index1, index2, ...")
	}
}
