package service

import (
	"context"
	"ec-recommend/internal/dto"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
)

// RecommendationService implements the RecommendationServiceInterface
type RecommendationService struct {
	repo        RecommendationRepositoryInterface
	chatService ChatServiceInterface
	modelID     string
}

// NewRecommendationService creates a new recommendation service instance
func NewRecommendationService(repo RecommendationRepositoryInterface, chatService ChatServiceInterface, modelID string) *RecommendationService {
	return &RecommendationService{
		repo:        repo,
		chatService: chatService,
		modelID:     modelID,
	}
}

// GetRecommendations generates product recommendations based on the request type
func (rs *RecommendationService) GetRecommendations(ctx context.Context, req *dto.RecommendationRequest) (*dto.RecommendationResponse, error) {
	startTime := time.Now()

	// Set default values
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.RecommendationType == "" {
		req.RecommendationType = "hybrid"
	}
	if req.ContextType == "" {
		req.ContextType = "homepage"
	}

	// Get customer profile
	profile, err := rs.GetCustomerProfile(ctx, req.CustomerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer profile: %w", err)
	}

	var recommendations []dto.ProductRecommendation
	var algorithmVersion string

	// Generate recommendations based on type
	switch req.RecommendationType {
	case "similar":
		if req.ProductID == nil {
			return nil, fmt.Errorf("product_id is required for similar recommendations")
		}
		recommendations, err = rs.GetSimilarProducts(ctx, *req.ProductID, req.Limit)
		algorithmVersion = "similar_v1.0"
	case "collaborative":
		recommendations, err = rs.getCollaborativeRecommendations(ctx, profile, req.Limit)
		algorithmVersion = "collaborative_v1.0"
	case "content_based":
		recommendations, err = rs.getContentBasedRecommendations(ctx, profile, req.Limit)
		algorithmVersion = "content_based_v1.0"
	case "hybrid":
		recommendations, err = rs.getHybridRecommendations(ctx, profile, req)
		algorithmVersion = "hybrid_v1.0"
	default:
		return nil, fmt.Errorf("unsupported recommendation type: %s", req.RecommendationType)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to generate recommendations: %w", err)
	}

	// Filter out owned products if requested
	if req.ExcludeOwned {
		recommendations = rs.filterOwnedProducts(recommendations, profile.PurchaseHistory)
	}

	// Limit results
	if len(recommendations) > req.Limit {
		recommendations = recommendations[:req.Limit]
	}

	fmt.Println("recommendations", len(recommendations))
	for _, rec := range recommendations {
		fmt.Println(rec.Name)
	}

	// Generate AI-powered explanations and confidence scores
	recommendations, err = rs.enhanceWithAI(ctx, recommendations, profile, req.ContextType)
	if err != nil {
		// Log error but don't fail the request
		fmt.Printf("Warning: failed to enhance recommendations with AI: %v\n", err)
	} else {
		// Sort recommendations by AI-enhanced confidence scores in descending order
		sort.Slice(recommendations, func(i, j int) bool {
			return recommendations[i].ConfidenceScore > recommendations[j].ConfidenceScore
		})
	}

	fmt.Println("recommendations after AI enhancement and sorting", len(recommendations))
	for _, rec := range recommendations {
		fmt.Printf("%s (confidence: %.3f)\n", rec.Name, rec.ConfidenceScore)
	}

	// Log recommendation for analytics
	sessionID := uuid.New()
	productIDs := make([]uuid.UUID, len(recommendations))
	for i, rec := range recommendations {
		productIDs[i] = rec.ProductID
	}

	err = rs.repo.LogRecommendation(ctx, req.CustomerID, req.RecommendationType, req.ContextType, productIDs, sessionID)
	if err != nil {
		fmt.Printf("Warning: failed to log recommendation: %v\n", err)
	}

	processingTime := time.Since(startTime).Milliseconds()

	return &dto.RecommendationResponse{
		CustomerID:         req.CustomerID,
		Recommendations:    recommendations,
		RecommendationType: req.RecommendationType,
		ContextType:        req.ContextType,
		GeneratedAt:        time.Now(),
		Metadata: dto.RecommendationMetadata{
			AlgorithmVersion: algorithmVersion,
			ProcessingTimeMs: processingTime,
			TotalProducts:    len(recommendations),
			FilteredProducts: len(recommendations),
			AIModelUsed:      rs.modelID,
			SessionID:        sessionID,
		},
	}, nil
}

