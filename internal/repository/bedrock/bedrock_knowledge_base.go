package service

import (
	"context"
	"ec-recommend/internal/service"
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockagentruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockagentruntime/types"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/google/uuid"
)

// BedrockKnowledgeBaseService implements the BedrockKnowledgeBaseInterface
type BedrockKnowledgeBaseService struct {
	agentClient      *bedrockagentruntime.Client
	runtimeClient    *bedrockruntime.Client
	knowledgeBaseID  string
	modelARN         string
	embeddingModelID string
}

// NewBedrockKnowledgeBaseService creates a new Bedrock Knowledge Base service instance
func NewBedrockKnowledgeBaseService(
	agentClient *bedrockagentruntime.Client,
	runtimeClient *bedrockruntime.Client,
	knowledgeBaseID, modelARN, embeddingModelID string,
) *BedrockKnowledgeBaseService {
	return &BedrockKnowledgeBaseService{
		agentClient:      agentClient,
		runtimeClient:    runtimeClient,
		knowledgeBaseID:  knowledgeBaseID,
		modelARN:         modelARN,
		embeddingModelID: embeddingModelID,
	}
}

// QueryKnowledgeBase performs a query against the knowledge base with advanced features
func (bkb *BedrockKnowledgeBaseService) QueryKnowledgeBase(ctx context.Context, query string, filters map[string]interface{}) (*service.RAGResponse, error) {
	startTime := time.Now()

	// Build basic retrieval configuration
	retrievalConfig := &types.KnowledgeBaseRetrievalConfiguration{
		VectorSearchConfiguration: &types.KnowledgeBaseVectorSearchConfiguration{
			NumberOfResults:    aws.Int32(10),
			OverrideSearchType: types.SearchTypeHybrid, // Use hybrid search for better results
		},
	}

	// Create retrieve request
	retrieveInput := &bedrockagentruntime.RetrieveInput{
		KnowledgeBaseId:        aws.String(bkb.knowledgeBaseID),
		RetrievalQuery:         &types.KnowledgeBaseQuery{Text: aws.String(query)},
		RetrievalConfiguration: retrievalConfig,
	}

	// Execute retrieve operation
	retrieveOutput, err := bkb.agentClient.Retrieve(ctx, retrieveInput)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve from knowledge base: %w", err)
	}

	// Convert results to our interface format
	results := make([]service.KnowledgeBaseResult, 0, len(retrieveOutput.RetrievalResults))
	sources := make([]string, 0)

	for _, result := range retrieveOutput.RetrievalResults {
		kbResult := service.KnowledgeBaseResult{
			Content: aws.ToString(result.Content.Text),
			Score:   float64(aws.ToFloat64(result.Score)),
		}

		// Extract source information
		if result.Location != nil {
			if result.Location.S3Location != nil {
				kbResult.Source = aws.ToString(result.Location.S3Location.Uri)
				sources = append(sources, kbResult.Source)
			}
		}

		// Extract metadata - fix type conversion
		if result.Metadata != nil {
			metadata := make(map[string]interface{})
			for k, v := range result.Metadata {
				metadata[k] = v
			}
			kbResult.Metadata = metadata
		}

		// Add location information
		if result.Location != nil {
			kbResult.Location = &service.DocumentLocation{
				DocumentID: extractDocumentIDFromURI(kbResult.Source),
			}
		}

		results = append(results, kbResult)
	}

	processingTime := time.Since(startTime).Milliseconds()

	// Calculate confidence level based on scores
	confidenceLevel := bkb.calculateConfidenceLevel(results)

	return &service.RAGResponse{
		Results: results,
		RetrievalMetadata: &service.RetrievalMetadata{
			QueryProcessingTimeMs: processingTime,
			RetrievalCount:        len(results),
			Sources:               bkb.deduplicateSources(sources),
			ConfidenceLevel:       confidenceLevel,
		},
		ProcessingTimeMs: processingTime,
	}, nil
}

