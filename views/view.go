package views

import (
	"html/template"
	"path/filepath"
)

var (
	layoutDir   string = "views/layouts/"
	templateExt string = ".gohtml"
)

// NewView handles parsing of our templates
func NewView(layout string, files ...string) *View {
	files = append(files, layoutFiles()...)

	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}

	return &View{
		Template: t,
		Layout:   layout,
	}
}

// View struct for template view
type View struct {
	Template *template.Template
	Layout   string
}

// returns a slice of strings repping files within layout direction
func layoutFiles() []string {
	files, err := filepath.Glob(layoutDir + "*" + templateExt)
	if err != nil {
		panic(err)
	}
	return files
}
