package model

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Entry struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	UserID    uint           `json:"user_id"`
	Content   string         `gorm:"type:text" json:"content" binding:"required"`
	MoodType  string         `gorm:"type:text" json:"mood_type" binding:"required"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"  json:"-"`
}

type UpdatedEntryInput struct {
	ID      uint   `gorm:"primarykey" json:"id" binding:"required"`
	Content string `gorm:"type:text" json:"content" binding:"required"`
}

// Save: insert entry into database.
func (entry *Entry) Save(db *gorm.DB) (*Entry, error) {

	err := db.Create(&entry).Error

	if err != nil {
		return &Entry{}, err
	}
	return entry, nil
}

// Update: updates data for the associated entry record in database
func (entry *Entry) Update(db *gorm.DB, input UpdatedEntryInput) (*Entry, error) {

	originalEntry := *entry

	var result *gorm.DB

	// specify model you want to perfom operations on and update record
	if result = db.Model(entry).Updates(input); result.Error != nil {

		return &Entry{}, result.Error

	}

	log.Printf(
		"\nupdated entry from: %+v\n\nto: %+v\n\n",
		gin.H{
			"id":      originalEntry.ID,
			"content": originalEntry.Content,
		},
		gin.H{
			"id":      entry.ID,
			"content": entry.Content,
		},
	)

	log.Println("rows affected:", result.RowsAffected)

	return entry, nil

}

// Delete: remove entry record from database
func (entry *Entry) Delete(db *gorm.DB) (*Entry, error) {

	// ? NOTE: seems to mark as deleted but remains in the database but not retrieved if DeleteAt is not null
	if err := db.Delete(&entry).Error; err != nil {

		return &Entry{}, err
	}

	return entry, nil
}

// FindEntryById: query database for entry record with given id
func FindEntryById(db *gorm.DB, entryId uint) (Entry, error) {

	// destination struct pointer
	var entry Entry

	notFound := "not_found"

	// query with format string in where clause | initialize with not found if no record is found
	if err := db.Limit(1).Where("id = ?", entryId).Attrs(Entry{Content: notFound}).FirstOrInit(&entry).Error; err != nil {

		return Entry{}, err

	}
	// if destinaion struct has zero-value for id return not found to client
	if entry.Content == notFound {
		return Entry{}, errors.New(fmt.Sprintf("could not find a entry with the id: %d", entryId))
	}

	return entry, nil
}
