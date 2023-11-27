package model

import (
	"time"

	"gorm.io/gorm"
)

type Entry struct {
	ID        uint           `gorm:"primarykey" json:"-"`
	UserID    uint           `json:"user_id"`
	Content   string         `gorm:"type:text" json:"content" binding:"required"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"  json:"-"`
}

// Save: insert entry into database.
func (entry *Entry) Save(db *gorm.DB) (*Entry, error) {

	err := db.Create(&entry).Error

	if err != nil {
		return &Entry{}, err
	}
	return entry, nil
}
