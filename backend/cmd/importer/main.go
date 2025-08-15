package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"

	"finhub-backend/config"
	"finhub-backend/models"

	"gorm.io/gorm"

	"gorm.io/driver/postgres"
)

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

var companySizes []models.CompanySize

func getCompanySizes(db *gorm.DB, size int) *models.CompanySize {
	if companySizes == nil {
		db.Order("min_employees ASC").Find(&companySizes)
	}
	for _, s := range companySizes {
		if s.MinEmployees != nil && s.MaxEmployees != nil && *s.MinEmployees <= size && *s.MaxEmployees >= size {
			return &s
		}
	}
	return nil
}

func createPhoneNumber(db *gorm.DB, number string, entityID string, entityType string, tenantID string) *models.PhoneNumber {
	var phoneNumberType models.PhoneNumberType
	db.Where("code = ?", "WORK").First(&phoneNumberType)
	phoneNumber := models.PhoneNumber{
		Number:     number,
		EntityID:   entityID,
		TypeID:     &phoneNumberType.ID,
		IsPrimary:  true,
		EntityType: entityType,
		TenantID:   tenantID,
	}
	db.Create(&phoneNumber)
	return &phoneNumber
}

func createEmailAddress(db *gorm.DB, email string, entityID string, entityType string, tenantID string) *models.EmailAddress {

	var emailAddressType models.EmailAddressType
	db.Where("code = ?", "WORK").First(&emailAddressType)
	emailAddress := models.EmailAddress{
		Email:      email,
		EntityID:   entityID,
		TypeID:     &emailAddressType.ID,
		IsPrimary:  true,
		EntityType: entityType,
		TenantID:   tenantID,
	}
	db.Create(&emailAddress)
	return &emailAddress
}

