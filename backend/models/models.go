package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ============================================================================
// CORE USER MANAGEMENT
// ============================================================================

type User struct {
	ID        string  `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Email     string  `json:"email" gorm:"uniqueIndex;not null"`
	FirstName string  `json:"firstName" gorm:"column:first_name;not null"`
	LastName  string  `json:"lastName" gorm:"column:last_name;not null"`
	Avatar    *string `json:"avatar"`
	Password  string  `json:"-" gorm:"not null"`

	RoleID *string   `json:"roleId" gorm:"column:role_id;type:uuid"`
	Role   *UserRole `json:"role,omitempty" gorm:"foreignKey:RoleID"`

	IsActive bool `json:"isActive" gorm:"column:is_active;default:true"`

	TenantID string `json:"tenantId" gorm:"column:tenant_id;type:uuid;not null"`
	Tenant   Tenant `json:"tenant,omitempty" gorm:"foreignKey:TenantID"`

	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`
	CreatedBy *string   `json:"createdBy" gorm:"column:created_by;type:uuid"`
}

type UserRole struct {
	ID          string  `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Name        string  `json:"name" gorm:"not null"`
	Code        string  `json:"code" gorm:"not null"`
	Description *string `json:"description"`
	IsActive    bool    `json:"isActive" gorm:"column:is_active;default:true"`
	IsSystem    bool    `json:"isSystem" gorm:"column:is_system;default:false"`

	Permissions interface{} `json:"permissions" gorm:"type:jsonb"`

	TenantID string `json:"tenantId" gorm:"column:tenant_id;type:uuid;not null"`
	Tenant   Tenant `json:"tenant,omitempty" gorm:"foreignKey:TenantID"`

	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`

	// Remove the Users field to avoid circular reference
	// Users []User `json:"users,omitempty"`
}

type Tenant struct {
	ID        string `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Name      string `json:"name" gorm:"not null"`
	Subdomain string `json:"subdomain" gorm:"uniqueIndex;not null"`
	IsActive  bool   `json:"isActive" gorm:"column:is_active;default:true"`

	Settings interface{} `json:"settings" gorm:"type:jsonb"`

	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`

	// Remove these fields to avoid circular references
	// Users            []User            `json:"users,omitempty"`
	// Companies        []Company         `json:"companies,omitempty"`
	// Contacts         []Contact         `json:"contacts,omitempty"`
	// Leads            []Lead            `json:"leads,omitempty"`
	// Deals            []Deal            `json:"deals,omitempty"`
	// UserRoles        []UserRole        `json:"userRoles,omitempty"`
	// LeadStatuses     []LeadStatus      `json:"leadStatuses,omitempty"`
	// Industries       []Industry        `json:"industries,omitempty"`
	// CompanySizes     []CompanySize     `json:"companySizes,omitempty"`
	// LeadTemperatures []LeadTemperature `json:"leadTemperatures,omitempty"`
}

// ============================================================================
// LOOKUP TABLES
// ============================================================================

type LeadStatus struct {
	ID          string  `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Name        string  `json:"name" gorm:"not null"`
	Code        string  `json:"code" gorm:"not null"`
	Description *string `json:"description"`
	Color       *string `json:"color"`
	Order       int     `json:"order" gorm:"not null"`
	IsActive    bool    `json:"isActive" gorm:"column:is_active;default:true"`
	IsSystem    bool    `json:"isSystem" gorm:"column:is_system;default:false"`

	TenantID string `json:"tenantId" gorm:"column:tenant_id;type:uuid;not null"`
	Tenant   Tenant `json:"tenant,omitempty" gorm:"foreignKey:TenantID"`

	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`

	// Remove the Leads field to avoid circular reference
	// Leads []Lead `json:"leads,omitempty"`
}

type LeadTemperature struct {
	ID          string  `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Name        string  `json:"name" gorm:"not null"`
	Code        string  `json:"code" gorm:"not null"`
	Description *string `json:"description"`
	Color       *string `json:"color"`
	Order       int     `json:"order" gorm:"not null"`
	IsActive    bool    `json:"isActive" gorm:"column:is_active;default:true"`

	TenantID string `json:"tenantId" gorm:"column:tenant_id;type:uuid;not null"`
	Tenant   Tenant `json:"tenant,omitempty" gorm:"foreignKey:TenantID"`

	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`

	// Remove the Leads field to avoid circular reference
	// Leads []Lead `json:"leads,omitempty"`
}

type Industry struct {
	ID          string  `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Name        string  `json:"name" gorm:"not null"`
	Code        string  `json:"code" gorm:"not null"`
	Description *string `json:"description"`
	IsActive    bool    `json:"isActive" gorm:"column:is_active;default:true"`

	TenantID string `json:"tenantId" gorm:"column:tenant_id;type:uuid;not null"`
	Tenant   Tenant `json:"tenant,omitempty" gorm:"foreignKey:TenantID"`

	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`

	// Remove the Companies field to avoid circular reference
	// Companies []Company `json:"companies,omitempty"`
}

