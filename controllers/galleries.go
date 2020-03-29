package controllers

import (
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