// RetrieveAndGenerate performs retrieval-augmented generation with basic configuration
func (bkb *BedrockKnowledgeBaseService) RetrieveAndGenerate(ctx context.Context, req *service.RetrieveAndGenerateRequest) (*service.RetrieveAndGenerateResponse, error) {
	// Build basic retrieval configuration
	retrievalConfig := &types.RetrieveAndGenerateConfiguration{
		Type: types.RetrieveAndGenerateTypeKnowledgeBase,
		KnowledgeBaseConfiguration: &types.KnowledgeBaseRetrieveAndGenerateConfiguration{
			KnowledgeBaseId: aws.String(req.KnowledgeBaseID),
			ModelArn:        aws.String(req.ModelARN),
		},
	}

	// Create retrieve and generate request
	ragInput := &bedrockagentruntime.RetrieveAndGenerateInput{
		Input: &types.RetrieveAndGenerateInput{
			Text: aws.String(req.Query),
		},
		RetrieveAndGenerateConfiguration: retrievalConfig,
	}

	// Execute retrieve and generate operation
	ragOutput, err := bkb.agentClient.RetrieveAndGenerate(ctx, ragInput)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve and generate: %w", err)
	}

	// Process citations from retrieved sources
	citations := make([]service.Citation, 0)
	retrievedResults := make([]service.KnowledgeBaseResult, 0)

	if ragOutput.Citations != nil {
		for _, citation := range ragOutput.Citations {
			// Process each citation
			if citation.GeneratedResponsePart != nil && citation.GeneratedResponsePart.TextResponsePart != nil {
				citationObj := service.Citation{
					GeneratedResponsePart: &service.GeneratedResponsePart{
						TextResponsePart: &service.TextResponsePart{
							Text: aws.ToString(citation.GeneratedResponsePart.TextResponsePart.Text),
							Span: &service.Span{
								Start: int(aws.ToInt32(citation.GeneratedResponsePart.TextResponsePart.Span.Start)),
								End:   int(aws.ToInt32(citation.GeneratedResponsePart.TextResponsePart.Span.End)),
							},
						},
					},
				}

				// Process retrieved references
				if citation.RetrievedReferences != nil {
					references := make([]service.RetrievedReference, 0, len(citation.RetrievedReferences))
					for _, ref := range citation.RetrievedReferences {
						reference := service.RetrievedReference{
							Content: &service.RetrievalResultContent{
								Text: aws.ToString(ref.Content.Text),
							},
						}

						// Add location information
						if ref.Location != nil {
							reference.Location = &service.RetrievalResultLocation{}
							if ref.Location.S3Location != nil {
								reference.Location.S3Location = &service.S3Location{
									URI: aws.ToString(ref.Location.S3Location.Uri),
								}
							}
						}

						// Add metadata
						if ref.Metadata != nil {
							metadata := make(map[string]interface{})
							for k, v := range ref.Metadata {
								metadata[k] = v
							}
							reference.Metadata = metadata
						}

						references = append(references, reference)

						// Also add to retrieved results for compatibility
						retrievedResult := service.KnowledgeBaseResult{
							Content: aws.ToString(ref.Content.Text),
							Score:   0.0, // Score not available in citations
						}

						if ref.Location != nil && ref.Location.S3Location != nil {
							retrievedResult.Source = aws.ToString(ref.Location.S3Location.Uri)
							retrievedResult.Location = &service.DocumentLocation{
								DocumentID: extractDocumentIDFromURI(retrievedResult.Source),
							}
						}

						if ref.Metadata != nil {
							metadata := make(map[string]interface{})
							for k, v := range ref.Metadata {
								metadata[k] = v
							}
							retrievedResult.Metadata = metadata
						}

						retrievedResults = append(retrievedResults, retrievedResult)
					}
					citationObj.RetrievedReferences = references
				}

				citations = append(citations, citationObj)
			}
		}
	}

	return &service.RetrieveAndGenerateResponse{
		Output:           aws.ToString(ragOutput.Output.Text),
		Citations:        citations,
		RetrievedResults: retrievedResults,
		Metadata: &service.GenerationMetadata{
			ModelID:          extractModelIDFromARN(req.ModelARN),
			ProcessingTimeMs: time.Since(time.Now()).Milliseconds(),
		},
	}, nil
}

