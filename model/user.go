package model

import (
	"diary_api/database"
	"html"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"size:255;not null;unique" json:"username"`
	Password string `gorm:"size:255;not null" json:"-"`
	Entries  []Entry
}

//--------------------User Methods--------------------//

// Save: insert user into database.
func (user *User) Save() (*User, error) {

	// insert user into database if no errors are encountered
	err := database.Database.Create(&user).Error

	if err != nil {
		return &User{}, err
	}
	return user, nil
}

func (user *User) BeforeSave(*gorm.DB) error {

	// hash password before insertion into database for security purposes
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	user.Password = string(passwordHash)

	// trim whitespace from username
	user.Username = html.EscapeString(strings.TrimSpace(user.Username))

	return nil

}

// ValidatePassword: validate provided password for a given user.
func (user *User) ValidatePassword(password string) error {

	// generate and compare hash's
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}

//--------------------User Functions--------------------//

// FindUserByUsername: query database to find user with corresponding username.
func FindUserByUsername(username string) (User, error) {

	// define user object to be loaded
	var user User

	// query database to find user with matching username
	// if found load the user into the user object defined
	err := database.Database.Where("username=?", username).Find(&user).Error

	if err != nil {
		return User{}, err
	}

	return user, nil
}

// FindUserByID: query database to find user by ID
// return user and all associated entries.
func FindUserByID(id uint) (User, error) {
	var user User

	err := database.Database.Preload("Entries").Where("id=?", id).Find(&user).Error

	if err != nil {
		return User{}, err
	}

	return user, nil
}

// gorm.Model

//   - a struct defined in the gorm package
//   - meant to be embedded within model struct's
//   - the following fields are embedded into the parent struct:
//       ~ ID
//       ~ CreatedAt
//       ~ UpdatedAt
//       ~ DeletedAt

// Exclude Field in JSON

//   - the struct tag `json:"-"` will exclude the field from being written to JSON

// Database.Preload

//   - case-sensitive, the table name must be capitalized
//   - if `unsupported relations for schema` error is encounter check capitalization of Preload
