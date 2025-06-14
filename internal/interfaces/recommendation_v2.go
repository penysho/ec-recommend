package interfaces

import (
	"context"
	"ec-recommend/internal/dto"

	"github.com/google/uuid"
)

// RecommendationServiceV2Interface defines the enhanced interface for recommendation operations using RAG and vector search
type RecommendationServiceV2Interface interface {
	// GetRecommendationsV2 generates advanced product recommendations using RAG and vector search
	GetRecommendationsV2(ctx context.Context, req *dto.RecommendationRequestV2) (*dto.RecommendationResponseV2, error)

	// SemanticSearch performs semantic search using natural language queries
	SemanticSearch(ctx context.Context, req *dto.SemanticSearchRequest) (*dto.SemanticSearchResponse, error)

	// GetVectorSimilarProducts finds products similar to a given product using vector similarity
	GetVectorSimilarProducts(ctx context.Context, req *dto.VectorSimilarityRequest) (*dto.VectorSimilarityResponse, error)

	// GetKnowledgeBasedRecommendations generates recommendations based on comprehensive knowledge base analysis
	GetKnowledgeBasedRecommendations(ctx context.Context, req *dto.KnowledgeBasedRecommendationRequest) (*dto.KnowledgeBasedRecommendationResponse, error)

	// GetRecommendationExplanation provides detailed explanation for a specific recommendation
	GetRecommendationExplanation(ctx context.Context, recommendationID, customerID uuid.UUID) (*dto.RecommendationExplanationResponse, error)

	// GetTrendingProductsV2 returns trending products with AI-powered insights
	GetTrendingProductsV2(ctx context.Context, req *dto.TrendingProductsRequestV2) (*dto.TrendingProductsResponseV2, error)

	// GetCustomerProfile retrieves customer profile data for personalized recommendations
	GetCustomerProfile(ctx context.Context, customerID uuid.UUID) (*dto.CustomerProfile, error)

	// LogRecommendationInteraction logs customer interactions with recommendations for analytics
	LogRecommendationInteraction(ctx context.Context, analytics *dto.RecommendationAnalytics) error
}

// BedrockKnowledgeBaseInterface defines the interface for Amazon Bedrock Knowledge Base operations
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

// OpenSearchVectorInterface defines the interface for OpenSearch vector operations
type OpenSearchVectorInterface interface {
	// VectorSearch performs vector similarity search
	VectorSearch(ctx context.Context, req *VectorSearchRequest) (*VectorSearchResponse, error)

	// HybridSearch performs hybrid search combining vector and text search
	HybridSearch(ctx context.Context, req *HybridSearchRequest) (*HybridSearchResponse, error)

	// IndexVector indexes a vector with metadata
	IndexVector(ctx context.Context, req *IndexVectorRequest) error

	// DeleteVector removes a vector from the index
	DeleteVector(ctx context.Context, vectorID string) error

	// GetVectorMetadata retrieves metadata for a vector
	GetVectorMetadata(ctx context.Context, vectorID string) (*VectorMetadata, error)
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

// Span represents a span of text
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

// GenerationMetadata contains metadata about the generation process
type GenerationMetadata struct {
	Usage            *UsageMetadata `json:"usage,omitempty"`
	ModelID          string         `json:"model_id,omitempty"`
	ProcessingTimeMs int64          `json:"processing_time_ms"`
}

// UsageMetadata contains usage information
type UsageMetadata struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
	TotalTokens  int `json:"total_tokens"`
}

// SimilarDocumentsResponse represents the response from similar documents search
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

// VectorSearchRequest represents a request for vector search
type VectorSearchRequest struct {
	Vector          []float64              `json:"vector"`
	IndexName       string                 `json:"index_name"`
	Size            int                    `json:"size"`
	MinScore        float64                `json:"min_score,omitempty"`
	Filters         map[string]interface{} `json:"filters,omitempty"`
	IncludeMetadata bool                   `json:"include_metadata"`
}

// VectorSearchResponse represents the response from vector search
type VectorSearchResponse struct {
	Results          []VectorSearchResult `json:"results"`
	TotalFound       int                  `json:"total_found"`
	ProcessingTimeMs int64                `json:"processing_time_ms"`
}

