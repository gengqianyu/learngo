package view

import (
	"html/template"
	"io"
)

type View struct {
	Template *template.Template
}

func CreateView(templateName string) View {
	templateName = "src/command/view/" + templateName + ".html"
	return View{
		Template: template.Must(template.ParseFiles(templateName)),
	}
}

func (v View) Render(w io.Writer, data interface{}) error {
	return v.Template.Execute(w, data)
}
