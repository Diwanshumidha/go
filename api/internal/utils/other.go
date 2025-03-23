package utils

import "github.com/gin-gonic/gin"

func GetProtocol(c *gin.Context) string {
	if c.Request.TLS != nil {
		return "https"
	}
	return "http"
}
