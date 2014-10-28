package users

import (
	"errors"
	"github.com/Unknwon/com"
	"github.com/belogik/goes"
	"github.com/codegangsta/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	"net/http"
	"net/url"
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

func findUser(form *loginForm) error {
	conn := goes.NewConnection("localhost", "9200")

	var query = map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": map[string]interface{}{
					"term": map[string]interface{}{
						"email": map[string]interface{}{
							"value": form.Email,
						},
					},
				},
			},
		},
	}

	extraArgs := make(url.Values, 0)
	searchResults, err := conn.Search(query, []string{"users"}, []string{"user"}, extraArgs)

	if err != nil {
		return err
	}

	if searchResults.Hits.Total == 0 {
		return errors.New("User not found")
	}

	hash := com.Sha256(form.Password)
	if searchResults.Hits.Hits[0].Source["password"] != hash {
		return errors.New("Wrong password")
	}

	return nil
}

func Login(req *http.Request, r render.Render, session sessions.Session) {
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
				// TODO: auth, save session
				r.Redirect("/")
			}
		}
	}

	data := map[string]interface{}{
		"form":    form,
		"message": message,
	}
	r.HTML(200, "users/login", data)
}