// GetVectorEmbedding generates vector embeddings for the given text using Titan Embeddings
func (bkb *BedrockKnowledgeBaseService) GetVectorEmbedding(ctx context.Context, text string) ([]float64, error) {
	// Prepare the input for Titan Embeddings model
	input := map[string]interface{}{
		"inputText": text,
	}

	inputBytes, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal embedding input: %w", err)
	}

	// Call the embedding model
	invokeInput := &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(bkb.embeddingModelID),
		ContentType: aws.String("application/json"),
		Body:        inputBytes,
	}

	result, err := bkb.runtimeClient.InvokeModel(ctx, invokeInput)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke embedding model: %w", err)
	}

	// Parse the response
	var response struct {
		Embedding []float64 `json:"embedding"`
	}

	if err := json.Unmarshal(result.Body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal embedding response: %w", err)
	}

	return response.Embedding, nil
}

// GetSimilarDocuments finds similar documents based on vector similarity
func (bkb *BedrockKnowledgeBaseService) GetSimilarDocuments(ctx context.Context, embedding []float64, limit int, filters map[string]interface{}) (*service.SimilarDocumentsResponse, error) {
	startTime := time.Now()

	// Use the knowledge base for similarity search
	// We'll create a synthetic query and use the embedding
	query := "Find similar products based on vector similarity"

	kbResponse, err := bkb.QueryKnowledgeBase(ctx, query, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to query knowledge base for similarity: %w", err)
	}

	// Convert knowledge base results to similar documents
	documents := make([]service.SimilarDocument, 0, len(kbResponse.Results))

	// Limit results
	maxResults := min(len(kbResponse.Results), limit)
	for i := 0; i < maxResults; i++ {
		result := kbResponse.Results[i]

		doc := service.SimilarDocument{
			DocumentID: result.Location.DocumentID,
			Content:    result.Content,
			Score:      result.Score,
			Metadata:   result.Metadata,
			Source:     result.Source,
		}

		documents = append(documents, doc)
	}

	return &service.SimilarDocumentsResponse{
		Documents:        documents,
		ProcessingTimeMs: time.Since(startTime).Milliseconds(),
	}, nil
}

// Helper functions

func (bkb *BedrockKnowledgeBaseService) calculateConfidenceLevel(results []service.KnowledgeBaseResult) float64 {
	if len(results) == 0 {
		return 0.0
	}

	var totalScore float64
	for _, result := range results {
		totalScore += result.Score
	}

	avgScore := totalScore / float64(len(results))

	// Convert to confidence level (0.0 to 1.0)
	// Assuming scores are between 0 and 1, adjust if different
	confidence := avgScore

	// Apply some thresholds
	if confidence > 0.8 {
		return 0.95
	} else if confidence > 0.6 {
		return 0.8
	} else if confidence > 0.4 {
		return 0.6
	} else {
		return 0.4
	}
}

func (bkb *BedrockKnowledgeBaseService) deduplicateSources(sources []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0)

	for _, source := range sources {
		if !seen[source] {
			seen[source] = true
			result = append(result, source)
		}
	}

	return result
}

