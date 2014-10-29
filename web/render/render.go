// Package `renderer` makes more comfortable work
// with template rendering
package render

import (
	"github.com/codegangsta/martini-contrib/render"
	"github.com/go-martini/martini"
)

// Allows to write less code to pass data to tempalte.
// r.HTML("index". render.ViewData{"key": val})
type ViewData map[string]interface{}

type Render struct {
	renderer render.Render
	// Allows to add data before render
	// Used to add `global` values
	Context ViewData
}

func (r *Render) HTML(name string, bindings map[string]interface{}) {
	// Merge `global` data and local
	for key, val := range r.Context {
		bindings[key] = val
	}
	r.renderer.HTML(200, name, bindings)
}

func (r *Render) Redirect(url string) {
	r.renderer.Redirect(url)
}

// Render middleware
func RenderHandler(c martini.Context, mr render.Render) {
	r := &Render{mr, ViewData{}}
	c.Map(r)
}
