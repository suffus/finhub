package models

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

// CSVCompany represents a company from the CSV data
type CSVCompany struct {
	CompanyID   string
	CompanyName string
	Industry    string
	CompanySize string
	Website     string
	Phone       string
	Address     string
	City        string
	State       string
	ZipCode     string
	Country     string
}

// CSVContact represents a contact from the CSV data
type CSVContact struct {
	ContactID       string
	CompanyID       string
	FirstName       string
	LastName        string
	Email           string
	Phone           string
	JobTitle        string
	Department      string
	CreatedDate     string
	LastContactDate string
}

// CSVUser represents a user from the CSV data
type CSVUser struct {
	UserID    string
	FirstName string
	LastName  string
	Email     string
	Role      string
	HireDate  string
	Territory string
}

// CSVDeal represents a deal from the CSV data
type CSVDeal struct {
	DealID            string
	ContactID         string
	UserID            string
	DealName          string
	Stage             string
	Value             string
	Probability       string
	ExpectedCloseDate string
	CreatedDate       string
	LastUpdated       string
}

// CSVDealHistory represents a deal history entry from the CSV data
type CSVDealHistory struct {
	HistoryID   string
	DealID      string
	Stage       string
	Value       string
	Probability string
	ChangedDate string
	UserID      string
	Notes       string
}

// CSVLead represents a lead from the CSV data
type CSVLead struct {
	LeadID           string
	FirstName        string
	LastName         string
	Email            string
	Phone            string
	CompanyName      string
	ContactID        string
	UserID           string
	Source           string
	Status           string
	CreatedDate      string
	LastActivityDate string
	EstimatedValue   string
	Probability      string
	Notes            string
}

// LoadCSVCompanies loads companies from a CSV file
func LoadCSVCompanies(filePath string) ([]CSVCompany, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Skip header row
	if _, err := reader.Read(); err != nil {
		return nil, fmt.Errorf("error reading header: %v", err)
	}

	var companies []CSVCompany

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading record: %v", err)
		}

		if len(record) < 11 {
			return nil, fmt.Errorf("invalid record format, expected 11 fields, got %d", len(record))
		}

		company := CSVCompany{
			CompanyID:   record[0],
			CompanyName: record[1],
			Industry:    record[2],
			CompanySize: record[3],
			Website:     record[4],
			Phone:       record[5],
			Address:     record[6],
			City:        record[7],
			State:       record[8],
			ZipCode:     record[9],
			Country:     record[10],
		}

		companies = append(companies, company)
	}

	return companies, nil
}

// LoadCSVContacts loads contacts from a CSV file
func LoadCSVContacts(filePath string) ([]CSVContact, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Skip header row
	if _, err := reader.Read(); err != nil {
		return nil, fmt.Errorf("error reading header: %v", err)
	}

	var contacts []CSVContact

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading record: %v", err)
		}

		if len(record) < 10 {
			return nil, fmt.Errorf("invalid record format, expected 10 fields, got %d", len(record))
		}

		contact := CSVContact{
			ContactID:       record[0],
			CompanyID:       record[1],
			FirstName:       record[2],
			LastName:        record[3],
			Email:           record[4],
			Phone:           record[5],
			JobTitle:        record[6],
			Department:      record[7],
			CreatedDate:     record[8],
			LastContactDate: record[9],
		}

		contacts = append(contacts, contact)
	}

	return contacts, nil
}

// LoadCSVUsers loads users from a CSV file
func LoadCSVUsers(filePath string) ([]CSVUser, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Skip header row
	if _, err := reader.Read(); err != nil {
		return nil, fmt.Errorf("error reading header: %v", err)
	}

	var users []CSVUser

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading record: %v", err)
		}

		if len(record) < 7 {
			return nil, fmt.Errorf("invalid record format, expected 7 fields, got %d", len(record))
		}

		user := CSVUser{
			UserID:    record[0],
			FirstName: record[1],
			LastName:  record[2],
			Email:     record[3],
			Role:      record[4],
			HireDate:  record[5],
			Territory: record[6],
		}

		users = append(users, user)
	}

	return users, nil
}

// LoadCSVDeals loads deals from a CSV file
func LoadCSVDeals(filePath string) ([]CSVDeal, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Skip header row
	if _, err := reader.Read(); err != nil {
		return nil, fmt.Errorf("error reading header: %v", err)
	}

	var deals []CSVDeal

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading record: %v", err)
		}

		if len(record) < 10 {
			return nil, fmt.Errorf("invalid record format, expected 10 fields, got %d", len(record))
		}

		deal := CSVDeal{
			DealID:            record[0],
			ContactID:         record[1],
			UserID:            record[2],
			DealName:          record[3],
			Stage:             record[4],
			Value:             record[5],
			Probability:       record[6],
			ExpectedCloseDate: record[7],
			CreatedDate:       record[8],
			LastUpdated:       record[9],
		}

		deals = append(deals, deal)
	}

	return deals, nil
}

// LoadCSVDealHistory loads deal history from a CSV file
func LoadCSVDealHistories(filePath string) ([]CSVDealHistory, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Skip header row
	if _, err := reader.Read(); err != nil {
		return nil, fmt.Errorf("error reading header: %v", err)
	}

	var histories []CSVDealHistory

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading record: %v", err)
		}

		if len(record) < 8 {
			return nil, fmt.Errorf("invalid record format, expected 8 fields, got %d", len(record))
		}

		history := CSVDealHistory{
			HistoryID:   record[0],
			DealID:      record[1],
			Stage:       record[2],
			Value:       record[3],
			Probability: record[4],
			ChangedDate: record[5],
			UserID:      record[6],
			Notes:       record[7],
		}

		histories = append(histories, history)
	}

	return histories, nil
}

// LoadCSVLeads loads leads from a CSV file
func LoadCSVLeads(filePath string) ([]CSVLead, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Skip header row
	if _, err := reader.Read(); err != nil {
		return nil, fmt.Errorf("error reading header: %v", err)
	}

	var leads []CSVLead

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading record: %v", err)
		}

		if len(record) < 15 {
			return nil, fmt.Errorf("invalid record format, expected 15 fields, got %d", len(record))
		}

		lead := CSVLead{
			LeadID:           record[0],
			FirstName:        record[1],
			LastName:         record[2],
			Email:            record[3],
			Phone:            record[4],
			CompanyName:      record[5],
			ContactID:        record[6],
			UserID:           record[7],
			Source:           record[8],
			Status:           record[9],
			CreatedDate:      record[10],
			LastActivityDate: record[11],
			EstimatedValue:   record[12],
			Probability:      record[13],
			Notes:            record[14],
		}

		leads = append(leads, lead)
	}

	return leads, nil
}

// Helper functions for data conversion

// ParseFloat parses a string to float64, returning 0 if empty or invalid
func ParseFloat(s string) float64 {
	if s == "" {
		return 0
	}
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return val
}

// ParseInt parses a string to int, returning 0 if empty or invalid
func ParseInt(s string) int {
	if s == "" {
		return 0
	}
	val, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return val
}

// ParseDate parses a date string to time.Time
func ParseDate(dateStr string) (time.Time, error) {
	if dateStr == "" {
		return time.Time{}, nil
	}
	return time.Parse("2006-01-02", dateStr)
}

// IsEmpty checks if a string is empty
func IsEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}
