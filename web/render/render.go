// Package `renderer` makes more comfortable work
// with template rendering
package render

import (
	"github.com/codegangsta/martini-contrib/render"
	"github.com/firstrow/logvoyage/web/context"
	"github.com/go-martini/martini"
)

// Allows to write less code to pass data to tempalte.
// r.HTML("index". render.ViewData{"key": val})
type ViewData map[string]interface{}

type Render struct {
	renderer render.Render
	Context  *context.Context
}

func (r *Render) HTML(name string, bindings map[string]interface{}) {
	bindings["context"] = r.Context
	r.renderer.HTML(200, name, bindings)
}

func (r *Render) Redirect(url string) {
	r.renderer.Redirect(url)
}

// Render middleware
func RenderHandler(c martini.Context, mr render.Render, ctx *context.Context) {
	r := &Render{mr, ctx}
	c.Map(r)
}
