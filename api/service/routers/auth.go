package routers

import (
	"go-api/database/model"
	"go-api/entities"
	"go-api/internal/env"
	"go-api/internal/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthRouter struct {
	db *gorm.DB
}

func NewAuthRouter(db *gorm.DB) *AuthRouter {
	return &AuthRouter{db: db}
}

func (r *AuthRouter) RegisterRouter(router *gin.RouterGroup) {
	authRouter := router.Group("/auth")
	{
		authRouter.POST("/register", r.RegisterAccount)
		authRouter.POST("/login", r.LoginAccount)
	}
}

func (r *AuthRouter) RegisterAccount(c *gin.Context) {
	body, ok := utils.GetBody[entities.AuthRegisterRequestBody](c)
	if !ok {
		return
	}

	existingAccount, err := model.GetUserByEmail(r.db, body.Email)
	if err == nil || existingAccount != nil {
		c.JSON(http.StatusBadRequest, "Account already exists")
		return
	}

	encryptedPassword, err := utils.HashPassword(body.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Something went wrong.")
		return
	}

	user := model.User{
		Name:       body.Email,
		Email:      body.Email,
		Password:   encryptedPassword,
		ShortLinks: []model.ShortLink{},
	}

	usrId, err := model.CreateUserAndReturnID(r.db, &user)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Account Already Exists")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"userId": usrId,
	})

	return
}

func (r *AuthRouter) LoginAccount(c *gin.Context) {
	body, ok := utils.GetBody[entities.AuthLoginRequestBody](c)
	if !ok {
		return
	}

	user, err := model.GetUserByEmail(r.db, body.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Invalid Credentials")
		return
	}

	if !utils.CheckPassword(user.Password, body.Password) {
		c.JSON(http.StatusUnauthorized, "Invalid Credentials")
		return
	}

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Something went wrong.")
		return
	}

	tokenExpiration := time.Duration(env.GetInt("JWT_TOKEN_EXPIRATION", int(time.Hour*10))) // Default to 10 hours

	c.SetCookie("token", token, int(tokenExpiration.Seconds()), "/", "", false, false)
	c.JSON(http.StatusOK, gin.H{
		"userId": user.ID,
	})
}
