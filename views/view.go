package views

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"path/filepath"

	"github.com/gorilla/csrf"
	"github.com/sajicode/go-photo/context"
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

	t, err := template.New("").Funcs(
		template.FuncMap{
			"csrfField": func() (template.HTML, error) {
				return "", errors.New("csrfField is not implemented")
			},
		},
	).ParseFiles(files...)
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
	v.Render(w, r, nil)
}

// Render function helps with rendering view with pre-defined layout
func (v *View) Render(w http.ResponseWriter, r *http.Request, data interface{}) {
	w.Header().Set("Content-Type", "text/html")
	var vd Data
	switch d := data.(type) {
	case Data:
		vd = d
		// do nothing
	default:
		vd = Data{
			Yield: data,
		}
	}
	if alert := getAlert(r); alert != nil {
		vd.Alert = alert
		clearAlert(w)
	}

	// set user from context on view data
	vd.User = context.User(r.Context())
	var buf bytes.Buffer
	csrfField := csrf.TemplateField(r)
	tpl := v.Template.Funcs(template.FuncMap{
		"csrfField": func() template.HTML {
			return csrfField
		},
	})
	if err := tpl.ExecuteTemplate(&buf, v.Layout, vd); err != nil {
		fmt.Println(err)
		fmt.Fprintln(w, "Something went wrong. If the problem persists, please send us an email", http.StatusInternalServerError)
		return
	}
	io.Copy(w, &buf)
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
