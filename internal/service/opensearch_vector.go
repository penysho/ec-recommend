package service

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
)

// OpenSearchVectorService implements the OpenSearchVectorInterface for Amazon OpenSearch Service
type OpenSearchVectorService struct {
	client      *http.Client
	endpoint    string
	region      string
	credentials aws.CredentialsProvider
	signer      *v4.Signer
}

// NewOpenSearchVectorService creates a new OpenSearch vector service instance
func NewOpenSearchVectorService(endpoint, region string, credentials aws.CredentialsProvider) *OpenSearchVectorService {
	// Create HTTP client with timeout and TLS configuration
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: false,
			},
			MaxIdleConns:        10,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     30 * time.Second,
		},
	}

	return &OpenSearchVectorService{
		client:      client,
		endpoint:    strings.TrimSuffix(endpoint, "/"),
		region:      region,
		credentials: credentials,
		signer:      v4.NewSigner(),
	}
}

// VectorSearch performs vector similarity search using k-NN
func (os *OpenSearchVectorService) VectorSearch(ctx context.Context, req *VectorSearchRequest) (*VectorSearchResponse, error) {
	startTime := time.Now()

	// Build k-NN search query
	query := map[string]interface{}{
		"size": req.Size,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"knn": map[string]interface{}{
							"product_vector": map[string]interface{}{
								"vector": req.Vector,
								"k":      req.Size * 2, // Search more than needed for better results
							},
						},
					},
				},
			},
		},
		"_source": map[string]interface{}{
			"excludes": []string{"product_vector"}, // Exclude vector from response to reduce size
		},
	}

	// Add filters if provided
	if req.Filters != nil && len(req.Filters) > 0 {
		filterClauses := os.buildFilterClauses(req.Filters)
		if len(filterClauses) > 0 {
			query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"] = filterClauses
		}
	}

	// Add minimum score threshold
	if req.MinScore > 0 {
		query["min_score"] = req.MinScore
	}

	// Execute search
	results, err := os.executeSearch(ctx, req.IndexName, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute vector search: %w", err)
	}

	// Parse results
	vectorResults := make([]VectorSearchResult, 0)
	totalFound := 0

	if hits, ok := results["hits"].(map[string]interface{}); ok {
		if total, ok := hits["total"].(map[string]interface{}); ok {
			if value, ok := total["value"].(float64); ok {
				totalFound = int(value)
			}
		}

		if hitsList, ok := hits["hits"].([]interface{}); ok {
			for _, hit := range hitsList {
				if hitMap, ok := hit.(map[string]interface{}); ok {
					result := VectorSearchResult{
						ID:    fmt.Sprintf("%v", hitMap["_id"]),
						Score: float64(hitMap["_score"].(float64)),
					}

					// Extract metadata from source
					if source, ok := hitMap["_source"].(map[string]interface{}); ok {
						result.Metadata = source
						if content, ok := source["content"].(string); ok {
							result.Content = content
						}
					}

					// Include vector if requested
					if req.IncludeMetadata {
						// Vector is excluded from _source, but can be included if needed
						result.Vector = make([]float64, 0)
					}

					vectorResults = append(vectorResults, result)
				}
			}
		}
	}

	processingTime := time.Since(startTime).Milliseconds()

	return &VectorSearchResponse{
		Results:          vectorResults,
		TotalFound:       totalFound,
		ProcessingTimeMs: processingTime,
	}, nil
}

