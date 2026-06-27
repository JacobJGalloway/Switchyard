import { useAuth0 } from '@auth0/auth0-react'
import SKUMovementChart from '../components/SKUMovementChart'
import styles from './DashboardPage.module.css'

export default function DashboardPage() {
  const { user } = useAuth0()
  return (
    <main className={styles.page}>
      <h1 className={styles.greeting}>Hello, {user?.name}</h1>
      <SKUMovementChart />
    </main>
  )
}
