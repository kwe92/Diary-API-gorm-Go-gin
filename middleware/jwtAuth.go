package middleware

import (
	"diary_api/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware: validate request JWT before the request is handled.
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		err := helper.ValidateJWT(context)

		if err != nil {
			context.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication Required"})
			// prevent pending handlers from being called | useful for authentication middleware
			context.Abort()
		}
		// used inside middleware to execute pending handlers
		context.Next()
	}
}

// Middleware For Authenticated Endpoints

//   - some endpoints require a user to be authenticated
//   - requests to authenticated endpoints required a bearer token in the request Authorization Header
//   - if no bearer token is found then an error should be returned in the response body to the client

// What Authentication Middleware Should Do

//   - intercept a request before the request reaches the associated router handlers
//    - ensure the presence of a valid bearer token in Authorization Header
