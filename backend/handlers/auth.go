package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"finhub-backend/config"
	"finhub-backend/models"
)

// Helper function to convert string to *string
func stringPtr(s string) *string {
	return &s
}

type AuthHandler struct {
	db     *gorm.DB
	config *config.Config
}

type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	TenantID  string `json:"tenantId"` // Make this optional
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

func NewAuthHandler(db *gorm.DB, config *config.Config) *AuthHandler {
	return &AuthHandler{
		db:     db,
		config: config,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user already exists
	var existingUser models.User
	if err := h.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	// Determine tenant ID - create default tenant if none provided
	var tenantID string
	var defaultRoleID string
	if req.TenantID != "" {
		tenantID = req.TenantID
	} else {
		// Check if any tenant exists
		var existingTenant models.Tenant
		if err := h.db.First(&existingTenant).Error; err != nil {
			// No tenant exists, create a default one
			defaultTenant := models.Tenant{
				Name:      "Default Organization",
				Subdomain: "default",
				IsActive:  true,
				Settings:  map[string]interface{}{},
			}

			if err := h.db.Create(&defaultTenant).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create default tenant"})
				return
			}
			tenantID = defaultTenant.ID

			// Create a default admin role for the new tenant
			defaultRole := models.UserRole{
				Name:        "Administrator",
				Code:        "ADMIN",
				Description: stringPtr("Full system access"),
				IsActive:    true,
				IsSystem:    true,
				TenantID:    tenantID,
				Permissions: map[string]interface{}{
					"users":     []string{"create", "read", "update", "delete"},
					"companies": []string{"create", "read", "update", "delete"},
					"contacts":  []string{"create", "read", "update", "delete"},
					"leads":     []string{"create", "read", "update", "delete"},
					"deals":     []string{"create", "read", "update", "delete"},
				},
			}

			if err := h.db.Create(&defaultRole).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create default user role"})
				return
			}
			defaultRoleID = defaultRole.ID
		} else {
			tenantID = existingTenant.ID
			// Find existing admin role for this tenant
			var existingRole models.UserRole
			if err := h.db.Where("tenant_id = ? AND code = ?", tenantID, "ADMIN").First(&existingRole).Error; err == nil {
				defaultRoleID = existingRole.ID
			}
		}
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create user
	user := models.User{
		Email:     req.Email,
		Password:  string(hashedPassword),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		TenantID:  tenantID,
		RoleID:    stringPtr(defaultRoleID),
		IsActive:  true,
	}

	if err := h.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Generate JWT token
	token, err := h.generateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Hide password from response
	user.Password = ""

	c.JSON(http.StatusCreated, AuthResponse{
		Token: token,
		User:  user,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user
	var user models.User
	if err := h.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Check if user is active
	if !user.IsActive {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User account is deactivated"})
		return
	}

	// Generate JWT token
	token, err := h.generateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Hide password from response
	user.Password = ""

	c.JSON(http.StatusOK, AuthResponse{
		Token: token,
		User:  user,
	})
}

func (h *AuthHandler) generateJWT(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days
		"iat":     time.Now().Unix(),
	})

	return token.SignedString([]byte(h.config.JWTSecret))
}
