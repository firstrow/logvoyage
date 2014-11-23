package context

import (
	"net/http"

	"github.com/codegangsta/martini"
	"github.com/firstrow/logvoyage/common"
	"github.com/martini-contrib/sessions"
)

type Context struct {
	Session sessions.Session
	User    *common.User
	Request *http.Request
	IsGuest bool
}

func Contexter(c martini.Context, sess sessions.Session, req *http.Request) {
	var user *common.User

	if sess.Get("email") != nil {
		user = common.FindUserByEmail(sess.Get("email").(string))
	}

	ctx := &Context{
		Session: sess,
		IsGuest: user == nil,
		User:    user,
		Request: req,
	}
	c.Map(ctx)
}
