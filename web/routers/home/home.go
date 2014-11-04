package home

import (
	"github.com/belogik/goes"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/firstrow/logvoyage/common"
	"github.com/firstrow/logvoyage/web/render"
)

// Search logs in elastic.
func search(text string, indexes []string, size int, from int) goes.Response {
	conn := common.GetConnection()

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
		"size": size,
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
	size := 10
	from := 0
	page, _ := strconv.Atoi(req.URL.Query().Get("p"))
	if page > 1 {
		from = 10 * page
	}

	log.Println("From is: ", from)

	// Load records
	data := search(query_text, []string{user.GetIndexName()}, size, from)

	r.HTML("index", render.ViewData{
		"logs":       data.Hits.Hits,
		"total":      data.Hits.Total,
		"query_text": query_text,
	})
}
