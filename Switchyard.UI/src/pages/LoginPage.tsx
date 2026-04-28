import { useAuth0 } from '@auth0/auth0-react'
import { useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import fullLogoLight from '../assets/Digital Parts Full Logo Light Mode.svg'

export default function LoginPage() {
  const { loginWithRedirect, isAuthenticated, isLoading } = useAuth0()
  const navigate = useNavigate()

  useEffect(() => {
    if (!isLoading && isAuthenticated) {
      navigate('/', { replace: true })
    }
  }, [isAuthenticated, isLoading, navigate])

  return (
    <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center', minHeight: '60vh', gap: '3rem' }}>
      <img src={fullLogoLight} alt="Digital Parts Logistics" style={{ width: '450px' }} />
      <button onClick={() => loginWithRedirect()}>Log In</button>
    </div>
  )
}
