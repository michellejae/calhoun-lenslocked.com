package main

import (
	"fmt"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"gitlab.com/michellejae/lenslocked.com/models"
)

const (
	host   = "localhost"
	port   = 5432
	user   = "macbookprowoe"
	dbname = "lenslocked_dev"
)

// type User struct {
// 	gorm.Model
// 	Name   string
// 	Email  string `gorm:"not null;unique_index"`
// 	Color  string
// 	Orders []Order
// }

// type Order struct {
// 	gorm.Model
// 	UserID      uint
// 	Amount      int
// 	Description string
// }

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable", host, port, user, dbname)
	us, err := models.NewUserService(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer us.Close()
	// us.DestructiveReset()
	// user := models.User{
	// 	Name:  "Michael Scott",
	// 	Email: "michael@dundermifflin.com",
	// }

	// if err := us.Create(&user); err != nil {
	// 	panic(err)
	// }

	user, err := us.ByID(1)
	if err != nil {
		panic(err)
	}
	fmt.Println(user)

	// shows us what sql statemetns gorm is running
	//db.LogMode(true)
	//db.DropTableIfExists(&User{})
	//db.AutoMigrate(&User{}, &Order{})

	//var u User
	// QUERY 1
	// db.First(&u)
	// db.Last(&u)

	// QEURTY TWO
	// can add in the second parameter with rfirst to qeruty taht specific user with that id
	//db.First(&u, 3)
	// longer version of above
	// db = db.Where("id = ?", 3)
	// db.First(&u)

	//QUERY 3
	// note the ?, it's like the $ from psql and the . after each line
	// db.Where("color = ?", "blue").
	// 	Where("id > ?", 2).
	// 	First(&u)

	// QUERY 4
	// var u User = User{
	// 	Color: "blue",
	// 	Name:  "meeeee",
	// }
	// db.Where(u).First(&u)

	// QUERY 5
	// mulitple users
	// var users []User
	// db.Find(&users)
	// fmt.Println(len(users))
	// fmt.Println(users)

	//var u User
	// newDB := db.Where("email = ?", "blahs")
	// newDB = newDB.Or("color = ?", "blue")
	// newDB = newDB.First(&u)

	// ERROR NO WORK
	// THIS DOES NOT PANIC EVEN THOUGH IT SHOULD
	// db.Where("email = ?", "blah").First(&u)
	// if db.Error != nil {
	// 	panic(err)
	// }

	// ERROR 1
	// THIS DOES PANIC cause we assign it to a new variable and check THAT error
	// newDB := db.Where("email = ?", "blah").First(&u)
	// if newDB.Error != nil {
	// 	panic(err)
	// }

	// ERROR 2
	// db = db.Where("email = ?", "blah").First(&u)
	// errors := db.GetErrors()
	// if len(errors) > 0 {
	// 	fmt.Println(errors)
	// 	os.Exit(1)
	// }

	//ERROR 3
	// db = db.Where("email = ?", "blah").First(&u)
	// if db.RecordNotFound() {
	// 	fmt.Println("No user found")
	// } else if db.Error != nil {
	// 	panic(db.Error)
	// } else {
	// 	fmt.Println(u)
	// }

	// ERROR 4
	// db = db.Where("email = ?", "blah").First(&u)
	// if err := db.Where("email = ?", "blah").First(&u).Error; err != nil {
	// 	switch err {
	// 	case gorm.ErrRecordNotFound:
	// 		fmt.Println("no user found")
	// 	default:
	// 		panic(err)
	// 	}
	// }
	// fmt.Println(u)

	// RELATIONAL DATABASE USER STUFF
	// var u User
	// if err := db.Preload("Orders").First(&u).Error; err != nil {
	// 	panic(err)
	// }
	// fmt.Println(u)
	// fmt.Println(u.Orders)

	// added orders to DB
	// createOrder(db, u, 1001, "Chicken Wings")
	// createOrder(db, u, 999, "Fries")
	// createOrder(db, u, 100, "Ranch")
}

// func createOrder(db *gorm.DB, user User, amount int, desc string) {
// 	err := db.Create(&Order{
// 		UserID:      user.ID,
// 		Amount:      amount,
// 		Description: desc,
// 	}).Error
// 	if err != nil {
// 		panic(err)
// 	}
// }
