'use client'

import React, { useState } from 'react'
import { EntityList } from '../../components/EntityList'
import { Building2, Users, Target, TrendingUp } from 'lucide-react'

const entityTypes = [
  { key: 'companies', label: 'Companies', icon: Building2, description: 'Manage your business relationships' },
  { key: 'contacts', label: 'Contacts', icon: Users, description: 'Track individual contacts and prospects' },
  { key: 'leads', label: 'Leads', icon: Target, description: 'Monitor lead pipeline and conversions' },
  { key: 'deals', label: 'Deals', icon: TrendingUp, description: 'Track sales opportunities and revenue' },
]

export default function EntitiesPage() {
  const [selectedEntityType, setSelectedEntityType] = useState('companies')

  const handleEntityClick = (entity: Record<string, any>) => {
    console.log('Entity clicked:', entity)
    // Here you could navigate to a detail view or open a modal
  }

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900">Entity Management</h1>
          <p className="mt-2 text-gray-600">
            View and manage your CRM entities with advanced filtering, sorting, and pagination
          </p>
        </div>

        {/* Entity Type Selector */}
        <div className="mb-6">
          <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
            {entityTypes.map((entityType) => {
              const Icon = entityType.icon
              return (
                <button
                  key={entityType.key}
                  onClick={() => setSelectedEntityType(entityType.key)}
                  className={`p-4 rounded-lg border-2 transition-all ${
                    selectedEntityType === entityType.key
                      ? 'border-blue-500 bg-blue-50 text-blue-700'
                      : 'border-gray-200 bg-white text-gray-700 hover:border-gray-300 hover:bg-gray-50'
                  }`}
                >
                  <div className="flex items-center space-x-3">
                    <Icon className="w-6 h-6" />
                    <div className="text-left">
                      <div className="font-medium">{entityType.label}</div>
                      <div className="text-sm opacity-75">{entityType.description}</div>
                    </div>
                  </div>
                </button>
              )
            })}
          </div>
        </div>

        {/* Entity List */}
        <EntityList
          entityType={selectedEntityType}
          title={`${entityTypes.find(e => e.key === selectedEntityType)?.label}`}
          onEntityClick={handleEntityClick}
        />

        {/* Features Overview */}
        <div className="mt-12 grid grid-cols-1 md:grid-cols-3 gap-6">
          <div className="bg-white p-6 rounded-lg shadow">
            <h3 className="text-lg font-semibold text-gray-900 mb-3">Advanced Filtering</h3>
            <p className="text-gray-600">
              Filter entities by various criteria including search terms, date ranges, and custom fields.
            </p>
          </div>
          
          <div className="bg-white p-6 rounded-lg shadow">
            <h3 className="text-lg font-semibold text-gray-900 mb-3">Smart Sorting</h3>
            <p className="text-gray-600">
              Click any sortable column header to sort data. Toggle between ascending and descending order.
            </p>
          </div>
          
          <div className="bg-white p-6 rounded-lg shadow">
            <h3 className="text-lg font-semibold text-gray-900 mb-3">Multiple Views</h3>
            <p className="text-gray-600">
              Switch between different view configurations to see the data that matters most to you.
            </p>
          </div>
        </div>
      </div>
    </div>
  )
} 