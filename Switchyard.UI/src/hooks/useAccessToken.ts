import { useAuth0 } from '@auth0/auth0-react'

export function useAccessToken() {
  const { getAccessTokenSilently, isAuthenticated } = useAuth0()

  const getToken = async (): Promise<string | null> => {
    if (!isAuthenticated) return null

    try {
      return await getAccessTokenSilently()
    } catch (error) {
      console.error('Failed to retrieve access token:', error)
      return null
    }
  }

  return { getToken }
}
