package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const (
	host   = "localhost"
	port   = 5432
	user   = "macbookprowoe"
	dbname = "lenslocked_dev"
)

type User struct {
	gorm.Model
	Name  string
	Email string `gorm:"not null;unique_index"`
	Color string
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable", host, port, user, dbname)
	db, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := db.DB().Ping(); err != nil {
		panic(err)
	}

	// shows us what sql statemetns gorm is running
	db.LogMode(true)
	//db.DropTableIfExists(&User{})
	db.AutoMigrate(&User{})

	name, email, color := getInfo()
	u := User{
		Name:  name,
		Email: email,
		Color: color,
	}
	if err = db.Create(&u).Error; err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", u)

}

func getInfo() (name, email, color string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("What is your name?")
	name, _ = reader.ReadString('\n')
	fmt.Println("What is your email address?")
	email, _ = reader.ReadString('\n')
	fmt.Println("What is your fav color?")
	color, _ = reader.ReadString('\n')
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)
	color = strings.TrimSpace(color)
	return name, email, color

}