// GetCustomerProfile retrieves comprehensive customer profile data
func (rs *RecommendationService) GetCustomerProfile(ctx context.Context, customerID uuid.UUID) (*dto.CustomerProfile, error) {
	profile, err := rs.repo.GetCustomerByID(ctx, customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}

	// Get purchase history
	purchaseHistory, err := rs.repo.GetCustomerPurchaseHistory(ctx, customerID, 50)
	if err != nil {
		return nil, fmt.Errorf("failed to get purchase history: %w", err)
	}
	profile.PurchaseHistory = purchaseHistory

	// Get recent activities
	activities, err := rs.repo.GetCustomerActivities(ctx, customerID, 100)
	if err != nil {
		return nil, fmt.Errorf("failed to get activities: %w", err)
	}
	profile.RecentActivities = activities

	return profile, nil
}

// LogRecommendationInteraction logs customer interactions with recommendations
func (rs *RecommendationService) LogRecommendationInteraction(ctx context.Context, analytics *dto.RecommendationAnalytics) error {
	return rs.repo.LogRecommendationInteraction(ctx, analytics)
}

// GetSimilarProducts finds products similar to a given product using content-based filtering
func (rs *RecommendationService) GetSimilarProducts(ctx context.Context, productID uuid.UUID, limit int) ([]dto.ProductRecommendation, error) {
	// Get the target product to extract its tags
	products, err := rs.repo.GetProductsByIDs(ctx, []uuid.UUID{productID})
	if err != nil || len(products) == 0 {
		return nil, fmt.Errorf("failed to get target product: %w", err)
	}

	targetProduct := products[0]

	// Find similar products by tags
	similarProducts, err := rs.repo.GetSimilarProductsByTags(ctx, targetProduct.Tags, productID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get similar products: %w", err)
	}

	// Calculate similarity scores based on shared tags
	for i := range similarProducts {
		similarProducts[i].ConfidenceScore = rs.calculateTagSimilarity(targetProduct.Tags, similarProducts[i].Tags)
	}

	// Sort by confidence score
	sort.Slice(similarProducts, func(i, j int) bool {
		return similarProducts[i].ConfidenceScore > similarProducts[j].ConfidenceScore
	})

	return similarProducts, nil
}

// GetTrendingProducts returns currently trending products
func (rs *RecommendationService) GetTrendingProducts(ctx context.Context, categoryID *int, limit int) ([]dto.ProductRecommendation, error) {
	return rs.repo.GetTrendingProducts(ctx, categoryID, limit)
}

// GetPersonalizedRecommendations generates AI-powered personalized recommendations
func (rs *RecommendationService) GetPersonalizedRecommendations(ctx context.Context, profile *dto.CustomerProfile, limit int) ([]dto.ProductRecommendation, error) {
	// Create AI prompt based on customer profile
	prompt := rs.createPersonalizationPrompt(profile)

	// Get chat response
	chatResponse, err := rs.chatService.GenerateResponse(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat response: %w", err)
	}

	// Parse chat response to extract product recommendations
	productIDs, err := rs.parseAIRecommendations(chatResponse.Content)
	if err != nil {
		// Fallback to content-based recommendations
		return rs.getContentBasedRecommendations(ctx, profile, limit)
	}

	// Get product details
	recommendations, err := rs.repo.GetProductsByIDs(ctx, productIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get recommended products: %w", err)
	}

	return recommendations, nil
}

// getCollaborativeRecommendations implements collaborative filtering
func (rs *RecommendationService) getCollaborativeRecommendations(ctx context.Context, profile *dto.CustomerProfile, limit int) ([]dto.ProductRecommendation, error) {
	// Find customers with similar purchase patterns
	similarCustomers, err := rs.repo.GetCustomersWithSimilarPurchases(ctx, profile.CustomerID, 20)
	if err != nil {
		return nil, fmt.Errorf("failed to get similar customers: %w", err)
	}

	// Get products owned by the current customer
	ownedProductIDs := make([]uuid.UUID, len(profile.PurchaseHistory))
	for i, purchase := range profile.PurchaseHistory {
		ownedProductIDs[i] = purchase.ProductID
	}

	// Get popular products among similar customers
	recommendations, err := rs.repo.GetPopularProductsAmongSimilarCustomers(ctx, similarCustomers, ownedProductIDs, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get collaborative recommendations: %w", err)
	}

	return recommendations, nil
}

