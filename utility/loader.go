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

// *gorm.DB.AutoMigrate(pointer_to_a_struct_that_will_be_a_table)

//   - automatic schema migration for a given struct
//   - will create the table and column names if they don't exist

// Loading Environment Variables

//   - environment variables must be set or loaded in order to be used

// godotenv.Load(.env_file)

//   - used to load your environment variables
//   - should be called at the top of main
