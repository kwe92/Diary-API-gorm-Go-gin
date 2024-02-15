package model

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// TODO: Finish implementation | delete method

type LikedQuote struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	UserID    uint           `json:"user_id"`
	Author    string         `gorm:"type:text" json:"author" binding:"required"`
	Quote     string         `gorm:"type:text" json:"quote" binding:"required"`
	IsLiked   bool           `gorm:"type:bool" json:"is_liked" binding:"required"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"  json:"-"`
}

func (quote *LikedQuote) Save(db *gorm.DB) (*LikedQuote, error) {

	// insert liked quote into database if there are no errors are encountered
	err := db.Create(&quote).Error

	if err != nil {

		return &LikedQuote{}, err
	}
	return quote, nil
}

// Delete: PERMANENTLY delete quote from database.
func (quote *LikedQuote) Delete(db *gorm.DB) (*LikedQuote, error) {

	// select quote by id then delete permanently
	if err := db.Unscoped().Delete(&quote).Error; err != nil {
		return &LikedQuote{}, err
	}

	return quote, nil
}

// FindQuoteById: query database for quote record with given id
func FindQuoteById(db *gorm.DB, quoteId uint) (LikedQuote, error) {

	// declare destination struct
	var quote LikedQuote

	// define error message if quote not found
	notFound := "not_founds"

	// query with format string in where clause | initialize LikedQoute Quote field with not found if no record is found
	if err := db.Limit(1).Where("id = ?", quoteId).Attrs(LikedQuote{Quote: notFound}).FirstOrInit(&quote).Error; err != nil {
		return LikedQuote{}, err
	}

	// if destinaion struct has not found for quote content return not found error to caller
	if quote.Quote == notFound {
		return LikedQuote{}, errors.New(fmt.Sprintf("could not find a quote with the id: %d", quoteId))
	}

	return quote, nil

}