type CompanySize struct {
	ID           string  `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Name         string  `json:"name" gorm:"not null"`
	Code         string  `json:"code" gorm:"not null"`
	Description  *string `json:"description"`
	MinEmployees *int    `json:"minEmployees" gorm:"column:min_employees"`
	MaxEmployees *int    `json:"maxEmployees" gorm:"column:max_employees"`
	IsActive     bool    `json:"isActive" gorm:"column:is_active;default:true"`

	TenantID string `json:"tenantId" gorm:"column:tenant_id;type:uuid;not null"`
	Tenant   Tenant `json:"tenant,omitempty" gorm:"foreignKey:TenantID"`

	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`

	// Remove the Companies field to avoid circular reference
	// Companies []Company `json:"companies,omitempty"`
}

// ============================================================================
// COMPANIES AND CONTACTS
// ============================================================================

type Company struct {
	ID      string  `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Name    string  `json:"name" gorm:"not null"`
	Website *string `json:"website"`
	Domain  *string `json:"domain"`

	IndustryID *string   `json:"industryId" gorm:"column:industry_id;type:uuid"`
	Industry   *Industry `json:"industry,omitempty" gorm:"foreignKey:IndustryID"`

	SizeID *string      `json:"sizeId" gorm:"column:size_id;type:uuid"`
	Size   *CompanySize `json:"size,omitempty" gorm:"foreignKey:SizeID"`

	Revenue    *float64 `json:"revenue"`
	ExternalID *string  `json:"externalId" gorm:"column:external_id"`

	AssignedUserID *string `json:"assignedUserId" gorm:"column:assigned_user_id;type:uuid"`
	AssignedUser   *User   `json:"assignedUser,omitempty" gorm:"foreignKey:AssignedUserID"`

	TenantID string `json:"tenantId" gorm:"column:tenant_id;type:uuid;not null"`
	Tenant   Tenant `json:"tenant,omitempty" gorm:"foreignKey:TenantID"`

	CreatedAt time.Time  `json:"createdAt" gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time  `json:"updatedAt" gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`
	CreatedBy *string    `json:"createdBy" gorm:"column:created_by;type:uuid"`
	IsDeleted bool       `json:"isDeleted" gorm:"column:is_deleted;default:false"`
	DeletedAt *time.Time `json:"deletedAt" gorm:"column:deleted_at"`

	// Remove these fields to avoid circular references
	// Contacts []Contact `json:"contacts,omitempty"`
	// Deals    []Deal    `json:"deals,omitempty"`
	// Leads    []Lead    `json:"leads,omitempty"`
}

type Contact struct {
	ID        string  `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Title     *string `json:"title"`
	FirstName string  `json:"firstName" gorm:"column:first_name;not null"`
	Prefix    *string `json:"prefix"`
	Suffix    *string `json:"suffix"`
	LastName  string  `json:"lastName" gorm:"column:last_name;not null"`

	JobTitle *string     `json:"jobTitle"`
	StatusID *string     `json:"statusId" gorm:"column:status_id;type:uuid"`
	Status   *LeadStatus `json:"status,omitempty" gorm:"foreignKey:StatusID"`

	Department *string `json:"department"`

	CompanyID *string  `json:"companyId" gorm:"column:company_id;type:uuid"`
	Company   *Company `json:"company,omitempty" gorm:"foreignKey:CompanyID"`

	OriginalSource *string `json:"originalSource" gorm:"column:original_source"`

	EmailOptIn bool `json:"emailOptIn" gorm:"column:email_opt_in;default:true"`
	SmsOptIn   bool `json:"smsOptIn" gorm:"column:sms_opt_in;default:false"`
	CallOptIn  bool `json:"callOptIn" gorm:"column:call_opt_in;default:true"`

	TenantID string `json:"tenantId" gorm:"column:tenant_id;type:uuid;not null"`
	Tenant   Tenant `json:"tenant,omitempty" gorm:"foreignKey:TenantID"`

	CreatedAt time.Time  `json:"createdAt" gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time  `json:"updatedAt" gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`
	CreatedBy *string    `json:"createdBy" gorm:"column:created_by;type:uuid"`
	IsDeleted bool       `json:"isDeleted" gorm:"column:is_deleted;default:false"`
	DeletedAt *time.Time `json:"deletedAt" gorm:"column:deleted_at"`

	// Remove these fields to avoid circular references
	// Leads []Lead `json:"leads,omitempty"`
	// Deals []Deal `json:"deals,omitempty"`
}

// ============================================================================
// LEADS AND DEALS
// ============================================================================

