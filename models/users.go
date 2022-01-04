package models

import (
	"regexp"
	"strings"
	"time"

	"github.com/jinzhu/gorm"

	"gitlab.com/michellejae/lenslocked.com/hash"
	"gitlab.com/michellejae/lenslocked.com/rand"
	"golang.org/x/crypto/bcrypt"
)

// User represents the user model stored in our database
// this is used for user accounts storing an email adddress and a password
// so users can login and gain access to their content
type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"not null;unique_index"`
	Password     string `gorm:"-"` //- means ignore this and don't store in db
	PasswordHash string `gorm:"not null"`
	Remember     string `gorm:"-"`
	RememberHash string `gorm:"not null;unique_index"`
}

//userDB is used to interact with the usrs database

// If the user is found, we will return a nil error
// If the user is not found, we will return ErrNotFound
//If there is another error, we will return an error with more information
// about what went wrong. This may not be an error generated by the model package

// for single user queries, as a general rule, any error but ErrNotFound should probaly result in a 500 error

type UserDB interface {
	// methodd for querying or single Users
	ByID(id uint) (*User, error)
	ByEmail(email string) (*User, error)
	ByRemember(token string) (*User, error)

	// methods for altering users
	Create(user *User) error
	Update(user *User) error
	Delete(id uint) error
}

//UserService is a set of methods used to manipulate and work with the user model
type UserService interface {
	// Authenticate will verify provided email address and password are correct
	// If correct, the user corresponding to the given email will be returned
	// If they are wrong, you will receive either:
	// errorNotFound (user), errorInvalidPassword, or error if something goes wrong
	Authenticate(email, password string) (*User, error)
	// InitiateReset will start the initiate process by creating a reset token for user found with provided email address
	InitiateReset(email string) (string, error)
	CompleteReset(token, newPw string) (*User, error)
	UserDB
}

// user service now only handles authenticate
// accesses userGorm in order to acccess DB
// accesses userValid to validate and normalize
func NewUserService(db *gorm.DB, pepper, hmacKey string) UserService {
	ug := &userGorm{db}
	hmac := hash.NewHmac(hmacKey)
	uv := newUserValidator(ug, hmac, pepper)
	return &userService{
		UserDB:    uv,
		pepper:    pepper,
		pwResetDB: newPwResetValidator(&pwResetGorm{db}, hmac),
	}
}

var _ UserService = &userService{}

// UserService type that will chain into UserDB interface
type userService struct {
	UserDB
	pepper    string
	pwResetDB pwResetDB
}

// can be used to authenticate a user with the provided email address and password
// if the email provided is invalid this will return nil, ErrNotFound
// if the password provided is invalid this will return nil, ErrPasswordIncorrect
// if both are valid, will return user, nil
// if both are invalid will return nil, error

// we keep this as a userService method cause though we have to grab info from DB
// it's not really doing anything with DB. it's more about the USER
func (us *userService) Authenticate(email, password string) (*User, error) {
	// this is the uv byEmail
	foundUser, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash), []byte(password+us.pepper))
	if err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return nil, ErrPasswordIncorrect
		default:
			return nil, err
		}

	}
	return foundUser, nil
}

func (us *userService) InitiateReset(email string) (string, error) {
	user, err := us.ByEmail(email)
	if err != nil {
		return "", err
	}
	pwr := pwReset{
		UserID: user.ID,
	}
	if err := us.pwResetDB.Create(&pwr); err != nil {
		return "", err
	}
	return pwr.Token, nil
}

func (us *userService) CompleteReset(token, newPw string) (*User, error) {
	// 1. lookup pwReset using the token
	pwr, err := us.pwResetDB.ByToken(token)
	if err != nil {
		if err == ErrNotFound {
			return nil, ErrTokenInvalid
		}
		return nil, err
	}
	// 2. check that token is not expired (older then 12 hours)
	if time.Now().Sub(pwr.CreatedAt) > (12 * time.Hour) {
		return nil, ErrTokenInvalid
	}

	// 3. lookup user by the pwReset.UserID
	user, err := us.ByID(pwr.UserID)
	if err != nil {
		return nil, err
	}
	// 4. update the user's password w/ newPw
	user.Password = newPw
	err = us.Update(user)
	if err != nil {
		return nil, err
	}
	// 5. delete the pwReset
	us.pwResetDB.Delete(pwr.ID)
	// 6. return user
	return user, nil

}

type userValFunc func(*User) error

func runUserValFuncs(user *User, fns ...userValFunc) error {
	for _, fn := range fns {
		if err := fn(user); err != nil {
			return err
		}
	}
	return nil
}

var _ UserDB = &userValidator{}

func newUserValidator(udb UserDB, hmac hash.HMAC, pepper string) *userValidator {
	return &userValidator{
		UserDB: udb,
		hmac:   hmac,
		// emailRegex is used to match email addresses. it is not perfect but works okay for now
		emailRegex: regexp.MustCompile(`^[a-z0-9._%+\-]+@+[a-z0-9.\-]+\.[a-z]{2,16}$`),
		pepper:     pepper,
	}
}

type userValidator struct {
	UserDB
	hmac       hash.HMAC
	emailRegex *regexp.Regexp
	pepper     string
}

// ByEmail will normlize the email addresss before calling ByEmail on the UserField
func (uv *userValidator) ByEmail(email string) (*User, error) {
	user := User{
		Email: email,
	}
	if err := runUserValFuncs(&user, uv.normalizeEmail); err != nil {
		return nil, err
	}
	return uv.UserDB.ByEmail(user.Email)
}

// ByRemember will hash the remember token and then call ByRemember on the subsequent UserDB layer
// which i believe is the userGorm ByRemember
func (uv *userValidator) ByRemember(token string) (*User, error) {
	user := User{
		Remember: token,
	}
	if err := runUserValFuncs(&user, uv.hmacRemember); err != nil {
		return nil, err
	}

	return uv.UserDB.ByRemember(user.RememberHash) //
}

//Creat provided user and backfill data like the ID, CreatedAt and UpdatedAt fields
func (uv *userValidator) Create(user *User) error {
	// setRemember has to go before hmacRemember
	// bcrypt has to happen before passwordHashRequired
	// rememberMinBytes has to be before the hmacRemember
	// hmac has to be before rememberHashRequired
	err := runUserValFuncs(user,
		uv.passwordRequired,
		uv.passwordMinLength,
		uv.bcryptPassword,
		uv.passwordHashRequired,
		uv.setRememberIfUnset,
		uv.rememberMinBytes,
		uv.hmacRemember,
		uv.rememberHashRequired,
		uv.normalizeEmail,
		uv.requireEmail,
		uv.emailFormat,
		uv.emailIsAvail)
	if err != nil {
		return err
	}
	return uv.UserDB.Create(user)
}

// update will hash a remember token if it's provided
func (uv *userValidator) Update(user *User) error {

	err := runUserValFuncs(user,
		uv.passwordMinLength,
		uv.bcryptPassword,
		uv.passwordHashRequired,
		uv.rememberMinBytes,
		uv.hmacRemember,
		uv.rememberHashRequired,
		uv.normalizeEmail,
		uv.requireEmail,
		uv.emailFormat,
		uv.emailIsAvail)
	if err != nil {
		return err
	}
	return uv.UserDB.Update(user)
}

func (uv *userValidator) Delete(id uint) error {
	var user User
	user.ID = id
	err := runUserValFuncs(&user, uv.idGreatThan(0))
	if err != nil {
		return err
	}
	return uv.UserDB.Delete(id)
}

// method will only hash password
// will not validate if it's correct length or meets requirement
// uses predefined pepper
func (uv *userValidator) bcryptPassword(user *User) error {
	// if password is empty string return nil cause we will not hash empty
	if user.Password == "" {
		return nil
	}
	pwBytes := []byte(user.Password + uv.pepper)
	hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""
	return nil
}

func (uv *userValidator) hmacRemember(user *User) error {
	if user.Remember == "" {
		return nil
	}
	user.RememberHash = uv.hmac.Hash(user.Remember)
	return nil
}

// only will be used on Create method as if we update a user we don't want to update the remember token
func (uv *userValidator) setRememberIfUnset(user *User) error {
	if user.Remember != "" {
		return nil
	}
	token, err := rand.RememberToken()
	if err != nil {
		return err
	}
	user.Remember = token
	return nil
}

func (uv *userValidator) rememberMinBytes(user *User) error {
	if user.Remember == "" {
		return nil
	}
	n, err := rand.NBytes(user.Remember)
	if err != nil {
		return err
	}
	if n < 32 {
		return ErrRememberTooShort
	}
	return nil
}

func (uv *userValidator) rememberHashRequired(user *User) error {
	if user.RememberHash == "" {
		return ErrRememberRequired
	}

	return nil
}

// same as functino as the one below but it uses closures and functions as paramenters
func (uv *userValidator) idGreatThan(n uint) userValFunc {
	return userValFunc(func(user *User) error {
		if user.ID <= n {
			return ErrIDInvalid
		}
		return nil
	})
}

// func (uv *userValidator) idGreaterThanZero(user *User) error {
// 	if user.ID <= 0 {
// 		return ErrIDInvalid
// 	}
// 	return nil
// }

// don't have to use this on Authenticate cause we use ByEmail inside of it so it's normalized in that func
// otherwise this method wil be used on uv's ByEmail, Create and Update
// lower cases all emaill addresses from users
// trims white spaces
// then sends to the ByEmail method on ug to send to DB
func (uv *userValidator) normalizeEmail(user *User) error {
	user.Email = strings.ToLower(user.Email)
	user.Email = strings.TrimSpace(user.Email)
	return nil
}

func (uv *userValidator) requireEmail(user *User) error {
	if user.Email == "" {
		return ErrEmailRequired
	}
	return nil
}

func (uv *userValidator) emailFormat(user *User) error {
	if user.Email == "" {
		return nil
	}
	if !uv.emailRegex.MatchString(user.Email) {
		return ErrEmailInvalid
	}
	return nil
}

func (uv *userValidator) emailIsAvail(user *User) error {
	exisiting, err := uv.ByEmail(user.Email)
	if err == ErrNotFound {
		// email address is not taken
		return nil
	}
	if err != nil {
		return err
	}

	// found a user with this email
	// if the found user has the same ID as this user, it is an update and this is the same user
	if user.ID != exisiting.ID {
		return ErrEmailTaken
	}
	return nil
}

func (uv *userValidator) passwordMinLength(user *User) error {
	if user.Password == "" {
		return nil
	}
	if len(user.Password) < 8 {
		return ErrPasswordTooShort
	}
	return nil
}

func (uv *userValidator) passwordRequired(user *User) error {
	if user.Password == "" {
		return ErrPasswordRequired
	}

	return nil
}

func (uv *userValidator) passwordHashRequired(user *User) error {
	if user.PasswordHash == "" {
		return ErrPasswordRequired
	}

	return nil
}

// underscore tells compiler that this variable will never actually be used
// however setting UserDB to the userGorm ensures that the userGorm type always matches UserDB interface
var _ UserDB = &userGorm{}

type userGorm struct {
	db *gorm.DB
}

func (ug *userGorm) ByID(id uint) (*User, error) {
	var user User
	db := ug.db.Where("id = ?", id)
	err := first(db, &user)
	return &user, err
}

// same as above but byEMAIL
func (ug *userGorm) ByEmail(email string) (*User, error) {
	// this is the user where we play the return from the db
	var user User
	db := ug.db.Where("email = ?", email)
	err := first(db, &user)
	return &user, err
}

// looks up a user with the given remember token and retunr thats user
// this method expects the remember token to already be hashed
func (ug *userGorm) ByRemember(rememberHash string) (*User, error) {
	var user User
	// query has to use remember_hash even though we save it as RememberHash down below
	err := first(ug.db.Where("remember_hash = ?", rememberHash), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

//Creates the user in the actual DB with data provided
func (ug *userGorm) Create(user *User) error {
	return ug.db.Create(user).Error
}

// will update the provied user with the provided user object
// note that you ahve to provide ALL fields whether their data is updatd or not
// ie if the email is staying the same and you don't provide it, it will delete
func (ug *userGorm) Update(user *User) error {
	return ug.db.Save(user).Error
}

// will delete the user with the provided ID
func (ug *userGorm) Delete(id uint) error {
	// have to actually specify model to get acces to gorm.Model to get the ID
	user := User{Model: gorm.Model{ID: id}}
	return ug.db.Delete(&user).Error
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