// VectorSearchResult represents a single vector search result
type VectorSearchResult struct {
	ID       string                 `json:"id"`
	Score    float64                `json:"score"`
	Vector   []float64              `json:"vector,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	Content  string                 `json:"content,omitempty"`
}

// HybridSearchRequest represents a request for hybrid search
type HybridSearchRequest struct {
	Query           string                 `json:"query"`
	Vector          []float64              `json:"vector,omitempty"`
	IndexName       string                 `json:"index_name"`
	Size            int                    `json:"size"`
	VectorWeight    float64                `json:"vector_weight,omitempty"`
	TextWeight      float64                `json:"text_weight,omitempty"`
	Filters         map[string]interface{} `json:"filters,omitempty"`
	IncludeMetadata bool                   `json:"include_metadata"`
}

// HybridSearchResponse represents the response from hybrid search
type HybridSearchResponse struct {
	Results          []HybridSearchResult `json:"results"`
	TotalFound       int                  `json:"total_found"`
	ProcessingTimeMs int64                `json:"processing_time_ms"`
}

// HybridSearchResult represents a single hybrid search result
type HybridSearchResult struct {
	ID            string                 `json:"id"`
	VectorScore   float64                `json:"vector_score"`
	TextScore     float64                `json:"text_score"`
	CombinedScore float64                `json:"combined_score"`
	Vector        []float64              `json:"vector,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
	Content       string                 `json:"content,omitempty"`
}

// IndexVectorRequest represents a request to index a vector
type IndexVectorRequest struct {
	ID        string                 `json:"id"`
	Vector    []float64              `json:"vector"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	Content   string                 `json:"content,omitempty"`
	IndexName string                 `json:"index_name"`
}

// VectorMetadata represents metadata associated with a vector
type VectorMetadata struct {
	ID         string                 `json:"id"`
	IndexName  string                 `json:"index_name"`
	Metadata   map[string]interface{} `json:"metadata"`
	CreatedAt  int64                  `json:"created_at"`
	UpdatedAt  int64                  `json:"updated_at"`
	VectorSize int                    `json:"vector_size"`
}

// RecommendationRepositoryV2Interface defines the enhanced interface for recommendation data access
type RecommendationRepositoryV2Interface interface {
	// Enhanced customer and product methods
	GetCustomerByID(ctx context.Context, customerID uuid.UUID) (*dto.CustomerProfile, error)
	GetCustomerPurchaseHistory(ctx context.Context, customerID uuid.UUID, limit int) ([]dto.PurchaseItem, error)
	GetCustomerActivities(ctx context.Context, customerID uuid.UUID, limit int) ([]dto.ActivityItem, error)
	GetProductsByIDs(ctx context.Context, productIDs []uuid.UUID) ([]dto.ProductRecommendationV2, error)
	GetProductsByCategory(ctx context.Context, categoryID int, limit int) ([]dto.ProductRecommendationV2, error)

	// Vector and semantic search methods
	GetProductsWithVectorSearch(ctx context.Context, vector []float64, limit int, filters map[string]interface{}) ([]dto.ProductRecommendationV2, error)
	GetProductsWithSemanticSearch(ctx context.Context, query string, limit int, filters map[string]interface{}) ([]dto.ProductRecommendationV2, error)
	GetProductsWithHybridSearch(ctx context.Context, query string, vector []float64, limit int, filters map[string]interface{}) ([]dto.ProductRecommendationV2, error)

	// Enhanced analytics and trending
	GetTrendingProductsV2(ctx context.Context, categoryID *int, timeRange string, limit int) ([]dto.TrendingProductV2, error)
	GetProductPerformanceMetrics(ctx context.Context, productIDs []uuid.UUID) (map[uuid.UUID]*dto.ProductPerformanceMetrics, error)
	GetMarketAnalysis(ctx context.Context, categoryID *int, timeRange string) (*dto.MarketAnalysis, error)

	// Logging and analytics
	LogRecommendation(ctx context.Context, customerID uuid.UUID, recommendationType, contextType string, productIDs []uuid.UUID, sessionID uuid.UUID) error
	LogRecommendationInteraction(ctx context.Context, analytics *dto.RecommendationAnalytics) error
	LogSemanticSearch(ctx context.Context, customerID *uuid.UUID, query string, results []uuid.UUID, processingTimeMs int64) error

	// Cache management
	GetCachedRecommendations(ctx context.Context, key string) ([]dto.ProductRecommendationV2, error)
	SetCachedRecommendations(ctx context.Context, key string, recommendations []dto.ProductRecommendationV2, ttl int64) error
	InvalidateCache(ctx context.Context, pattern string) error
}
