package context

import (
	"net/http"

	"github.com/codegangsta/martini-contrib/render"
	"bitbucket.org/firstrow/logvoyage/common"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/sessions"
)

type ViewData map[string]interface{}

type Context struct {
	Session sessions.Session
	User    *common.User
	Request *http.Request
	Render  render.Render
	IsGuest bool
}

func (c *Context) HTML(view string, data ViewData, layout ...string) {
	data["context"] = c
	if c.Request.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// Disable layout for ajax requests
		c.Render.HTML(200, view, data, render.HTMLOptions{Layout: ""})
	} else {
		var l string
		if len(layout) > 0 {
			l = layout[0]
		} else {
			l = "layouts/main"
		}
		c.Render.HTML(200, view, data, render.HTMLOptions{Layout: l})
	}
}

func Contexter(c martini.Context, r render.Render, sess sessions.Session, req *http.Request) {
	email := sess.Get("email")
	var user *common.User

	if email != nil {
		user, _ = common.FindCachedUser(email.(string))
	} else {
		user = nil
	}

	ctx := &Context{
		Session: sess,
		IsGuest: user == nil,
		User:    user,
		Request: req,
		Render:  r,
	}
	c.Map(ctx)
}
