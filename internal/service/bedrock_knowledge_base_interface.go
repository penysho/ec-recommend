package service

import (
	"context"
)

// BedrockKnowledgeBaseInterface defines the interface for Amazon Bedrock Knowledge Base operations
// This interface is defined in the service package as it is consumed by services
type BedrockKnowledgeBaseInterface interface {
	// QueryKnowledgeBase performs a query against the knowledge base
	QueryKnowledgeBase(ctx context.Context, query string, filters map[string]interface{}) (*BedrockKnowledgeBaseResponse, error)

	// RetrieveAndGenerate performs retrieval-augmented generation
	RetrieveAndGenerate(ctx context.Context, req *RetrieveAndGenerateRequest) (*RetrieveAndGenerateResponse, error)

	// GetVectorEmbedding generates vector embeddings for the given text
	GetVectorEmbedding(ctx context.Context, text string) ([]float64, error)

	// GetSimilarDocuments finds similar documents based on vector similarity
	GetSimilarDocuments(ctx context.Context, embedding []float64, limit int, filters map[string]interface{}) (*SimilarDocumentsResponse, error)
}

// BedrockKnowledgeBaseResponse represents the response from knowledge base query
type BedrockKnowledgeBaseResponse struct {
	Results           []KnowledgeBaseResult `json:"results"`
	RetrievalMetadata *RetrievalMetadata    `json:"retrieval_metadata,omitempty"`
	ProcessingTimeMs  int64                 `json:"processing_time_ms"`
}

// KnowledgeBaseResult represents a single result from knowledge base
type KnowledgeBaseResult struct {
	Content  string                 `json:"content"`
	Score    float64                `json:"score"`
	Source   string                 `json:"source,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	Location *DocumentLocation      `json:"location,omitempty"`
}

// RetrievalMetadata contains metadata about the retrieval process
type RetrievalMetadata struct {
	QueryProcessingTimeMs int64    `json:"query_processing_time_ms"`
	RetrievalCount        int      `json:"retrieval_count"`
	Sources               []string `json:"sources"`
	ConfidenceLevel       float64  `json:"confidence_level"`
}

// DocumentLocation represents the location of content within a document
type DocumentLocation struct {
	DocumentID string `json:"document_id"`
	Page       int    `json:"page,omitempty"`
	Chapter    string `json:"chapter,omitempty"`
	Section    string `json:"section,omitempty"`
}

// RetrieveAndGenerateRequest represents a request for retrieval-augmented generation
type RetrieveAndGenerateRequest struct {
	Query            string                 `json:"query"`
	KnowledgeBaseID  string                 `json:"knowledge_base_id"`
	ModelARN         string                 `json:"model_arn"`
	RetrievalConfig  *RetrievalConfig       `json:"retrieval_config,omitempty"`
	GenerationConfig *GenerationConfig      `json:"generation_config,omitempty"`
	Filters          map[string]interface{} `json:"filters,omitempty"`
}

// RetrievalConfig represents configuration for retrieval
type RetrievalConfig struct {
	NumberOfResults    int                    `json:"number_of_results,omitempty"`
	SearchType         string                 `json:"search_type,omitempty"` // "semantic", "hybrid"
	VectorSearchConfig *VectorSearchConfig    `json:"vector_search_config,omitempty"`
	MetadataFilters    map[string]interface{} `json:"metadata_filters,omitempty"`
}

// VectorSearchConfig represents configuration for vector search
type VectorSearchConfig struct {
	NumberOfResults    int    `json:"number_of_results"`
	OverrideSearchType string `json:"override_search_type,omitempty"`
}

// GenerationConfig represents configuration for generation
type GenerationConfig struct {
	Temperature                  float64                `json:"temperature,omitempty"`
	TopP                         float64                `json:"top_p,omitempty"`
	TopK                         int                    `json:"top_k,omitempty"`
	MaxTokens                    int                    `json:"max_tokens,omitempty"`
	StopSequences                []string               `json:"stop_sequences,omitempty"`
	AdditionalModelRequestFields map[string]interface{} `json:"additional_model_request_fields,omitempty"`
}

// RetrieveAndGenerateResponse represents the response from retrieval-augmented generation
type RetrieveAndGenerateResponse struct {
	Output           string                `json:"output"`
	Citations        []Citation            `json:"citations,omitempty"`
	RetrievedResults []KnowledgeBaseResult `json:"retrieved_results,omitempty"`
	Metadata         *GenerationMetadata   `json:"metadata,omitempty"`
}

// Citation represents a citation in the generated response
type Citation struct {
	GeneratedResponsePart *GeneratedResponsePart `json:"generated_response_part,omitempty"`
	RetrievedReferences   []RetrievedReference   `json:"retrieved_references,omitempty"`
}

// GeneratedResponsePart represents a part of the generated response
type GeneratedResponsePart struct {
	TextResponsePart *TextResponsePart `json:"text_response_part,omitempty"`
}

// TextResponsePart represents the text part of the response
type TextResponsePart struct {
	Text string `json:"text"`
	Span *Span  `json:"span,omitempty"`
}

// Span represents a span in the text
type Span struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

// RetrievedReference represents a retrieved reference
type RetrievedReference struct {
	Content  *RetrievalResultContent  `json:"content,omitempty"`
	Location *RetrievalResultLocation `json:"location,omitempty"`
	Metadata map[string]interface{}   `json:"metadata,omitempty"`
}

// RetrievalResultContent represents the content of a retrieval result
type RetrievalResultContent struct {
	Text string `json:"text"`
}

// RetrievalResultLocation represents the location of a retrieval result
type RetrievalResultLocation struct {
	Type               string              `json:"type"`
	S3Location         *S3Location         `json:"s3_location,omitempty"`
	WebLocation        *WebLocation        `json:"web_location,omitempty"`
	ConfluenceLocation *ConfluenceLocation `json:"confluence_location,omitempty"`
	SalesforceLocation *SalesforceLocation `json:"salesforce_location,omitempty"`
	SharePointLocation *SharePointLocation `json:"share_point_location,omitempty"`
}

// S3Location represents an S3 location
type S3Location struct {
	URI string `json:"uri"`
}

// WebLocation represents a web location
type WebLocation struct {
	URL string `json:"url"`
}

// ConfluenceLocation represents a Confluence location
type ConfluenceLocation struct {
	URL string `json:"url"`
}

// SalesforceLocation represents a Salesforce location
type SalesforceLocation struct {
	URL string `json:"url"`
}

// SharePointLocation represents a SharePoint location
type SharePointLocation struct {
	URL string `json:"url"`
}

// GenerationMetadata represents metadata about generation
type GenerationMetadata struct {
	Usage            *UsageMetadata `json:"usage,omitempty"`
	ModelID          string         `json:"model_id,omitempty"`
	ProcessingTimeMs int64          `json:"processing_time_ms"`
}

// UsageMetadata represents usage metadata
type UsageMetadata struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
	TotalTokens  int `json:"total_tokens"`
}

// SimilarDocumentsResponse represents the response for similar documents
type SimilarDocumentsResponse struct {
	Documents        []SimilarDocument `json:"documents"`
	ProcessingTimeMs int64             `json:"processing_time_ms"`
}

// SimilarDocument represents a similar document
type SimilarDocument struct {
	DocumentID string                 `json:"document_id"`
	Content    string                 `json:"content"`
	Score      float64                `json:"score"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
	Source     string                 `json:"source,omitempty"`
}
