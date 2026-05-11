package services

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/JacobJGalloway/switchyard-go/internal/events"
	"github.com/JacobJGalloway/switchyard-go/internal/models"
	"github.com/JacobJGalloway/switchyard-go/internal/repository"
)

// HOSStatus is the color dot shown on in-delivery driver cards.
type HOSStatus string

const (
	HOSStatusGreen  HOSStatus = "green"
	HOSStatusYellow HOSStatus = "yellow"
	HOSStatusRed    HOSStatus = "red"
)

// AlertType classifies the operational condition behind a board alert.
type AlertType string

// AlertSeverity determines how the alert is surfaced to the dispatcher.
type AlertSeverity string

const (
	AlertTypeHOSWarning        AlertType = "hos_warning"
	AlertTypeHOSWeeklyLimit    AlertType = "hos_weekly_limit"
	AlertTypeRoadsideBreakdown AlertType = "roadside_breakdown"
	AlertTypeExpiringDeadhead  AlertType = "expiring_deadhead"

	AlertSeverityWarning  AlertSeverity = "warning"
	AlertSeverityCritical AlertSeverity = "critical"
)

// --- Workflow pipeline card types (left side of board) ---

// DraftCard represents a BOL that has been created but not yet claimed for route planning.
type DraftCard struct {
	PlanBOL *models.PlanBOLRecord
}

// PendingCard represents a BOL actively being route-planned by a dispatcher or route planner.
// Acts as a human-spinner — signals to other planners that this BOL is already in progress.
type PendingCard struct {
	PlanBOL *models.PlanBOLRecord
}

// LoadingReadyCard covers both the Loading and Ready phases (status: loading or validated).
// IsReady distinguishes the two: false = dock loading in progress, true = trailer loaded.
// IsLongWait flags BOLs that have been sitting past the configurable age threshold.
// Assignment and Driver are non-nil when a driver has been assigned but not yet departed.
type LoadingReadyCard struct {
	PlanBOL    *models.PlanBOLRecord
	Stops      []*models.PlanBOLStop
	IsReady    bool
	IsLongWait bool
	Assignment *models.DriverBOLAssignment
	Driver     *models.Driver
}

// InDeliveryCard is the primary driver card while a run is active.
// MandatedStopAt non-nil places the card in the Mandated Stop sub-section.
type InDeliveryCard struct {
	Assignment     *models.DriverBOLAssignment
	Driver         *models.Driver
	PlanBOL        *models.PlanBOLRecord
	Equipment      *models.Equipment
	CurrentStop    *models.PlanBOLStop
	HOSStatus      HOSStatus
	HOSPillTone    string // "ok" | "warn" | "danger"
	HOSPillLabel   string // "Healthy" | "Daily limit · 1h 30m" | "HOS limit"
	HOSWindow      *models.HOSWindow
	MandatedStopAt *time.Time
	ELDStopRef     *string
}

// BreakdownCard represents a roadside equipment failure with load still attached.
// Only roadside+load breakdowns appear in InDelivery — all others go to Maintenance.
type BreakdownCard struct {
	Equipment *models.Equipment
	Breakdown *models.BreakdownRecord
	Driver    *models.Driver
}

// InDeliveryColumn holds the three sub-sections within the In Delivery board column.
// InTransit is expanded by default. MandatedStop and Breakdown are collapsed with count badges.
type InDeliveryColumn struct {
	InTransit    []*InDeliveryCard
	MandatedStop []*InDeliveryCard
	Breakdown    []*BreakdownCard
}

// DeliveredCard appears after all stops are confirmed.
// Displays the dead-head search window countdown.
type DeliveredCard struct {
	Assignment              *models.DriverBOLAssignment
	Driver                  *models.Driver
	PlanBOL                 *models.PlanBOLRecord
	Equipment               *models.Equipment
	DeadheadWindowExpiresAt time.Time
	DeadheadWindowRemaining time.Duration
}

