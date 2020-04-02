package controllers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/sajicode/go-photo/context"
	"github.com/sajicode/go-photo/email"
	"github.com/sajicode/go-photo/models"
	"github.com/sajicode/go-photo/rand"
	"github.com/sajicode/go-photo/views"
)

// NewUsers is used to create a new user controller. should only be used at setup
func NewUsers(us models.UserService, emailer email.Client) *Users {
	return &Users{
		NewView:   views.NewView("bootstrap", "users/new"),
		LoginView: views.NewView("bootstrap", "users/login"),
		us:        us,
		emailer:   emailer,
	}
}

// Users struct
type Users struct {
	NewView   *views.View
	LoginView *views.View
	us        models.UserService
	emailer   email.Client
}

// SignupForm struct
type SignupForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

// New function to render user signup
// GET /signup
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	u.NewView.Render(w, r, nil)
}

// Create a new user
// POST /signup
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form SignupForm
	if err := parseForm(r, &form); err != nil {
		log.Println(err)
		vd.SetAlert(err)
		u.NewView.Render(w, r, vd)
		return
	}
	user := models.User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}
	if err := u.us.Create(&user); err != nil {
		vd.SetAlert(err)
		u.NewView.Render(w, r, vd)
		return
	}
	err := u.emailer.Welcome(user.Name, user.Email)

	if err != nil {
		log.Println(err)
	}
	err = u.signIn(w, &user)

	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)

		alert := views.Alert{
			Level:   views.AlertLvlSuccess,
			Message: "Shutters!",
		}
		views.RedirectAlert(w, r, "/galleries", http.StatusFound, alert)

	}
	http.Redirect(w, r, "/galleries", http.StatusFound)
}

// LoginForm struct
type LoginForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

// Login is used to verify a user's email & password
// POST /login
func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	vd := views.Data{}
	form := LoginForm{}
	if err := parseForm(r, &form); err != nil {
		log.Println(err)
		vd.SetAlert(err)
		u.LoginView.Render(w, r, vd)
		return
	}

	user, err := u.us.Authenticate(form.Email, form.Password)

	if err != nil {
		switch err {
		case models.ErrNotFound:
			vd.AlertError("Invalid email address")
		default:
			vd.SetAlert(err)
		}
		u.LoginView.Render(w, r, vd)
		return
	}
	err = u.signIn(w, user)

	if err != nil {
		vd.SetAlert(err)
		u.LoginView.Render(w, r, vd)

		return
	}
	//* we need to set the cookie before printing the user object
	http.Redirect(w, r, "/galleries", http.StatusFound)
}

func (u *Users) signIn(w http.ResponseWriter, user *models.User) error {
	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
		err = u.us.Update(user)
		if err != nil {
			return err
		}
	}
	cookie := http.Cookie{
		Name:     "remember_token",
		Value:    user.Remember,
		HttpOnly: true, //! remember to remove when we want to connect a frontend
	}
	http.SetCookie(w, &cookie)
	return nil
}

// Logout is used to delete a users session cookie (remember_token)
// and then will update the user resource with a new remmeber
// token.
//
// POST /logout
func (u *Users) Logout(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:     "remember_token",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true, //! remove when connecting to client apps
	}
	http.SetCookie(w, &cookie)

	user := context.User(r.Context())
	token, _ := rand.RememberToken()
	user.Remember = token
	u.us.Update(user)
	http.Redirect(w, r, "/", http.StatusFound)
}

// CookieTest is used to display cookies set on a current user
func (u *Users) CookieTest(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("remember_token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	user, err := u.us.ByRemember(cookie.Value)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	fmt.Fprintln(w, user)
}
