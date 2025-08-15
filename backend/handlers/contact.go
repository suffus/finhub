package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"finhub-backend/models"
)

type ContactHandler struct {
	db *gorm.DB
}

type CreateContactRequest struct {
	FirstName      string  `json:"firstName" binding:"required"`
	LastName       string  `json:"lastName" binding:"required"`
	Title          *string `json:"title"`
	Department     *string `json:"department"`
	CompanyID      *string `json:"companyId"`
	OriginalSource *string `json:"originalSource"`
	EmailOptIn     bool    `json:"emailOptIn"`
	SmsOptIn       bool    `json:"smsOptIn"`
	CallOptIn      bool    `json:"callOptIn"`
}

type UpdateContactRequest struct {
	FirstName      *string `json:"firstName"`
	LastName       *string `json:"lastName"`
	Title          *string `json:"title"`
	Department     *string `json:"department"`
	CompanyID      *string `json:"companyId"`
	OriginalSource *string `json:"originalSource"`
	EmailOptIn     *bool   `json:"emailOptIn"`
	SmsOptIn       *bool   `json:"smsOptIn"`
	CallOptIn      *bool   `json:"callOptIn"`
}

func NewContactHandler(db *gorm.DB) *ContactHandler {
	return &ContactHandler{db: db}
}

func (h *ContactHandler) GetContacts(c *gin.Context) {
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

	var contacts []models.Contact
	if err := h.db.Where("tenant_id = ? AND is_deleted = ?", user.TenantID, false).
		Preload("Company").
		Find(&contacts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch contacts"})
		return
	}

	c.JSON(http.StatusOK, contacts)
}

func (h *ContactHandler) CreateContact(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req CreateContactRequest
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
	contact := models.Contact{
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Title:          req.Title,
		Department:     req.Department,
		CompanyID:      req.CompanyID,
		OriginalSource: req.OriginalSource,
		EmailOptIn:     req.EmailOptIn,
		SmsOptIn:       req.SmsOptIn,
		CallOptIn:      req.CallOptIn,
		TenantID:       user.TenantID,
		CreatedBy:      &userIDStr,
	}

	if err := h.db.Create(&contact).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create contact"})
		return
	}

	c.JSON(http.StatusCreated, contact)
}

func (h *ContactHandler) GetContact(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	contactID := c.Param("id")

	var user models.User
	if err := h.db.Select("tenant_id").First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var contact models.Contact
	if err := h.db.Where("id = ? AND tenant_id = ? AND is_deleted = ?", contactID, user.TenantID, false).
		Preload("Company").
		First(&contact).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contact not found"})
		return
	}

	c.JSON(http.StatusOK, contact)
}

func (h *ContactHandler) UpdateContact(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	contactID := c.Param("id")
	var req UpdateContactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := h.db.Select("tenant_id").First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var contact models.Contact
	if err := h.db.Where("id = ? AND tenant_id = ? AND is_deleted = ?", contactID, user.TenantID, false).
		First(&contact).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contact not found"})
		return
	}

	// Update fields
	if req.FirstName != nil {
		contact.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		contact.LastName = *req.LastName
	}
	if req.Title != nil {
		contact.Title = req.Title
	}
	if req.Department != nil {
		contact.Department = req.Department
	}
	if req.CompanyID != nil {
		contact.CompanyID = req.CompanyID
	}
	if req.OriginalSource != nil {
		contact.OriginalSource = req.OriginalSource
	}
	if req.EmailOptIn != nil {
		contact.EmailOptIn = *req.EmailOptIn
	}
	if req.SmsOptIn != nil {
		contact.SmsOptIn = *req.SmsOptIn
	}
	if req.CallOptIn != nil {
		contact.CallOptIn = *req.CallOptIn
	}

	if err := h.db.Save(&contact).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update contact"})
		return
	}

	c.JSON(http.StatusOK, contact)
}

func (h *ContactHandler) DeleteContact(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	contactID := c.Param("id")

	var user models.User
	if err := h.db.Select("tenant_id").First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var contact models.Contact
	if err := h.db.Where("id = ? AND tenant_id = ? AND is_deleted = ?", contactID, user.TenantID, false).
		First(&contact).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contact not found"})
		return
	}

	// Soft delete
	contact.IsDeleted = true
	if err := h.db.Save(&contact).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete contact"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Contact deleted successfully"})
}