// getContentBasedRecommendations implements content-based filtering
func (rs *RecommendationService) getContentBasedRecommendations(ctx context.Context, profile *dto.CustomerProfile, limit int) ([]dto.ProductRecommendation, error) {
	var allRecommendations []dto.ProductRecommendation

	// Get recommendations based on preferred categories
	for _, categoryID := range profile.PreferredCategories {
		categoryProducts, err := rs.repo.GetProductsByCategory(ctx, categoryID, limit/len(profile.PreferredCategories)+1)
		if err != nil {
			continue
		}
		allRecommendations = append(allRecommendations, categoryProducts...)
	}

	// Get recommendations based on price range
	if profile.PriceRangeMin != nil && profile.PriceRangeMax != nil {
		priceRangeProducts, err := rs.repo.GetProductsInPriceRange(ctx, *profile.PriceRangeMin, *profile.PriceRangeMax, limit/2)
		if err == nil {
			allRecommendations = append(allRecommendations, priceRangeProducts...)
		}
	}

	// Remove duplicates and limit results
	uniqueProducts := rs.removeDuplicateProducts(allRecommendations)
	if len(uniqueProducts) > limit {
		uniqueProducts = uniqueProducts[:limit]
	}

	return uniqueProducts, nil
}

// getHybridRecommendations combines multiple recommendation strategies
func (rs *RecommendationService) getHybridRecommendations(ctx context.Context, profile *dto.CustomerProfile, req *dto.RecommendationRequest) ([]dto.ProductRecommendation, error) {
	var allRecommendations []dto.ProductRecommendation

	// Get collaborative recommendations (40% weight)
	collaborativeRecs, err := rs.getCollaborativeRecommendations(ctx, profile, req.Limit/2)
	if err == nil {
		for i := range collaborativeRecs {
			collaborativeRecs[i].ConfidenceScore = collaborativeRecs[i].ConfidenceScore * 0.4
		}
		allRecommendations = append(allRecommendations, collaborativeRecs...)
	}

	// Get content-based recommendations (40% weight)
	contentRecs, err := rs.getContentBasedRecommendations(ctx, profile, req.Limit/2)
	if err == nil {
		for i := range contentRecs {
			contentRecs[i].ConfidenceScore = contentRecs[i].ConfidenceScore * 0.4
		}
		allRecommendations = append(allRecommendations, contentRecs...)
	}

	// Get trending products (20% weight)
	var categoryID *int
	if req.CategoryID != nil {
		categoryID = req.CategoryID
	} else if len(profile.PreferredCategories) > 0 {
		categoryID = &profile.PreferredCategories[0]
	}

	trendingRecs, err := rs.GetTrendingProducts(ctx, categoryID, req.Limit/4)
	if err == nil {
		for i := range trendingRecs {
			trendingRecs[i].ConfidenceScore = trendingRecs[i].ConfidenceScore * 0.2
		}
		allRecommendations = append(allRecommendations, trendingRecs...)
	}

	// Remove duplicates and sort by confidence score
	uniqueProducts := rs.removeDuplicateProducts(allRecommendations)
	sort.Slice(uniqueProducts, func(i, j int) bool {
		return uniqueProducts[i].ConfidenceScore > uniqueProducts[j].ConfidenceScore
	})

	if len(uniqueProducts) > req.Limit {
		uniqueProducts = uniqueProducts[:req.Limit]
	}

	return uniqueProducts, nil
}

// enhanceWithAI adds AI-generated explanations and confidence scores
func (rs *RecommendationService) enhanceWithAI(ctx context.Context, recommendations []dto.ProductRecommendation, profile *dto.CustomerProfile, contextType string) ([]dto.ProductRecommendation, error) {
	if len(recommendations) == 0 {
		return recommendations, nil
	}

	// Create prompt for AI enhancement
	prompt := rs.createEnhancementPrompt(recommendations, profile, contextType)

	// Get chat response
	chatResponse, err := rs.chatService.GenerateResponse(ctx, prompt)
	if err != nil {
		return recommendations, err
	}

	// Parse chat response and enhance recommendations
	enhanced, err := rs.parseAIEnhancements(chatResponse.Content, recommendations)
	if err != nil {
		return recommendations, err
	}

	return enhanced, nil
}

