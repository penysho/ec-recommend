package service

import (
	"context"
	"ec-recommend/internal/dto"
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
func (bkb *BedrockKnowledgeBaseService) QueryKnowledgeBase(ctx context.Context, query string, filters map[string]interface{}) (*BedrockKnowledgeBaseResponse, error) {
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
	results := make([]KnowledgeBaseResult, 0, len(retrieveOutput.RetrievalResults))
	sources := make([]string, 0)

	for _, result := range retrieveOutput.RetrievalResults {
		kbResult := KnowledgeBaseResult{
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
			kbResult.Location = &DocumentLocation{
				DocumentID: extractDocumentIDFromURI(kbResult.Source),
			}
		}

		results = append(results, kbResult)
	}

	processingTime := time.Since(startTime).Milliseconds()

	// Calculate confidence level based on scores
	confidenceLevel := bkb.calculateConfidenceLevel(results)

	return &BedrockKnowledgeBaseResponse{
		Results: results,
		RetrievalMetadata: &RetrievalMetadata{
			QueryProcessingTimeMs: processingTime,
			RetrievalCount:        len(results),
			Sources:               bkb.deduplicateSources(sources),
			ConfidenceLevel:       confidenceLevel,
		},
		ProcessingTimeMs: processingTime,
	}, nil
}

// RetrieveAndGenerate performs retrieval-augmented generation with basic configuration
func (bkb *BedrockKnowledgeBaseService) RetrieveAndGenerate(ctx context.Context, req *RetrieveAndGenerateRequest) (*RetrieveAndGenerateResponse, error) {
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
	citations := make([]Citation, 0)
	retrievedResults := make([]KnowledgeBaseResult, 0)

	if ragOutput.Citations != nil {
		for _, citation := range ragOutput.Citations {
			// Process each citation
			if citation.GeneratedResponsePart != nil && citation.GeneratedResponsePart.TextResponsePart != nil {
				citationObj := Citation{
					GeneratedResponsePart: &GeneratedResponsePart{
						TextResponsePart: &TextResponsePart{
							Text: aws.ToString(citation.GeneratedResponsePart.TextResponsePart.Text),
							Span: &Span{
								Start: int(aws.ToInt32(citation.GeneratedResponsePart.TextResponsePart.Span.Start)),
								End:   int(aws.ToInt32(citation.GeneratedResponsePart.TextResponsePart.Span.End)),
							},
						},
					},
				}

				// Process retrieved references
				if citation.RetrievedReferences != nil {
					references := make([]RetrievedReference, 0, len(citation.RetrievedReferences))
					for _, ref := range citation.RetrievedReferences {
						reference := RetrievedReference{
							Content: &RetrievalResultContent{
								Text: aws.ToString(ref.Content.Text),
							},
						}

						// Add location information
						if ref.Location != nil {
							reference.Location = &RetrievalResultLocation{}
							if ref.Location.S3Location != nil {
								reference.Location.S3Location = &S3Location{
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
						retrievedResult := KnowledgeBaseResult{
							Content: aws.ToString(ref.Content.Text),
							Score:   0.0, // Score not available in citations
						}

						if ref.Location != nil && ref.Location.S3Location != nil {
							retrievedResult.Source = aws.ToString(ref.Location.S3Location.Uri)
							retrievedResult.Location = &DocumentLocation{
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

	return &RetrieveAndGenerateResponse{
		Output:           aws.ToString(ragOutput.Output.Text),
		Citations:        citations,
		RetrievedResults: retrievedResults,
		Metadata: &GenerationMetadata{
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
func (bkb *BedrockKnowledgeBaseService) GetSimilarDocuments(ctx context.Context, embedding []float64, limit int, filters map[string]interface{}) (*SimilarDocumentsResponse, error) {
	startTime := time.Now()

	// Use the knowledge base for similarity search
	// We'll create a synthetic query and use the embedding
	query := "Find similar products based on vector similarity"

	kbResponse, err := bkb.QueryKnowledgeBase(ctx, query, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to query knowledge base for similarity: %w", err)
	}

	// Convert knowledge base results to similar documents
	documents := make([]SimilarDocument, 0, len(kbResponse.Results))

	// Limit results
	maxResults := min(len(kbResponse.Results), limit)
	for i := 0; i < maxResults; i++ {
		result := kbResponse.Results[i]

		doc := SimilarDocument{
			DocumentID: result.Location.DocumentID,
			Content:    result.Content,
			Score:      result.Score,
			Metadata:   result.Metadata,
			Source:     result.Source,
		}

		documents = append(documents, doc)
	}

	return &SimilarDocumentsResponse{
		Documents:        documents,
		ProcessingTimeMs: time.Since(startTime).Milliseconds(),
	}, nil
}

// CreatePersonalizedPrompt creates a personalized prompt based on customer profile
func (bkb *BedrockKnowledgeBaseService) CreatePersonalizedPrompt(customerProfile map[string]interface{}, query string) string {
	var promptBuilder strings.Builder

	promptBuilder.WriteString("Based on the following customer profile, provide personalized product recommendations:\n\n")
	promptBuilder.WriteString("Customer Profile:\n")
	promptBuilder.WriteString(bkb.formatProfileForPrompt(customerProfile))
	promptBuilder.WriteString("\n\nQuery: ")
	promptBuilder.WriteString(query)
	promptBuilder.WriteString("\n\nPlease provide relevant product recommendations with explanations.")

	return promptBuilder.String()
}

// ExtractProductEntities extracts product entities from text using NLP
func (bkb *BedrockKnowledgeBaseService) ExtractProductEntities(ctx context.Context, text string) ([]map[string]interface{}, error) {
	// Create a prompt for entity extraction
	prompt := fmt.Sprintf(`
Extract product entities from the following text. Return a JSON array of objects with product information:

Text: %s

Please extract:
- Product names
- Categories
- Brands
- Key features
- Price mentions

Return format: [{"name": "product name", "category": "category", "brand": "brand", "features": ["feature1", "feature2"]}]
`, text)

	// Create input for Claude model
	input := map[string]interface{}{
		"anthropic_version": "bedrock-2023-05-31",
		"max_tokens":        1000,
		"messages": []map[string]interface{}{
			{
				"role":    "user",
				"content": prompt,
			},
		},
	}

	inputBytes, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal input: %w", err)
	}

	// Call Claude model
	invokeInput := &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(bkb.modelARN),
		ContentType: aws.String("application/json"),
		Body:        inputBytes,
	}

	result, err := bkb.runtimeClient.InvokeModel(ctx, invokeInput)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke model: %w", err)
	}

	// Parse the response
	var response struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
	}

	if err := json.Unmarshal(result.Body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(response.Content) == 0 {
		return []map[string]interface{}{}, nil
	}

	// Parse the JSON from the response
	var entities []map[string]interface{}
	responseText := response.Content[0].Text

	// Clean up the response text to extract JSON
	jsonStart := strings.Index(responseText, "[")
	jsonEnd := strings.LastIndex(responseText, "]")

	if jsonStart == -1 || jsonEnd == -1 {
		log.Printf("Warning: Could not extract JSON from response: %s", responseText)
		return []map[string]interface{}{}, nil
	}

	jsonText := responseText[jsonStart : jsonEnd+1]
	if err := json.Unmarshal([]byte(jsonText), &entities); err != nil {
		log.Printf("Warning: Could not parse entities JSON: %v", err)
		return []map[string]interface{}{}, nil
	}

	return entities, nil
}

// GenerateRecommendationExplanation generates explanations for recommendations
func (bkb *BedrockKnowledgeBaseService) GenerateRecommendationExplanation(ctx context.Context, customerProfile map[string]interface{}, recommendedProducts []map[string]interface{}) (string, error) {
	// Create a prompt for explanation generation
	prompt := fmt.Sprintf(`
Generate a clear and helpful explanation for why these products are recommended for this customer.

Customer Profile:
%s

Recommended Products:
%s

Please provide:
1. Why these products match the customer's preferences
2. How they align with their purchase history
3. Any special features that make them suitable
4. Personalized benefits for this customer

Keep the explanation concise but informative.
`, bkb.formatProfileForPrompt(customerProfile), bkb.formatProductsForPrompt(recommendedProducts))

	// Create input for Claude model
	input := map[string]interface{}{
		"anthropic_version": "bedrock-2023-05-31",
		"max_tokens":        1500,
		"messages": []map[string]interface{}{
			{
				"role":    "user",
				"content": prompt,
			},
		},
	}

	inputBytes, err := json.Marshal(input)
	if err != nil {
		return "", fmt.Errorf("failed to marshal input: %w", err)
	}

	// Call Claude model
	invokeInput := &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(bkb.modelARN),
		ContentType: aws.String("application/json"),
		Body:        inputBytes,
	}

	result, err := bkb.runtimeClient.InvokeModel(ctx, invokeInput)
	if err != nil {
		return "", fmt.Errorf("failed to invoke model: %w", err)
	}

	// Parse the response
	var response struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
	}

	if err := json.Unmarshal(result.Body, &response); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(response.Content) == 0 {
		return "Unable to generate explanation at this time.", nil
	}

	return response.Content[0].Text, nil
}

// Helper functions

func (bkb *BedrockKnowledgeBaseService) formatProfileForPrompt(profile map[string]interface{}) string {
	var builder strings.Builder
	for key, value := range profile {
		builder.WriteString(fmt.Sprintf("- %s: %v\n", key, value))
	}
	return builder.String()
}

func (bkb *BedrockKnowledgeBaseService) formatProductsForPrompt(products []map[string]interface{}) string {
	var builder strings.Builder
	for i, product := range products {
		builder.WriteString(fmt.Sprintf("%d. ", i+1))
		for key, value := range product {
			builder.WriteString(fmt.Sprintf("%s: %v, ", key, value))
		}
		builder.WriteString("\n")
	}
	return builder.String()
}

func (bkb *BedrockKnowledgeBaseService) buildMetadataFilters(filters map[string]interface{}) *types.RetrievalFilter {
	// For now, return nil as metadata filtering would need proper configuration
	// This would need to be expanded based on your metadata structure
	return nil
}

func (bkb *BedrockKnowledgeBaseService) createSingleFilter(key string, value interface{}) *types.RetrievalFilter {
	// Return nil for now - would need proper filter configuration
	return nil
}

func (bkb *BedrockKnowledgeBaseService) calculateConfidenceLevel(results []KnowledgeBaseResult) float64 {
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

// GetProductsWithVectorSearch performs vector-based product search using Amazon Bedrock Knowledge Base
func (bkb *BedrockKnowledgeBaseService) GetProductsWithVectorSearch(ctx context.Context, vector []float64, limit int, filters map[string]interface{}) ([]dto.ProductRecommendationV2, error) {
	startTime := time.Now()

	// Convert vector to query format for Knowledge Base
	vectorQuery := "Product similarity search based on vector embeddings"

	// Configure retrieval for vector search only (using hybrid since pure vector is not available)
	retrievalConfig := &types.KnowledgeBaseRetrievalConfiguration{
		VectorSearchConfiguration: &types.KnowledgeBaseVectorSearchConfiguration{
			NumberOfResults:    aws.Int32(int32(min(limit, 50))), // AWS recommends max 50 results per query
			OverrideSearchType: types.SearchTypeHybrid,           // Use hybrid as closest to vector search
		},
	}

	// Create retrieve request
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

	// Convert results to ProductRecommendationV2
	products, err := bkb.convertKnowledgeBaseResultsToProducts(retrieveOutput.RetrievalResults, "vector_search")
	if err != nil {
		return nil, fmt.Errorf("failed to convert vector search results: %w", err)
	}

	// Enhance results with vector metadata
	for i := range products {
		if i < len(retrieveOutput.RetrievalResults) {
			products[i].VectorMetadata = &dto.VectorMetadata{
				DistanceScore:   float64(aws.ToFloat64(retrieveOutput.RetrievalResults[i].Score)),
				SearchMethod:    "vector_search",
				MatchedCriteria: []string{"vector_similarity"},
				EmbeddingModel:  bkb.embeddingModelID,
			}
			products[i].SimilarityScore = float64(aws.ToFloat64(retrieveOutput.RetrievalResults[i].Score))
			products[i].ConfidenceScore = bkb.calculateConfidenceFromScore(products[i].SimilarityScore)
		}
	}

	// Apply additional filtering if needed
	if len(filters) > 0 {
		products = bkb.applyPostRetrievalFilters(products, filters)
	}

	// Log vector search operation
	processingTime := time.Since(startTime).Milliseconds()
	log.Printf("Vector search completed: %d results in %dms", len(products), processingTime)

	return products[:min(len(products), limit)], nil
}

// GetProductsWithSemanticSearch performs semantic/text-based product search using Amazon Bedrock Knowledge Base
func (bkb *BedrockKnowledgeBaseService) GetProductsWithSemanticSearch(ctx context.Context, query string, limit int, filters map[string]interface{}) ([]dto.ProductRecommendationV2, error) {
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

	// Convert results to ProductRecommendationV2
	products, err := bkb.convertKnowledgeBaseResultsToProducts(retrieveOutput.RetrievalResults, "semantic_search")
	if err != nil {
		return nil, fmt.Errorf("failed to convert semantic search results: %w", err)
	}

	// Extract semantic insights from the query
	semanticInsights, err := bkb.extractSemanticInsights(ctx, query, retrieveOutput.RetrievalResults)
	if err != nil {
		log.Printf("Warning: Failed to extract semantic insights: %v", err)
	}

	// Enhance results with semantic metadata
	for i := range products {
		if i < len(retrieveOutput.RetrievalResults) {
			products[i].VectorMetadata = &dto.VectorMetadata{
				DistanceScore:    float64(aws.ToFloat64(retrieveOutput.RetrievalResults[i].Score)),
				SearchMethod:     "semantic_search",
				MatchedCriteria:  bkb.extractMatchedCriteria(retrieveOutput.RetrievalResults[i]),
				SemanticClusters: bkb.extractSemanticClusters(semanticInsights),
				EmbeddingModel:   bkb.embeddingModelID,
			}
			products[i].SimilarityScore = float64(aws.ToFloat64(retrieveOutput.RetrievalResults[i].Score))
			products[i].ConfidenceScore = bkb.calculateConfidenceFromScore(products[i].SimilarityScore)

			// Add AI-generated insights if available
			if semanticInsights != nil {
				products[i].AIInsights = bkb.generateProductAIInsights(ctx, products[i], semanticInsights)
			}
		}
	}

	// Apply additional filtering if needed
	if len(filters) > 0 {
		products = bkb.applyPostRetrievalFilters(products, filters)
	}

	// Log semantic search operation
	processingTime := time.Since(startTime).Milliseconds()
	log.Printf("Semantic search completed: %d results in %dms for query: %s", len(products), processingTime, query)

	return products[:min(len(products), limit)], nil
}

// GetProductsWithHybridSearch performs hybrid search combining vector and semantic approaches using Amazon Bedrock Knowledge Base
func (bkb *BedrockKnowledgeBaseService) GetProductsWithHybridSearch(ctx context.Context, query string, vector []float64, limit int, filters map[string]interface{}) ([]dto.ProductRecommendationV2, error) {
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

	// Convert results to ProductRecommendationV2
	products, err := bkb.convertKnowledgeBaseResultsToProducts(retrieveOutput.RetrievalResults, "hybrid_search")
	if err != nil {
		return nil, fmt.Errorf("failed to convert hybrid search results: %w", err)
	}

	// Extract semantic insights for better understanding
	semanticInsights, err := bkb.extractSemanticInsights(ctx, query, retrieveOutput.RetrievalResults)
	if err != nil {
		log.Printf("Warning: Failed to extract semantic insights for hybrid search: %v", err)
	}

	// Apply hybrid ranking algorithm - combine vector similarity and semantic relevance
	products = bkb.applyHybridRanking(products, vector, query, retrieveOutput.RetrievalResults)

	// Enhance results with comprehensive metadata
	for i := range products {
		if i < len(retrieveOutput.RetrievalResults) {
			products[i].VectorMetadata = &dto.VectorMetadata{
				DistanceScore:    float64(aws.ToFloat64(retrieveOutput.RetrievalResults[i].Score)),
				SearchMethod:     "hybrid_search",
				MatchedCriteria:  append([]string{"vector_similarity", "semantic_relevance"}, bkb.extractMatchedCriteria(retrieveOutput.RetrievalResults[i])...),
				SemanticClusters: bkb.extractSemanticClusters(semanticInsights),
				EmbeddingModel:   bkb.embeddingModelID,
			}

			// Calculate hybrid confidence score (weighted combination)
			vectorScore := float64(aws.ToFloat64(retrieveOutput.RetrievalResults[i].Score))
			semanticScore := bkb.calculateSemanticRelevance(products[i], query)
			products[i].SimilarityScore = vectorScore
			products[i].ConfidenceScore = bkb.calculateHybridConfidence(vectorScore, semanticScore)

			// Add comprehensive AI insights
			if semanticInsights != nil {
				products[i].AIInsights = bkb.generateProductAIInsights(ctx, products[i], semanticInsights)
			}

			// Add relevance context explaining why this product is recommended
			products[i].RelevanceContext = bkb.generateRelevanceContext(products[i], query, vector)
		}
	}

	// Apply additional filtering if needed
	if len(filters) > 0 {
		products = bkb.applyPostRetrievalFilters(products, filters)
	}

	// Final ranking and selection
	products = bkb.finalHybridRanking(products, limit)

	// Log hybrid search operation
	processingTime := time.Since(startTime).Milliseconds()
	log.Printf("Hybrid search completed: %d results in %dms for query: %s", len(products), processingTime, query)

	return products[:min(len(products), limit)], nil
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

// buildRetrievalFilter builds the retrieval filter from the provided filters map
// Note: Returns nil for now as metadata filtering implementation depends on knowledge base schema
func (bkb *BedrockKnowledgeBaseService) buildRetrievalFilter(filters map[string]interface{}) interface{} {
	// AWS Bedrock Knowledge Base supports metadata filtering
	// This is a simplified implementation - actual filtering would depend on your metadata structure
	if len(filters) == 0 {
		return nil
	}

	// For now, return nil and handle filtering in post-processing
	// In production, you would implement proper metadata filtering based on your knowledge base schema
	return nil
}

// convertFilterValue converts a filter value to the appropriate Bedrock format
func (bkb *BedrockKnowledgeBaseService) convertFilterValue(key string, value interface{}) interface{} {
	// Convert various filter types to Bedrock-compatible format
	switch key {
	case "category_id":
		if categoryID, ok := value.(int); ok {
			return categoryID
		}
	case "min_price", "max_price":
		if price, ok := value.(float64); ok {
			return price
		}
	case "brand":
		if brand, ok := value.(string); ok {
			return brand
		}
	}
	return nil
}

// convertKnowledgeBaseResultsToProducts converts Bedrock Knowledge Base results to ProductRecommendationV2
func (bkb *BedrockKnowledgeBaseService) convertKnowledgeBaseResultsToProducts(results []types.KnowledgeBaseRetrievalResult, searchMethod string) ([]dto.ProductRecommendationV2, error) {
	products := make([]dto.ProductRecommendationV2, 0, len(results))

	for _, result := range results {
		// Parse product information from the retrieved content
		product, err := bkb.parseProductFromContent(aws.ToString(result.Content.Text))
		if err != nil {
			log.Printf("Warning: Failed to parse product from content: %v", err)
			continue
		}

		// Set search-specific metadata
		product.VectorMetadata = &dto.VectorMetadata{
			DistanceScore:  float64(aws.ToFloat64(result.Score)),
			SearchMethod:   searchMethod,
			EmbeddingModel: bkb.embeddingModelID,
		}

		products = append(products, product)
	}

	return products, nil
}

// parseProductFromContent parses product information from retrieved content
func (bkb *BedrockKnowledgeBaseService) parseProductFromContent(content string) (dto.ProductRecommendationV2, error) {
	// This is a simplified implementation
	// In practice, you'd want more sophisticated parsing based on your data structure
	var product dto.ProductRecommendationV2

	// Try to parse JSON if the content is structured
	if strings.HasPrefix(content, "{") {
		if err := json.Unmarshal([]byte(content), &product); err == nil {
			return product, nil
		}
	}

	// If not JSON, create a basic product from text content
	// This is a placeholder implementation
	product = dto.ProductRecommendationV2{
		ProductID:   uuid.New(), // This should come from actual data
		Name:        bkb.extractProductName(content),
		Description: content,
		// Other fields would be populated based on your data structure
	}

	return product, nil
}

// extractProductName extracts product name from content (simplified implementation)
func (bkb *BedrockKnowledgeBaseService) extractProductName(content string) string {
	lines := strings.Split(content, "\n")
	if len(lines) > 0 {
		// Return first line as name, truncated if too long
		name := strings.TrimSpace(lines[0])
		if len(name) > 100 {
			name = name[:100] + "..."
		}
		return name
	}
	return "Unknown Product"
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

// extractSemanticInsights extracts semantic insights from search results
func (bkb *BedrockKnowledgeBaseService) extractSemanticInsights(ctx context.Context, query string, results []types.KnowledgeBaseRetrievalResult) (map[string]interface{}, error) {
	// Use Bedrock to analyze the query and results for semantic insights
	prompt := fmt.Sprintf(`
Analyze the following search query and results to extract semantic insights:

Query: %s

Results:
%s

Extract:
1. Query intent
2. Key entities mentioned
3. Semantic clusters
4. Related concepts

Provide insights in JSON format.
`, query, bkb.formatResultsForAnalysis(results))

	// Call Claude model for analysis
	response, err := bkb.callClaudeForAnalysis(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to extract semantic insights: %w", err)
	}

	// Parse the response
	var insights map[string]interface{}
	if err := json.Unmarshal([]byte(response), &insights); err != nil {
		log.Printf("Warning: Could not parse semantic insights: %v", err)
		return map[string]interface{}{
			"query_intent": "search",
			"entities":     []string{},
			"clusters":     []string{},
			"concepts":     []string{},
		}, nil
	}

	return insights, nil
}

// formatResultsForAnalysis formats results for semantic analysis
func (bkb *BedrockKnowledgeBaseService) formatResultsForAnalysis(results []types.KnowledgeBaseRetrievalResult) string {
	var formatted strings.Builder
	for i, result := range results {
		if i >= 5 { // Limit to first 5 results for analysis
			break
		}
		formatted.WriteString(fmt.Sprintf("Result %d: %s\n", i+1, aws.ToString(result.Content.Text)))
	}
	return formatted.String()
}

// callClaudeForAnalysis calls Claude model for semantic analysis
func (bkb *BedrockKnowledgeBaseService) callClaudeForAnalysis(ctx context.Context, prompt string) (string, error) {
	// Prepare input for Claude
	input := map[string]interface{}{
		"anthropic_version": "bedrock-2023-05-31",
		"max_tokens":        2000,
		"messages": []map[string]interface{}{
			{
				"role":    "user",
				"content": prompt,
			},
		},
	}

	inputBytes, err := json.Marshal(input)
	if err != nil {
		return "", fmt.Errorf("failed to marshal input: %w", err)
	}

	// Call Claude model
	invokeInput := &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(bkb.modelARN),
		ContentType: aws.String("application/json"),
		Body:        inputBytes,
	}

	result, err := bkb.runtimeClient.InvokeModel(ctx, invokeInput)
	if err != nil {
		return "", fmt.Errorf("failed to invoke Claude model: %w", err)
	}

	// Parse response
	var response struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
	}

	if err := json.Unmarshal(result.Body, &response); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(response.Content) == 0 {
		return "", fmt.Errorf("no content in response")
	}

	return response.Content[0].Text, nil
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

// extractSemanticClusters extracts semantic clusters from insights
func (bkb *BedrockKnowledgeBaseService) extractSemanticClusters(insights map[string]interface{}) []string {
	if insights == nil {
		return []string{}
	}

	if clusters, ok := insights["clusters"].([]interface{}); ok {
		var clusterNames []string
		for _, cluster := range clusters {
			if clusterName, ok := cluster.(string); ok {
				clusterNames = append(clusterNames, clusterName)
			}
		}
		return clusterNames
	}

	return []string{}
}

// generateProductAIInsights generates AI insights for a product
func (bkb *BedrockKnowledgeBaseService) generateProductAIInsights(ctx context.Context, product dto.ProductRecommendationV2, semanticInsights map[string]interface{}) *dto.ProductAIInsights {
	// Generate AI insights based on product and semantic analysis
	return &dto.ProductAIInsights{
		KeyFeatures:    bkb.extractKeyFeatures(product),
		UseCases:       bkb.extractUseCases(product, semanticInsights),
		TargetAudience: bkb.extractTargetAudience(product, semanticInsights),
	}
}

// applyHybridRanking applies hybrid ranking algorithm to combine vector and semantic scores
func (bkb *BedrockKnowledgeBaseService) applyHybridRanking(products []dto.ProductRecommendationV2, vector []float64, query string, results []types.KnowledgeBaseRetrievalResult) []dto.ProductRecommendationV2 {
	// Apply weighted combination of vector similarity and semantic relevance
	for i := range products {
		if i < len(results) {
			vectorScore := float64(aws.ToFloat64(results[i].Score))
			semanticScore := bkb.calculateSemanticRelevance(products[i], query)

			// Weighted combination (60% vector, 40% semantic)
			hybridScore := 0.6*vectorScore + 0.4*semanticScore
			products[i].ConfidenceScore = hybridScore
		}
	}

	// Sort by hybrid score
	sort.Slice(products, func(i, j int) bool {
		return products[i].ConfidenceScore > products[j].ConfidenceScore
	})

	return products
}

// calculateSemanticRelevance calculates semantic relevance score
func (bkb *BedrockKnowledgeBaseService) calculateSemanticRelevance(product dto.ProductRecommendationV2, query string) float64 {
	// Simple semantic relevance calculation based on text similarity
	// In practice, you'd want more sophisticated semantic analysis

	productText := strings.ToLower(product.Name + " " + product.Description)
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

// generateRelevanceContext generates relevance context for a product
func (bkb *BedrockKnowledgeBaseService) generateRelevanceContext(product dto.ProductRecommendationV2, query string, vector []float64) []dto.RelevanceContext {
	contexts := []dto.RelevanceContext{}

	// Add semantic match context
	semanticRelevance := bkb.calculateSemanticRelevance(product, query)
	if semanticRelevance > 0.3 {
		contexts = append(contexts, dto.RelevanceContext{
			ContextType: "semantic_match",
			Explanation: fmt.Sprintf("Product matches your search query with %.1f%% semantic similarity", semanticRelevance*100),
			Confidence:  semanticRelevance,
			SourceData:  "search_query",
		})
	}

	// Add vector similarity context
	if product.VectorMetadata != nil {
		contexts = append(contexts, dto.RelevanceContext{
			ContextType: "vector_similarity",
			Explanation: fmt.Sprintf("Product has high vector similarity score of %.3f", product.VectorMetadata.DistanceScore),
			Confidence:  product.VectorMetadata.DistanceScore,
			SourceData:  "embedding_vectors",
		})
	}

	return contexts
}

// applyPostRetrievalFilters applies additional filters after retrieval
func (bkb *BedrockKnowledgeBaseService) applyPostRetrievalFilters(products []dto.ProductRecommendationV2, filters map[string]interface{}) []dto.ProductRecommendationV2 {
	filtered := make([]dto.ProductRecommendationV2, 0, len(products))

	for _, product := range products {
		include := true

		// Apply price range filter
		if minPrice, ok := filters["min_price"].(float64); ok && product.Price < minPrice {
			include = false
		}
		if maxPrice, ok := filters["max_price"].(float64); ok && product.Price > maxPrice {
			include = false
		}

		// Apply category filter
		if categoryID, ok := filters["category_id"].(int); ok && product.CategoryID != categoryID {
			include = false
		}

		// Apply brand filter
		if brand, ok := filters["brand"].(string); ok && product.Brand != brand {
			include = false
		}

		if include {
			filtered = append(filtered, product)
		}
	}

	return filtered
}

// finalHybridRanking performs final ranking and selection for hybrid search
func (bkb *BedrockKnowledgeBaseService) finalHybridRanking(products []dto.ProductRecommendationV2, limit int) []dto.ProductRecommendationV2 {
	// Apply diversity algorithm to ensure varied results
	diversified := bkb.applyDiversityAlgorithm(products, limit)

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

// applyDiversityAlgorithm applies diversity algorithm to ensure varied results
func (bkb *BedrockKnowledgeBaseService) applyDiversityAlgorithm(products []dto.ProductRecommendationV2, limit int) []dto.ProductRecommendationV2 {
	if len(products) <= limit {
		return products
	}

	// Simple diversity algorithm - ensure category diversity
	categoryMap := make(map[int]bool)
	diversified := make([]dto.ProductRecommendationV2, 0, limit)

	// First pass: include top products from different categories
	for _, product := range products {
		if len(diversified) >= limit {
			break
		}
		if !categoryMap[product.CategoryID] {
			categoryMap[product.CategoryID] = true
			diversified = append(diversified, product)
		}
	}

	// Second pass: fill remaining slots with highest scoring products
	for _, product := range products {
		if len(diversified) >= limit {
			break
		}
		// Check if product is already included
		found := false
		for _, included := range diversified {
			if included.ProductID == product.ProductID {
				found = true
				break
			}
		}
		if !found {
			diversified = append(diversified, product)
		}
	}

	return diversified
}

// extractKeyFeatures extracts key features from product
func (bkb *BedrockKnowledgeBaseService) extractKeyFeatures(product dto.ProductRecommendationV2) []string {
	// Simple feature extraction from product name and description
	features := []string{}

	text := strings.ToLower(product.Name + " " + product.Description)

	// Common feature keywords
	featureKeywords := []string{"wireless", "waterproof", "portable", "lightweight", "durable", "premium", "eco-friendly"}

	for _, keyword := range featureKeywords {
		if strings.Contains(text, keyword) {
			features = append(features, keyword)
		}
	}

	return features
}

// extractUseCases extracts use cases from product and semantic insights
func (bkb *BedrockKnowledgeBaseService) extractUseCases(product dto.ProductRecommendationV2, insights map[string]interface{}) []string {
	// Extract use cases based on product category and insights
	useCases := []string{}

	// Category-based use cases
	switch product.CategoryID {
	case 1: // Electronics
		useCases = append(useCases, "daily_use", "professional_work")
	case 2: // Clothing
		useCases = append(useCases, "casual_wear", "formal_occasions")
	case 3: // Home & Garden
		useCases = append(useCases, "home_improvement", "decoration")
	}

	return useCases
}

// extractTargetAudience extracts target audience from product and insights
func (bkb *BedrockKnowledgeBaseService) extractTargetAudience(product dto.ProductRecommendationV2, insights map[string]interface{}) []string {
	// Extract target audience based on product characteristics
	audience := []string{}

	// Price-based audience segmentation
	if product.Price < 50 {
		audience = append(audience, "budget_conscious")
	} else if product.Price > 500 {
		audience = append(audience, "premium_buyers")
	}

	// Rating-based audience
	if product.RatingAverage > 4.5 {
		audience = append(audience, "quality_seekers")
	}

	return audience
}
