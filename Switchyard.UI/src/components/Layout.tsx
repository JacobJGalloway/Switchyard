import Header from './Header'
import NavMenu from './NavMenu'

export default function Layout({ children }: { children: React.ReactNode }) {
  return (
    <>
      <Header />
      <NavMenu />
      {children}
    </>
  )
}
