package users

import (
	"fmt"
	"net/url"
	"time"

	"github.com/belogik/goes"
	"github.com/nu7hatch/gouuid"

	c "github.com/firstrow/logvoyage/configuration"
	"github.com/firstrow/logvoyage/common"
	"github.com/firstrow/logvoyage/web/context"
)

type registerForm struct {
	*common.EnableValidation
	Email    string
	Password string
}

func (r *registerForm) IsValid() bool {
	user, _ := common.FindUserByEmail(r.Email)
	if user != nil {
		r.Valid.SetError("Email", "This email is already taken")
		return false
	}
	return true
}

func (r *registerForm) SetupValidation() {
	r.Valid.Required(r.Email, "Email")
	r.Valid.Email(r.Email, "Email")
	r.Valid.Required(r.Password, "Password")
	r.Valid.MinSize(r.Password, 5, "Password")
	r.Valid.MaxSize(r.Password, 25, "Password")
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

		if !form.EnableValidation.Valid.HasErrors() && form.IsValid() {
			password, err := common.HashPassword(form.Password)

			if err != nil {
				panic(err.Error())
			}

			doc := goes.Document{
				Index: c.ReadConf().Indexes.User,
				Type:  "user",
				Fields: map[string]string{
					"email":    form.Email,
					"password": password,
					"apiKey":   buildApiKey(form.Email),
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

func buildApiKey(email string) string {
	t := fmt.Sprintf("%%", email, time.Now().Nanosecond())
	apiKey, _ := uuid.NewV5(uuid.NamespaceURL, []byte(t))
	return apiKey.String()
}
