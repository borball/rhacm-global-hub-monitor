import { describe, it, expect, vi, beforeEach } from 'vitest'
import axios from 'axios'
import { hubsAPI, healthAPI } from './api'

vi.mock('axios')

describe('API Services', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('healthAPI', () => {
    it('should fetch health status', async () => {
      const mockResponse = {
        data: { status: 'healthy', version: '1.0.0', timestamp: '2024-01-01T00:00:00Z' },
      }
      vi.mocked(axios.create).mockReturnValue({
        get: vi.fn().mockResolvedValue(mockResponse),
      } as any)

      const result = await healthAPI.getHealth()
      expect(result.status).toBe('healthy')
    })
  })

  describe('hubsAPI', () => {
    it('should list hubs', async () => {
      const mockResponse = {
        data: {
          success: true,
          data: [
            { name: 'hub1', status: 'Ready' },
            { name: 'hub2', status: 'Ready' },
          ],
        },
      }
      vi.mocked(axios.create).mockReturnValue({
        get: vi.fn().mockResolvedValue(mockResponse),
      } as any)

      const result = await hubsAPI.listHubs()
      expect(result).toHaveLength(2)
    })
  })
})

