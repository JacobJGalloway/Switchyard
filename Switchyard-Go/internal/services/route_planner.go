package services

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/JacobJGalloway/switchyard-go/internal/integrations"
	"github.com/JacobJGalloway/switchyard-go/internal/models"
	"github.com/JacobJGalloway/switchyard-go/internal/repository"
)

// hosChecker is the minimal HOS interface the route planner requires.
type hosChecker interface {
	CanAssign(ctx context.Context, driverID uuid.UUID, estimatedRunHours float64, stateCode, cycleLabel string) error
}

// StopRequest is a single store stop and what must be delivered there.
type StopRequest struct {
	LocationID string
	Items      map[string]int // sku_marker → quantity
}

// PlanRouteInput is everything needed to build and validate a route plan.
type PlanRouteInput struct {
	DriverID               uuid.UUID
	OriginWarehouseID      string
	AdditionalWarehouseIDs []string
	StoreStops             []StopRequest
	EstimatedRunHours      float64
	StateCode              string
	CycleLabel             string
}

// RoutePlannerService builds valid BOL route plans and enforces truck inventory constraints.
type RoutePlannerService struct {
	planRepo  repository.PlanBOLRepository
	invClient integrations.InventoryClient
	hos       hosChecker
}

func NewRoutePlannerService(
	planRepo repository.PlanBOLRepository,
	invClient integrations.InventoryClient,
	hos hosChecker,
) *RoutePlannerService {
	return &RoutePlannerService{planRepo: planRepo, invClient: invClient, hos: hos}
}

// PlanRoute resolves a valid stop sequence, persists the PlanBOLRecord and its stops,
// and validates HOS eligibility. Returns an error if any hard constraint fails.
func (s *RoutePlannerService) PlanRoute(ctx context.Context, in PlanRouteInput) (*models.PlanBOLRecord, error) {
	if err := s.hos.CanAssign(ctx, in.DriverID, in.EstimatedRunHours, in.StateCode, in.CycleLabel); err != nil {
		return nil, fmt.Errorf("HOS check: %w", err)
	}

	allWhIDs := append([]string{in.OriginWarehouseID}, in.AdditionalWarehouseIDs...)
	whInventory, err := s.fetchWarehouseInventory(ctx, allWhIDs)
	if err != nil {
		return nil, err
	}

	resolved, violations := solvePlan(in.OriginWarehouseID, in.AdditionalWarehouseIDs, in.StoreStops, whInventory)
	if len(violations) > 0 {
		return nil, fmt.Errorf("route plan unsatisfiable: %s", strings.Join(violations, "; "))
	}

	plan := &models.PlanBOLRecord{
		ID:              uuid.New(),
		DriverID:        in.DriverID,
		OriginatingWhID: in.OriginWarehouseID,
		Status:          models.PlanBOLStatusDraft,
		CreatedAt:       time.Now().UTC(),
	}
	if err := s.planRepo.Create(ctx, plan); err != nil {
		return nil, fmt.Errorf("creating plan BOL record: %w", err)
	}

	now := time.Now().UTC()
	for i, rs := range resolved {
		stop := &models.PlanBOLStop{
			ID:            uuid.New(),
			PlanBOLID:     plan.ID,
			Sequence:      i + 1,
			LocationID:    rs.locationID,
			StopType:      rs.stopType,
			DeliveryItems: rs.items,
		}
		// Stop 1 is the origin warehouse — auto-processed at BOL creation (§4.1 invariant).
		if i == 0 {
			stop.IsProcessed = true
			stop.ProcessedAt = &now
		}
		if err := s.planRepo.CreateStop(ctx, stop); err != nil {
			return nil, fmt.Errorf("creating stop %d (%s): %w", i+1, rs.locationID, err)
		}
	}

	return plan, nil
}

