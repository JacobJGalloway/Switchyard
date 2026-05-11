package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/JacobJGalloway/switchyard-go/internal/integrations"
)

type stubInventoryClient struct {
	items []integrations.InventoryItem
	err   error
}

func (s *stubInventoryClient) GetByLocation(_ context.Context, _ string) ([]integrations.InventoryItem, error) {
	return s.items, s.err
}

// --- GetRegional ---

func TestRegionalInventory_ClientError_Returns502(t *testing.T) {
	h := NewRegionalInventoryHandler(&stubInventoryClient{err: errors.New("upstream down")}, []string{"WH001"})
	req := httptest.NewRequest(http.MethodGet, "/api/inventory/region", nil)
	rec := httptest.NewRecorder()
	h.GetRegional(rec, req)
	assert.Equal(t, http.StatusBadGateway, rec.Code)
}

func TestRegionalInventory_Success_Returns200(t *testing.T) {
	items := []integrations.InventoryItem{
		{SKUMarker: "SKU-A", Category: "clothing", LocationID: "WH001"},
		{SKUMarker: "SKU-A", Category: "clothing", LocationID: "WH001"},
		{SKUMarker: "SKU-B", Category: "ppe", LocationID: "WH001"},
	}
	h := NewRegionalInventoryHandler(&stubInventoryClient{items: items}, []string{"WH001", "WH002"})
	req := httptest.NewRequest(http.MethodGet, "/api/inventory/region", nil)
	rec := httptest.NewRecorder()
	h.GetRegional(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestRegionalInventory_SKUFilter_Returns200(t *testing.T) {
	items := []integrations.InventoryItem{
		{SKUMarker: "SKU-A", Category: "clothing", LocationID: "WH001"},
		{SKUMarker: "SKU-B", Category: "ppe", LocationID: "WH001"},
	}
	h := NewRegionalInventoryHandler(&stubInventoryClient{items: items}, []string{"WH001"})
	req := httptest.NewRequest(http.MethodGet, "/api/inventory/region?sku=SKU-A", nil)
	rec := httptest.NewRecorder()
	h.GetRegional(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}
