package model

import (
	"diary_api/database"

	"gorm.io/gorm"
)

type Entry struct {
	gorm.Model
	Content string `gorm:"type:text" json:"content"`
	UserID  uint
}

// Save: insert entry into database.
func (entry *Entry) Save() (*Entry, error) {

	err := database.Database.Create(&entry).Error

	if err != nil {
		return &Entry{}, nil
	}
	return entry, nil
}
