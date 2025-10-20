import axios from 'axios'
import type { APIResponse, ManagedHub, ManagedCluster, HealthResponse } from '@/types'

const api = axios.create({
  baseURL: '/api',
  headers: {
    'Content-Type': 'application/json',
  },
})

// Add request interceptor to include auth token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Add response interceptor for error handling
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // Handle unauthorized - redirect to login or show error
      console.error('Unauthorized access')
    }
    return Promise.reject(error)
  }
)

export const healthAPI = {
  getHealth: async (): Promise<HealthResponse> => {
    const response = await api.get<HealthResponse>('/health')
    return response.data
  },
  getReady: async (): Promise<HealthResponse> => {
    const response = await api.get<HealthResponse>('/ready')
    return response.data
  },
  getLive: async (): Promise<HealthResponse> => {
    const response = await api.get<HealthResponse>('/live')
    return response.data
  },
}

export const hubsAPI = {
  listHubs: async (): Promise<ManagedHub[]> => {
    const response = await api.get<APIResponse<ManagedHub[]>>('/hubs')
    return response.data.data || []
  },
  getHub: async (name: string): Promise<ManagedHub> => {
    const response = await api.get<APIResponse<ManagedHub>>(`/hubs/${name}`)
    if (!response.data.data) {
      throw new Error('Hub not found')
    }
    return response.data.data
  },
  getHubClusters: async (hubName: string): Promise<ManagedCluster[]> => {
    const response = await api.get<APIResponse<ManagedCluster[]>>(`/hubs/${hubName}/clusters`)
    return response.data.data || []
  },
}

export default api

