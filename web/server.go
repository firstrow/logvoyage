package main

import (
	"github.com/codegangsta/martini-contrib/render"
	"github.com/firstrow/logvoyage/web/routers/home"
	"github.com/go-martini/martini"
	"html/template"
	"log"
	"time"
)

func main() {
	log.Println("Starting server")
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

	m.Get("/", home.Index)
	m.Run()
}
