package users

import (
	"github.com/Unknwon/com"
	"github.com/astaxie/beego/validation"
	"github.com/belogik/goes"
	"github.com/codegangsta/martini-contrib/render"
	"net/http"
	"net/url"
)

type EnableValidation struct {
	Valid validation.Validation
}

func (this *EnableValidation) GetError(key string) string {
	for _, err := range this.Valid.Errors {
		if err.Key == key {
			return err.Message
		}
	}
	return ""
}

type registerForm struct {
	Email    string
	Password string
	*EnableValidation
}

func (this *registerForm) SetupValidation() {
	this.Valid.Required(this.Email, "Email")
	this.Valid.Email(this.Email, "Email")
	this.Valid.Required(this.Password, "Password")
	this.Valid.MinSize(this.Password, 5, "Password")
	this.Valid.MaxSize(this.Password, 25, "Password")
}

func Register(req *http.Request, r render.Render) {
	req.ParseForm()
	form := &registerForm{
		EnableValidation: &EnableValidation{},
	}

	if req.Method == "POST" {
		form.Email = req.Form.Get("email")
		form.Password = req.Form.Get("password")
		form.SetupValidation()

		if !form.EnableValidation.Valid.HasErrors() {
			// perform register
			conn := goes.NewConnection("localhost", "9200")
			doc := goes.Document{
				Index: "users",
				Type:  "user",
				Id:    "W_FwNJqpSLaYb-DfT7tG7Q",
				Fields: map[string]string{
					"email":    form.Email,
					"password": com.Sha256(form.Password),
				},
			}
			extraArgs := make(url.Values, 0)
			conn.Index(doc, extraArgs)
		}
	}

	data := map[string]interface{}{
		"form": form,
	}
	r.HTML(200, "users/register", data)
}
