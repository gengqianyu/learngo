package view

import (
	"crawler/frontend/model"
	"html/template"
	"io"
)

type View struct {
	Template *template.Template
}

func FactoryView(filename string) View {
	return View{Template: template.Must(template.ParseFiles(filename))}
}

// render template 渲染模板 并写入writer
func (v *View) Render(w io.Writer, d model.SearchResult) error {

	return v.Template.Execute(w, d)
}
