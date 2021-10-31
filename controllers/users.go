package controllers

import (
	"fmt"
	"net/http"

	"gitlab.com/michellejae/lenslocked.com/views"
)

// NewUsers is used to create a new Users Controller
//this function will panic if templates are not parsed correctly and should only be used during initial setup
func NewUsers() *Users {
	return &Users{
		NewView: views.NewView("bootstrap", "users/new"),
	}
}

type Users struct {
	NewView *views.View
}

// New is used to render the form where a user can create a new user account
// Get /signup
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	if err := u.NewView.Render(w, nil); err != nil {
		panic(err)
	}
}

// struct tags are not checked by compilier or something
type SignupForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

// Create is u sed to process the signup form when a user submits it. This is used to create a new user account
// POST /signup
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	var form SignupForm
	if err := u.parseForm(r, &form); err != nil {
		panic(err)
	}
	fmt.Fprintln(w, form)

}
