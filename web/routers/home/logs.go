package home

import (
	"encoding/json"
	"github.com/belogik/goes"
	"net/http"
	"net/url"

	"bitbucket.org/firstrow/logvoyage/common"
	"bitbucket.org/firstrow/logvoyage/web/context"
	"github.com/go-martini/martini"
)

// View log record
func View(res http.ResponseWriter, ctx *context.Context, params martini.Params) {
	conn := common.GetConnection()
	response, err := conn.Get(ctx.User.GetIndexName(), params["type"], params["id"], url.Values{})

	if err != nil {
		res.WriteHeader(404)
	}

	j, err := json.Marshal(response.Source)

	if err != nil {
		res.WriteHeader(503)
	}

	res.Write(j)
}

// Delete log record
func Delete(res http.ResponseWriter, ctx *context.Context, params martini.Params) {
	conn := common.GetConnection()
	d := goes.Document{
		Index: ctx.User.GetIndexName(),
		Type:  params["type"],
		Id:    params["id"],
	}
	_, err := conn.Delete(d, url.Values{})

	if err != nil {
		res.WriteHeader(503)
	}
}