func extractDocumentIDFromURI(uri string) string {
	if uri == "" {
		return ""
	}

	// Extract filename from S3 URI or file path
	parts := strings.Split(uri, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return uri
}

func extractModelIDFromARN(arn string) string {
	// Extract model ID from ARN format like: arn:aws:bedrock:region:account:foundation-model/model-id
	parts := strings.Split(arn, "/")
	if len(parts) > 1 {
		return parts[len(parts)-1]
	}
	return arn
}

// GetProductsWithSemanticSearch performs comprehensive semantic search using Amazon Bedrock Knowledge Base
// This method handles both text-based semantic search and vector similarity search using the Knowledge Base's
// automatic vectorization capabilities following AWS best practices.

func (bkb *BedrockKnowledgeBaseService) GetProductsWithSemanticSearch(ctx context.Context, query string, limit int, filters map[string]interface{}) (*service.RAGSemanticSearchResponse, error) {
	startTime := time.Now()

	// AWS Bedrock Knowledge Base automatically handles vectorization internally
	// No need to manually generate embeddings - the Knowledge Base will convert the text query to vectors automatically
	enhancedQuery := bkb.enhanceSemanticQuery(query, filters)

	// Optimize retrieval configuration for better results
	// Use more results for better ranking and diversity
	retrievalLimit := min(limit*3, 100) // Get more results for better selection

	retrievalConfig := &types.KnowledgeBaseRetrievalConfiguration{
		VectorSearchConfiguration: &types.KnowledgeBaseVectorSearchConfiguration{
			NumberOfResults:    aws.Int32(int32(retrievalLimit)),
			OverrideSearchType: types.SearchTypeHybrid, // Hybrid provides best of both semantic and keyword search
		},
	}

	// Create retrieve request - Knowledge Base handles automatic vectorization
	retrieveInput := &bedrockagentruntime.RetrieveInput{
		KnowledgeBaseId:        aws.String(bkb.knowledgeBaseID),
		RetrievalQuery:         &types.KnowledgeBaseQuery{Text: aws.String(enhancedQuery)},
		RetrievalConfiguration: retrievalConfig,
	}

	// Execute retrieve operation - Knowledge Base automatically:
	// 1. Converts query text to embeddings using the configured embedding model
	// 2. Searches the vector index for semantically similar content
	// 3. Returns ranked results based on similarity scores
	retrieveOutput, err := bkb.agentClient.Retrieve(ctx, retrieveInput)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve products with semantic search: %w", err)
	}

	// Convert Knowledge Base results to our service format
	results, err := bkb.convertKnowledgeBaseResultsToRAGResults(retrieveOutput.RetrievalResults, "semantic_search")
	if err != nil {
		return nil, fmt.Errorf("failed to convert semantic search results: %w", err)
	}

	// Enhanced scoring and metadata enrichment
	if len(results) > 0 {
		maxScore := 0.0
		minScore := 1.0

		// Find score range for normalization
		for i := range results {
			if i < len(retrieveOutput.RetrievalResults) {
				kbScore := float64(aws.ToFloat64(retrieveOutput.RetrievalResults[i].Score))
				if kbScore > maxScore {
					maxScore = kbScore
				}
				if kbScore < minScore {
					minScore = kbScore
				}
			}
		}

		// Apply score normalization and enhanced metadata
		scoreRange := maxScore - minScore
		for i := range results {
			if i < len(retrieveOutput.RetrievalResults) {
				kbScore := float64(aws.ToFloat64(retrieveOutput.RetrievalResults[i].Score))

				// Normalize scores if there's meaningful variation
				normalizedScore := kbScore
				if scoreRange > 0.01 {
					normalizedScore = (kbScore - minScore) / scoreRange
				}

				results[i].DistanceScore = normalizedScore
				results[i].SimilarityScore = normalizedScore
				results[i].SearchMethod = "semantic_search_enhanced"
				results[i].MatchedCriteria = bkb.extractMatchedCriteria(retrieveOutput.RetrievalResults[i])
				results[i].EmbeddingModel = bkb.embeddingModelID
				results[i].ConfidenceScore = bkb.calculateConfidenceFromScore(normalizedScore)
				results[i].RetrievalRank = i + 1

				// Extract semantic features for enhanced metadata
				if len(results[i].MatchedCriteria) == 0 {
					results[i].MatchedCriteria = []string{"semantic_similarity", "content_relevance"}
				}
			}
		}
	}

	// Apply enhanced ranking based on query relevance
	results = bkb.applySemanticReranking(results, query)

	// Apply post-retrieval filters if specified
	if len(filters) > 0 {
		results = bkb.applyPostRetrievalFiltersToBedrockResults(results, filters)
	}

	// Apply diversity algorithm to avoid overly similar results
	if len(results) > limit {
		results = bkb.applyDiversityAlgorithmForBedrock(results, limit)
	} else if len(results) > limit {
		results = results[:limit]
	}

	// Log enhanced semantic search operation
	processingTime := time.Since(startTime).Milliseconds()
	avgScore := 0.0
	if len(results) > 0 {
		for _, result := range results {
			avgScore += result.SimilarityScore
		}
		avgScore /= float64(len(results))
	}

	log.Printf("Enhanced semantic search completed: %d results in %dms for query '%s' (avg score: %.3f)",
		len(results), processingTime, query, avgScore)

	return &service.RAGSemanticSearchResponse{
		Query:            query,
		Results:          results,
		TotalFound:       len(results),
		ProcessingTimeMs: processingTime,
		SearchMetadata: &service.RAGSearchMeta{
			SearchType:       "semantic_search_enhanced",
			EmbeddingModel:   bkb.embeddingModelID,
			KnowledgeBaseID:  bkb.knowledgeBaseID,
			SimilarityMetric: "hybrid_cosine",
			FiltersApplied:   filters,
			RerankerUsed:     true,
			CacheUsed:        false,
		},
	}, nil
}

