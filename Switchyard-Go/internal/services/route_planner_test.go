package services

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/JacobJGalloway/switchyard-go/internal/models"
)

// --- mock for ValidatePlan tests ---

type mockPlanBOLRepo struct {
	stops []*models.PlanBOLStop
	err   error
}

func (m *mockPlanBOLRepo) GetStops(ctx context.Context, id uuid.UUID) ([]*models.PlanBOLStop, error) {
	return m.stops, m.err
}

func (m *mockPlanBOLRepo) Create(_ context.Context, _ *models.PlanBOLRecord) error {
	panic("not implemented")
}
func (m *mockPlanBOLRepo) GetByID(_ context.Context, _ uuid.UUID) (*models.PlanBOLRecord, error) {
	panic("not implemented")
}
func (m *mockPlanBOLRepo) GetByStatus(_ context.Context, _ models.PlanBOLStatus) ([]*models.PlanBOLRecord, error) {
	panic("not implemented")
}
func (m *mockPlanBOLRepo) UpdateStatus(_ context.Context, _ uuid.UUID, _ models.PlanBOLStatus) error {
	panic("not implemented")
}
func (m *mockPlanBOLRepo) SetSubmittedTransactionID(_ context.Context, _ uuid.UUID, _ string) error {
	panic("not implemented")
}
func (m *mockPlanBOLRepo) CreateStop(_ context.Context, _ *models.PlanBOLStop) error {
	panic("not implemented")
}
func (m *mockPlanBOLRepo) GetStopByID(_ context.Context, _ uuid.UUID) (*models.PlanBOLStop, error) {
	panic("not implemented")
}
func (m *mockPlanBOLRepo) MarkStopProcessed(_ context.Context, _ uuid.UUID, _ time.Time) error {
	panic("not implemented")
}
func (m *mockPlanBOLRepo) CreateSnapshot(_ context.Context, _ *models.TruckInventorySnapshot) error {
	panic("not implemented")
}
func (m *mockPlanBOLRepo) GetSnapshots(_ context.Context, _ uuid.UUID) ([]*models.TruckInventorySnapshot, error) {
	panic("not implemented")
}
func (m *mockPlanBOLRepo) GetStatusHistory(_ context.Context, _ uuid.UUID) ([]*models.BOLStatusHistory, error) {
	panic("not implemented")
}

// --- solvePlan tests (pure function — no DB, no I/O) ---

func TestSolvePlan_PreferredStrategy(t *testing.T) {
	tests := []struct {
		name        string
		originWhID  string
		additionals []string
		storeStops  []StopRequest
		inventory   map[string]map[string]int
		wantStops   int
		firstIsWH   bool
	}{
		{
			name:       "single warehouse covers all demand",
			originWhID: "wh-1",
			storeStops: []StopRequest{
				{LocationID: "store-1", Items: map[string]int{"SKU-A": 3, "SKU-B": 2}},
			},
			inventory: map[string]map[string]int{
				"wh-1": {"SKU-A": 10, "SKU-B": 5},
			},
			wantStops: 2, // origin WH + store
			firstIsWH: true,
		},
		{
			name:        "two warehouses split demand",
			originWhID:  "wh-1",
			additionals: []string{"wh-2"},
			storeStops: []StopRequest{
				{LocationID: "store-1", Items: map[string]int{"SKU-A": 5, "SKU-B": 3}},
			},
			inventory: map[string]map[string]int{
				"wh-1": {"SKU-A": 5},
				"wh-2": {"SKU-B": 3},
			},
			wantStops: 3, // origin WH + second WH + store
			firstIsWH: true,
		},
		{
			name:       "multiple stores single warehouse",
			originWhID: "wh-1",
			storeStops: []StopRequest{
				{LocationID: "store-1", Items: map[string]int{"SKU-A": 2}},
				{LocationID: "store-2", Items: map[string]int{"SKU-A": 3}},
			},
			inventory: map[string]map[string]int{
				"wh-1": {"SKU-A": 10},
			},
			wantStops: 3, // origin WH + store-1 + store-2
			firstIsWH: true,
		},
		{
			name:        "origin loads exactly what each store needs",
			originWhID:  "wh-1",
			additionals: []string{"wh-2"},
			storeStops: []StopRequest{
				{LocationID: "store-1", Items: map[string]int{"SKU-A": 4}},
				{LocationID: "store-2", Items: map[string]int{"SKU-B": 6}},
			},
			inventory: map[string]map[string]int{
				"wh-1": {"SKU-A": 4},
				"wh-2": {"SKU-B": 6},
			},
			wantStops: 4, // wh-1 + wh-2 + store-1 + store-2
			firstIsWH: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stops, violations := solvePlan(tt.originWhID, tt.additionals, tt.storeStops, tt.inventory)
			assert.Empty(t, violations)
			require.Len(t, stops, tt.wantStops)
			if tt.firstIsWH {
				assert.Equal(t, models.StopTypeWarehouse, stops[0].stopType)
				assert.Equal(t, tt.originWhID, stops[0].locationID)
			}
		})
	}
}

