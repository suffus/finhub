# FinHub Picklist System

This document describes the comprehensive picklist system implemented in FinHub, which provides dynamic dropdown selection capabilities with search, caching, and multi-select support.

## Overview

The picklist system allows users to select values from predefined lists that can be either hardcoded or dynamically loaded from database tables. It supports both single and multi-select modes, with built-in search capabilities and efficient caching.

## Features

- **Dynamic Data Loading**: Picklists can be populated from database tables
- **Search Capabilities**: Built-in search with database-level filtering for large lists
- **Caching**: Automatic caching of picklist data for performance
- **Multi-Select Support**: Can be configured for single or multiple selection
- **Active/Inactive Filtering**: Only shows active items by default
- **Tenant Isolation**: Data is properly isolated per tenant
- **Responsive UI**: Modern dropdown interface with search and keyboard navigation

## Backend Implementation

### Models

The system uses existing models for picklist data:

- `Industry`: Company industry classifications
- `CompanySize`: Company size ranges
- `LeadStatus`: Lead status options
- `LeadTemperature`: Lead temperature indicators

### Picklist Handler

Located at `backend/handlers/picklist.go`, the handler provides:

- **GetIndustries()**: Fetch all active industries
- **GetCompanySizes()**: Fetch all active company sizes
- **GetLeadStatuses()**: Fetch all active lead statuses
- **GetLeadTemperatures()**: Fetch all active lead temperatures
- **SearchPicklist()**: Generic search across any picklist entity
- **GetPicklistByEntity()**: Generic endpoint for any picklist type

### API Endpoints

```
GET  /api/picklists/:entity          # Get all items for a specific entity
POST /api/picklists/search           # Search picklist items
```

Supported entity types:
- `industries`
- `companysizes`
- `leadstatuses`
- `leadtemperatures`

### Search Parameters

The search endpoint accepts:

```json
{
  "query": "search term",
  "limit": 50,
  "offset": 0,
  "entityType": "industry"
}
```

## Frontend Implementation

### PicklistSelect Component

A reusable React component (`frontend/components/PicklistSelect.tsx`) that provides:

- Single and multi-select modes
- Built-in search functionality
- Loading states
- Error handling
- Keyboard navigation
- Responsive design

### Usage Example

```tsx
import { PicklistSelect } from '@/components/PicklistSelect'

<PicklistSelect
  items={industries}
  value={selectedIndustry}
  onChange={setSelectedIndustry}
  placeholder="Select an industry"
  label="Industry"
  searchable={true}
  multiple={false}
/>
```

### usePicklist Hook

A custom React hook (`frontend/hooks/usePicklist.ts`) that manages:

- Data fetching and caching
- Search functionality
- Pagination support
- Error handling
- Loading states

### Hook Usage

```tsx
import { useIndustries, useCompanySizes } from '@/hooks/usePicklist'

function MyComponent() {
  const { items: industries, loading, error, search } = useIndustries()
  const { items: sizes } = useCompanySizes()
  
  // Use the data...
}
```

## Database Seeding

The system includes a seeding script (`backend/cmd/seed/main.go`) that populates the database with sample data:

### Industries
- Technology
- Healthcare
- Finance
- Manufacturing
- Retail
- Education
- Real Estate
- Consulting
- Media & Entertainment
- Transportation & Logistics

### Company Sizes
- Startup (1-10 employees)
- Small Business (11-50 employees)
- Medium Business (51-200 employees)
- Large Business (201-1000 employees)
- Enterprise (1000+ employees)

### Lead Statuses
- New, Contacted, Qualified, Proposal, Negotiation, Converted, Lost

### Lead Temperatures
- Hot, Warm, Cold

## Setup and Usage

### 1. Backend Setup

```bash
cd backend

# Build and seed the database
make setup

# Or run individually:
make build-seed
make seed
```

### 2. Frontend Usage

