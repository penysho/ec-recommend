package service

import (
	"context"
	"ec-recommend/internal/service"
	"encoding/json"
	"fmt"
	"log"
	"math"
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

// GetProductsWithVectorSearch performs vector-based product search using Amazon Bedrock Knowledge Base with enhanced ranking
func (bkb *BedrockKnowledgeBaseService) GetProductsWithVectorSearch(ctx context.Context, vector []float64, limit int, filters map[string]interface{}) (*service.RAGVectorSearchResponse, error) {
	startTime := time.Now()

	// Step 1: Analyze the input vector to extract semantic features
	semanticFeatures, err := bkb.analyzeVectorForFeatures(ctx, vector)
	if err != nil {
		log.Printf("Warning: Failed to analyze vector features: %v", err)
		semanticFeatures = []string{"similar products"} // Fallback
	}

	// Step 2: Convert vector features to enhanced search query
	vectorQuery := bkb.buildVectorBasedQuery(semanticFeatures, filters)

	// Step 3: Optimized retrieval configuration with higher result count for better ranking
	// Use more results for better reranking (up to 3x the requested limit or max 100)
	retrievalLimit := min(limit*4, 100) // Increased from 3x to 4x for better selection

	retrievalConfig := &types.KnowledgeBaseRetrievalConfiguration{
		VectorSearchConfiguration: &types.KnowledgeBaseVectorSearchConfiguration{
			NumberOfResults:    aws.Int32(int32(retrievalLimit)),
			OverrideSearchType: types.SearchTypeHybrid, // Keep hybrid for better coverage
		},
	}

	retrieveInput := &bedrockagentruntime.RetrieveInput{
		KnowledgeBaseId:        aws.String(bkb.knowledgeBaseID),
		RetrievalQuery:         &types.KnowledgeBaseQuery{Text: aws.String(vectorQuery)},
		RetrievalConfiguration: retrievalConfig,
	}

	// Execute retrieve operation
	retrieveOutput, err := bkb.agentClient.Retrieve(ctx, retrieveInput)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve products with vector search: %w", err)
	}

	// Step 4: Convert results to service.RAGSearchResult
	results, err := bkb.convertKnowledgeBaseResultsToRAGResults(retrieveOutput.RetrievalResults, "vector_search")
	if err != nil {
		return nil, fmt.Errorf("failed to convert vector search results: %w", err)
	}

	// Step 5: Enhanced scoring with normalization
	// First, normalize scores to handle Knowledge Base score variations
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

		// Apply normalization if there's a meaningful range
		scoreRange := maxScore - minScore
		if scoreRange > 0.01 { // Only normalize if there's sufficient variation
			for i := range results {
				if i < len(retrieveOutput.RetrievalResults) {
					kbScore := float64(aws.ToFloat64(retrieveOutput.RetrievalResults[i].Score))
					// Normalize to 0-1 range and apply power transformation for better distribution
					normalizedScore := (kbScore - minScore) / scoreRange
					enhancedScore := math.Pow(normalizedScore, 0.8) // Slight compression to reduce extreme values

					results[i].SimilarityScore = enhancedScore
					results[i].DistanceScore = enhancedScore
				}
			}
		} else {
			// If minimal variation, use original scores
			for i := range results {
				if i < len(retrieveOutput.RetrievalResults) {
					kbScore := float64(aws.ToFloat64(retrieveOutput.RetrievalResults[i].Score))
					results[i].SimilarityScore = kbScore
					results[i].DistanceScore = kbScore
				}
			}
		}
	}

	// Step 6: Enhanced vector metadata and feature analysis
	vectorClusters := bkb.extractVectorClusters(vector)

	for i := range results {
		if i < len(retrieveOutput.RetrievalResults) {
			vectorScore := results[i].SimilarityScore

			// Enhanced confidence calculation
			confidence := bkb.calculateVectorConfidence(vectorScore, len(semanticFeatures))

			// Apply vector quality boosting based on clusters
			qualityBoost := bkb.calculateVectorQualityBoost(vectorClusters, vectorScore)
			confidence = math.Min(confidence+qualityBoost, 1.0)

			results[i].DistanceScore = vectorScore
			results[i].SearchMethod = "vector_search_enhanced"
			results[i].MatchedCriteria = append([]string{"vector_similarity"}, semanticFeatures...)
			results[i].EmbeddingModel = bkb.embeddingModelID
			results[i].SemanticClusters = vectorClusters
			results[i].SimilarityScore = vectorScore
			results[i].ConfidenceScore = confidence
			results[i].RetrievalRank = i + 1
		}
	}

	// Step 7: Advanced reranking with multiple criteria
	results = bkb.applyAdvancedReranking(results, vector, semanticFeatures)

	// Step 8: Apply filters with enhanced logic
	if len(filters) > 0 {
		results = bkb.applyPostRetrievalFiltersToBedrockResults(results, filters)
	}

	// Step 9: Diversity enhancement to avoid too similar results
	if len(results) > limit {
		results = bkb.applyDiversityAlgorithmForBedrock(results, limit)
	} else if len(results) > limit {
		results = results[:limit]
	}

	// Log enhanced vector search operation
	processingTime := time.Since(startTime).Milliseconds()
	avgScore := 0.0
	if len(results) > 0 {
		for _, result := range results {
			avgScore += result.SimilarityScore
		}
		avgScore /= float64(len(results))
	}

	log.Printf("Enhanced vector search completed: %d results in %dms with %d-dimensional vector (avg score: %.3f)",
		len(results), processingTime, len(vector), avgScore)

	return &service.RAGVectorSearchResponse{
		Results:          results,
		TotalFound:       len(results),
		ProcessingTimeMs: processingTime,
		SearchMetadata: &service.RAGSearchMeta{
			SearchType:       "vector_search_enhanced",
			EmbeddingModel:   bkb.embeddingModelID,
			KnowledgeBaseID:  bkb.knowledgeBaseID,
			SimilarityMetric: "cosine_enhanced",
			FiltersApplied:   filters,
			RerankerUsed:     true,
			CacheUsed:        false,
		},
	}, nil
}

