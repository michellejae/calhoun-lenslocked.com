package controllers

import (
	"fmt"
	"log"
	"net/http"

	"gitlab.com/michellejae/lenslocked.com/context"
	"gitlab.com/michellejae/lenslocked.com/models"
	"gitlab.com/michellejae/lenslocked.com/views"
)

func NewGalleries(gs models.GalleryService) *Galleries {
	return &Galleries{
		New: views.NewView("bootstrap", "galleries/new"),
		gs:  gs,
	}
}

type Galleries struct {
	New *views.View
	gs  models.GalleryService
}

type GalleryForm struct {
	Title string `schema:"title"`
}

// Create is used to process the signup form when a user submits it. This is used to create a new user account
// POST /signup
func (g *Galleries) Create(w http.ResponseWriter, r *http.Request) {

	var vd views.Data
	var form GalleryForm
	if err := parseForm(r, &form); err != nil {
		log.Println(err)
		vd.SetAlert(err)
		g.New.Render(w, vd)
		return
	}
	user := context.User(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	fmt.Println("crate got the user:", user)
	gallery := models.Gallery{
		Title:  form.Title,
		UserID: user.ID,
	}
	if err := g.gs.Create(&gallery); err != nil {
		vd.SetAlert(err)
		g.New.Render(w, vd)
		return
	}
	fmt.Fprintln(w, gallery)
}