```tsx
import { PicklistSelect } from '@/components/PicklistSelect'
import { useIndustries } from '@/hooks/usePicklist'

function CompanyForm() {
  const { items: industries, loading } = useIndustries()
  
  return (
    <PicklistSelect
      items={industries}
      value={industryId}
      onChange={setIndustryId}
      placeholder="Select industry"
      label="Industry"
      loading={loading}
    />
  )
}
```

## Configuration Options

### PicklistSelect Props

- `items`: Array of picklist items
- `value`: Selected value(s)
- `onChange`: Change handler
- `placeholder`: Placeholder text
- `label`: Field label
- `multiple`: Enable multi-select
- `searchable`: Enable search functionality
- `disabled`: Disable the component
- `error`: Error message
- `loading`: Show loading state

### usePicklist Options

- `entityType`: Type of picklist to fetch
- `searchable`: Enable search capabilities
- `cacheKey`: Custom cache key
- `preload`: Preload data on mount

## Performance Considerations

### Caching Strategy

- **In-Memory Cache**: 5-minute cache duration for frequently accessed data
- **Database Queries**: Optimized with proper indexing
- **Lazy Loading**: Data loaded only when needed

### Search Optimization

- **Database-Level Search**: Uses ILIKE for efficient text search
- **Pagination**: Supports large datasets with offset/limit
- **Indexing**: Ensure proper database indexes on search fields

## Extending the System

### Adding New Picklist Types

1. **Create Model**: Add new model to `models.go`
2. **Add Handler Methods**: Implement in `picklist.go`
3. **Update Routes**: Add to main.go
4. **Frontend Hook**: Create convenience hook in `usePicklist.ts`

### Example: Adding Product Categories

```go
// 1. Add to models.go
type ProductCategory struct {
  ID          string  `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
  Name        string  `json:"name" gorm:"not null"`
  Code        string  `json:"code" gorm:"not null"`
  Description *string `json:"description"`
  IsActive    bool    `json:"isActive" gorm:"column:is_active;default:true"`
  TenantID    string  `json:"tenantId" gorm:"column:tenant_id;type:uuid;not null"`
  Tenant      Tenant  `json:"tenant,omitempty" gorm:"foreignKey:TenantID"`
  CreatedAt   time.Time `json:"createdAt" gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
  UpdatedAt   time.Time `json:"updatedAt" gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`
}

// 2. Add handler methods
func (h *PicklistHandler) GetProductCategories(c *gin.Context) {
  // Implementation...
}

// 3. Add to routes
api.GET("/picklists/productcategories", picklistHandler.GetProductCategories)

// 4. Add frontend hook
export function useProductCategories(options?: Omit<UsePicklistOptions, 'entityType'>) {
  return usePicklist({ ...options, entityType: 'productcategories' })
}
```

## Best Practices

1. **Always use the hook**: Don't manually fetch picklist data
2. **Leverage caching**: The system automatically caches data
3. **Handle loading states**: Show appropriate loading indicators
4. **Error handling**: Always handle potential errors gracefully
5. **Search optimization**: Use search for large datasets
6. **Tenant isolation**: Ensure data is properly isolated

## Troubleshooting

### Common Issues

1. **Data not loading**: Check authentication and tenant setup
2. **Search not working**: Verify database indexes are created
3. **Caching issues**: Clear browser cache or restart the application
4. **Performance problems**: Check database query performance and indexing

### Debug Mode

Enable debug logging in the backend to see detailed picklist operations:

```bash
export DEBUG=true
make run
```

## Future Enhancements

- **Redis Caching**: Replace in-memory cache with Redis
- **Real-time Updates**: WebSocket support for live picklist updates
- **Advanced Search**: Full-text search with Elasticsearch
- **Custom Fields**: Support for custom picklist item properties
- **Bulk Operations**: Support for bulk picklist item management
- **Audit Logging**: Track picklist usage and changes 