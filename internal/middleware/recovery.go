package middleware

import (
	"net/http"
	"time"

	"ec-recommend/internal/dto"

	"github.com/gin-gonic/gin"
)

// Recovery returns a panic recovery middleware
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:     "Internal server error: " + err,
				Code:      http.StatusInternalServerError,
				Timestamp: time.Now(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:     "Internal server error",
				Code:      http.StatusInternalServerError,
				Timestamp: time.Now(),
			})
		}
		c.Abort()
	})
}
