import { useAuth0 } from '@auth0/auth0-react'
import styles from './DashboardPage.module.css'

export default function DashboardPage() {
  const { user } = useAuth0()
  return (
    <main className={styles.page}>
      <h1 className={styles.greeting}>Hello, {user?.name}</h1>
      <p className={styles.coming}>Company charts and operational reports coming soon.</p>
    </main>
  )
}
