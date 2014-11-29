package context

import (
	"net/http"

	"github.com/codegangsta/martini-contrib/render"
	"github.com/firstrow/logvoyage/common"
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

func (c *Context) HTML(view string, data ViewData) {
	data["context"] = c
	if c.Request.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		c.Render.HTML(200, view, data, render.HTMLOptions{Layout: ""})
	} else {
		c.Render.HTML(200, view, data)
	}
}

func Contexter(c martini.Context, r render.Render, sess sessions.Session, req *http.Request) {
	email := sess.Get("email")
	var user *common.User

	if email != nil {
		user = common.FindCachedUser(email.(string))
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
