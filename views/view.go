package views

import "html/template"

// NewView handles parsing of our templates
func NewView(layout string, files ...string) *View {
	files = append(files, "views/layouts/bootstrap.gohtml", "views/layouts/navbar.gohtml", "views/layouts/footer.gohtml")

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
