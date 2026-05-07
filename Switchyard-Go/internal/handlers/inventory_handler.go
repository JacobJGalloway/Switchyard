package handlers

import (
	"net/http"
	"strings"
	"sync"

	"github.com/JacobJGalloway/switchyard-go/internal/integrations"
)

// RegionalInventoryHandler fans out inventory queries across all warehouses in
// the configured network and returns per-warehouse SKU counts plus region totals.
type RegionalInventoryHandler struct {
	inv          integrations.InventoryClient
	warehouseIDs []string
}

func NewRegionalInventoryHandler(inv integrations.InventoryClient, warehouseIDs []string) *RegionalInventoryHandler {
	return &RegionalInventoryHandler{inv: inv, warehouseIDs: warehouseIDs}
}

type warehouseInventory struct {
	WarehouseID string         `json:"warehouse_id"`
	SKUs        map[string]int `json:"skus"`
}

type regionalInventoryResponse struct {
	Warehouses   []warehouseInventory `json:"warehouses"`
	RegionTotals map[string]int       `json:"region_totals"`
}

// GetRegional handles GET /api/inventory/region
// Optional query param: ?sku=SKU001 to filter to a single SKU.
func (h *RegionalInventoryHandler) GetRegional(w http.ResponseWriter, r *http.Request) {
	skuFilter := strings.ToUpper(r.URL.Query().Get("sku"))

	type whResult struct {
		whID  string
		skus  map[string]int
		err   error
	}

	results := make([]whResult, len(h.warehouseIDs))
	var wg sync.WaitGroup

	for i, whID := range h.warehouseIDs {
		wg.Add(1)
		go func(idx int, id string) {
			defer wg.Done()
			items, err := h.inv.GetByLocation(r.Context(), id)
			if err != nil {
				results[idx] = whResult{whID: id, err: err}
				return
			}
			counts := make(map[string]int)
			for _, item := range items {
				if skuFilter == "" || strings.ToUpper(item.SKUMarker) == skuFilter {
					counts[item.SKUMarker]++
				}
			}
			results[idx] = whResult{whID: id, skus: counts}
		}(i, whID)
	}

	wg.Wait()

	regionTotals := make(map[string]int)
	warehouses := make([]warehouseInventory, 0, len(h.warehouseIDs))

	for _, res := range results {
		if res.err != nil {
			writeError(w, http.StatusBadGateway, "inventory fetch failed for "+res.whID+": "+res.err.Error())
			return
		}
		warehouses = append(warehouses, warehouseInventory{WarehouseID: res.whID, SKUs: res.skus})
		for sku, qty := range res.skus {
			regionTotals[sku] += qty
		}
	}

	writeJSON(w, http.StatusOK, regionalInventoryResponse{
		Warehouses:   warehouses,
		RegionTotals: regionTotals,
	})
}
