import { useState, useRef, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { useAuth0 } from '@auth0/auth0-react'
import { UserCircle, Sun, Moon } from 'lucide-react'
import { useTheme } from '../contexts/ThemeContext'
import { resolveAsset, isClientOverride } from '../utils/assetResolver'
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

  const nameOnly = resolveAsset('logo-name-only', theme)
  const clientName = (import.meta.env.VITE_CLIENT_NAME as string | undefined) ?? 'Switchyard'

  return (
    <header className={styles.header}>
      <Link to="/" className={styles.logoLink}>
        {nameOnly
          ? (
            <>
              <img src={resolveAsset('logo-detail', theme)!} alt="" className={styles.icon} />
              <img src={nameOnly} alt={clientName} className={styles.name} />
            </>
          )
          : <img src={resolveAsset('logo-full-name', theme)!} alt={clientName} className={styles.logo} />
        }
        {isClientOverride && (
          <img src={resolveAsset('logo-powered-by', theme)!} alt="Powered by Switchyard" className={styles.poweredBy} />
        )}
      </Link>

      <div className={styles.controls}>
        <button
          className={styles.themeButton}
          onClick={toggleTheme}
          aria-label={theme === 'dark' ? 'Switch to light mode' : 'Switch to dark mode'}
          title={theme === 'dark' ? 'Switch to light mode' : 'Switch to dark mode'}
        >
          {theme === 'dark' ? <Sun size={18} /> : <Moon size={18} />}
        </button>

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
                onClick={() => logout({ logoutParams: { returnTo: window.location.origin + '/login' } })}
              >
                Log Out
              </button>
            </div>
          )}
        </div>
      </div>
    </header>
  )
}
