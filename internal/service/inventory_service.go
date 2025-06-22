package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"ec-recommend/internal/dto"
	"ec-recommend/internal/repository/db/models"

	"github.com/google/uuid"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// InventoryServiceInterface defines the interface for inventory service operations
type InventoryServiceInterface interface {
	// GetInventoryList retrieves a paginated list of inventory items
	GetInventoryList(ctx context.Context, req *dto.InventorySearchRequest) (*dto.InventoryListResponse, error)

	// GetInventoryByProductID retrieves inventory information for a specific product
	GetInventoryByProductID(ctx context.Context, productID uuid.UUID) (*dto.InventoryItem, error)

	// UpdateInventory updates the stock quantity for a product
	UpdateInventory(ctx context.Context, req *dto.InventoryUpdateRequest, userID string) error

	// PerformInventoryTransaction performs an inventory transaction (increase, decrease, set)
	PerformInventoryTransaction(ctx context.Context, req *dto.InventoryTransactionRequest, userID string) error

	// BatchUpdateInventory updates multiple inventory items
	BatchUpdateInventory(ctx context.Context, req *dto.InventoryBatchUpdateRequest, userID string) error

	// GetInventoryHistory retrieves inventory adjustment history for a product
	GetInventoryHistory(ctx context.Context, productID uuid.UUID, page, pageSize int) (*dto.InventoryHistoryResponse, error)

	// GetInventoryAlerts retrieves low stock alerts
	GetInventoryAlerts(ctx context.Context, threshold int) (*dto.InventoryAlertsResponse, error)

	// GetInventoryStats retrieves overall inventory statistics
	GetInventoryStats(ctx context.Context) (*dto.InventoryStatsResponse, error)
}

// InventoryService implements inventory business logic
type InventoryService struct {
	db *sql.DB
}

// NewInventoryService creates a new inventory service
func NewInventoryService(db *sql.DB) *InventoryService {
	return &InventoryService{
		db: db,
	}
}

// GetInventoryList retrieves a paginated list of inventory items
func (s *InventoryService) GetInventoryList(ctx context.Context, req *dto.InventorySearchRequest) (*dto.InventoryListResponse, error) {
	// Build query modifiers
	queryMods := []qm.QueryMod{
		qm.Select("id", "name", "sku", "stock_quantity", "created_at", "updated_at"),
	}

	// Add search filters
	if req.Query != "" {
		queryMods = append(queryMods, qm.Where("name ILIKE ?", "%"+req.Query+"%"))
	}

	if req.CategoryID != nil {
		queryMods = append(queryMods, qm.Where("category_id = ?", *req.CategoryID))
	}

	if req.SKU != "" {
		queryMods = append(queryMods, qm.Where("sku = ?", req.SKU))
	}

	if req.MinStock != nil {
		queryMods = append(queryMods, qm.Where("stock_quantity >= ?", *req.MinStock))
	}

	if req.MaxStock != nil {
		queryMods = append(queryMods, qm.Where("stock_quantity <= ?", *req.MaxStock))
	}

	// Add stock status filter
	if req.StockStatus != "" {
		switch req.StockStatus {
		case "in_stock":
			queryMods = append(queryMods, qm.Where("stock_quantity > 10"))
		case "low_stock":
			queryMods = append(queryMods, qm.Where("stock_quantity > 0 AND stock_quantity <= 10"))
		case "out_of_stock":
			queryMods = append(queryMods, qm.Where("stock_quantity = 0"))
		}
	}

	// Add sorting
	if req.SortBy != "" {
		direction := "ASC"
		if req.SortDirection == "desc" {
			direction = "DESC"
		}
		queryMods = append(queryMods, qm.OrderBy(req.SortBy+" "+direction))
	} else {
		queryMods = append(queryMods, qm.OrderBy("name ASC"))
	}

	// Count total items
	countMods := make([]qm.QueryMod, len(queryMods)-2) // Exclude SELECT and ORDER BY
	copy(countMods, queryMods[1:len(queryMods)-1])
	total, err := models.Products(countMods...).Count(ctx, s.db)
	if err != nil {
		return nil, fmt.Errorf("failed to count inventory items: %w", err)
	}

	// Add pagination
	offset := (req.Page - 1) * req.PageSize
	queryMods = append(queryMods, qm.Limit(req.PageSize), qm.Offset(offset))

	// Execute query
	products, err := models.Products(queryMods...).All(ctx, s.db)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve inventory items: %w", err)
	}

	// Convert to DTOs
	items := make([]dto.InventoryItem, len(products))
	for i, product := range products {
		items[i] = dto.InventoryItem{
			ProductID:     uuid.MustParse(product.ID),
			ProductName:   product.Name,
			SKU:           product.Sku.String,
			StockQuantity: product.StockQuantity.Int,
			CreatedAt:     product.CreatedAt.Time,
			UpdatedAt:     product.UpdatedAt.Time,
		}
	}

	totalPages := int(total+int64(req.PageSize)-1) / req.PageSize

	return &dto.InventoryListResponse{
		Items:      items,
		Total:      int(total),
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// GetInventoryByProductID retrieves inventory information for a specific product
func (s *InventoryService) GetInventoryByProductID(ctx context.Context, productID uuid.UUID) (*dto.InventoryItem, error) {
	product, err := models.FindProduct(ctx, s.db, productID.String())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product not found: %s", productID.String())
		}
		return nil, fmt.Errorf("failed to find product: %w", err)
	}

	return &dto.InventoryItem{
		ProductID:     uuid.MustParse(product.ID),
		ProductName:   product.Name,
		SKU:           product.Sku.String,
		StockQuantity: product.StockQuantity.Int,
		CreatedAt:     product.CreatedAt.Time,
		UpdatedAt:     product.UpdatedAt.Time,
	}, nil
}