// --- Resource pool card types (right side of board) ---

// AvailableDriverCard represents a driver available for assignment.
// HOSStatus is green (ample hours) or yellow (approaching limit but still assignable).
type AvailableDriverCard struct {
	Driver       *models.Driver
	HOSWindow    *models.HOSWindow
	HOSStatus    HOSStatus
	HOSPillTone  string // "ok" | "warn"
	HOSPillLabel string // "Healthy" | "Daily limit · 1h 30m"
}

// RestingDriverCard represents a driver on mandated rest.
// RestEndsAt is the max of daily and weekly reset times — the driver exits rest fully legal.
// RestType is "daily" or "weekly" to label the countdown on the card.
type RestingDriverCard struct {
	Driver         *models.Driver
	HOSWindow      *models.HOSWindow
	RestEndsAt     time.Time
	RestType       string
	TimeUntilReset time.Duration
}

// AvailableColumn is the driver pool. Default view shows AvailableNow; toggle shows Resting.
type AvailableColumn struct {
	AvailableNow []*AvailableDriverCard
	Resting      []*RestingDriverCard
}

// MaintenanceCard covers both scheduled maintenance and non-load-attached breakdowns.
// Exactly one of Maintenance or Breakdown is non-nil depending on why the equipment is out.
type MaintenanceCard struct {
	Equipment   *models.Equipment
	Maintenance *models.MaintenanceRecord
	Breakdown   *models.BreakdownRecord
}

// --- Board state ---

// BoardState is the complete snapshot returned by GetBoardState.
// Left zone (workflow pipeline): Draft → Pending → LoadingReady → InDelivery → Delivered.
// Right zone (resource pools): Available (driver pool), Maintenance (equipment out of rotation).
type BoardState struct {
	Draft        []*DraftCard
	Pending      []*PendingCard
	LoadingReady []*LoadingReadyCard
	InDelivery   InDeliveryColumn
	Delivered    []*DeliveredCard
	Available    AvailableColumn
	Maintenance  []*MaintenanceCard
	GeneratedAt  time.Time
}

// BoardAlert is an operational condition that requires dispatcher attention.
type BoardAlert struct {
	ID           uuid.UUID
	AlertType    AlertType
	Severity     AlertSeverity
	Message      string
	DriverID     *uuid.UUID
	EquipmentID  *uuid.UUID
	AssignmentID *uuid.UUID
	CreatedAt    time.Time
}

// WhiteboardService assembles the dispatch board and surfaces active alerts.
// Board state is assembled dynamically from live repo data on every request.
type WhiteboardService struct {
	driverRepo            repository.DriverRepository
	assignRepo            repository.AssignmentRepository
	bolRepo               repository.PlanBOLRepository
	equipRepo             repository.EquipmentRepository
	hosRepo               repository.HOSRepository
	warningThresholdHours float64
	deadheadSearchWindow  time.Duration
	loadingAgeThreshold   time.Duration
	defaultCycleLabel     string
}

func NewWhiteboardService(
	driverRepo repository.DriverRepository,
	assignRepo repository.AssignmentRepository,
	bolRepo repository.PlanBOLRepository,
	equipRepo repository.EquipmentRepository,
	hosRepo repository.HOSRepository,
	warningThresholdHours float64,
	deadheadSearchWindowHours float64,
	loadingAgeThresholdHours float64,
	defaultCycleLabel string,
) *WhiteboardService {
	return &WhiteboardService{
		driverRepo:            driverRepo,
		assignRepo:            assignRepo,
		bolRepo:               bolRepo,
		equipRepo:             equipRepo,
		hosRepo:               hosRepo,
		warningThresholdHours: warningThresholdHours,
		deadheadSearchWindow:  time.Duration(deadheadSearchWindowHours * float64(time.Hour)),
		loadingAgeThreshold:   time.Duration(loadingAgeThresholdHours * float64(time.Hour)),
		defaultCycleLabel:     defaultCycleLabel,
	}
}

