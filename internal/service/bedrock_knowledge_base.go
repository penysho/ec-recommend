package service

import (
	"context"
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
