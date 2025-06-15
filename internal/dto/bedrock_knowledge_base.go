package dto

import (
	"github.com/google/uuid"
)

// BedrockSearchResult represents a search result from Bedrock Knowledge Base
type BedrockSearchResult struct {
	ProductID        uuid.UUID              `json:"product_id"`
	DistanceScore    float64                `json:"distance_score"`
	SimilarityScore  float64                `json:"similarity_score"`
	ConfidenceScore  float64                `json:"confidence_score"`
	SearchMethod     string                 `json:"search_method"`
	EmbeddingModel   string                 `json:"embedding_model"`
	MatchedCriteria  []string               `json:"matched_criteria"`
	SemanticClusters []string               `json:"semantic_clusters"`
	Metadata         map[string]interface{} `json:"metadata"`
	Source           string                 `json:"source"`
	RetrievalRank    int                    `json:"retrieval_rank"`
}

// BedrockVectorSearchResponse represents the response from vector search
type BedrockVectorSearchResponse struct {
	Results          []BedrockSearchResult `json:"results"`
	TotalFound       int                   `json:"total_found"`
	ProcessingTimeMs int64                 `json:"processing_time_ms"`
	SearchMetadata   *BedrockSearchMeta    `json:"search_metadata"`
}

// BedrockSemanticSearchResponse represents the response from semantic search
type BedrockSemanticSearchResponse struct {
	Query            string                `json:"query"`
	Results          []BedrockSearchResult `json:"results"`
	TotalFound       int                   `json:"total_found"`
	ProcessingTimeMs int64                 `json:"processing_time_ms"`
	SearchMetadata   *BedrockSearchMeta    `json:"search_metadata"`
}

// BedrockHybridSearchResponse represents the response from hybrid search
type BedrockHybridSearchResponse struct {
	Query            string                `json:"query"`
	Vector           []float64             `json:"vector,omitempty"`
	Results          []BedrockSearchResult `json:"results"`
	TotalFound       int                   `json:"total_found"`
	ProcessingTimeMs int64                 `json:"processing_time_ms"`
	SearchMetadata   *BedrockSearchMeta    `json:"search_metadata"`
}

// BedrockSearchMeta contains metadata about the search operation
type BedrockSearchMeta struct {
	SearchType         string                 `json:"search_type"`
	EmbeddingModel     string                 `json:"embedding_model"`
	KnowledgeBaseID    string                 `json:"knowledge_base_id"`
	SimilarityMetric   string                 `json:"similarity_metric"`
	FiltersApplied     map[string]interface{} `json:"filters_applied"`
	RerankerUsed       bool                   `json:"reranker_used"`
	CacheUsed          bool                   `json:"cache_used"`
	HybridSearchWeight map[string]float64     `json:"hybrid_search_weight,omitempty"`
}

// BedrockEmbeddingResponse represents the response from embedding generation
type BedrockEmbeddingResponse struct {
	Embedding        []float64 `json:"embedding"`
	ModelID          string    `json:"model_id"`
	ProcessingTimeMs int64     `json:"processing_time_ms"`
}

// BedrockKnowledgeBaseQueryResponse represents the response from knowledge base query
type BedrockKnowledgeBaseQueryResponse struct {
	Query             string                    `json:"query"`
	Results           []BedrockKnowledgeResult  `json:"results"`
	ProcessingTimeMs  int64                     `json:"processing_time_ms"`
	RetrievalMetadata *BedrockRetrievalMetadata `json:"retrieval_metadata"`
	ConfidenceLevel   float64                   `json:"confidence_level"`
	Sources           []string                  `json:"sources"`
}

// BedrockKnowledgeResult represents a single result from knowledge base
type BedrockKnowledgeResult struct {
	Content    string                 `json:"content"`
	Score      float64                `json:"score"`
	Source     string                 `json:"source"`
	Metadata   map[string]interface{} `json:"metadata"`
	DocumentID string                 `json:"document_id"`
}

// BedrockRetrievalMetadata contains metadata about the retrieval operation
type BedrockRetrievalMetadata struct {
	QueryProcessingTimeMs int64    `json:"query_processing_time_ms"`
	RetrievalCount        int      `json:"retrieval_count"`
	Sources               []string `json:"sources"`
	ConfidenceLevel       float64  `json:"confidence_level"`
}

// BedrockRAGResponse represents the response from Retrieve and Generate operation
type BedrockRAGResponse struct {
	Query            string                   `json:"query"`
	GeneratedText    string                   `json:"generated_text"`
	Citations        []BedrockCitation        `json:"citations"`
	RetrievedResults []BedrockKnowledgeResult `json:"retrieved_results"`
	ProcessingTimeMs int64                    `json:"processing_time_ms"`
	Metadata         *BedrockRAGMetadata      `json:"metadata"`
}

// BedrockCitation represents a citation in the generated response
type BedrockCitation struct {
	GeneratedTextSpan *BedrockTextSpan         `json:"generated_text_span"`
	References        []BedrockReferenceSource `json:"references"`
}

// BedrockTextSpan represents a span of text in the generated response
type BedrockTextSpan struct {
	Text  string `json:"text"`
	Start int    `json:"start"`
	End   int    `json:"end"`
}

// BedrockReferenceSource represents a reference source for a citation
type BedrockReferenceSource struct {
	Content    string                 `json:"content"`
	Source     string                 `json:"source"`
	Metadata   map[string]interface{} `json:"metadata"`
	DocumentID string                 `json:"document_id"`
}

// BedrockRAGMetadata contains metadata about the RAG operation
type BedrockRAGMetadata struct {
	ModelID          string `json:"model_id"`
	ProcessingTimeMs int64  `json:"processing_time_ms"`
	TokensUsed       int    `json:"tokens_used"`
	CitationCount    int    `json:"citation_count"`
}

// BedrockProductSearchRequest represents a request for product search
type BedrockProductSearchRequest struct {
	Query      string                 `json:"query,omitempty"`
	Vector     []float64              `json:"vector,omitempty"`
	SearchType string                 `json:"search_type"` // "vector", "semantic", "hybrid"
	Limit      int                    `json:"limit"`
	Filters    map[string]interface{} `json:"filters,omitempty"`
	CustomerID uuid.UUID              `json:"customer_id,omitempty"`
}
