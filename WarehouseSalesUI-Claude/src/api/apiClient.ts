const inventoryApiUrl = import.meta.env.VITE_INVENTORY_API_URL
const logisticsApiUrl = import.meta.env.VITE_LOGISTICS_API_URL

async function request<T>(
  baseUrl: string,
  path: string,
  getToken: () => Promise<string | null>,
  options: RequestInit = {}
): Promise<T> {
  const token = await getToken()

  const response = await fetch(`${baseUrl}${path}`, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
      ...options.headers,
    },
  })

  if (!response.ok) {
    throw new Error(`${response.status} ${response.statusText}`)
  }

  return response.json() as Promise<T>
}

export function createInventoryClient(getToken: () => Promise<string | null>) {
  return {
    get: <T>(path: string) => request<T>(inventoryApiUrl, path, getToken),
    post: <T>(path: string, body: unknown) =>
      request<T>(inventoryApiUrl, path, getToken, { method: 'POST', body: JSON.stringify(body) }),
    put: <T>(path: string, body: unknown) =>
      request<T>(inventoryApiUrl, path, getToken, { method: 'PUT', body: JSON.stringify(body) }),
    patch: <T>(path: string, body: unknown) =>
      request<T>(inventoryApiUrl, path, getToken, { method: 'PATCH', body: JSON.stringify(body) }),
    delete: <T>(path: string) =>
      request<T>(inventoryApiUrl, path, getToken, { method: 'DELETE' }),
  }
}

export function createLogisticsClient(getToken: () => Promise<string | null>) {
  return {
    get: <T>(path: string) => request<T>(logisticsApiUrl, path, getToken),
    post: <T>(path: string, body: unknown) =>
      request<T>(logisticsApiUrl, path, getToken, { method: 'POST', body: JSON.stringify(body) }),
    put: <T>(path: string, body: unknown) =>
      request<T>(logisticsApiUrl, path, getToken, { method: 'PUT', body: JSON.stringify(body) }),
    patch: <T>(path: string, body: unknown) =>
      request<T>(logisticsApiUrl, path, getToken, { method: 'PATCH', body: JSON.stringify(body) }),
    delete: <T>(path: string) =>
      request<T>(logisticsApiUrl, path, getToken, { method: 'DELETE' }),
  }
}