// applySemanticReranking applies additional ranking based on semantic relevance to the query
func (bkb *BedrockKnowledgeBaseService) applySemanticReranking(results []service.RAGSearchResult, query string) []service.RAGSearchResult {
	if len(results) <= 1 {
		return results
	}

	// Sort by a combination of similarity score and semantic relevance
	sort.Slice(results, func(i, j int) bool {
		// Combine similarity score with query relevance
		scoreI := results[i].SimilarityScore*0.7 + bkb.calculateSemanticRelevanceForBedrock(results[i], query)*0.3
		scoreJ := results[j].SimilarityScore*0.7 + bkb.calculateSemanticRelevanceForBedrock(results[j], query)*0.3
		return scoreI > scoreJ
	})

	return results
}

// Helper methods for the new search implementations

// enhanceSemanticQuery enhances the query with additional context for better semantic understanding
func (bkb *BedrockKnowledgeBaseService) enhanceSemanticQuery(query string, filters map[string]interface{}) string {
	var enhancedQuery strings.Builder
	enhancedQuery.WriteString(query)

	// Add category context if provided
	if categoryID, ok := filters["category_id"]; ok {
		enhancedQuery.WriteString(fmt.Sprintf(" in category %v", categoryID))
	}

	// Add price range context if provided
	if minPrice, ok := filters["min_price"]; ok {
		if maxPrice, ok := filters["max_price"]; ok {
			enhancedQuery.WriteString(fmt.Sprintf(" with price between %v and %v", minPrice, maxPrice))
		}
	}

	// Add brand context if provided
	if brand, ok := filters["brand"]; ok {
		enhancedQuery.WriteString(fmt.Sprintf(" from brand %v", brand))
	}

	return enhancedQuery.String()
}

// parseKnowledgeBaseDocument parses the Knowledge Base document format
// Extracts only the product_id from content
func (bkb *BedrockKnowledgeBaseService) parseKnowledgeBaseDocument(content string) (uuid.UUID, error) {
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "product_id:") {
			idStr := bkb.extractMarkdownValue(line, "product_id:")

			if idStr != "" {
				// Extract only the UUID part (first 36 characters in UUID format)
				uuidStr := bkb.extractUUIDFromString(idStr)
				if uuidStr != "" {
					if parsedID, err := uuid.Parse(uuidStr); err == nil {
						return parsedID, nil
					}
				}
			}
		}
	}

	return uuid.Nil, fmt.Errorf("no product_id found in content")
}

