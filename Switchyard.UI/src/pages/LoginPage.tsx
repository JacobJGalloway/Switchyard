import { useAuth0 } from '@auth0/auth0-react'
import { useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { useTheme } from '../contexts/ThemeContext'
import { resolveAsset } from '../utils/assetResolver'

export default function LoginPage() {
  const { loginWithRedirect, isAuthenticated, isLoading } = useAuth0()
  const navigate = useNavigate()
  const { theme } = useTheme()

  const clientName = (import.meta.env.VITE_CLIENT_NAME as string | undefined) ?? 'Switchyard'

  useEffect(() => {
    if (!isLoading && isAuthenticated) {
      navigate('/', { replace: true })
    }
  }, [isAuthenticated, isLoading, navigate])

  return (
    <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center', minHeight: '60vh', gap: '3rem' }}>
      <img src={resolveAsset('logo-full-name', theme)!} alt={`${clientName} Logistics`} style={{ width: '450px' }} />
      <button onClick={() => loginWithRedirect()}>Log In</button>
    </div>
  )
}
