package controllers

import (
	"net/http"

	"github.com/sajicode/go-photo/views"
)

// NewUsers is used to create a new user controller. should only be used at setup
func NewUsers() *Users {
	return &Users{
		NewView: views.NewView("bootstrap", "views/users/new.gohtml"),
	}
}

// Users struct
type Users struct {
	NewView *views.View
}

// New function to render user signup
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	if err := u.NewView.Render(w, nil); err != nil {
		panic(err)
	}
}
