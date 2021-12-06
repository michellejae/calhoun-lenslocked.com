package main

import (
	"fmt"

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
	//us.AutoMigrate()
	user := models.User{
		Name:     "Jon Jon",
		Email:    "jon@jon.com",
		Password: "jon",
		Remember: "abc123",
	}

	err = us.Create(&user)
	if err != nil {
		panic(err)
	}
	fmt.Println(user)

	user2, err := us.ByRemember("abc123")
	if err != nil {
		panic(err)
	}

	fmt.Println(user2)
}