type Lead struct {
	ID        string  `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	FirstName *string `json:"firstName" gorm:"column:first_name"`
	LastName  *string `json:"lastName" gorm:"column:last_name"`
	Title     *string `json:"title"`

	StatusID *string     `json:"statusId" gorm:"column:status_id;type:uuid"`
	Status   *LeadStatus `json:"status,omitempty" gorm:"foreignKey:StatusID"`

	TemperatureID *string          `json:"temperatureId" gorm:"column:temperature_id;type:uuid"`
	Temperature   *LeadTemperature `json:"temperature,omitempty" gorm:"foreignKey:TemperatureID"`

	Source   *string `json:"source"`
	Campaign *string `json:"campaign"`
	Score    int     `json:"score" gorm:"default:0"`

	CompanyID *string  `json:"companyId" gorm:"column:company_id;type:uuid"`
	Company   *Company `json:"company,omitempty" gorm:"foreignKey:CompanyID"`

	ContactID *string  `json:"contactId" gorm:"column:contact_id;type:uuid;uniqueIndex"`
	Contact   *Contact `json:"contact,omitempty" gorm:"foreignKey:ContactID"`

	AssignedUserID *string `json:"assignedUserId" gorm:"column:assigned_user_id;type:uuid"`
	AssignedUser   *User   `json:"assignedUser,omitempty" gorm:"foreignKey:AssignedUserID"`

	ConvertedAt       *time.Time `json:"convertedAt" gorm:"column:converted_at"`
	ConvertedToDealID *string    `json:"convertedToDealId" gorm:"column:converted_to_deal_id;type:uuid"`

	TenantID string `json:"tenantId" gorm:"column:tenant_id;type:uuid;not null"`
	Tenant   Tenant `json:"tenant,omitempty" gorm:"foreignKey:TenantID"`

	CreatedAt time.Time  `json:"createdAt" gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time  `json:"updatedAt" gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`
	CreatedBy *string    `json:"createdBy" gorm:"column:created_by;type:uuid"`
	IsDeleted bool       `json:"isDeleted" gorm:"column:is_deleted;default:false"`
	DeletedAt *time.Time `json:"deletedAt" gorm:"column:deleted_at"`
}

type Deal struct {
	ID          string   `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Name        string   `json:"name" gorm:"not null"`
	Amount      *float64 `json:"amount"`
	Currency    string   `json:"currency" gorm:"default:USD"`
	Probability int      `json:"probability" gorm:"default:0"`

	PipelineID string   `json:"pipelineId" gorm:"column:pipeline_id;type:uuid;not null"`
	Pipeline   Pipeline `json:"pipeline,omitempty" gorm:"foreignKey:PipelineID"`

	StageID string `json:"stageId" gorm:"column:stage_id;type:uuid;not null"`
	Stage   Stage  `json:"stage,omitempty" gorm:"foreignKey:StageID"`

	ExpectedCloseDate *time.Time `json:"expectedCloseDate" gorm:"column:expected_close_date"`
	ActualCloseDate   *time.Time `json:"actualCloseDate" gorm:"column:actual_close_date"`

	CompanyID *string  `json:"companyId" gorm:"column:company_id;type:uuid"`
	Company   *Company `json:"company,omitempty" gorm:"foreignKey:CompanyID"`

	ContactID *string  `json:"contactId" gorm:"column:contact_id;type:uuid"`
	Contact   *Contact `json:"contact,omitempty" gorm:"foreignKey:ContactID"`

	AssignedUserID *string `json:"assignedUserId" gorm:"column:assigned_user_id;type:uuid"`
	AssignedUser   *User   `json:"assignedUser,omitempty" gorm:"foreignKey:AssignedUserID"`

	TenantID string `json:"tenantId" gorm:"column:tenant_id;type:uuid;not null"`
	Tenant   Tenant `json:"tenant,omitempty" gorm:"foreignKey:TenantID"`

	CreatedAt time.Time  `json:"createdAt" gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time  `json:"updatedAt" gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`
	CreatedBy *string    `json:"createdBy" gorm:"column:created_by;type:uuid"`
	IsDeleted bool       `json:"isDeleted" gorm:"column:is_deleted;default:false"`
	DeletedAt *time.Time `json:"deletedAt" gorm:"column:deleted_at"`
}

type Pipeline struct {
	ID       string `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Name     string `json:"name" gorm:"not null"`
	IsActive bool   `json:"isActive" gorm:"column:is_active;default:true"`

	TenantID string `json:"tenantId" gorm:"column:tenant_id;type:uuid;not null"`
	Tenant   Tenant `json:"tenant,omitempty" gorm:"foreignKey:TenantID"`

	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`

	// Remove these fields to avoid circular references
	// Stages []Stage `json:"stages,omitempty"`
	// Deals  []Deal  `json:"deals,omitempty"`
}

