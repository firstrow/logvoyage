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

type DateTimeRange struct {
	Start string
	Stop  string
}

// Represents search request to perform in ES
type SearchRequest struct {
	Text      string   // test to search
	Indexes   []string // ES indexeses to perform search
	Types     []string // search types
	Size      int      // home much objects ES must return
	From      int      // how much objects should ES skip from first
	TimeRange DateTimeRange
}

func buildSearchRequest(text string, indexes []string, size int, from int, datetime DateTimeRange) SearchRequest {
	req := SearchRequest{
		Text:      text,
		Indexes:   indexes,
		From:      from,
		Types:     []string{"logs"},
		Size:      10,
		TimeRange: DateTimeRange{},
	}
	return req
}

// Search logs in elastic.
func search(searchRequest SearchRequest) (goes.Response, error) {
	conn := common.GetConnection()

	if len(searchRequest.Text) > 0 {
		strconv.Quote(searchRequest.Text)
	} else {
		searchRequest.Text = "*"
	}

	var query = map[string]interface{}{
		"query": map[string]interface{}{
			"query_string": map[string]string{
				"default_field": "message",
				"query":         searchRequest.Text,
			},
		},
		"from": searchRequest.From,
		"size": searchRequest.Size,
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
	searchResults, err := conn.Search(query, searchRequest.Indexes, searchRequest.Types, extraArgs)

	if err != nil {
		return goes.Response{}, errors.New("No records found.")
	} else {
		return searchResults, nil
	}
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
	searchRequest := buildSearchRequest(
		query_text,
		[]string{user.GetIndexName()},
		pagination.GetPerPage(),
		pagination.DetectFrom(),
		DateTimeRange{},
	)
	data, err := search(searchRequest)

	pagination.SetTotalRecords(data.Hits.Total)

	var viewName string
	if data.Hits.Total > 0 && err == nil {
		viewName = "home/index"
	} else {
		viewName = "home/no_records"
	}

	r.HTML(viewName, render.ViewData{
		"logs":       data.Hits.Hits,
		"total":      data.Hits.Total,
		"query_text": query_text,
		"pagination": pagination,
	})
}
