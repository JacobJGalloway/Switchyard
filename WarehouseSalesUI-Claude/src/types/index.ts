export interface Warehouse {
  warehouseId: string
  city: string
  state: string
}

export interface Store {
  partitionKey: string
  storeId: string
  baseWarehouseId: string
  city: string
  state: string
}

export interface InventoryItem {
  partitionKey: string
  rowKey: string
  locationId: string
  skuMarker: string
  unloadedDate: string
  projected: boolean
}

export type InventoryType = 'Clothing' | 'PPE' | 'Tool'
export type LocationType = 'Warehouse' | 'Store'
