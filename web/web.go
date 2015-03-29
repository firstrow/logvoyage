package main

import (
	"html/template"
	"reflect"
	"runtime"
	"time"

	"github.com/Unknwon/com"
	"github.com/codegangsta/martini-contrib/render"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/sessions"

	"github.com/firstrow/logvoyage/web/context"
	"github.com/firstrow/logvoyage/web/middleware"
	"github.com/firstrow/logvoyage/web/routers/home"
	"github.com/firstrow/logvoyage/web/routers/live"
	"github.com/firstrow/logvoyage/web/routers/profile"
	"github.com/firstrow/logvoyage/web/routers/projects"
	"github.com/firstrow/logvoyage/web/routers/users"
	"github.com/firstrow/logvoyage/web/widgets"
	"github.com/firstrow/logvoyage/web_socket"
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
		"isEmpty": func(i interface{}) bool {
			switch reflect.TypeOf(i).Kind() {
			case reflect.Slice:
				v := reflect.ValueOf(i)
				return v.Len() == 0
			}
			return true
		},
		"eq":                 reflect.DeepEqual,
		"isSliceContainsStr": com.IsSliceContainsStr,
		"buildLogLine":       widgets.BuildLogLine,
	}

	m := martini.Classic()
	// Template
	m.Use(render.Renderer(render.Options{
		Funcs:  []template.FuncMap{templateFunc},
		Layout: "layouts/main",
	}))
	// Serve static files
	m.Use(martini.Static("../static", martini.StaticOptions{
		Prefix:      "static",
		SkipLogging: true,
	}))
	// Sessions
	store := sessions.NewCookieStore([]byte("super_secret_key"))
	m.Use(sessions.Sessions("default", store))

	m.Use(context.Contexter)

	// Routes
	m.Any("/register", middleware.RedirectIfAuthorized, users.Register)
	m.Any("/login", middleware.RedirectIfAuthorized, users.Login)
	m.Any("/logout", middleware.Authorize, users.Logout)
	m.Get("/maintenance", func(ctx *context.Context) {
		ctx.HTML("maintenance", context.ViewData{}, "layouts/simple")
	})
	// Auth routes
	m.Any("/", middleware.Authorize, home.ProjectSearch)
	m.Any("/project/:id", middleware.Authorize, home.ProjectSearch)
	m.Any("/profile", middleware.Authorize, profile.Index)
	m.Any("/live", middleware.Authorize, live.Index)
	// Logs
	m.Get("/log/:id/type/:type", middleware.Authorize, home.View)
	m.Delete("/log/:id/type/:type", middleware.Authorize, home.Delete)
	// Projects
	m.Group("/projects", func(r martini.Router) {
		r.Any("", projects.Index)
		r.Any("/new", projects.New)
		r.Any("/edit/:id", projects.Edit)
		r.Any("/delete/:id", projects.Delete)
		// Types
		r.Any("/types", projects.Types)
		r.Any("/types/delete/:name", projects.DeleteType)
	}, middleware.Authorize)

	go web_socket.StartServer()
	m.Run()
}