// HybridSearch performs hybrid search combining vector and text search
func (os *OpenSearchVectorService) HybridSearch(ctx context.Context, req *HybridSearchRequest) (*HybridSearchResponse, error) {
	startTime := time.Now()

	// Set default weights if not provided
	vectorWeight := req.VectorWeight
	textWeight := req.TextWeight
	if vectorWeight == 0 && textWeight == 0 {
		vectorWeight = 0.7
		textWeight = 0.3
	}

	// Build hybrid search query
	query := map[string]interface{}{
		"size": req.Size,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"should": []map[string]interface{}{},
			},
		},
		"_source": map[string]interface{}{
			"excludes": []string{"product_vector"},
		},
	}

	shouldClauses := []map[string]interface{}{}

	// Add vector search component if vector provided
	if req.Vector != nil && len(req.Vector) > 0 {
		knnClause := map[string]interface{}{
			"knn": map[string]interface{}{
				"product_vector": map[string]interface{}{
					"vector": req.Vector,
					"k":      req.Size * 2,
				},
			},
			"boost": vectorWeight,
		}
		shouldClauses = append(shouldClauses, knnClause)
	}

	// Add text search component
	if req.Query != "" {
		textClause := map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  req.Query,
				"fields": []string{"name^3", "description^2", "brand", "category_name", "tags"},
				"type":   "best_fields",
				"boost":  textWeight,
			},
		}
		shouldClauses = append(shouldClauses, textClause)

		// Add phrase matching for exact matches
		phraseClause := map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  req.Query,
				"fields": []string{"name^2", "description"},
				"type":   "phrase",
				"boost":  textWeight * 1.5,
			},
		}
		shouldClauses = append(shouldClauses, phraseClause)
	}

	query["query"].(map[string]interface{})["bool"].(map[string]interface{})["should"] = shouldClauses

	// Add filters if provided
	if req.Filters != nil && len(req.Filters) > 0 {
		filterClauses := os.buildFilterClauses(req.Filters)
		if len(filterClauses) > 0 {
			query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"] = filterClauses
		}
	}

	// Execute search
	results, err := os.executeSearch(ctx, req.IndexName, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute hybrid search: %w", err)
	}

	// Parse results with separate scoring
	hybridResults := make([]HybridSearchResult, 0)
	totalFound := 0

	if hits, ok := results["hits"].(map[string]interface{}); ok {
		if total, ok := hits["total"].(map[string]interface{}); ok {
			if value, ok := total["value"].(float64); ok {
				totalFound = int(value)
			}
		}

		if hitsList, ok := hits["hits"].([]interface{}); ok {
			for _, hit := range hitsList {
				if hitMap, ok := hit.(map[string]interface{}); ok {
					result := HybridSearchResult{
						ID:            fmt.Sprintf("%v", hitMap["_id"]),
						CombinedScore: float64(hitMap["_score"].(float64)),
						VectorScore:   0.0, // Would need separate vector scoring
						TextScore:     0.0, // Would need separate text scoring
					}

					// Extract metadata from source
					if source, ok := hitMap["_source"].(map[string]interface{}); ok {
						result.Metadata = source
						if content, ok := source["content"].(string); ok {
							result.Content = content
						}
					}

					hybridResults = append(hybridResults, result)
				}
			}
		}
	}

	processingTime := time.Since(startTime).Milliseconds()

	return &HybridSearchResponse{
		Results:          hybridResults,
		TotalFound:       totalFound,
		ProcessingTimeMs: processingTime,
	}, nil
}

