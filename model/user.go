package model

import (
	"diary_api/database"
	"html"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	// embed gorm.Model for associated fields
	gorm.Model
	Username string `gorm:"size:255;not null;unique" json:"username"`
	Password string `gorm:"size:255;not null" json:"-"`
	Entries  []Entry
}

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

// ValidatePassword: validates a provided password for a given user.
func (user *User) ValidatePassword(password string) error {

	// generate and compare hash's
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}

// FindUserByUsername: query database to find user with the corresponding username.
func FindUserByUsername(username string) (User, error) {
	var user User

	// query database to find user with matching username
	// if found load the user into the User struct defined
	err := database.Database.Where("username=?", username).Find(&user).Error

	if err != nil {
		return User{}, err
	}

	return user, nil
}

// FindUserByID: query database to find user with the corresponding ID and extract all entries populating them into the user struct.
func FindUserByID(id uint) (User, error) {
	var user User

	// TODO: test with Entries and ID lowercase
	err := database.Database.Preload("Entries").Where("ID=?", id).Find(&user).Error

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
