package home

import (
	"github.com/belogik/goes"
	_ "log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/firstrow/logvoyage/common"
	"github.com/firstrow/logvoyage/web/render"
	"github.com/firstrow/logvoyage/web/widgets"
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
	pagination := widgets.NewPagination(req)
	pagination.SetPerPage(10)

	// Load records
	data := search(
		query_text,
		[]string{user.GetIndexName()},
		pagination.GetPerPage(),
		pagination.DetectFrom(),
	)

	pagination.SetTotalRecords(data.Hits.Total)

	r.HTML("index", render.ViewData{
		"logs":       data.Hits.Hits,
		"total":      data.Hits.Total,
		"query_text": query_text,
		"pagination": pagination,
	})
}
