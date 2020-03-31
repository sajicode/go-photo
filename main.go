package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sajicode/go-photo/controllers"
	"github.com/sajicode/go-photo/middleware"
	"github.com/sajicode/go-photo/models"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "kamikaze"
	dbname   = "gophotos_dev"
	dbDriver = "postgres"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	services, err := models.NewServices(dbDriver, psqlInfo)
	must(err)
	defer services.Close()
	//! to clear db
	// services.DestructiveReset()

	services.AutoMigrate()

	r := mux.NewRouter()
	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(services.User)
	galleriesC := controllers.NewGalleries(services.Gallery, services.Image, r)
	userMw := middleware.User{
		UserService: services.User,
	}

	// user middleware
	requireUserMw := middleware.RequireUser{
		User: userMw,
	}

	r.NotFoundHandler = http.HandlerFunc(notFound)
	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")

	// User routes
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	r.HandleFunc("/cookietest", usersC.CookieTest).Methods("GET")
	r.Handle("/login", usersC.LoginView).Methods("GET")
	r.HandleFunc("/login", usersC.Login).Methods("POST")

	r.HandleFunc("/faq", faq).Methods("GET")

	// Image routes
	// * so far a route has an image prefix, run the accompanying function
	imageHandler := http.FileServer(http.Dir("./images/"))
	r.PathPrefix("/images/").Handler(http.StripPrefix("/images/", imageHandler))

	// * named routes are useful for when we want to redirect to a particular route after an action
	// Gallery routes
	r.Handle("/galleries", requireUserMw.ApplyFn(galleriesC.Index)).Methods("GET")
	r.Handle("/galleries/new", requireUserMw.Apply(galleriesC.New)).Methods("GET")
	r.HandleFunc("/galleries", requireUserMw.ApplyFn(galleriesC.Create)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/edit", requireUserMw.ApplyFn(galleriesC.Edit)).Methods("GET").Name(controllers.EditGallery)
	r.HandleFunc("/galleries/{id:[0-9]+}/update", requireUserMw.ApplyFn(galleriesC.Update)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/delete", requireUserMw.ApplyFn(galleriesC.Delete)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/images", requireUserMw.ApplyFn(galleriesC.ImageUpload)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}", galleriesC.Show).Methods("GET").Name(controllers.ShowGallery)

	fmt.Println("Starting Server on PORT 4500")
	// apply middleware on all routes
	http.ListenAndServe(":4500", userMw.Apply(r))
}

func faq(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "What questions do you have? Share them here and we would do our best to answer. :)")
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "Sorry, we couldn't get the page you requested")
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
