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
			utility.SendBadRequestResponse(ctx, err)
			return
		}

		// retrieve currently authenticated user instance as a Struct from the database using the request header information (jwt)
		user, err := utility.CurrentUser(ctx, db, false)

		if err != nil {
			utility.SendBadRequestResponse(ctx, err)
			return
		}

		// call Update method on object to update record in database
		if updatedUser, err := user.Update(db, updatedUserInput); err != nil {
			// if the updated email already exists respond with an error message
			if strings.Contains(err.Error(), "duplicate key") {
				utility.SendBadRequestResponse(ctx, errors.New("email already exists"))
			} else {
				// if any other error occurs respond with the error message received
				utility.SendBadRequestResponse(ctx, err)
				return
			}

		} else {
			// if the update was successful send updated user info

			ctx.JSON(http.StatusOK, gin.H{"user": map[string]any{
				"first_name":   updatedUser.Fname,
				"last_name":    updatedUser.Lname,
				"email":        updatedUser.Email,
				"phone_number": updatedUser.Phone,
			},
			})

			// ctx.JSON(http.StatusOK, gin.H{"user": updatedUser})

			// first_name
			// last_name
			// email
			// phone_number
		}

	}
}

// DeleteAccount: http handler that deletes the currently authenticated user and ALL associated data PERMENANTLY
func DeleteAccount(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// retrieve currently authenticated user instance as a Struct from the database using the request header information (jwt)
		user, err := utility.CurrentUser(ctx, db, false)

		if err != nil {
			utility.SendBadRequestResponse(ctx, err)
			return
		}

		// Delete user instance and association instances from database permanently
		if err := user.Delete(db); err != nil {
			utility.SendBadRequestResponse(ctx, err)
			return
		}

		log.Println("\ndeleted account:", user.Email)

		ctx.JSON(http.StatusOK, gin.H{"deleted_user_email": user.Email})

	}
}
