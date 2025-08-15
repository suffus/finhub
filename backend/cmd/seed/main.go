package main

import (
	"fmt"
	"log"
	"strings"

	"finhub-backend/config"
	"finhub-backend/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := initDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Create a default tenant if it doesn't exist
	var tenant models.Tenant
	if err := db.Where("subdomain = ?", "default").First(&tenant).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			tenant = models.Tenant{
				Name:      "Default Tenant",
				Subdomain: "default",
				IsActive:  true,
			}
			if err := db.Create(&tenant).Error; err != nil {
				log.Fatal("Failed to create default tenant:", err)
			}
			log.Println("Created default tenant with ID:", tenant.ID)
		} else {
			log.Fatal("Failed to check for tenant:", err)
		}
	} else {
		log.Println("Using existing tenant with ID:", tenant.ID)
	}

	// Seed industries
	industries := []models.Industry{
		{
			Name:        "Technology",
			Code:        "TECH",
			Description: stringPtr("Software, hardware, and IT services"),
			IsActive:    true,
			TenantID:    tenant.ID,
		},
		{
			Name:        "Healthcare",
			Code:        "HEALTH",
			Description: stringPtr("Medical devices, pharmaceuticals, and healthcare services"),
			IsActive:    true,
			TenantID:    tenant.ID,
		},
		{
			Name:        "Finance",
			Code:        "FINANCE",
			Description: stringPtr("Banking, insurance, and financial services"),
			IsActive:    true,
			TenantID:    tenant.ID,
		},
		{
			Name:        "Manufacturing",
			Code:        "MFG",
			Description: stringPtr("Industrial manufacturing and production"),
			IsActive:    true,
			TenantID:    tenant.ID,
		},
		{
			Name:        "Retail",
			Code:        "RETAIL",
			Description: stringPtr("Consumer goods and retail services"),
			IsActive:    true,
			TenantID:    tenant.ID,
		},
		{
			Name:        "Education",
			Code:        "EDU",
			Description: stringPtr("Educational institutions and training services"),
			IsActive:    true,
			TenantID:    tenant.ID,
		},
		{
			Name:        "Real Estate",
			Code:        "REAL_ESTATE",
			Description: stringPtr("Property development and real estate services"),
			IsActive:    true,
			TenantID:    tenant.ID,
		},
		{
			Name:        "Consulting",
			Code:        "CONSULTING",
			Description: stringPtr("Business consulting and advisory services"),
			IsActive:    true,
			TenantID:    tenant.ID,
		},
		{
			Name:        "Media & Entertainment",
			Code:        "MEDIA",
			Description: stringPtr("Content creation, publishing, and entertainment"),
			IsActive:    true,
			TenantID:    tenant.ID,
		},
		{
			Name:        "Transportation & Logistics",
			Code:        "TRANSPORT",
			Description: stringPtr("Shipping, logistics, and transportation services"),
			IsActive:    true,
			TenantID:    tenant.ID,
		},
	}

	for _, industry := range industries {
		var existing models.Industry
		if err := db.Where("code = ? AND tenant_id = ?", industry.Code, tenant.ID).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&industry).Error; err != nil {
					log.Printf("Failed to create industry %s: %v", industry.Name, err)
				} else {
					log.Printf("Created industry: %s", industry.Name)
				}
			} else {
				log.Printf("Error checking industry %s: %v", industry.Name, err)
			}
		} else {
			log.Printf("Industry %s already exists", industry.Name)
		}
	}

	// Seed company sizes
	companySizes := []models.CompanySize{
		{
			Name:         "Startup (1-10)",
			Code:         "STARTUP",
			Description:  stringPtr("Early-stage companies with 1-10 employees"),
			MinEmployees: intPtr(1),
			MaxEmployees: intPtr(10),
			IsActive:     true,
			TenantID:     tenant.ID,
		},
		{
			Name:         "Small Business (11-50)",
			Code:         "SMALL",
			Description:  stringPtr("Small businesses with 11-50 employees"),
			MinEmployees: intPtr(11),
			MaxEmployees: intPtr(50),
			IsActive:     true,
			TenantID:     tenant.ID,
		},
		{
			Name:         "Medium Business (51-200)",
			Code:         "MEDIUM",
			Description:  stringPtr("Medium-sized businesses with 51-200 employees"),
			MinEmployees: intPtr(51),
			MaxEmployees: intPtr(200),
			IsActive:     true,
			TenantID:     tenant.ID,
		},
		{
			Name:         "Large Business (201-1000)",
			Code:         "LARGE",
			Description:  stringPtr("Large businesses with 201-1000 employees"),
			MinEmployees: intPtr(201),
			MaxEmployees: intPtr(1000),
			IsActive:     true,
			TenantID:     tenant.ID,
		},
		{
			Name:         "Enterprise (1000+)",
			Code:         "ENTERPRISE",
			Description:  stringPtr("Enterprise companies with 1000+ employees"),
			MinEmployees: intPtr(1000),
			MaxEmployees: nil,
			IsActive:     true,
			TenantID:     tenant.ID,
		},
	}

	for _, size := range companySizes {
		var existing models.CompanySize
		if err := db.Where("code = ? AND tenant_id = ?", size.Code, tenant.ID).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&size).Error; err != nil {
					log.Printf("Failed to create company size %s: %v", size.Name, err)
				} else {
					log.Printf("Created company size: %s", size.Name)
				}
			} else {
				log.Printf("Error checking company size %s: %v", size.Name, err)
			}
		} else {
			log.Printf("Company size %s already exists", size.Name)
		}
	}

	// Seed lead statuses
	leadStatuses := []models.LeadStatus{
		{
			Name:        "New",
			Code:        "NEW",
			Description: stringPtr("Newly created lead"),
			Color:       stringPtr("#3B82F6"),
			Order:       1,
			IsActive:    true,
			TenantID:    tenant.ID,
		},
		{
			Name:        "Contacted",
			Code:        "CONTACTED",
			Description: stringPtr("Initial contact made"),
			Color:       stringPtr("#10B981"),
			Order:       2,
			IsActive:    true,
			TenantID:    tenant.ID,
		},
		{
			Name:        "Qualified",
			Code:        "QUALIFIED",
			Description: stringPtr("Lead has been qualified"),
			Color:       stringPtr("#F59E0B"),
			Order:       3,
			IsActive:    true,
			TenantID:    tenant.ID,
		},
		{
			Name:        "Proposal",
			Code:        "PROPOSAL",
			Description: stringPtr("Proposal sent to lead"),
			Color:       stringPtr("#8B5CF6"),
			Order:       4,
			IsActive:    true,
			TenantID:    tenant.ID,
		},
		{
			Name:        "Negotiation",
			Code:        "NEGOTIATION",
			Description: stringPtr("In negotiation phase"),
			Color:       stringPtr("#EF4444"),
			Order:       5,
			IsActive:    true,
			TenantID:    tenant.ID,
		},
		{
			Name:        "Converted",
			Code:        "CONVERTED",
			Description: stringPtr("Successfully converted to customer"),
			Color:       stringPtr("#059669"),
			Order:       6,
			IsActive:    true,
			TenantID:    tenant.ID,
		},
		{
			Name:        "Lost",
			Code:        "LOST",
			Description: stringPtr("Lead was lost"),
			Color:       stringPtr("#6B7280"),
			Order:       7,
			IsActive:    true,
			TenantID:    tenant.ID,
		},
	}

	for _, status := range leadStatuses {
		var existing models.LeadStatus
		if err := db.Where("code = ? AND tenant_id = ?", status.Code, tenant.ID).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&status).Error; err != nil {
					log.Printf("Failed to create lead status %s: %v", status.Name, err)
				} else {
					log.Printf("Created lead status: %s", status.Name)
				}
			} else {
				log.Printf("Error checking lead status %s: %v", status.Name, err)
			}
		} else {
			log.Printf("Lead status %s already exists", status.Name)
		}
	}

	// Seed lead temperatures
	leadTemperatures := []models.LeadTemperature{
		{
			Name:        "Hot",
			Code:        "HOT",
			Description: stringPtr("High probability of conversion"),
			Color:       stringPtr("#EF4444"),
			Order:       1,
			IsActive:    true,
			TenantID:    tenant.ID,
		},
		{
			Name:        "Warm",
			Code:        "WARM",
			Description: stringPtr("Medium probability of conversion"),
			Color:       stringPtr("#F59E0B"),
			Order:       2,
			IsActive:    true,
			TenantID:    tenant.ID,
		},
		{
			Name:        "Cold",
			Code:        "COLD",
			Description: stringPtr("Low probability of conversion"),
			Color:       stringPtr("#3B82F6"),
			Order:       3,
			IsActive:    true,
			TenantID:    tenant.ID,
		},
	}

	// Seed pipelines
	pipelines := []models.Pipeline{
		{
			Name:     "Sales Pipeline",
			IsActive: true,
			TenantID: tenant.ID,
		},
		{
			Name:     "Lead Qualification",
			IsActive: true,
			TenantID: tenant.ID,
		},
	}

	for _, pipeline := range pipelines {
		var existing models.Pipeline
		if err := db.Where("name = ? AND tenant_id = ?", pipeline.Name, tenant.ID).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&pipeline).Error; err != nil {
					log.Printf("Failed to create pipeline %s: %v", pipeline.Name, err)
				} else {
					log.Printf("Created pipeline: %s", pipeline.Name)
				}
			} else {
				log.Printf("Error checking pipeline %s: %v", pipeline.Name, err)
			}
		} else {
			log.Printf("Pipeline %s already exists", pipeline.Name)
		}
	}

	// Seed stages for the sales pipeline
	var salesPipeline models.Pipeline
	if err := db.Where("name = ? AND tenant_id = ?", "Sales Pipeline", tenant.ID).First(&salesPipeline).Error; err == nil {
		stages := []models.Stage{
			{
				Name:         "Prospecting",
				Order:        1,
				Probability:  10,
				IsClosedWon:  false,
				IsClosedLost: false,
				Color:        stringPtr("#3B82F6"),
				PipelineID:   salesPipeline.ID,
				TenantID:     tenant.ID,
			},
			{
				Name:         "Qualification",
				Order:        2,
				Probability:  25,
				IsClosedWon:  false,
				IsClosedLost: false,
				Color:        stringPtr("#10B981"),
				PipelineID:   salesPipeline.ID,
				TenantID:     tenant.ID,
			},
			{
				Name:         "Proposal",
				Order:        3,
				Probability:  50,
				IsClosedWon:  false,
				IsClosedLost: false,
				Color:        stringPtr("#F59E0B"),
				PipelineID:   salesPipeline.ID,
				TenantID:     tenant.ID,
			},
			{
				Name:         "Negotiation",
				Order:        4,
				Probability:  75,
				IsClosedWon:  false,
				IsClosedLost: false,
				Color:        stringPtr("#8B5CF6"),
				PipelineID:   salesPipeline.ID,
				TenantID:     tenant.ID,
			},
			{
				Name:         "Closed Won",
				Order:        5,
				Probability:  100,
				IsClosedWon:  true,
				IsClosedLost: false,
				Color:        stringPtr("#059669"),
				PipelineID:   salesPipeline.ID,
				TenantID:     tenant.ID,
			},
			{
				Name:         "Closed Lost",
				Order:        6,
				Probability:  0,
				IsClosedWon:  false,
				IsClosedLost: true,
				Color:        stringPtr("#6B7280"),
				PipelineID:   salesPipeline.ID,
				TenantID:     tenant.ID,
			},
		}

		for _, stage := range stages {
			var existing models.Stage
			if err := db.Where("name = ? AND pipeline_id = ? AND tenant_id = ?", stage.Name, stage.PipelineID, tenant.ID).First(&existing).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					if err := db.Create(&stage).Error; err != nil {
						log.Printf("Failed to create stage %s: %v", stage.Name, err)
					} else {
						log.Printf("Created stage: %s", stage.Name)
					}
				} else {
					log.Printf("Error checking stage %s: %v", stage.Name, err)
				}
			} else {
				log.Printf("Stage %s already exists", stage.Name)
			}
		}
	}

	for _, temp := range leadTemperatures {
		var existing models.LeadTemperature
		if err := db.Where("code = ? AND tenant_id = ?", temp.Code, tenant.ID).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&temp).Error; err != nil {
					log.Printf("Failed to create lead temperature %s: %v", temp.Name, err)
				} else {
					log.Printf("Created lead temperature: %s", temp.Name)
				}
			} else {
				log.Printf("Error checking lead temperature %s: %v", temp.Name, err)
			}
		} else {
			log.Printf("Lead temperature %s already exists", temp.Name)
		}
	}

	// Seed marketing source types
	marketingSourceTypes := []models.MarketingSourceType{
		{
			Name:        "Digital Advertising",
			Code:        "DIGITAL_AD",
			Description: stringPtr("Online advertising including PPC, display ads, and social media ads"),
			Color:       stringPtr("#3B82F6"),
			IsActive:    true,
			TenantID:    tenant.ID,
		},
		{
			Name:        "Content Marketing",
			Code:        "CONTENT",
			Description: stringPtr("Blog posts, whitepapers, ebooks, and other content"),
			Color:       stringPtr("#10B981"),
			IsActive:    true,
			TenantID:    tenant.ID,
		},
		{
			Name:        "Email Marketing",
			Code:        "EMAIL",
			Description: stringPtr("Email campaigns and newsletters"),
			Color:       stringPtr("#F59E0B"),
			IsActive:    true,
			TenantID:    tenant.ID,
		},
		{
			Name:        "Social Media",
			Code:        "SOCIAL",
			Description: stringPtr("Social media marketing and engagement"),
			Color:       stringPtr("#8B5CF6"),
			IsActive:    true,
			TenantID:    tenant.ID,
		},
		{
			Name:        "Events",
			Code:        "EVENTS",
			Description: stringPtr("Trade shows, conferences, and webinars"),
			Color:       stringPtr("#EF4444"),
			IsActive:    true,
			TenantID:    tenant.ID,
		},
		{
			Name:        "Referral",
			Code:        "REFERRAL",
			Description: stringPtr("Customer referrals and word-of-mouth"),
			Color:       stringPtr("#059669"),
			IsActive:    true,
			TenantID:    tenant.ID,
		},
	}

	for _, mst := range marketingSourceTypes {
		var existing models.MarketingSourceType
		if err := db.Where("code = ? AND tenant_id = ?", mst.Code, tenant.ID).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&mst).Error; err != nil {
					log.Printf("Failed to create marketing source type %s: %v", mst.Name, err)
				} else {
					log.Printf("Created marketing source type: %s", mst.Name)
				}
			} else {
				log.Printf("Error checking marketing source type %s: %v", mst.Name, err)
			}
		} else {
			log.Printf("Marketing source type %s already exists", mst.Name)
		}
	}

	// Seed marketing sources
	marketingSources := []models.MarketingSource{
		{
			Name:     "Google Ads",
			TypeID:   getTypeID(db, "DIGITAL_AD", "marketing_source_types", tenant.ID),
			Medium:   stringPtr("PPC"),
			Campaign: stringPtr("Brand Awareness"),
			Source:   stringPtr("Google"),
			Content:  stringPtr("Search ads"),
			Term:     stringPtr("CRM software"),
			TenantID: tenant.ID,
		},
		{
			Name:     "LinkedIn Ads",
			TypeID:   getTypeID(db, "DIGITAL_AD", "marketing_source_types", tenant.ID),
			Medium:   stringPtr("Social"),
			Campaign: stringPtr("B2B Lead Generation"),
			Source:   stringPtr("LinkedIn"),
			Content:  stringPtr("Sponsored content"),
			TenantID: tenant.ID,
		},
		{
			Name:     "Company Blog",
			TypeID:   getTypeID(db, "CONTENT", "marketing_source_types", tenant.ID),
			Medium:   stringPtr("Content"),
			Campaign: stringPtr("Thought Leadership"),
			Source:   stringPtr("Organic"),
			Content:  stringPtr("Blog posts"),
			TenantID: tenant.ID,
		},
		{
			Name:     "Trade Show",
			TypeID:   getTypeID(db, "EVENTS", "marketing_source_types", tenant.ID),
			Medium:   stringPtr("Events"),
			Campaign: stringPtr("Industry Presence"),
			Source:   stringPtr("Direct"),
			Content:  stringPtr("Booth presence"),
			TenantID: tenant.ID,
		},
		{
			Name:     "Website",
			TypeID:   getTypeID(db, "EVENTS", "marketing_source_types", tenant.ID),
			Medium:   stringPtr("Events"),
			Campaign: stringPtr("Industry Presence"),
			Source:   stringPtr("Direct"),
			Content:  stringPtr("Booth presence"),
			TenantID: tenant.ID,
		},
	}

	for _, ms := range marketingSources {
		var existing models.MarketingSource
		if err := db.Where("name = ? AND tenant_id = ?", ms.Name, tenant.ID).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&ms).Error; err != nil {
					log.Printf("Failed to create marketing source %s: %v", ms.Name, err)
				} else {
					log.Printf("Created marketing source: %s", ms.Name)
				}
			} else {
				log.Printf("Error checking marketing source %s: %v", ms.Name, err)
			}
		} else {
			log.Printf("Marketing source %s already exists", ms.Name)
		}
	}

	// Seed marketing asset types
	marketingAssetTypes := []models.MarketingAssetType{
		{
			Name:        "Whitepaper",
			Code:        "WHITEPAPER",
			Description: stringPtr("In-depth technical or business documents"),
			Color:       stringPtr("#3B82F6"),
			IsActive:    true,
			TenantID:    tenant.ID,
		},
		{
			Name:        "Case Study",
			Code:        "CASE_STUDY",
			Description: stringPtr("Customer success stories and results"),
			Color:       stringPtr("#10B981"),
			IsActive:    true,
			TenantID:    tenant.ID,
		},
		{
			Name:        "Video",
			Code:        "VIDEO",
			Description: stringPtr("Product demos, testimonials, and explainer videos"),
			Color:       stringPtr("#F59E0B"),
			IsActive:    true,
			TenantID:    tenant.ID,
		},
		{
			Name:        "Infographic",
			Code:        "INFOGRAPHIC",
			Description: stringPtr("Visual representations of data and concepts"),
			Color:       stringPtr("#8B5CF6"),
			IsActive:    true,
			TenantID:    tenant.ID,
		},
		{
			Name:        "Webinar",
			Code:        "WEBINAR",
			Description: stringPtr("Online presentations and training sessions"),
			Color:       stringPtr("#EF4444"),
			IsActive:    true,
			TenantID:    tenant.ID,
		},
	}

	for _, mat := range marketingAssetTypes {
		var existing models.MarketingAssetType
		if err := db.Where("code = ? AND tenant_id = ?", mat.Code, tenant.ID).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&mat).Error; err != nil {
					log.Printf("Failed to create marketing asset type %s: %v", mat.Name, err)
				} else {
					log.Printf("Created marketing asset type: %s", mat.Name)
				}
			} else {
				log.Printf("Error checking marketing asset type %s: %v", mat.Name, err)
			}
		} else {
			log.Printf("Marketing asset type %s already exists", mat.Name)
		}
	}

	// Seed communication types
	communicationTypes := []models.CommunicationType{
		{
			Name:        "Email",
			Code:        "EMAIL",
			Description: stringPtr("Email communications"),
			Icon:        stringPtr("mail"),
			IsActive:    true,
			TenantID:    tenant.ID,
		},
		{
			Name:        "Phone Call",
			Code:        "PHONE",
			Description: stringPtr("Telephone conversations"),
			Icon:        stringPtr("phone"),
			IsActive:    true,
			TenantID:    tenant.ID,
		},
		{
			Name:        "Meeting",
			Code:        "MEETING",
			Description: stringPtr("In-person or virtual meetings"),
			Icon:        stringPtr("users"),
			IsActive:    true,
			TenantID:    tenant.ID,
		},
		{
			Name:        "Text Message",
			Code:        "SMS",
			Description: stringPtr("SMS and messaging app communications"),
			Icon:        stringPtr("message-circle"),
			IsActive:    true,
			TenantID:    tenant.ID,
		},
		{
			Name:        "LinkedIn Message",
			Code:        "LINKEDIN",
			Description: stringPtr("LinkedIn messaging"),
			Icon:        stringPtr("linkedin"),
			IsActive:    true,
			TenantID:    tenant.ID,
		},
	}

	for _, ct := range communicationTypes {
		var existing models.CommunicationType
		if err := db.Where("code = ? AND tenant_id = ?", ct.Code, tenant.ID).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&ct).Error; err != nil {
					log.Printf("Failed to create communication type %s: %v", ct.Name, err)
				} else {
					log.Printf("Created communication type: %s", ct.Name)
				}
			} else {
				log.Printf("Error checking communication type %s: %v", ct.Name, err)
			}
		} else {
			log.Printf("Communication type %s already exists", ct.Name)
		}
	}

	// Seed task types
	taskTypes := []models.TaskType{
		{
			Name:        "Follow Up",
			Code:        "FOLLOW_UP",
			Description: stringPtr("Follow-up tasks and reminders"),
			Color:       stringPtr("#3B82F6"),
			Icon:        stringPtr("refresh-cw"),
			IsActive:    true,
			TenantID:    tenant.ID,
		},
		{
			Name:        "Proposal",
			Code:        "PROPOSAL",
			Description: stringPtr("Creating and sending proposals"),
			Color:       stringPtr("#10B981"),
			Icon:        stringPtr("file-text"),
			IsActive:    true,
			TenantID:    tenant.ID,
		},
		{
			Name:        "Demo",
			Code:        "DEMO",
			Description: stringPtr("Product demonstrations"),
			Color:       stringPtr("#F59E0B"),
			Icon:        stringPtr("play"),
			IsActive:    true,
			TenantID:    tenant.ID,
		},
		{
			Name:        "Contract Review",
			Code:        "CONTRACT",
			Description: stringPtr("Contract review and negotiation"),
			Color:       stringPtr("#8B5CF6"),
			Icon:        stringPtr("file-check"),
			IsActive:    true,
			TenantID:    tenant.ID,
		},
		{
			Name:        "Research",
			Code:        "RESEARCH",
			Description: stringPtr("Market and competitive research"),
			Color:       stringPtr("#EF4444"),
			Icon:        stringPtr("search"),
			IsActive:    true,
			TenantID:    tenant.ID,
		},
	}

	for _, tt := range taskTypes {
		var existing models.TaskType
		if err := db.Where("code = ? AND tenant_id = ?", tt.Code, tenant.ID).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&tt).Error; err != nil {
					log.Printf("Failed to create task type %s: %v", tt.Name, err)
				} else {
					log.Printf("Created task type: %s", tt.Name)
				}
			} else {
				log.Printf("Error checking task type %s: %v", tt.Name, err)
			}
		} else {
			log.Printf("Task type %s already exists", tt.Name)
		}
	}

	// Seed address types
	addressTypes := []models.AddressType{
		{
			Name:        "Billing",
			Code:        "BILLING",
			Description: stringPtr("Billing address"),
			IsActive:    true,
			TenantID:    tenant.ID,
		},
		{
			Name:        "Shipping",
			Code:        "SHIPPING",
			Description: stringPtr("Shipping address"),
			IsActive:    true,
			TenantID:    tenant.ID,
		},
		{
			Name:        "Office",
			Code:        "OFFICE",
			Description: stringPtr("Office location"),
			IsActive:    true,
			TenantID:    tenant.ID,
		},
		{
			Name:        "Home",
			Code:        "HOME",
			Description: stringPtr("Home address"),
			IsActive:    true,
			TenantID:    tenant.ID,
		},
	}

	for _, at := range addressTypes {
		var existing models.AddressType
		if err := db.Where("code = ? AND tenant_id = ?", at.Code, tenant.ID).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&at).Error; err != nil {
					log.Printf("Failed to create address type %s: %v", at.Name, err)
				} else {
					log.Printf("Created address type: %s", at.Name)
				}
			} else {
				log.Printf("Error checking address type %s: %v", at.Name, err)
			}
		} else {
			log.Printf("Address type %s already exists", at.Name)
		}
	}

	// Seed email address types
	emailAddressTypes := []models.EmailAddressType{
		{
			Name:        "Work",
			Code:        "WORK",
			Description: stringPtr("Work email address"),
			IsActive:    true,
			TenantID:    tenant.ID,
		},
		{
			Name:        "Personal",
			Code:        "PERSONAL",
			Description: stringPtr("Personal email address"),
			IsActive:    true,
			TenantID:    tenant.ID,
		},
		{
			Name:        "Billing",
			Code:        "BILLING",
			Description: stringPtr("Billing email address"),
			IsActive:    true,
			TenantID:    tenant.ID,
		},
		{
			Name:        "Support",
			Code:        "SUPPORT",
			Description: stringPtr("Support email address"),
			IsActive:    true,
			TenantID:    tenant.ID,
		},
	}

	for _, eat := range emailAddressTypes {
		var existing models.EmailAddressType
		if err := db.Where("code = ? AND tenant_id = ?", eat.Code, tenant.ID).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&eat).Error; err != nil {
					log.Printf("Failed to create email address type %s: %v", eat.Name, err)
				} else {
					log.Printf("Created email address type: %s", eat.Name)
				}
			} else {
				log.Printf("Error checking email address type %s: %v", eat.Name, err)
			}
		} else {
			log.Printf("Email address type %s already exists", eat.Name)
		}
	}

	// Seed phone number types
	phoneNumberTypes := []models.PhoneNumberType{
		{
			Name:        "Mobile",
			Code:        "MOBILE",
			Description: stringPtr("Mobile phone number"),
			IsActive:    true,
			TenantID:    tenant.ID,
		},
		{
			Name:        "Work",
			Code:        "WORK",
			Description: stringPtr("Work phone number"),
			IsActive:    true,
			TenantID:    tenant.ID,
		},
		{
			Name:        "Home",
			Code:        "HOME",
			Description: stringPtr("Home phone number"),
			IsActive:    true,
			TenantID:    tenant.ID,
		},
		{
			Name:        "Fax",
			Code:        "FAX",
			Description: stringPtr("Fax number"),
			IsActive:    true,
			TenantID:    tenant.ID,
		},
	}

	for _, pnt := range phoneNumberTypes {
		var existing models.PhoneNumberType
		if err := db.Where("code = ? AND tenant_id = ?", pnt.Code, tenant.ID).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&pnt).Error; err != nil {
					log.Printf("Failed to create phone number type %s: %v", pnt.Name, err)
				} else {
					log.Printf("Created phone number type: %s", pnt.Name)
				}
			} else {
				log.Printf("Error checking phone number type %s: %v", pnt.Name, err)
			}
		} else {
			log.Printf("Phone number type %s already exists", pnt.Name)
		}
	}

	// Seed territory types
	territoryTypes := []models.TerritoryType{
		{
			Name:        "Geographic",
			Code:        "GEOGRAPHIC",
			Description: stringPtr("Geographic territories by region, state, or country"),
			IsActive:    true,
			TenantID:    tenant.ID,
		},
		{
			Name:        "Industry",
			Code:        "INDUSTRY",
			Description: stringPtr("Industry-specific territories"),
			IsActive:    true,
			TenantID:    tenant.ID,
		},
		{
			Name:        "Account Size",
			Code:        "ACCOUNT_SIZE",
			Description: stringPtr("Territories based on company size"),
			IsActive:    true,
			TenantID:    tenant.ID,
		},
		{
			Name:        "Product Line",
			Code:        "PRODUCT_LINE",
			Description: stringPtr("Territories based on product lines"),
			IsActive:    true,
			TenantID:    tenant.ID,
		},
	}

	for _, tt := range territoryTypes {
		var existing models.TerritoryType
		if err := db.Where("code = ? AND tenant_id = ?", tt.Code, tenant.ID).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&tt).Error; err != nil {
					log.Printf("Failed to create territory type %s: %v", tt.Name, err)
				} else {
					log.Printf("Created territory type: %s", tt.Name)
				}
			} else {
				log.Printf("Error checking territory type %s: %v", tt.Name, err)
			}
		} else {
			log.Printf("Territory type %s already exists", tt.Name)
		}
	}

	// Seed social media types
	socialMediaTypes := []models.SocialMediaType{
		{
			Name:     "LinkedIn",
			Code:     "LINKEDIN",
			Icon:     stringPtr("linkedin"),
			BaseURL:  stringPtr("https://linkedin.com/in/"),
			IsActive: true,
			TenantID: tenant.ID,
		},
		{
			Name:     "Twitter",
			Code:     "TWITTER",
			Icon:     stringPtr("twitter"),
			BaseURL:  stringPtr("https://twitter.com/"),
			IsActive: true,
			TenantID: tenant.ID,
		},
		{
			Name:     "Facebook",
			Code:     "FACEBOOK",
			Icon:     stringPtr("facebook"),
			BaseURL:  stringPtr("https://facebook.com/"),
			IsActive: true,
			TenantID: tenant.ID,
		},
		{
			Name:     "Instagram",
			Code:     "INSTAGRAM",
			Icon:     stringPtr("instagram"),
			BaseURL:  stringPtr("https://instagram.com/"),
			IsActive: true,
			TenantID: tenant.ID,
		},
		{
			Name:     "YouTube",
			Code:     "YOUTUBE",
			Icon:     stringPtr("youtube"),
			BaseURL:  stringPtr("https://youtube.com/"),
			IsActive: true,
			TenantID: tenant.ID,
		},
	}

	for _, smt := range socialMediaTypes {
		var existing models.SocialMediaType
		if err := db.Where("code = ? AND tenant_id = ?", smt.Code, tenant.ID).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&smt).Error; err != nil {
					log.Printf("Failed to create social media type %s: %v", smt.Name, err)
				} else {
					log.Printf("Created social media type: %s", smt.Name)
				}
			} else {
				log.Printf("Error checking social media type %s: %v", smt.Name, err)
			}
		} else {
			log.Printf("Social media type %s already exists", smt.Name)
		}
	}

	// Seed some sample territories (temporarily disabled due to JSONB issues)
	log.Println("Skipping territory seeding due to JSONB field handling complexity")
	/*
		territories := []models.Territory{
			{
				Name:        "North America",
				TypeID:      getTypeID(db, "GEOGRAPHIC", "territory_types", tenant.ID),
				Countries:   []string{"US", "CA", "MX"},
				States:      []string{"CA", "NY", "TX", "FL", "IL"},
				Cities:      []string{},
				PostalCodes: []string{},
				Industries:  []string{"TECH", "FINANCE", "HEALTHCARE"},
				CompanySize: []string{"LARGE", "ENTERPRISE"},
				TenantID:    tenant.ID,
			},
			{
				Name:        "Europe",
				TypeID:      getTypeID(db, "GEOGRAPHIC", "territory_types", tenant.ID),
				Countries:   []string{"GB", "DE", "FR", "IT", "ES"},
				States:      []string{},
				Cities:      []string{},
				PostalCodes: []string{},
				Industries:  []string{"MANUFACTURING", "FINANCE", "RETAIL"},
				CompanySize: []string{"MEDIUM", "LARGE"},
				TenantID:    tenant.ID,
			},
			{
				Name:        "Technology Sector",
				TypeID:      getTypeID(db, "INDUSTRY", "territory_types", tenant.ID),
				Countries:   []string{},
				States:      []string{},
				Cities:      []string{},
				PostalCodes: []string{},
				Industries:  []string{"TECH"},
				CompanySize: []string{"STARTUP", "SMALL", "MEDIUM", "LARGE", "ENTERPRISE"},
				TenantID:    tenant.ID,
			},
		}

		for _, territory := range territories {
			var existing models.Territory
			if err := db.Where("name = ? AND tenant_id = ?", territory.Name, tenant.ID).First(&existing).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					if err := createTerritoryWithJSONB(db, territory); err != nil {
						log.Printf("Failed to create territory %s: %v", territory.Name, err)
					} else {
						log.Printf("Created territory: %s", territory.Name)
					}
				} else {
					log.Printf("Error checking territory %s: %v", territory.Name, err)
				}
			} else {
				log.Printf("Territory %s already exists", territory.Name)
			}
		}
	*/

	log.Println("Database seeding completed successfully!")
}

