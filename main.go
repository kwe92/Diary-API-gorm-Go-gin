package main

import (
	"diary_api/controller"
	"diary_api/database"
	"diary_api/middleware"
	"diary_api/model"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	loadEnv()
	loadDatabase()

	router := SetupRouter()

	router.Run(":8000")

	fmt.Println("Server running on port 8000")

}

// loadEnv: Loads environment variables.
func loadEnv() {
	err := godotenv.Load(".env.local")

	if err != nil {
		log.Fatal("error loading .env.local:", err)
	}
}

// loadDatabase: loads the postgres database associated with the environment variables.
func loadDatabase() {

	database.Connect()

	err := database.Database.AutoMigrate(&model.User{})

	checkErr(err)

	err = database.Database.AutoMigrate(&model.Entry{})

	checkErr(err)

}

func SetupRouter() *gin.Engine {
	router := gin.Default()

	publicRoutes := router.Group("/auth")

	publicRoutes.POST("/register", controller.Register)

	publicRoutes.POST("/login", controller.Login)

	privateRoutes := router.Group("/api")

	// add middleware to group
	privateRoutes.Use(middleware.JWTAuthMiddleware())

	privateRoutes.POST("/entry", controller.AddEntry)

	privateRoutes.GET("/entry", controller.GetAllEntries)

	return router
}

func checkErr(err error) {

	if err != nil {
		log.Fatal(err)
	}
}

// *gorm.DB.AutoMigrate(pointer_to_a_struct_that_will_be_a_table)

//   - automatic schema migration for a given struct
//   - will create the table and column names if they don't exist

// Loading Environment Variables

//   - environment variables must be set or loaded in order to be used

// godotenv.Load(.env_file)

//   - used to load your environment variables
//   - should be called at the top of main

// Gin Web Framework

//   - Build API's simply and quickly
//   - easy request parsing and validation
//   - provides the ability to group endpoints
//   - provides the ability to add middleware to groups of endpoints
