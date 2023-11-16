package middleware

import (
	"diary_api/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware: validate request JWT before request handling.
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {

		_, err := helper.ValidateJWT(context)

		if err != nil {
			context.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication Required"})

			// prevent pending handler execution
			context.Abort()
		}

		// executes pending handlers inside middleware
		context.Next()
	}
}

// Middleware For Authenticated Endpoints

//   - some endpoints require user authentication
//   - requests to authenticated endpoints handled by JWT's
//     require a bearer token in request Authorization Header
//   - in the absence of a bearer token an error should be written to the response body

// What Authentication Middleware Should Do

//   - intercept request before request reaches associated router handlers
//   - ensure the presence of a valid bearer token in Authorization Header