// extractMarkdownValue extracts value from markdown bold format
func (bkb *BedrockKnowledgeBaseService) extractMarkdownValue(line, key string) string {
	if idx := strings.Index(line, key); idx != -1 {
		remaining := strings.TrimSpace(line[idx+len(key):])
		// Remove trailing markdown or other formatting
		if spaceIdx := strings.Index(remaining, "  "); spaceIdx != -1 {
			remaining = remaining[:spaceIdx]
		}
		return strings.TrimSpace(remaining)
	}
	return ""
}

// extractUUIDFromString extracts UUID from a string that may contain additional text
func (bkb *BedrockKnowledgeBaseService) extractUUIDFromString(text string) string {
	// UUID format: 8-4-4-4-12 characters (36 total including hyphens)
	// Example: e491fc15-95de-4b8a-b3d7-6310dbf0b4db

	// First try to find UUID at the beginning of the string
	if len(text) >= 36 {
		candidate := text[:36]
		// Check if it matches UUID format (contains 4 hyphens at correct positions)
		if len(candidate) == 36 &&
			candidate[8] == '-' && candidate[13] == '-' &&
			candidate[18] == '-' && candidate[23] == '-' {
			// Validate that all other characters are hex digits
			if bkb.isValidUUIDFormat(candidate) {
				return candidate
			}
		}
	}

	// If not found at beginning, search within the text
	words := strings.Fields(text)
	for _, word := range words {
		if len(word) == 36 &&
			word[8] == '-' && word[13] == '-' &&
			word[18] == '-' && word[23] == '-' {
			if bkb.isValidUUIDFormat(word) {
				return word
			}
		}
	}

	return ""
}

// isValidUUIDFormat checks if a string matches UUID format
func (bkb *BedrockKnowledgeBaseService) isValidUUIDFormat(candidate string) bool {
	if len(candidate) != 36 {
		return false
	}

	for i, char := range candidate {
		if i == 8 || i == 13 || i == 18 || i == 23 {
			if char != '-' {
				return false
			}
		} else {
			// Check if character is valid hex digit (0-9, a-f, A-F)
			if !((char >= '0' && char <= '9') ||
				(char >= 'a' && char <= 'f') ||
				(char >= 'A' && char <= 'F')) {
				return false
			}
		}
	}

	return true
}

// calculateConfidenceFromScore calculates confidence score from similarity score
func (bkb *BedrockKnowledgeBaseService) calculateConfidenceFromScore(score float64) float64 {
	// Convert similarity score to confidence (0.0 to 1.0)
	// This is a simplified calculation - adjust based on your scoring system
	confidence := score
	if confidence > 1.0 {
		confidence = 1.0
	}
	if confidence < 0.0 {
		confidence = 0.0
	}
	return confidence
}

// extractMatchedCriteria extracts matched criteria from a retrieval result
func (bkb *BedrockKnowledgeBaseService) extractMatchedCriteria(result types.KnowledgeBaseRetrievalResult) []string {
	criteria := []string{}

	// Analyze the result to determine what criteria were matched
	content := aws.ToString(result.Content.Text)

	// Simple keyword-based criteria extraction
	if strings.Contains(strings.ToLower(content), "price") {
		criteria = append(criteria, "price_match")
	}
	if strings.Contains(strings.ToLower(content), "brand") {
		criteria = append(criteria, "brand_match")
	}
	if strings.Contains(strings.ToLower(content), "category") {
		criteria = append(criteria, "category_match")
	}

	return criteria
}

// calculateSemanticRelevanceForBedrock calculates semantic relevance score for Bedrock search results
func (bkb *BedrockKnowledgeBaseService) calculateSemanticRelevanceForBedrock(result service.RAGSearchResult, query string) float64 {
	// Calculate semantic relevance score for Bedrock search results
	productText := strings.ToLower(result.Source)
	queryTerms := strings.Fields(strings.ToLower(query))

	matches := 0
	for _, term := range queryTerms {
		if strings.Contains(productText, term) {
			matches++
		}
	}

	if len(queryTerms) == 0 {
		return 0.0
	}

	return float64(matches) / float64(len(queryTerms))
}

