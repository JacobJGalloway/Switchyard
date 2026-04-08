import { useState, useEffect, useRef } from 'react'
import { useApiClients } from '../hooks/useApiClients'
import type { Warehouse, Store, InventoryItem, InventoryType, LocationType } from '../types'

const INVENTORY_TYPES: InventoryType[] = ['Clothing', 'PPE', 'Tool']

function MultiSelectDropdown({
  options,
  selected,
  onToggle,
}: {
  options: InventoryType[]
  selected: InventoryType[]
  onToggle: (type: InventoryType) => void
}) {
  const [open, setOpen] = useState(false)
  const ref = useRef<HTMLDivElement>(null)

  useEffect(() => {
    function handleClickOutside(e: MouseEvent) {
      if (ref.current && !ref.current.contains(e.target as Node)) setOpen(false)
    }
    document.addEventListener('mousedown', handleClickOutside)
    return () => document.removeEventListener('mousedown', handleClickOutside)
  }, [])

  const label = selected.length === 0
    ? '-- Select --'
    : selected.length === options.length
      ? 'All Types'
      : selected.join(', ')

  return (
    <div ref={ref} style={{ position: 'relative', display: 'inline-block' }}>
      <button
        type="button"
        onClick={() => setOpen(o => !o)}
        style={{ minWidth: '160px', textAlign: 'left', padding: '0.2rem 0.4rem', cursor: 'pointer' }}
      >
        {label} ▾
      </button>
      {open && (
        <div style={{
          position: 'absolute', top: '100%', left: 0, zIndex: 10,
          background: 'var(--surface)', border: '1px solid var(--table-border)', borderRadius: '4px',
          padding: '0.4rem 0', minWidth: '160px', boxShadow: '0 2px 6px rgba(0,0,0,0.15)'
        }}>
          {options.map(type => (
            <label key={type} style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', padding: '0.3rem 0.8rem', cursor: 'pointer' }}>
              <input
                type="checkbox"
                checked={selected.includes(type)}
                onChange={() => onToggle(type)}
              />
              {type}
            </label>
          ))}
        </div>
      )}
    </div>
  )
}

export default function InventoryViewer() {
  const { inventory, logistics } = useApiClients()

  const [locationType, setLocationType] = useState<LocationType | ''>('')
  const [warehouses, setWarehouses] = useState<Warehouse[]>([])
  const [stores, setStores] = useState<Store[]>([])
  const [selectedLocationId, setSelectedLocationId] = useState('')
  const [selectedTypes, setSelectedTypes] = useState<InventoryType[]>([...INVENTORY_TYPES])
  const [projectedFilter, setProjectedFilter] = useState<'All' | 'Yes' | 'No'>('All')
  const [results, setResults] = useState<{ type: InventoryType; items: InventoryItem[] }[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  useEffect(() => {
    if (locationType === 'Warehouse') {
      logistics.get<Warehouse[]>('/api/Warehouse').then(setWarehouses).catch(() => setWarehouses([]))
    } else if (locationType === 'Store') {
      logistics.get<Store[]>('/api/Store').then(setStores).catch(() => setStores([]))
    }
    setSelectedLocationId('')
    setResults([])
  }, [locationType])

  function toggleType(type: InventoryType) {
    setSelectedTypes(prev =>
      prev.includes(type) ? prev.filter(t => t !== type) : [...prev, type]
    )
  }

  async function handleSearch() {
    if (!selectedLocationId || selectedTypes.length === 0) return
    setLoading(true)
    setError('')
    try {
      const fetches = selectedTypes.map(async type => {
        const endpoint = `/api/${type}/location/${selectedLocationId}`
        const items = await inventory.get<InventoryItem[]>(endpoint)
        return { type, items }
      })
      const raw = await Promise.all(fetches)
      setResults(raw.map(({ type, items }) => ({
        type,
        items: projectedFilter === 'All' ? items
          : items.filter(i => i.projected === (projectedFilter === 'Yes')),
      })))
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

  return (
    <div style={{ marginTop: '1.5rem' }}>
      <h2>Inventory Viewer</h2>

      <div style={{ display: 'flex', gap: '1rem', flexWrap: 'wrap', alignItems: 'flex-end' }}>
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

        <div>
          Inventory Types
          <br />
          <MultiSelectDropdown
            options={INVENTORY_TYPES}
            selected={selectedTypes}
            onToggle={toggleType}
          />
        </div>

        <button onClick={handleSearch} disabled={!selectedLocationId || selectedTypes.length === 0 || loading}>
          {loading ? 'Loading...' : 'Search'}
        </button>
      </div>

      {error && <p style={{ color: 'red' }}>{error}</p>}

      {results.map(({ type, items }) => (
        <div key={type} style={{ marginTop: '1.5rem' }}>
          <h3>{type} ({items.length})</h3>
          {items.length === 0 ? (
            <p>No items found.</p>
          ) : (
            <table style={{ borderCollapse: 'collapse', width: '100%' }}>
              <thead>
                <tr>
                  {['SKU', 'Unloaded Date', 'Projected', 'Location'].map(h => (
                    <th key={h} style={{ border: '1px solid var(--table-border)', padding: '0.4rem 0.8rem', textAlign: 'left', background: 'var(--surface)' }}>{h}</th>
                  ))}
                </tr>
              </thead>
              <tbody>
                {items.map(item => (
                  <tr key={item.partitionKey}>
                    <td style={{ border: '1px solid var(--table-border)', padding: '0.4rem 0.8rem', background: 'var(--surface)' }}>{item.skuMarker}</td>
                    <td style={{ border: '1px solid var(--table-border)', padding: '0.4rem 0.8rem', background: 'var(--surface)' }}>{new Date(item.unloadedDate).toLocaleDateString()}</td>
                    <td style={{ border: '1px solid var(--table-border)', padding: '0.4rem 0.8rem', background: 'var(--surface)' }}>{item.projected ? 'Yes' : 'No'}</td>
                    <td style={{ border: '1px solid var(--table-border)', padding: '0.4rem 0.8rem', background: 'var(--surface)' }}>{item.locationId}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>
      ))}
    </div>
  )
}
