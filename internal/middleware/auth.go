package middleware

import (
	"net/http"
	"project-1/internal/auth"
	"project-1/internal/config"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(
	cfg *config.Config,
) gin.HandlerFunc {
	return func(c *gin.Context) {

		header := c.GetHeader("Authorization")

		if header == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "authorization header required.",
			})
			c.Abort()
			return
		}

		parts := strings.Split(header, " ")

		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization format.",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		token, err := auth.ValidateToken(tokenString, cfg)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token.",
			})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*auth.Claims)

		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token claims.",
			})
			c.Abort()
			return
		}

		c.Set("UserID", claims.UserID)

		c.Next()
	}
}
