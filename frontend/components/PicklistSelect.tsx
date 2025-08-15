'use client'

import React, { useState, useEffect, useRef } from 'react'
import { ChevronDown, X, Search } from 'lucide-react'

export interface PicklistItem {
  id: string
  name: string
  code: string
  description?: string
  isActive: boolean
}

interface PicklistSelectProps {
  items: PicklistItem[]
  value: string | string[] | null
  onChange: (value: string | string[] | null) => void
  placeholder?: string
  label?: string
  multiple?: boolean
  searchable?: boolean
  disabled?: boolean
  error?: string
  className?: string
  onSearch?: (query: string) => void
  loading?: boolean
  hasMore?: boolean
  onLoadMore?: () => void
}

export function PicklistSelect({
  items,
  value,
  onChange,
  placeholder = "Select an option",
  label,
  multiple = false,
  searchable = true,
  disabled = false,
  error,
  className = "",
  onSearch,
  loading = false,
  hasMore = false,
  onLoadMore
}: PicklistSelectProps) {
  const [isOpen, setIsOpen] = useState(false)
  const [searchQuery, setSearchQuery] = useState("")
  const [filteredItems, setFilteredItems] = useState<PicklistItem[]>(items)
  const dropdownRef = useRef<HTMLDivElement>(null)
  const searchInputRef = useRef<HTMLInputElement>(null)

  useEffect(() => {
    setFilteredItems(items)
  }, [items])

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        setIsOpen(false)
      }
    }

    document.addEventListener('mousedown', handleClickOutside)
    return () => document.removeEventListener('mousedown', handleClickOutside)
  }, [])

  useEffect(() => {
    if (isOpen && searchable && searchInputRef.current) {
      searchInputRef.current.focus()
    }
  }, [isOpen, searchable])

  const handleSearch = (query: string) => {
    setSearchQuery(query)
    
    if (onSearch) {
      onSearch(query)
    } else {
      // Local filtering
      const filtered = items.filter(item =>
        item.name.toLowerCase().includes(query.toLowerCase()) ||
        item.code.toLowerCase().includes(query.toLowerCase())
      )
      setFilteredItems(filtered)
    }
  }

  const handleSelect = (item: PicklistItem) => {
    if (multiple) {
      const currentValue = Array.isArray(value) ? value : []
      const newValue = currentValue.includes(item.id)
        ? currentValue.filter(id => id !== item.id)
        : [...currentValue, item.id]
      onChange(newValue.length > 0 ? newValue : null)
    } else {
      onChange(item.id)
      setIsOpen(false)
    }
  }

  const handleRemove = (itemId: string) => {
    if (multiple && Array.isArray(value)) {
      const newValue = value.filter(id => id !== itemId)
      onChange(newValue.length > 0 ? newValue : null)
    }
  }

  const getSelectedItems = (): PicklistItem[] => {
    if (!value) return []
    
    if (multiple && Array.isArray(value)) {
      return items.filter(item => value.includes(item.id))
    } else if (!multiple && typeof value === 'string') {
      const item = items.find(item => item.id === value)
      return item ? [item] : []
    }
    
    return []
  }

  const getDisplayValue = (): string => {
    const selected = getSelectedItems()
    
    if (selected.length === 0) return ""
    if (selected.length === 1) return selected[0].name
    if (multiple) return `${selected.length} items selected`
    
    return ""
  }

  const isSelected = (itemId: string): boolean => {
    if (!value) return false
    
    if (multiple && Array.isArray(value)) {
      return value.includes(itemId)
    } else if (!multiple && typeof value === 'string') {
      return value === itemId
    }
    
    return false
  }

  const handleLoadMore = () => {
    if (onLoadMore && hasMore && !loading) {
      onLoadMore()
    }
  }

  return (
    <div className={`relative ${className}`}>
      {label && (
        <label className="block text-sm font-medium text-gray-700 mb-1">
          {label}
        </label>
      )}
      
      <div className="relative">
        <div
          ref={dropdownRef}
          className={`relative border rounded-md shadow-sm cursor-pointer ${
            error ? 'border-red-300' : 'border-gray-300'
          } ${disabled ? 'bg-gray-50 cursor-not-allowed' : 'bg-white hover:border-gray-400'}`}
        >
          <div
            className="flex items-center justify-between p-3"
            onClick={() => !disabled && setIsOpen(!isOpen)}
          >
            <div className="flex-1 min-w-0">
              {multiple && Array.isArray(value) && value.length > 0 ? (
                <div className="flex flex-wrap gap-1">
                  {getSelectedItems().map(item => (
                    <span
                      key={item.id}
                      className="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-blue-100 text-blue-800"
                    >
                      {item.name}
                      <button
                        type="button"
                        onClick={(e) => {
                          e.stopPropagation()
                          handleRemove(item.id)
                        }}
                        className="ml-1 inline-flex items-center justify-center w-4 h-4 rounded-full text-blue-400 hover:bg-blue-200 hover:text-blue-500"
                      >
                        <X className="w-3 h-3" />
                      </button>
                    </span>
                  ))}
                </div>
              ) : (
                <span className={getDisplayValue() ? 'text-gray-900' : 'text-gray-500'}>
                  {getDisplayValue() || placeholder}
                </span>
              )}
            </div>
            
            <ChevronDown
              className={`w-5 h-5 text-gray-400 transition-transform ${
                isOpen ? 'rotate-180' : ''
              }`}
            />
          </div>

          {isOpen && !disabled && (
            <div className="absolute z-50 w-full mt-1 bg-white border border-gray-300 rounded-md shadow-lg max-h-60 overflow-hidden">
              {searchable && (
                <div className="p-2 border-b border-gray-200">
                  <div className="relative">
                    <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-gray-400" />
                    <input
                      ref={searchInputRef}
                      type="text"
                      placeholder="Search..."
                      value={searchQuery}
                      onChange={(e) => handleSearch(e.target.value)}
                      className="w-full pl-10 pr-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                      onClick={(e) => e.stopPropagation()}
                    />
                  </div>
                </div>
              )}

              <div className="max-h-48 overflow-y-auto">
                {filteredItems.length === 0 ? (
                  <div className="px-3 py-2 text-sm text-gray-500">
                    {loading ? 'Loading...' : 'No options found'}
                  </div>
                ) : (
                  <>
                    {filteredItems.map((item) => (
                      <div
                        key={item.id}
                        className={`px-3 py-2 cursor-pointer hover:bg-gray-100 ${
                          isSelected(item.id) ? 'bg-blue-50 text-blue-900' : 'text-gray-900'
                        }`}
                        onClick={() => handleSelect(item)}
                      >
                        <div className="flex items-center justify-between">
                          <div>
                            <div className="font-medium">{item.name}</div>
                            {item.description && (
                              <div className="text-sm text-gray-500">{item.description}</div>
                            )}
                          </div>
                          {isSelected(item.id) && (
                            <div className="text-blue-600">
                              {multiple ? (
                                <div className="w-4 h-4 bg-blue-600 rounded-sm flex items-center justify-center">
                                  <svg className="w-3 h-3 text-white" fill="currentColor" viewBox="0 0 20 20">
                                    <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
                                  </svg>
                                </div>
                              ) : (
                                <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
                                  <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
                                </svg>
                              )}
                            </div>
                          )}
                        </div>
                      </div>
                    ))}
                    
                    {hasMore && (
                      <div className="px-3 py-2 border-t border-gray-200">
                        <button
                          type="button"
                          onClick={handleLoadMore}
                          disabled={loading}
                          className="w-full text-sm text-blue-600 hover:text-blue-800 disabled:opacity-50 disabled:cursor-not-allowed"
                        >
                          {loading ? 'Loading...' : 'Load more...'}
                        </button>
                      </div>
                    )}
                  </>
                )}
              </div>
            </div>
          )}
        </div>
      </div>

      {error && (
        <p className="mt-1 text-sm text-red-600">{error}</p>
      )}
    </div>
  )
} 