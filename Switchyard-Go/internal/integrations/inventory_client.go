package integrations

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

// InventoryItem is the Go representation of a single inventory unit returned
// by the Switchyard .NET Inventory API. Category distinguishes the source
// endpoint (clothing, ppe, tool) — the SKUMarker is the canonical identifier.
type InventoryItem struct {
	SKUMarker    string    `json:"sku_marker"`
	Category     string    `json:"category"`
	LocationID   string    `json:"location_id"`
	UnloadedDate time.Time `json:"unloaded_date"`
	Projected    bool      `json:"projected"`
}

// InventoryClient is the only surface in the Go backend that reads inventory
// from the Switchyard .NET Inventory API. No other package may call that API.
type InventoryClient interface {
	// GetByLocation fetches all inventory at a location across Clothing, PPE,
	// and Tool categories concurrently. Returns a combined list.
	GetByLocation(ctx context.Context, locationID string) ([]InventoryItem, error)
}

type HTTPInventoryClient struct {
	baseURL       string
	tokenProvider TokenProvider
	httpClient    *http.Client
}

func NewInventoryClient(baseURL string, tp TokenProvider) *HTTPInventoryClient {
	return &HTTPInventoryClient{
		baseURL:       baseURL,
		tokenProvider: tp,
		httpClient:    newHTTPClient(),
	}
}

func (c *HTTPInventoryClient) GetByLocation(ctx context.Context, locationID string) ([]InventoryItem, error) {
	categories := []string{"Clothing", "PPE", "Tool"}

	type result struct {
		items []InventoryItem
		err   error
	}

	results := make([]result, len(categories))
	var wg sync.WaitGroup

	for i, cat := range categories {
		wg.Add(1)
		go func(idx int, category string) {
			defer wg.Done()
			items, err := c.fetchCategory(ctx, category, locationID)
			results[idx] = result{items: items, err: err}
		}(i, cat)
	}

	wg.Wait()

	var all []InventoryItem
	for i, r := range results {
		if r.err != nil {
			return nil, fmt.Errorf("fetching %s inventory at %s: %w", categories[i], locationID, r.err)
		}
		all = append(all, r.items...)
	}

	return all, nil
}

// dotnetInventoryItem mirrors the Clothing/PPE/Tool JSON shape from the .NET API.
// TODO: verify exact JSON field names against live API response — ASP.NET Core
// camelCase policy renders SKUMarker as "sKUMarker" which is non-standard.
type dotnetInventoryItem struct {
	LocationId   string    `json:"locationId"`
	SKUMarker    string    `json:"sKUMarker"`
	UnloadedDate time.Time `json:"unloadedDate"`
	Projected    bool      `json:"projected"`
}

func (c *HTTPInventoryClient) fetchCategory(ctx context.Context, category, locationID string) ([]InventoryItem, error) {
	url := fmt.Sprintf("%s/api/%s/location/%s", c.baseURL, category, locationID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.tokenProvider())
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s inventory endpoint returned %d for location %s", category, resp.StatusCode, locationID)
	}

	var raw []dotnetInventoryItem
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("decoding %s inventory response: %w", category, err)
	}

	items := make([]InventoryItem, len(raw))
	for i, r := range raw {
		items[i] = InventoryItem{
			SKUMarker:    r.SKUMarker,
			Category:     strings.ToLower(category),
			LocationID:   r.LocationId,
			UnloadedDate: r.UnloadedDate,
			Projected:    r.Projected,
		}
	}

	return items, nil
}
