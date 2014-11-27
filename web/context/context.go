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
	c.Render.HTML(200, view, data)
}

// Cache all authorized user in memmory
// TODO: Clear cache in lastActivity < now() - 1minute
var userCache = make(map[string]*common.User)

func loadUser(email string) *common.User {
	if email != "" {
		if val, ok := userCache[email]; ok {
			return val
		} else {
			userCache[email] = common.FindUserByEmail(email)
			return userCache[email]
		}
	}
	return nil
}

func Contexter(c martini.Context, r render.Render, sess sessions.Session, req *http.Request) {
	email := sess.Get("email")
	var user *common.User

	if email != nil {
		user = loadUser(sess.Get("email").(string))
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
