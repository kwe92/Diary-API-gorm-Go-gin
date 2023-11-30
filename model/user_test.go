package model

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// TODO: Review sqlmock
// TODO: figureout how to test inserts
// TODO: all methods that interact with the database should take the database as an argument
func TestSaveUser(t *testing.T) {

	mockDB, _, _ := sqlmock.New()

	defer mockDB.Close()

	dialector := postgres.New(postgres.Config{
		Conn:       mockDB,
		DriverName: "postgres",
	})

	user := User{
		Email:    "kwe92@gmail.com",
		Password: "Ronin",
	}

	db, _ := gorm.Open(dialector, &gorm.Config{})

	user.Save(db)

}
