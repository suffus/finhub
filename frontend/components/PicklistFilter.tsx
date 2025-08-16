'use client'

import { useState, useEffect } from 'react'
import { usePicklist } from '@/hooks/usePicklist'
import { ChevronDown, X, Check, Search } from 'lucide-react'

interface PicklistFilterProps {
  entityType: string // The picklist entity type (e.g., 'industries', 'companysizes')
  value: string | null // Current selected value ID
  onChange: (value: string | null) => void
  placeholder?: string
  label?: string
  className?: string
  operator?: string
  onOperatorChange?: (operator: string) => void
  showOperator?: boolean
  disabled?: boolean
}

export function PicklistFilter({
  entityType,
  value,
  onChange,
  placeholder = 'Select...',
  label,
  className = '',
  operator = 'eq',
  onOperatorChange,
  showOperator = false,
  disabled = false
}: PicklistFilterProps) {
  const [isOpen, setIsOpen] = useState(false)
  const [searchQuery, setSearchQuery] = useState('')
  
  // Get picklist items using the usePicklist hook
  const {
    items,
    loading,
    error,
    search,
    refresh
  } = usePicklist({
    entityType: entityType as any, // Cast to the expected type
    searchable: true,
    preload: true
  })

  // Find the selected item's name for display
  const selectedItem = items.find(item => item.id === value)
  
  // Close dropdown if clicked outside
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      const target = event.target as HTMLElement
      if (isOpen && !target.closest(`[data-picklist="${entityType}"]`)) {
        setIsOpen(false)
      }
    }
    
    document.addEventListener('mousedown', handleClickOutside)
    return () => document.removeEventListener('mousedown', handleClickOutside)
  }, [isOpen, entityType])

  // Handle search input
  const handleSearchChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const query = e.target.value
    setSearchQuery(query)
    
    // Debounce search
    const timeoutId = setTimeout(() => {
      search(query)
    }, 300)
    
    return () => clearTimeout(timeoutId)
  }

  // Handle item selection
  const handleSelect = (itemId: string | null) => {
    onChange(itemId)
    setIsOpen(false)
    setSearchQuery('')
  }

  const operators = [
    { value: 'eq', label: 'Equals' },
    { value: 'neq', label: 'Not Equals' }
  ]

  return (
    <div className={`relative ${className}`} data-picklist={entityType}>
      {label && (
        <label className="block text-xs font-medium text-gray-500 mb-1">
          {label}
        </label>
      )}
      
      <div className="flex space-x-2">
        {/* Operator selector */}
        {showOperator && (
          <select
            value={operator}
            onChange={(e) => onOperatorChange?.(e.target.value)}
            disabled={disabled}
            className="flex-shrink-0 w-24 border border-gray-300 rounded px-1 py-1 text-xs"
          >
            {operators.map(op => (
              <option key={op.value} value={op.value}>{op.label}</option>
            ))}
          </select>
        )}
        
        {/* Picklist selector */}
        <div className="relative w-full">
          <button
            type="button"
            onClick={() => !disabled && setIsOpen(!isOpen)}
            disabled={disabled}
            className={`w-full flex items-center justify-between border rounded px-2 py-1 text-xs ${
              disabled 
                ? 'bg-gray-100 text-gray-400 cursor-not-allowed'
                : 'bg-white hover:bg-gray-50 text-gray-900'
            } ${
              isOpen ? 'border-blue-500 ring-1 ring-blue-500' : 'border-gray-300'
            }`}
          >
            <span className={`${!selectedItem && 'text-gray-500'}`}>
              {selectedItem ? selectedItem.name : placeholder}
            </span>
            <div className="flex items-center">
              {value && (
                <button
                  type="button"
                  onClick={(e) => {
                    e.stopPropagation()
                    handleSelect(null)
                  }}
                  className="mr-1 text-gray-400 hover:text-gray-600"
                >
                  <X className="h-3 w-3" />
                </button>
              )}
              <ChevronDown className="h-3 w-3 text-gray-400" />
            </div>
          </button>
          
          {isOpen && (
            <div className="absolute z-10 w-full mt-1 bg-white border border-gray-300 rounded shadow-lg">
              {/* Search input */}
              <div className="p-2 border-b border-gray-200">
                <div className="relative">
                  <input
                    type="text"
                    value={searchQuery}
                    onChange={handleSearchChange}
                    placeholder="Search..."
                    className="w-full border border-gray-300 rounded pl-8 pr-2 py-1 text-xs"
                    autoFocus
                  />
                  <Search className="h-3 w-3 text-gray-400 absolute left-2 top-1.5" />
                </div>
              </div>
              
              {/* List of items */}
              <div className="max-h-48 overflow-y-auto py-1">
                {loading ? (
                  <div className="text-center p-2 text-xs text-gray-500">
                    Loading...
                  </div>
                ) : error ? (
                  <div className="text-center p-2 text-xs text-red-500">
                    {error}
                  </div>
                ) : items.length === 0 ? (
                  <div className="text-center p-2 text-xs text-gray-500">
                    No items found
                  </div>
                ) : (
                  items.map((item) => (
                    <button
                      key={item.id}
                      type="button"
                      onClick={() => handleSelect(item.id)}
                      className="w-full text-left px-3 py-2 text-xs hover:bg-gray-100 flex items-center justify-between"
                    >
                      <span>{item.name}</span>
                      {value === item.id && (
                        <Check className="h-3 w-3 text-blue-500" />
                      )}
                    </button>
                  ))
                )}
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}