import { useState, useEffect, useCallback } from 'react'
import { apiClient, EntityQueryRequest, EntityQueryResponse, EntityViewConfig } from '../services/api'

interface UseEntityListOptions {
  entityType: string
  initialPageSize?: number
  initialSortBy?: string
  initialSortOrder?: 'asc' | 'desc'
  initialFilters?: Record<string, any>
  view?: string
}

import { useRef } from 'react'

// Interface to track previous values for dependency comparison
interface PrevValues {
  page: number
  pageSize: number
  sortBy: string
  sortOrder: 'asc' | 'desc'
  view: string | null
}

export function useEntityList({
  entityType,
  initialPageSize = 20,
  initialSortBy,
  initialSortOrder = 'asc',
  initialFilters = {},
  view
}: UseEntityListOptions) {
  // For infinite scrolling, we need to accumulate entities rather than replace them
  const [entities, setEntities] = useState<Record<string, any>[]>([])
  const [accumulatedEntities, setAccumulatedEntities] = useState<Record<string, any>[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [views, setViews] = useState<EntityViewConfig[]>([])
  const [currentView, setCurrentView] = useState<EntityViewConfig | null>(null)
  
  // Ref to track previous values for dependency comparison
  const prevRef = useRef<PrevValues>({
    page: 1,
    pageSize: initialPageSize,
    sortBy: initialSortBy || '',
    sortOrder: initialSortOrder,
    view: null
  })
  
  // Pagination state
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(initialPageSize)
  const [totalCount, setTotalCount] = useState(0)
  const [totalPages, setTotalPages] = useState(0)
  const [hasMore, setHasMore] = useState(false)
  
  // Sorting and filtering state
  const [sortBy, setSortBy] = useState(initialSortBy || '')
  const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>(initialSortOrder)
  const [filters, setFilters] = useState<Record<string, any>>(initialFilters)

  // Load available views for the entity type
  useEffect(() => {
    const loadViews = async () => {
      try {
        const response = await apiClient.getEntityViews(entityType)
        if (response.data) {
          setViews(response.data.views)
          // Set default view
          if (response.data.views.length > 0) {
            const defaultView = view 
              ? response.data.views.find(v => v.name === view)
              : response.data.views[0]
            if (defaultView) {
              setCurrentView(defaultView)
              if (!initialSortBy && defaultView.defaultSort) {
                setSortBy(defaultView.defaultSort)
                setSortOrder(defaultView.defaultOrder as 'asc' | 'desc')
              }
            }
          }
        }
      } catch (err) {
        console.error('Failed to load entity views:', err)
      }
    }

    loadViews()
  }, [entityType, view, initialSortBy])

  // Fetch entities
  const fetchEntities = useCallback(async (resetPage = false, appendResults = false) => {
    if (!currentView) return

    setLoading(true)
    setError(null)

    try {
      const request: EntityQueryRequest = {
        entityType,
        page: resetPage ? 1 : page,
        pageSize,
        sortBy,
        sortOrder,
        filters,
        view: currentView.name
      }
      
      console.log('Fetching entities with request:', request) // Debug log

      const response = await apiClient.queryEntities(request)
      
      if (response.data) {
        setEntities(response.data.entities)
        setTotalCount(response.data.totalCount)
        setTotalPages(response.data.totalPages)
        setHasMore(response.data.hasMore)
        if (resetPage) {
          setPage(1)
        }
      } else if (response.error) {
        setError(response.error)
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch entities')
    } finally {
      setLoading(false)
    }
  }, [entityType, currentView, page, pageSize, sortBy, sortOrder, filters])

  // Load entities when dependencies change
  useEffect(() => {
    if (currentView) {
      fetchEntities(false) // Don't reset page when dependencies change
    }
  }, [currentView, page, pageSize, sortBy, sortOrder, filters])

  // Change page
  const goToPage = useCallback((newPage: number) => {
    if (newPage >= 1 && newPage <= totalPages) {
      setPage(newPage)
    }
  }, [totalPages])

  // Change page size
  const changePageSize = useCallback((newPageSize: number) => {
    setPageSize(newPageSize)
    setPage(1) // Reset to first page
  }, [])

  // Change sorting
  const changeSorting = useCallback((newSortBy: string) => {
    if (sortBy === newSortBy) {
      // Toggle sort order if same column
      setSortOrder(sortOrder === 'asc' ? 'desc' : 'asc')
    } else {
      // New column, default to ascending
      setSortBy(newSortBy)
      setSortOrder('asc')
    }
    // Reset entities and page for new sorting
    setEntities([])
    setPage(1)
  }, [sortBy, sortOrder])

  // Change filters
  const changeFilters = useCallback((newFilters: Record<string, any>) => {
    setFilters(newFilters)
    // Reset entities and page for new filters
    setEntities([])
    setPage(1)
  }, [])

  // Update prevRef when values change
  useEffect(() => {
    prevRef.current = {
      page,
      pageSize,
      sortBy,
      sortOrder,
      view: currentView?.name || null
    }
  }, [page, pageSize, sortBy, sortOrder, currentView])

  // Change view
  const changeView = useCallback((viewName: string) => {
    const newView = views.find(v => v.name === viewName)
    if (newView) {
      setCurrentView(newView)
      if (newView.defaultSort && newView.defaultSort !== sortBy) {
        setSortBy(newView.defaultSort)
        setSortOrder(newView.defaultOrder as 'asc' | 'desc')
      }
    }
  }, [views, sortBy])

  // Refresh data
  const refresh = useCallback(() => {
    fetchEntities(true)
  }, [fetchEntities])

  // Load more (for infinite scroll)
  const loadMore = useCallback(() => {
    if (hasMore && !loading) {
      const nextPage = page + 1;
      setPage(nextPage);
      
      // We need to fetch entities directly here instead of relying on useEffect
      // to avoid race conditions with the page state update
      if (currentView) {
        const request: EntityQueryRequest = {
          entityType,
          page: nextPage,
          pageSize,
          sortBy,
          sortOrder,
          filters,
          view: currentView.name
        };
        
        setLoading(true);
        apiClient.queryEntities(request).then(response => {
          // Use optional chaining to safely access properties
          const entities = response.data?.entities || [];
          const totalCount = response.data?.totalCount || 0;
          const totalPages = response.data?.totalPages || 0;
          const hasMore = response.data?.hasMore || false;
          
          setEntities(prev => [...prev, ...entities]);
          setTotalCount(totalCount);
          setTotalPages(totalPages);
          setHasMore(hasMore);
          setLoading(false);
        }).catch(err => {
          setLoading(false);
          console.error('Failed to load more entities:', err);
        });
      }
    }
  }, [hasMore, loading, page, pageSize, sortBy, sortOrder, filters, currentView, entityType])

  return {
    // Data
    entities,
    views,
    currentView,
    
    // Pagination
    page,
    pageSize,
    totalCount,
    totalPages,
    hasMore,
    
    // Sorting and filtering
    sortBy,
    sortOrder,
    filters,
    
    // State
    loading,
    error,
    
    // Actions
    goToPage,
    changePageSize,
    changeSorting,
    changeFilters,
    changeView,
    refresh,
    loadMore,
    
    // Computed
    startIndex: (page - 1) * pageSize + 1,
    endIndex: Math.min(page * pageSize, totalCount)
  }
} 