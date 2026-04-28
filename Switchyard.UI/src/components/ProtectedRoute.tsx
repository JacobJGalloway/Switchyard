import { useAuth0 } from '@auth0/auth0-react'
import { Navigate } from 'react-router-dom'
import { usePermissions, type Permissions } from '../hooks/usePermissions'

interface ProtectedRouteProps {
  children: React.ReactNode
  permission?: keyof Omit<Permissions, 'role'>
}

export default function ProtectedRoute({ children, permission }: ProtectedRouteProps) {
  const { isAuthenticated, isLoading } = useAuth0()
  const permissions = usePermissions()

  if (isLoading) return <div>Loading...</div>
  if (!isAuthenticated) return <Navigate to="/login" replace />

  if (permission && !permissions[permission]) {
    return <Navigate to="/" replace />
  }

  return <>{children}</>
}
