package middleware

import (
	"github.com/firstrow/logvoyage/web/render"
	"github.com/martini-contrib/sessions"
)

// Add some defaults to tempalte data each request
func PopulateAppContext(r *render.Render, sess sessions.Session) {
	r.Context["email"] = sess.Get("email")

	if sess.Get("email") != nil {
		r.Context["isGuest"] = false
	}
}

// Check user authentication
func Authorize(r *render.Render, sess sessions.Session) {
	email := sess.Get("email")
	if email == nil {
		r.Redirect("/login")
	}
}

// Redirect user to Dashboard if authorized
func RedirectIfAuthorized(r *render.Render, sess sessions.Session) {
	email := sess.Get("email")
	if email != nil {
		r.Redirect("/dashboard")
	}
}
