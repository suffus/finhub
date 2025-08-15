package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"finhub-backend/models"
)

type DealHandler struct {
	db *gorm.DB
}

type CreateDealRequest struct {
	Name              string   `json:"name" binding:"required"`
	Amount            *float64 `json:"amount"`
	Currency          string   `json:"currency"`
	Probability       int      `json:"probability"`
	PipelineID        string   `json:"pipelineId" binding:"required"`
	StageID           string   `json:"stageId" binding:"required"`
	ExpectedCloseDate *string  `json:"expectedCloseDate"`
	CompanyID         *string  `json:"companyId"`
	ContactID         *string  `json:"contactId"`
	AssignedUserID    *string  `json:"assignedUserId"`
}

type UpdateDealRequest struct {
	Name              *string  `json:"name"`
	Amount            *float64 `json:"amount"`
	Currency          *string  `json:"currency"`
	Probability       *int     `json:"probability"`
	PipelineID        *string  `json:"pipelineId"`
	StageID           *string  `json:"stageId"`
	ExpectedCloseDate *string  `json:"expectedCloseDate"`
	CompanyID         *string  `json:"companyId"`
	ContactID         *string  `json:"contactId"`
	AssignedUserID    *string  `json:"assignedUserId"`
}

func NewDealHandler(db *gorm.DB) *DealHandler {
	return &DealHandler{db: db}
}

func (h *DealHandler) GetDeals(c *gin.Context) {
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

	var deals []models.Deal
	if err := h.db.Where("tenant_id = ? AND is_deleted = ?", user.TenantID, false).
		Preload("Pipeline").Preload("Stage").Preload("Company").Preload("Contact").Preload("AssignedUser").
		Find(&deals).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch deals"})
		return
	}

	c.JSON(http.StatusOK, deals)
}

func (h *DealHandler) CreateDeal(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req CreateDealRequest
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
	deal := models.Deal{
		Name:           req.Name,
		Amount:         req.Amount,
		Currency:       req.Currency,
		Probability:    req.Probability,
		PipelineID:     req.PipelineID,
		StageID:        req.StageID,
		CompanyID:      req.CompanyID,
		ContactID:      req.ContactID,
		AssignedUserID: req.AssignedUserID,
		TenantID:       user.TenantID,
		CreatedBy:      &userIDStr,
	}

	if err := h.db.Create(&deal).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create deal"})
		return
	}

	c.JSON(http.StatusCreated, deal)
}

func (h *DealHandler) GetDeal(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	dealID := c.Param("id")

	var user models.User
	if err := h.db.Select("tenant_id").First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var deal models.Deal
	if err := h.db.Where("id = ? AND tenant_id = ? AND is_deleted = ?", dealID, user.TenantID, false).
		Preload("Pipeline").Preload("Stage").Preload("Company").Preload("Contact").Preload("AssignedUser").
		First(&deal).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Deal not found"})
		return
	}

	c.JSON(http.StatusOK, deal)
}

func (h *DealHandler) UpdateDeal(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	dealID := c.Param("id")
	var req UpdateDealRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := h.db.Select("tenant_id").First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var deal models.Deal
	if err := h.db.Where("id = ? AND tenant_id = ? AND is_deleted = ?", dealID, user.TenantID, false).
		First(&deal).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Deal not found"})
		return
	}

	// Update fields
	if req.Name != nil {
		deal.Name = *req.Name
	}
	if req.Amount != nil {
		deal.Amount = req.Amount
	}
	if req.Currency != nil {
		deal.Currency = *req.Currency
	}
	if req.Probability != nil {
		deal.Probability = *req.Probability
	}
	if req.PipelineID != nil {
		deal.PipelineID = *req.PipelineID
	}
	if req.StageID != nil {
		deal.StageID = *req.StageID
	}
	if req.CompanyID != nil {
		deal.CompanyID = req.CompanyID
	}
	if req.ContactID != nil {
		deal.ContactID = req.ContactID
	}
	if req.AssignedUserID != nil {
		deal.AssignedUserID = req.AssignedUserID
	}

	if err := h.db.Save(&deal).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update deal"})
		return
	}

	c.JSON(http.StatusOK, deal)
}

func (h *DealHandler) DeleteDeal(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	dealID := c.Param("id")

	var user models.User
	if err := h.db.Select("tenant_id").First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var deal models.Deal
	if err := h.db.Where("id = ? AND tenant_id = ? AND is_deleted = ?", dealID, user.TenantID, false).
		First(&deal).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Deal not found"})
		return
	}

	// Soft delete
	deal.IsDeleted = true
	if err := h.db.Save(&deal).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete deal"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deal deleted successfully"})
}
