package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		admToken := c.GetHeader("token")
		if admToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"description": "Пользователь не авторизован"})
			c.Abort()
			return
		}

		if admToken != "admin_token" { // example
			c.JSON(http.StatusUnauthorized, gin.H{"description": "Пользователь не имеет доступа"})
			c.Abort()
			return
		}

		c.Next()
	}
}
