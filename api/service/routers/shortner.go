package routers

import (
	"fmt"
	"go-api/database/model"
	"go-api/entities"
	"go-api/internal/auth"
	"go-api/internal/middleware"
	"go-api/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ShortenerRouter struct {
	db *gorm.DB
}

func NewShortenerRouter(db *gorm.DB) *ShortenerRouter {
	return &ShortenerRouter{db: db}
}

func (r *ShortenerRouter) RegisterBaseRoutes(router *gin.Engine) {
	router.GET("/short/:uid", r.GetShortener)
}

func (r *ShortenerRouter) RegisterRouter(router *gin.RouterGroup) {
	router.POST("/short", middleware.AuthMiddleware(), r.PostShortener)
}

func (r *ShortenerRouter) GetShortener(c *gin.Context) {
	params, ok := utils.GetParams[entities.ShortenerParams](c)
	if !ok {
		return
	}

	shortUrl, err := model.GetShortLinkByID(r.db, uint(params.UID))
	if err != nil {
		c.JSON(http.StatusBadRequest, "Invalid URL")
		return
	}

	if shortUrl == nil {
		c.JSON(http.StatusBadRequest, "Invalid URL")
		return
	}

	c.Redirect(http.StatusMovedPermanently, shortUrl.URL)
	return
}

func (r *ShortenerRouter) PostShortener(c *gin.Context) {
	body, ok := utils.GetBody[entities.ShortenerPost](c)
	if !ok {
		return
	}

	userId := auth.GetCurrentUserID(c)

	shortUrl := model.ShortLink{
		UserID: int(userId),
		URL:    body.Url,
	}

	data, err := model.CreateShortLink(r.db, &shortUrl)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Invalid URL or url already exists")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"longUrl":  data.URL,
		"shortUrl": fmt.Sprintf("%s://%s/short/%d", utils.GetProtocol(c), c.Request.Host, data.ID),
	})
}
