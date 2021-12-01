package main

import (
	"fmt"

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

	// shows us what sql statemetns gorm is running
	db.LogMode(true)
	//db.DropTableIfExists(&User{})
	db.AutoMigrate(&User{})

	//var u User

	// db.First(&u)
	// db.Last(&u)

	// can add in the second parameter with rfirst to qeruty taht specific user with that id
	//db.First(&u, 3)
	// longer version of above
	// db = db.Where("id = ?", 3)
	// db.First(&u)

	// note the ?, it's like the $ from psql and the . after each line
	// db.Where("color = ?", "blue").
	// 	Where("id > ?", 2).
	// 	First(&u)

	// var u User = User{
	// 	Color: "blue",
	// 	Name:  "meeeee",
	// }
	// db.Where(u).First(&u)

	// mulitple users
	var users []User
	db.Find(&users)
	fmt.Println(len(users))
	fmt.Println(users)

}
