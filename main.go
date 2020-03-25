package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sajicode/go-photo/controllers"
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
	us, err := models.NewUserService(dbDriver, psqlInfo)
	must(err)

	defer us.Close()
	//! to clear db
	// us.DestructiveReset()

	us.AutoMigrate()

	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(us)

	r := mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(notFound)
	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	r.HandleFunc("/cookietest", usersC.CookieTest).Methods("GET")
	r.Handle("/login", usersC.LoginView).Methods("GET")
	r.HandleFunc("/login", usersC.Login).Methods("POST")
	r.HandleFunc("/faq", faq).Methods("GET")
	fmt.Println("Starting Server on PORT 4500")
	http.ListenAndServe(":4500", r)
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
