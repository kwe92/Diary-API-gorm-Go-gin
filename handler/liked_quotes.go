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

// AddQuote: route handler to add a new liked quote of an authorized user
func AddQuote(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// declare expected requrest body input
		var quote model.LikedQuote

		// unmarshel request body into struct
		err := ctx.ShouldBindJSON(&quote)

		if err != nil {

			utility.SendBadRequestResponse(ctx, err)

			return
		}

		// retrieve currently authenticated user from request header
		user, err := utility.CurrentUser(ctx, db, false)

		if err != nil {

			utility.SendBadRequestResponse(ctx, err)

			return
		}

		// map currently authenticated user id to quote
		quote.UserID = user.ID

		// insert quote into database
		savedQuote, err := quote.Save(db)

		if err != nil {
			utility.SendBadRequestResponse(ctx, err)
			return
		}
		// write liked quote to response body
		ctx.JSON(http.StatusOK, gin.H{"success": savedQuote.Quote})

	}
}

// DeleteQuote: http handler that deletes a single quote in the database
func DeleteQuote(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// parse quote id from URL parameters
		quoteId, err := strconv.Atoi(ctx.Param("id"))

		// check for error after parsing quote id

		if err != nil {
			utility.SendBadRequestResponse(ctx, err)
			return
		}

		//  find record to delete and load record into destination struct

		if quote, err := model.FindQuoteById(db, uint(quoteId)); err != nil {
			utility.SendBadRequestResponse(ctx, err)
			return
		} else {
			if deletedQuote, err := quote.Delete(db); err != nil {

				utility.SendBadRequestResponse(ctx, err)
				return

			} else {
				log.Println("\ndeleted quote:", *deletedQuote)

				// return deleted object with DeletedAt time
				ctx.JSON(http.StatusOK, gin.H{"success": gin.H{
					"quote":      deletedQuote.Quote,
					"deleted_at": deletedQuote.DeletedAt,
				}})
			}

		}

	}
}

// GetAllLikedQuotes: route handler that retrieves current user and returns all associated liked quotes as a response
func GetAllLikedQuotes(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// retrieve current authenticated user from request header
		user, err := utility.CurrentUser(ctx, db, true)

		if err != nil {
			utility.SendBadRequestResponse(ctx, err)
			return
		}

		log.Println(user.LikedQuotes)

		// write current user liked quotes to response body
		ctx.JSON(http.StatusOK, gin.H{"data": user.LikedQuotes})
	}
}
