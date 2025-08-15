'use client'

import React, { useState, useEffect } from 'react'
import { Modal } from './Modal'
import { apiClient, Company } from '@/services/api'

interface AddContactModalProps {
  isOpen: boolean
  onClose: () => void
  onSuccess: () => void
}

export function AddContactModal({ isOpen, onClose, onSuccess }: AddContactModalProps) {
  const [formData, setFormData] = useState({
    firstName: '',
    lastName: '',
    title: '',
    department: '',
    companyId: '',
    originalSource: '',
    emailOptIn: false,
    smsOptIn: false,
    callOptIn: false
  })
  const [companies, setCompanies] = useState<Company[]>([])
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState('')

  // Load companies for the dropdown
  useEffect(() => {
    if (isOpen) {
      loadCompanies()
    }
  }, [isOpen])

  const loadCompanies = async () => {
    try {
      const response = await apiClient.getCompanies()
      if (response.data) {
        setCompanies(response.data)
      }
    } catch (error) {
      console.error('Failed to load companies:', error)
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setIsLoading(true)
    setError('')

    try {
      const response = await apiClient.createContact({
        firstName: formData.firstName,
        lastName: formData.lastName,
        title: formData.title || undefined,
        department: formData.department || undefined,
        companyId: formData.companyId || undefined,
        originalSource: formData.originalSource || undefined,
        emailOptIn: formData.emailOptIn,
        smsOptIn: formData.smsOptIn,
        callOptIn: formData.callOptIn
      })

      if (response.data) {
        onSuccess()
        onClose()
        // Reset form
        setFormData({
          firstName: '',
          lastName: '',
          title: '',
          department: '',
          companyId: '',
          originalSource: '',
          emailOptIn: false,
          smsOptIn: false,
          callOptIn: false
        })
      } else {
        setError(response.error || 'Failed to create contact')
      }
    } catch (error) {
      setError('An error occurred while creating the contact')
    } finally {
      setIsLoading(false)
    }
  }

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) => {
    setFormData(prev => ({
      ...prev,
      [e.target.name]: e.target.value
    }))
  }

  return (
    <Modal isOpen={isOpen} onClose={onClose} title="Add New Contact">
      <form onSubmit={handleSubmit} className="space-y-4">
        {error && (
          <div className="bg-red-50 border border-red-200 text-red-600 px-3 py-2 rounded-md text-sm">
            {error}
          </div>
        )}

        <div className="grid grid-cols-2 gap-4">
          <div>
            <label htmlFor="firstName" className="block text-sm font-medium text-gray-700 mb-1">
              First Name *
            </label>
            <input
              type="text"
              id="firstName"
              name="firstName"
              required
              value={formData.firstName}
              onChange={handleChange}
              className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
              placeholder="Enter first name"
            />
          </div>

          <div>
            <label htmlFor="lastName" className="block text-sm font-medium text-gray-700 mb-1">
              Last Name *
            </label>
            <input
              type="text"
              id="lastName"
              name="lastName"
              required
              value={formData.lastName}
              onChange={handleChange}
              className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
              placeholder="Enter last name"
            />
          </div>
        </div>

        <div className="grid grid-cols-2 gap-4">
          <div>
            <label htmlFor="title" className="block text-sm font-medium text-gray-700 mb-1">
              Job Title
            </label>
            <input
              type="text"
              id="title"
              name="title"
              value={formData.title}
              onChange={handleChange}
              className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
              placeholder="e.g., Sales Manager"
            />
          </div>

          <div>
            <label htmlFor="department" className="block text-sm font-medium text-gray-700 mb-1">
              Department
            </label>
            <input
              type="text"
              id="department"
              name="department"
              value={formData.department}
              onChange={handleChange}
              className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
              placeholder="e.g., Sales"
            />
          </div>
        </div>

        <div>
          <label htmlFor="companyId" className="block text-sm font-medium text-gray-700 mb-1">
            Company
          </label>
          <select
            id="companyId"
            name="companyId"
            value={formData.companyId}
            onChange={handleChange}
            className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
          >
            <option value="">Select company (optional)</option>
            {companies.map(company => (
              <option key={company.id} value={company.id}>
                {company.name}
              </option>
            ))}
          </select>
        </div>

        <div className="grid grid-cols-2 gap-4">
          <div>
            <label htmlFor="title" className="block text-sm font-medium text-gray-700 mb-1">
              Job Title
            </label>
            <input
              type="text"
              id="title"
              name="title"
              value={formData.title}
              onChange={handleChange}
              className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
              placeholder="e.g., Sales Manager"
            />
          </div>

          <div>
            <label htmlFor="department" className="block text-sm font-medium text-gray-700 mb-1">
              Department
            </label>
            <input
              type="text"
              id="department"
              name="department"
              value={formData.department}
              onChange={handleChange}
              className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
              placeholder="e.g., Sales"
            />
          </div>
        </div>

        <div>
          <label htmlFor="originalSource" className="block text-sm font-medium text-gray-700 mb-1">
            Source
          </label>
          <input
            type="text"
            id="originalSource"
            name="originalSource"
            value={formData.originalSource}
            onChange={handleChange}
            className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            placeholder="e.g., Website, Referral, Trade Show"
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Communication Preferences
          </label>
          <div className="space-y-2">
            <label className="flex items-center">
              <input
                type="checkbox"
                name="emailOptIn"
                checked={formData.emailOptIn}
                onChange={(e) => setFormData(prev => ({ ...prev, emailOptIn: e.target.checked }))}
                className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
              />
              <span className="ml-2 text-sm text-gray-700">Email communications</span>
            </label>
            <label className="flex items-center">
              <input
                type="checkbox"
                name="smsOptIn"
                checked={formData.smsOptIn}
                onChange={(e) => setFormData(prev => ({ ...prev, smsOptIn: e.target.checked }))}
                className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
              />
              <span className="ml-2 text-sm text-gray-700">SMS communications</span>
            </label>
            <label className="flex items-center">
              <input
                type="checkbox"
                name="callOptIn"
                checked={formData.callOptIn}
                onChange={(e) => setFormData(prev => ({ ...prev, callOptIn: e.target.checked }))}
                className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
              />
              <span className="ml-2 text-sm text-gray-700">Phone calls</span>
            </label>
          </div>
        </div>



        <div className="flex justify-end space-x-3 pt-4">
          <button
            type="button"
            onClick={onClose}
            className="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
          >
            Cancel
          </button>
          <button
            type="submit"
            disabled={isLoading}
            className="px-4 py-2 text-sm font-medium text-white bg-blue-600 border border-transparent rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {isLoading ? 'Creating...' : 'Create Contact'}
          </button>
        </div>
      </form>
    </Modal>
  )
} 