// UpdateInventory updates the stock quantity for a product
func (s *InventoryService) UpdateInventory(ctx context.Context, req *dto.InventoryUpdateRequest, userID string) error {
	// Start transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Get current product
	product, err := models.FindProduct(ctx, tx, req.ProductID.String())
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("product not found: %s", req.ProductID.String())
		}
		return fmt.Errorf("failed to find product: %w", err)
	}

	previousStock := product.StockQuantity.Int

	// Update stock quantity
	product.StockQuantity = null.NewInt(req.StockQuantity, true)
	product.UpdatedAt = null.NewTime(time.Now(), true)

	if _, err := product.Update(ctx, tx, boil.Infer()); err != nil {
		return fmt.Errorf("failed to update product stock: %w", err)
	}

	// Log the inventory adjustment (in a real system, you'd have an inventory_logs table)
	// For now, we'll skip this part since we don't have the table defined

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// PerformInventoryTransaction performs an inventory transaction (increase, decrease, set)
func (s *InventoryService) PerformInventoryTransaction(ctx context.Context, req *dto.InventoryTransactionRequest, userID string) error {
	// Start transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Get current product
	product, err := models.FindProduct(ctx, tx, req.ProductID.String())
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("product not found: %s", req.ProductID.String())
		}
		return fmt.Errorf("failed to find product: %w", err)
	}

	previousStock := product.StockQuantity.Int
	var newStock int

	// Calculate new stock based on transaction type
	switch req.TransactionType {
	case dto.TransactionTypeIncrease:
		newStock = previousStock + req.Quantity
	case dto.TransactionTypeDecrease:
		newStock = previousStock - req.Quantity
		if newStock < 0 {
			return fmt.Errorf("insufficient stock: current stock is %d, cannot decrease by %d", previousStock, req.Quantity)
		}
	case dto.TransactionTypeSet:
		newStock = req.Quantity
	default:
		return fmt.Errorf("invalid transaction type: %s", req.TransactionType)
	}

	// Update stock quantity
	product.StockQuantity = null.NewInt(newStock, true)
	product.UpdatedAt = null.NewTime(time.Now(), true)

	if _, err := product.Update(ctx, tx, boil.Infer()); err != nil {
		return fmt.Errorf("failed to update product stock: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// BatchUpdateInventory updates multiple inventory items
func (s *InventoryService) BatchUpdateInventory(ctx context.Context, req *dto.InventoryBatchUpdateRequest, userID string) error {
	// Start transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Process each update
	for _, update := range req.Updates {
		product, err := models.FindProduct(ctx, tx, update.ProductID.String())
		if err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("product not found: %s", update.ProductID.String())
			}
			return fmt.Errorf("failed to find product: %w", err)
		}

		// Update stock quantity
		product.StockQuantity = null.NewInt(update.StockQuantity, true)
		product.UpdatedAt = null.NewTime(time.Now(), true)

		if _, err := product.Update(ctx, tx, boil.Infer()); err != nil {
			return fmt.Errorf("failed to update product stock for %s: %w", update.ProductID.String(), err)
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit batch update transaction: %w", err)
	}

	return nil
}

// GetInventoryHistory retrieves inventory adjustment history for a product
func (s *InventoryService) GetInventoryHistory(ctx context.Context, productID uuid.UUID, page, pageSize int) (*dto.InventoryHistoryResponse, error) {
	// In a real system, you would query from an inventory_logs table
	// For now, we'll return a mock response since we don't have this table
	return &dto.InventoryHistoryResponse{
		ProductID: productID,
		Logs:      []dto.InventoryAdjustmentLog{},
		Total:     0,
		Page:      page,
		PageSize:  pageSize,
	}, nil
}

// GetInventoryAlerts retrieves low stock alerts
func (s *InventoryService) GetInventoryAlerts(ctx context.Context, threshold int) (*dto.InventoryAlertsResponse, error) {
	products, err := models.Products(
		qm.Select("id", "name", "sku", "stock_quantity"),
		qm.Where("stock_quantity <= ?", threshold),
		qm.OrderBy("stock_quantity ASC"),
	).All(ctx, s.db)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve low stock alerts: %w", err)
	}

	alerts := make([]dto.InventoryStockAlert, len(products))
	for i, product := range products {
		status := "low_stock"
		if product.StockQuantity.Int == 0 {
			status = "out_of_stock"
		}

		alerts[i] = dto.InventoryStockAlert{
			ProductID:        uuid.MustParse(product.ID),
			ProductName:      product.Name,
			SKU:              product.Sku.String,
			CurrentStock:     product.StockQuantity.Int,
			MinimumThreshold: threshold,
			Status:           status,
		}
	}

	return &dto.InventoryAlertsResponse{
		Alerts: alerts,
		Total:  len(alerts),
	}, nil
}

