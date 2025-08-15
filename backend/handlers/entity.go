package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"finhub-backend/models"
)

type EntityHandler struct {
	db *gorm.DB
}

type EntityQueryRequest struct {
	EntityType string                 `json:"entityType" binding:"required"`
	Page       int                    `json:"page" binding:"min=1"`
	PageSize   int                    `json:"pageSize" binding:"min=1,max=100"`
	SortBy     string                 `json:"sortBy"`
	SortOrder  string                 `json:"sortOrder"` // "asc" or "desc"
	Filters    map[string]interface{} `json:"filters"`
	View       string                 `json:"view"` // view configuration name
}

type EntityQueryResponse struct {
	Entities   []map[string]interface{} `json:"entities"`
	TotalCount int64                    `json:"totalCount"`
	Page       int                      `json:"page"`
	PageSize   int                      `json:"pageSize"`
	TotalPages int                      `json:"totalPages"`
	HasMore    bool                     `json:"hasMore"`
	SortBy     string                   `json:"sortBy"`
	SortOrder  string                   `json:"sortOrder"`
}

type EntityViewConfig struct {
	Name         string   `json:"name"`
	DisplayName  string   `json:"displayName"`
	Columns      []Column `json:"columns"`
	DefaultSort  string   `json:"defaultSort"`
	DefaultOrder string   `json:"defaultOrder"`
}

type Column struct {
	Key        string `json:"key"`
	Label      string `json:"label"`
	Type       string `json:"type"` // "text", "number", "date", "boolean", "link", "status"
	Sortable   bool   `json:"sortable"`
	Filterable bool   `json:"filterable"`
	Width      string `json:"width"`
	Align      string `json:"align"`  // "left", "center", "right"
	Format     string `json:"format"` // "currency", "percentage", "date", etc.
}

func NewEntityHandler(db *gorm.DB) *EntityHandler {
	return &EntityHandler{db: db}
}

// GetEntityList handles generic entity queries with pagination, filtering, and sorting
func (h *EntityHandler) GetEntityList(c *gin.Context) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Convert to string if it's not already
	var userID string
	switch v := userIDValue.(type) {
	case string:
		userID = v
	case float64:
		userID = fmt.Sprintf("%v", v)
	default:
		userID = fmt.Sprintf("%v", userIDValue)
	}

	var req EntityQueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set defaults
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	if req.SortOrder == "" {
		req.SortOrder = "asc"
	}

	// Get user's tenant
	var user models.User
	if err := h.db.Select("tenant_id").First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Build query based on entity type
	query, err := h.buildEntityQuery(req.EntityType, user.TenantID, req.Filters)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Apply sorting
	if req.SortBy != "" {
		query = h.applySorting(query, req.SortBy, req.SortOrder)
	}

	// Get total count
	var totalCount int64
	if err := query.Count(&totalCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count entities"})
		return
	}

	// Apply pagination
	offset := (req.Page - 1) * req.PageSize
	query = query.Offset(offset).Limit(req.PageSize)

	// Execute query and get results
	entities, err := h.executeEntityQuery(req.EntityType, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch entities"})
		return
	}

	// Calculate pagination info
	totalPages := int((totalCount + int64(req.PageSize) - 1) / int64(req.PageSize))
	hasMore := req.Page < totalPages

	response := EntityQueryResponse{
		Entities:   entities,
		TotalCount: totalCount,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
		HasMore:    hasMore,
		SortBy:     req.SortBy,
		SortOrder:  req.SortOrder,
	}

	c.JSON(http.StatusOK, response)
}

// GetEntityViews returns available view configurations for an entity type
func (h *EntityHandler) GetEntityViews(c *gin.Context) {
	entityType := c.Param("entityType")

	// No need to check user authentication for view configurations
	// as they are not tenant-specific and don't contain sensitive data
	views := h.getDefaultViews(entityType)
	c.JSON(http.StatusOK, gin.H{"views": views})
}

