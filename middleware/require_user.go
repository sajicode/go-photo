package middleware

import (
	"net/http"

	"github.com/sajicode/go-photo/context"
	"github.com/sajicode/go-photo/models"
)

// User struct
type User struct {
	models.UserService
}

// Apply middleware takes http handler as arg and returns ApplyFn function
func (u *User) Apply(next http.Handler) http.HandlerFunc {
	return u.ApplyFn(next.ServeHTTP)
}

// ApplyFn middleware to controller
func (u *User) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("remember_token")
		if err != nil {
			next(w, r)
			return
		}
		user, err := u.UserService.ByRemember(cookie.Value)
		if err != nil {
			next(w, r)
			return
		}
		ctx := r.Context()
		ctx = context.WithUser(ctx, user)
		r = r.WithContext(ctx)
		next(w, r)
	})
}

// RequireUser struct holds the fields required
type RequireUser struct {
	User
}

// Apply assumes that User middleware has already been run,
// otherwise it will not work correctly
func (mw *RequireUser) Apply(next http.Handler) http.HandlerFunc {
	return mw.ApplyFn(next.ServeHTTP)
}

// ApplyFn assumes that User middleware has already been yun
func (mw *RequireUser) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := context.User(r.Context())
		if user == nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		next(w, r)
	})
}
