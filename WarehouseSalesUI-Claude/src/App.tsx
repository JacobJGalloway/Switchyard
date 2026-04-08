import { Routes, Route, Navigate } from 'react-router-dom'
import ProtectedRoute from './components/ProtectedRoute'
import Layout from './components/Layout'
import LoginPage from './pages/LoginPage'
import DashboardPage from './pages/DashboardPage'
import InventoryPage from './pages/InventoryPage'
import BillsOfLadingPage from './pages/BillsOfLadingPage'
import UsersPage from './pages/UsersPage'

export default function App() {
  return (
    <Routes>
      <Route path="/login" element={<LoginPage />} />
      <Route path="/" element={
        <ProtectedRoute><Layout><DashboardPage /></Layout></ProtectedRoute>
      } />
      <Route path="/inventory" element={
        <ProtectedRoute permission="canReadInventory">
          <Layout><InventoryPage /></Layout>
        </ProtectedRoute>
      } />
      <Route path="/bills-of-lading" element={
        <ProtectedRoute permission="canReadBOL">
          <Layout><BillsOfLadingPage /></Layout>
        </ProtectedRoute>
      } />
      <Route path="/users" element={
        <ProtectedRoute permission="canManageUsers">
          <Layout><UsersPage /></Layout>
        </ProtectedRoute>
      } />
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  )
}
