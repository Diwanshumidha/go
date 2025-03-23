package middleware

import (
	"go-api/internal/auth"
	"go-api/internal/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("token")
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{
				"message": "Unauthorized",
			})
			c.Abort()
			return
		}

		if token == "" {
			c.AbortWithStatusJSON(401, gin.H{
				"message": "Unauthorized",
			})
			c.Abort()
			return
		}

		tokenClaims, err := utils.ValidateJWT(token)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{
				"message": "Unauthorized",
			})
			c.Abort()
			return
		}

		if tokenClaims.UserID == 0 {
			c.AbortWithStatusJSON(401, gin.H{
				"message": "Unauthorized",
			})
			c.Abort()
			return
		}

		println(tokenClaims.UserID)
		c.Set(auth.UserIdKey, tokenClaims.UserID)
		c.Next()
	}
}
