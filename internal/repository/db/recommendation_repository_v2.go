package repository

import (
	"context"
	"database/sql"
	"ec-recommend/internal/dto"
	"ec-recommend/internal/repository/db/models"
	"ec-recommend/internal/service"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// RecommendationRepositoryV2 implements the RecommendationRepositoryV2Interface
type RecommendationRepositoryV2 struct {
	db boil.ContextExecutor
}

// NewRecommendationRepositoryV2 creates a new recommendation repository v2 instance
func NewRecommendationRepositoryV2(db *sql.DB) service.RecommendationRepositoryV2Interface {
	return &RecommendationRepositoryV2{
		db: db,
	}
}

// GetCustomerByID retrieves customer profile data by ID
func (r *RecommendationRepositoryV2) GetCustomerByID(ctx context.Context, customerID uuid.UUID) (*dto.CustomerProfile, error) {
	query := `
		SELECT
			id, email, first_name, last_name, preferred_categories,
			price_range_min, price_range_max, preferred_brands,
			lifestyle_tags, is_premium, total_spent, order_count
		FROM customers
		WHERE id = $1
	`

	var profile dto.CustomerProfile
	var firstName, lastName sql.NullString
	var priceRangeMin, priceRangeMax sql.NullFloat64
	var isPremium sql.NullBool
	var totalSpent sql.NullFloat64
	var orderCount sql.NullInt64
	var customerIDStr string
	var preferredCategoriesArray pq.Int64Array

	db := r.db.(*sql.DB)
	err := db.QueryRowContext(ctx, query, customerID.String()).Scan(
		&customerIDStr,
		&profile.Email,
		&firstName,
		&lastName,
		&preferredCategoriesArray,
		&priceRangeMin,
		&priceRangeMax,
		pq.Array(&profile.PreferredBrands),
		pq.Array(&profile.LifestyleTags),
		&isPremium,
		&totalSpent,
		&orderCount,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("customer not found: %s", customerID)
		}
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}

	// Convert UUID string to UUID
	profile.CustomerID = customerID

	// Convert Int64Array to []int for PreferredCategories
	profile.PreferredCategories = make([]int, len(preferredCategoriesArray))
	for i, v := range preferredCategoriesArray {
		profile.PreferredCategories[i] = int(v)
	}

	// Handle nullable fields
	profile.IsPremium = isPremium.Bool
	profile.TotalSpent = totalSpent.Float64
	profile.OrderCount = int(orderCount.Int64)

	if priceRangeMin.Valid {
		profile.PriceRangeMin = &priceRangeMin.Float64
	}
	if priceRangeMax.Valid {
		profile.PriceRangeMax = &priceRangeMax.Float64
	}

	return &profile, nil
}

// GetCustomerPurchaseHistory retrieves customer's purchase history
func (r *RecommendationRepositoryV2) GetCustomerPurchaseHistory(ctx context.Context, customerID uuid.UUID, limit int) ([]dto.PurchaseItem, error) {
	orderItems, err := models.OrderItems(
		qm.InnerJoin("orders o ON order_items.order_id = o.id"),
		qm.InnerJoin("products p ON order_items.product_id = p.id"),
		qm.Where("o.customer_id = ? AND o.status = ?", customerID.String(), "delivered"),
		qm.OrderBy("o.ordered_at DESC"),
		qm.Limit(limit),
		qm.Load("Order"),
		qm.Load("Product"),
	).All(ctx, r.db)

	if err != nil {
		return nil, fmt.Errorf("failed to query purchase history: %w", err)
	}

	purchases := make([]dto.PurchaseItem, len(orderItems))
	for i, item := range orderItems {
		// Convert unit price from types.Decimal to float64
		unitPrice, ok := item.UnitPrice.Float64()
		if !ok {
			return nil, fmt.Errorf("failed to convert unit price to float64")
		}

		// Get ordered date from the loaded Order relation
		var orderedAt time.Time
		if item.R != nil && item.R.Order != nil && item.R.Order.OrderedAt.Valid {
			orderedAt = item.R.Order.OrderedAt.Time
		}

		// Get category ID from the loaded Product relation
		var categoryID int
		if item.R != nil && item.R.Product != nil {
			categoryID = item.R.Product.CategoryID
		}

		purchases[i] = dto.PurchaseItem{
			ProductID:   uuid.MustParse(item.ProductID),
			CategoryID:  categoryID,
			Price:       unitPrice,
			Quantity:    item.Quantity,
			PurchasedAt: orderedAt,
		}
	}

	return purchases, nil
}

