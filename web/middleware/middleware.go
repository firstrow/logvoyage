package middleware

import (
	"github.com/codegangsta/martini-contrib/render"
	"bitbucket.org/firstrow/logvoyage/common"
	"github.com/martini-contrib/sessions"
)

// Check user authentication
func Authorize(r render.Render, sess sessions.Session) {
	email := sess.Get("email")
	if email != nil {
		user, err := common.FindCachedUser(email.(string))
		if err == common.ErrSendingElasticSearchRequest {
			r.Redirect("/maintenance")
			return
		}
		if user == nil {
			r.Redirect("/login")
			return
		}
	} else {
		r.Redirect("/login")
	}
}

// Redirect user to Dashboard if authorized
func RedirectIfAuthorized(r render.Render, sess sessions.Session) {
	email := sess.Get("email")
	if email != nil {
		r.Redirect("/")
	}
}
