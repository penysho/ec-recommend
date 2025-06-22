package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"ec-recommend/internal/dto"

	"github.com/gin-gonic/gin"
)

// UserRole represents user roles for authorization
type UserRole string

const (
	RoleAdmin     UserRole = "admin"
	RoleEmployee  UserRole = "employee"
	RoleCustomer  UserRole = "customer"
)

// UserClaims represents the authenticated user information
type UserClaims struct {
	UserID   string   `json:"user_id"`
	Email    string   `json:"email"`
	Role     UserRole `json:"role"`
	CustomerID string `json:"customer_id,omitempty"` // Only for customer role
}

// RequireAuth returns a middleware that requires authentication
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error:     "Authorization header is required",
				Code:      http.StatusUnauthorized,
				Timestamp: time.Now(),
			})
			c.Abort()
			return
		}

		// Check Bearer token format
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error:     "Invalid authorization header format",
				Code:      http.StatusUnauthorized,
				Timestamp: time.Now(),
			})
			c.Abort()
			return
		}

		token := tokenParts[1]

		// For demonstration purposes, we'll use a simple token validation
		// In production, you would validate JWT tokens or use proper authentication service
		userClaims, err := validateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error:     "Invalid or expired token: " + err.Error(),
				Code:      http.StatusUnauthorized,
				Timestamp: time.Now(),
			})
			c.Abort()
			return
		}

		// Store user claims in context for use by handlers
		c.Set("user_claims", userClaims)
		c.Next()
	}
}

// RequireRole returns a middleware that requires specific roles
func RequireRole(allowedRoles ...UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		// This middleware should be used after RequireAuth
		userClaimsInterface, exists := c.Get("user_claims")
		if !exists {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error:     "User not authenticated",
				Code:      http.StatusUnauthorized,
				Timestamp: time.Now(),
			})
			c.Abort()
			return
		}

		userClaims, ok := userClaimsInterface.(*UserClaims)
		if !ok {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:     "Invalid user claims",
				Code:      http.StatusInternalServerError,
				Timestamp: time.Now(),
			})
			c.Abort()
			return
		}

		// Check if user has required role
		hasValidRole := false
		for _, allowedRole := range allowedRoles {
			if userClaims.Role == allowedRole {
				hasValidRole = true
				break
			}
		}

		if !hasValidRole {
			c.JSON(http.StatusForbidden, dto.ErrorResponse{
				Error:     "Insufficient permissions",
				Code:      http.StatusForbidden,
				Timestamp: time.Now(),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// validateToken validates the provided token and returns user claims
// This is a simplified implementation for demonstration
func validateToken(token string) (*UserClaims, error) {
	// In production, you would:
	// 1. Verify JWT signature
	// 2. Check token expiration
	// 3. Validate against user database
	// 4. Handle token refresh logic

	// For demonstration, we'll use predefined tokens
	switch token {
	case "admin-token-123":
		return &UserClaims{
			UserID: "admin-001",
			Email:  "admin@example.com",
			Role:   RoleAdmin,
		}, nil
	case "employee-token-456":
		return &UserClaims{
			UserID: "emp-001",
			Email:  "employee@example.com",
			Role:   RoleEmployee,
		}, nil
	case "customer-token-789":
		return &UserClaims{
			UserID:     "cust-001",
			Email:      "customer@example.com",
			Role:       RoleCustomer,
			CustomerID: "550e8400-e29b-41d4-a716-446655440000", // Example customer UUID
		}, nil
	default:
		// In production, parse and validate JWT token here
		return nil, fmt.Errorf("invalid token")
	}
}

// GetUserClaims is a helper function to get user claims from gin context
func GetUserClaims(c *gin.Context) (*UserClaims, bool) {
	userClaimsInterface, exists := c.Get("user_claims")
	if !exists {
		return nil, false
	}

	userClaims, ok := userClaimsInterface.(*UserClaims)
	return userClaims, ok
}
