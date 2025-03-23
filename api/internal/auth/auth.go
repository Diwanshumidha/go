package auth

import (
	"go-api/database/model"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	UserIdKey = "userID"
)

func GetCurrentUser(c *gin.Context, db *gorm.DB) (*model.User, bool) {
	userId := c.GetInt(UserIdKey)

	if userId == 0 {
		return nil, false
	}

	user, err := model.GetUserByID(db, uint(userId))
	if err != nil {
		log.Fatal(err)
		return nil, false
	}

	return user, true
}

func GetCurrentUserID(c *gin.Context) uint {
	return uint(c.GetInt(UserIdKey))
}