// buildEntityQuery creates the base query for the specified entity type
func (h *EntityHandler) buildEntityQuery(entityType string, tenantID string, filters map[string]interface{}) (*gorm.DB, error) {
	var query *gorm.DB

	switch strings.ToLower(entityType) {
	case "companies":
		query = h.db.Model(&models.Company{}).
			Select(`
				companies.id, companies.name, companies.website, companies.domain,
				companies.revenue, companies.created_at, companies.updated_at,
				industries.name as industry_name,
				company_sizes.name as size_name,
				phone_numbers.number as phone,
				email_addresses.email as email,
				COUNT(DISTINCT contacts.id) as contact_count,
				COUNT(DISTINCT leads.id) as lead_count,
				COUNT(DISTINCT deals.id) as deal_count
			`).
			Joins("LEFT JOIN industries ON companies.industry_id = industries.id").
			Joins("LEFT JOIN company_sizes ON companies.size_id = company_sizes.id").
			Joins("LEFT JOIN phone_numbers ON companies.id = phone_numbers.entity_id AND phone_numbers.entity_type = 'company' AND phone_numbers.is_primary = true").
			Joins("LEFT JOIN email_addresses ON companies.id = email_addresses.entity_id AND email_addresses.entity_type = 'company' AND email_addresses.is_primary = true").
			Joins("LEFT JOIN contacts ON companies.id = contacts.company_id").
			Joins("LEFT JOIN leads ON companies.id = leads.company_id").
			Joins("LEFT JOIN deals ON companies.id = deals.company_id").
			Where("companies.tenant_id = ?", tenantID).
			Group("companies.id, industries.name, company_sizes.name, phone_numbers.number, email_addresses.email")

	case "contacts":
		query = h.db.Model(&models.Contact{}).
			Select(`
				contacts.id, contacts.first_name, contacts.last_name, contacts.title, 
				contacts.department, contacts.created_at, contacts.updated_at,
				companies.name as company_name,
				phone_numbers.number as phone,
				email_addresses.email as email,
				COUNT(DISTINCT leads.id) as lead_count,
				COUNT(DISTINCT deals.id) as deal_count
			`).
			Joins("LEFT JOIN companies ON contacts.company_id = companies.id").
			Joins("LEFT JOIN phone_numbers ON contacts.id = phone_numbers.entity_id AND phone_numbers.entity_type = 'contact' AND phone_numbers.is_primary = true").
			Joins("LEFT JOIN email_addresses ON contacts.id = email_addresses.entity_id AND email_addresses.entity_type = 'contact' AND email_addresses.is_primary = true").
			Joins("LEFT JOIN leads ON contacts.id = leads.contact_id").
			Joins("LEFT JOIN deals ON contacts.id = deals.contact_id").
			Where("contacts.tenant_id = ?", tenantID).
			Group("contacts.id, companies.name, phone_numbers.number, email_addresses.email")

	case "leads":
		// Add a virtual company_name field that combines the actual company name (from joined table)
		// or a constant string if no company is associated
		query = h.db.Model(&models.Lead{}).
			Select(`
				leads.id, leads.first_name, leads.last_name, leads.title, leads.score,
				leads.source, leads.campaign, leads.created_at, leads.updated_at,
				lead_statuses.name as status_name,
				lead_temperatures.name as temperature_name,
				COALESCE(companies.name, 'Unknown Company') as company_name,
				contacts.first_name as contact_first_name,
				contacts.last_name as contact_last_name,
				phone_numbers.number as phone,
				email_addresses.email as email
			`).
			Joins("LEFT JOIN lead_statuses ON leads.status_id = lead_statuses.id").
			Joins("LEFT JOIN lead_temperatures ON leads.temperature_id = lead_temperatures.id").
			Joins("LEFT JOIN companies ON leads.company_id = companies.id").
			Joins("LEFT JOIN contacts ON leads.contact_id = contacts.id").
			Joins("LEFT JOIN phone_numbers ON leads.id = phone_numbers.entity_id AND phone_numbers.entity_type = 'lead' AND phone_numbers.is_primary = true").
			Joins("LEFT JOIN email_addresses ON leads.id = email_addresses.entity_id AND email_addresses.entity_type = 'lead' AND email_addresses.is_primary = true").
			Where("leads.tenant_id = ?", tenantID)

	case "deals":
		query = h.db.Model(&models.Deal{}).
			Select(`
				deals.id, deals.name, deals.amount, deals.currency, deals.expected_close_date,
				deals.probability, deals.created_at, deals.updated_at,
				stages.name as stage_name,
				companies.name as company_name,
				contacts.first_name as contact_first_name,
				contacts.last_name as contact_last_name,
				users.first_name as owner_first_name,
				users.last_name as owner_last_name
			`).
			Joins("LEFT JOIN stages ON deals.stage_id = stages.id").
			Joins("LEFT JOIN companies ON deals.company_id = companies.id").
			Joins("LEFT JOIN contacts ON deals.contact_id = contacts.id").
			Joins("LEFT JOIN users ON deals.assigned_user_id = users.id").
			Where("deals.tenant_id = ?", tenantID)

	default:
		return nil, fmt.Errorf("unsupported entity type: %s", entityType)
	}

	// Apply filters
	query = h.applyFilters(query, filters)

	return query, nil
}

