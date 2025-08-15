'use client'

import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react'
import { apiClient, AuthResponse, User } from '@/services/api'

interface AuthContextType {
  user: User | null
  isAuthenticated: boolean
  isLoading: boolean
  login: (email: string, password: string) => Promise<{ success: boolean; error?: string }>
  register: (email: string, password: string, firstName: string, lastName: string) => Promise<{ success: boolean; error?: string }>
  logout: () => void
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function useAuth() {
  const context = useContext(AuthContext)
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}

interface AuthProviderProps {
  children: ReactNode
}

export function AuthProvider({ children }: AuthProviderProps) {
  const [user, setUser] = useState<User | null>(null)
  const [isLoading, setIsLoading] = useState(true)

  // Function to validate token and get user info
  const validateToken = async (token: string): Promise<User | null> => {
    try {
      // Try to get current user info from the backend
      const response = await apiClient.getCurrentUser()
      if (response.data) {
        return response.data
      } else {
        return null
      }
    } catch (error) {
      console.error('Token validation failed:', error)
      return null
    }
  }

  useEffect(() => {
    const initializeAuth = async () => {
      try {
        const token = localStorage.getItem('auth_token')
        if (token) {
          // Set the token in the API client
          apiClient.setToken(token)
          
          // Validate the token by getting current user info
          const userData = await validateToken(token)
          if (userData) {
            setUser(userData)
          } else {
            // Token is invalid, remove it
            localStorage.removeItem('auth_token')
            apiClient.clearToken()
          }
        }
      } catch (error) {
        console.error('Auth initialization failed:', error)
        // Clear invalid token
        localStorage.removeItem('auth_token')
        apiClient.clearToken()
      } finally {
        setIsLoading(false)
      }
    }

    initializeAuth()
  }, [])

  const login = async (email: string, password: string) => {
    try {
      setIsLoading(true)
      const response = await apiClient.login(email, password)
      
      if (response.data) {
        setUser(response.data.user)
        return { success: true }
      } else {
        return { success: false, error: response.error || 'Login failed' }
      }
    } catch (error) {
      return { success: false, error: 'Login failed' }
    } finally {
      setIsLoading(false)
    }
  }

  const register = async (email: string, password: string, firstName: string, lastName: string) => {
    try {
      setIsLoading(true)
      const response = await apiClient.register(email, password, firstName, lastName)
      
      if (response.data) {
        setUser(response.data.user)
        return { success: true }
      } else {
        return { success: false, error: response.error || 'Registration failed' }
      }
    } catch (error) {
      return { success: false, error: 'Registration failed' }
    } finally {
      setIsLoading(false)
    }
  }

  const logout = () => {
    setUser(null)
    apiClient.clearToken()
  }

  const value: AuthContextType = {
    user,
    isAuthenticated: !!user,
    isLoading,
    login,
    register,
    logout,
  }

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  )
} 