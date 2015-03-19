package main

import (
	"bitbucket.org/firstrow/logvoyage/common"
	"github.com/codegangsta/cli"
	"log"
	"os"
)

var createUsersIndex = cli.Command{
	Name:        "create_users_index",
	Usage:       "Will create `user` index in ES",
	Description: "",
	Action:      createUsersIndexFunc,
	Flags:       []cli.Flag{},
}

var deleteIndex = cli.Command{
	Name:        "delete_index",
	Usage:       "Will delete elastic search index",
	Description: "",
	Action:      deleteIndexFunc,
	Flags:       []cli.Flag{},
}

var createIndex = cli.Command{
	Name:        "create_index",
	Usage:       "create search index",
	Description: "",
	Action:      createIndexFunc,
	Flags:       []cli.Flag{},
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

func main() {
	app := cli.NewApp()
	app.Name = "LogVoyage"
	app.Commands = []cli.Command{
		createUsersIndex,
		deleteIndex,
		createIndex,
	}
	app.Run(os.Args)
}
