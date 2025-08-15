package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"finhub-backend/models"
)

type LeadHandler struct {
	db *gorm.DB
}

type CreateLeadRequest struct {
	FirstName      *string `json:"firstName"`
	LastName       *string `json:"lastName"`
	Title          *string `json:"title"`
	StatusID       *string `json:"statusId"`
	TemperatureID  *string `json:"temperatureId"`
	Source         *string `json:"source"`
	Campaign       *string `json:"campaign"`
	Score          int     `json:"score"`
	CompanyID      *string `json:"companyId"`
	AssignedUserID *string `json:"assignedUserId"`
}

type UpdateLeadRequest struct {
	FirstName      *string `json:"firstName"`
	LastName       *string `json:"lastName"`
	Title          *string `json:"title"`
	StatusID       *string `json:"statusId"`
	TemperatureID  *string `json:"temperatureId"`
	Source         *string `json:"source"`
	Campaign       *string `json:"campaign"`
	Score          *int    `json:"score"`
	CompanyID      *string `json:"companyId"`
	AssignedUserID *string `json:"assignedUserId"`
}

func NewLeadHandler(db *gorm.DB) *LeadHandler {
	return &LeadHandler{db: db}
}

func (h *LeadHandler) GetLeads(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var user models.User
	if err := h.db.Select("tenant_id").First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var leads []models.Lead
	if err := h.db.Where("tenant_id = ? AND is_deleted = ?", user.TenantID, false).
		Preload("Status").Preload("Temperature").Preload("Company").Preload("AssignedUser").
		Find(&leads).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch leads"})
		return
	}

	c.JSON(http.StatusOK, leads)
}

func (h *LeadHandler) CreateLead(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req CreateLeadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := h.db.Select("tenant_id").First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	userIDStr := userID.(string)
	lead := models.Lead{
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Title:          req.Title,
		StatusID:       req.StatusID,
		TemperatureID:  req.TemperatureID,
		Source:         req.Source,
		Campaign:       req.Campaign,
		Score:          req.Score,
		CompanyID:      req.CompanyID,
		AssignedUserID: req.AssignedUserID,
		TenantID:       user.TenantID,
		CreatedBy:      &userIDStr,
	}

	if err := h.db.Create(&lead).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create lead"})
		return
	}

	c.JSON(http.StatusCreated, lead)
}

func (h *LeadHandler) GetLead(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	leadID := c.Param("id")

	var user models.User
	if err := h.db.Select("tenant_id").First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var lead models.Lead
	if err := h.db.Where("id = ? AND tenant_id = ? AND is_deleted = ?", leadID, user.TenantID, false).
		Preload("Status").Preload("Temperature").Preload("Company").Preload("AssignedUser").
		First(&lead).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Lead not found"})
		return
	}

	c.JSON(http.StatusOK, lead)
}

func (h *LeadHandler) UpdateLead(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	leadID := c.Param("id")
	var req UpdateLeadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := h.db.Select("tenant_id").First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var lead models.Lead
	if err := h.db.Where("id = ? AND tenant_id = ? AND is_deleted = ?", leadID, user.TenantID, false).
		First(&lead).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Lead not found"})
		return
	}

	// Update fields
	if req.FirstName != nil {
		lead.FirstName = req.FirstName
	}
	if req.LastName != nil {
		lead.LastName = req.LastName
	}
	if req.Title != nil {
		lead.Title = req.Title
	}
	if req.StatusID != nil {
		lead.StatusID = req.StatusID
	}
	if req.TemperatureID != nil {
		lead.TemperatureID = req.TemperatureID
	}
	if req.Source != nil {
		lead.Source = req.Source
	}
	if req.Campaign != nil {
		lead.Campaign = req.Campaign
	}
	if req.Score != nil {
		lead.Score = *req.Score
	}
	if req.CompanyID != nil {
		lead.CompanyID = req.CompanyID
	}
	if req.AssignedUserID != nil {
		lead.AssignedUserID = req.AssignedUserID
	}

	if err := h.db.Save(&lead).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update lead"})
		return
	}

	c.JSON(http.StatusOK, lead)
}

func (h *LeadHandler) DeleteLead(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	leadID := c.Param("id")

	var user models.User
	if err := h.db.Select("tenant_id").First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var lead models.Lead
	if err := h.db.Where("id = ? AND tenant_id = ? AND is_deleted = ?", leadID, user.TenantID, false).
		First(&lead).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Lead not found"})
		return
	}

	// Soft delete
	lead.IsDeleted = true
	if err := h.db.Save(&lead).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete lead"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Lead deleted successfully"})
}
