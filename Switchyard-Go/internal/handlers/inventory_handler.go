package handlers

import (
	"net/http"
	"strings"
	"sync"

	"github.com/JacobJGalloway/switchyard-go/internal/integrations"
	"github.com/JacobJGalloway/switchyard-go/internal/repository"
)

// RegionalInventoryHandler fans out inventory queries across all warehouses in
// the configured region and returns per-warehouse SKU counts plus region totals.
// Warehouse list is loaded from the database at request time — no WAREHOUSE_IDS config needed.
type RegionalInventoryHandler struct {
	inv      integrations.InventoryClient
	whRepo   repository.WarehouseRepository
}

func NewRegionalInventoryHandler(inv integrations.InventoryClient, whRepo repository.WarehouseRepository) *RegionalInventoryHandler {
	return &RegionalInventoryHandler{inv: inv, whRepo: whRepo}
}

type warehouseInventory struct {
	WarehouseID string         `json:"warehouse_id"`
	Region      *string        `json:"region,omitempty"`
	SKUs        map[string]int `json:"skus"`
}

type regionalInventoryResponse struct {
	Warehouses   []warehouseInventory `json:"warehouses"`
	RegionTotals map[string]int       `json:"region_totals"`
}

// GetRegional handles GET /api/inventory/region
// Optional query params: ?sku=SKU001 to filter to a single SKU; ?region=MIDWEST to scope to a region.
func (h *RegionalInventoryHandler) GetRegional(w http.ResponseWriter, r *http.Request) {
	skuFilter := strings.ToUpper(r.URL.Query().Get("sku"))
	regionFilter := r.URL.Query().Get("region")

	var warehouses []*struct {
		id     string
		region *string
	}

	if regionFilter != "" {
		whs, err := h.whRepo.GetByRegion(r.Context(), regionFilter)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to load warehouses")
			return
		}
		for _, wh := range whs {
			wh := wh
			warehouses = append(warehouses, &struct {
				id     string
				region *string
			}{id: wh.ID, region: wh.Region})
		}
	} else {
		whs, err := h.whRepo.GetAll(r.Context())
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to load warehouses")
			return
		}
		for _, wh := range whs {
			wh := wh
			warehouses = append(warehouses, &struct {
				id     string
				region *string
			}{id: wh.ID, region: wh.Region})
		}
	}

	type whResult struct {
		whID   string
		region *string
		skus   map[string]int
		err    error
	}

	results := make([]whResult, len(warehouses))
	var wg sync.WaitGroup

	for i, wh := range warehouses {
		wg.Add(1)
		go func(idx int, id string, region *string) {
			defer wg.Done()
			items, err := h.inv.GetByLocation(r.Context(), id)
			if err != nil {
				results[idx] = whResult{whID: id, region: region, err: err}
				return
			}
			counts := make(map[string]int)
			for _, item := range items {
				if skuFilter == "" || strings.ToUpper(item.SKUMarker) == skuFilter {
					counts[item.SKUMarker]++
				}
			}
			results[idx] = whResult{whID: id, region: region, skus: counts}
		}(i, wh.id, wh.region)
	}

	wg.Wait()

	regionTotals := make(map[string]int)
	out := make([]warehouseInventory, 0, len(warehouses))

	for _, res := range results {
		if res.err != nil {
			writeError(w, http.StatusBadGateway, "inventory fetch failed for "+res.whID+": "+res.err.Error())
			return
		}
		out = append(out, warehouseInventory{WarehouseID: res.whID, Region: res.region, SKUs: res.skus})
		for sku, qty := range res.skus {
			regionTotals[sku] += qty
		}
	}

	writeJSON(w, http.StatusOK, regionalInventoryResponse{
		Warehouses:   out,
		RegionTotals: regionTotals,
	})
}
