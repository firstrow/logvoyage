package home

import (
	_ "log"
	"net/http"
	"net/url"

	"github.com/firstrow/logvoyage/common"
	"github.com/firstrow/logvoyage/web/render"
)

// View log record by id
func View(req *http.Request, r *render.Render) {
	user := common.FindUserByEmail(r.Context["email"].(string))
	docId := req.URL.Query().Get("id")
	conn := common.GetConnection()

	response, _ := conn.Get(user.GetIndexName(), "logs", docId, url.Values{})

	r.HTML("home/view", render.ViewData{
		"record": response.Source,
	})
}
