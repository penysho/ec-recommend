package dto

import "time"

// QuestionResponse represents the response structure for question answers
type QuestionResponse struct {
	Answer    string    `json:"answer"`
	Usage     UsageInfo `json:"usage,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

// ChatResponse represents the response structure for chat conversations
type ChatResponse struct {
	Message   string    `json:"message"`
	Usage     UsageInfo `json:"usage,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

// UsageInfo represents token usage information
type UsageInfo struct {
	InputTokens  int `json:"input_tokens,omitempty"`
	OutputTokens int `json:"output_tokens,omitempty"`
}

// ErrorResponse represents error response structure
type ErrorResponse struct {
	Error     string    `json:"error"`
	Code      int       `json:"code"`
	Timestamp time.Time `json:"timestamp"`
}

// HealthResponse represents health check response
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version,omitempty"`
}
