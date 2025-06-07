package dto

// QuestionRequest represents the request structure for asking questions
type QuestionRequest struct {
	Question string `json:"question" binding:"required,min=1,max=2000" validate:"required"`
}

// ChatRequest represents the request structure for chat conversations
type ChatRequest struct {
	Messages []ChatMessage `json:"messages" binding:"required,min=1" validate:"required"`
}

// ChatMessage represents a single message in a chat conversation
type ChatMessage struct {
	Role    string `json:"role" binding:"required,oneof=user assistant system" validate:"required"`
	Content string `json:"content" binding:"required,min=1,max=2000" validate:"required"`
}
