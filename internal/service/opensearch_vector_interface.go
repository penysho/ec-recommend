package service

import (
	"context"
)

// OpenSearchVectorInterface defines the interface for OpenSearch vector operations
// This interface is defined in the service package as it is consumed by services
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

// VectorSearchRequest represents a vector search request
type VectorSearchRequest struct {
	Vector          []float64              `json:"vector"`
	IndexName       string                 `json:"index_name"`
	Size            int                    `json:"size"`
	MinScore        float64                `json:"min_score,omitempty"`
	Filters         map[string]interface{} `json:"filters,omitempty"`
	IncludeMetadata bool                   `json:"include_metadata"`
}

// VectorSearchResponse represents a vector search response
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

// HybridSearchRequest represents a hybrid search request
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

// HybridSearchResponse represents a hybrid search response
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

// VectorMetadata represents metadata for a vector
type VectorMetadata struct {
	ID         string                 `json:"id"`
	IndexName  string                 `json:"index_name"`
	Metadata   map[string]interface{} `json:"metadata"`
	CreatedAt  int64                  `json:"created_at"`
	UpdatedAt  int64                  `json:"updated_at"`
	VectorSize int                    `json:"vector_size"`
}