func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

// getTypeID retrieves the ID of a type by its code and table name
func getTypeID(db *gorm.DB, code, tableName string, tenantID string) *string {
	var id string
	query := fmt.Sprintf("SELECT id FROM %s WHERE code = ? AND tenant_id = ? LIMIT 1", tableName)

	if err := db.Raw(query, code, tenantID).Scan(&id).Error; err != nil {
		log.Printf("Warning: Could not find %s with code %s: %v", tableName, code, err)
		return nil
	}

	if id == "" {
		log.Printf("Warning: No ID found for %s with code %s", tableName, code)
		return nil
	}

	return &id
}

// createTerritoryWithJSONB creates a territory with proper JSONB handling
func createTerritoryWithJSONB(db *gorm.DB, territory models.Territory) error {
	// Use raw SQL to properly handle JSONB fields
	query := `
		INSERT INTO territories (
			id, name, type_id, countries, states, cities, postal_codes, 
			industries, company_size, tenant_id, created_at, updated_at
		) VALUES (
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
		)
	`

	// Convert slices to JSON strings for proper JSONB insertion
	countriesJSON := "[]"
	if len(territory.Countries) > 0 {
		countriesJSON = fmt.Sprintf(`["%s"]`, strings.Join(territory.Countries, `","`))
	}

	statesJSON := "[]"
	if len(territory.States) > 0 {
		statesJSON = fmt.Sprintf(`["%s"]`, strings.Join(territory.States, `","`))
	}

	citiesJSON := "[]"
	if len(territory.Cities) > 0 {
		citiesJSON = fmt.Sprintf(`["%s"]`, strings.Join(territory.Cities, `","`))
	}

	postalCodesJSON := "[]"
	if len(territory.PostalCodes) > 0 {
		postalCodesJSON = fmt.Sprintf(`["%s"]`, strings.Join(territory.PostalCodes, `","`))
	}

	industriesJSON := "[]"
	if len(territory.Industries) > 0 {
		industriesJSON = fmt.Sprintf(`["%s"]`, strings.Join(territory.Industries, `","`))
	}

	companySizeJSON := "[]"
	if len(territory.CompanySize) > 0 {
		companySizeJSON = fmt.Sprintf(`["%s"]`, strings.Join(territory.CompanySize, `","`))
	}

	return db.Exec(query,
		territory.ID, territory.Name, territory.TypeID,
		countriesJSON, statesJSON, citiesJSON, postalCodesJSON,
		industriesJSON, companySizeJSON, territory.TenantID,
	).Error
}

func initDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := cfg.DatabaseURL

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Enable UUID extension for PostgreSQL
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Create UUID extension if it doesn't exist
	_, err = sqlDB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	if err != nil {
		return nil, err
	}

	return db, nil
}
