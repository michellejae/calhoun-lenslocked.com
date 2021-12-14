package models

// import (
// 	"fmt"
// 	"testing"
// 	"time"
// )

// func testingUserService() (*UserService, error) {

// 	const (
// 		host   = "localhost"
// 		port   = 5432
// 		user   = "macbookprowoe"
// 		dbname = "lenslocked_test"
// 	)

// 	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable", host, port, user, dbname)

// 	us, err := NewUserService(psqlInfo)
// 	if err != nil {
// 		return nil, err
// 	}
// 	us.db.LogMode(false)

// 	// clear the users table between tests
// 	us.DestructiveReset()
// 	return us, nil
// }

// func TestCreateUser(t *testing.T) {
// 	us, err := testingUserService()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	user := User{
// 		Name:  "Michael Scott",
// 		Email: "michael@dundermifflin.com",
// 	}

// 	// check to see if any error came back fcrom create
// 	err = us.Create(&user)
// 	if err != nil {
// 		// fatal will stop the rest of the test
// 		t.Fatal(err)
// 	}
// 	// id will only be a positive number cause it's a uint so we check to make sure it's not 0 as that woud
// 	// be the only random edgecase
// 	if user.ID == 0 {
// 		// errorF will log the error but conitnue running them
// 		t.Errorf("Expeted ID > 0. Recieed %d", user.ID)
// 	}
// 	// createdAt & updatedAt shuold be created within the last 5 seconds for  user
// 	if time.Since(user.CreatedAt) > time.Duration(5*time.Second) {
// 		t.Errorf("Expected CreatedAt to be recent. Received %s", user.CreatedAt)
// 	}

// 	if time.Since(user.UpdatedAt) > time.Duration(5*time.Second) {
// 		t.Errorf("Expected UpdatedAt to be recent. Received %s", user.UpdatedAt)
// 	}

// }
