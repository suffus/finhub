'use client'

import { useState } from 'react'
import { EntityView } from './EntityView'

// Export the EntityType so it can be reused elsewhere
export type EntityType = 'companies' | 'contacts' | 'leads' | 'deals' | 
                  'industries' | 'companysizes' | 'leadstatuses' | 'leadtemperatures'

interface EntityPageProps {
  params: {
    entityType: EntityType
  }
}

// Mapping of entity types to display names
const entityDisplayNames: Record<EntityType, string> = {
  companies: 'Companies',
  contacts: 'Contacts',
  leads: 'Leads',
  deals: 'Deals',
  industries: 'Industries',
  companysizes: 'Company Sizes',
  leadstatuses: 'Lead Statuses',
  leadtemperatures: 'Lead Temperatures'
}

export function EntityPage({ params }: EntityPageProps) {
  const { entityType } = params
  const [selectedEntity, setSelectedEntity] = useState<Record<string, any> | null>(null)
  
  const displayName = entityDisplayNames[entityType as EntityType] || 
                      entityType.charAt(0).toUpperCase() + entityType.slice(1)
  
  // Handle entity selection for viewing details
  const handleEntityClick = (entity: Record<string, any>) => {
    setSelectedEntity(entity)
  }
  
  // Handle back button click to return to list
  const handleBackClick = () => {
    setSelectedEntity(null)
  }

  return (
    <div className="h-screen flex flex-col">
      <div className="flex-1 overflow-hidden">
        {selectedEntity ? (
          // Show entity details when selected
          <div className="h-full">
            {/* This would be a detail view component, not implemented yet */}
            <div className="bg-white h-full p-6">
              <button 
                onClick={handleBackClick}
                className="mb-4 px-3 py-1 text-sm flex items-center text-blue-600 hover:text-blue-800"
              >
                ← Back to list
              </button>
              <h2 className="text-2xl font-semibold mb-4">
                {selectedEntity.name || selectedEntity.firstName + ' ' + selectedEntity.lastName || `${displayName} Details`}
              </h2>
              
              <pre className="bg-gray-100 p-4 rounded text-sm overflow-auto max-h-96">
                {JSON.stringify(selectedEntity, null, 2)}
              </pre>
            </div>
          </div>
        ) : (
          // Show entity list
          <EntityView 
            entityType={entityType} 
            title={displayName}
            onEntityClick={handleEntityClick}
          />
        )}
      </div>
    </div>
  )
}