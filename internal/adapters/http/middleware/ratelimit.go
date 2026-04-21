package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"rate-limited-api/internal/core/ports"
)

// RateLimitMiddleware returns a Gin middleware that checks the rate limit for each user.
func RateLimitMiddleware(limiter ports.RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Query("user_id")
		if userID == "" {
			// In a real app, we might check a Header or JWT claim
			// For this assignment, we'll try to get it from query or JSON body in the handler
			// But for the middleware to be effective, let's look for it in a header or query
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

		// Store userID in context for subsequent handlers
		c.Set("user_id", userID)
		c.Next()
	}
}