func TestSolvePlan_WarehouseLoadsMatchDeliveries(t *testing.T) {
	// The origin warehouse should load exactly what the stores collectively need —
	// no more, no less (empty truck rule).
	stops, violations := solvePlan("wh-1", nil, []StopRequest{
		{LocationID: "store-1", Items: map[string]int{"SKU-A": 3}},
		{LocationID: "store-2", Items: map[string]int{"SKU-A": 2}},
	}, map[string]map[string]int{
		"wh-1": {"SKU-A": 20},
	})

	require.Empty(t, violations)
	require.Len(t, stops, 3)
	assert.Equal(t, map[string]int{"SKU-A": 5}, stops[0].items, "origin should load exactly the total demand")
	assert.Equal(t, map[string]int{"SKU-A": 3}, stops[1].items, "store-1 delivery unchanged")
	assert.Equal(t, map[string]int{"SKU-A": 2}, stops[2].items, "store-2 delivery unchanged")
}

func TestSolvePlan_Violations(t *testing.T) {
	tests := []struct {
		name        string
		originWhID  string
		additionals []string
		storeStops  []StopRequest
		inventory   map[string]map[string]int
		wantSubstr  string
	}{
		{
			name:       "no warehouse stocks the required SKU",
			originWhID: "wh-1",
			storeStops: []StopRequest{
				{LocationID: "store-1", Items: map[string]int{"SKU-Z": 5}},
			},
			inventory: map[string]map[string]int{
				"wh-1": {"SKU-A": 10},
			},
			wantSubstr: "SKU-Z",
		},
		{
			name:        "combined warehouse stock insufficient for demand",
			originWhID:  "wh-1",
			additionals: []string{"wh-2"},
			storeStops: []StopRequest{
				{LocationID: "store-1", Items: map[string]int{"SKU-A": 10}},
			},
			inventory: map[string]map[string]int{
				"wh-1": {"SKU-A": 3},
				"wh-2": {"SKU-A": 2},
			},
			wantSubstr: "SKU-A",
		},
		{
			name:       "warehouse has zero stock",
			originWhID: "wh-1",
			storeStops: []StopRequest{
				{LocationID: "store-1", Items: map[string]int{"SKU-A": 1}},
			},
			inventory: map[string]map[string]int{
				"wh-1": {},
			},
			wantSubstr: "SKU-A",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, violations := solvePlan(tt.originWhID, tt.additionals, tt.storeStops, tt.inventory)
			require.NotEmpty(t, violations)
			assert.Contains(t, violations[0], tt.wantSubstr)
		})
	}
}

// --- canDeliver pure-function tests ---

func TestCanDeliver_EmptyItems_ReturnsTrue(t *testing.T) {
	// items is empty — loop never executes, returns true directly.
	assert.True(t, canDeliver(map[string]int{"SKU-A": 5}, map[string]int{}))
}

// --- tryFallbackStrategy direct tests ---

func TestSolvePlan_FallbackSucceeds_WhenPreferredFails(t *testing.T) {
	// Origin has only 3 of SKU-A, but store needs 5. Preferred fails.
	// wh-2 can cover the 2-unit shortfall, so fallback succeeds.
	stops, violations := solvePlan("wh-1", []string{"wh-2"},
		[]StopRequest{
			{LocationID: "store-1", Items: map[string]int{"SKU-A": 5}},
		},
		map[string]map[string]int{
			"wh-1": {"SKU-A": 3},
			"wh-2": {"SKU-A": 2},
		},
	)
	assert.Empty(t, violations)
	require.Len(t, stops, 3) // wh-1 load + wh-2 top-up + store-1 delivery
}