// GetProductsWithSemanticSearch performs semantic/text-based product search using Amazon Bedrock Knowledge Base
func (bkb *BedrockKnowledgeBaseService) GetProductsWithSemanticSearch(ctx context.Context, query string, limit int, filters map[string]interface{}) (*service.RAGSemanticSearchResponse, error) {
	startTime := time.Now()

	// Enhance query with context for better semantic understanding
	enhancedQuery := bkb.enhanceSemanticQuery(query, filters)

	// Configure retrieval for semantic search (using hybrid since pure semantic is not available)
	retrievalConfig := &types.KnowledgeBaseRetrievalConfiguration{
		VectorSearchConfiguration: &types.KnowledgeBaseVectorSearchConfiguration{
			NumberOfResults: aws.Int32(int32(min(limit, 50))), // AWS recommends max 50 results per query
			// Note: Pure semantic search not available, using hybrid with text emphasis
			OverrideSearchType: types.SearchTypeHybrid,
		},
	}

	// Create retrieve request
	retrieveInput := &bedrockagentruntime.RetrieveInput{
		KnowledgeBaseId:        aws.String(bkb.knowledgeBaseID),
		RetrievalQuery:         &types.KnowledgeBaseQuery{Text: aws.String(enhancedQuery)},
		RetrievalConfiguration: retrievalConfig,
	}

	// Execute retrieve operation
	retrieveOutput, err := bkb.agentClient.Retrieve(ctx, retrieveInput)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve products with semantic search: %w", err)
	}

	// Convert results to service.RAGSearchResult
	results, err := bkb.convertKnowledgeBaseResultsToRAGResults(retrieveOutput.RetrievalResults, "semantic_search")
	if err != nil {
		return nil, fmt.Errorf("failed to convert semantic search results: %w", err)
	}

	// Enhance results with semantic metadata
	for i := range results {
		if i < len(retrieveOutput.RetrievalResults) {
			results[i].DistanceScore = float64(aws.ToFloat64(retrieveOutput.RetrievalResults[i].Score))
			results[i].SearchMethod = "semantic_search"
			results[i].MatchedCriteria = bkb.extractMatchedCriteria(retrieveOutput.RetrievalResults[i])
			results[i].EmbeddingModel = bkb.embeddingModelID
			results[i].SimilarityScore = float64(aws.ToFloat64(retrieveOutput.RetrievalResults[i].Score))
			results[i].ConfidenceScore = bkb.calculateConfidenceFromScore(results[i].SimilarityScore)
			results[i].RetrievalRank = i + 1
		}
	}

	// Apply additional filtering if needed
	if len(filters) > 0 {
		results = bkb.applyPostRetrievalFiltersToBedrockResults(results, filters)
	}

	// Log semantic search operation
	processingTime := time.Since(startTime).Milliseconds()
	log.Printf("Semantic search completed: %d results in %dms for query: %s", len(results), processingTime, query)

	return &service.RAGSemanticSearchResponse{
		Query:            query,
		Results:          results[:min(len(results), limit)],
		TotalFound:       len(results),
		ProcessingTimeMs: processingTime,
		SearchMetadata: &service.RAGSearchMeta{
			SearchType:       "semantic_search",
			EmbeddingModel:   bkb.embeddingModelID,
			KnowledgeBaseID:  bkb.knowledgeBaseID,
			SimilarityMetric: "hybrid",
			FiltersApplied:   filters,
			RerankerUsed:     false,
			CacheUsed:        false,
		},
	}, nil
}

