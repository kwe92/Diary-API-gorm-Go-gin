package middleware

import (
	"journal_api/utility"
	"net/http"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware: validate request JWT before request handling.
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		_, err := utility.ValidateJWT(ctx)

		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication Required"})

			// prevent pending handler execution
			ctx.Abort()
		}

		// execute pending handlers inside middleware
		ctx.Next()
	}
}

// Middleware For Authenticated Endpoints

//   - some endpoints require user authentication
//   - requests to authenticated endpoints are handled by JWT's
//     require a bearer token in request Authorization Header
//   - in the absence of a bearer token an error should be written to the response body

// What Authentication Middleware Should Do

//   - intercept request before request reaches associated router handlers
//   - ensure the presence of a valid bearer token in Authorization Header