// Helper methods

func (rs *RecommendationService) calculateTagSimilarity(tags1, tags2 []string) float64 {
	if len(tags1) == 0 || len(tags2) == 0 {
		return 0.0
	}

	tagSet1 := make(map[string]bool)
	for _, tag := range tags1 {
		tagSet1[tag] = true
	}

	commonTags := 0
	for _, tag := range tags2 {
		if tagSet1[tag] {
			commonTags++
		}
	}

	// Jaccard similarity
	totalUniqueTags := len(tags1) + len(tags2) - commonTags
	if totalUniqueTags == 0 {
		return 0.0
	}

	return float64(commonTags) / float64(totalUniqueTags)
}

func (rs *RecommendationService) filterOwnedProducts(recommendations []dto.ProductRecommendation, purchaseHistory []dto.PurchaseItem) []dto.ProductRecommendation {
	ownedProducts := make(map[uuid.UUID]bool)
	for _, purchase := range purchaseHistory {
		ownedProducts[purchase.ProductID] = true
	}

	var filtered []dto.ProductRecommendation
	for _, rec := range recommendations {
		if !ownedProducts[rec.ProductID] {
			filtered = append(filtered, rec)
		}
	}

	return filtered
}

func (rs *RecommendationService) removeDuplicateProducts(recommendations []dto.ProductRecommendation) []dto.ProductRecommendation {
	seen := make(map[uuid.UUID]bool)
	var unique []dto.ProductRecommendation

	for _, rec := range recommendations {
		if !seen[rec.ProductID] {
			seen[rec.ProductID] = true
			unique = append(unique, rec)
		}
	}

	return unique
}

func (rs *RecommendationService) createPersonalizationPrompt(profile *dto.CustomerProfile) string {
	return fmt.Sprintf(`
Based on the following customer profile, recommend products that would be most relevant:

Customer Profile:
- Total Spent: %.2f
- Order Count: %d
- Preferred Categories: %v
- Preferred Brands: %v
- Lifestyle Tags: %v
- Is Premium: %t

Recent Purchase History (last 10 items):
%s

Please provide a JSON array of product IDs that would be most suitable for this customer.
Format: ["product-id-1", "product-id-2", ...]
`,
		profile.TotalSpent,
		profile.OrderCount,
		profile.PreferredCategories,
		profile.PreferredBrands,
		profile.LifestyleTags,
		profile.IsPremium,
		rs.formatPurchaseHistory(profile.PurchaseHistory),
	)
}

func (rs *RecommendationService) createEnhancementPrompt(recommendations []dto.ProductRecommendation, profile *dto.CustomerProfile, contextType string) string {
	return fmt.Sprintf(`
You are an advanced product recommendation AI. Analyze each product recommendation and provide personalized confidence scores and explanations based on comprehensive customer data.

## Customer Profile Analysis:
**Financial Profile:**
- Total Spent: %.2f
- Order Count: %d
- Average Order Value: %.2f
- Is Premium Customer: %t
- Price Range Preference: %s

**Behavioral Profile:**
- Preferred Categories: %v
- Preferred Brands: %v
- Lifestyle Tags: %v

**Purchase Patterns:**
%s

**Recent Activity Patterns:**
%s

## Recommendation Context: %s

## Products to Analyze:
%s

## Instructions:
1. For each product, calculate a confidence score (0.0-1.0) based on:
   - Purchase history alignment (30%%)
   - Category/brand preferences (25%%)
   - Price range compatibility (20%%)
   - Lifestyle and behavioral fit (15%%)
   - Context appropriateness (10%%)

2. Provide detailed reasoning considering:
   - Why this product matches their buying patterns
   - How it fits their lifestyle and preferences
   - Timing and context relevance
   - Potential value and satisfaction

3. Rank products by recommendation priority and provide nuanced confidence scores.

Please provide a JSON response with the following format:
{
  "enhancements": [
    {
      "product_id": "uuid",
      "confidence_score": 0.85,
      "reason": "Detailed personalized explanation including purchase pattern analysis, lifestyle fit, and value proposition"
    }
  ]
}
`,
		profile.TotalSpent,
		profile.OrderCount,
		rs.calculateAverageOrderValue(profile),
		profile.IsPremium,
		rs.formatPriceRangePreference(profile),
		profile.PreferredCategories,
		profile.PreferredBrands,
		profile.LifestyleTags,
		rs.formatDetailedPurchaseHistory(profile.PurchaseHistory),
		rs.formatRecentActivities(profile.RecentActivities),
		contextType,
		rs.formatEnhancedRecommendationsForAI(recommendations),
	)
}

