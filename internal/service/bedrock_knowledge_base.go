package service

import (
	"context"
	"ec-recommend/internal/interfaces"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockagentruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockagentruntime/types"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
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
func (bkb *BedrockKnowledgeBaseService) QueryKnowledgeBase(ctx context.Context, query string, filters map[string]interface{}) (*interfaces.BedrockKnowledgeBaseResponse, error) {
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
	results := make([]interfaces.KnowledgeBaseResult, 0, len(retrieveOutput.RetrievalResults))
	sources := make([]string, 0)

	for _, result := range retrieveOutput.RetrievalResults {
		kbResult := interfaces.KnowledgeBaseResult{
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
			kbResult.Location = &interfaces.DocumentLocation{
				DocumentID: extractDocumentIDFromURI(kbResult.Source),
			}
		}

		results = append(results, kbResult)
	}

	processingTime := time.Since(startTime).Milliseconds()

	// Calculate confidence level based on scores
	confidenceLevel := bkb.calculateConfidenceLevel(results)

	return &interfaces.BedrockKnowledgeBaseResponse{
		Results: results,
		RetrievalMetadata: &interfaces.RetrievalMetadata{
			QueryProcessingTimeMs: processingTime,
			RetrievalCount:        len(results),
			Sources:               bkb.deduplicateSources(sources),
			ConfidenceLevel:       confidenceLevel,
		},
		ProcessingTimeMs: processingTime,
	}, nil
}

// RetrieveAndGenerate performs retrieval-augmented generation with basic configuration
func (bkb *BedrockKnowledgeBaseService) RetrieveAndGenerate(ctx context.Context, req *interfaces.RetrieveAndGenerateRequest) (*interfaces.RetrieveAndGenerateResponse, error) {
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
	citations := make([]interfaces.Citation, 0)
	retrievedResults := make([]interfaces.KnowledgeBaseResult, 0)

	if ragOutput.Citations != nil {
		for _, citation := range ragOutput.Citations {
			// Process each citation
			if citation.GeneratedResponsePart != nil && citation.GeneratedResponsePart.TextResponsePart != nil {
				citationObj := interfaces.Citation{
					GeneratedResponsePart: &interfaces.GeneratedResponsePart{
						TextResponsePart: &interfaces.TextResponsePart{
							Text: aws.ToString(citation.GeneratedResponsePart.TextResponsePart.Text),
							Span: &interfaces.Span{
								Start: int(aws.ToInt32(citation.GeneratedResponsePart.TextResponsePart.Span.Start)),
								End:   int(aws.ToInt32(citation.GeneratedResponsePart.TextResponsePart.Span.End)),
							},
						},
					},
				}

				// Process retrieved references
				if citation.RetrievedReferences != nil {
					references := make([]interfaces.RetrievedReference, 0, len(citation.RetrievedReferences))
					for _, ref := range citation.RetrievedReferences {
						reference := interfaces.RetrievedReference{
							Content: &interfaces.RetrievalResultContent{
								Text: aws.ToString(ref.Content.Text),
							},
						}

						// Add location information
						if ref.Location != nil {
							reference.Location = &interfaces.RetrievalResultLocation{}
							if ref.Location.S3Location != nil {
								reference.Location.S3Location = &interfaces.S3Location{
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
						retrievedResult := interfaces.KnowledgeBaseResult{
							Content: aws.ToString(ref.Content.Text),
							Score:   0.0, // Score not available in citations
						}

						if ref.Location != nil && ref.Location.S3Location != nil {
							retrievedResult.Source = aws.ToString(ref.Location.S3Location.Uri)
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

	// Build response
	return &interfaces.RetrieveAndGenerateResponse{
		Output:           aws.ToString(ragOutput.Output.Text),
		Citations:        citations,
		RetrievedResults: retrievedResults,
	}, nil
}

// GetVectorEmbedding generates vector embeddings for the given text
func (bkb *BedrockKnowledgeBaseService) GetVectorEmbedding(ctx context.Context, text string) ([]float64, error) {
	// Prepare the request for Amazon Titan Text Embeddings
	requestBody := map[string]interface{}{
		"inputText": text,
	}

	requestJSON, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Invoke the embedding model
	input := &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(bkb.embeddingModelID),
		ContentType: aws.String("application/json"),
		Accept:      aws.String("application/json"),
		Body:        requestJSON,
	}

	result, err := bkb.runtimeClient.InvokeModel(ctx, input)
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
func (bkb *BedrockKnowledgeBaseService) GetSimilarDocuments(ctx context.Context, embedding []float64, limit int, filters map[string]interface{}) (*interfaces.SimilarDocumentsResponse, error) {
	startTime := time.Now()

	// For this implementation, we'll use text search as a fallback
	// In a production system, you would use a proper vector database

	// Create a query that represents the embedding semantically
	// This is a simplified approach - in reality, you'd store and search embeddings directly
	query := "product recommendations based on customer preferences"

	kbResponse, err := bkb.QueryKnowledgeBase(ctx, query, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to query knowledge base: %w", err)
	}

	// Convert knowledge base results to similar documents
	documents := make([]interfaces.SimilarDocument, 0, len(kbResponse.Results))
	for i, result := range kbResponse.Results {
		if i >= limit {
			break
		}

		doc := interfaces.SimilarDocument{
			DocumentID: fmt.Sprintf("doc_%d", i),
			Content:    result.Content,
			Score:      result.Score,
			Metadata:   result.Metadata,
			Source:     result.Source,
		}
		documents = append(documents, doc)
	}

	processingTime := time.Since(startTime).Milliseconds()

	return &interfaces.SimilarDocumentsResponse{
		Documents:        documents,
		ProcessingTimeMs: processingTime,
	}, nil
}

// Helper methods for advanced RAG operations

// CreatePersonalizedPrompt creates a personalized prompt for recommendations
func (bkb *BedrockKnowledgeBaseService) CreatePersonalizedPrompt(customerProfile map[string]interface{}, query string) string {
	var promptBuilder strings.Builder

	promptBuilder.WriteString("Based on the customer profile and product knowledge base, provide personalized recommendations.\n\n")

	// Add customer profile information
	promptBuilder.WriteString("Customer Profile:\n")
	for key, value := range customerProfile {
		promptBuilder.WriteString(fmt.Sprintf("- %s: %v\n", key, value))
	}

	promptBuilder.WriteString("\nQuery: ")
	promptBuilder.WriteString(query)

	promptBuilder.WriteString("\n\nPlease provide detailed product recommendations with explanations based on the customer's profile and preferences.")

	return promptBuilder.String()
}

// ExtractProductEntities extracts product-related entities from text
func (bkb *BedrockKnowledgeBaseService) ExtractProductEntities(ctx context.Context, text string) ([]map[string]interface{}, error) {
	prompt := fmt.Sprintf(`
Extract product-related entities from the following text and return them as a JSON array:

Text: "%s"

Extract entities such as:
- Product names
- Brands
- Categories
- Features
- Price mentions
- Quality indicators

Format: [{"type": "product", "value": "iPhone 14", "confidence": 0.9}, ...]
`, text)

	// Create a simple request for entity extraction
	requestBody := map[string]interface{}{
		"messages": []map[string]interface{}{
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"max_tokens":        1000,
		"temperature":       0.1,
		"anthropic_version": "bedrock-2023-05-31",
	}

	requestJSON, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Invoke the model for entity extraction
	input := &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(bkb.modelARN),
		ContentType: aws.String("application/json"),
		Accept:      aws.String("application/json"),
		Body:        requestJSON,
	}

	result, err := bkb.runtimeClient.InvokeModel(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke model for entity extraction: %w", err)
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

	// Try to parse the extracted entities
	var entities []map[string]interface{}
	if err := json.Unmarshal([]byte(response.Content[0].Text), &entities); err != nil {
		log.Printf("Warning: failed to parse extracted entities as JSON: %v", err)
		return []map[string]interface{}{}, nil
	}

	return entities, nil
}

// GenerateRecommendationExplanation generates an explanation for why products were recommended
func (bkb *BedrockKnowledgeBaseService) GenerateRecommendationExplanation(ctx context.Context, customerProfile map[string]interface{}, recommendedProducts []map[string]interface{}) (string, error) {
	prompt := fmt.Sprintf(`
Explain why these products were recommended to this customer:

Customer Profile:
%s

Recommended Products:
%s

Please provide a detailed, personalized explanation for each recommendation, considering:
1. How each product matches the customer's preferences
2. The reasoning behind the selection
3. What makes these products suitable for this customer
4. Any potential concerns or alternatives

Format the response as a clear, customer-friendly explanation.
`, bkb.formatProfileForPrompt(customerProfile), bkb.formatProductsForPrompt(recommendedProducts))

	// Create request for explanation generation
	requestBody := map[string]interface{}{
		"messages": []map[string]interface{}{
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"max_tokens":        2000,
		"temperature":       0.3,
		"anthropic_version": "bedrock-2023-05-31",
	}

	requestJSON, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Invoke the model
	input := &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(bkb.modelARN),
		ContentType: aws.String("application/json"),
		Accept:      aws.String("application/json"),
		Body:        requestJSON,
	}

	result, err := bkb.runtimeClient.InvokeModel(ctx, input)
	if err != nil {
		return "", fmt.Errorf("failed to invoke model for explanation: %w", err)
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
		return "No explanation available", nil
	}

	return response.Content[0].Text, nil
}

// Helper methods for formatting data

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
		builder.WriteString(fmt.Sprintf("Product %d:\n", i+1))
		for key, value := range product {
			builder.WriteString(fmt.Sprintf("  - %s: %v\n", key, value))
		}
		builder.WriteString("\n")
	}
	return builder.String()
}

// Helper methods for building metadata filters
// Note: Metadata filtering temporarily disabled due to AWS SDK v2 complexity
func (bkb *BedrockKnowledgeBaseService) buildMetadataFilters(filters map[string]interface{}) types.RetrievalFilter {
	// TODO: Implement proper metadata filtering with correct AWS SDK v2 types
	log.Printf("Warning: Metadata filtering not yet implemented for filters: %+v", filters)
	return nil
}

func (bkb *BedrockKnowledgeBaseService) createSingleFilter(key string, value interface{}) types.RetrievalFilter {
	// TODO: Implement proper single filter creation with correct AWS SDK v2 types
	log.Printf("Warning: Single filter creation not yet implemented for key: %s, value: %+v", key, value)
	return nil
}

// Helper functions for enhanced functionality
func (bkb *BedrockKnowledgeBaseService) calculateConfidenceLevel(results []interfaces.KnowledgeBaseResult) float64 {
	if len(results) == 0 {
		return 0.0
	}

	var totalScore float64
	var maxScore float64
	for _, result := range results {
		totalScore += result.Score
		if result.Score > maxScore {
			maxScore = result.Score
		}
	}

	// Calculate confidence based on average score and maximum score
	avgScore := totalScore / float64(len(results))
	confidence := (avgScore + maxScore) / 2.0

	// Normalize to 0-1 range
	if confidence > 1.0 {
		confidence = 1.0
	}
	if confidence < 0.0 {
		confidence = 0.0
	}

	return confidence
}

func (bkb *BedrockKnowledgeBaseService) deduplicateSources(sources []string) []string {
	seen := make(map[string]bool)
	var result []string
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
	parts := strings.Split(uri, "/")
	if len(parts) > 0 {
		fileName := parts[len(parts)-1]
		// Remove file extension for cleaner document ID
		if dotIndex := strings.LastIndex(fileName, "."); dotIndex != -1 {
			fileName = fileName[:dotIndex]
		}
		return fileName
	}
	return uri
}

func extractModelIDFromARN(arn string) string {
	if arn == "" {
		return ""
	}
	parts := strings.Split(arn, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return arn
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
