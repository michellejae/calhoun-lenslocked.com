package models

import (
	"errors"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"
)

var (
	// return when a resource cannot be found in DB
	ErrNotFound = errors.New("models: resource not found")

	// returned when an invalid id is provided to a mehod like delete
	ErrInvalidID = errors.New("models: id provided was invalid")
)

const userPwPepper = "booopity-beep-berp"

// how we connect to DB
func NewUserService(connectionInfo string) (*UserService, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	return &UserService{
		db: db,
	}, nil
}

// access gorm to give acces to DB model
type UserService struct {
	db *gorm.DB
}

// ByID will look up a user with provided ID.
// If the user is found, we will return a nil error
// If the user is not found, we will return ErrNotFound
//If there is another error, we will return an error with more information
// about what went wrong. This may not be an error generated by the model package

// as a general rule, any error but ErrNotFound should probaly result in a 500 error
func (us *UserService) ByID(id uint) (*User, error) {
	var user User
	db := us.db.Where("id = ?", id)
	err := first(db, &user)
	return &user, err
}

// same as above but byEMAIL
func (us *UserService) ByEmail(email string) (*User, error) {
	// this is the user where we play the return from the db
	var user User
	db := us.db.Where("email = ?", email)
	err := first(db, &user)
	return &user, err
}

// first will query using provided gorm.db and will
// get first item return and place it into dst
// if nothing is return (or fround in query) it will return err not found
func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	return err
}

//Creat provided user and backfill data like the ID, CreatedAt and UpdatedAt fields
func (us *UserService) Create(user *User) error {
	pwBytes := []byte(user.Password + userPwPepper)
	hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""
	return us.db.Create(user).Error
}

// will update the provied user with the provided user object
// note that you ahve to provide ALL fields whether their data is updatd or not
// ie if the email is staying the same and you don't provide it, it will delete
func (us *UserService) Update(user *User) error {
	return us.db.Save(user).Error
}

// will delete the user with the provided ID
func (us *UserService) Delete(id uint) error {
	if id == 0 {
		return ErrInvalidID
	}

	// have to actually specify model to get acces to gorm.Model to get the ID
	user := User{Model: gorm.Model{ID: id}}
	return us.db.Delete(&user).Error
}

// closes the UserService DB connection
func (us *UserService) Close() error {
	return us.db.Close()
}

// method only used in development to reset database just to make sure we are starting fresh
// drops user table and rebuilds it
func (us *UserService) DestructiveReset() error {
	if err := us.db.DropTableIfExists(&User{}).Error; err != nil {
		return err
	}
	return us.AutoMigrate()

}

// Automigrate will attempt to automatiaclly migrate the users table
func (us *UserService) AutoMigrate() error {
	if err := us.db.AutoMigrate(&User{}).Error; err != nil {
		return err
	}
	return nil
}

type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"not null;unique_index"`
	Password     string `gorm:"-"` //- means ignore this and don't store in db
	PasswordHash string `gorm:"not null"`
}
