package model

import (
	"errors"
	"fmt"
	"html"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type User struct {
	gorm.Model
	Fname string `gorm:"size:255;not null" json:"first_name" binding:"required"`
	Lname string `gorm:"size:255;not null" json:"last_name" binding:"required"`
	// TODO: figure out why generated column does not work
	// FullName string `gorm:"->;type:GENERATED ALWAYS AS (concat(fname,' ',lname));"`
	Email       string `gorm:"size:255;not null;unique" json:"email" binding:"required,email"`
	Phone       string `json:"phone_number" binding:"required,e164"`
	Password    string `gorm:"size:255;not null" json:"-" binding:"required"`
	Entries     []Entry
	LikedQuotes []LikedQuote
}

type UpdatedUser struct {
	gorm.Model
	Fname       string `gorm:"size:255;not null" json:"first_name" binding:"required"`
	Lname       string `gorm:"size:255;not null" json:"last_name" binding:"required"`
	Email       string `gorm:"size:255;not null;unique" json:"email" binding:"required,email"`
	Phone       string `json:"phone_number" binding:"required,e164"`
	Entries     []Entry
	LikedQuotes []LikedQuote
}

type UserEmail struct {
	Email string `json:"email" binding:"required,email"`
}

//--------------------User Methods--------------------//

// Save: insert user into database.
func (user *User) Save(db *gorm.DB) (*User, error) {

	// insert user into database if no errors are encountered
	err := db.Create(&user).Error

	if err != nil {
		if strings.Contains(err.Error(), "users_email_key") {
			return &User{}, errors.New("a user already exists with the associated email address")
		}
		return &User{}, err
	}
	return user, nil
}

// BeforeSave: gorm hook invoked before a user is inserted into the database
// hash user password for security
func (user *User) BeforeSave(*gorm.DB) error {

	// hash password before insertion into database for security purposes
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	user.Password = string(passwordHash)

	// trim whitespace from email
	user.Email = html.EscapeString(strings.TrimSpace(user.Email))

	return nil

}

// ValidatePassword: validate provided password for a given user.
func (user *User) ValidatePassword(password string) error {

	// generate and compare hash's
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return errors.New("password incorrect")
	}

	return nil
}

// Update: update user record instance in database
func (user *User) Update(db *gorm.DB, updatedUser UpdatedUser) (User, error) {

	originalUser := *user

	var result *gorm.DB

	// specify model you want to perfom operations on and update record | ensure you pass a new instance of the model to updates or the UpdatedAt field will not be updated
	if result = db.Model(user).Updates(User{
		Fname: updatedUser.Fname,
		Lname: updatedUser.Lname,
		Email: updatedUser.Email,
		Phone: updatedUser.Phone,
	}); result.Error != nil {

		return User{}, result.Error

	}

	log.Printf(
		"\nupdated user data from: %+v\n\nto: %+v\n\n",
		gin.H{
			"first_name": originalUser.Fname,
			"last_name":  originalUser.Lname,
			"email":      originalUser.Email,
			"phone":      originalUser.Phone,
		},
		gin.H{
			"first_name": user.Fname,
			"last_name":  user.Lname,
			"email":      user.Email,
			"phone":      user.Phone,
		},
	)

	log.Println("rows affected:", result.RowsAffected)

	return *user, nil

}

// Delete: PERMANENTLY delete user instance and all associated instance entries from the database.
func (user *User) Delete(db *gorm.DB) error {

	// select user by id and all associations then delete that user and all associations  permanently
	if err := db.Unscoped().Select(clause.Associations).Delete(&user).Error; err != nil {
		return err
	}

	return nil
}

//--------------------User Functions--------------------//

// FindUserByUsername: query database to find user with corresponding email.
func FindUserByEmail(email string, db *gorm.DB) (User, error) {

	// define user object to be loaded
	var user User

	fmt.Println("\n\nUsername:", email)

	// query database to find user with matching email
	// if found load the user into the user object defined
	if err := db.Where("email=?", email).First(&user).Error; err != nil {
		return User{}, errors.New("user not found")
	}

	fmt.Println("\n\nUser:", user)

	return user, nil
}

// FindUserByID: query database to find user by ID
// return user and all associated entries.
func FindUserByID(id uint, db *gorm.DB, preloadEntries bool) (User, error) {
	var user User

	if preloadEntries {
		// initialize user and associated entries with Preload
		if err := db.Preload("Entries").Preload("LikedQuotes").Where("id=?", id).First(&user).Error; err != nil {

			return User{}, err

		}
	} else {
		// initialize user and without associated entries with Preload
		if err := db.Where("id=?", id).First(&user).Error; err != nil {

			return User{}, err

		}
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

// gorm Hooks

//   - methods you can add to a model to do pre and post proccessing
//     before and after some database action

// Delete Single Record Instance (Entity - Row - Class - Struct) and All Associations (Instance Records - Rows that are linked)

//   - db.Unscoped().Select(clause.Associations).Delete(&PointerToStructInstanceYouWouldLikeToDelete)
//   - CAUTION: if the struct you pass in does not have a primary key a batch delete will be executed
