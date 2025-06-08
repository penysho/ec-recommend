package dto

import "time"

// ChatRequest represents the request structure for chat conversations
type ChatRequest struct {
	Messages []ChatMessage `json:"messages" binding:"required,min=1" validate:"required"`
}

// ChatMessage represents a single message in a chat conversation
type ChatMessage struct {
	Role    string `json:"role" binding:"required,oneof=user assistant system" validate:"required"`
	Content string `json:"content" binding:"required,min=1,max=2000" validate:"required"`
}

// QuestionRequest represents the request structure for asking questions
type QuestionRequest struct {
	Question string `json:"question" binding:"required,min=1,max=2000" validate:"required"`
}

// ChatResponse represents the response structure for chat conversations
type ChatResponse struct {
	Message   string    `json:"message"`
	Usage     UsageInfo `json:"usage,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

// QuestionResponse represents the response structure for question answers
type QuestionResponse struct {
	Answer    string    `json:"answer"`
	Usage     UsageInfo `json:"usage,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}
