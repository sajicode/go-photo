package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if r.URL.Path == "/" {
		fmt.Fprint(w, "<h1>Welcome to my cool site</h1>")
	} else if r.URL.Path == "/contact" {
		fmt.Fprint(w, "To get in touch, please send an email to <a href=\"mailto:sajioloye@gmail.com\">support@lenslocked.com</a>.")
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "<h1>We could not find the page you requested :(</h1><p>Please send an email if you keep getting sent to an invalid page.</p>")
	}
}

func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}

func main() {
	router := httprouter.New()
	router.GET("/hello/:name", Hello)

	// http.HandleFunc("/", handlerFunc)
	http.ListenAndServe(":4500", router)
}