// GetCustomerActivities retrieves customer's recent activities
func (r *RecommendationRepositoryV2) GetCustomerActivities(ctx context.Context, customerID uuid.UUID, limit int) ([]dto.ActivityItem, error) {
	activities, err := models.CustomerActivities(
		models.CustomerActivityWhere.CustomerID.EQ(customerID.String()),
		qm.OrderBy("created_at DESC"),
		qm.Limit(limit),
	).All(ctx, r.db)

	if err != nil {
		return nil, fmt.Errorf("failed to query activities: %w", err)
	}

	result := make([]dto.ActivityItem, len(activities))
	for i, activity := range activities {
		item := dto.ActivityItem{
			ActivityType: activity.ActivityType,
			CreatedAt:    activity.CreatedAt.Time,
		}

		if activity.ProductID.Valid {
			productID := uuid.MustParse(activity.ProductID.String)
			item.ProductID = &productID
		}

		if activity.SearchQuery.Valid {
			item.SearchQuery = activity.SearchQuery.String
		}

		result[i] = item
	}

	return result, nil
}

// GetProductsByIDs retrieves products by their IDs and returns ProductRecommendationV2
func (r *RecommendationRepositoryV2) GetProductsByIDs(ctx context.Context, productIDs []uuid.UUID) ([]dto.ProductRecommendationV2, error) {
	if len(productIDs) == 0 {
		return []dto.ProductRecommendationV2{}, nil
	}

	// Convert UUID slice to string slice
	stringIDs := make([]interface{}, len(productIDs))
	for i, id := range productIDs {
		stringIDs[i] = id.String()
	}

	products, err := models.Products(
		qm.InnerJoin("categories c ON products.category_id = c.id"),
		qm.WhereIn("products.id IN ?", stringIDs...),
		models.ProductWhere.IsActive.EQ(null.BoolFrom(true)),
		qm.Load("Category"),
	).All(ctx, r.db)

	if err != nil {
		return nil, fmt.Errorf("failed to query products: %w", err)
	}

	return r.convertToProductRecommendationsV2(products), nil
}

// GetProductsByCategory retrieves products by category and returns ProductRecommendationV2
func (r *RecommendationRepositoryV2) GetProductsByCategory(ctx context.Context, categoryID int, limit int) ([]dto.ProductRecommendationV2, error) {
	products, err := models.Products(
		qm.InnerJoin("categories c ON products.category_id = c.id"),
		models.ProductWhere.CategoryID.EQ(categoryID),
		models.ProductWhere.IsActive.EQ(null.BoolFrom(true)),
		qm.OrderBy("popularity_score DESC"),
		qm.Limit(limit),
		qm.Load("Category"),
	).All(ctx, r.db)

	if err != nil {
		return nil, fmt.Errorf("failed to query products by category: %w", err)
	}

	return r.convertToProductRecommendationsV2(products), nil
}

// convertToProductRecommendationsV2 converts model products to ProductRecommendationV2 DTOs
func (r *RecommendationRepositoryV2) convertToProductRecommendationsV2(products models.ProductSlice) []dto.ProductRecommendationV2 {
	recommendations := make([]dto.ProductRecommendationV2, len(products))

	for i, product := range products {
		var categoryName string
		if product.R != nil && product.R.Category != nil {
			categoryName = product.R.Category.Name
		}

		// Convert price from types.Decimal to float64
		price, _ := product.Price.Float64()

		var originalPrice *float64
		if !product.OriginalPrice.IsZero() {
			origPrice, _ := product.OriginalPrice.Float64()
			originalPrice = &origPrice
		}

		// Convert rating average
		var ratingAverage float64
		if !product.RatingAverage.IsZero() {
			ratingAverage, _ = product.RatingAverage.Float64()
		}

		// Convert rating count
		var ratingCount int
		if product.RatingCount.Valid {
			ratingCount = product.RatingCount.Int
		}

		// Convert popularity score
		var popularityScore int
		if product.PopularityScore.Valid {
			popularityScore = product.PopularityScore.Int
		}

		// Convert tags
		var tags []string
		if len(product.Tags) > 0 {
			tags = product.Tags
		} else {
			tags = []string{}
		}

		recommendations[i] = dto.ProductRecommendationV2{
			ProductID:       uuid.MustParse(product.ID),
			Name:            product.Name,
			Description:     product.Description.String,
			Price:           price,
			OriginalPrice:   originalPrice,
			Brand:           product.Brand.String,
			CategoryID:      product.CategoryID,
			CategoryName:    categoryName,
			RatingAverage:   ratingAverage,
			RatingCount:     ratingCount,
			PopularityScore: popularityScore,
			ConfidenceScore: 0.8,                                // Default confidence score - would be enhanced with AI models
			SimilarityScore: 0.0,                                // Default similarity score - would be enhanced with vector search
			Reason:          "Product matches your preferences", // Default reason - would be enhanced with AI
			Tags:            tags,
			ImageURL:        "", // ImageURL field doesn't exist in the model, using empty string
		}
	}

	return recommendations
}