// GetBoardState assembles the full board from live data.
//
// Pipeline placement rules:
//   - status draft          → Draft
//   - status plan-progress  → Pending
//   - status loading/validated, no active assignment or assigned+not departed → LoadingReady
//   - status submitted, departed_at set, fulfilled_at nil → InDelivery
//   - status submitted, fulfilled_at set, deadhead_confirmed_at nil → Delivered
//   - deadhead_confirmed_at set → archived, not shown
//
// Resource pool placement:
//   - Off-run driver, no mandated stop → Available.AvailableNow
//   - Off-run driver, mandated_stop_at set → Available.Resting
//   - Equipment status maintenance → Maintenance
//   - Equipment status breakdown, roadside+load attached → InDelivery.Breakdown
//   - Equipment status breakdown, all other → Maintenance
func (s *WhiteboardService) GetBoardState(ctx context.Context) (*BoardState, error) {
	board := &BoardState{GeneratedAt: time.Now()}

	// --- Active assignments: drive InDelivery, Delivered, and LoadingReady assignment overlay ---

	activeAssignments, err := s.assignRepo.GetAllActive(ctx)
	if err != nil {
		return nil, fmt.Errorf("loading active assignments: %w", err)
	}

	// assignedBOLs maps planBOLID to its assignment for any active state.
	// Used to annotate LoadingReady cards and skip BOLs from pipeline columns.
	assignedBOLs := make(map[uuid.UUID]*models.DriverBOLAssignment, len(activeAssignments))
	activeDriverIDs := make(map[uuid.UUID]bool, len(activeAssignments))

	for _, a := range activeAssignments {
		assignedBOLs[a.PlanBOLID] = a
		activeDriverIDs[a.DriverID] = true

		if a.DepartedAt == nil {
			// Assigned but not yet departed — BOL stays in LoadingReady with assignment overlay.
			continue
		}

		driver, err := s.driverRepo.GetByID(ctx, a.DriverID)
		if err != nil {
			return nil, fmt.Errorf("loading driver %s: %w", a.DriverID, err)
		}
		bol, err := s.bolRepo.GetByID(ctx, a.PlanBOLID)
		if err != nil {
			return nil, fmt.Errorf("loading plan BOL %s: %w", a.PlanBOLID, err)
		}
		equip, err := s.equipRepo.GetByID(ctx, a.EquipmentID)
		if err != nil {
			return nil, fmt.Errorf("loading equipment %s: %w", a.EquipmentID, err)
		}

		switch {
		case a.FulfilledAt == nil:
			hosWindow, _ := s.hosRepo.GetWindowByDriver(ctx, a.DriverID)
			hosStatus, pillTone, pillLabel := s.hosStateForWindow(ctx, hosWindow, driver.LicenseState)
			card := &InDeliveryCard{
				Assignment:   a,
				Driver:       driver,
				PlanBOL:      bol,
				Equipment:    equip,
				CurrentStop:  s.firstUnprocessedStop(ctx, a.PlanBOLID),
				HOSStatus:    hosStatus,
				HOSPillTone:  pillTone,
				HOSPillLabel: pillLabel,
				HOSWindow:    hosWindow,
			}
			if hosWindow != nil && hosWindow.MandatedStopAt != nil {
				card.MandatedStopAt = hosWindow.MandatedStopAt
				card.ELDStopRef = hosWindow.ELDStopRef
				board.InDelivery.MandatedStop = append(board.InDelivery.MandatedStop, card)
			} else {
				board.InDelivery.InTransit = append(board.InDelivery.InTransit, card)
			}

		default:
			expiresAt := a.FulfilledAt.Add(s.deadheadSearchWindow)
			board.Delivered = append(board.Delivered, &DeliveredCard{
				Assignment:              a,
				Driver:                  driver,
				PlanBOL:                 bol,
				Equipment:               equip,
				DeadheadWindowExpiresAt: expiresAt,
				DeadheadWindowRemaining: time.Until(expiresAt),
			})
		}
	}

	// --- Workflow pipeline: Draft, Pending, LoadingReady ---

	draftBOLs, err := s.bolRepo.GetByStatus(ctx, models.PlanBOLStatusDraft)
	if err != nil {
		return nil, fmt.Errorf("loading draft BOLs: %w", err)
	}
	for _, bol := range draftBOLs {
		board.Draft = append(board.Draft, &DraftCard{PlanBOL: bol})
	}

	pendingBOLs, err := s.bolRepo.GetByStatus(ctx, models.PlanBOLStatusPlanProgress)
	if err != nil {
		return nil, fmt.Errorf("loading pending BOLs: %w", err)
	}
	for _, bol := range pendingBOLs {
		board.Pending = append(board.Pending, &PendingCard{PlanBOL: bol})
	}

	loadingBOLs, err := s.bolRepo.GetByStatus(ctx, models.PlanBOLStatusLoading)
	if err != nil {
		return nil, fmt.Errorf("loading loading-phase BOLs: %w", err)
	}
	readyBOLs, err := s.bolRepo.GetByStatus(ctx, models.PlanBOLStatusValidated)
	if err != nil {
		return nil, fmt.Errorf("loading ready BOLs: %w", err)
	}
	for _, bol := range append(loadingBOLs, readyBOLs...) {
		stops, _ := s.bolRepo.GetStops(ctx, bol.ID)
		card := &LoadingReadyCard{
			PlanBOL:    bol,
			Stops:      stops,
			IsReady:    bol.Status == models.PlanBOLStatusValidated,
			IsLongWait: time.Since(bol.CreatedAt) > s.loadingAgeThreshold,
		}
		if a, ok := assignedBOLs[bol.ID]; ok && a.DepartedAt == nil {
			driver, _ := s.driverRepo.GetByID(ctx, a.DriverID)
			card.Assignment = a
			card.Driver = driver
		}
		board.LoadingReady = append(board.LoadingReady, card)
	}

	// --- Equipment: Maintenance column and InDelivery.Breakdown ---

	allEquipment, err := s.equipRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("loading equipment: %w", err)
	}
	for _, e := range allEquipment {
		switch e.Status {
		case models.EquipmentStatusMaintenance:
			rec, _ := s.equipRepo.GetActiveMaintenanceByEquipment(ctx, e.ID)
			board.Maintenance = append(board.Maintenance, &MaintenanceCard{
				Equipment:   e,
				Maintenance: rec,
			})

		case models.EquipmentStatusBreakdown:
			rec, _ := s.equipRepo.GetActiveBreakdownByEquipment(ctx, e.ID)
			if rec != nil && rec.BreakdownType == models.BreakdownTypeRoadside && rec.LoadAttached {
				bdCard := &BreakdownCard{Equipment: e, Breakdown: rec}
				if rec.DriverID != nil {
					driver, err := s.driverRepo.GetByID(ctx, *rec.DriverID)
					if err != nil {
						return nil, fmt.Errorf("loading driver for breakdown on %s: %w", e.UnitID, err)
					}
					bdCard.Driver = driver
				}
				board.InDelivery.Breakdown = append(board.InDelivery.Breakdown, bdCard)
			} else {
				board.Maintenance = append(board.Maintenance, &MaintenanceCard{
					Equipment: e,
					Breakdown: rec,
				})
			}
		}
	}

	// --- Available: all off-run active drivers, split by rest state ---

	allDrivers, err := s.driverRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("loading drivers: %w", err)
	}
	for _, d := range allDrivers {
		if !d.IsActive || activeDriverIDs[d.ID] {
			continue
		}
		window, _ := s.hosRepo.GetWindowByDriver(ctx, d.ID)
		limit, _ := s.hosRepo.GetLimitByStateAndCycle(ctx, d.LicenseState, s.defaultCycleLabel)

		if window != nil && window.MandatedStopAt != nil {
			restEndsAt, restType := computeRestEnd(window, limit)
			board.Available.Resting = append(board.Available.Resting, &RestingDriverCard{
				Driver:         d,
				HOSWindow:      window,
				RestEndsAt:     restEndsAt,
				RestType:       restType,
				TimeUntilReset: time.Until(restEndsAt),
			})
		} else {
			hosStatus, pillTone, pillLabel := s.hosStateForWindow(ctx, window, d.LicenseState)
			board.Available.AvailableNow = append(board.Available.AvailableNow, &AvailableDriverCard{
				Driver:       d,
				HOSWindow:    window,
				HOSStatus:    hosStatus,
				HOSPillTone:  pillTone,
				HOSPillLabel: pillLabel,
			})
		}
	}

	return board, nil
}