// ValidatePlan re-runs truck inventory constraints over a persisted plan's stops.
// Returns a list of violations; an empty slice means the plan is valid.
func (s *RoutePlannerService) ValidatePlan(ctx context.Context, planBOLID uuid.UUID) ([]string, error) {
	stops, err := s.planRepo.GetStops(ctx, planBOLID)
	if err != nil {
		return nil, fmt.Errorf("fetching stops: %w", err)
	}

	truck := make(map[string]int)
	var violations []string

	for _, stop := range stops {
		switch stop.StopType {
		case models.StopTypeWarehouse:
			for sku, qty := range stop.DeliveryItems {
				truck[sku] += qty
			}
		case models.StopTypeStore:
			for sku, qty := range stop.DeliveryItems {
				if truck[sku] < qty {
					violations = append(violations, fmt.Sprintf(
						"stop %d (%s): need %d %s, have %d on truck",
						stop.Sequence, stop.LocationID, qty, sku, truck[sku],
					))
				}
				truck[sku] -= qty
				if truck[sku] <= 0 {
					delete(truck, sku)
				}
			}
		}
	}

	if len(violations) == 0 && len(truck) > 0 {
		var leftover []string
		for sku, qty := range truck {
			leftover = append(leftover, fmt.Sprintf("%s×%d", sku, qty))
		}
		sort.Strings(leftover)
		violations = append(violations, "truck not empty at final stop: "+strings.Join(leftover, ", "))
	}

	return violations, nil
}

func (s *RoutePlannerService) fetchWarehouseInventory(ctx context.Context, whIDs []string) (map[string]map[string]int, error) {
	out := make(map[string]map[string]int, len(whIDs))
	for _, id := range whIDs {
		items, err := s.invClient.GetByLocation(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("fetching inventory at warehouse %s: %w", id, err)
		}
		counts := make(map[string]int, len(items))
		for _, item := range items {
			counts[item.SKUMarker]++
		}
		out[id] = counts
	}
	return out, nil
}

// resolvedStop is the solver's internal representation of a planned stop.
// items holds what to load (warehouse stops) or deliver (store stops).
type resolvedStop struct {
	locationID string
	stopType   models.StopType
	items      map[string]int
}

// solvePlan first attempts the preferred strategy (all warehouse loads before any store
// deliveries), then falls back to mid-route warehouse insertion if necessary.
func solvePlan(
	originWhID string,
	additionalWhIDs []string,
	storeStops []StopRequest,
	whInventory map[string]map[string]int,
) ([]resolvedStop, []string) {
	demand := make(map[string]int)
	for _, s := range storeStops {
		for sku, qty := range s.Items {
			demand[sku] += qty
		}
	}

	whOrder := append([]string{originWhID}, additionalWhIDs...)

	if stops, ok := tryPreferredStrategy(originWhID, additionalWhIDs, whOrder, storeStops, demand, whInventory); ok {
		return stops, nil
	}
	return tryFallbackStrategy(originWhID, additionalWhIDs, storeStops, demand, whInventory)
}

// tryPreferredStrategy sequences all warehouse pickups before any store deliveries.
// Returns (stops, true) if the full demand can be satisfied; (nil, false) otherwise.
func tryPreferredStrategy(
	originWhID string,
	additionalWhIDs []string,
	whOrder []string,
	storeStops []StopRequest,
	demand map[string]int,
	whInventory map[string]map[string]int,
) ([]resolvedStop, bool) {
	remaining := copyIntMap(demand)
	whLoads := make(map[string]map[string]int)

	for _, whID := range whOrder {
		inv := whInventory[whID]
		load := make(map[string]int)
		for sku, need := range remaining {
			if avail := inv[sku]; avail > 0 {
				loadQty := need
				if avail < need {
					loadQty = avail
				}
				load[sku] = loadQty
				remaining[sku] -= loadQty
				if remaining[sku] == 0 {
					delete(remaining, sku)
				}
			}
		}
		if len(load) > 0 {
			whLoads[whID] = load
		}
	}

	if len(remaining) > 0 {
		return nil, false
	}

	var stops []resolvedStop

	originLoad := whLoads[originWhID]
	if originLoad == nil {
		originLoad = make(map[string]int)
	}
	stops = append(stops, resolvedStop{locationID: originWhID, stopType: models.StopTypeWarehouse, items: originLoad})

	for _, whID := range additionalWhIDs {
		if load, ok := whLoads[whID]; ok {
			stops = append(stops, resolvedStop{locationID: whID, stopType: models.StopTypeWarehouse, items: load})
		}
	}

	for _, s := range storeStops {
		stops = append(stops, resolvedStop{locationID: s.LocationID, stopType: models.StopTypeStore, items: s.Items})
	}

	return stops, true
}