// IndexVector indexes a vector with metadata
func (os *OpenSearchVectorService) IndexVector(ctx context.Context, req *IndexVectorRequest) error {
	// Prepare document for indexing
	document := map[string]interface{}{
		"product_vector": req.Vector,
		"content":        req.Content,
	}

	// Add metadata fields
	for key, value := range req.Metadata {
		document[key] = value
	}

	// Convert to JSON
	docBytes, err := json.Marshal(document)
	if err != nil {
		return fmt.Errorf("failed to marshal document: %w", err)
	}

	// Create HTTP request
	url := fmt.Sprintf("%s/%s/_doc/%s", os.endpoint, req.IndexName, req.ID)
	httpReq, err := http.NewRequestWithContext(ctx, "PUT", url, bytes.NewReader(docBytes))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	// Sign request
	if err := os.signRequest(ctx, httpReq); err != nil {
		return fmt.Errorf("failed to sign request: %w", err)
	}

	// Execute request
	resp, err := os.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("indexing failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// DeleteVector removes a vector from the index
func (os *OpenSearchVectorService) DeleteVector(ctx context.Context, vectorID string) error {
	// This requires knowing the index name - for now, we'll use a default
	// In production, you'd want to pass the index name or store it with the vector
	indexName := "products"

	url := fmt.Sprintf("%s/%s/_doc/%s", os.endpoint, indexName, vectorID)
	httpReq, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Sign request
	if err := os.signRequest(ctx, httpReq); err != nil {
		return fmt.Errorf("failed to sign request: %w", err)
	}

	// Execute request
	resp, err := os.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 404 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("deletion failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// GetVectorMetadata retrieves metadata for a vector
func (os *OpenSearchVectorService) GetVectorMetadata(ctx context.Context, vectorID string) (*VectorMetadata, error) {
	// This requires knowing the index name - for now, we'll use a default
	indexName := "products"

	url := fmt.Sprintf("%s/%s/_doc/%s", os.endpoint, indexName, vectorID)
	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Sign request
	if err := os.signRequest(ctx, httpReq); err != nil {
		return nil, fmt.Errorf("failed to sign request: %w", err)
	}

	// Execute request
	resp, err := os.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("vector not found")
	}

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("get failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Extract metadata
	source, ok := response["_source"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}

	// Remove vector from metadata to get clean metadata
	metadata := make(map[string]interface{})
	for key, value := range source {
		if key != "product_vector" {
			metadata[key] = value
		}
	}

	vectorSize := 0
	if vector, ok := source["product_vector"].([]interface{}); ok {
		vectorSize = len(vector)
	}

	return &VectorMetadata{
		ID:         vectorID,
		IndexName:  indexName,
		Metadata:   metadata,
		CreatedAt:  time.Now().Unix(), // Would need to be stored in the document
		UpdatedAt:  time.Now().Unix(), // Would need to be stored in the document
		VectorSize: vectorSize,
	}, nil
}

// Helper methods

func (os *OpenSearchVectorService) executeSearch(ctx context.Context, indexName string, query map[string]interface{}) (map[string]interface{}, error) {
	// Convert query to JSON
	queryBytes, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query: %w", err)
	}

	// Create HTTP request
	url := fmt.Sprintf("%s/%s/_search", os.endpoint, indexName)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(queryBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	// Sign request
	if err := os.signRequest(ctx, httpReq); err != nil {
		return nil, fmt.Errorf("failed to sign request: %w", err)
	}

	// Execute request
	resp, err := os.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("search failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var results map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return results, nil
}

func (os *OpenSearchVectorService) buildFilterClauses(filters map[string]interface{}) []map[string]interface{} {
	var filterClauses []map[string]interface{}

	for key, value := range filters {
		switch key {
		case "category_id", "brand_id", "price_min", "price_max":
			if key == "price_min" {
				filterClauses = append(filterClauses, map[string]interface{}{
					"range": map[string]interface{}{
						"price": map[string]interface{}{
							"gte": value,
						},
					},
				})
			} else if key == "price_max" {
				filterClauses = append(filterClauses, map[string]interface{}{
					"range": map[string]interface{}{
						"price": map[string]interface{}{
							"lte": value,
						},
					},
				})
			} else {
				filterClauses = append(filterClauses, map[string]interface{}{
					"term": map[string]interface{}{
						key: value,
					},
				})
			}
		default:
			// Generic filter
			filterClauses = append(filterClauses, map[string]interface{}{
				"term": map[string]interface{}{
					key: value,
				},
			})
		}
	}

	return filterClauses
}

func (os *OpenSearchVectorService) signRequest(ctx context.Context, req *http.Request) error {
	// Get credentials
	creds, err := os.credentials.Retrieve(ctx)
	if err != nil {
		return fmt.Errorf("failed to retrieve credentials: %w", err)
	}

	// Sign request
	err = os.signer.SignHTTP(ctx, creds, req, "", "es", os.region, time.Now())
	if err != nil {
		return fmt.Errorf("failed to sign request: %w", err)
	}

	return nil
}

// CreateIndex creates an index with vector field mapping
func (os *OpenSearchVectorService) CreateIndex(ctx context.Context, indexName string, vectorDimension int) error {
	mapping := map[string]interface{}{
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"product_vector": map[string]interface{}{
					"type":      "knn_vector",
					"dimension": vectorDimension,
					"method": map[string]interface{}{
						"name":   "hnsw",
						"engine": "nmslib",
						"parameters": map[string]interface{}{
							"ef_construction": 128,
							"m":               24,
						},
					},
				},
				"name": map[string]interface{}{
					"type":     "text",
					"analyzer": "standard",
				},
				"description": map[string]interface{}{
					"type":     "text",
					"analyzer": "standard",
				},
				"brand": map[string]interface{}{
					"type": "keyword",
				},
				"category_name": map[string]interface{}{
					"type": "keyword",
				},
				"price": map[string]interface{}{
					"type": "double",
				},
				"tags": map[string]interface{}{
					"type": "keyword",
				},
			},
		},
		"settings": map[string]interface{}{
			"index": map[string]interface{}{
				"knn":                      true,
				"knn.algo_param.ef_search": 100,
			},
		},
	}

	// Convert to JSON
	mappingBytes, err := json.Marshal(mapping)
	if err != nil {
		return fmt.Errorf("failed to marshal mapping: %w", err)
	}

	// Create HTTP request
	url := fmt.Sprintf("%s/%s", os.endpoint, indexName)
	httpReq, err := http.NewRequestWithContext(ctx, "PUT", url, bytes.NewReader(mappingBytes))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	// Sign request
	if err := os.signRequest(ctx, httpReq); err != nil {
		return fmt.Errorf("failed to sign request: %w", err)
	}

	// Execute request
	resp, err := os.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("index creation failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
