package users

import (
	"errors"

	"github.com/Unknwon/com"
	"github.com/firstrow/logvoyage/common"
	"github.com/firstrow/logvoyage/web/context"
)

type loginForm struct {
	*common.EnableValidation
	Email    string
	Password string
}

func (this *loginForm) SetupValidation() {
	this.Valid.Required(this.Email, "Email")
	this.Valid.Email(this.Email, "Email")
	this.Valid.Required(this.Password, "Password")
	this.Valid.MinSize(this.Password, 5, "Password")
	this.Valid.MaxSize(this.Password, 25, "Password")
}

// Search user by login and password
func userExists(form *loginForm) error {
	user, _ := common.FindUserByEmail(form.Email)

	if user == nil {
		return errors.New("User not found")
	}

	hash := com.Sha256(form.Password)
	if user.Password != hash {
		return errors.New("Wrong password")
	}

	return nil
}

func Login(ctx *context.Context) {
	message := ""
	ctx.Request.ParseForm()
	form := &loginForm{
		EnableValidation: &common.EnableValidation{},
	}

	if ctx.Request.Method == "POST" {
		form.Email = ctx.Request.Form.Get("email")
		form.Password = ctx.Request.Form.Get("password")
		form.SetupValidation()

		if !form.EnableValidation.Valid.HasErrors() {
			// find user
			err := userExists(form)
			if err != nil {
				message = "User not found or wrong password"
			} else {
				ctx.Session.Set("email", form.Email)
				ctx.Render.Redirect("/")
			}
		}
	}

	ctx.HTML("users/login", context.ViewData{
		"form":    form,
		"message": message,
	})
}
