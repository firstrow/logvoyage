package home

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/firstrow/logvoyage/common"
	"github.com/firstrow/logvoyage/web/render"
)

// View log record
func View(req *http.Request, r *render.Render) {
	user := common.FindUserByEmail(r.Context["email"].(string))
	conn := common.GetConnection()

	docId := req.URL.Query().Get("id")
	docType := req.URL.Query().Get("type")

	response, err := conn.Get(user.GetIndexName(), docType, docId, url.Values{})

	if err != nil {
		r.HTML("home/message", render.ViewData{
			"message": "Record not found",
		})
	}

	j, _ := json.Marshal(response.Source)

	if err != nil {
		r.HTML("home/no_records", render.ViewData{
			"message": "Error encoding json",
		})
	}

	r.HTML("home/view", render.ViewData{
		"record": string(j),
	})
}
