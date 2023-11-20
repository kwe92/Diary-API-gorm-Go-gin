package handler

import (
	"diary_api/model"
	"diary_api/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AddEntry: route handler to add a new entry of an authorized user
func AddEntry(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// declare expected request body input
		var entry model.Entry

		// unmarshal request body into struct
		err := ctx.ShouldBindJSON(&entry)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// retrieve currently authenticated user from request header
		user, err := utility.CurrentUser(ctx, db)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// map current user id to entry
		entry.UserID = user.ID

		// insert entry into the database
		savedEntry, err := entry.Save(db)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// write new entry to response body
		ctx.JSON(http.StatusCreated, gin.H{"data": savedEntry})
	}
}

// GetAllEntries: route handler that retrieves current user and returns all associated entries as a response
func GetAllEntries(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// retrieve current authenticated user from request header
		user, err := utility.CurrentUser(ctx, db)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// write current user entries to response body
		ctx.JSON(http.StatusOK, gin.H{"data": user.Entries})
	}
}
