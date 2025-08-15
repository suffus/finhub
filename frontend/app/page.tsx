'use client'

import { useState } from 'react'
import { TopToolbar } from '@/components/TopToolbar'
import { LeftSidebar } from '@/components/LeftSidebar'
import { Dashboard } from '@/components/Dashboard'
import { Login } from '@/components/Login'
import { useAuth } from '@/contexts/AuthContext'

export default function Home() {
  const [sidebarOpen, setSidebarOpen] = useState(true)
  const { isAuthenticated, isLoading } = useAuth()

  // Show loading state
  if (isLoading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
          <p className="text-gray-600">Loading FinHub CRM...</p>
        </div>
      </div>
    )
  }

  // Show login if not authenticated
  if (!isAuthenticated) {
    return <Login />
  }

  // Show main dashboard if authenticated
  return (
    <div className="h-screen flex flex-col bg-gray-50">
      {/* Top Toolbar */}
      <TopToolbar onMenuClick={() => setSidebarOpen(!sidebarOpen)} />
      
      <div className="flex flex-1 overflow-hidden">
        {/* Left Sidebar */}
        <LeftSidebar isOpen={sidebarOpen} />
        
        {/* Main Content */}
        <main className="flex-1 overflow-auto">
          <Dashboard />
        </main>
      </div>
    </div>
  )
} 