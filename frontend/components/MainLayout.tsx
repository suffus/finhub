'use client'

import { useState, useEffect } from 'react'
import { Dashboard } from '@/components/Dashboard'
import { EntityView } from '@/components/EntityView'
import { EntityType } from '@/components/EntityPage'

interface MainLayoutProps {
  initialEntity?: {
    type: EntityType
    id?: string
  }
  onEntityChange?: (entity: {type: EntityType, id?: string} | null) => void
}

export function MainLayout({ initialEntity, onEntityChange }: MainLayoutProps) {
  const [activeEntity, setActiveEntity] = useState<{type: EntityType, id?: string} | null>(
    initialEntity || null
  )

  // Update internal state when props change
  useEffect(() => {
    setActiveEntity(initialEntity || null)
  }, [initialEntity])

  const handleEntitySelect = (entityType: EntityType) => {
    const newEntity = { type: entityType }
    setActiveEntity(newEntity)
    if (onEntityChange) {
      onEntityChange(newEntity)
    }
  }

  const handleEntityBack = () => {
    setActiveEntity(null)
    if (onEntityChange) {
      onEntityChange(null)
    }
  }

  return (
    <>
      {activeEntity ? (
        // Show entity view when an entity is selected
        <EntityView 
          entityType={activeEntity.type}
          onBack={handleEntityBack}
        />
      ) : (
        // Show dashboard when no entity is selected
        <Dashboard onEntitySelect={handleEntitySelect} />
      )}
    </>
  )
}