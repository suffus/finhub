const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api'

export interface ApiResponse<T> {
  data?: T
  error?: string
  message?: string
}

export interface User {
  id: string
  email: string
  firstName: string
  lastName: string
  avatar?: string
  roleId?: string
  isActive: boolean
  tenantId: string
  createdAt: string
  updatedAt: string
}

export interface Company {
  id: string
  name: string
  website?: string
  domain?: string
  industryId?: string
  sizeId?: string
  revenue?: number
  externalId?: string
  tenantId: string
  createdAt: string
  updatedAt: string
  isDeleted: boolean
  industry?: PicklistItem
  size?: PicklistItem
}

export interface PicklistItem {
  id: string
  name: string
  code: string
  description?: string
  isActive: boolean
}

export interface PicklistResponse {
  items: PicklistItem[]
  totalCount: number
  hasMore: boolean
}

export interface PicklistSearchRequest {
  query: string
  limit: number
  offset: number
  entityType: string
}

// Entity management interfaces
export interface EntityQueryRequest {
  entityType: string
  page: number
  pageSize: number
  sortBy?: string
  sortOrder?: 'asc' | 'desc'
  filters?: Record<string, any>
  view?: string
}

export interface EntityQueryResponse {
  entities: Record<string, any>[]
  totalCount: number
  page: number
  pageSize: number
  totalPages: number
  hasMore: boolean
  sortBy?: string
  sortOrder?: string
}

export interface EntityViewConfig {
  name: string
  displayName: string
  columns: Column[]
  defaultSort: string
  defaultOrder: string
}

export interface Column {
  key: string
  label: string
  type: string
  sortable: boolean
  filterable: boolean
  width: string
  align?: string
  format?: string
}

export interface Contact {
  id: string
  firstName: string
  lastName: string
  title?: string
  department?: string
  companyId?: string
  originalSource?: string
  emailOptIn: boolean
  smsOptIn: boolean
  callOptIn: boolean
  tenantId: string
  createdAt: string
  updatedAt: string
  isDeleted: boolean
}

export interface Lead {
  id: string
  firstName?: string
  lastName?: string
  title?: string
  statusId?: string
  temperatureId?: string
  source?: string
  campaign?: string
  score: number
  companyId?: string
  contactId?: string
  assignedUserId?: string
  convertedAt?: string
  convertedToDealId?: string
  tenantId: string
  createdAt: string
  updatedAt: string
  isDeleted: boolean
}

export interface Deal {
  id: string
  name: string
  amount?: number
  currency: string
  probability: number
  pipelineId: string
  stageId: string
  expectedCloseDate?: string
  actualCloseDate?: string
  companyId?: string
  contactId?: string
  assignedUserId?: string
  tenantId: string
  createdAt: string
  updatedAt: string
  isDeleted: boolean
}

export interface DashboardStats {
  totalDeals: number
  totalLeads: number
  totalCompanies: number
  totalContacts: number
  pipelineStages: Array<{
    name: string
    count: number
    color: string
  }>
}

export interface AuthResponse {
  token: string
  user: User
}

