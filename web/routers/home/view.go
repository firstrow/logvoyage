package home

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/firstrow/logvoyage/common"
	"github.com/firstrow/logvoyage/web/render"
)

// View log record
func View(req *http.Request, res http.ResponseWriter, r *render.Render) {
	user := common.FindUserByEmail(r.Context["email"].(string))
	conn := common.GetConnection()

	docId := req.URL.Query().Get("id")
	docType := req.URL.Query().Get("type")

	response, err := conn.Get(user.GetIndexName(), docType, docId, url.Values{})

	if err != nil {
		res.WriteHeader(404)
	}

	j, err := json.Marshal(response.Source)

	if err != nil {
		res.WriteHeader(503)
	}

	res.Write(j)
}
