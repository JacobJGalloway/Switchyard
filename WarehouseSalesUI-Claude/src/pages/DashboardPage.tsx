import { useAuth0 } from '@auth0/auth0-react'
import InventoryViewer from '../components/InventoryViewer'

export default function DashboardPage() {
  const { user, logout } = useAuth0()

  return (
    <div style={{ padding: '1rem' }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <h1>Dashboard</h1>
        <div>
          <span style={{ marginRight: '1rem' }}>{user?.name}</span>
          <button onClick={() => logout({ logoutParams: { returnTo: window.location.origin + '/login' } })}>
            Log Out
          </button>
        </div>
      </div>
      <InventoryViewer />
    </div>
  )
}
