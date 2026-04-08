import { createInventoryClient, createLogisticsClient } from '../api/apiClient'
import { useAccessToken } from './useAccessToken'

export function useApiClients() {
  const { getToken } = useAccessToken()

return {
    inventory: createInventoryClient(getToken),
    logistics: createLogisticsClient(getToken),
  }
}
