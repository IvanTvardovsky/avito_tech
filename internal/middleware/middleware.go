package middleware

import (
	"avito_tech/internal/tkn"
	"github.com/gin-gonic/gin"
	"net/http"
)

func TokenTypeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("token")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"description": "Пользователь не авторизован"})
			c.Abort()
			return
		}

		tokenType := "none"
		if tkn.IsAdminToken(token) {
			tokenType = "admin"
		}

		if tkn.IsUserToken(token) {
			tokenType = "user"
		}

		if tokenType == "none" {
			c.JSON(http.StatusUnauthorized, gin.H{"description": "Пользователь не авторизован"})
			c.Abort()
			return
		}

		c.Set("token_type", tokenType)
		c.Next()
	}
}

func AuthAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if val, _ := c.Get("token_type"); val != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"description": "Пользователь не имеет доступа"})
			c.Abort()
			return
		}
		c.Next()
	}
}
