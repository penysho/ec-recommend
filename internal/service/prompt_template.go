package service

import (
	"context"
	"ec-recommend/internal/dto"
)

// PromptTemplate represents a template for generating prompts
type PromptTemplate struct {
	ID            string        `json:"id"`
	Name          string        `json:"name"`
	BasePrompt    string        `json:"base_prompt"`
	Examples      []ExampleCase `json:"examples"`
	ContextRules  []ContextRule `json:"context_rules"`
	OutputSchema  OutputSchema  `json:"output_schema"`
	Version       string        `json:"version"`
	Category      string        `json:"category"`
	Description   string        `json:"description"`
}

// ExampleCase represents a few-shot learning example
type ExampleCase struct {
	Input           string `json:"input"`
	ExpectedOutput  string `json:"expected_output"`
	Explanation     string `json:"explanation"`
	Context         string `json:"context"`
	Priority        int    `json:"priority"`
}

// ContextRule defines rules for modifying prompts based on context
type ContextRule struct {
	Condition    string `json:"condition"`
	Modification string `json:"modification"`
	Priority     int    `json:"priority"`
	ApplyWhen    string `json:"apply_when"`
}

// OutputSchema defines the expected output format
type OutputSchema struct {
	Format      string        `json:"format"`
	Fields      []SchemaField `json:"fields"`
	Constraints []string      `json:"constraints"`
	Example     string        `json:"example"`
}

// SchemaField represents a field in the output schema
type SchemaField struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Required    bool   `json:"required"`
	Description string `json:"description"`
	Default     string `json:"default,omitempty"`
}

// PromptTemplateManager manages prompt templates
type PromptTemplateManager interface {
	// GetTemplate retrieves a template by ID
	GetTemplate(ctx context.Context, templateID string) (*PromptTemplate, error)

	// GetTemplateByContext retrieves the most appropriate template for the given context
	GetTemplateByContext(ctx context.Context, contextType string, profile *dto.CustomerProfile) (*PromptTemplate, error)

	// ListTemplates returns all available templates
	ListTemplates(ctx context.Context) ([]*PromptTemplate, error)

	// RegisterTemplate registers a new template
	RegisterTemplate(ctx context.Context, template *PromptTemplate) error

	// UpdateTemplate updates an existing template
	UpdateTemplate(ctx context.Context, template *PromptTemplate) error

	// ValidateTemplate validates a template structure
	ValidateTemplate(ctx context.Context, template *PromptTemplate) error
}

// PromptContext contains context information for prompt generation
type PromptContext struct {
	CustomerProfile *dto.CustomerProfile           `json:"customer_profile"`
	Products        []dto.ProductRecommendationV2  `json:"products"`
	ContextType     string                         `json:"context_type"`
	CurrentProduct  *dto.ProductRecommendationV2   `json:"current_product,omitempty"`
	SessionInfo     map[string]interface{}         `json:"session_info,omitempty"`
	Metadata        map[string]interface{}         `json:"metadata,omitempty"`
}

// PromptVariables contains template variables for prompt generation
type PromptVariables struct {
	CustomerProfile string `json:"customer_profile"`
	Products        string `json:"products"`
	CurrentProduct  string `json:"current_product"`
	OutputSchema    string `json:"output_schema"`
	Examples        string `json:"examples"`
	ReasoningChain  string `json:"reasoning_chain"`
	ContextInfo     string `json:"context_info"`
}
