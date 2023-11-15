package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// database struct defined in gorm
var Database *gorm.DB

func Connect() {

	start := time.Now()

	var err error

	// retrieve Environment Variables
	host := os.Getenv("DB_HOST")
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	databaseName := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	// Open Database Connectivity Data Source Name | ODBC dsn
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=US/Central", host, username, password, databaseName, port)

	// open connection to postgres database
	Database, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}
	elapased := time.Since(start)

	log.Printf("Successfully connected to database | connection time taken: %v", elapased)

}

// ORM: Object Relational Mapper

//   - convert objects in different type systems
//   - map application domain model objects to relational database tables
//   - convert data between incompatable type systems
//   - removes the need to write many DDL and DML SQL statements
//   - objects are inserted as rows in relation database tables
//   - rows of a table can also be converted to objects

// gorm

//   - Object Relational Mapping Library for GO
//   - comes with database drivers
//   - a database you can interact with as a struct
//   - mapping of GO structs and maps to relational database tables
