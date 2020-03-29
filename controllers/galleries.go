package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sajicode/go-photo/context"
	"github.com/sajicode/go-photo/models"
	"github.com/sajicode/go-photo/views"
)

// NewGalleries is used to create a new gallery controller. should only be used at setup
func NewGalleries(gs models.GalleryService) *Galleries {
	return &Galleries{
		New: views.NewView("bootstrap", "galleries/new"),
		gs:  gs,
	}
}

// Galleries struct
type Galleries struct {
	New *views.View
	gs  models.GalleryService
}

// GalleryForm struct
type GalleryForm struct {
	Title string `schema:"title"`
}

// Create a new gallery
// POST /galleries
func (g *Galleries) Create(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form GalleryForm
	if err := parseForm(r, &form); err != nil {
		log.Println(err)
		vd.SetAlert(err)
		g.New.Render(w, vd)
		return
	}
	user := context.User(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
	}
	fmt.Println("Create got the user:", user)
	gallery := models.Gallery{
		Title:  form.Title,
		UserID: user.ID,
	}
	if err := g.gs.Create(&gallery); err != nil {
		log.Println(err)
		vd.SetAlert(err)
		g.New.Render(w, vd)
		return
	}

	fmt.Fprintln(w, gallery)
}
