package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.com/michellejae/lenslocked.com/controllers"
	"gitlab.com/michellejae/lenslocked.com/models"
)

const (
	host   = "localhost"
	port   = 5432
	user   = "macbookprowoe"
	dbname = "lenslocked_dev"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable", host, port, user, dbname)

	us, err := models.NewUserService(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer us.Close()
	us.DestructiveReset()

	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(us)

	r := mux.NewRouter()
	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	fmt.Println("Starting the server on :3000")
	http.ListenAndServe(":3000", r)
}
