package projects

import (
	"github.com/firstrow/logvoyage/common"
	"github.com/firstrow/logvoyage/web/context"
	"github.com/Unknwon/com"
	"github.com/go-martini/martini"
)

type projectForm struct {
	*common.EnableValidation
	Name        string
	Description string
	Types       []string
	Id          string
}

func (s *projectForm) HasType(typeName string) bool {
	return com.IsSliceContainsStr(s.Types, typeName)
}

func (s *projectForm) SetupValidation() {
	s.Valid.Required(s.Name, "Name")
	s.Valid.MaxSize(s.Name, 25, "Name")
	s.Valid.MaxSize(s.Description, 250, "Description")
}

func Index(ctx *context.Context) {
	ctx.HTML("projects/index", context.ViewData{})
}

func New(ctx *context.Context) {
	form := buildForm(ctx)
	update(ctx, form)
}

func Edit(ctx *context.Context, params martini.Params) {
	form := buildForm(ctx)
	group, err := ctx.User.GetProject(params["id"])

	if err != nil {
		ctx.Render.Error(404)
	}

	form.Id = group.Id
	form.Name = group.Name
	form.Description = group.Description
	form.Types = group.Types
	update(ctx, form)
}

func Delete(ctx *context.Context, params martini.Params) {
	ctx.User.DeleteProject(params["id"])
	ctx.User.Save()
	ctx.Session.AddFlash("Project has been successfully deleted.", "success")
	ctx.Render.Redirect("/projects")
}

func buildForm(ctx *context.Context) *projectForm {
	ctx.Request.ParseForm()
	form := &projectForm{
		EnableValidation: &common.EnableValidation{},
	}
	return form
}

func update(ctx *context.Context, form *projectForm) {
	if ctx.Request.Method == "POST" {
		form.Name = ctx.Request.Form.Get("name")
		form.Description = ctx.Request.Form.Get("description")
		form.Types = ctx.Request.PostForm["types"]
		form.SetupValidation()

		if !form.EnableValidation.Valid.HasErrors() {
			group := &common.Project{
				Id:          form.Id,
				Name:        form.Name,
				Description: form.Description,
				Types:       form.Types,
			}
			ctx.User.AddProject(group).Save()
			ctx.Session.AddFlash("Project has been successfully saved.", "success")
			ctx.Render.Redirect("/projects")
		}
	}

	ctx.HTML("projects/new", context.ViewData{
		"form": form,
	})
}
