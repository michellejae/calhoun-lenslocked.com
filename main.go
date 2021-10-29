package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.com/michellejae/lenslocked.com/controllers"
	"gitlab.com/michellejae/lenslocked.com/views"
)

var (
	homeView    *views.View
	contactView *views.View
)

// response for homePage
func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	must(homeView.Render(w, nil))

}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	must(contactView.Render(w, nil))
}


func main() {
	// we run these two funcs here to parse the html files when we first set up and start our project.
	// this way we know right away if we have issues with our html pages (the errors/panic happen on view.go)
	// vs waiting until someone hits the page and then finding out there we have error from our templates
	// they will not be "execute/render" until each page is hit, which is fine.
	homeView = views.NewView("bootstrap", "views/home.gohtml") // bootstrap is from the views/layouts/folder
	contactView = views.NewView("bootstrap", "views/contact.gohtml")
	usersC := controllers.NewUsers()

	r := mux.NewRouter()
	r.HandleFunc("/", home)
	r.HandleFunc("/contact", contact)
	r.HandleFunc("/signup", usersC.New)
	http.ListenAndServe(":3000", r)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
