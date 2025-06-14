package handler

import (
	"context"
	"net/http"
	"time"

	"ec-recommend/internal/dto"
	"ec-recommend/internal/types"

	"github.com/gin-gonic/gin"
)

// ChatServiceInterface defines the interface for chat service operations used by handler
type ChatServiceInterface interface {
	// GenerateResponse generates a response for the given prompt
	GenerateResponse(ctx context.Context, prompt string) (*types.AIResponse, error)
}

// ChatHandler handles chat-related HTTP requests
type ChatHandler struct {
	chatService ChatServiceInterface
}

// NewChatHandler creates a new chat handler
func NewChatHandler(chatService ChatServiceInterface) *ChatHandler {
	return &ChatHandler{
		chatService: chatService,
	}
}

// AskQuestion handles question-answer requests
func (h *ChatHandler) AskQuestion(c *gin.Context) {
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
	chatResponse, err := h.chatService.GenerateResponse(c.Request.Context(), req.Question)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:     "Failed to generate chat response: " + err.Error(),
			Code:      http.StatusInternalServerError,
			Timestamp: time.Now(),
		})
		return
	}

	// Create response
	response := dto.QuestionResponse{
		Answer: chatResponse.Content,
		Usage: dto.UsageInfo{
			InputTokens:  chatResponse.Usage.InputTokens,
			OutputTokens: chatResponse.Usage.OutputTokens,
		},
		Timestamp: time.Now(),
	}

	c.JSON(http.StatusOK, response)
}

// Chat handles chat conversation requests
func (h *ChatHandler) Chat(c *gin.Context) {
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
	chatResponse, err := h.chatService.GenerateResponse(c.Request.Context(), lastUserMessage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:     "Failed to generate chat response: " + err.Error(),
			Code:      http.StatusInternalServerError,
			Timestamp: time.Now(),
		})
		return
	}

	// Create response
	response := dto.ChatResponse{
		Message: chatResponse.Content,
		Usage: dto.UsageInfo{
			InputTokens:  chatResponse.Usage.InputTokens,
			OutputTokens: chatResponse.Usage.OutputTokens,
		},
		Timestamp: time.Now(),
	}

	c.JSON(http.StatusOK, response)
}
