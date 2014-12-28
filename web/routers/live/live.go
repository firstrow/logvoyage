package live

import (
	"github.com/firstrow/logvoyage/web/context"
)

func Index(ctx *context.Context) {
	ctx.HTML("live/index", context.ViewData{}, "layouts/simple")
}
