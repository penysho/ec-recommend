package repository

import (
	"context"
	"database/sql"
	"ec-recommend/internal/dto"
	"ec-recommend/internal/models"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// RecommendationRepository implements the interfaces.RecommendationRepositoryInterface
type RecommendationRepository struct {
	db boil.ContextExecutor
}

// NewRecommendationRepository creates a new recommendation repository instance
func NewRecommendationRepository(db *sql.DB) *RecommendationRepository {
	return &RecommendationRepository{
		db: db,
	}
}

// GetCustomerByID retrieves customer profile data by ID
func (r *RecommendationRepository) GetCustomerByID(ctx context.Context, customerID uuid.UUID) (*dto.CustomerProfile, error) {
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
func (r *RecommendationRepository) GetCustomerPurchaseHistory(ctx context.Context, customerID uuid.UUID, limit int) ([]dto.PurchaseItem, error) {
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
func (r *RecommendationRepository) GetCustomerActivities(ctx context.Context, customerID uuid.UUID, limit int) ([]dto.ActivityItem, error) {
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

// GetProductsByIDs retrieves products by their IDs
func (r *RecommendationRepository) GetProductsByIDs(ctx context.Context, productIDs []uuid.UUID) ([]dto.ProductRecommendation, error) {
	if len(productIDs) == 0 {
		return []dto.ProductRecommendation{}, nil
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

	return r.convertToProductRecommendations(products), nil
}

// GetProductsByCategory retrieves products by category
func (r *RecommendationRepository) GetProductsByCategory(ctx context.Context, categoryID int, limit int) ([]dto.ProductRecommendation, error) {
	products, err := models.Products(
		models.ProductWhere.CategoryID.EQ(categoryID),
		models.ProductWhere.IsActive.EQ(null.BoolFrom(true)),
		qm.Load("Category"),
		qm.OrderBy("rating_average DESC, popularity_score DESC"),
		qm.Limit(limit),
	).All(ctx, r.db)

	if err != nil {
		return nil, fmt.Errorf("failed to query products by category: %w", err)
	}

	return r.convertToProductRecommendations(products), nil
}

// GetTrendingProducts retrieves trending products
func (r *RecommendationRepository) GetTrendingProducts(ctx context.Context, categoryID *int, limit int) ([]dto.ProductRecommendation, error) {
	var mods []qm.QueryMod

	mods = append(mods,
		models.ProductWhere.IsActive.EQ(null.BoolFrom(true)),
		qm.Load("Category"),
		qm.OrderBy("popularity_score DESC, rating_average DESC"),
		qm.Limit(limit),
	)

	if categoryID != nil {
		mods = append(mods, models.ProductWhere.CategoryID.EQ(*categoryID))
	}

	products, err := models.Products(mods...).All(ctx, r.db)
	if err != nil {
		return nil, fmt.Errorf("failed to query trending products: %w", err)
	}

	return r.convertToProductRecommendations(products), nil
}

// GetSimilarProductsByTags retrieves similar products based on tags
func (r *RecommendationRepository) GetSimilarProductsByTags(ctx context.Context, tags []string, excludeProductID uuid.UUID, limit int) ([]dto.ProductRecommendation, error) {
	if len(tags) == 0 {
		return []dto.ProductRecommendation{}, nil
	}

	products, err := models.Products(
		qm.Where("tags && ?", pq.Array(tags)),
		models.ProductWhere.ID.NEQ(excludeProductID.String()),
		models.ProductWhere.IsActive.EQ(null.BoolFrom(true)),
		qm.Load("Category"),
		qm.OrderBy("rating_average DESC"),
		qm.Limit(limit),
	).All(ctx, r.db)

	if err != nil {
		return nil, fmt.Errorf("failed to query similar products: %w", err)
	}

	return r.convertToProductRecommendations(products), nil
}

// GetProductsInPriceRange retrieves products within a price range
func (r *RecommendationRepository) GetProductsInPriceRange(ctx context.Context, minPrice, maxPrice float64, limit int) ([]dto.ProductRecommendation, error) {
	products, err := models.Products(
		qm.Where("price BETWEEN ? AND ?", minPrice, maxPrice),
		models.ProductWhere.IsActive.EQ(null.BoolFrom(true)),
		qm.Load("Category"),
		qm.OrderBy("rating_average DESC, popularity_score DESC"),
		qm.Limit(limit),
	).All(ctx, r.db)

	if err != nil {
		return nil, fmt.Errorf("failed to query products in price range: %w", err)
	}

	return r.convertToProductRecommendations(products), nil
}

// LogRecommendation logs a recommendation event
func (r *RecommendationRepository) LogRecommendation(ctx context.Context, customerID uuid.UUID, recommendationType, contextType string, productIDs []uuid.UUID, sessionID uuid.UUID) error {
	// Convert UUID slice to string slice for PostgreSQL array
	stringIDs := make([]string, len(productIDs))
	for i, id := range productIDs {
		stringIDs[i] = id.String()
	}

	log := &models.RecommendationLog{
		CustomerID:          customerID.String(),
		SessionID:           null.StringFrom(sessionID.String()),
		RecommendationType:  recommendationType,
		ContextType:         null.StringFrom(contextType),
		RecommendedProducts: stringIDs,
		ClickedProducts:     []string{},
		PurchasedProducts:   []string{},
		AlgorithmVersion:    null.StringFrom("v1.0"),
	}

	return log.Insert(ctx, r.db, boil.Infer())
}

// LogRecommendationInteraction logs interaction analytics
func (r *RecommendationRepository) LogRecommendationInteraction(ctx context.Context, analytics *dto.RecommendationAnalytics) error {
	// Convert UUID slices to string slices
	recommendedProducts := make([]string, len(analytics.RecommendedProducts))
	for i, id := range analytics.RecommendedProducts {
		recommendedProducts[i] = id.String()
	}

	clickedProducts := make([]string, len(analytics.ClickedProducts))
	for i, id := range analytics.ClickedProducts {
		clickedProducts[i] = id.String()
	}

	purchasedProducts := make([]string, len(analytics.PurchasedProducts))
	for i, id := range analytics.PurchasedProducts {
		purchasedProducts[i] = id.String()
	}

	log := &models.RecommendationLog{
		ID:                  analytics.RecommendationID.String(),
		CustomerID:          analytics.CustomerID.String(),
		RecommendedProducts: recommendedProducts,
		ClickedProducts:     clickedProducts,
		PurchasedProducts:   purchasedProducts,
	}

	_, err := log.Update(ctx, r.db, boil.Infer())
	return err
}

// GetCustomersWithSimilarPurchases finds customers with similar purchase patterns
func (r *RecommendationRepository) GetCustomersWithSimilarPurchases(ctx context.Context, customerID uuid.UUID, limit int) ([]uuid.UUID, error) {
	// Complex query for finding similar customers based on purchase overlap
	query := `
		WITH customer_products AS (
			SELECT DISTINCT oi.product_id
			FROM order_items oi
			JOIN orders o ON oi.order_id = o.id
			WHERE o.customer_id = $1 AND o.status = 'delivered'
		),
		similar_customers AS (
			SELECT
				o.customer_id,
				COUNT(DISTINCT oi.product_id) as shared_products,
				COUNT(DISTINCT cp.product_id) as total_customer_products
			FROM order_items oi
			JOIN orders o ON oi.order_id = o.id
			JOIN customer_products cp ON oi.product_id = cp.product_id
			WHERE o.customer_id != $1 AND o.status = 'delivered'
			GROUP BY o.customer_id
			HAVING COUNT(DISTINCT oi.product_id) >= 2
		)
		SELECT customer_id
		FROM similar_customers
		ORDER BY shared_products DESC, total_customer_products DESC
		LIMIT $2
	`

	db := r.db.(*sql.DB)
	rows, err := db.QueryContext(ctx, query, customerID.String(), limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query similar customers: %w", err)
	}
	defer rows.Close()

	var similarCustomers []uuid.UUID
	for rows.Next() {
		var customerIDStr string
		if err := rows.Scan(&customerIDStr); err != nil {
			return nil, fmt.Errorf("failed to scan similar customer: %w", err)
		}

		customerUUID, err := uuid.Parse(customerIDStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse customer UUID: %w", err)
		}

		similarCustomers = append(similarCustomers, customerUUID)
	}

	return similarCustomers, nil
}

// GetPopularProductsAmongSimilarCustomers retrieves popular products among similar customers
func (r *RecommendationRepository) GetPopularProductsAmongSimilarCustomers(ctx context.Context, similarCustomerIDs []uuid.UUID, excludeOwned []uuid.UUID, limit int) ([]dto.ProductRecommendation, error) {
	if len(similarCustomerIDs) == 0 {
		return []dto.ProductRecommendation{}, nil
	}

	// Convert UUID slices to string slices
	similarCustomerStrings := make([]interface{}, len(similarCustomerIDs))
	for i, id := range similarCustomerIDs {
		similarCustomerStrings[i] = id.String()
	}

	excludeOwnedStrings := make([]interface{}, len(excludeOwned))
	for i, id := range excludeOwned {
		excludeOwnedStrings[i] = id.String()
	}

	var mods []qm.QueryMod

	mods = append(mods,
		qm.InnerJoin("order_items oi ON products.id = oi.product_id"),
		qm.InnerJoin("orders o ON oi.order_id = o.id"),
		qm.WhereIn("o.customer_id IN ?", similarCustomerStrings...),
		qm.Where("o.status = ?", "delivered"),
		models.ProductWhere.IsActive.EQ(null.BoolFrom(true)),
		qm.Load("Category"),
	)

	if len(excludeOwned) > 0 {
		mods = append(mods, qm.WhereNotIn("products.id NOT IN ?", excludeOwnedStrings...))
	}

	mods = append(mods,
		qm.GroupBy("products.id, products.name, products.description, products.price, products.original_price, products.brand, products.category_id, products.rating_average, products.rating_count, products.popularity_score, products.tags"),
		qm.OrderBy("COUNT(DISTINCT o.customer_id) DESC, products.rating_average DESC"),
		qm.Limit(limit),
	)

	products, err := models.Products(mods...).All(ctx, r.db)
	if err != nil {
		return nil, fmt.Errorf("failed to query popular products among similar customers: %w", err)
	}

	return r.convertToProductRecommendations(products), nil
}

// convertToProductRecommendations converts SQLBoiler models to DTO
func (r *RecommendationRepository) convertToProductRecommendations(products models.ProductSlice) []dto.ProductRecommendation {
	recommendations := make([]dto.ProductRecommendation, len(products))

	for i, product := range products {
		price, _ := product.Price.Float64()
		ratingAverage, _ := product.RatingAverage.Float64()

		var tags []string
		if product.Tags != nil {
			tags = product.Tags
		}

		rec := dto.ProductRecommendation{
			ProductID:       uuid.MustParse(product.ID),
			Name:            product.Name,
			Price:           price,
			CategoryID:      product.CategoryID,
			RatingAverage:   ratingAverage,
			RatingCount:     product.RatingCount.Int,
			PopularityScore: product.PopularityScore.Int,
			Tags:            tags,
		}

		if product.Description.Valid {
			rec.Description = product.Description.String
		}

		if originalPrice, ok := product.OriginalPrice.Float64(); ok {
			rec.OriginalPrice = &originalPrice
		}

		if product.Brand.Valid {
			rec.Brand = product.Brand.String
		}

		if product.R != nil && product.R.Category != nil {
			rec.CategoryName = product.R.Category.Name
		}

		recommendations[i] = rec
	}

	return recommendations
}
