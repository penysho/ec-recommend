package router

import (
	"ec-recommend/internal/handler"
	"ec-recommend/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter configures and returns the HTTP router
func SetupRouter(aiHandler *handler.AIHandler, healthHandler *handler.HealthHandler) *gin.Engine {
	// Set Gin mode based on environment
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// Add middlewares
	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())
	router.Use(middleware.CORS())

	// Health check endpoint
	router.GET("/health", healthHandler.Check)

	// API routes
	v1 := router.Group("/api/v1")
	{
		// AI endpoints
		ai := v1.Group("/ai")
		{
			ai.POST("/ask", aiHandler.AskQuestion)
			ai.POST("/chat", aiHandler.Chat)
		}
	}

	return router
}
