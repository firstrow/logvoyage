package live

import (
	"../../../web/context"
)

func Index(ctx *context.Context) {
	ctx.HTML("live/index", context.ViewData{}, "layouts/simple")
}
