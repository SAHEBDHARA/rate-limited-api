package handler

import (
	"net/http"

	"rate-limited-api/internal/core/domain"
	"rate-limited-api/internal/core/ports"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RequestHandler struct {
	service ports.RequestService
}

func NewRequestHandler(service ports.RequestService) *RequestHandler {
	return &RequestHandler{service: service}
}

type postRequestInput struct {
	UserID  string `json:"user_id" binding:"required"`
	Payload string `json:"payload" binding:"required"`
}

func (h *RequestHandler) HandlePostRequest(c *gin.Context) {
	var input postRequestInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create domain request
	req, err := domain.NewRequest(uuid.New().String(), input.UserID, input.Payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.ProcessRequest(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message":    "request accepted",
		"request_id": req.ID,
	})
}
