package handler

import (
	"diary_api/database"
	"diary_api/model"
	"diary_api/utility"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AddEntry: route handler to add a new entry of an authorized user
func AddEntry(context *gin.Context) {
	// declare expected request body
	var entry model.Entry

	// unmarshal request body into struct
	err := context.ShouldBindJSON(&entry)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// retrieve currently authenticated user from request header
	user, err := utility.CurrentUser(context)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// map current user id to entry
	entry.UserID = user.ID

	// insert entry into the database
	savedEntry, err := entry.Save(database.Database)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// write new entry to response body
	context.JSON(http.StatusCreated, gin.H{"data": savedEntry})
}

// GetAllEntries: route handler that retrieves current user and returns all associated entries as a response
func GetAllEntries(context *gin.Context) {

	// retrieve current authenticated user from request header
	user, err := utility.CurrentUser(context)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// write current user entries to response body
	context.JSON(http.StatusOK, gin.H{"data": user.Entries})

}
