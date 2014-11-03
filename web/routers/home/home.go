package home

import (
	"github.com/belogik/goes"
	"net/http"
	"net/url"
	"strconv"

	"github.com/firstrow/logvoyage/common"
	"github.com/firstrow/logvoyage/web/render"
)

func getConnection() *goes.Connection {
	return goes.NewConnection("localhost", "9200")
}

// Search logs in elastic.
func search(text string, indexes []string, from int) goes.Response {
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
		"from": from,
		"size": 10,
		"sort": map[string]string{
			"datetime": "desc",
		},
	}

	extraArgs := make(url.Values, 1)
	searchResults, err := conn.Search(query, indexes, []string{"logs"}, extraArgs)

	if err != nil {
		panic(err)
	}

	return searchResults
}

func Index(req *http.Request, r *render.Render) {
	query_text := req.URL.Query().Get("q")
	user := common.FindUserByEmail(r.Context["email"].(string))

	// Pagination
	from := 0
	page, _ := strconv.Atoi(req.URL.Query().Get("page"))
	if page > 0 {
		from = 10 * page
	}

	// Load records
	data := search(query_text, []string{user.GetIndexName()}, from)

	r.HTML("index", render.ViewData{
		"logs":       data.Hits.Hits,
		"total":      data.Hits.Total,
		"query_text": query_text,
	})
}