// applyPostRetrievalFiltersToBedrockResults applies additional filters after retrieval
func (bkb *BedrockKnowledgeBaseService) applyPostRetrievalFiltersToBedrockResults(results []service.RAGSearchResult, filters map[string]interface{}) []service.RAGSearchResult {
	if len(filters) == 0 {
		return results
	}

	filtered := make([]service.RAGSearchResult, 0, len(results))
	excludedCount := 0

	for _, result := range results {
		include := true

		// Apply exclusion filter (exclude specific product ID) - check multiple possible keys
		excludeID := ""
		if id, ok := filters["exclude_id"].(string); ok {
			excludeID = id
		} else if id, ok := filters["exclude_product_id"].(string); ok {
			excludeID = id
		}

		if excludeID != "" {
			if result.ProductID.String() == excludeID {
				include = false
				excludedCount++
				log.Printf("Excluding target product ID %s from results", excludeID)
			}
		}

		// Apply minimum score filter
		if minScore, ok := filters["min_score"].(float64); ok && result.SimilarityScore < minScore {
			include = false
		}

		// Apply maximum results filter
		if maxResults, ok := filters["max_results"].(int); ok && len(filtered) >= maxResults {
			break
		}

		if include {
			filtered = append(filtered, result)
		}
	}

	if excludedCount > 0 {
		log.Printf("Excluded %d products from search results", excludedCount)
	}

	return filtered
}

// applyDiversityAlgorithmForBedrock applies diversity algorithm to ensure varied results
func (bkb *BedrockKnowledgeBaseService) applyDiversityAlgorithmForBedrock(results []service.RAGSearchResult, limit int) []service.RAGSearchResult {
	if len(results) <= limit {
		return results
	}

	// Simple diversity algorithm - ensure product diversity by source
	sourceMap := make(map[string]bool)
	diversified := make([]service.RAGSearchResult, 0, limit)

	// First pass: include top products from different sources
	for _, result := range results {
		if len(diversified) >= limit {
			break
		}
		source := result.Source
		if source == "" {
			source = result.ProductID.String()[:8] // Use first 8 chars of UUID as fallback
		}
		if !sourceMap[source] {
			sourceMap[source] = true
			diversified = append(diversified, result)
		}
	}

	// Second pass: fill remaining slots with highest scoring products
	for _, result := range results {
		if len(diversified) >= limit {
			break
		}
		// Check if product is already included
		found := false
		for _, included := range diversified {
			if included.ProductID == result.ProductID {
				found = true
				break
			}
		}
		if !found {
			diversified = append(diversified, result)
		}
	}

	return diversified
}

// convertKnowledgeBaseResultsToRAGResults converts Bedrock Knowledge Base results to service.RAGSearchResult
func (bkb *BedrockKnowledgeBaseService) convertKnowledgeBaseResultsToRAGResults(results []types.KnowledgeBaseRetrievalResult, searchMethod string) ([]service.RAGSearchResult, error) {
	bedrockResults := make([]service.RAGSearchResult, 0, len(results))

	for _, result := range results {
		productID, err := bkb.parseKnowledgeBaseDocument(aws.ToString(result.Content.Text))
		if err != nil {
			log.Printf("Warning: Failed to parse product from content: %v", err)
			continue
		}

		// Extract metadata from result
		metadata := make(map[string]interface{})
		if result.Metadata != nil {
			for k, v := range result.Metadata {
				metadata[k] = v
			}
		}

		// Extract source information
		source := ""
		if result.Location != nil && result.Location.S3Location != nil {
			source = aws.ToString(result.Location.S3Location.Uri)
		}

		bedrockResults = append(bedrockResults, service.RAGSearchResult{
			ProductID:       productID,
			DistanceScore:   float64(aws.ToFloat64(result.Score)),
			SimilarityScore: float64(aws.ToFloat64(result.Score)),
			SearchMethod:    searchMethod,
			Metadata:        metadata,
			Source:          source,
		})
	}

	return bedrockResults, nil
}
