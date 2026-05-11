package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/JacobJGalloway/switchyard-go/internal/events"
	"github.com/JacobJGalloway/switchyard-go/internal/models"
)

// --- minimal stubs ---

type stubAssignRepo struct {
	assignment *models.DriverBOLAssignment
}

func (r *stubAssignRepo) GetAllActive(_ context.Context) ([]*models.DriverBOLAssignment, error) {
	return nil, nil
}
func (r *stubAssignRepo) GetByID(_ context.Context, _ uuid.UUID) (*models.DriverBOLAssignment, error) {
	if r.assignment == nil {
		return nil, errNotFound
	}
	return r.assignment, nil
}
func (r *stubAssignRepo) Create(_ context.Context, _ *models.DriverBOLAssignment) error { return nil }
func (r *stubAssignRepo) GetByPlanBOL(_ context.Context, _ uuid.UUID) (*models.DriverBOLAssignment, error) {
	return nil, nil
}
func (r *stubAssignRepo) GetActiveByDriver(_ context.Context, _ uuid.UUID) (*models.DriverBOLAssignment, error) {
	return nil, nil
}
func (r *stubAssignRepo) MarkDeparted(_ context.Context, _ uuid.UUID, _ time.Time) error {
	return nil
}
func (r *stubAssignRepo) MarkFulfilled(_ context.Context, _ uuid.UUID, _ time.Time) error {
	return nil
}
func (r *stubAssignRepo) ConfirmDeadhead(_ context.Context, _ uuid.UUID, _ time.Time) error {
	return nil
}

type stubDriverRepo struct {
	driver    *models.Driver
	all       []*models.Driver
	err       error
	getAllErr  error
}

func (r *stubDriverRepo) GetAll(_ context.Context) ([]*models.Driver, error) {
	return r.all, r.getAllErr
}
func (r *stubDriverRepo) GetByID(_ context.Context, _ uuid.UUID) (*models.Driver, error) {
	if r.err != nil {
		return nil, r.err
	}
	return r.driver, nil
}
func (r *stubDriverRepo) Create(_ context.Context, _ *models.Driver) error { return nil }
func (r *stubDriverRepo) Update(_ context.Context, _ *models.Driver) error { return nil }

type stubBOLRepo struct {
	bol     *models.PlanBOLRecord
	stop    *models.PlanBOLStop
	stopErr error
	stops   []*models.PlanBOLStop
}

func (r *stubBOLRepo) Create(_ context.Context, _ *models.PlanBOLRecord) error { return nil }
func (r *stubBOLRepo) GetByID(_ context.Context, _ uuid.UUID) (*models.PlanBOLRecord, error) {
	if r.bol == nil {
		return nil, errNotFound
	}
	return r.bol, nil
}
func (r *stubBOLRepo) GetByStatus(_ context.Context, _ models.PlanBOLStatus) ([]*models.PlanBOLRecord, error) {
	return nil, nil
}
func (r *stubBOLRepo) UpdateStatus(_ context.Context, _ uuid.UUID, _ models.PlanBOLStatus) error {
	return nil
}
func (r *stubBOLRepo) SetSubmittedTransactionID(_ context.Context, _ uuid.UUID, _ string) error {
	return nil
}
func (r *stubBOLRepo) CreateStop(_ context.Context, _ *models.PlanBOLStop) error { return nil }
func (r *stubBOLRepo) GetStops(_ context.Context, _ uuid.UUID) ([]*models.PlanBOLStop, error) {
	return r.stops, nil
}
func (r *stubBOLRepo) GetStopByID(_ context.Context, _ uuid.UUID) (*models.PlanBOLStop, error) {
	if r.stopErr != nil {
		return nil, r.stopErr
	}
	return r.stop, nil
}
func (r *stubBOLRepo) MarkStopProcessed(_ context.Context, _ uuid.UUID, _ time.Time) error {
	return nil
}
func (r *stubBOLRepo) CreateSnapshot(_ context.Context, _ *models.TruckInventorySnapshot) error {
	return nil
}
func (r *stubBOLRepo) GetSnapshots(_ context.Context, _ uuid.UUID) ([]*models.TruckInventorySnapshot, error) {
	return nil, nil
}
func (r *stubBOLRepo) GetStatusHistory(_ context.Context, _ uuid.UUID) ([]*models.BOLStatusHistory, error) {
	return nil, nil
}