// GetInventoryStats retrieves overall inventory statistics
func (s *InventoryService) GetInventoryStats(ctx context.Context) (*dto.InventoryStatsResponse, error) {
	// Total products
	totalProducts, err := models.Products().Count(ctx, s.db)
	if err != nil {
		return nil, fmt.Errorf("failed to count total products: %w", err)
	}

	// In stock products (> 0)
	inStockProducts, err := models.Products(qm.Where("stock_quantity > 0")).Count(ctx, s.db)
	if err != nil {
		return nil, fmt.Errorf("failed to count in-stock products: %w", err)
	}

	// Low stock products (1-10)
	lowStockProducts, err := models.Products(qm.Where("stock_quantity > 0 AND stock_quantity <= 10")).Count(ctx, s.db)
	if err != nil {
		return nil, fmt.Errorf("failed to count low-stock products: %w", err)
	}

	// Out of stock products (= 0)
	outOfStockProducts, err := models.Products(qm.Where("stock_quantity = 0")).Count(ctx, s.db)
	if err != nil {
		return nil, fmt.Errorf("failed to count out-of-stock products: %w", err)
	}

	// Calculate average stock level
	var avgStock float64
	if totalProducts > 0 {
		row := s.db.QueryRowContext(ctx, "SELECT AVG(COALESCE(stock_quantity, 0)) FROM products")
		if err := row.Scan(&avgStock); err != nil {
			return nil, fmt.Errorf("failed to calculate average stock level: %w", err)
		}
	}

	// Calculate total stock value (would need price information)
	// For now, we'll use a placeholder calculation
	var totalValue float64
	row := s.db.QueryRowContext(ctx, `
		SELECT COALESCE(SUM(COALESCE(stock_quantity, 0) * price), 0)
		FROM products
		WHERE stock_quantity > 0 AND price IS NOT NULL
	`)
	if err := row.Scan(&totalValue); err != nil {
		return nil, fmt.Errorf("failed to calculate total stock value: %w", err)
	}

	return &dto.InventoryStatsResponse{
		TotalProducts:      int(totalProducts),
		InStockProducts:    int(inStockProducts),
		LowStockProducts:   int(lowStockProducts),
		OutOfStockProducts: int(outOfStockProducts),
		TotalStockValue:    totalValue,
		AverageStockLevel:  avgStock,
		LastUpdated:        time.Now(),
	}, nil
}
