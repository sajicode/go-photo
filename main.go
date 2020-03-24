package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sajicode/go-photo/controllers"
)

func faq(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "What questions do you have? Share them here and we would do our best to answer. :)")
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "Sorry, we couldn't get the page you requested")
}

func main() {
	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers()

	r := mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(notFound)
	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	r.HandleFunc("/faq", faq).Methods("GET")
	http.ListenAndServe(":4500", r)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