// applyFilters applies the provided filters to the query
func (h *EntityHandler) applyFilters(query *gorm.DB, filters map[string]interface{}) *gorm.DB {
	if filters == nil {
		return query
	}

	for key, value := range filters {
		if value == nil || value == "" {
			continue
		}

		switch key {
		case "search":
			// Generic search across name/email fields
			searchTerm := fmt.Sprintf("%%%v%%", value)
			query = query.Where(
				"companies.name ILIKE ? OR contacts.first_name ILIKE ? OR contacts.last_name ILIKE ? OR leads.first_name ILIKE ? OR leads.last_name ILIKE ? OR deals.name ILIKE ?",
				searchTerm, searchTerm, searchTerm, searchTerm, searchTerm, searchTerm,
			)
		case "industry_id":
			query = query.Where("companies.industry_id = ?", value)
		case "size_id":
			query = query.Where("companies.size_id = ?", value)
		case "status_id":
			query = query.Where("leads.status_id = ?", value)
		case "temperature_id":
			query = query.Where("leads.temperature_id = ?", value)
		case "stage_id":
			query = query.Where("deals.stage_id = ?", value)
		case "owner_id":
			query = query.Where("deals.assigned_user_id = ?", value)
		case "company_id":
			query = query.Where("contacts.company_id = ? OR leads.company_id = ? OR deals.company_id = ?", value, value, value)
		case "created_after":
			query = query.Where("created_at >= ?", value)
		case "created_before":
			query = query.Where("created_at <= ?", value)
		case "amount_min":
			query = query.Where("deals.amount >= ?", value)
		case "amount_max":
			query = query.Where("deals.amount <= ?", value)
		}
	}

	return query
}

// applySorting applies sorting to the query
func (h *EntityHandler) applySorting(query *gorm.DB, sortBy, sortOrder string) *gorm.DB {
	// Map frontend column keys to database fields
	fieldMap := map[string]string{
		"name":           "companies.name",
		"company_name":   "companies.name",
		"first_name":     "contacts.first_name",
		"last_name":      "contacts.last_name",
		"email":          "contacts.email",
		"phone":          "contacts.phone",
		"title":          "contacts.title",
		"industry":       "industries.name",
		"size":           "company_sizes.name",
		"status":         "lead_statuses.name",
		"temperature":    "lead_temperatures.name",
		"stage":          "stages.name",
		"amount":         "deals.amount",
		"probability":    "deals.probability",
		"expected_close": "deals.expected_close_date",
		"created_at":     "created_at",
		"updated_at":     "updated_at",
		"contact_count":  "contact_count",
		"lead_count":     "lead_count",
		"deal_count":     "deal_count",
	}

	if dbField, exists := fieldMap[sortBy]; exists {
		if strings.ToLower(sortOrder) == "desc" {
			query = query.Order(clause.OrderByColumn{Column: clause.Column{Name: dbField}, Desc: true})
		} else {
			query = query.Order(clause.OrderByColumn{Column: clause.Column{Name: dbField}, Desc: false})
		}
	}

	return query
}

