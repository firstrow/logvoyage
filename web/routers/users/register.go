package users

import (
	"net/url"

	"github.com/Unknwon/com"
	"github.com/belogik/goes"
	"github.com/firstrow/logvoyage/common"
	"github.com/firstrow/logvoyage/web/context"
	"github.com/nu7hatch/gouuid"
)

type registerForm struct {
	*common.EnableValidation
	Email    string
	Password string
}

func (this *registerForm) SetupValidation() {
	this.Valid.Required(this.Email, "Email")
	this.Valid.Email(this.Email, "Email")
	this.Valid.Required(this.Password, "Password")
	this.Valid.MinSize(this.Password, 5, "Password")
	this.Valid.MaxSize(this.Password, 25, "Password")
}

func Register(ctx *context.Context) {
	ctx.Request.ParseForm()
	form := &registerForm{
		EnableValidation: &common.EnableValidation{},
	}

	if ctx.Request.Method == "POST" {
		form.Email = ctx.Request.Form.Get("email")
		form.Password = ctx.Request.Form.Get("password")
		form.SetupValidation()

		if !form.EnableValidation.Valid.HasErrors() {
			apiKey, _ := uuid.NewV5(uuid.NamespaceURL, []byte(form.Email))

			doc := goes.Document{
				Index: "users",
				Type:  "user",
				Fields: map[string]string{
					"email":    form.Email,
					"password": com.Sha256(form.Password),
					"apiKey":   apiKey.String(),
				},
			}
			extraArgs := make(url.Values, 0)
			common.GetConnection().Index(doc, extraArgs)
			ctx.Render.Redirect("/login")
		}
	}

	ctx.HTML("users/register", context.ViewData{
		"form": form,
	})
}