// ==== TODO: Below methods are placeholder implementations ====

// GetTrendingProductsV2 - TODO: Implement enhanced trending products analysis
func (r *RecommendationRepositoryV2) GetTrendingProductsV2(ctx context.Context, categoryID *int, timeRange string, limit int) ([]dto.TrendingProductV2, error) {
	// TODO: Implement advanced trending analysis with AI insights
	return nil, fmt.Errorf("trending products v2 not implemented yet")
}

// GetProductPerformanceMetrics - TODO: Implement product performance analytics
func (r *RecommendationRepositoryV2) GetProductPerformanceMetrics(ctx context.Context, productIDs []uuid.UUID) (map[uuid.UUID]*dto.ProductPerformanceMetrics, error) {
	// TODO: Implement product performance metrics calculation
	return nil, fmt.Errorf("product performance metrics not implemented yet")
}

// GetMarketAnalysis - TODO: Implement market analysis functionality
func (r *RecommendationRepositoryV2) GetMarketAnalysis(ctx context.Context, categoryID *int, timeRange string) (*dto.MarketAnalysis, error) {
	// TODO: Implement market analysis with competitor insights
	return nil, fmt.Errorf("market analysis not implemented yet")
}

// LogRecommendation - TODO: Implement enhanced recommendation logging
func (r *RecommendationRepositoryV2) LogRecommendation(ctx context.Context, customerID uuid.UUID, recommendationType, contextType string, productIDs []uuid.UUID, sessionID uuid.UUID) error {
	// TODO: Implement enhanced recommendation logging with session tracking
	return fmt.Errorf("recommendation logging not implemented yet")
}

// LogRecommendationInteraction - TODO: Implement recommendation interaction analytics
func (r *RecommendationRepositoryV2) LogRecommendationInteraction(ctx context.Context, analytics *dto.RecommendationAnalytics) error {
	// TODO: Implement interaction analytics logging
	return fmt.Errorf("recommendation interaction logging not implemented yet")
}

// LogSemanticSearch - TODO: Implement semantic search logging
func (r *RecommendationRepositoryV2) LogSemanticSearch(ctx context.Context, customerID *uuid.UUID, query string, results []uuid.UUID, processingTimeMs int64) error {
	// TODO: Implement semantic search result logging
	return fmt.Errorf("semantic search logging not implemented yet")
}

// GetCachedRecommendations - TODO: Implement cache retrieval
func (r *RecommendationRepositoryV2) GetCachedRecommendations(ctx context.Context, key string) ([]dto.ProductRecommendationV2, error) {
	// TODO: Implement Redis or similar cache retrieval
	return nil, fmt.Errorf("cache retrieval not implemented yet")
}

// SetCachedRecommendations - TODO: Implement cache storage
func (r *RecommendationRepositoryV2) SetCachedRecommendations(ctx context.Context, key string, recommendations []dto.ProductRecommendationV2, ttl int64) error {
	// TODO: Implement Redis or similar cache storage
	return fmt.Errorf("cache storage not implemented yet")
}

// InvalidateCache - TODO: Implement cache invalidation
func (r *RecommendationRepositoryV2) InvalidateCache(ctx context.Context, pattern string) error {
	// TODO: Implement cache pattern invalidation
	return fmt.Errorf("cache invalidation not implemented yet")
}
