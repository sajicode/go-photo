package views

import (
	"html/template"
	"net/http"
	"path/filepath"
)

var (
	layoutDir   string = "views/layouts/"
	templateDir string = "views/"
	templateExt string = ".gohtml"
)

// NewView handles parsing of our templates
func NewView(layout string, files ...string) *View {
	addTemplatePath(files)
	addTemplateExt(files)
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

// ServeHTTP function
func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := v.Render(w, nil); err != nil {
		panic(err)
	}
}

// Render function helps with rendering view with pre-defined layout
func (v *View) Render(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "text/html")
	switch data.(type) {
	case Data:
		// do nothing
	default:
		data = Data{
			Yield: data,
		}
	}
	return v.Template.ExecuteTemplate(w, v.Layout, data)
}

// returns a slice of strings repping files within layout direction
func layoutFiles() []string {
	files, err := filepath.Glob(layoutDir + "*" + templateExt)
	if err != nil {
		panic(err)
	}
	return files
}

// addTemplatePath takes in a slice of strings representing file paths for templates
// and it prepends the TemplateDir directory to each strinf in the slice'
// eg {"home"} ==> {"views/home"}
func addTemplatePath(files []string) {
	for i, f := range files {
		files[i] = templateDir + f
	}
}

// addTemplateExt takes in a slice of strings representing file paths for templates and
// it appends the template ext to each string in the slice e.g {home} ==> {home.gohtml}
func addTemplateExt(files []string) {
	for i, f := range files {
		files[i] = f + templateExt
	}
}
