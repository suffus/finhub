import { useState, useEffect, useCallback } from 'react'
import { apiClient } from '@/services/api'
import { PicklistItem } from '@/components/PicklistSelect'

export type PicklistEntityType = 'industries' | 'companysizes' | 'leadstatuses' | 'leadtemperatures'

interface UsePicklistOptions {
  entityType: PicklistEntityType
  searchable?: boolean
  cacheKey?: string
  preload?: boolean
}

interface UsePicklistReturn {
  items: PicklistItem[]
  loading: boolean
  error: string | null
  search: (query: string) => void
  loadMore: () => void
  refresh: () => void
  hasMore: boolean
  searchQuery: string
}

// Simple in-memory cache
const picklistCache = new Map<string, { items: PicklistItem[]; timestamp: number }>()
const CACHE_DURATION = 5 * 60 * 1000 // 5 minutes

export function usePicklist({
  entityType,
  searchable = true,
  cacheKey,
  preload = true
}: UsePicklistOptions): UsePicklistReturn {
  const [items, setItems] = useState<PicklistItem[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [searchQuery, setSearchQuery] = useState('')
  const [hasMore, setHasMore] = useState(false)
  const [offset, setOffset] = useState(0)
  const [totalCount, setTotalCount] = useState(0)

  const cacheKeyFinal = cacheKey || `picklist_${entityType}`

  const fetchItems = useCallback(async (query = '', newOffset = 0, append = false) => {
    setLoading(true)
    setError(null)

    try {
      let response

      if (searchable && query) {
        // Use search endpoint for queries
        response = await apiClient.searchPicklist({
          query,
          entityType: entityType.replace('s', ''), // Remove plural
          limit: 50,
          offset: newOffset
        })
      } else {
        // Use simple endpoint for all items
        response = await apiClient.getPicklist(entityType)
      }

      if (response.data) {
        const newItems = response.data.items || []
        
        if (append) {
          setItems(prev => [...prev, ...newItems])
        } else {
          setItems(newItems)
        }
        
        setTotalCount(response.data.totalCount || 0)
        setHasMore(response.data.hasMore || false)
        setOffset(newOffset)
        
        // Cache the results if no search query
        if (!query && newOffset === 0) {
          picklistCache.set(cacheKeyFinal, {
            items: newItems,
            timestamp: Date.now()
          })
        }
      } else {
        setError(response.error || 'Failed to fetch items')
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred')
    } finally {
      setLoading(false)
    }
  }, [entityType, searchable, cacheKeyFinal])

  const search = useCallback((query: string) => {
    setSearchQuery(query)
    if (query.trim()) {
      fetchItems(query, 0, false)
    } else {
      // Reset to cached/preloaded items
      const cached = picklistCache.get(cacheKeyFinal)
      if (cached && Date.now() - cached.timestamp < CACHE_DURATION) {
        setItems(cached.items)
        setHasMore(false)
        setOffset(0)
      } else {
        fetchItems('', 0, false)
      }
    }
  }, [fetchItems, cacheKeyFinal])

  const loadMore = useCallback(() => {
    if (hasMore && !loading) {
      const newOffset = offset + 50
      fetchItems(searchQuery, newOffset, true)
    }
  }, [hasMore, loading, offset, searchQuery, fetchItems])

  const refresh = useCallback(() => {
    // Clear cache and refetch
    picklistCache.delete(cacheKeyFinal)
    setSearchQuery('')
    setOffset(0)
    fetchItems('', 0, false)
  }, [fetchItems, cacheKeyFinal])

  // Preload data on mount
  useEffect(() => {
    if (preload) {
      const cached = picklistCache.get(cacheKeyFinal)
      
      if (cached && Date.now() - cached.timestamp < CACHE_DURATION) {
        // Use cached data
        setItems(cached.items)
        setHasMore(false)
        setOffset(0)
      } else {
        // Fetch fresh data
        fetchItems('', 0, false)
      }
    }
  }, [preload, cacheKeyFinal, fetchItems])

  return {
    items,
    loading,
    error,
    search,
    loadMore,
    refresh,
    hasMore,
    searchQuery
  }
}

// Convenience hooks for specific entity types
export function useIndustries(options?: Omit<UsePicklistOptions, 'entityType'>) {
  return usePicklist({ ...options, entityType: 'industries' })
}

export function useCompanySizes(options?: Omit<UsePicklistOptions, 'entityType'>) {
  return usePicklist({ ...options, entityType: 'companysizes' })
}

export function useLeadStatuses(options?: Omit<UsePicklistOptions, 'entityType'>) {
  return usePicklist({ ...options, entityType: 'leadstatuses' })
}

export function useLeadTemperatures(options?: Omit<UsePicklistOptions, 'entityType'>) {
  return usePicklist({ ...options, entityType: 'leadtemperatures' })
} 