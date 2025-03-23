package routers

import (
	"go-api/entities"
	"go-api/internal/middleware"
	"go-api/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthRouter struct{}

func NewHealthRouter() *HealthRouter {
	return &HealthRouter{}
}

func (r *HealthRouter) RegisterRouter(router *gin.RouterGroup) {
	router.GET("/ping", r.GetHealth)
	router.POST("/ping", r.PostHealth)
	router.GET("/ping/:quantity", middleware.AuthMiddleware(), r.GetHealthWithParams)
}

// Handler function that retrieves and uses validated data
func (r *HealthRouter) GetHealth(c *gin.Context) {
	// Use the validated data
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func (r *HealthRouter) PostHealth(c *gin.Context) {
	body, ok := utils.GetBody[entities.HealthPost](c)
	if !ok {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": body.Message,
	})
}

func (r *HealthRouter) GetHealthWithParams(c *gin.Context) {
	params, ok := utils.GetParams[entities.HealthParams](c)
	if !ok {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": params.Quantity,
	})
}