func TestSolvePlan_FallbackSkipsIrrelevantWarehouse(t *testing.T) {
	// wh-bad has SKU-B only — cannot cover the SKU-A shortfall.
	// wh-good has SKU-A — is used instead. Exercises !contributesAny continue branch.
	stops, violations := solvePlan("wh-1", []string{"wh-bad", "wh-good"},
		[]StopRequest{
			{LocationID: "store-1", Items: map[string]int{"SKU-A": 5}},
		},
		map[string]map[string]int{
			"wh-1":   {"SKU-A": 3},
			"wh-bad": {"SKU-B": 10},
			"wh-good": {"SKU-A": 2},
		},
	)
	assert.Empty(t, violations)
	require.Len(t, stops, 3) // wh-1 + wh-good + store-1 (wh-bad skipped)
	locationIDs := []string{stops[0].locationID, stops[1].locationID, stops[2].locationID}
	assert.NotContains(t, locationIDs, "wh-bad")
	assert.Contains(t, locationIDs, "wh-good")
}

// --- ValidatePlan tests ---

func TestValidatePlan_ValidPlan(t *testing.T) {
	planID := uuid.New()
	repo := &mockPlanBOLRepo{
		stops: []*models.PlanBOLStop{
			{Sequence: 1, LocationID: "wh-1", StopType: models.StopTypeWarehouse, DeliveryItems: map[string]int{"SKU-A": 5}},
			{Sequence: 2, LocationID: "store-1", StopType: models.StopTypeStore, DeliveryItems: map[string]int{"SKU-A": 3}},
			{Sequence: 3, LocationID: "store-2", StopType: models.StopTypeStore, DeliveryItems: map[string]int{"SKU-A": 2}},
		},
	}
	svc := NewRoutePlannerService(repo, nil, nil)
	violations, err := svc.ValidatePlan(context.Background(), planID)
	require.NoError(t, err)
	assert.Empty(t, violations)
}

func TestValidatePlan_ShortfallAtStore(t *testing.T) {
	planID := uuid.New()
	repo := &mockPlanBOLRepo{
		stops: []*models.PlanBOLStop{
			{Sequence: 1, LocationID: "wh-1", StopType: models.StopTypeWarehouse, DeliveryItems: map[string]int{"SKU-A": 2}},
			{Sequence: 2, LocationID: "store-1", StopType: models.StopTypeStore, DeliveryItems: map[string]int{"SKU-A": 5}},
		},
	}
	svc := NewRoutePlannerService(repo, nil, nil)
	violations, err := svc.ValidatePlan(context.Background(), planID)
	require.NoError(t, err)
	require.Len(t, violations, 1)
	assert.Contains(t, violations[0], "store-1")
	assert.Contains(t, violations[0], "SKU-A")
}

func TestValidatePlan_TruckNotEmptyAtEnd(t *testing.T) {
	// Warehouse loads more than stores collectively receive — violates empty truck rule.
	planID := uuid.New()
	repo := &mockPlanBOLRepo{
		stops: []*models.PlanBOLStop{
			{Sequence: 1, LocationID: "wh-1", StopType: models.StopTypeWarehouse, DeliveryItems: map[string]int{"SKU-A": 10}},
			{Sequence: 2, LocationID: "store-1", StopType: models.StopTypeStore, DeliveryItems: map[string]int{"SKU-A": 3}},
		},
	}
	svc := NewRoutePlannerService(repo, nil, nil)
	violations, err := svc.ValidatePlan(context.Background(), planID)
	require.NoError(t, err)
	require.Len(t, violations, 1)
	assert.Contains(t, violations[0], "truck not empty")
}

func TestValidatePlan_MultipleViolations(t *testing.T) {
	// Two store stops each with a shortfall should produce two violation entries.
	planID := uuid.New()
	repo := &mockPlanBOLRepo{
		stops: []*models.PlanBOLStop{
			{Sequence: 1, LocationID: "wh-1", StopType: models.StopTypeWarehouse, DeliveryItems: map[string]int{"SKU-A": 1}},
			{Sequence: 2, LocationID: "store-1", StopType: models.StopTypeStore, DeliveryItems: map[string]int{"SKU-A": 3}},
			{Sequence: 3, LocationID: "store-2", StopType: models.StopTypeStore, DeliveryItems: map[string]int{"SKU-A": 3}},
		},
	}
	svc := NewRoutePlannerService(repo, nil, nil)
	violations, err := svc.ValidatePlan(context.Background(), planID)
	require.NoError(t, err)
	assert.Len(t, violations, 2)
}