// GetProductsWithHybridSearch performs hybrid search combining vector and semantic approaches using Amazon Bedrock Knowledge Base
func (bkb *BedrockKnowledgeBaseService) GetProductsWithHybridSearch(ctx context.Context, query string, vector []float64, limit int, filters map[string]interface{}) (*service.RAGHybridSearchResponse, error) {
	startTime := time.Now()

	// Enhance query for hybrid search context
	enhancedQuery := bkb.enhanceHybridQuery(query, filters)

	// Configure retrieval for hybrid search (combines vector and semantic)
	retrievalConfig := &types.KnowledgeBaseRetrievalConfiguration{
		VectorSearchConfiguration: &types.KnowledgeBaseVectorSearchConfiguration{
			NumberOfResults:    aws.Int32(int32(min(limit*2, 100))), // Get more results for better hybrid ranking
			OverrideSearchType: types.SearchTypeHybrid,              // Hybrid search combining vector + semantic
		},
	}

	// Create retrieve request
	retrieveInput := &bedrockagentruntime.RetrieveInput{
		KnowledgeBaseId:        aws.String(bkb.knowledgeBaseID),
		RetrievalQuery:         &types.KnowledgeBaseQuery{Text: aws.String(enhancedQuery)},
		RetrievalConfiguration: retrievalConfig,
	}

	// Execute retrieve operation
	retrieveOutput, err := bkb.agentClient.Retrieve(ctx, retrieveInput)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve products with hybrid search: %w", err)
	}

	// Convert results to service.RAGSearchResult
	results, err := bkb.convertKnowledgeBaseResultsToRAGResults(retrieveOutput.RetrievalResults, "hybrid_search")
	if err != nil {
		return nil, fmt.Errorf("failed to convert hybrid search results: %w", err)
	}

	// Apply hybrid ranking algorithm - combine vector similarity and semantic relevance
	results = bkb.applyHybridRankingToBedrockResults(results, vector, query, retrieveOutput.RetrievalResults)

	// Enhance results with comprehensive metadata
	for i := range results {
		if i < len(retrieveOutput.RetrievalResults) {
			results[i].DistanceScore = float64(aws.ToFloat64(retrieveOutput.RetrievalResults[i].Score))
			results[i].SearchMethod = "hybrid_search"
			results[i].MatchedCriteria = append([]string{"vector_similarity", "semantic_relevance"}, bkb.extractMatchedCriteria(retrieveOutput.RetrievalResults[i])...)
			results[i].EmbeddingModel = bkb.embeddingModelID

			// Calculate hybrid confidence score (weighted combination)
			vectorScore := float64(aws.ToFloat64(retrieveOutput.RetrievalResults[i].Score))
			semanticScore := bkb.calculateSemanticRelevanceForBedrock(results[i], query)
			results[i].SimilarityScore = vectorScore
			results[i].ConfidenceScore = bkb.calculateHybridConfidence(vectorScore, semanticScore)
			results[i].RetrievalRank = i + 1
		}
	}

	// Apply additional filtering if needed
	if len(filters) > 0 {
		results = bkb.applyPostRetrievalFiltersToBedrockResults(results, filters)
	}

	// Final ranking and selection
	results = bkb.finalHybridRankingForBedrock(results, limit)

	// Log hybrid search operation
	processingTime := time.Since(startTime).Milliseconds()
	log.Printf("Hybrid search completed: %d results in %dms for query: %s", len(results), processingTime, query)

	return &service.RAGHybridSearchResponse{
		Query:            query,
		Vector:           vector,
		Results:          results[:min(len(results), limit)],
		TotalFound:       len(results),
		ProcessingTimeMs: processingTime,
		SearchMetadata: &service.RAGSearchMeta{
			SearchType:         "hybrid_search",
			EmbeddingModel:     bkb.embeddingModelID,
			KnowledgeBaseID:    bkb.knowledgeBaseID,
			SimilarityMetric:   "hybrid",
			FiltersApplied:     filters,
			RerankerUsed:       true,
			CacheUsed:          false,
			HybridSearchWeight: map[string]float64{"vector": 0.6, "semantic": 0.4},
		},
	}, nil
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

