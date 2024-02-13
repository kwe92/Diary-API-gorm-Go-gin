package handler

import (
	"journal_api/model"
	"journal_api/utility"
	"log"
	"net/http"
	"strconv"

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
		user, err := utility.CurrentUser(ctx, db, false)

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

		log.Println("\nnew entry added:", *savedEntry)

	}
}

// UpdateEntry: http handler that updates a single entry in database
func UpdateEntry(db *gorm.DB) gin.HandlerFunc {

	return func(ctx *gin.Context) {

		// expected request input struct
		var updateEntryInput model.UpdatedEntryInput

		// declare destination struct
		var entry model.Entry

		var err error

		entryId, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// find record to update and load record into destination struct
		if entry, err = model.FindEntryById(db, uint(entryId)); err != nil {

			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return

		}

		// read `deserialize` request body buffer into expected input struct
		if err := ctx.ShouldBindJSON(&updateEntryInput); err != nil {

			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

			return
		}

		// call Update method on object to update record in database
		if updatedEntry, err := entry.Update(db, updateEntryInput); err != nil {

			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return

		} else {

			ctx.JSON(http.StatusOK, gin.H{"updated_entry": updatedEntry})

		}

	}
}

// DeleteEntry: http handler that deletes a single entry in database
func DeleteEntry(db *gorm.DB) gin.HandlerFunc {

	return func(ctx *gin.Context) {

		entryId, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// find record to delete and load record into destination struct
		if entry, err := model.FindEntryById(db, uint(entryId)); err != nil {

			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return

		} else {

			// call Delete method on object to delete record from database
			if deletedEntry, err := entry.Delete(db); err != nil {

				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return

			} else {

				log.Println("\ndeleted entry:", *deletedEntry)

				// return deleted object with DeletedAt time
				ctx.JSON(http.StatusOK, gin.H{"deleted_entry": gin.H{
					"id":        deletedEntry.ID,
					"content":   deletedEntry.Content,
					"DeletedAt": deletedEntry.DeletedAt,
				}})

			}

		}

	}

}

// GetAllEntries: route handler that retrieves current user and returns all associated entries as a response
func GetAllEntries(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// retrieve current authenticated user from request header
		user, err := utility.CurrentUser(ctx, db, true)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// write current user entries to response body
		ctx.JSON(http.StatusOK, gin.H{"data": user.Entries})
	}
}
