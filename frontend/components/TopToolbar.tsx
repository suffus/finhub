'use client'

import { Search, Bell, Bot, User, Menu, Settings, Plus, Filter, BarChart3, LogOut } from 'lucide-react'
import { useAuth } from '@/contexts/AuthContext'

interface TopToolbarProps {
  onMenuClick: () => void
}

export function TopToolbar({ onMenuClick }: TopToolbarProps) {
  const { user, logout } = useAuth()

  const handleLogout = () => {
    logout()
  }

  return (
    <div className="bg-white border-b border-gray-200 px-4 py-3 flex items-center justify-between">
      {/* Left side - Menu button and Search */}
      <div className="flex items-center space-x-4">
        <button
          onClick={onMenuClick}
          className="p-2 rounded-md hover:bg-gray-100 transition-colors"
        >
          <Menu className="h-5 w-5 text-gray-600" />
        </button>
        
        {/* Logo/Brand */}
        <div className="flex items-center space-x-2">
          <div className="h-8 w-8 bg-gradient-to-r from-blue-600 to-purple-600 rounded-lg flex items-center justify-center">
            <span className="text-white font-bold text-sm">FH</span>
          </div>
          <span className="text-lg font-semibold text-gray-900">FinHub CRM</span>
        </div>
        
        {/* Search Bar */}
        <div className="relative">
          <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" />
          <input
            type="text"
            placeholder="Search companies, contacts, deals, leads..."
            className="pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent w-96"
          />
          <button className="absolute right-2 top-1/2 transform -translate-y-1/2 p-1 hover:bg-gray-100 rounded">
            <Filter className="h-4 w-4 text-gray-400" />
          </button>
        </div>
      </div>

      {/* Right side - Actions and User */}
      <div className="flex items-center space-x-4">
        {/* Quick Actions */}
        <div className="flex items-center space-x-2">
          <button className="flex items-center space-x-2 bg-blue-600 text-white px-3 py-2 rounded-lg hover:bg-blue-700 transition-colors">
            <Plus className="h-4 w-4" />
            <span className="text-sm font-medium">New</span>
          </button>
        </div>

        {/* AI Assistant */}
        <button className="p-2 rounded-md hover:bg-gray-100 transition-colors relative group">
          <Bot className="h-5 w-5 text-blue-600" />
          <div className="absolute bottom-full right-0 mb-2 px-2 py-1 bg-gray-800 text-white text-xs rounded opacity-0 group-hover:opacity-100 transition-opacity whitespace-nowrap">
            AI Assistant
          </div>
        </button>

        {/* Analytics */}
        <button className="p-2 rounded-md hover:bg-gray-100 transition-colors relative group">
          <BarChart3 className="h-5 w-5 text-gray-600" />
          <div className="absolute bottom-full right-0 mb-2 px-2 py-1 bg-gray-800 text-white text-xs rounded opacity-0 group-hover:opacity-100 transition-opacity whitespace-nowrap">
            Analytics
          </div>
        </button>

        {/* Notifications */}
        <button className="p-2 rounded-md hover:bg-gray-100 transition-colors relative">
          <Bell className="h-5 w-5 text-gray-600" />
          <span className="absolute -top-1 -right-1 bg-red-500 text-white text-xs rounded-full h-4 w-4 flex items-center justify-center">
            0
          </span>
        </button>

        {/* Settings */}
        <button className="p-2 rounded-md hover:bg-gray-100 transition-colors relative group">
          <Settings className="h-5 w-5 text-gray-600" />
          <div className="absolute bottom-full right-0 mb-2 px-2 py-1 bg-gray-800 text-white text-xs rounded opacity-0 group-hover:opacity-100 transition-opacity whitespace-nowrap">
            Settings
          </div>
        </button>

        {/* User Profile */}
        <div className="flex items-center space-x-2 p-2 rounded-md hover:bg-gray-100 transition-colors cursor-pointer group">
          <div className="h-8 w-8 bg-gradient-to-r from-blue-600 to-purple-600 rounded-full flex items-center justify-center">
            <User className="h-4 w-4 text-white" />
          </div>
          <div className="text-left">
            <span className="text-sm font-medium text-gray-700 block">
              {user ? `${user.firstName} ${user.lastName}` : 'User'}
            </span>
            <span className="text-xs text-gray-500">
              {user ? user.email : 'FinHub CRM'}
            </span>
          </div>
          
          {/* Logout dropdown */}
          <div className="relative">
            <button
              onClick={handleLogout}
              className="p-1 rounded hover:bg-gray-200 transition-colors"
              title="Logout"
            >
              <LogOut className="h-4 w-4 text-gray-600" />
            </button>
          </div>
        </div>
      </div>
    </div>
  )
} 