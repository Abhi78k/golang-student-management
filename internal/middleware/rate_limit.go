package middleware

import (
	"net/http"
	"project-1/internal/cache"

	"github.com/gin-gonic/gin"
)

func RateLimitMiddleware(
	limiter *cache.RateLimiter,
) gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx := c.Request.Context()

		userID, exists := c.Get("UserID")

		if !exists {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "user not found.",
			})
			c.Abort()
			return
		}

		allowed, err := limiter.Allow(
			ctx,
			userID.(string),
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			c.Abort()
			return
		}

		if !allowed {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
