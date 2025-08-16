'use client'

import { 
  Users, 
  Building2, 
  Target, 
  TrendingUp, 
  Megaphone, 
  FileText, 
  Settings, 
  BarChart3,
  Bookmark,
  ChevronRight,
  Phone,
  Mail,
  Calendar,
  MapPin,
  Tag
} from 'lucide-react'

import { EntityType } from './EntityPage'

interface LeftSidebarProps {
  isOpen: boolean
  activeEntity?: EntityType | null
  onEntitySelect?: (entityType: EntityType) => void
}

export function LeftSidebar({ isOpen, activeEntity = null, onEntitySelect }: LeftSidebarProps) {
  // Map URL paths to entity types
  const entityTypeMap: Record<string, EntityType> = {
    '/companies': 'companies',
    '/contacts': 'contacts',
    '/leads': 'leads',
    '/deals': 'deals',
    '/industries': 'industries',
    '/sizes': 'companysizes',
    '/lead-statuses': 'leadstatuses',
    '/temperatures': 'leadtemperatures'
  }

  // Handle entity item click
  const handleEntityClick = (e: React.MouseEvent, href: string) => {
    e.preventDefault()
    
    const entityType = entityTypeMap[href]
    if (entityType && onEntitySelect) {
      onEntitySelect(entityType)
    }
  }

  const navigationItems = [
    {
      title: 'Core CRM',
      icon: Building2,
      items: [
        { name: 'Companies', href: '/companies', count: 0, description: 'Business entities' },
        { name: 'Contacts', href: '/contacts', count: 0, description: 'People records' },
        { name: 'Leads', href: '/leads', count: 0, description: 'Sales prospects' },
        { name: 'Deals', href: '/deals', count: 0, description: 'Sales opportunities' },
      ]
    },
    {
      title: 'Sales Pipeline',
      icon: Target,
      items: [
        { name: 'Pipeline', href: '/pipeline', count: 0, description: 'Sales stages' },
        { name: 'Forecasts', href: '/forecasts', count: 0, description: 'Revenue projections' },
        { name: 'Tasks', href: '/tasks', count: 0, description: 'Follow-up activities' },
      ]
    },
    {
      title: 'Marketing',
      icon: Megaphone,
      items: [
        { name: 'Lead Sources', href: '/sources', count: 0, description: 'Lead generation' },
        { name: 'Campaigns', href: '/campaigns', count: 0, description: 'Marketing efforts' },
        { name: 'Assets', href: '/assets', count: 0, description: 'Marketing materials' },
      ]
    },
    {
      title: 'Data Management',
      icon: Settings,
      items: [
        { name: 'Industries', href: '/industries', count: 0, description: 'Business sectors' },
        { name: 'Company Sizes', href: '/sizes', count: 0, description: 'Business scale' },
        { name: 'Lead Statuses', href: '/lead-statuses', count: 0, description: 'Lead stages' },
        { name: 'Lead Temperatures', href: '/temperatures', count: 0, description: 'Lead quality' },
      ]
    },
    {
      title: 'Communication',
      icon: Phone,
      items: [
        { name: 'Phone Numbers', href: '/phone-numbers', count: 0, description: 'Contact info' },
        { name: 'Email Addresses', href: '/emails', count: 0, description: 'Email contacts' },
        { name: 'Addresses', href: '/addresses', count: 0, description: 'Physical locations' },
        { name: 'Social Media', href: '/social', count: 0, description: 'Online presence' },
      ]
    },
    {
      title: 'Analytics',
      icon: BarChart3,
      items: [
        { name: 'Reports', href: '/reports', count: 0, description: 'Business insights' },
        { name: 'Activity Log', href: '/activity', count: 0, description: 'System activity' },
        { name: 'Performance', href: '/performance', count: 0, description: 'Team metrics' },
      ]
    }
  ]

  const bookmarks = [
    { name: 'Hot Leads', icon: TrendingUp, color: 'text-red-500', href: '/leads?filter=hot' },
    { name: 'This Month', icon: Calendar, color: 'text-blue-500', href: '/deals?period=month' },
    { name: 'Follow-ups', icon: Target, color: 'text-green-500', href: '/tasks?type=followup' },
    { name: 'New Companies', icon: Building2, color: 'text-purple-500', href: '/companies?filter=new' }
  ]

  if (!isOpen) {
    return (
      <div className="w-16 bg-white border-r border-gray-200 flex flex-col items-center py-4 space-y-4">
        {navigationItems.map((item, index) => (
          <div key={index} className="relative group">
            <button className="p-3 rounded-lg hover:bg-gray-100 transition-colors">
              <item.icon className="h-5 w-5 text-gray-600" />
            </button>
            <div className="absolute left-full ml-2 px-2 py-1 bg-gray-800 text-white text-sm rounded opacity-0 group-hover:opacity-100 transition-opacity whitespace-nowrap z-50">
              {item.title}
            </div>
          </div>
        ))}
      </div>
    )
  }

  return (
    <div className="w-64 bg-white border-r border-gray-200 flex flex-col">
      {/* Bookmarks Section */}
      <div className="p-4 border-b border-gray-200">
        <h3 className="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-3">
          Quick Access
        </h3>
        <div className="space-y-2">
          {bookmarks.map((bookmark, index) => (
            <a
              key={index}
              href={bookmark.href}
              className="w-full flex items-center space-x-3 px-3 py-2 rounded-lg hover:bg-gray-100 transition-colors text-left"
            >
              <bookmark.icon className={`h-4 w-4 ${bookmark.color}`} />
              <span className="text-sm text-gray-700">{bookmark.name}</span>
            </a>
          ))}
        </div>
      </div>

      {/* Navigation Menu */}
      <div className="flex-1 overflow-y-auto py-4">
        <nav className="space-y-1">
          {navigationItems.map((section, sectionIndex) => (
            <div key={sectionIndex} className="px-4">
              <div className="flex items-center space-x-2 mb-2">
                <section.icon className="h-4 w-4 text-gray-500" />
                <h3 className="text-xs font-semibold text-gray-500 uppercase tracking-wider">
                  {section.title}
                </h3>
              </div>
              <div className="space-y-1 mb-4">
                {section.items.map((item, itemIndex) => (
                  <a
                    key={itemIndex}
                    href={item.href}
                    onClick={(e) => entityTypeMap[item.href] && handleEntityClick(e, item.href)}
                    className={`flex items-start justify-between px-3 py-2 rounded-lg hover:bg-gray-100 transition-colors text-sm group
                      ${activeEntity === entityTypeMap[item.href] 
                        ? 'bg-blue-50 text-blue-700 border-l-4 border-blue-500' 
                        : 'text-gray-700'}
                    `}
                  >
                    <div className="flex-1">
                      <span className="block">{item.name}</span>
                      <span className="text-xs text-gray-500">{item.description}</span>
                    </div>
                    <div className="flex items-center space-x-2">
                      <span className="text-xs text-gray-500 bg-gray-100 px-2 py-1 rounded-full">
                        {item.count}
                      </span>
                      <ChevronRight className="h-3 w-3 text-gray-400 group-hover:translate-x-1 transition-transform" />
                    </div>
                  </a>
                ))}
              </div>
            </div>
          ))}
        </nav>
      </div>

      {/* System Info */}
      <div className="p-4 border-t border-gray-200 bg-gray-50">
        <div className="text-xs text-gray-500">
          <div className="flex items-center space-x-2 mb-2">
            <div className="w-2 h-2 bg-green-500 rounded-full"></div>
            <span>System Online</span>
          </div>
          <div className="text-gray-400">
            FinHub CRM v1.0
          </div>
        </div>
      </div>
    </div>
  )
} 