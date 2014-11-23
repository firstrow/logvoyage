package home

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/firstrow/logvoyage/common"
	"github.com/firstrow/logvoyage/web/context"
)

// View log record
func View(res http.ResponseWriter, ctx *context.Context) {
	conn := common.GetConnection()

	docId := ctx.Request.URL.Query().Get("id")
	docType := ctx.Request.URL.Query().Get("type")

	response, err := conn.Get(ctx.User.GetIndexName(), docType, docId, url.Values{})

	if err != nil {
		res.WriteHeader(404)
	}

	j, err := json.Marshal(response.Source)

	if err != nil {
		res.WriteHeader(503)
	}

	res.Write(j)
}
