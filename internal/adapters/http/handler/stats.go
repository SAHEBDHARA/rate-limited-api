package handler

import (
	"net/http"

	"rate-limited-api/internal/core/ports"

	"github.com/gin-gonic/gin"
)

type StatsHandler struct {
	service ports.RequestService
}

func NewStatsHandler(service ports.RequestService) *StatsHandler {
	return &StatsHandler{service: service}
}

func (h *StatsHandler) HandleGetStats(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id query parameter is required"})
		return
	}

	stats, err := h.service.GetStats(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}
