'use client'

import { useState, useRef, useEffect } from 'react'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { EntityList } from './EntityList'
import { useEntityList } from '../hooks/useEntityList'
import { ChevronLeft, Settings } from 'lucide-react'

interface EntityViewProps {
  entityType: string
  title?: string
  onBack?: () => void
  onEntityClick?: (entity: Record<string, any>) => void
}

export function EntityView({ entityType, title, onBack, onEntityClick }: EntityViewProps) {
  const [activeView, setActiveView] = useState<string | null>(null)
  const listRef = useRef<HTMLDivElement>(null)
  
  // Get entity data and views
  const {
    views,
    currentView,
    changeView,
    hasMore,
    loading,
    loadMore
  } = useEntityList({
    entityType
  })

  // Set active view when views load
  useEffect(() => {
    if (views.length > 0 && !activeView) {
      setActiveView(views[0].name)
    }
  }, [views, activeView])

  // Handle infinite scroll
  useEffect(() => {
    const handleScroll = () => {
      if (!listRef.current || loading || !hasMore) return

      const { scrollTop, clientHeight, scrollHeight } = listRef.current
      
      // Load more when user scrolls to bottom (with a buffer of 200px)
      if (scrollTop + clientHeight >= scrollHeight - 200) {
        loadMore()
      }
    }

    const currentRef = listRef.current
    if (currentRef) {
      currentRef.addEventListener('scroll', handleScroll)
    }
    
    return () => {
      if (currentRef) {
        currentRef.removeEventListener('scroll', handleScroll)
      }
    }
  }, [loading, hasMore, loadMore])

  // Display name for the entity type
  const displayName = title || `${entityType.charAt(0).toUpperCase() + entityType.slice(1)}`

  return (
    <div className="h-full flex flex-col">
      {/* Header */}
      <div className="bg-white border-b border-gray-200 px-6 py-4 flex items-center justify-between">
        <div className="flex items-center space-x-4">
          {onBack && (
            <button 
              onClick={onBack} 
              className="p-2 rounded-full hover:bg-gray-100"
            >
              <ChevronLeft className="w-5 h-5 text-gray-500" />
              <span className="sr-only">Back to Dashboard</span>
            </button>
          )}
          <h1 className="text-2xl font-semibold text-gray-900">{displayName}</h1>
        </div>
        
        <div className="flex items-center space-x-2">
          <button className="px-4 py-2 border border-gray-300 rounded-md text-sm text-gray-700 hover:bg-gray-50 flex items-center space-x-1">
            <Settings className="w-4 h-4" />
            <span>Settings</span>
          </button>
        </div>
      </div>

      {/* Views Tabs */}
      {views.length > 0 && (
        <div className="border-b border-gray-200 bg-white">
          <Tabs 
            value={activeView || views[0].name} 
            onValueChange={tab => {
              setActiveView(tab)
              changeView(tab)
            }}
            className="px-6"
          >
            <TabsList>
              {views.map(view => (
                <TabsTrigger key={view.name} value={view.name}>
                  {view.displayName}
                </TabsTrigger>
              ))}
            </TabsList>
          </Tabs>
        </div>
      )}

      {/* Content Area with Infinite Scroll */}
      <div 
        ref={listRef}
        className="flex-1 overflow-y-auto"
      >
        {activeView && (
          <div className="p-6">
            <EntityList 
              entityType={entityType}
              title={undefined} // No title needed as we have the header
              onEntityClick={onEntityClick}
            />
            
            {/* Loading indicator for infinite scroll */}
            {loading && hasMore && (
              <div className="flex justify-center p-4">
                <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  )
}