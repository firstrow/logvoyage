package main

import (
	"github.com/codegangsta/martini-contrib/render"
	"github.com/go-martini/martini"
)

func main() {
	m := martini.Classic()
	m.Use(render.Renderer())

	m.Get("/", func(r render.Render) {
		data := map[string]interface{}{"name": "tester"}
		r.HTML(200, "index", data)
	})
	m.Run()
}
