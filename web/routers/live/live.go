package live

import (
	"bitbucket.org/firstrow/logvoyage/web/context"
)

func Index(ctx *context.Context) {
	ctx.HTML("live/index", context.ViewData{}, "layouts/simple")
}
