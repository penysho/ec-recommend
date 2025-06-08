package dto

import (
	"time"

	"github.com/google/uuid"
)

// RecommendationRequest represents a request for product recommendations
type RecommendationRequest struct {
	CustomerID         uuid.UUID  `json:"customer_id" binding:"required"`
	RecommendationType string     `json:"recommendation_type,omitempty"` // "similar", "collaborative", "content_based", "hybrid"
	ContextType        string     `json:"context_type,omitempty"`        // "homepage", "product_page", "cart", "checkout"
	ProductID          *uuid.UUID `json:"product_id,omitempty"`          // For product-based recommendations
	CategoryID         *int       `json:"category_id,omitempty"`         // For category-based recommendations
	Limit              int        `json:"limit,omitempty"`               // Number of recommendations to return (default: 10)
	ExcludeOwned       bool       `json:"exclude_owned,omitempty"`       // Exclude already purchased products
}

// RecommendationResponse represents the response containing product recommendations
type RecommendationResponse struct {
	CustomerID         uuid.UUID               `json:"customer_id"`
	Recommendations    []ProductRecommendation `json:"recommendations"`
	RecommendationType string                  `json:"recommendation_type"`
	ContextType        string                  `json:"context_type"`
	GeneratedAt        time.Time               `json:"generated_at"`
	Metadata           RecommendationMetadata  `json:"metadata"`
}

// ProductRecommendation represents a single product recommendation
type ProductRecommendation struct {
	ProductID       uuid.UUID `json:"product_id"`
	Name            string    `json:"name"`
	Description     string    `json:"description,omitempty"`
	Price           float64   `json:"price"`
	OriginalPrice   *float64  `json:"original_price,omitempty"`
	Brand           string    `json:"brand,omitempty"`
	CategoryID      int       `json:"category_id"`
	CategoryName    string    `json:"category_name"`
	RatingAverage   float64   `json:"rating_average"`
	RatingCount     int       `json:"rating_count"`
	PopularityScore int       `json:"popularity_score"`
	ConfidenceScore float64   `json:"confidence_score"` // AI-generated confidence
	Reason          string    `json:"reason"`           // AI-generated explanation
	Tags            []string  `json:"tags,omitempty"`
	ImageURL        string    `json:"image_url,omitempty"`
}

// RecommendationMetadata contains additional information about the recommendation process
type RecommendationMetadata struct {
	AlgorithmVersion string    `json:"algorithm_version"`
	ProcessingTimeMs int64     `json:"processing_time_ms"`
	TotalProducts    int       `json:"total_products"`
	FilteredProducts int       `json:"filtered_products"`
	AIModelUsed      string    `json:"ai_model_used,omitempty"`
	SessionID        uuid.UUID `json:"session_id,omitempty"`
}

// CustomerProfile represents customer data used for recommendations
type CustomerProfile struct {
	CustomerID          uuid.UUID      `json:"customer_id"`
	Email               string         `json:"email"`
	PreferredCategories []int          `json:"preferred_categories,omitempty"`
	PriceRangeMin       *float64       `json:"price_range_min,omitempty"`
	PriceRangeMax       *float64       `json:"price_range_max,omitempty"`
	PreferredBrands     []string       `json:"preferred_brands,omitempty"`
	LifestyleTags       []string       `json:"lifestyle_tags,omitempty"`
	IsPremium           bool           `json:"is_premium"`
	TotalSpent          float64        `json:"total_spent"`
	OrderCount          int            `json:"order_count"`
	PurchaseHistory     []PurchaseItem `json:"purchase_history,omitempty"`
	RecentActivities    []ActivityItem `json:"recent_activities,omitempty"`
}

// PurchaseItem represents a purchased product in customer history
type PurchaseItem struct {
	ProductID   uuid.UUID `json:"product_id"`
	CategoryID  int       `json:"category_id"`
	Price       float64   `json:"price"`
	Quantity    int       `json:"quantity"`
	PurchasedAt time.Time `json:"purchased_at"`
}

// ActivityItem represents customer activity data
type ActivityItem struct {
	ActivityType string     `json:"activity_type"`
	ProductID    *uuid.UUID `json:"product_id,omitempty"`
	SearchQuery  string     `json:"search_query,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}

// RecommendationAnalytics represents analytics data for recommendation performance
type RecommendationAnalytics struct {
	RecommendationID    uuid.UUID   `json:"recommendation_id"`
	CustomerID          uuid.UUID   `json:"customer_id"`
	RecommendedProducts []uuid.UUID `json:"recommended_products"`
	ClickedProducts     []uuid.UUID `json:"clicked_products,omitempty"`
	PurchasedProducts   []uuid.UUID `json:"purchased_products,omitempty"`
	ClickThroughRate    float64     `json:"click_through_rate"`
	ConversionRate      float64     `json:"conversion_rate"`
	CreatedAt           time.Time   `json:"created_at"`
}
