package main

import (
	"github.com/codegangsta/martini-contrib/render"
	"github.com/go-martini/martini"
	"html/template"
	"runtime"
	"time"

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
	m.Use(render.Renderer(render.Options{
		Funcs: []template.FuncMap{templateFunc},
	}))
	m.Use(martini.Static("../static"))

	// Routes
	m.Get("/", home.Index)
	m.Any("/register", users.Register)
	m.Any("/login", users.Login)

	m.Run()
}
