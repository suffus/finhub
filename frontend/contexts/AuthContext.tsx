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
        console.log("Token validation failed: User data not returned")
        apiClient.clearToken() // Clear the token if the response doesn't include user data
        return null
      }
    } catch (error) {
      console.error('Token validation failed:', error)
      apiClient.clearToken() // Clear the token on any error
      return null
    }
  }

  useEffect(() => {
    const initializeAuth = async () => {
      try {
        const token = localStorage.getItem('auth_token')
        // Try to get cached user data first
        const cachedUserData = localStorage.getItem('auth_user')
        
        if (token) {
          // Set the token in the API client
          apiClient.setToken(token)
          
          // First try to use cached user data
          if (cachedUserData) {
            try {
              const parsedUser = JSON.parse(cachedUserData)
              setUser(parsedUser)
              
              // Validate the token in the background
              validateToken(token).then(userData => {
                if (userData) {
                  // Update with fresh data if available
                  setUser(userData)
                  localStorage.setItem('auth_user', JSON.stringify(userData))
                } else {
                  // Token is invalid, remove everything
                  setUser(null)
                  console.log("Token is invalid, removing it (3)")
                  localStorage.removeItem('auth_token')
                  localStorage.removeItem('auth_user')
                  apiClient.clearToken()
                }
              })
            } catch (e) {
              console.error('Error parsing cached user data:', e)
              // Fallback to token validation
              const userData = await validateToken(token)
              if (userData) {
                setUser(userData)
                localStorage.setItem('auth_user', JSON.stringify(userData))
              } else {
                // Token is invalid, remove it
                console.log("Token is invalid, removing it (1)")
                localStorage.removeItem('auth_token')
                localStorage.removeItem('auth_user')
                apiClient.clearToken()
              }
            }
          } else {
            // No cached user data, validate token
            const userData = await validateToken(token)
            if (userData) {
              setUser(userData)
              localStorage.setItem('auth_user', JSON.stringify(userData))
            } else {
              // Token is invalid, remove it
              console.log("Token is invalid, removing it (2)")
              localStorage.removeItem('auth_token')
              apiClient.clearToken()
            }
          }
        }
      } catch (error) {
        console.error('Auth initialization failed:', error)
        // Clear invalid token and user data
        console.log("Auth initialization failed, removing token and user data")
        localStorage.removeItem('auth_token')
        localStorage.removeItem('auth_user')
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
        // Save the user data in state
        setUser(response.data.user)
        // Also save user data in localStorage for persistence
        localStorage.setItem('auth_user', JSON.stringify(response.data.user))
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
        // Save the user data in state
        setUser(response.data.user)
        // Also save user data in localStorage for persistence
        localStorage.setItem('auth_user', JSON.stringify(response.data.user))
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
    // Also remove user data from localStorage
    localStorage.removeItem('auth_user')
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