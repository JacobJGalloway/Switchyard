import { useAuth0 } from '@auth0/auth0-react'

export default function DashboardPage() {
  const { user } = useAuth0()
  return (
    <main style={{ padding: '1rem', textAlign: 'left' }}>
      <h1>Hello, {user?.name}</h1>
      <p>Company reports coming soon.</p>
    </main>
  )
}
