package profile

import (
	"github.com/firstrow/logvoyage/web/context"
)

func Index(ctx *context.Context) {
	ctx.HTML("profile/index", context.ViewData{})
}
