package home

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/belogik/goes"
	"github.com/go-martini/martini"

	"github.com/firstrow/logvoyage/common"
	"github.com/firstrow/logvoyage/web/context"
	"github.com/firstrow/logvoyage/web/widgets"
)

const (
	timeLayout = "2006/01/02 15:04" // Users input time format
	perPage    = 100
)

type DateTimeRange struct {
	Start string
	Stop  string
}

func (this *DateTimeRange) IsValid() bool {
	return this.Start != "" || this.Stop != ""
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

func buildSearchRequest(text string, indexes []string, types []string, size int, from int, datetime DateTimeRange) SearchRequest {
	return SearchRequest{
		Text:      text,
		Indexes:   indexes,
		From:      from,
		Types:     types,
		Size:      perPage,
		TimeRange: datetime,
	}
}

// Detects time range from request and returns
// elastic compatible format string
func buildTimeRange(req *http.Request) DateTimeRange {
	var timeRange DateTimeRange

	switch req.URL.Query().Get("time") {
	case "15m":
		timeRange.Start = "now-15m"
	case "30m":
		timeRange.Start = "now-30m"
	case "60m":
		timeRange.Start = "now-60m"
	case "12h":
		timeRange.Start = "now-12h"
	case "24h":
		timeRange.Start = "now-24h"
	case "week":
		timeRange.Start = "now-7d"
	case "custom":
		timeStart, err := time.Parse(timeLayout, req.URL.Query().Get("time_start"))
		if err == nil {
			timeRange.Start = timeStart.Format(time.RFC3339)
		}
		timeStop, err := time.Parse(timeLayout, req.URL.Query().Get("time_stop"))
		if err == nil {
			timeRange.Stop = timeStop.Format(time.RFC3339)
		}
	}

	return timeRange
}

// Search logs in elastic.
func search(searchRequest SearchRequest) (goes.Response, error) {
	conn := common.GetConnection()

	var query = map[string]interface{}{
		"from": searchRequest.From,
		"size": searchRequest.Size,
		"sort": map[string]string{
			"datetime": "desc",
		},
	}

	if len(searchRequest.Text) > 0 {
		strconv.Quote(searchRequest.Text)
		query["query"] = map[string]interface{}{
			"query_string": map[string]string{
				"default_field": "message",
				"query":         searchRequest.Text,
			},
		}
	}

	// Build time range query
	if searchRequest.TimeRange.IsValid() {
		datetime := make(map[string]string)
		if searchRequest.TimeRange.Start != "" {
			datetime["gte"] = searchRequest.TimeRange.Start
		}
		if searchRequest.TimeRange.Stop != "" {
			datetime["lte"] = searchRequest.TimeRange.Stop
		}
		query["filter"] = map[string]interface{}{
			"range": map[string]interface{}{
				"datetime": datetime,
			},
		}
	}

	extraArgs := make(url.Values, 0)
	searchResults, err := conn.Search(query, searchRequest.Indexes, searchRequest.Types, extraArgs)

	if err != nil {
		return goes.Response{}, errors.New("No records found.")
	} else {
		return *searchResults, nil
	}
}

// This function handles two routes "/" and "/project/:id"
func ProjectSearch(ctx *context.Context, params martini.Params) {
	var types []string
	var project *common.Project

	query_text := ctx.Request.URL.Query().Get("q")
	selected_types := ctx.Request.URL.Query()["types"]

	// Project scope
	if _, err := params["id"]; err {
		project, err := ctx.User.GetProject(params["id"])
		if err != nil {
			ctx.HTML("shared/error", context.ViewData{
				"message": "Project not found",
			})
			return
		}
		if len(project.Types) == 0 {
			ctx.HTML("home/empty_project", context.ViewData{
				"project": project,
			})
			return
		}
		if len(selected_types) > 0 {
			types = selected_types
		} else {
			types = project.Types
		}
	}

	// Pagination
	pagination := widgets.NewPagination(ctx.Request)
	pagination.SetPerPage(perPage)

	// Load records
	searchRequest := buildSearchRequest(
		query_text,
		[]string{ctx.User.GetIndexName()},
		types,
		pagination.GetPerPage(),
		pagination.DetectFrom(),
		buildTimeRange(ctx.Request),
	)
	// Search data in elastic
	data, _ := search(searchRequest)

	pagination.SetTotalRecords(data.Hits.Total)

	var viewName string
	viewData := context.ViewData{
		"project":    project,
		"logs":       data.Hits.Hits,
		"total":      data.Hits.Total,
		"took":       data.Took,
		"types":      types,
		"time":       ctx.Request.URL.Query().Get("time"),
		"time_start": ctx.Request.URL.Query().Get("time_start"),
		"time_stop":  ctx.Request.URL.Query().Get("time_stop"),
		"query_text": query_text,
		"pagination": pagination,
	}

	if ctx.Request.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		viewName = "home/table"
	} else {
		viewName = "home/index"
	}

	ctx.HTML(viewName, viewData)
}
