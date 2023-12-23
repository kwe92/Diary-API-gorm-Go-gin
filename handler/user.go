package handler

import (
	"errors"
	"journal_api/model"
	"journal_api/utility"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UpdateUser: http handler that updates a single user in database
func UpdateUser(db *gorm.DB) gin.HandlerFunc {

	return func(ctx *gin.Context) {

		// expected request input struct
		var updatedUserInput model.UpdatedUser

		// read `deserialize` request body buffer into expected input struct
		if err := ctx.ShouldBindJSON(&updatedUserInput); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// retrieve currently authenticated user instance as a Struct from the database using the request header information (jwt)
		user, err := utility.CurrentUser(ctx, db)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// call Update method on object to update record in database
		if updatedUser, err := user.Update(db, updatedUserInput); err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": errors.New("email already exists").Error()})

			} else {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

		} else {
			ctx.JSON(http.StatusOK, gin.H{"user": updatedUser})
		}

	}
}

// DeleteAccount: http handler that deletes the currently authenticated user and ALL associated data PERMENANTLY
func DeleteAccount(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// retrieve currently authenticated user instance as a Struct from the database using the request header information (jwt)
		user, err := utility.CurrentUser(ctx, db)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Delete user instance and association instances from database permanently
		if err := user.Delete(db); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		log.Println("\ndeleted account:", user.Email)

		ctx.JSON(http.StatusOK, gin.H{"deleted_user_email": user.Email})

	}
}
