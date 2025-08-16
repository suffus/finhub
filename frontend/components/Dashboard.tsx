'use client'

import { useState, useEffect } from 'react'
import { TrendingUp, Users, Target, DollarSign, Activity, Building2, Phone, Mail, Calendar } from 'lucide-react'
import { apiClient, DashboardStats } from '@/services/api'
import { AddCompanyModal } from './AddCompanyModal'
import { AddContactModal } from './AddContactModal'
import { EntityType } from './EntityPage'

interface DashboardProps {
  onEntitySelect?: (entityType: EntityType) => void
}

export function Dashboard({ onEntitySelect }: DashboardProps) {
  const [stats, setStats] = useState<DashboardStats | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [showCompanyModal, setShowCompanyModal] = useState(false)
  const [showContactModal, setShowContactModal] = useState(false)

  useEffect(() => {
    fetchDashboardStats()
  }, [])

  const fetchDashboardStats = async () => {
    try {
      setLoading(true)
      const response = await apiClient.getDashboardStats()
      if (response.data) {
        setStats(response.data)
      } else {
        setError(response.error || 'Failed to load dashboard data')
      }
    } catch (error) {
      setError('Failed to load dashboard data')
    } finally {
      setLoading(false)
    }
  }

  const handleCompanyCreated = () => {
    // Refresh dashboard stats to show the new company
    fetchDashboardStats()
  }

  const handleContactCreated = () => {
    // Refresh dashboard stats to show the new contact
    fetchDashboardStats()
  }

  const dashboardStats = [
    { title: 'Total Deals', value: stats?.totalDeals?.toString() || '0', change: '0%', changeType: 'neutral' as 'neutral' | 'positive' | 'negative', icon: Target, color: 'bg-blue-500', description: 'Active deals in pipeline' },
    { title: 'Total Leads', value: stats?.totalLeads?.toString() || '0', change: '0%', changeType: 'neutral' as 'neutral' | 'positive' | 'negative', icon: Users, color: 'bg-green-500', description: 'Prospects in the system' },
    { title: 'Total Companies', value: stats?.totalCompanies?.toString() || '0', change: '0%', changeType: 'neutral' as 'neutral' | 'positive' | 'negative', icon: Building2, color: 'bg-purple-500', description: 'Business entities' },
    { title: 'Total Contacts', value: stats?.totalContacts?.toString() || '0', change: '0%', changeType: 'neutral' as 'neutral' | 'positive' | 'negative', icon: Users, color: 'bg-orange-500', description: 'People in the system' },
  ]

  const pipelineStages = stats?.pipelineStages || [
    { name: 'Qualification', count: 0, color: 'bg-gray-500' },
    { name: 'Proposal', count: 0, color: 'bg-blue-500' },
    { name: 'Negotiation', count: 0, color: 'bg-yellow-500' },
    { name: 'Closed Won', count: 0, color: 'bg-green-500' },
  ]

  const quickActions = [
    { name: 'Add Company', icon: Building2, onClick: () => setShowCompanyModal(true), color: 'bg-blue-600 hover:bg-blue-700' },
    { name: 'Add Contact', icon: Users, onClick: () => setShowContactModal(true), color: 'bg-green-600 hover:bg-green-700' },
    { name: 'Create Lead', icon: Target, onClick: () => {}, color: 'bg-purple-600 hover:bg-purple-700' },
    { name: 'New Deal', icon: DollarSign, onClick: () => {}, color: 'bg-orange-600 hover:bg-orange-700' },
  ]

  const recentActivity = [
    { type: 'contact', message: 'New contact added', time: '2 hours ago', icon: Users },
    { type: 'company', message: 'Company updated', time: '4 hours ago', icon: Building2 },
    { type: 'deal', message: 'Deal moved to next stage', time: '1 day ago', icon: Target },
    { type: 'lead', message: 'Lead status changed', time: '2 days ago', icon: TrendingUp },
  ]

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="bg-red-50 border border-red-200 rounded-md p-4">
        <div className="flex">
          <div className="text-red-600">
            <p className="text-sm">{error}</p>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="p-6 space-y-6">
      {/* Dashboard Stats */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        {dashboardStats.map((stat, index) => (
          <div key={index} className="bg-white rounded-lg shadow p-6">
            <div className="flex items-center">
              <div className={`p-3 rounded-lg ${stat.color}`}>
                <stat.icon className="h-6 w-6 text-white" />
              </div>
              <div className="ml-4">
                <p className="text-sm font-medium text-gray-600">{stat.title}</p>
                <p className="text-2xl font-semibold text-gray-900">{stat.value}</p>
                <p className="text-xs text-gray-500">{stat.description}</p>
              </div>
            </div>
          </div>
        ))}
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Pipeline Overview */}
        <div className="lg:col-span-2 bg-white rounded-lg shadow p-6">
          <h3 className="text-lg font-medium text-gray-900 mb-4">Sales Pipeline</h3>
          <div className="space-y-3">
            {pipelineStages.map((stage, index) => (
              <div key={index} className="flex items-center justify-between">
                <div className="flex items-center">
                  <div className={`w-3 h-3 rounded-full ${stage.color} mr-3`}></div>
                  <span className="text-sm font-medium text-gray-700">{stage.name}</span>
                </div>
                <span className="text-sm text-gray-500">{stage.count}</span>
              </div>
            ))}
          </div>
        </div>

        {/* Quick Actions */}
        <div className="bg-white rounded-lg shadow p-6">
          <h3 className="text-lg font-medium text-gray-900 mb-4">Quick Actions</h3>
          <div className="space-y-3">
            {quickActions.map((action, index) => (
              <button
                key={index}
                onClick={action.onClick}
                className={`w-full ${action.color} text-white px-4 py-2 rounded-md text-sm font-medium transition-colors duration-200 flex items-center`}
              >
                <action.icon className="h-4 w-4 mr-2" />
                {action.name}
              </button>
            ))}
          </div>
        </div>
      </div>

      {/* Recent Activity */}
      <div className="bg-white rounded-lg shadow p-6">
        <h3 className="text-lg font-medium text-gray-900 mb-4">Recent Activity</h3>
        <div className="space-y-4">
          {recentActivity.map((activity, index) => (
            <div key={index} className="flex items-center space-x-3">
              <div className="flex-shrink-0">
                <activity.icon className="h-5 w-5 text-gray-400" />
              </div>
              <div className="flex-1 min-w-0">
                <p className="text-sm font-medium text-gray-900">{activity.message}</p>
                <p className="text-sm text-gray-500">{activity.time}</p>
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* Modals */}
      <AddCompanyModal
        isOpen={showCompanyModal}
        onClose={() => setShowCompanyModal(false)}
        onSuccess={handleCompanyCreated}
      />
      
      <AddContactModal
        isOpen={showContactModal}
        onClose={() => setShowContactModal(false)}
        onSuccess={handleContactCreated}
      />
    </div>
  )
} 