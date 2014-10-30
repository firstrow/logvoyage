package main

import (
	"github.com/firstrow/logvoyage/web/render"
	"github.com/martini-contrib/sessions"
)

// Add some defaults to tempalte data each request
func templateContext(r *render.Render, sess sessions.Session) {
	r.Context["email"] = sess.Get("email")

	if sess.Get("email") != nil {
		r.Context["isGuest"] = false
	}
}

// Check user authentication
func authorize(r *render.Render, sess sessions.Session) {
	email := sess.Get("email")
	if email == nil {
		r.Redirect("/login")
	}
}

// Redirect user to Dashboard if authorized
func redirectIfAuthorized(r *render.Render, sess sessions.Session) {
	email := sess.Get("email")
	if email != nil {
		r.Redirect("/dashboard")
	}
}