// GetAlerts returns active operational conditions requiring dispatcher attention.
func (s *WhiteboardService) GetAlerts(ctx context.Context) ([]*BoardAlert, error) {
	board, err := s.GetBoardState(ctx)
	if err != nil {
		return nil, err
	}

	var alerts []*BoardAlert
	now := time.Now()

	for _, card := range board.InDelivery.InTransit {
		if card.HOSStatus == HOSStatusYellow || card.HOSStatus == HOSStatusRed {
			driverID := card.Driver.ID
			assignID := card.Assignment.ID
			alerts = append(alerts, &BoardAlert{
				ID:           uuid.New(),
				AlertType:    AlertTypeHOSWarning,
				Severity:     AlertSeverityWarning,
				Message:      fmt.Sprintf("Driver %s is approaching HOS limit", card.Driver.Name),
				DriverID:     &driverID,
				AssignmentID: &assignID,
				CreatedAt:    now,
			})
		}
	}

	for _, card := range board.Available.Resting {
		driverID := card.Driver.ID
		alerts = append(alerts, &BoardAlert{
			ID:        uuid.New(),
			AlertType: AlertTypeHOSWeeklyLimit,
			Severity:  AlertSeverityWarning,
			Message: fmt.Sprintf("Driver %s is on %s rest — available in %.0f minutes",
				card.Driver.Name, card.RestType, card.TimeUntilReset.Minutes()),
			DriverID:  &driverID,
			CreatedAt: now,
		})
	}

	for _, card := range board.InDelivery.Breakdown {
		if card.Breakdown == nil {
			continue
		}
		equipID := card.Equipment.ID
		var driverID *uuid.UUID
		driverRef := ""
		if card.Driver != nil {
			id := card.Driver.ID
			driverID = &id
			driverRef = fmt.Sprintf(" — driver: %s", card.Driver.Name)
		}
		alerts = append(alerts, &BoardAlert{
			ID:          uuid.New(),
			AlertType:   AlertTypeRoadsideBreakdown,
			Severity:    AlertSeverityCritical,
			Message:     fmt.Sprintf("Roadside breakdown with load attached: %s%s", card.Equipment.UnitID, driverRef),
			EquipmentID: &equipID,
			DriverID:    driverID,
			CreatedAt:   now,
		})
	}

	const expiringThreshold = time.Hour
	for _, card := range board.Delivered {
		if card.DeadheadWindowRemaining <= expiringThreshold {
			driverID := card.Driver.ID
			assignID := card.Assignment.ID
			alerts = append(alerts, &BoardAlert{
				ID:           uuid.New(),
				AlertType:    AlertTypeExpiringDeadhead,
				Severity:     AlertSeverityWarning,
				Message:      fmt.Sprintf("Dead-head window expiring for driver %s (%.0f minutes remaining)", card.Driver.Name, card.DeadheadWindowRemaining.Minutes()),
				DriverID:     &driverID,
				AssignmentID: &assignID,
				CreatedAt:    now,
			})
		}
	}

	return alerts, nil
}

