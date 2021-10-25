package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if r.URL.Path == "/" {
		fmt.Fprintf(w, "<h1> Welcome to my site </h1>")
	} else if r.URL.Path == "/contact" {
		fmt.Fprintf(w, "To get in touch, please send an email to <a href=\"mailto:support@lenslocked.com\">support@lenslocked.com</a>")
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "<h1>We could not find the page you were looking for :(</h1> <p>Please email us if needed</p>")
	}

}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", handlerFunc)
	r.HandleFunc("/contact", handlerFunc)
	http.ListenAndServe(":3000", r)
}
