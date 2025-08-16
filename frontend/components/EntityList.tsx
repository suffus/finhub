import React, { useState } from 'react'
import { ChevronUp, ChevronDown, ChevronLeft, ChevronRight, Filter, RefreshCw, Eye, Settings } from 'lucide-react'
import { useEntityList } from '../hooks/useEntityList'
import { Column, EntityViewConfig } from '../services/api'
import { PicklistFilter } from './PicklistFilter'

interface EntityListProps {
  entityType: string
  title?: string
  onEntityClick?: (entity: Record<string, any>) => void
  className?: string
}

export function EntityList({ 
  entityType, 
  title, 
  onEntityClick,
  className = '' 
}: EntityListProps) {
  const {
    entities,
    views,
    currentView,
    page,
    pageSize,
    totalCount,
    totalPages,
    hasMore,
    sortBy,
    sortOrder,
    filters,
    loading,
    error,
    goToPage,
    changePageSize,
    changeSorting,
    changeFilters,
    changeView,
    refresh,
    startIndex,
    endIndex
  } = useEntityList({ entityType })

  const [showFilters, setShowFilters] = useState(false)
  const [localFilters, setLocalFilters] = useState<Record<string, any>>(filters)
  const [columnFilters, setColumnFilters] = useState<Record<string, {value: any, operator: string}>>({})

  // Apply filters
  const applyFilters = () => {
    // Merge global search with column-specific filters
    const mergedFilters = { ...localFilters }
    
    // Only include non-empty column filters
    Object.entries(columnFilters).forEach(([key, filter]) => {
      if (filter.value !== '' && filter.value !== null && filter.value !== undefined) {
        mergedFilters[key] = filter
      }
    })
    
    changeFilters(mergedFilters)
    setShowFilters(false)
  }

  // Reset filters
  const resetFilters = () => {
    setLocalFilters({})
    setColumnFilters({}) 
    changeFilters({})
  }

  // Render appropriate filter control based on column type
  const renderFilterControl = (column: Column) => {
    const filter = columnFilters[column.key] || { value: '', operator: 'eq' };
    
    const updateFilter = (value: any, operator = filter.operator) => {
      setColumnFilters(prev => ({
        ...prev,
        [column.key]: { value, operator }
      }));
      
      // Update the localFilters object which is used when applying filters
      setLocalFilters(prev => ({
        ...prev,
        [column.key]: { value, operator }
      }));
    };
    
    // Default operators
    const operators = [
      { value: 'eq', label: 'Equals' },
      { value: 'neq', label: 'Not Equals' },
      { value: 'contains', label: 'Contains' },
      { value: 'starts', label: 'Starts With' },
      { value: 'ends', label: 'Ends With' }
    ];
    
    // Number-specific operators
    const numberOperators = [
      { value: 'eq', label: 'Equals' },
      { value: 'neq', label: 'Not Equals' },
      { value: 'gt', label: 'Greater Than' },
      { value: 'gte', label: 'Greater/Equal' },
      { value: 'lt', label: 'Less Than' },
      { value: 'lte', label: 'Less/Equal' }
    ];
    
    // Date-specific operators
    const dateOperators = [
      { value: 'eq', label: 'Equals' },
      { value: 'neq', label: 'Not Equals' },
      { value: 'gt', label: 'After' },
      { value: 'gte', label: 'On or After' },
      { value: 'lt', label: 'Before' },
      { value: 'lte', label: 'On or Before' }
    ];
    
    switch (column.type) {
      case 'number':
      case 'currency':
      case 'percentage':
        return (
          <div className="flex space-x-2">
            <select
              value={filter.operator}
              onChange={(e) => updateFilter(filter.value, e.target.value)}
              className="flex-shrink-0 w-24 border border-gray-300 rounded px-1 py-1 text-xs"
            >
              {numberOperators.map(op => (
                <option key={op.value} value={op.value}>{op.label}</option>
              ))}
            </select>
            <input
              type="number"
              value={filter.value || ''}
              onChange={(e) => updateFilter(e.target.value)}
              className="w-full border border-gray-300 rounded px-2 py-1 text-xs"
              placeholder={`Filter by ${column.label.toLowerCase()}`}
            />
          </div>
        );
      
      case 'date':
        return (
          <div className="flex space-x-2">
            <select
              value={filter.operator}
              onChange={(e) => updateFilter(filter.value, e.target.value)}
              className="flex-shrink-0 w-24 border border-gray-300 rounded px-1 py-1 text-xs"
            >
              {dateOperators.map(op => (
                <option key={op.value} value={op.value}>{op.label}</option>
              ))}
            </select>
            <input
              type="date"
              value={filter.value || ''}
              onChange={(e) => updateFilter(e.target.value)}
              className="w-full border border-gray-300 rounded px-2 py-1 text-xs"
            />
          </div>
        );
      
      case 'status':
      case 'select':
        // For picklist fields like industry, company size, lead status, etc.
        const picklistType = column.key === 'industryId' ? 'industries' :
                            column.key === 'sizeId' ? 'companysizes' :
                            column.key === 'statusId' ? 'leadstatuses' :
                            column.key === 'temperatureId' ? 'leadtemperatures' :
                            null;
                            
        if (picklistType) {
          return (
            <PicklistFilter
              entityType={picklistType}
              value={filter.value || null}
              onChange={(value) => updateFilter(value)}
              placeholder={`Select ${column.label.toLowerCase()}`}
              operator={filter.operator}
              onOperatorChange={(op) => updateFilter(filter.value, op)}
              showOperator={true}
            />
          );
        }
        
        // Fallback for other select fields
        return (
          <div className="flex space-x-2">
            <select
              value={filter.operator}
              onChange={(e) => updateFilter(filter.value, e.target.value)}
              className="flex-shrink-0 w-24 border border-gray-300 rounded px-1 py-1 text-xs"
            >
              {operators.slice(0, 2).map(op => (
                <option key={op.value} value={op.value}>{op.label}</option>
              ))}
            </select>
            <input
              type="text"
              value={filter.value || ''}
              onChange={(e) => updateFilter(e.target.value)}
              className="w-full border border-gray-300 rounded px-2 py-1 text-xs"
              placeholder={`Select ${column.label.toLowerCase()}`}
            />
          </div>
        );
      
      case 'boolean':
        return (
          <select
            value={filter.value || ''}
            onChange={(e) => updateFilter(e.target.value)}
            className="w-full border border-gray-300 rounded px-2 py-1 text-xs"
          >
            <option value="">Any</option>
            <option value="true">Yes</option>
            <option value="false">No</option>
          </select>
        );
      
      default: // text, link, etc.
        return (
          <div className="flex space-x-2">
            <select
              value={filter.operator}
              onChange={(e) => updateFilter(filter.value, e.target.value)}
              className="flex-shrink-0 w-24 border border-gray-300 rounded px-1 py-1 text-xs"
            >
              {operators.map(op => (
                <option key={op.value} value={op.value}>{op.label}</option>
              ))}
            </select>
            <input
              type="text"
              value={filter.value || ''}
              onChange={(e) => updateFilter(e.target.value)}
              className="w-full border border-gray-300 rounded px-2 py-1 text-xs"
              placeholder={`Filter by ${column.label.toLowerCase()}`}
            />
          </div>
        );
    }
  };

  // Format cell value based on column type
  const formatCellValue = (value: any, column: Column) => {
    if (value === null || value === undefined) {
      return '-'
    }

    switch (column.type) {
      case 'date':
        return new Date(value).toLocaleDateString()
      case 'currency':
        return new Intl.NumberFormat('en-US', {
          style: 'currency',
          currency: 'USD'
        }).format(Number(value))
      case 'percentage':
        return `${value}%`
      case 'number':
        return new Intl.NumberFormat('en-US').format(Number(value))
      case 'link':
        return value ? (
          <a 
            href={value.startsWith('http') ? value : `https://${value}`} 
            target="_blank" 
            rel="noopener noreferrer"
            className="text-blue-600 hover:text-blue-800 underline"
          >
            {value}
          </a>
        ) : '-'
      case 'status':
        return (
          <span className={`px-2 py-1 rounded-full text-xs font-medium ${
            value === 'Hot' ? 'bg-red-100 text-red-800' :
            value === 'Warm' ? 'bg-yellow-100 text-yellow-800' :
            value === 'Cold' ? 'bg-blue-100 text-blue-800' :
            'bg-gray-100 text-gray-800'
          }`}>
            {value}
          </span>
        )
      default:
        return String(value)
    }
  }

  // Render sort indicator
  const renderSortIndicator = (columnKey: string) => {
    if (!columnKey || sortBy !== columnKey) {
      return null
    }
    return sortOrder === 'asc' ? (
      <ChevronUp className="w-4 h-4 ml-1" />
    ) : (
      <ChevronDown className="w-4 h-4 ml-1" />
    )
  }

  if (error) {
    return (
      <div className={`bg-white rounded-lg shadow p-6 ${className}`}>
        <div className="text-center text-red-600">
          <p className="text-lg font-medium">Error loading {entityType}</p>
          <p className="text-sm mt-2">{error}</p>
          <button
            onClick={refresh}
            className="mt-4 px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
          >
            Try Again
          </button>
        </div>
      </div>
    )
  }

  return (
    <div className={`bg-white rounded-lg shadow ${className}`}>
      {/* Header */}
      <div className="px-6 py-4 border-b border-gray-200">
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-4">
            <h2 className="text-xl font-semibold text-gray-900">
              {title || `${entityType.charAt(0).toUpperCase() + entityType.slice(1)}`}
            </h2>
            <span className="text-sm text-gray-500">
              {totalCount > 0 ? `${startIndex}-${endIndex} of ${totalCount}` : 'No results'}
            </span>
          </div>
          
          <div className="flex items-center space-x-2">
            {/* View selector */}
            {views.length > 1 && (
              <div className="flex items-center space-x-2">
                <Eye className="w-4 h-4 text-gray-400" />
                <select
                  value={currentView?.name || ''}
                  onChange={(e) => changeView(e.target.value)}
                  className="border border-gray-300 rounded px-3 py-1 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
                >
                  {views.map(view => (
                    <option key={view.name} value={view.name}>
                      {view.displayName}
                    </option>
                  ))}
                </select>
              </div>
            )}

            {/* Filter toggle */}
            <button
              onClick={() => setShowFilters(!showFilters)}
              className={`p-2 rounded ${
                showFilters ? 'bg-blue-100 text-blue-600' : 'text-gray-400 hover:text-gray-600'
              }`}
            >
              <Filter className="w-4 h-4" />
            </button>

            {/* Refresh */}
            <button
              onClick={refresh}
              disabled={loading}
              className="p-2 text-gray-400 hover:text-gray-600 disabled:opacity-50"
            >
              <RefreshCw className={`w-4 h-4 ${loading ? 'animate-spin' : ''}`} />
            </button>
          </div>
        </div>
      </div>

      {/* Filters */}
      {showFilters && (
        <div className="px-6 py-4 border-b border-gray-200 bg-gray-50">
          <div className="mb-4">
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Global Search
                </label>
                <input
                  type="text"
                  placeholder="Search across all fields..."
                  value={localFilters.search || ''}
                  onChange={(e) => setLocalFilters(prev => ({ ...prev, search: e.target.value }))}
                  className="w-full border border-gray-300 rounded px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>
              
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Page Size
                </label>
                <select
                  value={pageSize}
                  onChange={(e) => changePageSize(Number(e.target.value))}
                  className="w-full border border-gray-300 rounded px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
                >
                  <option value={10}>10 per page</option>
                  <option value={20}>20 per page</option>
                  <option value={50}>50 per page</option>
                  <option value={100}>100 per page</option>
                </select>
              </div>

              <div className="flex items-end space-x-2">
                <button
                  onClick={applyFilters}
                  className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 text-sm"
                >
                  Apply
                </button>
                <button
                  onClick={resetFilters}
                  className="px-4 py-2 border border-gray-300 text-gray-700 rounded hover:bg-gray-50 text-sm"
                >
                  Reset
                </button>
              </div>
            </div>
          </div>
          
          {/* Column-specific filters */}
          {currentView && currentView.columns && currentView.columns.filter(col => col.filterable).length > 0 && (
            <div>
              <h3 className="text-sm font-medium text-gray-700 mb-3 border-b pb-2">Column Filters</h3>
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                {currentView.columns
                  .filter(column => column.filterable)
                  .map(column => (
                    <div key={column.key} className="space-y-1">
                      <label className="block text-xs font-medium text-gray-500">
                        {column.label}
                      </label>
                      {renderFilterControl(column)}
                    </div>
                  ))
                }
              </div>
            </div>
          )}
        </div>
      )}

      {/* Table */}
      <div className="overflow-x-auto">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              {currentView && currentView.columns && currentView.columns.map((column) => (
                <th
                  key={column.key}
                  className={`px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider ${
                    column.align === 'center' ? 'text-center' :
                    column.align === 'right' ? 'text-right' : 'text-left'
                  }`}
                  style={{ width: column.width }}
                >
                  {column.sortable ? (
                    <button
                      onClick={() => changeSorting(column.key)}
                      className="flex items-center hover:text-gray-700 focus:outline-none focus:text-gray-700"
                    >
                      {column.label}
                      {renderSortIndicator(column.key)}
                    </button>
                  ) : (
                    column.label
                  )}
                </th>
              ))}
            </tr>
          </thead>
          <tbody className="bg-white divide-y divide-gray-200">
            {loading ? (
              <tr>
                <td colSpan={(currentView && currentView.columns ? currentView.columns.length : 0) || 1} className="px-6 py-12 text-center">
                  <div className="flex items-center justify-center">
                    <RefreshCw className="w-6 h-6 animate-spin text-gray-400 mr-2" />
                    <span className="text-gray-500">Loading...</span>
                  </div>
                </td>
              </tr>
            ) : entities.length === 0 ? (
              <tr>
                <td colSpan={(currentView && currentView.columns ? currentView.columns.length : 0) || 1} className="px-6 py-12 text-center">
                  <span className="text-gray-500">No {entityType} found</span>
                </td>
              </tr>
            ) : (
              entities.map((entity, index) => (
                <tr
                  key={entity.id || index}
                  onClick={() => onEntityClick?.(entity)}
                  className={`hover:bg-gray-50 ${
                    onEntityClick ? 'cursor-pointer' : ''
                  }`}
                >
                  {currentView && currentView.columns && currentView.columns.map((column) => (
                    <td
                      key={column.key}
                      className={`px-6 py-4 whitespace-nowrap text-sm text-gray-900 ${
                        column.align === 'center' ? 'text-center' :
                        column.align === 'right' ? 'text-right' : 'text-left'
                      }`}
                    >
                      {formatCellValue(entity[column.key], column)}
                    </td>
                  ))}
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>

      {/* Pagination */}
      {totalPages > 1 && (
        <div className="px-6 py-4 border-t border-gray-200">
          <div className="flex items-center justify-between">
            <div className="text-sm text-gray-700">
              Showing {startIndex} to {endIndex} of {totalCount} results
            </div>
            
            <div className="flex items-center space-x-2">
              <button
                onClick={() => goToPage(page - 1)}
                disabled={page <= 1}
                className="p-2 border border-gray-300 rounded disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50"
              >
                <ChevronLeft className="w-4 h-4" />
              </button>
              
              {Array.from({ length: Math.min(5, totalPages) }, (_, i) => {
                let pageNum
                if (totalPages <= 5) {
                  pageNum = i + 1
                } else if (page <= 3) {
                  pageNum = i + 1
                } else if (page >= totalPages - 2) {
                  pageNum = totalPages - 4 + i
                } else {
                  pageNum = page - 2 + i
                }
                
                return (
                  <button
                    key={pageNum}
                    onClick={() => goToPage(pageNum)}
                    className={`px-3 py-2 text-sm border rounded ${
                      page === pageNum
                        ? 'bg-blue-600 text-white border-blue-600'
                        : 'border-gray-300 text-gray-700 hover:bg-gray-50'
                    }`}
                  >
                    {pageNum}
                  </button>
                )
              })}
              
              <button
                onClick={() => goToPage(page + 1)}
                disabled={page >= totalPages}
                className="p-2 border border-gray-300 rounded disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50"
              >
                <ChevronRight className="w-4 h-4" />
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
} 