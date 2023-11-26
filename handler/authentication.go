package handler

import (
	"diary_api/database"
	"diary_api/model"
	"diary_api/utility"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Register: validates JSON request, creates new user,
// and writes details of saved user to JSON response.
func Register(ctx *gin.Context) {

	// expected authentication input from request body
	var authInput model.AuthenticationInput

	// unmarshal request body into expected input
	err := ctx.ShouldBindJSON(&authInput)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// instantiate user
	user := model.User{
		Username: authInput.Username,
		Password: authInput.Password,
	}

	// save user to database
	savedUser, err := user.Save(database.Database)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// write saved user to response body
	ctx.JSON(http.StatusCreated, gin.H{"user": savedUser})
}

// Login: validates request, locates user if exists, validates password, generates JWT and writes the token to response body.
func Login(ctx *gin.Context) {

	// define expected authentication input from request body
	var authInput model.AuthenticationInput

	// unmarshal request body into expected input
	err := ctx.ShouldBindJSON(&authInput)

	fmt.Println("\n\nAUTH Input:", authInput)

	fmt.Println("\n\nAUTH Username:", authInput.Username)

	fmt.Println("\n\nAUTH Password:", authInput.Password)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// locate existing user by username
	user, err := model.FindUserByUsername(authInput.Username, database.Database)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// validate user password
	err = user.ValidatePassword(authInput.Password)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// generate JWT based on the user attempting to signin
	jwt, err := utility.GenerateJWT(user)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// write jwt to response body
	ctx.JSON(http.StatusOK, gin.H{"jwt": jwt})

}