// hosStatusForWindow is a thin wrapper kept for backward compatibility.
func (s *WhiteboardService) hosStatusForWindow(ctx context.Context, window *models.HOSWindow, stateCode string) HOSStatus {
	status, _, _ := s.hosStateForWindow(ctx, window, stateCode)
	return status
}

// hosStateForWindow computes the HOS status, pill tone, and pill label for a driver card.
// Returns Green/"ok"/"Healthy" when limit data cannot be fetched — the board renders either way.
func (s *WhiteboardService) hosStateForWindow(ctx context.Context, window *models.HOSWindow, stateCode string) (HOSStatus, string, string) {
	if window == nil {
		return HOSStatusGreen, "ok", "Healthy"
	}
	limit, err := s.hosRepo.GetLimitByStateAndCycle(ctx, stateCode, s.defaultCycleLabel)
	if err != nil || limit == nil {
		return HOSStatusGreen, "ok", "Healthy"
	}
	remaining := limit.DailyDrivingLimitHours - window.DailyHoursUsed
	switch {
	case remaining <= 0:
		return HOSStatusRed, "danger", "HOS limit"
	case remaining <= s.warningThresholdHours:
		h := int(remaining)
		m := int((remaining - float64(h)) * 60)
		if m < 0 {
			m = 0
		}
		return HOSStatusYellow, "warn", fmt.Sprintf("Daily limit · %dh %02dm", h, m)
	default:
		return HOSStatusGreen, "ok", "Healthy"
	}
}

