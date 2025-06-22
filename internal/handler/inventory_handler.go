package handler

import (
	"net/http"
	"strconv"
	"time"

	"ec-recommend/internal/dto"
	"ec-recommend/internal/middleware"
	"ec-recommend/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// InventoryHandler handles inventory-related HTTP requests
type InventoryHandler struct {
	inventoryService service.InventoryServiceInterface
}

// NewInventoryHandler creates a new inventory handler instance
func NewInventoryHandler(inventoryService service.InventoryServiceInterface) *InventoryHandler {
	return &InventoryHandler{
		inventoryService: inventoryService,
	}
}

// GetInventoryList handles GET /api/v1/inventory
// Retrieves a paginated list of inventory items with search and filtering capabilities
func (h *InventoryHandler) GetInventoryList(c *gin.Context) {
	// Build search request from query parameters
	req := &dto.InventorySearchRequest{
		Query:         c.Query("query"),
		SKU:           c.Query("sku"),
		StockStatus:   c.Query("stock_status"),
		SortBy:        c.Query("sort_by"),
		SortDirection: c.Query("sort_direction"),
		Page:          1,
		PageSize:      20,
	}

	// Parse category_id
	if categoryIDStr := c.Query("category_id"); categoryIDStr != "" {
		categoryID, err := strconv.Atoi(categoryIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:     "Invalid category_id format",
				Code:      http.StatusBadRequest,
				Timestamp: time.Now(),
			})
			return
		}
		req.CategoryID = &categoryID
	}

	// Parse min_stock
	if minStockStr := c.Query("min_stock"); minStockStr != "" {
		minStock, err := strconv.Atoi(minStockStr)
		if err != nil || minStock < 0 {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:     "min_stock must be a non-negative integer",
				Code:      http.StatusBadRequest,
				Timestamp: time.Now(),
			})
			return
		}
		req.MinStock = &minStock
	}

	// Parse max_stock
	if maxStockStr := c.Query("max_stock"); maxStockStr != "" {
		maxStock, err := strconv.Atoi(maxStockStr)
		if err != nil || maxStock < 0 {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:     "max_stock must be a non-negative integer",
				Code:      http.StatusBadRequest,
				Timestamp: time.Now(),
			})
			return
		}
		req.MaxStock = &maxStock
	}

	// Parse page
	if pageStr := c.Query("page"); pageStr != "" {
		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:     "page must be a positive integer",
				Code:      http.StatusBadRequest,
				Timestamp: time.Now(),
			})
			return
		}
		req.Page = page
	}

	// Parse page_size
	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		pageSize, err := strconv.Atoi(pageSizeStr)
		if err != nil || pageSize < 1 || pageSize > 100 {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:     "page_size must be between 1 and 100",
				Code:      http.StatusBadRequest,
				Timestamp: time.Now(),
			})
			return
		}
		req.PageSize = pageSize
	}

	// Validate min_stock and max_stock relationship
	if req.MinStock != nil && req.MaxStock != nil && *req.MinStock > *req.MaxStock {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:     "min_stock cannot be greater than max_stock",
			Code:      http.StatusBadRequest,
			Timestamp: time.Now(),
		})
		return
	}

	// Get inventory list
	response, err := h.inventoryService.GetInventoryList(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:     "Failed to retrieve inventory list: " + err.Error(),
			Code:      http.StatusInternalServerError,
			Timestamp: time.Now(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetInventoryByProductID handles GET /api/v1/inventory/products/:product_id
// Retrieves inventory information for a specific product
func (h *InventoryHandler) GetInventoryByProductID(c *gin.Context) {
	productIDStr := c.Param("product_id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:     "Invalid product_id format",
			Code:      http.StatusBadRequest,
			Timestamp: time.Now(),
		})
		return
	}

	inventory, err := h.inventoryService.GetInventoryByProductID(c.Request.Context(), productID)
	if err != nil {
		if err.Error() == "product not found: "+productID.String() {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:     "Product not found",
				Code:      http.StatusNotFound,
				Timestamp: time.Now(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:     "Failed to retrieve inventory: " + err.Error(),
			Code:      http.StatusInternalServerError,
			Timestamp: time.Now(),
		})
		return
	}

	c.JSON(http.StatusOK, inventory)
}

