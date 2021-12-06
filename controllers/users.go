package controllers

import (
	"fmt"
	"net/http"

	"gitlab.com/michellejae/lenslocked.com/models"
	"gitlab.com/michellejae/lenslocked.com/views"
)

// NewUsers is used to create a new Users Controller
//this function will panic if templates are not parsed correctly and should only be used during initial setup
func NewUsers(us *models.UserService) *Users {
	return &Users{
		NewView:   views.NewView("bootstrap", "users/new"),
		LoginView: views.NewView("bootstrap", "users/login"),
		us:        us,
	}
}

type Users struct {
	NewView   *views.View
	LoginView *views.View
	us        *models.UserService
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
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

// Create is u sed to process the signup form when a user submits it. This is used to create a new user account
// POST /signup
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {

	var form SignupForm
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}
	user := models.User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}
	if err := u.us.Create(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, user)
}

type LoginForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

// login used to verify the provided emil address and password
// and then login the user if correct
// POST / login
func (u *Users) Login(w http.ResponseWriter, r *http.Request) {

	form := LoginForm{}

	if err := parseForm(r, &form); err != nil {
		panic(err)
	}

	user, err := u.us.Authenticate(form.Email, form.Password)
	if err != nil {
		switch err {
		case models.ErrNotFound:
			fmt.Fprintln(w, "Invalid Email Address")
		case models.ErrInvalidPassword:
			fmt.Fprintln(w, "Invalid Password")
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	}
	cookie := http.Cookie{
		Name:  "email",
		Value: user.Email,
	}
	http.SetCookie(w, &cookie)
	fmt.Fprintln(w, user)

}