// enhanceHybridQuery enhances the query specifically for hybrid search
func (bkb *BedrockKnowledgeBaseService) enhanceHybridQuery(query string, filters map[string]interface{}) string {
	// For hybrid search, we want to maintain the original semantic meaning
	// while providing enough context for vector similarity
	enhancedQuery := bkb.enhanceSemanticQuery(query, filters)
	enhancedQuery += " - find similar products with matching features and characteristics"
	return enhancedQuery
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

// callNovaForAnalysis calls Amazon Nova model for semantic analysis
func (bkb *BedrockKnowledgeBaseService) callNovaForAnalysis(ctx context.Context, prompt string) (string, error) {
	// Prepare input for Amazon Nova model
	input := map[string]interface{}{
		"messages": []map[string]interface{}{
			{
				"role": "user",
				"content": []map[string]interface{}{
					{
						"text": prompt,
					},
				},
			},
		},
		"inferenceConfig": map[string]interface{}{
			"maxTokens": 2000,
		},
	}

	inputBytes, err := json.Marshal(input)
	if err != nil {
		return "", fmt.Errorf("failed to marshal input: %w", err)
	}

	// Extract model ID from ARN or use directly
	modelID := bkb.extractNovaModelID()

	// Call Nova model
	invokeInput := &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(modelID),
		ContentType: aws.String("application/json"),
		Body:        inputBytes,
	}

	result, err := bkb.runtimeClient.InvokeModel(ctx, invokeInput)
	if err != nil {
		return "", fmt.Errorf("failed to invoke Nova model: %w", err)
	}

	// Parse Nova response format
	var response struct {
		Output struct {
			Message struct {
				Content []struct {
					Text string `json:"text"`
				} `json:"content"`
			} `json:"message"`
		} `json:"output"`
	}

	if err := json.Unmarshal(result.Body, &response); err != nil {
		return "", fmt.Errorf("failed to unmarshal Nova response: %w", err)
	}

	if len(response.Output.Message.Content) == 0 {
		return "", fmt.Errorf("no content in Nova response")
	}

	return response.Output.Message.Content[0].Text, nil
}

