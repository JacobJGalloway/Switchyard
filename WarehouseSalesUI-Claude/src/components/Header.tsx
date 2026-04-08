import { useState, useRef, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { useAuth0 } from '@auth0/auth0-react'
import { UserCircle } from 'lucide-react'
import { useTheme } from '../contexts/ThemeContext'
import iconLight from '../assets/Digital Parts Icon Light Mode.svg'
import iconDark from '../assets/Digital Parts Icon Dark Mode.svg'
import nameLight from '../assets/Digital Parts Name Light Mode.svg'
import nameDark from '../assets/Digital Parts Name Dark Mode.svg'
import styles from './Header.module.css'

export default function Header() {
  const { logout } = useAuth0()
  const { theme, toggleTheme } = useTheme()
  const [open, setOpen] = useState(false)
  const wrapperRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    function handleClickOutside(e: MouseEvent) {
      if (wrapperRef.current && !wrapperRef.current.contains(e.target as Node)) {
        setOpen(false)
      }
    }
    document.addEventListener('mousedown', handleClickOutside)
    return () => document.removeEventListener('mousedown', handleClickOutside)
  }, [])

  return (
    <header className={styles.header}>
      <Link to="/" className={styles.logo}>
        <img
          src={theme === 'dark' ? iconDark : iconLight}
          alt="Digital Parts icon"
          className={styles.icon}
        />
        <img
          src={theme === 'dark' ? nameDark : nameLight}
          alt="Digital Parts"
          className={styles.name}
        />
      </Link>

      <div className={styles.profileWrapper} ref={wrapperRef}>
        <button
          className={styles.profileButton}
          onClick={() => setOpen(o => !o)}
          aria-label="Profile menu"
        >
          <UserCircle size={32} />
        </button>

        {open && (
          <div className={styles.dropdown}>
            <button
              className={styles.dropdownItem}
              onClick={() => { toggleTheme(); setOpen(false) }}
            >
              {theme === 'light' ? 'Switch to Dark Mode' : 'Switch to Light Mode'}
            </button>
            <hr className={styles.dropdownDivider} />
            <button
              className={styles.dropdownItem}
              onClick={() => logout({ logoutParams: { returnTo: window.location.origin + '/login' } })}
            >
              Log Out
            </button>
          </div>
        )}
      </div>
    </header>
  )
}
