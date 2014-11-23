package profile

import (
	"github.com/firstrow/logvoyage/common"
	"github.com/firstrow/logvoyage/web/context"
)

type profileForm struct {
	*common.EnableValidation
	Email     string
	FirstName string
	LastName  string
}

func (this *profileForm) SetupValidation() {
	this.Valid.Required(this.FirstName, "FirstName")
	this.Valid.Required(this.LastName, "LastName")
	this.Valid.MaxSize(this.FirstName, 25, "FristName")
	this.Valid.MaxSize(this.LastName, 25, "LastName")
}

func Index(ctx *context.Context) {
	ctx.Request.ParseForm()
	form := &profileForm{
		EnableValidation: &common.EnableValidation{},
		FirstName:        ctx.User.FirstName,
		LastName:         ctx.User.LastName,
	}

	if ctx.Request.Method == "POST" {
		form.FirstName = ctx.Request.Form.Get("firstName")
		form.LastName = ctx.Request.Form.Get("lastName")
		form.SetupValidation()

		if !form.EnableValidation.Valid.HasErrors() {
			ctx.User.FirstName = form.FirstName
			ctx.User.LastName = form.LastName
			ctx.User.Save()
		}
	}

	ctx.HTML("profile/index", context.ViewData{
		"form": form,
	})
}
