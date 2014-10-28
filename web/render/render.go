// Package `renderer` makes more comfortable work
// with tempalte rendering
package render

import (
	"github.com/codegangsta/martini-contrib/render"
	"github.com/go-martini/martini"
)

type ViewData map[string]interface{}

type Render struct {
	renderer render.Render
}

func (r *Render) HTML(name string, bindings map[string]interface{}) {
	// bindings["key"] = "val"
	r.renderer.HTML(200, name, bindings)
}

func (r *Render) Redirect(url string) {
	r.renderer.Redirect(url)
}

func RenderHandler(c martini.Context, mr render.Render) {
	r := &Render{mr}
	c.Map(r)
}
