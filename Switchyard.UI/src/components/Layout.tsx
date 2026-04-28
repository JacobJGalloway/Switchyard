import Header from './Header'
import NavMenu from './NavMenu'

export default function Layout({ children }: { children: React.ReactNode }) {
  return (
    <>
      <Header />
      <NavMenu />
      <hr style={{ margin: 0, border: 'none', borderTop: '2px solid var(--text-h)' }} />
      {children}
    </>
  )
}
