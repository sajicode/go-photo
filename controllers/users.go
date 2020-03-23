package controllers

import (
	"fmt"
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

// SignupForm struct
type SignupForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

// New function to render user signup
// GET /signup
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	if err := u.NewView.Render(w, nil); err != nil {
		panic(err)
	}
}

// Create a new user
// POST /signup
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	var form SignupForm
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}
	fmt.Fprintln(w, form)
}
