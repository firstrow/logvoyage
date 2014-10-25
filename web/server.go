package main

import (
	"github.com/codegangsta/martini-contrib/render"
	"github.com/go-martini/martini"
)

func indexPage(r render.Render) {
	data := map[string]interface{}{}
	r.HTML(200, "index", data)
}

func main() {
	m := martini.Classic()
	m.Use(render.Renderer())
	m.Use(martini.Static("../../static"))

	m.Get("/", indexPage)
	m.Run()
}
