package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

// BedrockService defines the interface for Bedrock AI operations
type BedrockService interface {
	GenerateResponse(ctx context.Context, prompt string) (*AIResponse, error)
}

// BedrockClient wraps the AWS Bedrock runtime client
type BedrockClient struct {
	client  *bedrockruntime.Client
	modelID string
}

// AIResponse represents the response from the AI model
type AIResponse struct {
	Content string `json:"content"`
	Usage   Usage  `json:"usage,omitempty"`
}

// Usage represents token usage information
type Usage struct {
	InputTokens  int `json:"input_tokens,omitempty"`
	OutputTokens int `json:"output_tokens,omitempty"`
}

// ClaudeRequest represents the request structure for Claude models
type ClaudeRequest struct {
	Messages         []Message `json:"messages"`
	MaxTokens        int       `json:"max_tokens"`
	Temperature      float64   `json:"temperature,omitempty"`
	TopP             float64   `json:"top_p,omitempty"`
	AnthropicVersion string    `json:"anthropic_version"`
}

// Message represents a message in the conversation
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// NovaMessage represents a message for Nova models
type NovaMessage struct {
	Role    string        `json:"role"`
	Content []NovaContent `json:"content"`
}

// NovaContent represents content structure for Nova
type NovaContent struct {
	Text string `json:"text"`
}

// ClaudeResponse represents the response structure from Claude models
type ClaudeResponse struct {
	Content []ContentBlock `json:"content"`
	Usage   Usage          `json:"usage"`
}

// ContentBlock represents a content block in the response
type ContentBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// TitanRequest represents the request structure for Titan models
type TitanRequest struct {
	InputText            string             `json:"inputText"`
	TextGenerationConfig TitanTextGenConfig `json:"textGenerationConfig"`
}

// TitanTextGenConfig represents text generation configuration for Titan
type TitanTextGenConfig struct {
	MaxTokenCount int      `json:"maxTokenCount"`
	Temperature   float64  `json:"temperature,omitempty"`
	TopP          float64  `json:"topP,omitempty"`
	StopSequences []string `json:"stopSequences,omitempty"`
}

// TitanResponse represents the response structure from Titan models
type TitanResponse struct {
	Results             []TitanResult `json:"results"`
	InputTextTokenCount int           `json:"inputTextTokenCount"`
}

// TitanResult represents a result from Titan
type TitanResult struct {
	TokenCount       int    `json:"tokenCount"`
	OutputText       string `json:"outputText"`
	CompletionReason string `json:"completionReason"`
}

// TitanUsage represents usage info from Titan
type TitanUsage struct {
	InputTextTokenCount int `json:"inputTextTokenCount"`
	TotalTokenCount     int `json:"totalTokenCount"`
}

// NovaRequest represents the request structure for Nova models
type NovaRequest struct {
	Messages        []NovaMessage       `json:"messages"`
	InferenceConfig NovaInferenceConfig `json:"inferenceConfig,omitempty"`
}

// NovaInferenceConfig represents inference configuration for Nova
type NovaInferenceConfig struct {
	MaxTokens int `json:"max_new_tokens"`
	// Temperature float64 `json:"temperature,omitempty"`
	// TopP        float64 `json:"top_p,omitempty"`
}

// NovaResponse represents the response structure from Nova models
type NovaResponse struct {
	Output NovaOutput `json:"output"`
	Usage  NovaUsage  `json:"usage"`
}

// NovaOutput represents output from Nova
type NovaOutput struct {
	Message NovaMessage `json:"message"`
}

// NovaUsage represents usage info from Nova
type NovaUsage struct {
	InputTokens  int `json:"inputTokens"`
	OutputTokens int `json:"outputTokens"`
	TotalTokens  int `json:"totalTokens"`
}

// NewBedrockClient creates a new Bedrock client wrapper
func NewBedrockClient(client *bedrockruntime.Client, modelID string) *BedrockClient {
	return &BedrockClient{
		client:  client,
		modelID: modelID,
	}
}

// GenerateResponse generates a response using the specified model
func (bc *BedrockClient) GenerateResponse(ctx context.Context, prompt string) (*AIResponse, error) {
	if prompt == "" {
		return nil, fmt.Errorf("prompt cannot be empty")
	}

	// Prepare request based on model type
	requestBody, err := bc.prepareRequest(prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare request: %w", err)
	}

	// Invoke the model
	input := &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(bc.modelID),
		ContentType: aws.String("application/json"),
		Accept:      aws.String("application/json"),
		Body:        requestBody,
	}

	result, err := bc.client.InvokeModel(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke model: %w", err)
	}

	// Parse response
	response, err := bc.parseResponse(result.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return response, nil
}