func main() {
	// Parse command line arguments
	contactPath := flag.String("contact", "", "Path to the contact CSV file to import")

	companyPath := flag.String("company", "", "Path to the company CSV file to import")
	dealPath := flag.String("deal", "", "Path to the deal CSV file to import")
	leadPath := flag.String("lead", "", "Path to the lead CSV file to import")
	userPath := flag.String("user", "", "Path to the user CSV file to import")
	dealHistoryPath := flag.String("dealHistory", "", "Path to the deal history CSV file to import")
	databaseUrl := flag.String("databaseUrl", "", "URL of the database to import data into")

	flag.Parse()
	// now we need the name of the model to populate the data

	if *contactPath == "" || *companyPath == "" || *dealPath == "" || *leadPath == "" || *userPath == "" || *dealHistoryPath == "" || *databaseUrl == "" {
		log.Fatal("Please provide all required arguments")
	}

	// now we need to load the csv files
	contactData, err := models.LoadCSVContacts(*contactPath)
	if err != nil {
		log.Fatal("Failed to load contacts:", err)
	}
	companyData, err := models.LoadCSVCompanies(*companyPath)
	if err != nil {
		log.Fatal("Failed to load companies:", err)
	}
	dealData, err := models.LoadCSVDeals(*dealPath)
	if err != nil {
		log.Fatal("Failed to load deals:", err)
	}
	leadData, err := models.LoadCSVLeads(*leadPath)
	if err != nil {
		log.Fatal("Failed to load leads:", err)
	}
	userData, err := models.LoadCSVUsers(*userPath)
	if err != nil {
		log.Fatal("Failed to load users:", err)
	}
	dealHistoryData, err := models.LoadCSVDealHistories(*dealHistoryPath)
	if err != nil {
		log.Fatal("Failed to load deal histories:", err)
	}

	if len(dealHistoryData) > 0 {
		fmt.Println("Deal histories are not supported yet")
	}

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
		}
	} else {
		log.Println("Using existing tenant with ID:", tenant.ID)
	}

	// now load the company sizes from the database and order them by MinEmployees

	dbUsers := map[string]models.User{}
	// first load the users
	for _, user := range userData {
		// find user by email
		var dbUser models.User
		db.Where("email = ?", user.Email).First(&dbUser)
		if dbUser.ID == "" {
			dbUser = models.User{
				FirstName: user.FirstName,
				LastName:  user.LastName,
				Email:     user.Email,
				TenantID:  tenant.ID,
			}
			db.Create(&dbUser)
		}

		dbUsers[user.UserID] = dbUser
	}

	// now get work phone number type and work email type - we need the IDs
	var workPhoneNumberType models.PhoneNumberType
	db.Where("code = ?", "WORK").First(&workPhoneNumberType)
	var workEmailType models.EmailAddressType
	db.Where("code = ?", "WORK").First(&workEmailType)
	// and the WORK address type
	var workAddressType models.AddressType
	db.Where("code = ?", "OFFICE").First(&workAddressType)

	dbCompanies := make(map[string]models.Company)

	// delete all companies

	// delete all deals
	db.Exec("DELETE FROM deals")
	db.Exec("DELETE FROM deal_stage_histories")
	db.Exec("DELETE FROM leads")
	db.Exec("DELETE FROM contacts")
	db.Exec("DELETE FROM companies")
	db.Exec("DELETE from email_addresses")
	db.Exec("DELETE from phone_numbers")
	db.Exec("DELETE from addresses")

	// now load the companies
	for _, company := range companyData {
		// first create the phone number and email address

		// get the company size
		// convert companySize to an int
		companySizeInt, err := strconv.Atoi(company.CompanySize)
		if err != nil {
			log.Fatal("Failed to convert company size to int:", err)
		}

		companySize := getCompanySizes(db, companySizeInt)
		// now create the company
		dbCompany := models.Company{
			Name:     company.CompanyName,
			Website:  &company.Website,
			TenantID: tenant.ID,
		}
		if companySize != nil {
			dbCompany.SizeID = &companySize.ID
		}
		fmt.Println("Creating company:", company.CompanyName, company.Website, company.Phone, company.Address, company.City, company.State, company.ZipCode, company.Country, company.CompanySize)
		db.Create(&dbCompany)

		dbCompanies[company.CompanyID] = dbCompany

		createPhoneNumber(db, company.Phone, dbCompany.ID, "company", tenant.ID)

		// now create the address
		address := models.Address{
			Street1:    &company.Address,
			City:       &company.City,
			State:      &company.State,
			PostalCode: &company.ZipCode,
			Country:    &company.Country,
			TypeID:     &workAddressType.ID,
			EntityID:   dbCompany.ID,
			EntityType: "company",
			TenantID:   tenant.ID,
		}
		db.Create(&address)

	}

	// delete all contacts

	dbContacts := make(map[string]models.Contact)

	for _, contact := range contactData {
		companyRecord, ok := dbCompanies[contact.CompanyID]
		if !ok {
			log.Fatal("Company not found for contact:", contact.CompanyID)
		}
		dbContact := models.Contact{
			FirstName:  contact.FirstName,
			LastName:   contact.LastName,
			CompanyID:  &companyRecord.ID,
			JobTitle:   &contact.JobTitle,
			Department: &contact.Department,
			TenantID:   tenant.ID,
		}
		fmt.Println("Creating contact:", contact.FirstName, contact.LastName, contact.JobTitle, contact.Department, contact.CompanyID, contact.Phone, contact.Email)
		db.Create(&dbContact)
		dbContacts[contact.ContactID] = dbContact

		// now create the phone number
		createPhoneNumber(db, contact.Phone, dbContact.ID, "contact", tenant.ID)

		// now create the email address
		createEmailAddress(db, contact.Email, dbContact.ID, "contact", tenant.ID)

	}

	// now do the leads
	dbLeads := map[string]models.Lead{}

	for _, lead := range leadData {
		var status models.LeadStatus
		db.Where("code = ?", lead.Status).First(&status)
		if status.ID == "" {
			status = models.LeadStatus{
				Name:     lead.Status,
				Code:     lead.Status,
				TenantID: tenant.ID,
			}
			db.Create(&status)
		}

		dbLead := models.Lead{
			FirstName: &lead.FirstName,
			LastName:  &lead.LastName,
			StatusID:  &status.ID,
			Source:    &lead.Source,
			TenantID:  tenant.ID,
		}
		if lead.ContactID != "" {
			contactRecord, ok := dbContacts[lead.ContactID]
			if !ok {
				log.Fatal("Company not found for lead:")
			}
			dbLead.ContactID = &contactRecord.ID
			dbLead.CompanyID = contactRecord.CompanyID
		}
		fmt.Println("Creating lead:", lead.FirstName, lead.LastName, lead.Status)
		db.Create(&dbLead)
		dbLeads[lead.LeadID] = dbLead

		// now create the phone number
		createPhoneNumber(db, lead.Phone, dbLead.ID, "lead", tenant.ID)
		// now create the email address
		createEmailAddress(db, lead.Email, dbLead.ID, "lead", tenant.ID)

	}

	// now do the deals
	dbDeals := map[string]models.Deal{}

	var pipeline models.Pipeline
	db.Where("name = ?", "Sales Pipeline").First(&pipeline)
	if pipeline.ID == "" {
		log.Fatal("Sales Pipeline Not Found")
	}

	for _, deal := range dealData {
		dealAmount, _ := strconv.ParseFloat(deal.Value, 64)
		dealProbability, _ := strconv.Atoi(deal.Probability)
		dealExpectedCloseDate, _ := models.ParseDate(deal.ExpectedCloseDate)
		var stage models.Stage
		err := db.Where("name = ? and tenant_id = ?", deal.Stage, tenant.ID).First(&stage).Error
		if err == gorm.ErrRecordNotFound {
			stage = models.Stage{
				Name:       deal.Stage,
				PipelineID: pipeline.ID,
				TenantID:   tenant.ID,
			}
			db.Create(&stage)
		}

		// get the user and contact

		dbUser, ok := dbUsers[deal.UserID]
		if !ok {
			log.Fatal("User not found for deal:", deal.UserID)
		}

		dbContact, ok := dbContacts[deal.ContactID]
		if !ok {
			log.Fatal("Contact not found for deal:", deal.ContactID)
		}

		dbCompanyId := *dbContact.CompanyID
		dbDeal := models.Deal{
			Name:              deal.DealName,
			Amount:            &dealAmount,
			Currency:          "USD",
			Probability:       dealProbability,
			CompanyID:         &dbCompanyId,
			ContactID:         &dbContact.ID,
			PipelineID:        pipeline.ID,
			AssignedUserID:    &dbUser.ID,
			ExpectedCloseDate: &dealExpectedCloseDate,
			StageID:           stage.ID,
			TenantID:          tenant.ID,
		}
		fmt.Println("Creating deal:", deal.DealName, deal.Value, deal.Probability, deal.ExpectedCloseDate, deal.Stage, deal.UserID, deal.ContactID)
		db.Create(&dbDeal)
		dbDeals[deal.DealID] = dbDeal

	}

}