type Stage struct {
	ID           string  `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Name         string  `json:"name" gorm:"not null"`
	Order        int     `json:"order" gorm:"not null"`
	Probability  int     `json:"probability" gorm:"default:0"`
	IsClosedWon  bool    `json:"isClosedWon" gorm:"column:is_closed_won;default:false"`
	IsClosedLost bool    `json:"isClosedLost" gorm:"column:is_closed_lost;default:false"`
	Color        *string `json:"color"`

	PipelineID string   `json:"pipelineId" gorm:"column:pipeline_id;type:uuid;not null"`
	Pipeline   Pipeline `json:"pipeline,omitempty" gorm:"foreignKey:PipelineID"`

	TenantID string `json:"tenantId" gorm:"column:tenant_id;type:uuid;not null"`
	Tenant   Tenant `json:"tenant,omitempty" gorm:"foreignKey:TenantID"`

	// Remove the Deals field to avoid circular reference
	// Deals []Deal `json:"deals,omitempty"`
}

// ============================================================================
// MARKETING AND COMMUNICATIONS
// ============================================================================

type MarketingSourceType struct {
	ID          string  `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Name        string  `json:"name" gorm:"not null"`
	Code        string  `json:"code" gorm:"not null"`
	Description *string `json:"description"`
	Color       *string `json:"color"`
	IsActive    bool    `json:"isActive" gorm:"column:is_active;default:true"`

	TenantID string `json:"tenantId" gorm:"column:tenant_id;type:uuid;not null"`
	Tenant   Tenant `json:"tenant,omitempty" gorm:"foreignKey:TenantID"`

	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`
}

type MarketingSource struct {
	ID       string  `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Name     string  `json:"name" gorm:"not null"`
	TypeID   *string `json:"typeId" gorm:"column:type_id;type:uuid"`
	Medium   *string `json:"medium"`
	Campaign *string `json:"campaign"`
	Source   *string `json:"source"`
	Content  *string `json:"content"`
	Term     *string `json:"term"`

	TenantID string `json:"tenantId" gorm:"column:tenant_id;type:uuid;not null"`
	Tenant   Tenant `json:"tenant,omitempty" gorm:"foreignKey:TenantID"`

	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`
}

type MarketingAssetType struct {
	ID          string  `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Name        string  `json:"name" gorm:"not null"`
	Code        string  `json:"code" gorm:"not null"`
	Description *string `json:"description"`
	Color       *string `json:"color"`
	IsActive    bool    `json:"isActive" gorm:"column:is_active;default:true"`

	TenantID string `json:"tenantId" gorm:"column:tenant_id;type:uuid;not null"`
	Tenant   Tenant `json:"tenant,omitempty" gorm:"foreignKey:TenantID"`

	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`
}

type MarketingAsset struct {
	ID          string  `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Name        string  `json:"name" gorm:"not null"`
	TypeID      *string `json:"typeId" gorm:"column:type_id;type:uuid"`
	URL         *string `json:"url"`
	Content     *string `json:"content"`
	Views       int     `json:"views" gorm:"default:0"`
	Clicks      int     `json:"clicks" gorm:"default:0"`
	Conversions int     `json:"conversions" gorm:"default:0"`

	TenantID string `json:"tenantId" gorm:"column:tenant_id;type:uuid;not null"`
	Tenant   Tenant `json:"tenant,omitempty" gorm:"foreignKey:TenantID"`

	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`
}

type CommunicationType struct {
	ID          string  `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Name        string  `json:"name" gorm:"not null"`
	Code        string  `json:"code" gorm:"not null"`
	Description *string `json:"description"`
	Icon        *string `json:"icon"`
	IsActive    bool    `json:"isActive" gorm:"column:is_active;default:true"`

	TenantID string `json:"tenantId" gorm:"column:tenant_id;type:uuid;not null"`
	Tenant   Tenant `json:"tenant,omitempty" gorm:"foreignKey:TenantID"`

	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`
}

// ============================================================================
// TASKS AND COMMUNICATIONS
// ============================================================================

type TaskType struct {
	ID          string  `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Name        string  `json:"name" gorm:"not null"`
	Code        string  `json:"code" gorm:"not null"`
	Description *string `json:"description"`
	Color       *string `json:"color"`
	Icon        *string `json:"icon"`
	IsActive    bool    `json:"isActive" gorm:"column:is_active;default:true"`

	TenantID string `json:"tenantId" gorm:"column:tenant_id;type:uuid;not null"`
	Tenant   Tenant `json:"tenant,omitempty" gorm:"foreignKey:TenantID"`

	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`
}

type Task struct {
	ID          string  `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Title       string  `json:"title" gorm:"not null"`
	Description *string `json:"description"`

	Priority TaskPriority `json:"priority" gorm:"default:MEDIUM"`
	Status   TaskStatus   `json:"status" gorm:"default:PENDING"`

	DueDate     *time.Time `json:"dueDate" gorm:"column:due_date"`
	CompletedAt *time.Time `json:"completedAt" gorm:"column:completed_at"`

	AssignedUserID *string `json:"assignedUserId" gorm:"column:assigned_user_id;type:uuid"`
	AssignedUser   *User   `json:"assignedUser,omitempty" gorm:"foreignKey:AssignedUserID"`

	CreatedBy *string `json:"createdBy" gorm:"column:created_by;type:uuid"`

	LeadID *string `json:"leadId" gorm:"column:lead_id;type:uuid"`
	Lead   *Lead   `json:"lead,omitempty" gorm:"foreignKey:LeadID"`

	DealID *string `json:"dealId" gorm:"column:deal_id;type:uuid"`
	Deal   *Deal   `json:"deal,omitempty" gorm:"foreignKey:DealID"`

	TenantID string `json:"tenantId" gorm:"column:tenant_id;type:uuid;not null"`
	Tenant   Tenant `json:"tenant,omitempty" gorm:"foreignKey:TenantID"`

	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`
}

type Communication struct {
	ID        string                 `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Subject   *string                `json:"subject"`
	Content   *string                `json:"content"`
	Direction CommunicationDirection `json:"direction" gorm:"not null"`

	ScheduledAt *time.Time `json:"scheduledAt" gorm:"column:scheduled_at"`
	SentAt      *time.Time `json:"sentAt" gorm:"column:sent_at"`
	ReceivedAt  *time.Time `json:"receivedAt" gorm:"column:received_at"`

	ExternalID *string `json:"externalId" gorm:"column:external_id"`

	UserID *string `json:"userId" gorm:"column:user_id;type:uuid"`
	User   *User   `json:"user,omitempty" gorm:"foreignKey:UserID"`

	ContactID *string  `json:"contactId" gorm:"column:contact_id;type:uuid"`
	Contact   *Contact `json:"contact,omitempty" gorm:"foreignKey:ContactID"`

	LeadID *string `json:"leadId" gorm:"column:lead_id;type:uuid"`
	Lead   *Lead   `json:"lead,omitempty" gorm:"foreignKey:LeadID"`

	DealID *string `json:"dealId" gorm:"column:deal_id;type:uuid"`
	Deal   *Deal   `json:"deal,omitempty" gorm:"foreignKey:DealID"`

	TenantID string `json:"tenantId" gorm:"column:tenant_id;type:uuid;not null"`
	Tenant   Tenant `json:"tenant,omitempty" gorm:"foreignKey:TenantID"`

	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`
}

// ============================================================================
// ENUMS
// ============================================================================

type TaskPriority string

const (
	TaskPriorityLow    TaskPriority = "LOW"
	TaskPriorityMedium TaskPriority = "MEDIUM"
	TaskPriorityHigh   TaskPriority = "HIGH"
	TaskPriorityUrgent TaskPriority = "URGENT"
)

type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "PENDING"
	TaskStatusInProgress TaskStatus = "IN_PROGRESS"
	TaskStatusCompleted  TaskStatus = "COMPLETED"
	TaskStatusCancelled  TaskStatus = "CANCELLED"
	TaskStatusOverdue    TaskStatus = "OVERDUE"
)

type CommunicationDirection string

const (
	CommunicationDirectionInbound  CommunicationDirection = "INBOUND"
	CommunicationDirectionOutbound CommunicationDirection = "OUTBOUND"
	CommunicationDirectionInternal CommunicationDirection = "INTERNAL"
)

// ============================================================================
// HOOKS
// ============================================================================

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}

func (r *UserRole) BeforeCreate(tx *gorm.DB) error {
	if r.ID == "" {
		r.ID = uuid.New().String()
	}
	return nil
}

func (t *Tenant) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	return nil
}

func (l *LeadStatus) BeforeCreate(tx *gorm.DB) error {
	if l.ID == "" {
		l.ID = uuid.New().String()
	}
	return nil
}

func (lt *LeadTemperature) BeforeCreate(tx *gorm.DB) error {
	if lt.ID == "" {
		lt.ID = uuid.New().String()
	}
	return nil
}

func (i *Industry) BeforeCreate(tx *gorm.DB) error {
	if i.ID == "" {
		i.ID = uuid.New().String()
	}
	return nil
}

func (cs *CompanySize) BeforeCreate(tx *gorm.DB) error {
	if cs.ID == "" {
		cs.ID = uuid.New().String()
	}
	return nil
}

func (c *Company) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return nil
}

func (c *Contact) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return nil
}

func (l *Lead) BeforeCreate(tx *gorm.DB) error {
	if l.ID == "" {
		l.ID = uuid.New().String()
	}
	return nil
}

func (d *Deal) BeforeCreate(tx *gorm.DB) error {
	if d.ID == "" {
		d.ID = uuid.New().String()
	}
	return nil
}

func (p *Pipeline) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return nil
}

func (s *Stage) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	return nil
}

func (t *Task) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	return nil
}

func (c *Communication) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return nil
}

func (mst *MarketingSourceType) BeforeCreate(tx *gorm.DB) error {
	if mst.ID == "" {
		mst.ID = uuid.New().String()
	}
	return nil
}

func (ms *MarketingSource) BeforeCreate(tx *gorm.DB) error {
	if ms.ID == "" {
		ms.ID = uuid.New().String()
	}
	return nil
}

func (mat *MarketingAssetType) BeforeCreate(tx *gorm.DB) error {
	if mat.ID == "" {
		mat.ID = uuid.New().String()
	}
	return nil
}

func (ma *MarketingAsset) BeforeCreate(tx *gorm.DB) error {
	if ma.ID == "" {
		ma.ID = uuid.New().String()
	}
	return nil
}

func (ct *CommunicationType) BeforeCreate(tx *gorm.DB) error {
	if ct.ID == "" {
		ct.ID = uuid.New().String()
	}
	return nil
}

func (tt *TaskType) BeforeCreate(tx *gorm.DB) error {
	if tt.ID == "" {
		tt.ID = uuid.New().String()
	}
	return nil
}

func (tt *TerritoryType) BeforeCreate(tx *gorm.DB) error {
	if tt.ID == "" {
		tt.ID = uuid.New().String()
	}
	return nil
}

func (t *Territory) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	return nil
}

func (cf *CustomField) BeforeCreate(tx *gorm.DB) error {
	if cf.ID == "" {
		cf.ID = uuid.New().String()
	}
	return nil
}

func (cfv *CustomFieldValue) BeforeCreate(tx *gorm.DB) error {
	if cfv.ID == "" {
		cfv.ID = uuid.New().String()
	}
	return nil
}

func (co *CustomObject) BeforeCreate(tx *gorm.DB) error {
	if co.ID == "" {
		co.ID = uuid.New().String()
	}
	return nil
}

func (al *ActivityLog) BeforeCreate(tx *gorm.DB) error {
	if al.ID == "" {
		al.ID = uuid.New().String()
	}
	return nil
}

func (dsh *DealStageHistory) BeforeCreate(tx *gorm.DB) error {
	if dsh.ID == "" {
		dsh.ID = uuid.New().String()
	}
	return nil
}

func (pnt *PhoneNumberType) BeforeCreate(tx *gorm.DB) error {
	if pnt.ID == "" {
		pnt.ID = uuid.New().String()
	}
	return nil
}

func (pn *PhoneNumber) BeforeCreate(tx *gorm.DB) error {
	if pn.ID == "" {
		pn.ID = uuid.New().String()
	}
	return nil
}

func (eat *EmailAddressType) BeforeCreate(tx *gorm.DB) error {
	if eat.ID == "" {
		eat.ID = uuid.New().String()
	}
	return nil
}

func (ea *EmailAddress) BeforeCreate(tx *gorm.DB) error {
	if ea.ID == "" {
		ea.ID = uuid.New().String()
	}
	return nil
}

func (at *AddressType) BeforeCreate(tx *gorm.DB) error {
	if at.ID == "" {
		at.ID = uuid.New().String()
	}
	return nil
}

func (a *Address) BeforeCreate(tx *gorm.DB) error {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	return nil
}

func (smt *SocialMediaType) BeforeCreate(tx *gorm.DB) error {
	if smt.ID == "" {
		smt.ID = uuid.New().String()
	}
	return nil
}

func (sma *SocialMediaAccount) BeforeCreate(tx *gorm.DB) error {
	if sma.ID == "" {
		sma.ID = uuid.New().String()
	}
	return nil
}

func (ca *CommunicationAttachment) BeforeCreate(tx *gorm.DB) error {
	if ca.ID == "" {
		ca.ID = uuid.New().String()
	}
	return nil
}

// ============================================================================
// TERRITORIES AND EXTENSIBILITY
// ============================================================================

type TerritoryType struct {
	ID          string  `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Name        string  `json:"name" gorm:"not null"`
	Code        string  `json:"code" gorm:"not null"`
	Description *string `json:"description"`
	IsActive    bool    `json:"isActive" gorm:"column:is_active;default:true"`

	TenantID string `json:"tenantId" gorm:"column:tenant_id;type:uuid;not null"`
	Tenant   Tenant `json:"tenant,omitempty" gorm:"foreignKey:TenantID"`

	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`
}

type Territory struct {
	ID          string   `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Name        string   `json:"name" gorm:"not null"`
	TypeID      *string  `json:"typeId" gorm:"column:type_id;type:uuid"`
	Countries   []string `json:"countries" gorm:"type:jsonb"`
	States      []string `json:"states" gorm:"type:jsonb"`
	Cities      []string `json:"cities" gorm:"type:jsonb"`
	PostalCodes []string `json:"postalCodes" gorm:"column:postal_codes;type:jsonb"`
	Industries  []string `json:"industries" gorm:"type:jsonb"`
	CompanySize []string `json:"companySize" gorm:"column:company_size;type:jsonb"`

	TenantID string `json:"tenantId" gorm:"column:tenant_id;type:uuid;not null"`
	Tenant   Tenant `json:"tenant,omitempty" gorm:"foreignKey:TenantID"`

	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`
}

type CustomField struct {
	ID           string      `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Name         string      `json:"name" gorm:"not null"`
	Label        string      `json:"label" gorm:"not null"`
	Type         string      `json:"type" gorm:"not null"`
	EntityType   string      `json:"entityType" gorm:"column:entity_type;not null"`
	IsRequired   bool        `json:"isRequired" gorm:"column:is_required;default:false"`
	IsUnique     bool        `json:"isUnique" gorm:"column:is_unique;default:false"`
	DefaultValue *string     `json:"defaultValue" gorm:"column:default_value"`
	Options      interface{} `json:"options" gorm:"type:jsonb"`
	LookupEntity *string     `json:"lookupEntity" gorm:"column:lookup_entity"`
	Validation   interface{} `json:"validation" gorm:"type:jsonb"`

	TenantID string `json:"tenantId" gorm:"column:tenant_id;type:uuid;not null"`
	Tenant   Tenant `json:"tenant,omitempty" gorm:"foreignKey:TenantID"`

	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`
}

type CustomFieldValue struct {
	ID           string      `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	FieldID      string      `json:"fieldId" gorm:"column:field_id;type:uuid;not null"`
	EntityID     string      `json:"entityId" gorm:"column:entity_id;not null"`
	EntityType   string      `json:"entityType" gorm:"column:entity_type;not null"`
	TextValue    *string     `json:"textValue" gorm:"column:text_value"`
	NumberValue  *int        `json:"numberValue" gorm:"column:number_value"`
	DecimalValue *float64    `json:"decimalValue" gorm:"column:decimal_value"`
	BooleanValue *bool       `json:"booleanValue" gorm:"column:boolean_value"`
	DateValue    *time.Time  `json:"dateValue" gorm:"column:date_value"`
	JSONValue    interface{} `json:"jsonValue" gorm:"column:json_value;type:jsonb"`

	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`
}

type CustomObject struct {
	ID          string `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Name        string `json:"name" gorm:"not null"`
	Label       string `json:"label" gorm:"not null"`
	PluralLabel string `json:"pluralLabel" gorm:"column:plural_label;not null"`

	TenantID string `json:"tenantId" gorm:"column:tenant_id;type:uuid;not null"`
	Tenant   Tenant `json:"tenant,omitempty" gorm:"foreignKey:TenantID"`

	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`
}

type ActivityLog struct {
	ID         string  `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	EntityType string  `json:"entityType" gorm:"column:entity_type;not null"`
	EntityID   string  `json:"entityId" gorm:"column:entity_id;not null"`
	Action     string  `json:"action" gorm:"not null"`
	FieldName  *string `json:"fieldName" gorm:"column:field_name"`
	OldValue   *string `json:"oldValue" gorm:"column:old_value"`
	NewValue   *string `json:"newValue" gorm:"column:new_value"`
	UserID     *string `json:"userId" gorm:"column:user_id;type:uuid"`
	IPAddress  *string `json:"ipAddress" gorm:"column:ip_address"`
	UserAgent  *string `json:"userAgent" gorm:"column:user_agent"`

	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
}

