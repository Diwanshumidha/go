package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name       string      `json:"name"`
	Email      string      `json:"email" gorm:"unique"`
	Password   string      `json:"password"`
	ShortLinks []ShortLink `gorm:"foreignKey:UserID;references:ID"`
}

func GetUserByID(db *gorm.DB, id uint) (*User, error) {
	var user User
	if err := db.First(&user, id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUserByEmail fetches a user by Email with enhanced error handling
func GetUserByEmail(db *gorm.DB, email string) (*User, error) {
	var user User
	if err := db.First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func CreateUser(db *gorm.DB, user *User) error {
	return db.Create(user).Error
}

func CreateUserAndReturnID(db *gorm.DB, user *User) (uint, error) {
	if err := db.Create(user).Error; err != nil {
		return 0, err
	}
	return user.ID, nil
}
