package middleware

import (
	"net/http"

	"rate-limited-api/internal/core/ports"

	"github.com/gin-gonic/gin"
)

func RateLimitMiddleware(limiter ports.RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Query("user_id")
		if userID == "" {
			userID = c.GetHeader("X-User-ID")
		}

		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required (use X-User-ID header or user_id query param)"})
			c.Abort()
			return
		}

		allowed, err := limiter.Allow(c.Request.Context(), userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "rate limit check failed"})
			c.Abort()
			return
		}

		if !allowed {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded. Max 5 requests per minute."})
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}
