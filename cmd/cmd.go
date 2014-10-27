package main

import (
	"github.com/codegangsta/cli"
	"github.com/firstrow/logvoyage/common"
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
			"users" : {
				"_source" : {"enabled" : true},
				"properties" : {
					"email" : {"type" : "string", "index": "not_analyzed" },
					"password" : {"type" : "string", "index": "not_analyzed" },
					"tokens" : {"type" : "string", "index": "not_analyzed" }
				}
			}
		}
	}`
	result, _ := common.SendToElastic("users", "PUT", []byte(settings))
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
	}
	app.Run(os.Args)
}
