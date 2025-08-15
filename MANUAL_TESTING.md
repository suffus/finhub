# Manual Testing Guide for FinHub Picklist System

This guide will walk you through testing the picklist system manually, step by step.

## Prerequisites

1. **Backend running**: Make sure the backend is running on port 8080
2. **Database seeded**: Run the seed command to populate sample data
3. **Frontend running**: Start the frontend development server

## Step 1: Setup and Seed Database

```bash
# Navigate to backend directory
cd backend

# Build and seed the database
make setup

# Or run individually:
make build-seed
make seed
```

This will create:
- A default tenant
- Sample industries (Technology, Healthcare, Finance, etc.)
- Sample company sizes (Startup, Small Business, etc.)
- Sample lead statuses and temperatures

## Step 2: Test Backend Endpoints

### 2.1 Test Authentication (Required for Picklist Endpoints)

First, create a test user and get an authentication token:

```bash
# Register a test user
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "testpass123",
    "firstName": "Test",
    "lastName": "User"
  }'

# Login to get a token
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "testpass123"
  }'
```

Copy the token from the response.

### 2.2 Test Picklist Endpoints

Use the token to test the picklist endpoints:

```bash
# Replace YOUR_TOKEN_HERE with the actual token
TOKEN="YOUR_TOKEN_HERE"

# Test industries endpoint
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/picklists/industries

# Test company sizes endpoint
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/picklists/companysizes

# Test search endpoint
curl -X POST http://localhost:8080/api/picklists/search \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "query": "tech",
    "limit": 10,
    "offset": 0,
    "entityType": "industry"
  }'
```

## Step 3: Test Frontend

### 3.1 Start Frontend

```bash
# In a new terminal, navigate to frontend
cd frontend

# Install dependencies (if not already done)
npm install

# Start development server
npm run dev
```

### 3.2 Test Picklist Components

1. **Open the application**: Navigate to `http://localhost:3000`
2. **Login**: Use the test user credentials from step 2.1
3. **Navigate to Companies**: Click on the Companies section
4. **Create a Company**: Click "Add Company" button
5. **Test Industry Picklist**:
   - Click on the Industry dropdown
   - Verify the list shows sample industries
   - Try searching for "tech" or "health"
   - Select an industry
6. **Test Company Size Picklist**:
   - Click on the Company Size dropdown
   - Verify the list shows sample company sizes
   - Try searching for "startup" or "enterprise"
   - Select a company size
7. **Submit the form**: Fill in the company name and submit

## Step 4: Verify Data Persistence

### 4.1 Check Database

```bash
# Connect to your database and verify the company was created
# with the correct industry and size IDs

psql your_database_name
SELECT c.name, i.name as industry, cs.name as company_size 
FROM companies c 
LEFT JOIN industries i ON c.industry_id = i.id 
LEFT JOIN company_sizes cs ON c.size_id = cs.id;
```

### 4.2 Check API Response

```bash
# Get the created company (replace COMPANY_ID with actual ID)
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/companies/COMPANY_ID
```

Verify that the response includes the industry and size data.

## Step 5: Test Search Functionality

### 5.1 Test Industry Search

```bash
# Search for industries containing "tech"
curl -X POST http://localhost:8080/api/picklists/search \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "query": "tech",
    "limit": 10,
    "offset": 0,
    "entityType": "industry"
  }'

# Search for industries containing "health"
curl -X POST http://localhost:8080/api/picklists/search \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "query": "health",
    "limit": 10,
    "offset": 0,
    "entityType": "industry"
  }'
```

### 5.2 Test Company Size Search

```bash
# Search for company sizes containing "startup"
curl -X POST http://localhost:8080/api/picklists/search \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "query": "startup",
    "limit": 10,
    "offset": 0,
    "entityType": "companysize"
  }'
```

## Step 6: Test Edge Cases

### 6.1 Test Empty Search

```bash
# Search with empty query
curl -X POST http://localhost:8080/api/picklists/search \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "query": "",
    "limit": 10,
    "offset": 0,
    "entityType": "industry"
  }'
```

### 6.2 Test Pagination

```bash
# Test with offset
curl -X POST http://localhost:8080/api/picklists/search \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "query": "",
    "limit": 5,
    "offset": 5,
    "entityType": "industry"
  }'
```

### 6.3 Test Invalid Entity Type

```bash
# Test with invalid entity type
curl -X POST http://localhost:8080/api/picklists/search \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "query": "test",
    "limit": 10,
    "offset": 0,
    "entityType": "invalid_type"
  }'
```

## Step 7: Performance Testing

### 7.1 Test Caching

1. **First request**: Make a request to `/api/picklists/industries`
2. **Second request**: Make the same request immediately
3. **Compare response times**: The second request should be faster due to caching

### 7.2 Test Large Datasets

If you have many industries/sizes, test search performance:

```bash
# Test search with large limit
curl -X POST http://localhost:8080/api/picklists/search \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "query": "",
    "limit": 100,
    "offset": 0,
    "entityType": "industry"
  }'
```

## Troubleshooting

### Common Issues

1. **"Authorization header required"**: Make sure you're including the Bearer token
2. **"User not found"**: Verify the user exists and the token is valid
3. **Empty results**: Check if the database has been seeded
4. **CORS errors**: Ensure the backend CORS middleware is working

### Debug Steps

1. **Check backend logs**: Look for error messages in the backend console
2. **Verify database**: Check if the seed data was created
3. **Check authentication**: Verify the JWT token is valid
4. **Test endpoints individually**: Test each endpoint separately to isolate issues

## Success Criteria

The picklist system is working correctly if:

✅ **Backend endpoints return data** with proper authentication  
✅ **Frontend dropdowns populate** with industry and company size data  
✅ **Search functionality works** for finding specific items  
✅ **Data persists correctly** when creating companies  
✅ **Caching improves performance** for subsequent requests  
✅ **Error handling works** for invalid requests  
✅ **Multi-select works** (if implemented)  
✅ **Responsive design works** on different screen sizes  

## Next Steps

Once testing is complete:

1. **Document any issues** found during testing
2. **Optimize performance** if needed
3. **Add more test cases** for edge cases
4. **Implement additional features** like bulk operations
5. **Add monitoring** for production use 