type stubAssignEquipRepo struct{ equipment *models.Equipment }

func (r *stubAssignEquipRepo) GetAll(_ context.Context) ([]*models.Equipment, error) { return nil, nil }
func (r *stubAssignEquipRepo) GetByID(_ context.Context, _ uuid.UUID) (*models.Equipment, error) {
	if r.equipment == nil {
		return nil, errNotFound
	}
	return r.equipment, nil
}
func (r *stubAssignEquipRepo) Create(_ context.Context, _ *models.Equipment) error { return nil }
func (r *stubAssignEquipRepo) UpdateStatus(_ context.Context, _ uuid.UUID, _ models.EquipmentStatus) error {
	return nil
}
func (r *stubAssignEquipRepo) CreateMaintenanceRecord(_ context.Context, _ *models.MaintenanceRecord) error {
	return nil
}
func (r *stubAssignEquipRepo) ResolveMaintenanceRecord(_ context.Context, _ uuid.UUID, _ time.Time) error {
	return nil
}
func (r *stubAssignEquipRepo) GetActiveMaintenanceByEquipment(_ context.Context, _ uuid.UUID) (*models.MaintenanceRecord, error) {
	return nil, nil
}
func (r *stubAssignEquipRepo) CreateBreakdownRecord(_ context.Context, _ *models.BreakdownRecord) error {
	return nil
}
func (r *stubAssignEquipRepo) ResolveBreakdownRecord(_ context.Context, _ uuid.UUID, _ time.Time) error {
	return nil
}
func (r *stubAssignEquipRepo) GetActiveBreakdownByEquipment(_ context.Context, _ uuid.UUID) (*models.BreakdownRecord, error) {
	return nil, nil
}

type stubHOSSvc struct{ err error }

func (s *stubHOSSvc) CanAssign(_ context.Context, _ uuid.UUID, _ float64, _, _ string) error {
	return s.err
}

type stubWBSvc struct{}

func (s *stubWBSvc) OnAssignmentDeparted(_ context.Context, _ events.AssignmentPayload) error {
	return nil
}
func (s *stubWBSvc) OnAssignmentFulfilled(_ context.Context, _ events.AssignmentPayload) error {
	return nil
}
func (s *stubWBSvc) OnDeadheadConfirmed(_ context.Context, _ events.AssignmentPayload) error {
	return nil
}

type stubAssignNotifySvc struct{}

func (s *stubAssignNotifySvc) OnBOLWorkflowCompleted(_ context.Context, _ events.BOLCompletedPayload) error {
	return nil
}

func newAssignHandler(
	assign *stubAssignRepo,
	driver *stubDriverRepo,
	bol *stubBOLRepo,
	equip *stubAssignEquipRepo,
	hos *stubHOSSvc,
) *AssignmentHandler {
	return NewAssignmentHandler(assign, driver, bol, equip, hos, &stubWBSvc{}, &stubAssignNotifySvc{})
}

// --- Get ---

