package handler

import (
	"net/http"
	"time"

	"ec-recommend/internal/dto"
	"ec-recommend/internal/service"

	"github.com/gin-gonic/gin"
)

// AIHandler handles AI-related HTTP requests
type AIHandler struct {
	bedrockService service.BedrockService
}

// NewAIHandler creates a new AI handler
func NewAIHandler(bedrockService service.BedrockService) *AIHandler {
	return &AIHandler{
		bedrockService: bedrockService,
	}
}

// AskQuestion handles question-answer requests
func (h *AIHandler) AskQuestion(c *gin.Context) {
	var req dto.QuestionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:     "Invalid request format: " + err.Error(),
			Code:      http.StatusBadRequest,
			Timestamp: time.Now(),
		})
		return
	}

	// Generate AI response
	aiResponse, err := h.bedrockService.GenerateResponse(c.Request.Context(), req.Question)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:     "Failed to generate AI response: " + err.Error(),
			Code:      http.StatusInternalServerError,
			Timestamp: time.Now(),
		})
		return
	}

	// Create response
	response := dto.QuestionResponse{
		Answer: aiResponse.Content,
		Usage: dto.UsageInfo{
			InputTokens:  aiResponse.Usage.InputTokens,
			OutputTokens: aiResponse.Usage.OutputTokens,
		},
		Timestamp: time.Now(),
	}

	c.JSON(http.StatusOK, response)
}

// Chat handles chat conversation requests
func (h *AIHandler) Chat(c *gin.Context) {
	var req dto.ChatRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:     "Invalid request format: " + err.Error(),
			Code:      http.StatusBadRequest,
			Timestamp: time.Now(),
		})
		return
	}

	// For now, use the last user message as the prompt
	// In a more advanced implementation, you would handle the full conversation context
	var lastUserMessage string
	for _, msg := range req.Messages {
		if msg.Role == "user" {
			lastUserMessage = msg.Content
		}
	}

	if lastUserMessage == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:     "No user message found in conversation",
			Code:      http.StatusBadRequest,
			Timestamp: time.Now(),
		})
		return
	}

	// Generate AI response
	aiResponse, err := h.bedrockService.GenerateResponse(c.Request.Context(), lastUserMessage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:     "Failed to generate AI response: " + err.Error(),
			Code:      http.StatusInternalServerError,
			Timestamp: time.Now(),
		})
		return
	}

	// Create response
	response := dto.ChatResponse{
		Message: aiResponse.Content,
		Usage: dto.UsageInfo{
			InputTokens:  aiResponse.Usage.InputTokens,
			OutputTokens: aiResponse.Usage.OutputTokens,
		},
		Timestamp: time.Now(),
	}

	c.JSON(http.StatusOK, response)
}
