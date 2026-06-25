import { useEffect, useMemo, useState } from 'react'
import {
  LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer
} from 'recharts'
import { useApiClients } from '../hooks/useApiClients'
import type { SKUCatalogItem, SKUMovementRow, Warehouse } from '../types'
import styles from './SKUMovementChart.module.css'

const LINE_COLORS = [
  '#4f8ef7', '#f7894f', '#4fc97a', '#f7cf4f', '#c44ff7', '#4ff7f0'
]

function formatDate(iso: string) {
  const d = new Date(iso)
  return d.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })
}

function formatCurrency(value: number) {
  return `$${value.toLocaleString('en-US', { minimumFractionDigits: 0, maximumFractionDigits: 0 })}`
}

export default function SKUMovementChart() {
  const { inventory, logistics } = useApiClients()

  const [movement, setMovement] = useState<SKUMovementRow[]>([])
  const [catalog, setCatalog] = useState<SKUCatalogItem[]>([])
  const [warehouses, setWarehouses] = useState<Warehouse[]>([])
  const [selectedWarehouses, setSelectedWarehouses] = useState<string[]>([])
  const [selectedSKUs, setSelectedSKUs] = useState<string[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    async function load() {
      const [mov, cat, wh] = await Promise.all([
        logistics.get<SKUMovementRow[]>('/api/analytics/sku-movement?days=30'),
        inventory.get<SKUCatalogItem[]>('/api/skucatalog'),
        logistics.get<Warehouse[]>('/api/warehouses'),
      ])
      setMovement(mov)
      setCatalog(cat)
      setWarehouses(wh)
      setSelectedWarehouses(wh.map(w => w.warehouseId))

      // Default: top 1 SKU per category by total value
      const priceMap = new Map(cat.map(c => [c.skuMarker, c.unitPrice]))
      const totals = new Map<string, number>()
      for (const row of mov) {
        const price = priceMap.get(row.skuMarker) ?? 0
        totals.set(row.skuMarker, (totals.get(row.skuMarker) ?? 0) + row.quantityMoved * price)
      }
      const catMap = new Map(cat.map(c => [c.skuMarker, c.category]))
      const topByCat = new Map<string, { sku: string; value: number }>()
      for (const [sku, value] of totals) {
        const category = catMap.get(sku) ?? 'Unknown'
        const current = topByCat.get(category)
        if (!current || value > current.value) topByCat.set(category, { sku, value })
      }
      setSelectedSKUs([...topByCat.values()].map(v => v.sku))
      setLoading(false)
    }
    load()
  }, [])

  // Re-fetch when warehouse selection changes
  useEffect(() => {
    if (loading) return
    const whParam = selectedWarehouses.length > 0
      ? `&warehouses=${selectedWarehouses.join(',')}`
      : ''
    logistics.get<SKUMovementRow[]>(`/api/analytics/sku-movement?days=30${whParam}`)
      .then(setMovement)
  }, [selectedWarehouses])

  // Build chart data: one entry per date, value per SKU as keys
  const chartData = useMemo(() => {
    const priceMap = new Map(catalog.map(c => [c.skuMarker, c.unitPrice]))
    const byDate = new Map<string, Record<string, number>>()
    for (const row of movement) {
      if (!selectedSKUs.includes(row.skuMarker)) continue
      const dateKey = formatDate(row.date)
      if (!byDate.has(dateKey)) byDate.set(dateKey, { date: dateKey } as Record<string, number>)
      const price = priceMap.get(row.skuMarker) ?? 0
      byDate.get(dateKey)![row.skuMarker] = (byDate.get(dateKey)![row.skuMarker] ?? 0) + row.quantityMoved * price
    }
    return [...byDate.values()]
  }, [movement, catalog, selectedSKUs])

  function toggleWarehouse(id: string) {
    setSelectedWarehouses(prev =>
      prev.includes(id) ? prev.filter(w => w !== id) : [...prev, id]
    )
  }

  function toggleSKU(sku: string) {
    setSelectedSKUs(prev =>
      prev.includes(sku) ? prev.filter(s => s !== sku) : [...prev, sku]
    )
  }

  if (loading) return <p className={styles.loading}>Loading chart data…</p>

  return (
    <div className={styles.widget}>
      <div className={styles.header}>
        <h2 className={styles.title}>SKU Value Moved — Last 30 Days</h2>
      </div>

      <div className={styles.filters}>
        <fieldset className={styles.fieldset}>
          <legend className={styles.legend}>Warehouses</legend>
          <div className={styles.checkGroup}>
            {warehouses.map(wh => (
              <label key={wh.warehouseId} className={styles.check}>
                <input
                  type="checkbox"
                  checked={selectedWarehouses.includes(wh.warehouseId)}
                  onChange={() => toggleWarehouse(wh.warehouseId)}
                />
                {wh.warehouseId}
              </label>
            ))}
          </div>
        </fieldset>

        <fieldset className={styles.fieldset}>
          <legend className={styles.legend}>SKUs</legend>
          <div className={styles.checkGroup}>
            {catalog.map(item => (
              <label key={item.skuMarker} className={styles.check}>
                <input
                  type="checkbox"
                  checked={selectedSKUs.includes(item.skuMarker)}
                  onChange={() => toggleSKU(item.skuMarker)}
                />
                {item.skuMarker}
              </label>
            ))}
          </div>
        </fieldset>
      </div>

      <ResponsiveContainer width="100%" height={320}>
        <LineChart data={chartData} margin={{ top: 8, right: 16, left: 8, bottom: 0 }}>
          <CartesianGrid strokeDasharray="3 3" stroke="var(--border)" />
          <XAxis dataKey="date" tick={{ fontSize: 12 }} />
          <YAxis tickFormatter={formatCurrency} tick={{ fontSize: 12 }} width={72} />
          <Tooltip formatter={(v: number) => formatCurrency(v)} />
          <Legend />
          {selectedSKUs.map((sku, i) => (
            <Line
              key={sku}
              type="monotone"
              dataKey={sku}
              stroke={LINE_COLORS[i % LINE_COLORS.length]}
              dot={{ r: 4 }}
              activeDot={{ r: 6 }}
              strokeWidth={2}
            />
          ))}
        </LineChart>
      </ResponsiveContainer>
    </div>
  )
}