// computeRestEnd calculates when a resting driver will be legally available again.
// Takes the higher of daily reset and weekly reset — ensures the driver exits rest
// fully legal, not just daily-legal. RestType "weekly" indicates both clocks reset.
func computeRestEnd(window *models.HOSWindow, limit *models.HOSLimit) (time.Time, string) {
	if window.MandatedStopAt == nil {
		return time.Time{}, ""
	}
	if limit == nil {
		return window.MandatedStopAt.Add(10 * time.Hour), "daily"
	}
	dailyEnd := window.MandatedStopAt.Add(time.Duration(limit.RestPeriodHours * float64(time.Hour)))
	if window.WeeklyHoursUsed >= limit.WeeklyLimitHours {
		weeklyEnd := window.MandatedStopAt.Add(time.Duration(limit.WeeklyResetHours * float64(time.Hour)))
		if weeklyEnd.After(dailyEnd) {
			return weeklyEnd, "weekly"
		}
	}
	return dailyEnd, "daily"
}

// firstUnprocessedStop returns the next stop to be completed, by sequence.
func (s *WhiteboardService) firstUnprocessedStop(ctx context.Context, planBOLID uuid.UUID) *models.PlanBOLStop {
	stops, err := s.bolRepo.GetStops(ctx, planBOLID)
	if err != nil {
		return nil
	}
	sort.Slice(stops, func(i, j int) bool { return stops[i].Sequence < stops[j].Sequence })
	for _, stop := range stops {
		if !stop.IsProcessed {
			return stop
		}
	}
	return nil
}

// --- events.WhiteboardService callbacks ---
// Board state is assembled dynamically on every request, so these are no-ops in v1.1.
// They satisfy the interface and provide hooks for future WebSocket push notifications.

func (s *WhiteboardService) OnAssignmentDeparted(_ context.Context, _ events.AssignmentPayload) error {
	return nil
}
func (s *WhiteboardService) OnAssignmentFulfilled(_ context.Context, _ events.AssignmentPayload) error {
	return nil
}
func (s *WhiteboardService) OnDeadheadConfirmed(_ context.Context, _ events.AssignmentPayload) error {
	return nil
}
func (s *WhiteboardService) OnMandatedStop(_ context.Context, _ events.MandatedStopPayload) error {
	return nil
}
func (s *WhiteboardService) OnEquipmentBreakdown(_ context.Context, _ events.EquipmentBreakdownPayload) error {
	return nil
}
func (s *WhiteboardService) OnEquipmentResolved(_ context.Context, _ events.EquipmentResolvedPayload) error {
	return nil
}
