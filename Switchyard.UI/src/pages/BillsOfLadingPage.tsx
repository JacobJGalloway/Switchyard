import { useState, useEffect } from 'react'
import { useApiClients } from '../hooks/useApiClients'
import type { BillOfLading, LineEntry } from '../types'
import styles from './BillsOfLadingPage.module.css'

interface ModalState {
  transactionId: string
  entries: LineEntry[]
}

export default function BillsOfLadingPage() {
  const { logistics } = useApiClients()
  const [bols, setBols] = useState<BillOfLading[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [modal, setModal] = useState<ModalState | null>(null)
  const [modalLoading, setModalLoading] = useState(false)

  useEffect(() => {
    logistics.get<BillOfLading[]>('/api/BillOfLading')
      .then(setBols)
      .catch(() => setError('Failed to load bills of lading.'))
      .finally(() => setLoading(false))
  }, [])

  async function openModal(transactionId: string) {
    setModalLoading(true)
    setModal({ transactionId, entries: [] })
    try {
      const entries = await logistics.get<LineEntry[]>(`/api/BillOfLading/${transactionId}/line-entry`)
      setModal({ transactionId, entries })
    } catch {
      setModal({ transactionId, entries: [] })
    } finally {
      setModalLoading(false)
    }
  }

  return (
    <main className={styles.page}>
      <h2>Bills of Lading</h2>

      {loading && <p>Loading...</p>}
      {error && <p style={{ color: 'red' }}>{error}</p>}

      {!loading && !error && bols.length === 0 && (
        <p>No bills of lading found.</p>
      )}

      {bols.length > 0 && (
        <table className={styles.table}>
          <thead>
            <tr>
              {['Transaction ID', 'Status', 'Customer', 'Destination', ''].map(h => (
                <th key={h} className={styles.th}>{h}</th>
              ))}
            </tr>
          </thead>
          <tbody>
            {bols.map(bol => (
              <tr key={bol.partitionKey}>
                <td className={styles.td}>{bol.transactionId}</td>
                <td className={styles.td}>{bol.status}</td>
                <td className={styles.td}>{bol.customerFirstName} {bol.customerLastName}</td>
                <td className={styles.td}>{bol.city}, {bol.state}</td>
                <td className={styles.tdCenter}>
                  <button className={styles.viewBtn} onClick={() => openModal(bol.transactionId)}>View</button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}

      {modal && (
        <div className={styles.overlay}>
          <div className={styles.modal}>
            <button
              className={styles.closeBtn}
              onClick={() => setModal(null)}
              aria-label="Close"
            >✕</button>

            <h2 className={styles.modalTitle}>
              Line Entries — {modal.transactionId}
            </h2>

            {modalLoading && <p>Loading...</p>}

            {!modalLoading && modal.entries.length === 0 && (
              <p>No line entries found.</p>
            )}

            {!modalLoading && modal.entries.length > 0 && (
              <table className={styles.table} style={{ marginTop: 0 }}>
                <thead>
                  <tr>
                    {['Location ID', 'SKU', 'Qty', 'Processed', 'Processed Date'].map(h => (
                      <th key={h} className={styles.th}>{h}</th>
                    ))}
                  </tr>
                </thead>
                <tbody>
                  {modal.entries.map(entry => (
                    <tr key={entry.partitionKey}>
                      <td className={styles.td}>{entry.locationId}</td>
                      <td className={styles.td}>{entry.skuMarker}</td>
                      <td className={styles.td}>{entry.quantity}</td>
                      <td className={styles.td}>{entry.isProcessed ? 'Yes' : 'No'}</td>
                      <td className={styles.td}>
                        {entry.processedDate
                          ? new Date(entry.processedDate).toLocaleDateString()
                          : '—'}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            )}

            <div className={styles.modalFooter}>
              <button onClick={() => setModal(null)}>Close</button>
            </div>
          </div>
        </div>
      )}
    </main>
  )
}
