package views

import "html/template"

// NewView handles parsing of our templates
func NewView(files ...string) *View {
	files = append(files, "views/layouts/footer.gohtml")

	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}

	return &View{
		Template: t,
	}
}

// View struct for template view
type View struct {
	Template *template.Template
}