type DealStageHistory struct {
	ID                    string     `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	DealID                string     `json:"dealId" gorm:"column:deal_id;type:uuid;not null"`
	FromStageID           *string    `json:"fromStageId" gorm:"column:from_stage_id;type:uuid"`
	ToStageID             string     `json:"toStageId" gorm:"column:to_stage_id;type:uuid;not null"`
	FromAmount            *float64   `json:"fromAmount" gorm:"column:from_amount"`
	ToAmount              *float64   `json:"toAmount" gorm:"column:to_amount"`
	FromProbability       *int       `json:"fromProbability" gorm:"column:from_probability"`
	ToProbability         *int       `json:"toProbability" gorm:"column:to_probability"`
	FromCurrency          *string    `json:"fromCurrency" gorm:"column:from_currency"`
	ToCurrency            *string    `json:"toCurrency" gorm:"column:to_currency"`
	FromExpectedCloseDate *time.Time `json:"fromExpectedCloseDate" gorm:"column:from_expected_close_date"`
	ToExpectedCloseDate   *time.Time `json:"toExpectedCloseDate" gorm:"column:to_expected_close_date"`
	ChangeReason          *string    `json:"changeReason" gorm:"column:change_reason"`
	Notes                 *string    `json:"notes"`
	MovedAt               time.Time  `json:"movedAt" gorm:"column:moved_at;default:CURRENT_TIMESTAMP"`
	MovedBy               *string    `json:"movedBy" gorm:"column:moved_by;type:uuid"`
}

// ============================================================================
// CONTACT INFORMATION
// ============================================================================

type PhoneNumberType struct {
	ID          string  `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Name        string  `json:"name" gorm:"not null"`
	Code        string  `json:"code" gorm:"not null"`
	Description *string `json:"description"`
	IsActive    bool    `json:"isActive" gorm:"column:is_active;default:true"`

	TenantID string `json:"tenantId" gorm:"column:tenant_id;type:uuid;not null"`
	Tenant   Tenant `json:"tenant,omitempty" gorm:"foreignKey:TenantID"`

	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`
}

type PhoneNumber struct {
	ID         string  `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Number     string  `json:"number" gorm:"not null"`
	Extension  *string `json:"extension"`
	IsPrimary  bool    `json:"isPrimary" gorm:"column:is_primary;default:false"`
	TypeID     *string `json:"typeId" gorm:"column:type_id;type:uuid"`
	EntityID   string  `json:"entityId" gorm:"column:entity_id;type:uuid;not null"`
	EntityType string  `json:"entityType" gorm:"column:entity_type;not null"`

	TenantID string `json:"tenantId" gorm:"column:tenant_id;type:uuid;not null"`
	Tenant   Tenant `json:"tenant,omitempty" gorm:"foreignKey:TenantID"`

	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`
}

type EmailAddressType struct {
	ID          string  `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Name        string  `json:"name" gorm:"not null"`
	Code        string  `json:"code" gorm:"not null"`
	Description *string `json:"description"`
	IsActive    bool    `json:"isActive" gorm:"column:is_active;default:true"`

	TenantID string `json:"tenantId" gorm:"column:tenant_id;type:uuid;not null"`
	Tenant   Tenant `json:"tenant,omitempty" gorm:"foreignKey:TenantID"`

	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`
}

type EmailAddress struct {
	ID         string  `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Email      string  `json:"email" gorm:"not null"`
	IsPrimary  bool    `json:"isPrimary" gorm:"column:is_primary;default:false"`
	IsVerified bool    `json:"isVerified" gorm:"column:is_verified;default:false"`
	TypeID     *string `json:"typeId" gorm:"column:type_id;type:uuid"`
	EntityID   string  `json:"entityId" gorm:"column:entity_id;type:uuid;not null"`
	EntityType string  `json:"entityType" gorm:"column:entity_type;not null"`

	TenantID string `json:"tenantId" gorm:"column:tenant_id;type:uuid;not null"`
	Tenant   Tenant `json:"tenant,omitempty" gorm:"foreignKey:TenantID"`

	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`
}

type AddressType struct {
	ID          string  `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Name        string  `json:"name" gorm:"not null"`
	Code        string  `json:"code" gorm:"not null"`
	Description *string `json:"description"`
	IsActive    bool    `json:"isActive" gorm:"column:is_active;default:true"`

	TenantID string `json:"tenantId" gorm:"column:tenant_id;type:uuid;not null"`
	Tenant   Tenant `json:"tenant,omitempty" gorm:"foreignKey:TenantID"`

	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`
}

type Address struct {
	ID         string  `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Street1    *string `json:"street1"`
	Street2    *string `json:"street2"`
	City       *string `json:"city"`
	State      *string `json:"state"`
	PostalCode *string `json:"postalCode" gorm:"column:postal_code"`
	Country    *string `json:"country"`
	IsPrimary  bool    `json:"isPrimary" gorm:"column:is_primary;default:false"`
	TypeID     *string `json:"typeId" gorm:"column:type_id;type:uuid"`
	EntityID   string  `json:"entityId" gorm:"column:entity_id;type:uuid;not null"`
	EntityType string  `json:"entityType" gorm:"column:entity_type;not null"`

	TenantID string `json:"tenantId" gorm:"column:tenant_id;type:uuid;not null"`
	Tenant   Tenant `json:"tenant,omitempty" gorm:"foreignKey:TenantID"`

	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`
}

type SocialMediaType struct {
	ID       string  `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Name     string  `json:"name" gorm:"not null"`
	Code     string  `json:"code" gorm:"not null"`
	Icon     *string `json:"icon"`
	BaseURL  *string `json:"baseUrl" gorm:"column:base_url"`
	IsActive bool    `json:"isActive" gorm:"column:is_active;default:true"`

	TenantID string `json:"tenantId" gorm:"column:tenant_id;type:uuid;not null"`
	Tenant   Tenant `json:"tenant,omitempty" gorm:"foreignKey:TenantID"`

	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`
}

type SocialMediaAccount struct {
	ID         string  `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Username   *string `json:"username"`
	URL        *string `json:"url"`
	IsPrimary  bool    `json:"isPrimary" gorm:"column:is_primary;default:false"`
	TypeID     *string `json:"typeId" gorm:"column:type_id;type:uuid"`
	EntityID   string  `json:"entityId" gorm:"column:entity_id;type:uuid;not null"`
	EntityType string  `json:"entityType" gorm:"column:entity_type;not null"`

	TenantID string `json:"tenantId" gorm:"column:tenant_id;type:uuid;not null"`
	Tenant   Tenant `json:"tenant,omitempty" gorm:"foreignKey:TenantID"`

	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`
}

type CommunicationAttachment struct {
	ID              string `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Filename        string `json:"filename" gorm:"not null"`
	OriginalName    string `json:"originalName" gorm:"column:original_name;not null"`
	MimeType        string `json:"mimeType" gorm:"column:mime_type;not null"`
	Size            int    `json:"size" gorm:"not null"`
	URL             string `json:"url" gorm:"not null"`
	CommunicationID string `json:"communicationId" gorm:"column:communication_id;type:uuid;not null"`

	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
}