func TestAssignmentGet_BadUUID_Returns400(t *testing.T) {
	h := newAssignHandler(&stubAssignRepo{}, &stubDriverRepo{}, &stubBOLRepo{}, &stubAssignEquipRepo{}, &stubHOSSvc{})
	req := withIDParam(httptest.NewRequest(http.MethodGet, "/", nil), "not-a-uuid")
	rec := httptest.NewRecorder()
	h.Get(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestAssignmentGet_NotFound_Returns404(t *testing.T) {
	h := newAssignHandler(&stubAssignRepo{}, &stubDriverRepo{}, &stubBOLRepo{}, &stubAssignEquipRepo{}, &stubHOSSvc{})
	req := withIDParam(httptest.NewRequest(http.MethodGet, "/", nil), uuid.New().String())
	rec := httptest.NewRecorder()
	h.Get(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestAssignmentGet_Success_Returns200(t *testing.T) {
	assign := &models.DriverBOLAssignment{ID: uuid.New(), DriverID: uuid.New(), PlanBOLID: uuid.New()}
	h := newAssignHandler(&stubAssignRepo{assignment: assign}, &stubDriverRepo{}, &stubBOLRepo{}, &stubAssignEquipRepo{}, &stubHOSSvc{})
	req := withIDParam(httptest.NewRequest(http.MethodGet, "/", nil), assign.ID.String())
	rec := httptest.NewRecorder()
	h.Get(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

// --- Create ---

func TestAssignmentCreate_BOLNotValidated_Returns409(t *testing.T) {
	bol := &models.PlanBOLRecord{ID: uuid.New(), Status: models.PlanBOLStatusLoading}
	equip := &models.Equipment{ID: uuid.New(), Status: models.EquipmentStatusAvailable}
	h := newAssignHandler(&stubAssignRepo{}, &stubDriverRepo{}, &stubBOLRepo{bol: bol}, &stubAssignEquipRepo{equipment: equip}, &stubHOSSvc{})
	body := map[string]any{
		"driver_id":    uuid.New().String(),
		"plan_bol_id":  bol.ID.String(),
		"equipment_id": equip.ID.String(),
		"state_code":   "IL",
		"cycle_label":  "60h/7d",
	}
	req := httptest.NewRequest(http.MethodPost, "/api/assignment", postBody(t, body))
	rec := httptest.NewRecorder()
	h.Create(rec, req)
	assert.Equal(t, http.StatusConflict, rec.Code)
}

func TestAssignmentCreate_EquipmentNotAvailable_Returns409(t *testing.T) {
	bol := &models.PlanBOLRecord{ID: uuid.New(), Status: models.PlanBOLStatusValidated}
	equip := &models.Equipment{ID: uuid.New(), Status: models.EquipmentStatusAssigned}
	h := newAssignHandler(&stubAssignRepo{}, &stubDriverRepo{}, &stubBOLRepo{bol: bol}, &stubAssignEquipRepo{equipment: equip}, &stubHOSSvc{})
	body := map[string]any{
		"driver_id":    uuid.New().String(),
		"plan_bol_id":  bol.ID.String(),
		"equipment_id": equip.ID.String(),
		"state_code":   "IL",
		"cycle_label":  "60h/7d",
	}
	req := httptest.NewRequest(http.MethodPost, "/api/assignment", postBody(t, body))
	rec := httptest.NewRecorder()
	h.Create(rec, req)
	assert.Equal(t, http.StatusConflict, rec.Code)
}

func TestAssignmentCreate_HOSViolation_Returns422(t *testing.T) {
	bol := &models.PlanBOLRecord{ID: uuid.New(), Status: models.PlanBOLStatusValidated}
	equip := &models.Equipment{ID: uuid.New(), Status: models.EquipmentStatusAvailable}
	hos := &stubHOSSvc{err: errNotFound} // any non-nil error simulates HOS rejection
	h := newAssignHandler(&stubAssignRepo{}, &stubDriverRepo{}, &stubBOLRepo{bol: bol}, &stubAssignEquipRepo{equipment: equip}, hos)
	body := map[string]any{
		"driver_id":           uuid.New().String(),
		"plan_bol_id":         bol.ID.String(),
		"equipment_id":        equip.ID.String(),
		"state_code":          "IL",
		"cycle_label":         "60h/7d",
		"estimated_run_hours": 12.0,
	}
	req := httptest.NewRequest(http.MethodPost, "/api/assignment", postBody(t, body))
	rec := httptest.NewRecorder()
	h.Create(rec, req)
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestAssignmentCreate_NilBody_Returns400(t *testing.T) {
	h := newAssignHandler(&stubAssignRepo{}, &stubDriverRepo{}, &stubBOLRepo{}, &stubAssignEquipRepo{}, &stubHOSSvc{})
	req := httptest.NewRequest(http.MethodPost, "/api/assignment", nil)
	rec := httptest.NewRecorder()
	h.Create(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestAssignmentCreate_InvalidDriverID_Returns400(t *testing.T) {
	h := newAssignHandler(&stubAssignRepo{}, &stubDriverRepo{}, &stubBOLRepo{}, &stubAssignEquipRepo{}, &stubHOSSvc{})
	body := map[string]any{
		"driver_id":    "not-a-uuid",
		"plan_bol_id":  uuid.New().String(),
		"equipment_id": uuid.New().String(),
		"state_code":   "IL",
		"cycle_label":  "60h/7d",
	}
	req := httptest.NewRequest(http.MethodPost, "/api/assignment", postBody(t, body))
	rec := httptest.NewRecorder()
	h.Create(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestAssignmentCreate_BOLNotFound_Returns404(t *testing.T) {
	h := newAssignHandler(&stubAssignRepo{}, &stubDriverRepo{}, &stubBOLRepo{}, &stubAssignEquipRepo{}, &stubHOSSvc{})
	body := map[string]any{
		"driver_id":    uuid.New().String(),
		"plan_bol_id":  uuid.New().String(),
		"equipment_id": uuid.New().String(),
		"state_code":   "IL",
		"cycle_label":  "60h/7d",
	}
	req := httptest.NewRequest(http.MethodPost, "/api/assignment", postBody(t, body))
	rec := httptest.NewRecorder()
	h.Create(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestAssignmentCreate_EquipmentNotFound_Returns404(t *testing.T) {
	bol := &models.PlanBOLRecord{ID: uuid.New(), Status: models.PlanBOLStatusValidated}
	h := newAssignHandler(&stubAssignRepo{}, &stubDriverRepo{}, &stubBOLRepo{bol: bol}, &stubAssignEquipRepo{}, &stubHOSSvc{})
	body := map[string]any{
		"driver_id":    uuid.New().String(),
		"plan_bol_id":  bol.ID.String(),
		"equipment_id": uuid.New().String(),
		"state_code":   "IL",
		"cycle_label":  "60h/7d",
	}
	req := httptest.NewRequest(http.MethodPost, "/api/assignment", postBody(t, body))
	rec := httptest.NewRecorder()
	h.Create(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestAssignmentCreate_Success_Returns201(t *testing.T) {
	bol := &models.PlanBOLRecord{ID: uuid.New(), Status: models.PlanBOLStatusValidated}
	equip := &models.Equipment{ID: uuid.New(), Status: models.EquipmentStatusAvailable}
	h := newAssignHandler(&stubAssignRepo{}, &stubDriverRepo{}, &stubBOLRepo{bol: bol}, &stubAssignEquipRepo{equipment: equip}, &stubHOSSvc{})
	body := map[string]any{
		"driver_id":           uuid.New().String(),
		"plan_bol_id":         bol.ID.String(),
		"equipment_id":        equip.ID.String(),
		"state_code":          "IL",
		"cycle_label":         "60h/7d",
		"estimated_run_hours": 4.0,
	}
	req := httptest.NewRequest(http.MethodPost, "/api/assignment", postBody(t, body))
	rec := httptest.NewRecorder()
	h.Create(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)
}

func TestAssignmentCreate_MissingStateCode_Returns400(t *testing.T) {
	bol := &models.PlanBOLRecord{ID: uuid.New(), Status: models.PlanBOLStatusValidated}
	equip := &models.Equipment{ID: uuid.New(), Status: models.EquipmentStatusAvailable}
	h := newAssignHandler(&stubAssignRepo{}, &stubDriverRepo{}, &stubBOLRepo{bol: bol}, &stubAssignEquipRepo{equipment: equip}, &stubHOSSvc{})
	body := map[string]any{
		"driver_id":    uuid.New().String(),
		"plan_bol_id":  bol.ID.String(),
		"equipment_id": equip.ID.String(),
		// state_code missing
	}
	req := httptest.NewRequest(http.MethodPost, "/api/assignment", postBody(t, body))
	rec := httptest.NewRecorder()
	h.Create(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// --- Depart ---

func TestAssignmentDepart_BadUUID_Returns400(t *testing.T) {
	h := newAssignHandler(&stubAssignRepo{}, &stubDriverRepo{}, &stubBOLRepo{}, &stubAssignEquipRepo{}, &stubHOSSvc{})
	req := withIDParam(httptest.NewRequest(http.MethodPatch, "/", nil), "not-a-uuid")
	rec := httptest.NewRecorder()
	h.Depart(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestAssignmentDepart_NotFound_Returns404(t *testing.T) {
	h := newAssignHandler(&stubAssignRepo{}, &stubDriverRepo{}, &stubBOLRepo{}, &stubAssignEquipRepo{}, &stubHOSSvc{})
	req := withIDParam(httptest.NewRequest(http.MethodPatch, "/", nil), uuid.New().String())
	rec := httptest.NewRecorder()
	h.Depart(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestAssignmentDepart_AlreadyDeparted_Returns409(t *testing.T) {
	departed := time.Now()
	assign := &models.DriverBOLAssignment{ID: uuid.New(), DepartedAt: &departed}
	h := newAssignHandler(&stubAssignRepo{assignment: assign}, &stubDriverRepo{}, &stubBOLRepo{bol: &models.PlanBOLRecord{}}, &stubAssignEquipRepo{}, &stubHOSSvc{})
	req := withIDParam(httptest.NewRequest(http.MethodPatch, "/", nil), assign.ID.String())
	rec := httptest.NewRecorder()
	h.Depart(rec, req)
	assert.Equal(t, http.StatusConflict, rec.Code)
}

func TestAssignmentDepart_Success_Returns204(t *testing.T) {
	assign := &models.DriverBOLAssignment{ID: uuid.New(), PlanBOLID: uuid.New()}
	h := newAssignHandler(&stubAssignRepo{assignment: assign}, &stubDriverRepo{}, &stubBOLRepo{bol: &models.PlanBOLRecord{}}, &stubAssignEquipRepo{}, &stubHOSSvc{})
	req := withIDParam(httptest.NewRequest(http.MethodPatch, "/", nil), assign.ID.String())
	rec := httptest.NewRecorder()
	h.Depart(rec, req)
	assert.Equal(t, http.StatusNoContent, rec.Code)
}

// --- Fulfill ---

func TestAssignmentFulfill_BadUUID_Returns400(t *testing.T) {
	h := newAssignHandler(&stubAssignRepo{}, &stubDriverRepo{}, &stubBOLRepo{}, &stubAssignEquipRepo{}, &stubHOSSvc{})
	req := withIDParam(httptest.NewRequest(http.MethodPatch, "/", nil), "not-a-uuid")
	rec := httptest.NewRecorder()
	h.Fulfill(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestAssignmentFulfill_NotFound_Returns404(t *testing.T) {
	h := newAssignHandler(&stubAssignRepo{}, &stubDriverRepo{}, &stubBOLRepo{}, &stubAssignEquipRepo{}, &stubHOSSvc{})
	req := withIDParam(httptest.NewRequest(http.MethodPatch, "/", nil), uuid.New().String())
	rec := httptest.NewRecorder()
	h.Fulfill(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestAssignmentFulfill_AlreadyFulfilled_Returns409(t *testing.T) {
	fulfilledAt := time.Now()
	assign := &models.DriverBOLAssignment{ID: uuid.New(), FulfilledAt: &fulfilledAt}
	h := newAssignHandler(&stubAssignRepo{assignment: assign}, &stubDriverRepo{}, &stubBOLRepo{bol: &models.PlanBOLRecord{}}, &stubAssignEquipRepo{}, &stubHOSSvc{})
	req := withIDParam(httptest.NewRequest(http.MethodPatch, "/", nil), assign.ID.String())
	rec := httptest.NewRecorder()
	h.Fulfill(rec, req)
	assert.Equal(t, http.StatusConflict, rec.Code)
}

func TestAssignmentFulfill_Success_Returns204(t *testing.T) {
	assign := &models.DriverBOLAssignment{ID: uuid.New(), PlanBOLID: uuid.New()}
	h := newAssignHandler(&stubAssignRepo{assignment: assign}, &stubDriverRepo{}, &stubBOLRepo{bol: &models.PlanBOLRecord{}}, &stubAssignEquipRepo{}, &stubHOSSvc{})
	req := withIDParam(httptest.NewRequest(http.MethodPatch, "/", nil), assign.ID.String())
	rec := httptest.NewRecorder()
	h.Fulfill(rec, req)
	assert.Equal(t, http.StatusNoContent, rec.Code)
}

// --- ConfirmDeadhead ---

func TestAssignmentConfirmDeadhead_BadUUID_Returns400(t *testing.T) {
	h := newAssignHandler(&stubAssignRepo{}, &stubDriverRepo{}, &stubBOLRepo{}, &stubAssignEquipRepo{}, &stubHOSSvc{})
	req := withIDParam(httptest.NewRequest(http.MethodPatch, "/", nil), "not-a-uuid")
	rec := httptest.NewRecorder()
	h.ConfirmDeadhead(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestAssignmentConfirmDeadhead_NotFound_Returns404(t *testing.T) {
	h := newAssignHandler(&stubAssignRepo{}, &stubDriverRepo{}, &stubBOLRepo{}, &stubAssignEquipRepo{}, &stubHOSSvc{})
	req := withIDParam(httptest.NewRequest(http.MethodPatch, "/", nil), uuid.New().String())
	rec := httptest.NewRecorder()
	h.ConfirmDeadhead(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestAssignmentConfirmDeadhead_NotYetFulfilled_Returns409(t *testing.T) {
	assign := &models.DriverBOLAssignment{ID: uuid.New(), FulfilledAt: nil}
	h := newAssignHandler(&stubAssignRepo{assignment: assign}, &stubDriverRepo{}, &stubBOLRepo{}, &stubAssignEquipRepo{}, &stubHOSSvc{})
	req := withIDParam(httptest.NewRequest(http.MethodPatch, "/", nil), assign.ID.String())
	rec := httptest.NewRecorder()
	h.ConfirmDeadhead(rec, req)
	assert.Equal(t, http.StatusConflict, rec.Code)
}

func TestAssignmentConfirmDeadhead_AlreadyConfirmed_Returns409(t *testing.T) {
	fulfilledAt := time.Now().Add(-1 * time.Hour)
	confirmedAt := time.Now()
	assign := &models.DriverBOLAssignment{
		ID:                  uuid.New(),
		FulfilledAt:         &fulfilledAt,
		DeadheadConfirmedAt: &confirmedAt,
	}
	h := newAssignHandler(&stubAssignRepo{assignment: assign}, &stubDriverRepo{}, &stubBOLRepo{}, &stubAssignEquipRepo{}, &stubHOSSvc{})
	req := withIDParam(httptest.NewRequest(http.MethodPatch, "/", nil), assign.ID.String())
	rec := httptest.NewRecorder()
	h.ConfirmDeadhead(rec, req)
	assert.Equal(t, http.StatusConflict, rec.Code)
}

func TestAssignmentConfirmDeadhead_Success_Returns204(t *testing.T) {
	fulfilledAt := time.Now().Add(-1 * time.Hour)
	assign := &models.DriverBOLAssignment{ID: uuid.New(), PlanBOLID: uuid.New(), FulfilledAt: &fulfilledAt}
	h := newAssignHandler(&stubAssignRepo{assignment: assign}, &stubDriverRepo{}, &stubBOLRepo{}, &stubAssignEquipRepo{}, &stubHOSSvc{})
	req := withIDParam(httptest.NewRequest(http.MethodPatch, "/", nil), assign.ID.String())
	rec := httptest.NewRecorder()
	h.ConfirmDeadhead(rec, req)
	assert.Equal(t, http.StatusNoContent, rec.Code)
}
