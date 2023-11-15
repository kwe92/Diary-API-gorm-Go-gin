package controller

import (
	"diary_api/helper"
	"diary_api/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Register: validates JSON request, creates a new user,
// and returns the details of the saved use as a JSON response.
func Register(context *gin.Context) {

	// struct of expected input from request
	var input model.AuthenticationInput

	// write request body into AuthenticationInput struct
	err := context.ShouldBindJSON(&input)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// create a user from request input
	user := model.User{
		Username: input.Username,
		Password: input.Password,
	}

	// ! should user.beforeSave be called here?

	// save user to database
	savedUser, err := user.Save()

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"user": savedUser})
}

// TODO: Add comments----------------------------------------------

func Login(context *gin.Context) {
	var input model.AuthenticationInput

	err := context.ShouldBindJSON(&input)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := model.FindUserByUsername(input.Username)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = user.ValidatePassword(input.Password)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jwt, err := helper.GenerateJWT(user)

	if err != nil {
		context.JSON(http.StatusBadRequest, err.Error())
		return
	}

	context.JSON(http.StatusOK, gin.H{"jwt": jwt})

}

// TODO END----------------------------------------------
