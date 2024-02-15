package utility

import (
	"journal_api/database"
	"journal_api/model"
	"log"

	"github.com/joho/godotenv"
)

// loadEnv: Loads environment variables.
func LoadEnv() {
	err := godotenv.Load(".env.local")

	if err != nil {
		log.Fatal("error loading .env.local:", err)
	}
}

// loadDatabase: loads the postgres database associated with the environment variables.
func LoadDatabase() {

	database.Connect()

	err := database.Database.AutoMigrate(&model.User{})

	checkErr(err)

	err = database.Database.AutoMigrate(&model.Entry{})

	checkErr(err)

	err = database.Database.AutoMigrate(&model.LikedQuote{})

	checkErr(err)

}

func checkErr(err error) {

	if err != nil {
		log.Fatal(err)
	}
}
