package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"finhub-backend/models"
)

type CompanyHandler struct {
	db *gorm.DB
}

type CreateCompanyRequest struct {
	Name       string   `json:"name" binding:"required"`
	Website    *string  `json:"website"`
	Domain     *string  `json:"domain"`
	IndustryID *string  `json:"industryId"`
	SizeID     *string  `json:"sizeId"`
	Revenue    *float64 `json:"revenue"`
}

type UpdateCompanyRequest struct {
	Name       *string  `json:"name"`
	Website    *string  `json:"website"`
	Domain     *string  `json:"domain"`
	IndustryID *string  `json:"industryId"`
	SizeID     *string  `json:"sizeId"`
	Revenue    *float64 `json:"revenue"`
}

func NewCompanyHandler(db *gorm.DB) *CompanyHandler {
	return &CompanyHandler{db: db}
}

func (h *CompanyHandler) GetCompanies(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get user's tenant
	var user models.User
	if err := h.db.Select("tenant_id").First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var companies []models.Company
	if err := h.db.Where("tenant_id = ? AND is_deleted = ?", user.TenantID, false).
		Preload("Industry").Preload("Size").
		Find(&companies).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch companies"})
		return
	}

	c.JSON(http.StatusOK, companies)
}

func (h *CompanyHandler) CreateCompany(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req CreateCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user's tenant
	var user models.User
	if err := h.db.Select("tenant_id").First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	userIDStr := userID.(string)
	company := models.Company{
		Name:       req.Name,
		Website:    req.Website,
		Domain:     req.Domain,
		IndustryID: req.IndustryID,
		SizeID:     req.SizeID,
		Revenue:    req.Revenue,
		TenantID:   user.TenantID,
		CreatedBy:  &userIDStr,
	}

	if err := h.db.Create(&company).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create company"})
		return
	}

	c.JSON(http.StatusCreated, company)
}

func (h *CompanyHandler) GetCompany(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	companyID := c.Param("id")

	// Get user's tenant
	var user models.User
	if err := h.db.Select("tenant_id").First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var company models.Company
	if err := h.db.Where("id = ? AND tenant_id = ? AND is_deleted = ?", companyID, user.TenantID, false).
		Preload("Industry").Preload("Size").
		First(&company).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
		return
	}

	c.JSON(http.StatusOK, company)
}

func (h *CompanyHandler) UpdateCompany(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	companyID := c.Param("id")
	var req UpdateCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user's tenant
	var user models.User
	if err := h.db.Select("tenant_id").First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var company models.Company
	if err := h.db.Where("id = ? AND tenant_id = ? AND is_deleted = ?", companyID, user.TenantID, false).
		First(&company).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
		return
	}

	// Update fields
	if req.Name != nil {
		company.Name = *req.Name
	}
	if req.Website != nil {
		company.Website = req.Website
	}
	if req.Domain != nil {
		company.Domain = req.Domain
	}
	if req.IndustryID != nil {
		company.IndustryID = req.IndustryID
	}
	if req.SizeID != nil {
		company.SizeID = req.SizeID
	}
	if req.Revenue != nil {
		company.Revenue = req.Revenue
	}

	if err := h.db.Save(&company).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update company"})
		return
	}

	c.JSON(http.StatusOK, company)
}

func (h *CompanyHandler) DeleteCompany(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	companyID := c.Param("id")

	// Get user's tenant
	var user models.User
	if err := h.db.Select("tenant_id").First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var company models.Company
	if err := h.db.Where("id = ? AND tenant_id = ? AND is_deleted = ?", companyID, user.TenantID, false).
		First(&company).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
		return
	}

	// Soft delete
	company.IsDeleted = true
	if err := h.db.Save(&company).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete company"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Company deleted successfully"})
}
