package handler

import (
	"net/http"
	"time"

	"ec-recommend/internal/dto"

	"github.com/gin-gonic/gin"
)

// HealthHandler handles health check requests
type HealthHandler struct{}

// NewHealthHandler creates a new health handler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Check handles health check requests
func (h *HealthHandler) Check(c *gin.Context) {
	response := dto.HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Version:   "1.0.0",
	}

	c.JSON(http.StatusOK, response)
}
