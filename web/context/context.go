package context

import (
	"github.com/codegangsta/martini-contrib/render"
	"net/http"

	"github.com/codegangsta/martini"
	"github.com/firstrow/logvoyage/common"
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
	c.Render.HTML(200, view, data)
}

func Contexter(c martini.Context, r render.Render, sess sessions.Session, req *http.Request) {
	var user *common.User

	if sess.Get("email") != nil {
		user = common.FindUserByEmail(sess.Get("email").(string))
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
