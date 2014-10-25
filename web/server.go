package main

import (
	"fmt"
	"github.com/belogik/goes"
	"github.com/codegangsta/martini-contrib/render"
	"github.com/go-martini/martini"
	"html/template"
	"net/url"
	"time"
)

func getConnection() *goes.Connection {
	return goes.NewConnection("localhost", "9200")
}

func search(text string, indexes []string) []goes.Hit {
	conn := getConnection()

	var query = map[string]interface{}{
		"query": map[string]interface{}{
			"query_string": map[string]string{
				"default_field": "message",
				"query":         text,
			},
		},
		"from": 0,
		"size": 1000,
		"sort": map[string]string{
			"datetime": "desc",
		},
	}

	extraArgs := make(url.Values, 1)
	searchResults, err := conn.Search(query, indexes, []string{"logs"}, extraArgs)

	if err != nil {
		panic(err)
	}

	fmt.Println("Took")
	fmt.Println(searchResults.Took)

	return searchResults.Hits.Hits
}

func indexPage(r render.Render) {
	data := map[string]interface{}{"logs": search("*", []string{"firstrow"})}
	r.HTML(200, "index", data)
}

func main() {
	templateFunc := template.FuncMap{
		"FormatTimeToHuman": func(s string) string {
			t, _ := time.Parse(time.RFC3339Nano, s)
			return t.Format("2006-01-02 15:04:05") + " UTC"
		},
	}

	m := martini.Classic()
	m.Use(render.Renderer(render.Options{
		Funcs: []template.FuncMap{templateFunc},
	}))
	m.Use(martini.Static("../static"))

	m.Get("/", indexPage)
	m.Run()
}