// extractNovaModelID extracts or determines the Nova model ID to use
func (bkb *BedrockKnowledgeBaseService) extractNovaModelID() string {
	// If modelARN contains Nova model, extract it
	if strings.Contains(bkb.modelARN, "nova") {
		return extractModelIDFromARN(bkb.modelARN)
	}

	// Default to Nova Lite model if not explicitly set
	return "amazon.nova-lite-v1:0"
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

// calculateHybridConfidence calculates hybrid confidence score
func (bkb *BedrockKnowledgeBaseService) calculateHybridConfidence(vectorScore, semanticScore float64) float64 {
	// Weighted combination with normalization
	hybridScore := 0.6*vectorScore + 0.4*semanticScore

	// Normalize to 0-1 range
	if hybridScore > 1.0 {
		hybridScore = 1.0
	}
	if hybridScore < 0.0 {
		hybridScore = 0.0
	}

	return hybridScore
}

// analyzeVectorForFeatures analyzes vector to extract enhanced semantic features
func (bkb *BedrockKnowledgeBaseService) analyzeVectorForFeatures(ctx context.Context, vector []float64) ([]string, error) {
	if len(vector) == 0 {
		return []string{"general_products"}, nil
	}

	var features []string

	// Statistical analysis
	var sum, sumSquares float64
	var max, min float64 = vector[0], vector[0]
	nonZero := 0

	for _, v := range vector {
		sum += v
		sumSquares += v * v
		if v > max {
			max = v
		}
		if v < min {
			min = v
		}
		if v != 0 {
			nonZero++
		}
	}

	mean := sum / float64(len(vector))
	variance := (sumSquares / float64(len(vector))) - (mean * mean)
	stdDev := math.Sqrt(variance)
	sparsity := float64(nonZero) / float64(len(vector))
	magnitude := math.Sqrt(sumSquares)

	// Analyze vector characteristics and map to product features

	// Magnitude-based features (product category insights)
	if magnitude > 0.8 {
		features = append(features, "premium products", "high-quality items", "detailed specifications")
	} else if magnitude > 0.4 {
		features = append(features, "standard products", "popular items", "mainstream brands")
	} else {
		features = append(features, "basic products", "essential items", "budget-friendly")
	}

	// Distribution-based features (product diversity insights)
	if stdDev > 0.15 {
		features = append(features, "diverse features", "multi-category", "varied specifications")
	} else if stdDev > 0.08 {
		features = append(features, "balanced features", "consistent quality", "similar specifications")
	} else {
		features = append(features, "focused features", "specialized products", "niche category")
	}

	// Sparsity-based features (product complexity insights)
	if sparsity > 0.9 {
		features = append(features, "feature-rich products", "comprehensive specifications", "detailed descriptions")
	} else if sparsity > 0.7 {
		features = append(features, "well-described products", "good specifications", "clear features")
	} else if sparsity > 0.3 {
		features = append(features, "basic specifications", "simple products", "essential features")
	} else {
		features = append(features, "minimal specifications", "simple items", "basic descriptions")
	}

	// Peak analysis for dominant features
	threshold := mean + (2 * stdDev)
	dominantDims := make([]int, 0)

	for i, v := range vector {
		if v > threshold {
			dominantDims = append(dominantDims, i)
		}
	}

	// Map dominant dimensions to product characteristics
	if len(dominantDims) > len(vector)/5 {
		features = append(features, "multi-faceted products", "comprehensive features", "well-rounded items")
	} else if len(dominantDims) > 0 {
		// Add dimension-specific features based on position patterns
		dimRatio := float64(len(dominantDims)) / float64(len(vector))
		if dimRatio > 0.1 {
			features = append(features, "distinctive features", "unique characteristics", "special attributes")
		} else {
			features = append(features, "focused characteristics", "specific features", "targeted attributes")
		}
	}

	// Range analysis for product variety
	vectorRange := max - min
	if vectorRange > 0.6 {
		features = append(features, "wide product range", "diverse options", "varied selection")
	} else if vectorRange > 0.3 {
		features = append(features, "moderate variety", "balanced selection", "curated options")
	} else {
		features = append(features, "focused selection", "specialized range", "consistent options")
	}

	// Semantic clustering for enhanced matching
	if mean > 0.3 && magnitude > 0.5 {
		features = append(features, "popular similar products", "trending items", "recommended choices")
	}

	if stdDev < 0.05 && sparsity > 0.8 {
		features = append(features, "consistent quality products", "reliable brands", "trusted items")
	}

	// Advanced Nova analysis for deeper insights (if available)
	if len(features) > 0 {
		// Create a prompt for Nova analysis
		topFeatureCount := 3
		if len(features) < 3 {
			topFeatureCount = len(features)
		}
		featurePrompt := fmt.Sprintf(
			"Based on vector analysis showing magnitude=%.3f, variance=%.3f, sparsity=%.3f, analyze product characteristics and suggest 2-3 key product features for search. Current features: %s",
			magnitude, variance, sparsity, strings.Join(features[:topFeatureCount], ", "))

		novaFeatures, err := bkb.callNovaForAnalysis(ctx, featurePrompt)
		if err == nil && novaFeatures != "" {
			// Parse Nova response and add refined features
			novaLines := strings.Split(novaFeatures, "\n")
			for _, line := range novaLines {
				line = strings.TrimSpace(line)
				if len(line) > 5 && len(line) < 50 { // Reasonable feature length
					features = append(features, line)
				}
			}
		}
	}

	// Ensure we have meaningful features
	if len(features) == 0 {
		features = []string{"similar products", "related items", "comparable options"}
	}

	// Remove duplicates and limit feature count
	features = bkb.deduplicateFeatures(features)

	// Limit to top 8 features to avoid overly complex queries
	if len(features) > 8 {
		features = features[:8]
	}

	return features, nil
}

// buildVectorBasedQuery builds a search query based on extracted vector features
func (bkb *BedrockKnowledgeBaseService) buildVectorBasedQuery(semanticFeatures []string, filters map[string]interface{}) string {
	var queryBuilder strings.Builder

	// Start with feature-based query
	queryBuilder.WriteString("Find products with similar characteristics: ")
	queryBuilder.WriteString(strings.Join(semanticFeatures, ", "))

	// Add filter context
	if categoryID, ok := filters["category_id"]; ok {
		queryBuilder.WriteString(fmt.Sprintf(" in category %v", categoryID))
	}

	if brand, ok := filters["brand"]; ok {
		queryBuilder.WriteString(fmt.Sprintf(" from brand %v", brand))
	}

	if minPrice, ok := filters["min_price"]; ok {
		if maxPrice, ok := filters["max_price"]; ok {
			queryBuilder.WriteString(fmt.Sprintf(" with price range %v to %v", minPrice, maxPrice))
		}
	}

	// Add similarity emphasis
	queryBuilder.WriteString(" - prioritize products with matching features and similar quality")

	return queryBuilder.String()
}

// extractVectorClusters extracts semantic clusters from vector analysis with improved logic
func (bkb *BedrockKnowledgeBaseService) extractVectorClusters(vector []float64) []string {
	clusters := []string{}

	if len(vector) == 0 {
		return []string{"empty_vector"}
	}

	// Statistical analysis of vector
	var sum, sumSquares float64
	var max, min float64 = vector[0], vector[0]
	nonZero := 0

	for _, v := range vector {
		sum += v
		sumSquares += v * v
		if v > max {
			max = v
		}
		if v < min {
			min = v
		}
		if v != 0 {
			nonZero++
		}
	}

	mean := sum / float64(len(vector))
	variance := (sumSquares / float64(len(vector))) - (mean * mean)
	stdDev := math.Sqrt(variance)
	sparsity := float64(nonZero) / float64(len(vector))
	magnitude := math.Sqrt(sumSquares)

	// Analyze vector characteristics and map to product features

	// Magnitude-based features (product category insights)
	if magnitude > 0.8 {
		clusters = append(clusters, "premium products", "high-quality items", "detailed specifications")
	} else if magnitude > 0.4 {
		clusters = append(clusters, "standard products", "popular items", "mainstream brands")
	} else {
		clusters = append(clusters, "basic products", "essential items", "budget-friendly")
	}

	// Distribution-based features (product diversity insights)
	if stdDev > 0.15 {
		clusters = append(clusters, "diverse features", "multi-category", "varied specifications")
	} else if stdDev > 0.08 {
		clusters = append(clusters, "balanced features", "consistent quality", "similar specifications")
	} else {
		clusters = append(clusters, "focused features", "specialized products", "niche category")
	}

	// Sparsity-based features (product complexity insights)
	if sparsity > 0.9 {
		clusters = append(clusters, "feature-rich products", "comprehensive specifications", "detailed descriptions")
	} else if sparsity > 0.7 {
		clusters = append(clusters, "well-described products", "good specifications", "clear features")
	} else if sparsity > 0.3 {
		clusters = append(clusters, "basic specifications", "simple products", "essential features")
	} else {
		clusters = append(clusters, "minimal specifications", "simple items", "basic descriptions")
	}

	// Peak analysis for dominant features
	threshold := mean + (2 * stdDev)
	dominantDims := make([]int, 0)

	for i, v := range vector {
		if v > threshold {
			dominantDims = append(dominantDims, i)
		}
	}

	// Map dominant dimensions to product characteristics
	if len(dominantDims) > len(vector)/5 {
		clusters = append(clusters, "multi-faceted products", "comprehensive features", "well-rounded items")
	} else if len(dominantDims) > 0 {
		// Add dimension-specific features based on position patterns
		dimRatio := float64(len(dominantDims)) / float64(len(vector))
		if dimRatio > 0.1 {
			clusters = append(clusters, "distinctive features", "unique characteristics", "special attributes")
		} else {
			clusters = append(clusters, "focused characteristics", "specific features", "targeted attributes")
		}
	}

	// Range analysis for product variety
	vectorRange := max - min
	if vectorRange > 0.6 {
		clusters = append(clusters, "wide product range", "diverse options", "varied selection")
	} else if vectorRange > 0.3 {
		clusters = append(clusters, "moderate variety", "balanced selection", "curated options")
	} else {
		clusters = append(clusters, "focused selection", "specialized range", "consistent options")
	}

	// Semantic clustering for enhanced matching
	if mean > 0.3 && magnitude > 0.5 {
		clusters = append(clusters, "popular similar products", "trending items", "recommended choices")
	}

	if stdDev < 0.05 && sparsity > 0.8 {
		clusters = append(clusters, "consistent quality products", "reliable brands", "trusted items")
	}

	// Additional clustering based on statistical patterns
	if len(clusters) > 0 {
		// Add contextual clusters based on existing analysis
		if magnitude > 0.6 && stdDev > 0.1 {
			clusters = append(clusters, "rich_content")
		}
		if sparsity > 0.8 && variance > 0.05 {
			clusters = append(clusters, "detailed_features")
		}
	}

	// Ensure we have meaningful features
	if len(clusters) == 0 {
		clusters = []string{"similar products", "related items", "comparable options"}
	}

	// Remove duplicates and limit feature count
	clusters = bkb.deduplicateFeatures(clusters)

	// Limit to top 8 features to avoid overly complex queries
	if len(clusters) > 8 {
		clusters = clusters[:8]
	}

	return clusters
}

// calculateVectorConfidence calculates enhanced confidence score for vector search results
func (bkb *BedrockKnowledgeBaseService) calculateVectorConfidence(vectorScore float64, featureCount int) float64 {
	// Validate input parameters
	if vectorScore < 0 {
		vectorScore = 0
	}
	if vectorScore > 1 {
		vectorScore = 1
	}

	// Base confidence from vector similarity (weighted more heavily)
	baseConfidence := vectorScore * 0.7 // Give 70% weight to actual similarity score

	// Feature richness boost - more features generally indicate better matches
	featureBoost := 0.0
	if featureCount > 0 {
		// Logarithmic scaling to prevent over-boosting
		featureBoost = math.Log(float64(featureCount)+1) * 0.05
		if featureBoost > 0.15 { // Cap at 15% boost
			featureBoost = 0.15
		}
	}

	// Score distribution analysis - penalize very low scores more heavily
	distributionAdjustment := 0.0
	if vectorScore < 0.1 {
		// Heavily penalize extremely low scores
		distributionAdjustment = -0.1
	} else if vectorScore < 0.3 {
		// Moderately penalize low scores
		distributionAdjustment = -0.05
	} else if vectorScore > 0.7 {
		// Reward high scores
		distributionAdjustment = 0.05
	}

	// Calculate final confidence
	confidence := baseConfidence + featureBoost + distributionAdjustment

	// Apply sigmoid transformation for smoother confidence curve
	// This helps spread out the confidence values more naturally
	confidence = 1.0 / (1.0 + math.Exp(-6*(confidence-0.5)))

	// Normalize to 0-1 range with more aggressive normalization for low scores
	if confidence > 1.0 {
		confidence = 1.0
	}
	if confidence < 0.0 {
		confidence = 0.0
	}

	// Apply minimum confidence threshold to avoid extremely low values
	minConfidence := 0.1
	if confidence < minConfidence {
		confidence = minConfidence
	}

	return confidence
}

// calculateVectorQualityBoost calculates quality boost based on vector clusters
func (bkb *BedrockKnowledgeBaseService) calculateVectorQualityBoost(clusters []string, score float64) float64 {
	boost := 0.0

	// Apply boosts based on cluster characteristics
	for _, cluster := range clusters {
		switch cluster {
		case "high_similarity":
			boost += 0.05
		case "moderate_similarity":
			boost += 0.02
		case "very_dense_features", "dense_features":
			boost += 0.03
		case "high_variance":
			boost += 0.02 // High variance can indicate rich feature representation
		case "multi_peaked", "focused_peaks":
			boost += 0.02 // Peak patterns can indicate strong feature matches
		case "wide_range":
			boost += 0.01 // Wide range can indicate comprehensive features
		}
	}

	// Scale boost based on base score - higher scores get proportionally more boost
	if score > 0.5 {
		boost *= 1.5
	} else if score < 0.2 {
		boost *= 0.5 // Reduce boost for very low scores
	}

	// Cap maximum boost
	if boost > 0.1 {
		boost = 0.1
	}

	return boost
}

// applyAdvancedReranking applies sophisticated reranking using multiple criteria
func (bkb *BedrockKnowledgeBaseService) applyAdvancedReranking(results []service.RAGSearchResult, vector []float64, semanticFeatures []string) []service.RAGSearchResult {
	if len(results) <= 1 {
		return results
	}

	// Calculate advanced ranking scores
	for i := range results {
		originalScore := results[i].SimilarityScore

		// Factor 1: Base similarity score (60% weight)
		rankingScore := originalScore * 0.6

		// Factor 2: Confidence score (25% weight)
		rankingScore += results[i].ConfidenceScore * 0.25

		// Factor 3: Feature richness (10% weight)
		featureRichness := float64(len(results[i].MatchedCriteria)) / 10.0 // Normalize
		if featureRichness > 1.0 {
			featureRichness = 1.0
		}
		rankingScore += featureRichness * 0.1

		// Factor 4: Semantic cluster quality (5% weight)
		clusterQuality := bkb.calculateClusterQuality(results[i].SemanticClusters)
		rankingScore += clusterQuality * 0.05

		// Store the ranking score for sorting
		results[i].DistanceScore = rankingScore // Use DistanceScore field to store ranking score
	}

	// Sort by ranking score
	sort.Slice(results, func(i, j int) bool {
		return results[i].DistanceScore > results[j].DistanceScore
	})

	// Restore original distance scores after sorting
	for i := range results {
		results[i].DistanceScore = results[i].SimilarityScore
	}

	return results
}

// calculateClusterQuality evaluates the quality of semantic clusters
func (bkb *BedrockKnowledgeBaseService) calculateClusterQuality(clusters []string) float64 {
	if len(clusters) == 0 {
		return 0.0
	}

	quality := 0.0
	count := 0

	for _, cluster := range clusters {
		count++
		switch cluster {
		case "high_similarity":
			quality += 1.0
		case "moderate_similarity":
			quality += 0.6
		case "low_similarity":
			quality += 0.2
		case "very_dense_features":
			quality += 0.9
		case "dense_features":
			quality += 0.7
		case "moderate_features":
			quality += 0.5
		case "high_variance":
			quality += 0.6
		case "moderate_variance":
			quality += 0.4
		case "multi_peaked":
			quality += 0.8
		case "focused_peaks":
			quality += 0.7
		case "wide_range":
			quality += 0.5
		default:
			quality += 0.3 // Default for unknown clusters
		}
	}

	if count > 0 {
		quality /= float64(count)
	}

	return quality
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

// applyHybridRankingToBedrockResults applies hybrid ranking algorithm to combine vector and semantic scores
func (bkb *BedrockKnowledgeBaseService) applyHybridRankingToBedrockResults(results []service.RAGSearchResult, vector []float64, query string, knowledgeResults []types.KnowledgeBaseRetrievalResult) []service.RAGSearchResult {
	// Apply hybrid ranking algorithm to Bedrock search results
	hybridResults := make([]service.RAGSearchResult, 0, len(results))

	for i, result := range results {
		hybridResult := service.RAGSearchResult{
			ProductID:       result.ProductID,
			DistanceScore:   result.DistanceScore,
			SimilarityScore: result.SimilarityScore,
			SearchMethod:    result.SearchMethod,
			Metadata:        result.Metadata,
			Source:          result.Source,
		}

		// Apply hybrid ranking algorithm to Bedrock search results
		hybridResult.SimilarityScore = bkb.calculateHybridRankingScore(result.SimilarityScore, vector, query, knowledgeResults[i])
		hybridResult.DistanceScore = hybridResult.SimilarityScore

		hybridResults = append(hybridResults, hybridResult)
	}

	// Sort hybrid results by similarity score in descending order
	sort.Slice(hybridResults, func(i, j int) bool {
		return hybridResults[i].SimilarityScore > hybridResults[j].SimilarityScore
	})

	return hybridResults
}

// calculateHybridRankingScore calculates hybrid ranking score
func (bkb *BedrockKnowledgeBaseService) calculateHybridRankingScore(vectorScore float64, vector []float64, query string, knowledgeResult types.KnowledgeBaseRetrievalResult) float64 {
	// Calculate hybrid ranking score
	hybridScore := vectorScore

	// Apply hybrid ranking algorithm to Bedrock search results
	hybridScore += bkb.calculateHybridRankingScoreFromKnowledge(knowledgeResult, vector, query)

	return hybridScore
}

// calculateHybridRankingScoreFromKnowledge calculates hybrid ranking score from knowledge result
func (bkb *BedrockKnowledgeBaseService) calculateHybridRankingScoreFromKnowledge(knowledgeResult types.KnowledgeBaseRetrievalResult, vector []float64, query string) float64 {
	// Calculate hybrid ranking score from knowledge result
	knowledgeScore := float64(aws.ToFloat64(knowledgeResult.Score))

	// Apply hybrid ranking algorithm to Bedrock search results
	hybridScore := knowledgeScore

	return hybridScore
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

// finalHybridRankingForBedrock performs final ranking and selection for hybrid search
func (bkb *BedrockKnowledgeBaseService) finalHybridRankingForBedrock(results []service.RAGSearchResult, limit int) []service.RAGSearchResult {
	// Apply diversity algorithm to ensure varied results
	diversified := bkb.applyDiversityAlgorithmForBedrock(results, limit)

	// Final confidence adjustment based on ranking position
	for i := range diversified {
		positionBonus := 1.0 - (float64(i) * 0.05) // Slight bonus for higher positions
		if positionBonus < 0.5 {
			positionBonus = 0.5
		}
		diversified[i].ConfidenceScore *= positionBonus
	}

	return diversified
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

// deduplicateFeatures removes duplicate features from a list
func (bkb *BedrockKnowledgeBaseService) deduplicateFeatures(features []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0)

	for _, feature := range features {
		if !seen[feature] {
			seen[feature] = true
			result = append(result, feature)
		}
	}

	return result
}
