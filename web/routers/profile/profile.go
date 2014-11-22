package profile

import (
	"github.com/firstrow/logvoyage/common"
	"github.com/firstrow/logvoyage/web/render"
	"net/http"
)

func Index(req *http.Request, r *render.Render) {
	user := common.FindUserByEmail(r.Context["email"].(string))
	userLogTypes, _ := common.GetTypes(user.GetIndexName())

	r.HTML("profile/index", render.ViewData{
		"userLogTypes": userLogTypes,
	})
}