// executeEntityQuery executes the query and returns the results
func (h *EntityHandler) executeEntityQuery(entityType string, query *gorm.DB) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	rows, err := query.Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	// Create a slice of interface{} to hold the values
	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	for rows.Next() {
		err := rows.Scan(valuePtrs...)
		if err != nil {
			return nil, err
		}

		// Convert the row to a map
		row := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			if val != nil {
				row[col] = val
			} else {
				row[col] = nil
			}
		}

		results = append(results, row)
	}

	return results, nil
}

// getDefaultViews returns the default view configurations for each entity type
func (h *EntityHandler) getDefaultViews(entityType string) []EntityViewConfig {
	switch strings.ToLower(entityType) {
	case "companies":
		return []EntityViewConfig{
			{
				Name:         "overview",
				DisplayName:  "Overview",
				DefaultSort:  "name",
				DefaultOrder: "asc",
				Columns: []Column{
					{Key: "name", Label: "Company Name", Type: "text", Sortable: true, Filterable: true, Width: "200px"},
					{Key: "industry_name", Label: "Industry", Type: "text", Sortable: true, Filterable: true, Width: "150px"},
					{Key: "size_name", Label: "Size", Type: "text", Sortable: true, Filterable: true, Width: "120px"},
					{Key: "website", Label: "Website", Type: "link", Sortable: false, Filterable: true, Width: "150px"},
					{Key: "phone", Label: "Phone", Type: "text", Sortable: false, Filterable: true, Width: "120px"},
					{Key: "email", Label: "Email", Type: "text", Sortable: false, Filterable: true, Width: "200px"},
					{Key: "revenue", Label: "Revenue", Type: "currency", Sortable: true, Filterable: true, Width: "120px", Align: "right"},
					{Key: "contact_count", Label: "Contacts", Type: "number", Sortable: true, Filterable: false, Width: "100px", Align: "center"},
					{Key: "lead_count", Label: "Leads", Type: "number", Sortable: true, Filterable: false, Width: "100px", Align: "center"},
					{Key: "deal_count", Label: "Deals", Type: "number", Sortable: true, Filterable: false, Width: "100px", Align: "center"},
					{Key: "created_at", Label: "Created", Type: "date", Sortable: true, Filterable: true, Width: "120px"},
				},
			},
			{
				Name:         "detailed",
				DisplayName:  "Detailed View",
				DefaultSort:  "created_at",
				DefaultOrder: "desc",
				Columns: []Column{
					{Key: "name", Label: "Company Name", Type: "text", Sortable: true, Filterable: true, Width: "200px"},
					{Key: "industry_name", Label: "Industry", Type: "text", Sortable: true, Filterable: true, Width: "150px"},
					{Key: "size_name", Label: "Size", Type: "text", Sortable: true, Filterable: true, Width: "120px"},
					{Key: "website", Label: "Website", Type: "link", Sortable: false, Filterable: true, Width: "150px"},
					{Key: "phone", Label: "Phone", Type: "text", Sortable: false, Filterable: true, Width: "120px"},
					{Key: "contact_count", Label: "Contacts", Type: "number", Sortable: true, Filterable: false, Width: "100px", Align: "center"},
					{Key: "lead_count", Label: "Leads", Type: "number", Sortable: true, Filterable: false, Width: "100px", Align: "center"},
					{Key: "deal_count", Label: "Deals", Type: "number", Sortable: true, Filterable: false, Width: "100px", Align: "center"},
					{Key: "created_at", Label: "Created", Type: "date", Sortable: true, Filterable: true, Width: "120px"},
					{Key: "updated_at", Label: "Updated", Type: "date", Sortable: true, Filterable: true, Width: "120px"},
				},
			},
		}

	case "contacts":
		return []EntityViewConfig{
			{
				Name:         "overview",
				DisplayName:  "Overview",
				DefaultSort:  "last_name",
				DefaultOrder: "asc",
				Columns: []Column{
					{Key: "first_name", Label: "First Name", Type: "text", Sortable: true, Filterable: true, Width: "120px"},
					{Key: "last_name", Label: "Last Name", Type: "text", Sortable: true, Filterable: true, Width: "120px"},
					{Key: "company_name", Label: "Company", Type: "text", Sortable: true, Filterable: true, Width: "200px"},
					{Key: "title", Label: "Title", Type: "text", Sortable: true, Filterable: true, Width: "150px"},
					{Key: "department", Label: "Department", Type: "text", Sortable: true, Filterable: true, Width: "120px"},
					{Key: "email", Label: "Email", Type: "text", Sortable: false, Filterable: true, Width: "200px"},
					{Key: "phone", Label: "Phone", Type: "text", Sortable: false, Filterable: true, Width: "120px"},
					{Key: "lead_count", Label: "Leads", Type: "number", Sortable: true, Filterable: false, Width: "80px", Align: "center"},
					{Key: "deal_count", Label: "Deals", Type: "number", Sortable: true, Filterable: false, Width: "80px", Align: "center"},
					{Key: "created_at", Label: "Created", Type: "date", Sortable: true, Filterable: true, Width: "120px"},
				},
			},
		}

	case "leads":
		return []EntityViewConfig{
			{
				Name:         "overview",
				DisplayName:  "Overview",
				DefaultSort:  "created_at",
				DefaultOrder: "desc",
				Columns: []Column{
					{Key: "first_name", Label: "First Name", Type: "text", Sortable: true, Filterable: true, Width: "120px"},
					{Key: "last_name", Label: "Last Name", Type: "text", Sortable: true, Filterable: true, Width: "120px"},
					{Key: "company_name", Label: "Company", Type: "text", Sortable: true, Filterable: true, Width: "200px"},
					{Key: "title", Label: "Title", Type: "text", Sortable: true, Filterable: true, Width: "150px"},
					{Key: "score", Label: "Score", Type: "number", Sortable: true, Filterable: true, Width: "80px", Align: "center"},
					{Key: "source", Label: "Source", Type: "text", Sortable: true, Filterable: true, Width: "120px"},
					{Key: "campaign", Label: "Campaign", Type: "text", Sortable: true, Filterable: true, Width: "120px"},
					{Key: "email", Label: "Email", Type: "text", Sortable: false, Filterable: true, Width: "200px"},
					{Key: "phone", Label: "Phone", Type: "text", Sortable: false, Filterable: true, Width: "120px"},
					{Key: "status_name", Label: "Status", Type: "status", Sortable: true, Filterable: true, Width: "120px"},
					{Key: "temperature_name", Label: "Temperature", Type: "status", Sortable: true, Filterable: true, Width: "120px"},
					{Key: "created_at", Label: "Created", Type: "date", Sortable: true, Filterable: true, Width: "120px"},
				},
			},
		}

	case "deals":
		return []EntityViewConfig{
			{
				Name:         "overview",
				DisplayName:  "Overview",
				DefaultSort:  "expected_close_date",
				DefaultOrder: "asc",
				Columns: []Column{
					{Key: "name", Label: "Deal Name", Type: "text", Sortable: true, Filterable: true, Width: "200px"},
					{Key: "company_name", Label: "Company", Type: "text", Sortable: true, Filterable: true, Width: "200px"},
					{Key: "stage_name", Label: "Stage", Type: "status", Sortable: true, Filterable: true, Width: "120px"},
					{Key: "amount", Label: "Amount", Type: "currency", Sortable: true, Filterable: true, Width: "120px", Align: "right"},
					{Key: "probability", Label: "Probability", Type: "percentage", Sortable: true, Filterable: true, Width: "100px", Align: "center"},
					{Key: "expected_close_date", Label: "Close Date", Type: "date", Sortable: true, Filterable: true, Width: "120px"},
					{Key: "owner_first_name", Label: "Owner", Type: "text", Sortable: true, Filterable: true, Width: "150px"},
					{Key: "created_at", Label: "Created", Type: "date", Sortable: true, Filterable: true, Width: "120px"},
				},
			},
		}

	default:
		return []EntityViewConfig{}
	}
}
