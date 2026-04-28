import { NavLink } from 'react-router-dom'
import { usePermissions } from '../hooks/usePermissions'
import styles from './NavMenu.module.css'

export default function NavMenu() {
  const { canReadInventory, canReadBOL, canManageUsers } = usePermissions()

  return (
    <nav className={styles.nav}>
      {canReadInventory && (
        <NavLink
          to="/inventory"
          className={({ isActive }) =>
            isActive ? `${styles.link} ${styles.active}` : styles.link
          }
        >
          Inventory
        </NavLink>
      )}
      {canReadBOL && (
        <NavLink
          to="/bills-of-lading"
          className={({ isActive }) =>
            isActive ? `${styles.link} ${styles.active}` : styles.link
          }
        >
          Bills of Lading
        </NavLink>
      )}
      {canManageUsers && (
        <NavLink
          to="/users"
          className={({ isActive }) =>
            isActive ? `${styles.link} ${styles.active}` : styles.link
          }
        >
          Users
        </NavLink>
      )}
    </nav>
  )
}
