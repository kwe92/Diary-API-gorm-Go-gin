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

	privateRoutes.POST("/entry", handler.AddEntry(database.Database))

	privateRoutes.POST("/update-entry/:id", handler.UpdateEntry(database.Database))

	privateRoutes.DELETE("/delete-entry/:id", handler.DeleteEntry(database.Database))

	privateRoutes.GET("/entry", handler.GetAllEntries(database.Database))

	privateRoutes.DELETE("/delete-account", handler.DeleteAccount(database.Database))

	return router
}