// UpdateInventory handles PUT /api/v1/inventory/products/:product_id
// Updates inventory for a specific product (Admin/Employee only)
func (h *InventoryHandler) UpdateInventory(c *gin.Context) {
	// Get user claims for logging purposes
	userClaims, exists := middleware.GetUserClaims(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:     "User not authenticated",
			Code:      http.StatusUnauthorized,
			Timestamp: time.Now(),
		})
		return
	}

	productIDStr := c.Param("product_id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:     "Invalid product_id format",
			Code:      http.StatusBadRequest,
			Timestamp: time.Now(),
		})
		return
	}

	var req dto.InventoryUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:     "Invalid request body: " + err.Error(),
			Code:      http.StatusBadRequest,
			Timestamp: time.Now(),
		})
		return
	}

	// Set product ID from URL parameter
	req.ProductID = productID

	// Update inventory
	if err := h.inventoryService.UpdateInventory(c.Request.Context(), &req, userClaims.UserID); err != nil {
		if err.Error() == "product not found: "+productID.String() {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:     "Product not found",
				Code:      http.StatusNotFound,
				Timestamp: time.Now(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:     "Failed to update inventory: " + err.Error(),
			Code:      http.StatusInternalServerError,
			Timestamp: time.Now(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Inventory updated successfully",
		"timestamp": time.Now(),
	})
}

// PerformInventoryTransaction handles POST /api/v1/inventory/transactions
// Performs inventory transactions (increase, decrease, set) (Admin/Employee only)
func (h *InventoryHandler) PerformInventoryTransaction(c *gin.Context) {
	// Get user claims for logging purposes
	userClaims, exists := middleware.GetUserClaims(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:     "User not authenticated",
			Code:      http.StatusUnauthorized,
			Timestamp: time.Now(),
		})
		return
	}

	var req dto.InventoryTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:     "Invalid request body: " + err.Error(),
			Code:      http.StatusBadRequest,
			Timestamp: time.Now(),
		})
		return
	}

	// Perform transaction
	if err := h.inventoryService.PerformInventoryTransaction(c.Request.Context(), &req, userClaims.UserID); err != nil {
		if err.Error() == "product not found: "+req.ProductID.String() {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:     "Product not found",
				Code:      http.StatusNotFound,
				Timestamp: time.Now(),
			})
			return
		}
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:     err.Error(),
			Code:      http.StatusBadRequest,
			Timestamp: time.Now(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Inventory transaction completed successfully",
		"timestamp": time.Now(),
	})
}

// BatchUpdateInventory handles POST /api/v1/inventory/batch-update
// Updates inventory for multiple products (Admin only)
func (h *InventoryHandler) BatchUpdateInventory(c *gin.Context) {
	// Get user claims for logging purposes
	userClaims, exists := middleware.GetUserClaims(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:     "User not authenticated",
			Code:      http.StatusUnauthorized,
			Timestamp: time.Now(),
		})
		return
	}

	var req dto.InventoryBatchUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:     "Invalid request body: " + err.Error(),
			Code:      http.StatusBadRequest,
			Timestamp: time.Now(),
		})
		return
	}

	// Validate that we have updates
	if len(req.Updates) == 0 {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:     "At least one update is required",
			Code:      http.StatusBadRequest,
			Timestamp: time.Now(),
		})
		return
	}

	// Batch update
	if err := h.inventoryService.BatchUpdateInventory(c.Request.Context(), &req, userClaims.UserID); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:     "Failed to perform batch update: " + err.Error(),
			Code:      http.StatusInternalServerError,
			Timestamp: time.Now(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Batch inventory update completed successfully",
		"updated_count": len(req.Updates),
		"timestamp":     time.Now(),
	})
}

// GetInventoryHistory handles GET /api/v1/inventory/products/:product_id/history
// Retrieves inventory adjustment history for a specific product
func (h *InventoryHandler) GetInventoryHistory(c *gin.Context) {
	productIDStr := c.Param("product_id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:     "Invalid product_id format",
			Code:      http.StatusBadRequest,
			Timestamp: time.Now(),
		})
		return
	}

	// Parse pagination parameters
	page := 1
	pageSize := 20

	if pageStr := c.Query("page"); pageStr != "" {
		parsedPage, err := strconv.Atoi(pageStr)
		if err != nil || parsedPage < 1 {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:     "page must be a positive integer",
				Code:      http.StatusBadRequest,
				Timestamp: time.Now(),
			})
			return
		}
		page = parsedPage
	}

	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		parsedPageSize, err := strconv.Atoi(pageSizeStr)
		if err != nil || parsedPageSize < 1 || parsedPageSize > 100 {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:     "page_size must be between 1 and 100",
				Code:      http.StatusBadRequest,
				Timestamp: time.Now(),
			})
			return
		}
		pageSize = parsedPageSize
	}

	// Get history
	response, err := h.inventoryService.GetInventoryHistory(c.Request.Context(), productID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:     "Failed to retrieve inventory history: " + err.Error(),
			Code:      http.StatusInternalServerError,
			Timestamp: time.Now(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetInventoryAlerts handles GET /api/v1/inventory/alerts
// Retrieves low stock alerts (Admin/Employee only)
func (h *InventoryHandler) GetInventoryAlerts(c *gin.Context) {
	// Parse threshold parameter
	threshold := 10 // Default threshold
	if thresholdStr := c.Query("threshold"); thresholdStr != "" {
		parsedThreshold, err := strconv.Atoi(thresholdStr)
		if err != nil || parsedThreshold < 0 {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:     "threshold must be a non-negative integer",
				Code:      http.StatusBadRequest,
				Timestamp: time.Now(),
			})
			return
		}
		threshold = parsedThreshold
	}

	// Get alerts
	response, err := h.inventoryService.GetInventoryAlerts(c.Request.Context(), threshold)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:     "Failed to retrieve inventory alerts: " + err.Error(),
			Code:      http.StatusInternalServerError,
			Timestamp: time.Now(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetInventoryStats handles GET /api/v1/inventory/stats
// Retrieves overall inventory statistics (Admin/Employee only)
func (h *InventoryHandler) GetInventoryStats(c *gin.Context) {
	// Get stats
	response, err := h.inventoryService.GetInventoryStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:     "Failed to retrieve inventory statistics: " + err.Error(),
			Code:      http.StatusInternalServerError,
			Timestamp: time.Now(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}
