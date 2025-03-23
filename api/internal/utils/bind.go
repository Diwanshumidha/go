package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetBody extracts and validates the request body into a struct T
func GetBody[T any](c *gin.Context) (T, bool) {
	var body T
	if err := c.ShouldBindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return body, false
	}
	return body, true
}

// GetParams extracts and validates URL parameters into a struct T
func GetParams[T any](c *gin.Context) (T, bool) {
	var params T
	if err := c.ShouldBindUri(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid path parameters: " + err.Error()})
		return params, false
	}
	return params, true
}

// GetSearchParams extracts and validates query parameters into a struct T
func GetSearchParams[T any](c *gin.Context) (T, bool) {
	var queryParams T
	if err := c.ShouldBindQuery(&queryParams); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters: " + err.Error()})
		return queryParams, false
	}
	return queryParams, true
}