func (rs *RecommendationService) formatPurchaseHistory(history []dto.PurchaseItem) string {
	if len(history) == 0 {
		return "No recent purchases"
	}

	result := ""
	limit := 10
	if len(history) < limit {
		limit = len(history)
	}

	for i := 0; i < limit; i++ {
		item := history[i]
		result += fmt.Sprintf("- Product ID: %s, Category: %d, Price: %.2f, Date: %s\n",
			item.ProductID.String(), item.CategoryID, item.Price, item.PurchasedAt.Format("2006-01-02"))
	}

	return result
}

func (rs *RecommendationService) formatRecommendationsForAI(recommendations []dto.ProductRecommendation) string {
	result := ""
	for _, rec := range recommendations {
		result += fmt.Sprintf("- ID: %s, Name: %s, Category: %d, Price: %.2f, Rating: %.1f\n",
			rec.ProductID.String(), rec.Name, rec.CategoryID, rec.Price, rec.RatingAverage)
	}
	return result
}

// calculateAverageOrderValue calculates the average order value for a customer
func (rs *RecommendationService) calculateAverageOrderValue(profile *dto.CustomerProfile) float64 {
	if profile.OrderCount == 0 {
		return 0.0
	}
	return profile.TotalSpent / float64(profile.OrderCount)
}

// formatPriceRangePreference formats the customer's price range preference
func (rs *RecommendationService) formatPriceRangePreference(profile *dto.CustomerProfile) string {
	if profile.PriceRangeMin != nil && profile.PriceRangeMax != nil {
		return fmt.Sprintf("%.2f - %.2f", *profile.PriceRangeMin, *profile.PriceRangeMax)
	} else if profile.PriceRangeMin != nil {
		return fmt.Sprintf("Above %.2f", *profile.PriceRangeMin)
	} else if profile.PriceRangeMax != nil {
		return fmt.Sprintf("Below %.2f", *profile.PriceRangeMax)
	}
	return "No specific preference"
}

// formatDetailedPurchaseHistory provides detailed analysis of purchase history
func (rs *RecommendationService) formatDetailedPurchaseHistory(history []dto.PurchaseItem) string {
	if len(history) == 0 {
		return "No purchase history available"
	}

	result := "Recent Purchase Analysis:\n"

	// Analyze recent purchases (last 10)
	limit := 10
	if len(history) < limit {
		limit = len(history)
	}

	categoryCount := make(map[int]int)
	totalValue := 0.0
	var recentPurchases []string

	for i := 0; i < limit; i++ {
		item := history[i]
		categoryCount[item.CategoryID]++
		totalValue += item.Price

		daysSince := int(time.Since(item.PurchasedAt).Hours() / 24)
		recentPurchases = append(recentPurchases,
			fmt.Sprintf("  - Category %d, Price: %.2f, %d days ago",
				item.CategoryID, item.Price, daysSince))
	}

	// Add purchase pattern summary
	result += fmt.Sprintf("- Total recent purchases: %d\n", limit)
	result += fmt.Sprintf("- Recent purchase value: %.2f\n", totalValue)
	result += fmt.Sprintf("- Average recent purchase: %.2f\n", totalValue/float64(limit))

	// Add category distribution
	result += "- Category frequency: "
	for categoryID, count := range categoryCount {
		result += fmt.Sprintf("Cat%d(%d) ", categoryID, count)
	}
	result += "\n"

	// Add individual purchases
	result += "Recent Purchases:\n"
	for _, purchase := range recentPurchases {
		result += purchase + "\n"
	}

	return result
}

