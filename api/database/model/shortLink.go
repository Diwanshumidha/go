package model

import (
	"gorm.io/gorm"
)

type ShortLink struct {
	gorm.Model
	UserID int    `gorm:"type:int;index"` // Ensure UUID consistency
	URL    string `json:"url"`
}

func GetShortLinkByID(db *gorm.DB, id uint) (*ShortLink, error) {
	var shortLink ShortLink
	if err := db.First(&shortLink, "id = ?", id).Error; err != nil {
		return nil, err
	}

	// TODO: ADD CACHING

	return &shortLink, nil
}

func CreateShortLink(db *gorm.DB, shortLink *ShortLink) (*ShortLink, error) {
	if err := db.Create(shortLink).Error; err != nil {
		return nil, err
	}
	return shortLink, nil
}
