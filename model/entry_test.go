package model

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestSaveEntry(t *testing.T) {

	mockDB, _, err := sqlmock.New()

	defer mockDB.Close()

	assert.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		Conn:       mockDB,
		DriverName: "postgres",
	})

	db, err := gorm.Open(dialector, &gorm.Config{})

	assert.NoError(t, err)

	entry := Entry{
		Content: "Begin, to begin is half the work let half still remain again beginthis and thou wilt have finished.",
		UserID:  532,
	}

	entry.Save(db)

	// mock.ExpectBegin()

	// rows := sqlmock.NewRows([]string{"id", "created_id", "updated_at", "deleted_at", "content", "user_id"}).AddRow(1, "2022-09-010 01:01:00", "2022-09-08 01:01:01", "2022-09-08 01:01:02", "emptiness in calmness, calmness in emptiness.", "1001")

	// mockedTransaction := mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `entries` WHERE id = ?")).WithArgs(1).WillReturnRows(rows)

	// fmt.Println(mockedTransaction)

	// mock.ExpectExec()

	// mock.ExpectCommit()
}
