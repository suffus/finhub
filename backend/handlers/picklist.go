package handlers

import (
	"net/http"

	"finhub-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PicklistHandler handles picklist-related requests
type PicklistHandler struct {
	db *gorm.DB
}

// NewPicklistHandler creates a new picklist handler
func NewPicklistHandler(db *gorm.DB) *PicklistHandler {
	return &PicklistHandler{db: db}
}

// PicklistSearchRequest defines the request body for searching picklists
type PicklistSearchRequest struct {
	Query      string `json:"query"`
	Limit      int    `json:"limit" binding:"required,min=1,max=100"`
	Offset     int    `json:"offset" binding:"min=0"`
	EntityType string `json:"entityType" binding:"required"`
}

// PicklistResponse defines the response for picklist queries
type PicklistResponse struct {
	Items      []gin.H `json:"items"`
	TotalCount int64   `json:"totalCount"`
	HasMore    bool    `json:"hasMore"`
}

// GetPicklistByEntity returns all active items for a specific entity type
func (h *PicklistHandler) GetPicklistByEntity(c *gin.Context) {
	entityType := c.Param("entity")

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get user from database
	var userObj models.User
	if err := h.db.Select("tenant_id").First(&userObj, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var response PicklistResponse

	switch entityType {
	case "industries":
		var industries []models.Industry
		if err := h.db.Where("tenant_id = ? AND is_active = ?", userObj.TenantID, true).
			Order("name ASC").Find(&industries).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch industries", "details": err.Error()})
			return
		}

		items := make([]gin.H, len(industries))
		for i, industry := range industries {
			items[i] = gin.H{
				"id":          industry.ID,
				"name":        industry.Name,
				"code":        industry.Code,
				"description": industry.Description,
			}
		}
		response.Items = items
		response.TotalCount = int64(len(items))
		response.HasMore = false

	case "companysizes":
		var sizes []models.CompanySize
		if err := h.db.Where("tenant_id = ? AND is_active = ?", userObj.TenantID, true).
			Order("name ASC").Find(&sizes).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch company sizes", "details": err.Error()})
			return
		}

		items := make([]gin.H, len(sizes))
		for i, size := range sizes {
			items[i] = gin.H{
				"id":           size.ID,
				"name":         size.Name,
				"code":         size.Code,
				"description":  size.Description,
				"minEmployees": size.MinEmployees,
				"maxEmployees": size.MaxEmployees,
			}
		}
		response.Items = items
		response.TotalCount = int64(len(items))
		response.HasMore = false

	case "leadstatuses":
		var statuses []models.LeadStatus
		if err := h.db.Where("tenant_id = ? AND is_active = ?", userObj.TenantID, true).
			Order("order ASC").Find(&statuses).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch lead statuses", "details": err.Error()})
			return
		}

		items := make([]gin.H, len(statuses))
		for i, status := range statuses {
			items[i] = gin.H{
				"id":          status.ID,
				"name":        status.Name,
				"code":        status.Code,
				"description": status.Description,
				"color":       status.Color,
				"order":       status.Order,
			}
		}
		response.Items = items
		response.TotalCount = int64(len(items))
		response.HasMore = false

	case "leadtemperatures":
		var temperatures []models.LeadTemperature
		if err := h.db.Where("tenant_id = ? AND is_active = ?", userObj.TenantID, true).
			Order("order ASC").Find(&temperatures).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch lead temperatures", "details": err.Error()})
			return
		}

		items := make([]gin.H, len(temperatures))
		for i, temp := range temperatures {
			items[i] = gin.H{
				"id":          temp.ID,
				"name":        temp.Name,
				"code":        temp.Code,
				"description": temp.Description,
				"color":       temp.Color,
				"order":       temp.Order,
			}
		}
		response.Items = items
		response.TotalCount = int64(len(items))
		response.HasMore = false

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entity type"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// SearchPicklist searches for items in a specific picklist entity
func (h *PicklistHandler) SearchPicklist(c *gin.Context) {
	var req PicklistSearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get user from database
	var userObj models.User
	if err := h.db.Select("tenant_id").First(&userObj, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var response PicklistResponse
	var totalCount int64

	switch req.EntityType {
	case "industry":
		var industries []models.Industry
		baseQuery := h.db.Model(&models.Industry{}).Where("tenant_id = ? AND is_active = ?", userObj.TenantID, true)

		// Count total records first
		if err := baseQuery.Count(&totalCount).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count industries", "details": err.Error()})
			return
		}

		// Build search query
		searchQuery := baseQuery
		if req.Query != "" {
			searchQuery = searchQuery.Where("LOWER(name) LIKE LOWER(?) OR LOWER(code) LIKE LOWER(?)", "%"+req.Query+"%", "%"+req.Query+"%")
		}

		// Apply pagination and ordering
		if err := searchQuery.Order("name ASC").Limit(req.Limit).Offset(req.Offset).Find(&industries).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch industries", "details": err.Error()})
			return
		}

		items := make([]gin.H, len(industries))
		for i, industry := range industries {
			items[i] = gin.H{
				"id":          industry.ID,
				"name":        industry.Name,
				"code":        industry.Code,
				"description": industry.Description,
			}
		}
		response.Items = items
		response.TotalCount = totalCount
		response.HasMore = int64(req.Offset+len(items)) < totalCount

	case "companysize":
		var sizes []models.CompanySize
		baseQuery := h.db.Model(&models.CompanySize{}).Where("tenant_id = ? AND is_active = ?", userObj.TenantID, true)

		// Count total records first
		if err := baseQuery.Count(&totalCount).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count company sizes", "details": err.Error()})
			return
		}

		// Build search query
		searchQuery := baseQuery
		if req.Query != "" {
			searchQuery = searchQuery.Where("LOWER(name) LIKE LOWER(?) OR LOWER(code) LIKE LOWER(?)", "%"+req.Query+"%", "%"+req.Query+"%")
		}

		// Apply pagination and ordering
		if err := searchQuery.Order("name ASC").Limit(req.Limit).Offset(req.Offset).Find(&sizes).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch company sizes", "details": err.Error()})
			return
		}

		items := make([]gin.H, len(sizes))
		for i, size := range sizes {
			items[i] = gin.H{
				"id":           size.ID,
				"name":         size.Name,
				"code":         size.Code,
				"description":  size.Description,
				"minEmployees": size.MinEmployees,
				"maxEmployees": size.MaxEmployees,
			}
		}
		response.Items = items
		response.TotalCount = totalCount
		response.HasMore = int64(req.Offset+len(items)) < totalCount

	case "leadstatus":
		var statuses []models.LeadStatus
		baseQuery := h.db.Model(&models.LeadStatus{}).Where("tenant_id = ? AND is_active = ?", userObj.TenantID, true)

		// Count total records first
		if err := baseQuery.Count(&totalCount).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count lead statuses", "details": err.Error()})
			return
		}

		// Build search query
		searchQuery := baseQuery
		if req.Query != "" {
			searchQuery = searchQuery.Where("LOWER(name) LIKE LOWER(?) OR LOWER(code) LIKE LOWER(?)", "%"+req.Query+"%", "%"+req.Query+"%")
		}

		// Apply pagination and ordering
		if err := searchQuery.Order("order ASC").Limit(req.Limit).Offset(req.Offset).Find(&statuses).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch lead statuses", "details": err.Error()})
			return
		}

		items := make([]gin.H, len(statuses))
		for i, status := range statuses {
			items[i] = gin.H{
				"id":          status.ID,
				"name":        status.Name,
				"code":        status.Code,
				"description": status.Description,
				"color":       status.Color,
				"order":       status.Order,
			}
		}
		response.Items = items
		response.TotalCount = totalCount
		response.HasMore = int64(req.Offset+len(items)) < totalCount

	case "leadtemperature":
		var temperatures []models.LeadTemperature
		baseQuery := h.db.Model(&models.LeadTemperature{}).Where("tenant_id = ? AND is_active = ?", userObj.TenantID, true)

		// Count total records first
		if err := baseQuery.Count(&totalCount).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count lead temperatures", "details": err.Error()})
			return
		}

		// Build search query
		searchQuery := baseQuery
		if req.Query != "" {
			searchQuery = searchQuery.Where("LOWER(name) LIKE LOWER(?) OR LOWER(code) LIKE LOWER(?)", "%"+req.Query+"%", "%"+req.Query+"%")
		}

		// Apply pagination and ordering
		if err := searchQuery.Order("order ASC").Limit(req.Limit).Offset(req.Offset).Find(&temperatures).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch lead temperatures", "details": err.Error()})
			return
		}

		items := make([]gin.H, len(temperatures))
		for i, temp := range temperatures {
			items[i] = gin.H{
				"id":          temp.ID,
				"name":        temp.Name,
				"code":        temp.Code,
				"description": temp.Description,
				"color":       temp.Color,
				"order":       temp.Order,
			}
		}
		response.Items = items
		response.TotalCount = totalCount
		response.HasMore = int64(req.Offset+len(items)) < totalCount

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entity type"})
		return
	}

	c.JSON(http.StatusOK, response)
}