class ApiClient {
  private baseUrl: string
  private token: string | null = null

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl
    // Try to get token from localStorage on initialization
    if (typeof window !== 'undefined') {
      this.token = localStorage.getItem('auth_token')
    }
  }

  setToken(token: string) {
    this.token = token
    if (typeof window !== 'undefined') {
      localStorage.setItem('auth_token', token)
    }
  }

  clearToken() {
    this.token = null
    if (typeof window !== 'undefined') {
      localStorage.removeItem('auth_token')
    }
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<ApiResponse<T>> {
    try {
      const url = `${this.baseUrl}${endpoint}`
      const headers: Record<string, string> = {
        'Content-Type': 'application/json',
        ...(options.headers as Record<string, string> || {}),
      }

      // Add authorization header if token exists
      if (this.token) {
        headers.Authorization = `Bearer ${this.token}`
      }

      const response = await fetch(url, {
        headers,
        ...options,
      })

      if (response.status === 401) {
        // Token expired or invalid
        this.clearToken()
        return { error: 'Authentication required' }
      }

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}))
        throw new Error(errorData.message || `HTTP error! status: ${response.status}`)
      }

      const data = await response.json()
      return { data }
    } catch (error) {
      console.error('API request failed:', error)
      return { error: error instanceof Error ? error.message : 'Unknown error' }
    }
  }

  // Authentication
  async login(email: string, password: string): Promise<ApiResponse<AuthResponse>> {
    const response = await this.request<AuthResponse>('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    })

    if (response.data?.token) {
      this.setToken(response.data.token)
    }

    return response
  }

  async register(email: string, password: string, firstName: string, lastName: string): Promise<ApiResponse<AuthResponse>> {
    const response = await this.request<AuthResponse>('/auth/register', {
      method: 'POST',
      body: JSON.stringify({ email, password, firstName, lastName }),
    })

    if (response.data?.token) {
      this.setToken(response.data.token)
    }

    return response
  }

  async getCurrentUser(): Promise<ApiResponse<User>> {
    return this.request<User>('/users/me')
  }

  // Dashboard - This will be implemented when backend supports it
  async getDashboardStats(): Promise<ApiResponse<DashboardStats>> {
    // For now, we'll aggregate data from individual endpoints
    try {
      const [companiesRes, contactsRes, leadsRes, dealsRes] = await Promise.all([
        this.getCompanies(),
        this.getContacts(),
        this.getLeads(),
        this.getDeals()
      ])

      const totalCompanies = companiesRes.data?.length || 0
      const totalContacts = contactsRes.data?.length || 0
      const totalLeads = leadsRes.data?.length || 0
      const totalDeals = dealsRes.data?.length || 0

      return {
        data: {
          totalDeals,
          totalLeads,
          totalCompanies,
          totalContacts,
          pipelineStages: [
            { name: 'Prospecting', count: 0, color: 'bg-blue-500' },
            { name: 'Qualification', count: 0, color: 'bg-yellow-500' },
            { name: 'Proposal', count: 0, color: 'bg-orange-500' },
            { name: 'Negotiation', count: 0, color: 'bg-purple-500' },
            { name: 'Closed Won', count: 0, color: 'bg-green-500' }
          ]
        }
      }
    } catch (error) {
      return { error: 'Failed to fetch dashboard data' }
    }
  }

  // Companies
  async getCompanies(): Promise<ApiResponse<Company[]>> {
    return this.request<Company[]>('/companies')
  }

  async createCompany(company: {
    name: string
    website?: string
    domain?: string
    industryId?: string
    sizeId?: string
    revenue?: number
  }): Promise<ApiResponse<Company>> {
    return this.request<Company>('/companies', {
      method: 'POST',
      body: JSON.stringify(company),
    })
  }

  async getCompany(id: string): Promise<ApiResponse<Company>> {
    return this.request<Company>(`/companies/${id}`)
  }

  async updateCompany(id: string, company: Partial<Company>): Promise<ApiResponse<Company>> {
    return this.request<Company>(`/companies/${id}`, {
      method: 'PUT',
      body: JSON.stringify(company),
    })
  }

  async deleteCompany(id: string): Promise<ApiResponse<void>> {
    return this.request<void>(`/companies/${id}`, {
      method: 'DELETE',
    })
  }

  // Contacts
  async getContacts(): Promise<ApiResponse<Contact[]>> {
    return this.request<Contact[]>('/contacts')
  }

  async createContact(contact: {
    firstName: string
    lastName: string
    title?: string
    department?: string
    companyId?: string
    originalSource?: string
    emailOptIn: boolean
    smsOptIn: boolean
    callOptIn: boolean
  }): Promise<ApiResponse<Contact>> {
    return this.request<Contact>('/contacts', {
      method: 'POST',
      body: JSON.stringify(contact),
    })
  }

  async getContact(id: string): Promise<ApiResponse<Contact>> {
    return this.request<Contact>(`/contacts/${id}`)
  }

  async updateContact(id: string, contact: Partial<Contact>): Promise<ApiResponse<Contact>> {
    return this.request<Contact>(`/contacts/${id}`, {
      method: 'PUT',
      body: JSON.stringify(contact),
    })
  }

  async deleteContact(id: string): Promise<ApiResponse<void>> {
    return this.request<void>(`/contacts/${id}`, {
      method: 'DELETE',
    })
  }

  // Leads
  async getLeads(): Promise<ApiResponse<Lead[]>> {
    return this.request<Lead[]>('/leads')
  }

  async createLead(lead: Omit<Lead, 'id' | 'createdAt' | 'updatedAt'>): Promise<ApiResponse<Lead>> {
    return this.request<Lead>('/leads', {
      method: 'POST',
      body: JSON.stringify(lead),
    })
  }

  async getLead(id: string): Promise<ApiResponse<Lead>> {
    return this.request<Lead>(`/leads/${id}`)
  }

  async updateLead(id: string, lead: Partial<Lead>): Promise<ApiResponse<Lead>> {
    return this.request<Lead>(`/leads/${id}`, {
      method: 'PUT',
      body: JSON.stringify(lead),
    })
  }

  async deleteLead(id: string): Promise<ApiResponse<void>> {
    return this.request<void>(`/leads/${id}`, {
      method: 'DELETE',
    })
  }

  // Deals
  async getDeals(): Promise<ApiResponse<Deal[]>> {
    return this.request<Deal[]>('/deals')
  }

  async createDeal(deal: Omit<Deal, 'id' | 'createdAt' | 'updatedAt'>): Promise<ApiResponse<Deal>> {
    return this.request<Deal>('/deals', {
      method: 'POST',
      body: JSON.stringify(deal),
    })
  }

  async getDeal(id: string): Promise<ApiResponse<Deal>> {
    return this.request<Deal>(`/deals/${id}`)
  }

  async updateDeal(id: string, deal: Partial<Deal>): Promise<ApiResponse<Deal>> {
    return this.request<Deal>(`/deals/${id}`, {
      method: 'PUT',
      body: JSON.stringify(deal),
    })
  }

  async deleteDeal(id: string): Promise<ApiResponse<void>> {
    return this.request<void>(`/deals/${id}`, {
      method: 'DELETE',
    })
  }

  // Picklists
async getPicklist(entityType: string): Promise<ApiResponse<PicklistResponse>> {
  return this.request<PicklistResponse>(`/picklists/${entityType}`)
}

async searchPicklist(request: PicklistSearchRequest): Promise<ApiResponse<PicklistResponse>> {
  return this.request<PicklistResponse>('/picklists/search', {
    method: 'POST',
    body: JSON.stringify(request),
  })
}

// Entities
async queryEntities(request: EntityQueryRequest): Promise<ApiResponse<EntityQueryResponse>> {
  return this.request<EntityQueryResponse>('/entities/query', {
    method: 'POST',
    body: JSON.stringify(request),
  })
}

async getEntityViews(entityType: string): Promise<ApiResponse<{ views: EntityViewConfig[] }>> {
  return this.request<{ views: EntityViewConfig[] }>(`/entities/${entityType}/views`)
}
}

export const apiClient = new ApiClient(API_BASE_URL) 