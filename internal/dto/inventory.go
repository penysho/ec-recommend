package dto

import (
	"time"

	"github.com/google/uuid"
)

// InventoryItem represents inventory information for a product
type InventoryItem struct {
	ProductID     uuid.UUID `json:"product_id"`
	ProductName   string    `json:"product_name"`
	SKU           string    `json:"sku,omitempty"`
	StockQuantity int       `json:"stock_quantity"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
	UpdatedAt     time.Time `json:"updated_at,omitempty"`
}

// InventoryListResponse represents a list of inventory items
type InventoryListResponse struct {
	Items      []InventoryItem `json:"items"`
	Total      int             `json:"total"`
	Page       int             `json:"page"`
	PageSize   int             `json:"page_size"`
	TotalPages int             `json:"total_pages"`
}

// InventoryUpdateRequest represents a request to update inventory
type InventoryUpdateRequest struct {
	ProductID     uuid.UUID `json:"product_id" binding:"required"`
	StockQuantity int       `json:"stock_quantity" binding:"gte=0"`
	Reason        string    `json:"reason,omitempty"` // Optional reason for the update
}

// InventoryBatchUpdateRequest represents a request to update multiple inventory items
type InventoryBatchUpdateRequest struct {
	Updates []InventoryUpdateRequest `json:"updates" binding:"required,dive"`
}

// InventoryTransactionType represents the type of inventory transaction
type InventoryTransactionType string

const (
	TransactionTypeIncrease InventoryTransactionType = "increase"
	TransactionTypeDecrease InventoryTransactionType = "decrease"
	TransactionTypeSet      InventoryTransactionType = "set"
)

// InventoryTransactionRequest represents a request to perform inventory transaction
type InventoryTransactionRequest struct {
	ProductID       uuid.UUID                `json:"product_id" binding:"required"`
	TransactionType InventoryTransactionType `json:"transaction_type" binding:"required,oneof=increase decrease set"`
	Quantity        int                      `json:"quantity" binding:"required,gt=0"`
	Reason          string                   `json:"reason,omitempty"`
}

// InventoryAdjustmentLog represents an inventory adjustment log entry
type InventoryAdjustmentLog struct {
	ID              uuid.UUID                `json:"id"`
	ProductID       uuid.UUID                `json:"product_id"`
	ProductName     string                   `json:"product_name"`
	TransactionType InventoryTransactionType `json:"transaction_type"`
	PreviousStock   int                      `json:"previous_stock"`
	Quantity        int                      `json:"quantity"`
	NewStock        int                      `json:"new_stock"`
	Reason          string                   `json:"reason,omitempty"`
	PerformedBy     string                   `json:"performed_by"`
	PerformedAt     time.Time                `json:"performed_at"`
}

// InventoryHistoryResponse represents inventory adjustment history
type InventoryHistoryResponse struct {
	ProductID uuid.UUID                `json:"product_id"`
	Logs      []InventoryAdjustmentLog `json:"logs"`
	Total     int                      `json:"total"`
	Page      int                      `json:"page"`
	PageSize  int                      `json:"page_size"`
}

// InventoryStockAlert represents a low stock alert
type InventoryStockAlert struct {
	ProductID        uuid.UUID `json:"product_id"`
	ProductName      string    `json:"product_name"`
	SKU              string    `json:"sku,omitempty"`
	CurrentStock     int       `json:"current_stock"`
	MinimumThreshold int       `json:"minimum_threshold"`
	Status           string    `json:"status"` // "low_stock", "out_of_stock"
}

// InventoryAlertsResponse represents a list of inventory alerts
type InventoryAlertsResponse struct {
	Alerts []InventoryStockAlert `json:"alerts"`
	Total  int                   `json:"total"`
}

// InventorySearchRequest represents a request to search inventory
type InventorySearchRequest struct {
	Query         string     `json:"query,omitempty"`
	CategoryID    *int       `json:"category_id,omitempty"`
	SKU           string     `json:"sku,omitempty"`
	MinStock      *int       `json:"min_stock,omitempty"`
	MaxStock      *int       `json:"max_stock,omitempty"`
	StockStatus   string     `json:"stock_status,omitempty"` // "in_stock", "low_stock", "out_of_stock"
	Page          int        `json:"page"`
	PageSize      int        `json:"page_size"`
	SortBy        string     `json:"sort_by,omitempty"`        // "stock_quantity", "product_name", "updated_at"
	SortDirection string     `json:"sort_direction,omitempty"` // "asc", "desc"
}

// InventoryStatsResponse represents inventory statistics
type InventoryStatsResponse struct {
	TotalProducts       int     `json:"total_products"`
	InStockProducts     int     `json:"in_stock_products"`
	LowStockProducts    int     `json:"low_stock_products"`
	OutOfStockProducts  int     `json:"out_of_stock_products"`
	TotalStockValue     float64 `json:"total_stock_value"`
	AverageStockLevel   float64 `json:"average_stock_level"`
	LastUpdated         time.Time `json:"last_updated"`
}
