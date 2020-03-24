package controllers

import "github.com/sajicode/go-photo/views"

// NewStatic function that helps render static pages
func NewStatic() *Static {
	return &Static{
		Home:    views.NewView("bootstrap", "static/home"),
		Contact: views.NewView("bootstrap", "static/contact"),
	}
}

// Static struct
type Static struct {
	Home    *views.View
	Contact *views.View
}
