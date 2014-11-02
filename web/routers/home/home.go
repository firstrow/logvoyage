package home

import (
	"github.com/belogik/goes"
	"github.com/firstrow/logvoyage/common"
	"github.com/firstrow/logvoyage/web/render"
	"net/http"
	"net/url"
	"strconv"
)

func getConnection() *goes.Connection {
	return goes.NewConnection("localhost", "9200")
}

// Search logs in elastic.
func search(text string, indexes []string) []goes.Hit {
	conn := getConnection()

	if len(text) > 0 {
		strconv.Quote(text)
	} else {
		text = "*"
	}

	var query = map[string]interface{}{
		"query": map[string]interface{}{
			"query_string": map[string]string{
				"default_field": "message",
				"query":         text,
			},
		},
		"from": 0,
		"size": 100,
		"sort": map[string]string{
			"datetime": "desc",
		},
	}

	extraArgs := make(url.Values, 1)
	searchResults, err := conn.Search(query, indexes, []string{"logs"}, extraArgs)

	if err != nil {
		panic(err)
	}

	return searchResults.Hits.Hits
}

func Index(req *http.Request, r *render.Render) {
	query_text := req.URL.Query().Get("q")
	user := common.FindUserByEmail(r.Context["email"].(string))
	r.HTML("index", render.ViewData{
		"logs":       search(query_text, []string{user.GetIndexName()}),
		"query_text": query_text,
	})
}
