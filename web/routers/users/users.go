package users

import (
	"github.com/astaxie/beego/validation"
	"github.com/codegangsta/martini-contrib/render"
	"log"
	"net/http"
)

type registerForm struct {
	Email      string
	Password   string
	Validation validation.Validation
}

func (this *registerForm) GetError(key string) string {
	for _, err := range this.Validation.Errors {
		if err.Key == key {
			return err.Message
		}
	}
	return ""
}

func (this *registerForm) SetupValidation() {
	this.Validation.Required(this.Email, "Email")
	this.Validation.Email(this.Email, "Email")
	this.Validation.Required(this.Password, "Password")
	this.Validation.MinSize(this.Password, 5, "Password")
	this.Validation.MaxSize(this.Password, 25, "Password")
}

func Register(req *http.Request, r render.Render) {
	req.ParseForm()
	form := &registerForm{
		Email:      req.Form.Get("email"),
		Password:   req.Form.Get("password"),
		Validation: validation.Validation{},
	}
	form.SetupValidation()

	log.Println(form.GetError("Email"))

	data := map[string]interface{}{
		"form": form,
	}
	r.HTML(200, "users/register", data)
}
