package controller

import (
	"diary_api/helper"
	"diary_api/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AddEntry: route handler to add a new entry of an authorized user
func AddEntry(context *gin.Context) {
	// declare expected request body struct
	var input model.Entry

	// unmarshal request body into struct
	err := context.ShouldBindJSON(&input)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// retrieve the currently authenticated user from request header
	user, err := helper.CurrentUser(context)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// map current user id to entry
	input.UserID = user.ID

	// insert the entry into the database
	savedEntry, err := input.Save()

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// return a response with newly saved entry
	context.JSON(http.StatusCreated, gin.H{"data": savedEntry})
}

// GetAllEntries: route handler that retrieves current user and returns all associated entries as a response
func GetAllEntries(context *gin.Context) {

	// retrieve the currently authenticated user from request header
	user, err := helper.CurrentUser(context)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// write to and send response body
	context.JSON(http.StatusOK, gin.H{"data": user.Entries})

}
