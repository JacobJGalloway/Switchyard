import { useState, useEffect } from 'react'
import { useApiClients } from '../hooks/useApiClients'
import type { AppUser, CreateUserRequest } from '../types'
import styles from './UsersPage.module.css'

const ROLES = ['Employee', 'Manager', 'Admin'] as const

const emptyForm: CreateUserRequest = {
  email: '',
  password: '',
  firstName: '',
  lastName: '',
  role: 'Employee',
  warehouseId: '',
  storeId: '',
}

export default function UsersPage() {
  const { logistics } = useApiClients()
  const [users, setUsers] = useState<AppUser[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [showModal, setShowModal] = useState(false)
  const [form, setForm] = useState<CreateUserRequest>(emptyForm)
  const [submitting, setSubmitting] = useState(false)
  const [formError, setFormError] = useState('')

  function loadUsers() {
    setLoading(true)
    logistics.get<AppUser[]>('/api/User')
      .then(setUsers)
      .catch(() => setError('Failed to load users.'))
      .finally(() => setLoading(false))
  }

  useEffect(() => {
    loadUsers()
  }, [])

  function openModal() {
    setForm(emptyForm)
    setFormError('')
    setShowModal(true)
  }

  function closeModal() {
    setShowModal(false)
    setFormError('')
  }

  function handleChange(e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) {
    setForm(prev => ({ ...prev, [e.target.name]: e.target.value }))
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setSubmitting(true)
    setFormError('')
    try {
      await logistics.post('/api/User', form)
      closeModal()
      loadUsers()
    } catch {
      setFormError('Failed to create user. Check all fields and try again.')
    } finally {
      setSubmitting(false)
    }
  }

  async function handleDeactivate(userId: string) {
    try {
      await logistics.patch(`/api/User/${encodeURIComponent(userId)}/deactivate`, null)
      loadUsers()
    } catch {
      setError('Failed to deactivate user.')
    }
  }

  return (
    <main className={styles.page}>
      <div className={styles.toolbar}>
        <h2 style={{ margin: 0 }}>Users</h2>
        <button onClick={openModal}>Create User</button>
      </div>

      {loading && <p>Loading...</p>}
      {error && <p style={{ color: 'red' }}>{error}</p>}

      {!loading && !error && users.length === 0 && (
        <p>No users found.</p>
      )}

      {users.length > 0 && (
        <table className={styles.table}>
          <thead>
            <tr>
              {['Name', 'Email', 'Role', 'Warehouse', 'Store', 'Actions'].map(h => (
                <th key={h} className={styles.th}>{h}</th>
              ))}
            </tr>
          </thead>
          <tbody>
            {users.map(user => {
              const muted = user.blocked
              const td = muted ? styles.tdMuted : styles.td
              const tdC = muted ? styles.tdCenterMuted : styles.tdCenter
              return (
                <tr key={user.userId}>
                  <td className={td}>{user.name}</td>
                  <td className={td}>{user.email}</td>
                  <td className={td}>{user.role}</td>
                  <td className={td}>{user.warehouseId || '—'}</td>
                  <td className={td}>{user.storeId || '—'}</td>
                  <td className={tdC}>
                    {!user.blocked && (
                      <button
                        className={styles.deactivateBtn}
                        onClick={() => handleDeactivate(user.userId)}
                      >
                        Deactivate
                      </button>
                    )}
                  </td>
                </tr>
              )
            })}
          </tbody>
        </table>
      )}

      {showModal && (
        <div className={styles.overlay}>
          <div className={styles.modal}>
            <button
              className={styles.closeBtn}
              onClick={closeModal}
              aria-label="Close"
            >✕</button>

            <h2 className={styles.modalTitle}>Create User</h2>

            <form onSubmit={handleSubmit}>
              <div className={styles.formGrid}>
                <div className={styles.formField}>
                  <label htmlFor="firstName">First Name</label>
                  <input
                    id="firstName"
                    name="firstName"
                    value={form.firstName}
                    onChange={handleChange}
                    required
                  />
                </div>
                <div className={styles.formField}>
                  <label htmlFor="lastName">Last Name</label>
                  <input
                    id="lastName"
                    name="lastName"
                    value={form.lastName}
                    onChange={handleChange}
                    required
                  />
                </div>
                <div className={styles.formField}>
                  <label htmlFor="email">Email</label>
                  <input
                    id="email"
                    name="email"
                    type="email"
                    value={form.email}
                    onChange={handleChange}
                    required
                  />
                </div>
                <div className={styles.formField}>
                  <label htmlFor="password">Password</label>
                  <input
                    id="password"
                    name="password"
                    type="password"
                    value={form.password}
                    onChange={handleChange}
                    required
                  />
                </div>
                <div className={styles.formField}>
                  <label htmlFor="role">Role</label>
                  <select
                    id="role"
                    name="role"
                    value={form.role}
                    onChange={handleChange}
                    required
                  >
                    {ROLES.map(r => (
                      <option key={r} value={r}>{r}</option>
                    ))}
                  </select>
                </div>
                <div className={styles.formField}>
                  <label htmlFor="warehouseId">Warehouse ID</label>
                  <input
                    id="warehouseId"
                    name="warehouseId"
                    value={form.warehouseId}
                    onChange={handleChange}
                  />
                </div>
                <div className={styles.formField}>
                  <label htmlFor="storeId">Store ID</label>
                  <input
                    id="storeId"
                    name="storeId"
                    value={form.storeId}
                    onChange={handleChange}
                  />
                </div>
              </div>

              {formError && <p style={{ color: 'red', marginTop: '0.75rem' }}>{formError}</p>}

              <div className={styles.modalFooter}>
                <button type="button" onClick={closeModal} disabled={submitting}>Cancel</button>
                <button type="submit" disabled={submitting}>
                  {submitting ? 'Creating...' : 'Create'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </main>
  )
}
