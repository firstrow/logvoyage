package home

import (
	"errors"
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
func search(text string, indexes []string, size int, from int) (goes.Response, error) {
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

	// query["filter"] = map[string]interface{}{
	// 	"range": map[string]interface{}{
	// 		"datetime": map[string]string{
	// 			"gt": "2014-11-02T12:08:52",
	// 		},
	// 	},
	// }

	extraArgs := make(url.Values, 1)
	searchResults, err := conn.Search(query, indexes, []string{"logs"}, extraArgs)

	if err != nil {
		return goes.Response{}, errors.New("No records found.")
	} else {
		return searchResults, nil
	}
}

type timeRange struct {
	gt string
	lt string
}

// Detects time range from request and returns
// elastic compatible format string
func buildTimeRange(req *http.Request) {
	// 2014-11-02T12:08:52
	// return gt and lt?
	// test

	// work flow:
	// get time value
	// use switch to build time
	// if time is custom
	//    get start and end times
	// return timeRange
}

func Index(req *http.Request, r *render.Render) {
	query_text := req.URL.Query().Get("q")
	user := common.FindUserByEmail(r.Context["email"].(string))

	// Pagination
	pagination := widgets.NewPagination(req)
	pagination.SetPerPage(100)

	// Load records
	data, err := search(
		query_text,
		[]string{user.GetIndexName()},
		pagination.GetPerPage(),
		pagination.DetectFrom(),
	)

	if err != nil {
		r.HTML("home/no_records", render.ViewData{})
	} else {
		pagination.SetTotalRecords(data.Hits.Total)

		r.HTML("home/index", render.ViewData{
			"logs":       data.Hits.Hits,
			"total":      data.Hits.Total,
			"query_text": query_text,
			"pagination": pagination,
		})
	}
}