// formatRecentActivities formats customer's recent activities
func (rs *RecommendationService) formatRecentActivities(activities []dto.ActivityItem) string {
	if len(activities) == 0 {
		return "No recent activity data available"
	}

	result := "Recent Activity Analysis:\n"

	limit := 10
	if len(activities) < limit {
		limit = len(activities)
	}

	activityTypes := make(map[string]int)
	var recentActivities []string

	for i := 0; i < limit; i++ {
		activity := activities[i]
		activityTypes[activity.ActivityType]++

		daysSince := int(time.Since(activity.CreatedAt).Hours() / 24)
		details := activity.SearchQuery
		if activity.ProductID != nil {
			details = fmt.Sprintf("Product: %s", activity.ProductID.String())
		}
		recentActivities = append(recentActivities,
			fmt.Sprintf("  - %s: %s (%d days ago)",
				activity.ActivityType, details, daysSince))
	}

	// Add activity summary
	result += fmt.Sprintf("- Total recent activities: %d\n", limit)
	result += "- Activity types: "
	for activityType, count := range activityTypes {
		result += fmt.Sprintf("%s(%d) ", activityType, count)
	}
	result += "\n"

	// Add individual activities
	result += "Recent Activities:\n"
	for _, activity := range recentActivities {
		result += activity + "\n"
	}

	return result
}

// formatEnhancedRecommendationsForAI provides detailed product information for AI analysis
func (rs *RecommendationService) formatEnhancedRecommendationsForAI(recommendations []dto.ProductRecommendation) string {
	if len(recommendations) == 0 {
		return "No products to analyze"
	}

	result := "Product Analysis Data:\n"
	for i, rec := range recommendations {
		result += fmt.Sprintf(`
Product %d:
- ID: %s
- Name: %s
- Category: %d
- Price: %.2f
- Rating: %.1f (from %d reviews)
- Tags: %v
- Current Confidence: %.3f
- Description: %s

`,
			i+1,
			rec.ProductID.String(),
			rec.Name,
			rec.CategoryID,
			rec.Price,
			rec.RatingAverage,
			rec.RatingCount,
			rec.Tags,
			rec.ConfidenceScore,
			rec.Description)
	}
	return result
}

func (rs *RecommendationService) parseAIRecommendations(content string) ([]uuid.UUID, error) {
	// Remove markdown code blocks if present
	cleanContent := rs.cleanJSONFromMarkdown(content)

	var productIDs []string
	err := json.Unmarshal([]byte(cleanContent), &productIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to parse AI recommendations: %w", err)
	}

	var uuids []uuid.UUID
	for _, idStr := range productIDs {
		id, err := uuid.Parse(idStr)
		if err != nil {
			continue
		}
		uuids = append(uuids, id)
	}

	return uuids, nil
}

func (rs *RecommendationService) parseAIEnhancements(content string, recommendations []dto.ProductRecommendation) ([]dto.ProductRecommendation, error) {
	// Remove markdown code blocks if present
	cleanContent := rs.cleanJSONFromMarkdown(content)

	var response struct {
		Enhancements []struct {
			ProductID       string  `json:"product_id"`
			ConfidenceScore float64 `json:"confidence_score"`
			Reason          string  `json:"reason"`
		} `json:"enhancements"`
	}

	err := json.Unmarshal([]byte(cleanContent), &response)
	if err != nil {
		return recommendations, fmt.Errorf("failed to parse AI enhancements JSON: %w", err)
	}

	// Create a map for quick lookup
	enhancementMap := make(map[string]struct {
		ConfidenceScore float64
		Reason          string
	})

	for _, enhancement := range response.Enhancements {
		enhancementMap[enhancement.ProductID] = struct {
			ConfidenceScore float64
			Reason          string
		}{
			ConfidenceScore: enhancement.ConfidenceScore,
			Reason:          enhancement.Reason,
		}
	}

	// Apply enhancements to recommendations
	for i := range recommendations {
		productIDStr := recommendations[i].ProductID.String()
		if enhancement, exists := enhancementMap[productIDStr]; exists {
			recommendations[i].ConfidenceScore = enhancement.ConfidenceScore
			recommendations[i].Reason = enhancement.Reason
		}
	}

	return recommendations, nil
}

// cleanJSONFromMarkdown removes markdown code blocks from JSON content
func (rs *RecommendationService) cleanJSONFromMarkdown(content string) string {
	// Remove markdown code block markers
	content = strings.ReplaceAll(content, "```json", "")
	content = strings.ReplaceAll(content, "```", "")

	// Trim whitespace
	content = strings.TrimSpace(content)

	return content
}