// prepareRequest prepares the request body based on the model type
func (bc *BedrockClient) prepareRequest(prompt string) ([]byte, error) {
	// For Claude models
	if bc.isClaudeModel() {
		request := ClaudeRequest{
			Messages: []Message{
				{
					Role:    "user",
					Content: prompt,
				},
			},
			MaxTokens:        4096,
			Temperature:      0.7,
			AnthropicVersion: "bedrock-2023-05-31",
		}
		return json.Marshal(request)
	}

	// For Titan models
	if bc.isTitanModel() {
		request := TitanRequest{
			InputText: prompt,
			TextGenerationConfig: TitanTextGenConfig{
				MaxTokenCount: 4096,
				Temperature:   0.7,
				TopP:          0.9,
				StopSequences: []string{},
			},
		}
		return json.Marshal(request)
	}

	// For Nova models
	if bc.isNovaModel() {
		request := NovaRequest{
			Messages: []NovaMessage{
				{
					Role: "user",
					Content: []NovaContent{
						{
							Text: prompt,
						},
					},
				},
			},
			InferenceConfig: NovaInferenceConfig{
				MaxTokens: 4096,
			},
		}
		return json.Marshal(request)
	}

	return nil, fmt.Errorf("unsupported model: %s", bc.modelID)
}

// parseResponse parses the response body based on the model type
func (bc *BedrockClient) parseResponse(body []byte) (*AIResponse, error) {
	// For Claude models
	if bc.isClaudeModel() {
		var claudeResp ClaudeResponse
		if err := json.Unmarshal(body, &claudeResp); err != nil {
			return nil, fmt.Errorf("failed to unmarshal Claude response: %w", err)
		}

		if len(claudeResp.Content) == 0 {
			return nil, fmt.Errorf("no content in response")
		}

		return &AIResponse{
			Content: claudeResp.Content[0].Text,
			Usage:   claudeResp.Usage,
		}, nil
	}

	// For Titan models
	if bc.isTitanModel() {
		var titanResp TitanResponse
		if err := json.Unmarshal(body, &titanResp); err != nil {
			return nil, fmt.Errorf("failed to unmarshal Titan response: %w", err)
		}

		if len(titanResp.Results) == 0 {
			return nil, fmt.Errorf("no results in Titan response")
		}

		return &AIResponse{
			Content: titanResp.Results[0].OutputText,
			Usage: Usage{
				InputTokens:  titanResp.InputTextTokenCount,
				OutputTokens: titanResp.Results[0].TokenCount,
			},
		}, nil
	}

	// For Nova models
	if bc.isNovaModel() {
		var novaResp NovaResponse
		if err := json.Unmarshal(body, &novaResp); err != nil {
			return nil, fmt.Errorf("failed to unmarshal Nova response: %w", err)
		}

		// Extract text from content array
		var content string
		if len(novaResp.Output.Message.Content) > 0 {
			content = novaResp.Output.Message.Content[0].Text
		}

		return &AIResponse{
			Content: content,
			Usage: Usage{
				InputTokens:  novaResp.Usage.InputTokens,
				OutputTokens: novaResp.Usage.OutputTokens,
			},
		}, nil
	}

	return nil, fmt.Errorf("unsupported model for response parsing: %s", bc.modelID)
}

// isClaudeModel checks if the current model is a Claude model
func (bc *BedrockClient) isClaudeModel() bool {
	return len(bc.modelID) > 10 && bc.modelID[:10] == "anthropic."
}

// isTitanModel checks if the current model is a Titan model
func (bc *BedrockClient) isTitanModel() bool {
	if len(bc.modelID) < 7 {
		return false
	}
	return bc.modelID[:7] == "amazon." &&
		(len(bc.modelID) > 12 && bc.modelID[7:12] == "titan")
}

// isNovaModel checks if the current model is a Nova model
func (bc *BedrockClient) isNovaModel() bool {
	if len(bc.modelID) < 7 {
		return false
	}
	return bc.modelID[:7] == "amazon." &&
		(len(bc.modelID) > 11 && bc.modelID[7:11] == "nova")
}