// tryFallbackStrategy inserts warehouse stops mid-route at the first shortfall point.
func tryFallbackStrategy(
	originWhID string,
	additionalWhIDs []string,
	storeStops []StopRequest,
	demand map[string]int,
	whInventory map[string]map[string]int,
) ([]resolvedStop, []string) {
	truck := make(map[string]int)

	originLoad := make(map[string]int)
	for sku, need := range demand {
		if avail := whInventory[originWhID][sku]; avail > 0 {
			loadQty := need
			if avail < need {
				loadQty = avail
			}
			originLoad[sku] = loadQty
			truck[sku] = loadQty
		}
	}

	stops := []resolvedStop{{locationID: originWhID, stopType: models.StopTypeWarehouse, items: originLoad}}
	whUsed := map[string]bool{originWhID: true}
	var violations []string

	for _, s := range storeStops {
		for !canDeliver(truck, s.Items) {
			sf := shortfall(truck, s.Items)
			inserted := false

			for _, whID := range additionalWhIDs {
				if whUsed[whID] {
					continue
				}
				inv := whInventory[whID]

				contributesAny := false
				for sku := range sf {
					if inv[sku] > 0 {
						contributesAny = true
						break
					}
				}
				if !contributesAny {
					continue
				}

				load := make(map[string]int)
				for sku, need := range sf {
					if avail := inv[sku]; avail > 0 {
						loadQty := need
						if avail < need {
							loadQty = avail
						}
						load[sku] = loadQty
						truck[sku] += loadQty
					}
				}
				stops = append(stops, resolvedStop{locationID: whID, stopType: models.StopTypeWarehouse, items: load})
				whUsed[whID] = true
				inserted = true
				break
			}

			if !inserted {
				missing := shortfallKeys(shortfall(truck, s.Items))
				violations = append(violations, fmt.Sprintf(
					"stop %s: no warehouse can cover missing SKUs [%s]",
					s.LocationID, strings.Join(missing, ", "),
				))
				break
			}
		}

		if canDeliver(truck, s.Items) {
			for sku, qty := range s.Items {
				truck[sku] -= qty
				if truck[sku] <= 0 {
					delete(truck, sku)
				}
			}
			stops = append(stops, resolvedStop{locationID: s.LocationID, stopType: models.StopTypeStore, items: s.Items})
		}
	}

	if len(violations) == 0 && len(truck) > 0 {
		var leftover []string
		for sku, qty := range truck {
			leftover = append(leftover, fmt.Sprintf("%s×%d", sku, qty))
		}
		sort.Strings(leftover)
		violations = append(violations, "truck not empty at final stop: "+strings.Join(leftover, ", "))
	}

	return stops, violations
}

func canDeliver(truck map[string]int, items map[string]int) bool {
	for sku, qty := range items {
		if truck[sku] < qty {
			return false
		}
	}
	return true
}

func shortfall(truck map[string]int, items map[string]int) map[string]int {
	sf := make(map[string]int)
	for sku, qty := range items {
		if have := truck[sku]; have < qty {
			sf[sku] = qty - have
		}
	}
	return sf
}

func shortfallKeys(sf map[string]int) []string {
	keys := make([]string, 0, len(sf))
	for sku := range sf {
		keys = append(keys, sku)
	}
	sort.Strings(keys)
	return keys
}

func copyIntMap(m map[string]int) map[string]int {
	c := make(map[string]int, len(m))
	for k, v := range m {
		c[k] = v
	}
	return c
}
