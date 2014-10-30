package users

import (
	"errors"
	"github.com/Unknwon/com"
	"github.com/firstrow/logvoyage/common"
	"github.com/firstrow/logvoyage/web/render"
	"github.com/martini-contrib/sessions"
	"net/http"
)

type loginForm struct {
	Email    string
	Password string
	*EnableValidation
}

func (this *loginForm) SetupValidation() {
	this.Valid.Required(this.Email, "Email")
	this.Valid.Email(this.Email, "Email")
	this.Valid.Required(this.Password, "Password")
	this.Valid.MinSize(this.Password, 5, "Password")
	this.Valid.MaxSize(this.Password, 25, "Password")
}

// Search user by login and password
func findUser(form *loginForm) error {
	user := common.FindUserByEmail(form.Email)

	if user == nil {
		return errors.New("User not found")
	}

	hash := com.Sha256(form.Password)
	if user.Password != hash {
		return errors.New("Wrong password")
	}

	return nil
}

func Login(req *http.Request, r *render.Render, sess sessions.Session) {
	message := ""
	req.ParseForm()
	form := &loginForm{
		EnableValidation: &EnableValidation{},
	}

	if req.Method == "POST" {
		form.Email = req.Form.Get("email")
		form.Password = req.Form.Get("password")
		form.SetupValidation()

		if !form.EnableValidation.Valid.HasErrors() {
			// find user
			err := findUser(form)
			if err != nil {
				message = "User not found or wrong password"
			} else {
				sess.Set("email", form.Email)
				r.Redirect("/dashboard")
			}
		}
	}

	r.HTML("users/login", render.ViewData{
		"form":    form,
		"message": message,
	})
}
