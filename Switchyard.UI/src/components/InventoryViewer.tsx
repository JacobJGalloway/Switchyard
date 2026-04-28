import { useState, useEffect } from 'react'
import { useApiClients } from '../hooks/useApiClients'
import type { Warehouse, Store, InventoryItem, InventoryType, LocationType } from '../types'
import styles from './InventoryViewer.module.css'

const INVENTORY_TYPES: InventoryType[] = ['Clothing', 'PPE', 'Tool']
const PAGE_SIZE = 50

export default function InventoryViewer() {
  const { inventory, logistics } = useApiClients()

  const [locationType, setLocationType] = useState<LocationType | ''>('')
  const [warehouses, setWarehouses] = useState<Warehouse[]>([])
  const [stores, setStores] = useState<Store[]>([])
  const [selectedLocationId, setSelectedLocationId] = useState('')
  const [projectedFilter, setProjectedFilter] = useState<'All' | 'Yes' | 'No'>('All')
  const [results, setResults] = useState<Map<InventoryType, InventoryItem[]>>(new Map())
  const [activeTab, setActiveTab] = useState<InventoryType>('Clothing')
  const [currentPage, setCurrentPage] = useState(1)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  useEffect(() => {
    if (locationType === 'Warehouse') {
      logistics.get<Warehouse[]>('/api/Warehouse').then(setWarehouses).catch(() => setWarehouses([]))
    } else if (locationType === 'Store') {
      logistics.get<Store[]>('/api/Store').then(setStores).catch(() => setStores([]))
    }
    setSelectedLocationId('')
    setResults(new Map())
  }, [locationType])

  function handleTabChange(tab: InventoryType) {
    setActiveTab(tab)
    setCurrentPage(1)
  }

  async function handleSearch() {
    if (!selectedLocationId) return
    setLoading(true)
    setError('')
    try {
      const fetches = INVENTORY_TYPES.map(async type => {
        const items = await inventory.get<InventoryItem[]>(`/api/${type}/location/${selectedLocationId}`)
        return [type, items] as const
      })
      const entries = await Promise.all(fetches)
      const filtered = new Map(
        entries.map(([type, items]) => [
          type,
          projectedFilter === 'All' ? items : items.filter(i => i.projected === (projectedFilter === 'Yes')),
        ])
      )
      setResults(filtered)
      setCurrentPage(1)
    } catch {
      setError('Failed to load inventory.')
    } finally {
      setLoading(false)
    }
  }

  const locationOptions =
    locationType === 'Warehouse'
      ? warehouses.map(w => ({ id: w.warehouseId, label: `${w.warehouseId} — ${w.city}, ${w.state}` }))
      : stores.map(s => ({ id: s.storeId, label: `${s.storeId} — ${s.city}, ${s.state}` }))

  const activeItems = results.get(activeTab) ?? []
  const totalPages = Math.max(1, Math.ceil(activeItems.length / PAGE_SIZE))
  const pageItems = activeItems.slice((currentPage - 1) * PAGE_SIZE, currentPage * PAGE_SIZE)
  const hasResults = results.size > 0

  return (
    <div className={styles.wrapper}>
      <h2>Inventory Viewer</h2>

      <div className={styles.controls}>
        <label>
          Location Type
          <br />
          <select value={locationType} onChange={e => setLocationType(e.target.value as LocationType)}>
            <option value="">-- Select --</option>
            <option value="Warehouse">Warehouse</option>
            <option value="Store">Store</option>
          </select>
        </label>

        <label>
          Location
          <br />
          <select
            value={selectedLocationId}
            onChange={e => setSelectedLocationId(e.target.value)}
            disabled={!locationType}
          >
            <option value="">-- Select --</option>
            {locationOptions.map(o => (
              <option key={o.id} value={o.id}>{o.label}</option>
            ))}
          </select>
        </label>

        <label>
          Projected
          <br />
          <select value={projectedFilter} onChange={e => setProjectedFilter(e.target.value as 'All' | 'Yes' | 'No')}>
            <option value="All">All</option>
            <option value="Yes">Yes</option>
            <option value="No">No</option>
          </select>
        </label>

        <button onClick={handleSearch} disabled={!selectedLocationId || loading}>
          {loading ? 'Loading...' : 'Search'}
        </button>
      </div>

      {error && <p className={styles.error}>{error}</p>}

      {hasResults && (
        <>
          <div className={styles.tabs}>
            {INVENTORY_TYPES.map(type => (
              <button
                key={type}
                type="button"
                className={`${styles.tab} ${activeTab === type ? styles.tabActive : ''}`}
                onClick={() => handleTabChange(type)}
              >
                {type}
                <span className={styles.tabCount}>({results.get(type)?.length ?? 0})</span>
              </button>
            ))}
          </div>

          {pageItems.length === 0 ? (
            <p>No items found.</p>
          ) : (
            <>
              <table className={styles.table}>
                <thead>
                  <tr>
                    {['SKU', 'Unloaded Date', 'Projected', 'Location'].map(h => (
                      <th key={h} className={styles.th}>{h}</th>
                    ))}
                  </tr>
                </thead>
                <tbody>
                  {pageItems.map(item => (
                    <tr key={item.partitionKey}>
                      <td className={styles.td}>{item.skuMarker}</td>
                      <td className={styles.td}>{new Date(item.unloadedDate).toLocaleDateString()}</td>
                      <td className={styles.td}>{item.projected ? 'Yes' : 'No'}</td>
                      <td className={styles.td}>{item.locationId}</td>
                    </tr>
                  ))}
                </tbody>
              </table>

              {totalPages > 1 && (
                <div className={styles.pagination}>
                  <button onClick={() => setCurrentPage(p => p - 1)} disabled={currentPage === 1}>
                    Prev
                  </button>
                  <span className={styles.pageInfo}>Page {currentPage} of {totalPages}</span>
                  <button onClick={() => setCurrentPage(p => p + 1)} disabled={currentPage === totalPages}>
                    Next
                  </button>
                </div>
              )}
            </>
          )}
        </>
      )}
    </div>
  )
}
