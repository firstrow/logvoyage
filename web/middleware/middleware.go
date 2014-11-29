package middleware

import (
	"github.com/codegangsta/martini-contrib/render"
	"github.com/firstrow/logvoyage/common"
	"github.com/martini-contrib/sessions"
)

// Check user authentication
func Authorize(r render.Render, sess sessions.Session) {
	email := sess.Get("email")
	if email != nil {
		user := common.FindCachedUser(email.(string))
		if user == nil {
			r.Redirect("/login")
		}
	}
}

// Redirect user to Dashboard if authorized
func RedirectIfAuthorized(r render.Render, sess sessions.Session) {
	email := sess.Get("email")
	if email != nil {
		r.Redirect("/dashboard")
	}
}
