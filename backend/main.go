package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"finhub-backend/config"
	"finhub-backend/handlers"
	"finhub-backend/middleware"
	"finhub-backend/models"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := initDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate models
	if err := db.AutoMigrate(
		// First, create tables without foreign keys
		&models.Tenant{},
		&models.UserRole{},
		&models.Industry{},
		&models.CompanySize{},
		&models.LeadStatus{},
		&models.LeadTemperature{},
		&models.Pipeline{},
		&models.Stage{},
		&models.MarketingSourceType{},
		&models.MarketingSource{},
		&models.MarketingAssetType{},
		&models.MarketingAsset{},
		&models.CommunicationType{},
		&models.TaskType{},
		&models.TerritoryType{},
		&models.Territory{},
		&models.CustomField{},
		&models.CustomFieldValue{},
		&models.CustomObject{},
		&models.ActivityLog{},
		&models.DealStageHistory{},
		&models.PhoneNumberType{},
		&models.PhoneNumber{},
		&models.EmailAddressType{},
		&models.EmailAddress{},
		&models.AddressType{},
		&models.Address{},
		&models.SocialMediaType{},
		&models.SocialMediaAccount{},
		&models.CommunicationAttachment{},

		// Then create tables with foreign keys
		&models.User{},
		&models.Company{},
		&models.Contact{},
		&models.Lead{},
		&models.Deal{},
		&models.Task{},
		&models.Communication{},
	); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(db, cfg)
	userHandler := handlers.NewUserHandler(db)
	companyHandler := handlers.NewCompanyHandler(db)
	contactHandler := handlers.NewContactHandler(db)
	leadHandler := handlers.NewLeadHandler(db)
	dealHandler := handlers.NewDealHandler(db)
	picklistHandler := handlers.NewPicklistHandler(db)
	entityHandler := handlers.NewEntityHandler(db)

	// Setup router
	r := gin.Default()

	// CORS middleware
	r.Use(middleware.CORS())

	// Public routes
	r.POST("/api/auth/register", authHandler.Register)
	r.POST("/api/auth/login", authHandler.Login)

	// Protected routes
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware(cfg.JWTSecret))

	// User routes
	api.GET("/users/me", userHandler.GetCurrentUser)
	api.PUT("/users/me", userHandler.UpdateCurrentUser)

	// Company routes
	api.GET("/companies", companyHandler.GetCompanies)
	api.POST("/companies", companyHandler.CreateCompany)
	api.GET("/companies/:id", companyHandler.GetCompany)
	api.PUT("/companies/:id", companyHandler.UpdateCompany)
	api.DELETE("/companies/:id", companyHandler.DeleteCompany)

	// Contact routes
	api.GET("/contacts", contactHandler.GetContacts)
	api.POST("/contacts", contactHandler.CreateContact)
	api.GET("/contacts/:id", contactHandler.GetContact)
	api.PUT("/contacts/:id", contactHandler.UpdateContact)
	api.DELETE("/contacts/:id", contactHandler.DeleteContact)

	// Lead routes
	api.GET("/leads", leadHandler.GetLeads)
	api.POST("/leads", leadHandler.CreateLead)
	api.GET("/leads/:id", leadHandler.GetLead)
	api.PUT("/leads/:id", leadHandler.UpdateLead)
	api.DELETE("/leads/:id", leadHandler.DeleteLead)

	// Deal routes
	api.GET("/deals", dealHandler.GetDeals)
	api.POST("/deals", dealHandler.CreateDeal)
	api.GET("/deals/:id", dealHandler.GetDeal)
	api.PUT("/deals/:id", dealHandler.UpdateDeal)
	api.DELETE("/deals/:id", dealHandler.DeleteDeal)

	// Picklist routes
	api.GET("/picklists/:entity", picklistHandler.GetPicklistByEntity)
	api.POST("/picklists/search", picklistHandler.SearchPicklist)

	// Entity routes
	api.POST("/entities/query", entityHandler.GetEntityList)
	api.GET("/entities/:entityType/views", entityHandler.GetEntityViews)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func initDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := cfg.DatabaseURL

	var logLevel logger.LogLevel
	if cfg.Debug {
		logLevel = logger.Info
	} else {
		logLevel = logger.Error
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logLevel),
		DisableForeignKeyConstraintWhenMigrating: false,
		PrepareStmt:                              true,
	})
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

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}
