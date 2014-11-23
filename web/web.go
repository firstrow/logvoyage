package main

import (
	"html/template"
	"reflect"
	"runtime"
	"time"

	"github.com/Unknwon/com"
	"github.com/codegangsta/martini-contrib/render"
	"github.com/firstrow/logvoyage/web/context"
	"github.com/firstrow/logvoyage/web/middleware"
	"github.com/firstrow/logvoyage/web/routers/home"
	"github.com/firstrow/logvoyage/web/routers/profile"
	"github.com/firstrow/logvoyage/web/routers/users"
	"github.com/firstrow/logvoyage/web_socket"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/sessions"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Template methods
	templateFunc := template.FuncMap{
		"FormatTimeToHuman": func(s ...string) string {
			if len(s) > 0 {
				t, _ := time.Parse(time.RFC3339Nano, s[0])
				return t.Format("2006-01-02 15:04:05") + " UTC"
			} else {
				return "Unknown"
			}
		},
		"eq":                 reflect.DeepEqual,
		"isSliceContainsStr": com.IsSliceContainsStr,
	}

	m := martini.Classic()
	// Template
	m.Use(render.Renderer(render.Options{
		Funcs:  []template.FuncMap{templateFunc},
		Layout: "layouts/main",
	}))
	// Serve static files
	m.Use(martini.Static("../static"))
	// Sessions
	store := sessions.NewCookieStore([]byte("super_secret_key"))
	m.Use(sessions.Sessions("default", store))

	m.Use(context.Contexter)

	// Routes
	m.Any("/register", middleware.RedirectIfAuthorized, users.Register)
	m.Any("/login", middleware.RedirectIfAuthorized, users.Login)
	// Auth routes
	m.Get("/dashboard", middleware.Authorize, home.Index)
	m.Get("/view", middleware.Authorize, home.View)
	m.Get("/profile", middleware.Authorize, profile.Index)

	go web_socket.StartServer()
	m.Run()
}
