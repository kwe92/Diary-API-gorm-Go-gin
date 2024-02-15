package router

import (
	"journal_api/database"
	"journal_api/handler"
	"journal_api/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	publicRoutes := router.Group("/auth")

	publicRoutes.POST("/available-email", handler.CheckAvailableEmail)

	publicRoutes.POST("/register", handler.Register)

	publicRoutes.POST("/login", handler.Login)

	privateRoutes := router.Group("/api")

	// add middleware to group
	privateRoutes.Use(middleware.JWTAuthMiddleware())

	// Journal Entry Endpoints

	privateRoutes.POST("/entry", handler.AddEntry(database.Database))

	privateRoutes.POST("/update-entry/:id", handler.UpdateEntry(database.Database))

	privateRoutes.DELETE("/delete-entry/:id", handler.DeleteEntry(database.Database))

	privateRoutes.GET("/entry", handler.GetAllEntries(database.Database))

	// User Account Endpoints

	privateRoutes.DELETE("/delete-account", handler.DeleteAccount(database.Database))

	privateRoutes.POST("/update-user-info", handler.UpdateUser(database.Database))

	// Liked Quotes Endpoints

	privateRoutes.POST("/liked-quotes", handler.AddQuote(database.Database))

	privateRoutes.DELETE("/delete-liked-quote/:id", handler.DeleteQuote(database.Database))

	privateRoutes.GET("/liked-quotes", handler.GetAllLikedQuotes(database.Database))

	return router
}

// Gin Web Framework

//   - Build API's simply and quickly
//   - easy request parsing and validation
//   - provides the ability to group endpoints
//   - provides the ability to add middleware to groups of endpoints
