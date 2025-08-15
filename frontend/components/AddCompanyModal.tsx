'use client'

import React, { useState } from 'react'
import { Modal } from './Modal'
import { apiClient } from '@/services/api'
import { PicklistSelect } from './PicklistSelect'
import { useIndustries, useCompanySizes } from '@/hooks/usePicklist'

interface AddCompanyModalProps {
  isOpen: boolean
  onClose: () => void
  onSuccess: () => void
}

export function AddCompanyModal({ isOpen, onClose, onSuccess }: AddCompanyModalProps) {
  const [formData, setFormData] = useState({
    name: '',
    website: '',
    domain: '',
    industryId: '',
    sizeId: '',
    revenue: '',
    description: ''
  })
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState('')

  // Use picklist hooks
  const { items: industries, loading: industriesLoading } = useIndustries()
  const { items: companySizes, loading: sizesLoading } = useCompanySizes()

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setIsLoading(true)
    setError('')

    try {
      const companyData: any = {
        name: formData.name,
        website: formData.website || undefined,
        domain: formData.domain || undefined,
        revenue: formData.revenue ? parseFloat(formData.revenue) : undefined
      }

      // Only add industryId and sizeId if they have actual values
      if (formData.industryId && formData.industryId.trim() !== '') {
        companyData.industryId = formData.industryId
      }
      if (formData.sizeId && formData.sizeId.trim() !== '') {
        companyData.sizeId = formData.sizeId
      }

      const response = await apiClient.createCompany(companyData)

      if (response.data) {
        onSuccess()
        onClose()
        // Reset form
        setFormData({
          name: '',
          website: '',
          domain: '',
          industryId: '',
          sizeId: '',
          revenue: '',
          description: ''
        })
      } else {
        setError(response.error || 'Failed to create company')
      }
    } catch (error) {
      setError('An error occurred while creating the company')
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
    <Modal isOpen={isOpen} onClose={onClose} title="Add New Company">
      <form onSubmit={handleSubmit} className="space-y-4">
        {error && (
          <div className="bg-red-50 border border-red-200 text-red-600 px-3 py-2 rounded-md text-sm">
            {error}
          </div>
        )}

        <div>
          <label htmlFor="name" className="block text-sm font-medium text-gray-700 mb-1">
            Company Name *
          </label>
          <input
            type="text"
            id="name"
            name="name"
            required
            value={formData.name}
            onChange={handleChange}
            className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            placeholder="Enter company name"
          />
        </div>

        <div className="grid grid-cols-2 gap-4">
          <div>
            <label htmlFor="website" className="block text-sm font-medium text-gray-700 mb-1">
              Website
            </label>
            <input
              type="url"
              id="website"
              name="website"
              value={formData.website}
              onChange={handleChange}
              className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
              placeholder="https://example.com"
            />
          </div>

          <div>
            <label htmlFor="domain" className="block text-sm font-medium text-gray-700 mb-1">
              Domain
            </label>
            <input
              type="text"
              id="domain"
              name="domain"
              value={formData.domain}
              onChange={handleChange}
              className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
              placeholder="example.com"
            />
          </div>
        </div>

        <div>
          <label htmlFor="revenue" className="block text-sm font-medium text-gray-700 mb-1">
            Annual Revenue
          </label>
          <input
            type="number"
            id="revenue"
            name="revenue"
            value={formData.revenue}
            onChange={handleChange}
            className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            placeholder="1000000"
          />
        </div>

        <div className="grid grid-cols-2 gap-4">
          <div>
            <PicklistSelect
              items={industries}
              value={formData.industryId || null}
              onChange={(value) => setFormData(prev => ({ ...prev, industryId: value as string || '' }))}
              placeholder="Select industry (can add later)"
              label="Industry (optional)"
              loading={industriesLoading}
              searchable={true}
            />
          </div>

          <div>
            <PicklistSelect
              items={companySizes}
              value={formData.sizeId || null}
              onChange={(value) => setFormData(prev => ({ ...prev, sizeId: value as string || '' }))}
              placeholder="Select size (can add later)"
              label="Company Size (optional)"
              loading={sizesLoading}
              searchable={true}
            />
          </div>
        </div>

        <div className="text-sm text-gray-500 bg-blue-50 p-3 rounded-md">
          <p>💡 <strong>Tip:</strong> Industry and Company Size are now powered by dynamic picklists with search capabilities. You can search through the available options or select from the dropdown. Only Company Name is required to get started.</p>
        </div>

        <div>
          <label htmlFor="description" className="block text-sm font-medium text-gray-700 mb-1">
            Description
          </label>
          <textarea
            id="description"
            name="description"
            rows={3}
            value={formData.description}
            onChange={handleChange}
            className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            placeholder="Brief description of the company"
          />
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
            {isLoading ? 'Creating...' : 'Create Company'}
          </button>
        </div>
      </form>
    </Modal>
  )
} 