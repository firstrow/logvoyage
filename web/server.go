package main

import (
	martiniRender "github.com/codegangsta/martini-contrib/render"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/sessions"
	"html/template"
	"runtime"
	"time"

	"github.com/firstrow/logvoyage/web/render"
	"github.com/firstrow/logvoyage/web/routers/home"
	"github.com/firstrow/logvoyage/web/routers/users"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	templateFunc := template.FuncMap{
		"FormatTimeToHuman": func(s string) string {
			t, _ := time.Parse(time.RFC3339Nano, s)
			return t.Format("2006-01-02 15:04:05") + " UTC"
		},
	}

	m := martini.Classic()
	// Template
	m.Use(martiniRender.Renderer(martiniRender.Options{
		Funcs: []template.FuncMap{templateFunc},
	}))
	// Serve static files
	m.Use(martini.Static("../static"))
	// Sessions
	store := sessions.NewCookieStore([]byte("super_secret_key"))
	m.Use(sessions.Sessions("default", store))

	// Application renderer
	m.Use(render.RenderHandler)
	m.Use(templateContext)

	// Routes
	m.Get("/dashboard", authorize, home.Index)
	m.Any("/register", users.Register)
	m.Any("/login", users.Login)

	m.Run()
}

func templateContext(r *render.Render, sess sessions.Session) {
	r.Context["email"] = sess.Get("email")
}

func authorize(r *render.Render, sess sessions.Session) {
	email := sess.Get("email")
	if email == nil {
		r.Redirect("/login")
	}
}
