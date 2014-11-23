package profile

import (
	"github.com/firstrow/logvoyage/web/context"
	"github.com/firstrow/logvoyage/web/render"
)

func Index(r *render.Render, ctx *context.Context) {
	r.HTML("profile/index", render.ViewData{})
}
