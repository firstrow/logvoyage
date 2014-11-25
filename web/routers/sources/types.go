package sources

import (
	"github.com/firstrow/logvoyage/common"
	"github.com/firstrow/logvoyage/web/context"
	"github.com/go-martini/martini"
)

func Types(ctx *context.Context) {
	ctx.HTML("sources/types", context.ViewData{
		"docCounter": common.CountTypeDocs,
	})
}

func DeleteType(ctx *context.Context, params martini.Params) {
	common.DeleteType(ctx.User.GetIndexName(), params["name"])
	ctx.Render.Redirect("/sources/types")
}
