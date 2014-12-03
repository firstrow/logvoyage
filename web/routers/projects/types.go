package projects

import (
	"github.com/firstrow/logvoyage/common"
	"github.com/firstrow/logvoyage/web/context"
	"github.com/go-martini/martini"
)

// Display list of ES types available to user.
func Types(ctx *context.Context) {
	ctx.HTML("projects/types", context.ViewData{
		"docCounter": common.CountTypeDocs,
	})
}

func DeleteType(ctx *context.Context, params martini.Params) {
	common.DeleteType(ctx.User.GetIndexName(), params["name"])
	ctx.Render.Redirect("/projects/types")
}
