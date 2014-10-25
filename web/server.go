package main

import (
	"fmt"
	"github.com/belogik/goes"
	"github.com/codegangsta/martini-contrib/render"
	"github.com/go-martini/martini"
	"net/url"
)

func getConnection() *goes.Connection {
	return goes.NewConnection("localhost", "9200")
}

func search() []goes.Hit {
	conn := getConnection()

	var query = map[string]interface{}{
		"query": map[string]interface{}{
			"query_string": map[string]string{
				"default_field": "message",
				"query":         "safari",
			},
		},
		"from": 0,
		"size": 1000,
	}

	extraArgs := make(url.Values, 1)
	searchResults, err := conn.Search(query, []string{"firstrow"}, []string{"logs"}, extraArgs)

	if err != nil {
		panic(err)
	}

	fmt.Println("Took")
	fmt.Println(searchResults.Took)

	return searchResults.Hits.Hits
}

func indexPage(r render.Render) {
	data := map[string]interface{}{"logs": search()}
	r.HTML(200, "index", data)
}

func main() {
	m := martini.Classic()
	m.Use(render.Renderer())
	m.Use(martini.Static("../static"))

	m.Get("/", indexPage)
	m.Run()
}
