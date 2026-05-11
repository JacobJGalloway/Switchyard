package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/JacobJGalloway/switchyard-go/internal/events"
	"github.com/JacobJGalloway/switchyard-go/internal/models"
	"github.com/JacobJGalloway/switchyard-go/internal/repository"
)

// --- stubs ---

type stubHOSRepo struct {
	window *models.HOSWindow
	err    error
}

func (r *stubHOSRepo) CreateLimit(_ context.Context, _ *models.HOSLimit) error  { return nil }
func (r *stubHOSRepo) GetLimitByStateAndCycle(_ context.Context, _, _ string) (*models.HOSLimit, error) {
	return nil, nil
}
func (r *stubHOSRepo) CreateWindow(_ context.Context, _ *models.HOSWindow) error { return nil }
func (r *stubHOSRepo) GetWindowByDriver(_ context.Context, _ uuid.UUID) (*models.HOSWindow, error) {
	if r.err != nil {
		return nil, r.err
	}
	return r.window, nil
}
func (r *stubHOSRepo) UpdateWindow(_ context.Context, _ *models.HOSWindow) error { return nil }

type stubDriverHOSSvc struct{ err error }

func (s *stubDriverHOSSvc) OnStopLogged(_ context.Context, _ events.StopLoggedPayload) error {
	return s.err
}

// activeAssignRepo overrides GetActiveByDriver to return a configured assignment.
type activeAssignRepo struct {
	stubAssignRepo
	active *models.DriverBOLAssignment
}

func (r *activeAssignRepo) GetActiveByDriver(_ context.Context, _ uuid.UUID) (*models.DriverBOLAssignment, error) {
	if r.active == nil {
		return nil, errNotFound
	}
	return r.active, nil
}

func newDriverHandler(
	driverRepo repository.DriverRepository,
	bolRepo repository.PlanBOLRepository,
	hosRepo repository.HOSRepository,
	assignRepo repository.AssignmentRepository,
	hosSvc driverHOSService,
) *DriverHandler {
	return NewDriverHandler(driverRepo, bolRepo, hosRepo, assignRepo, hosSvc, &stubLogisticsClient{}, nil)
}

// --- GetAll ---

func TestDriver_GetAll_RepoError_Returns500(t *testing.T) {
	h := newDriverHandler(&stubDriverRepo{getAllErr: errNotFound}, &stubBOLRepo{}, &stubHOSRepo{}, &stubAssignRepo{}, &stubDriverHOSSvc{})
	req := httptest.NewRequest(http.MethodGet, "/api/driver", nil)
	rec := httptest.NewRecorder()
	h.GetAll(rec, req)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestDriver_GetAll_Returns200(t *testing.T) {
	h := newDriverHandler(&stubDriverRepo{}, &stubBOLRepo{}, &stubHOSRepo{}, &stubAssignRepo{}, &stubDriverHOSSvc{})
	req := httptest.NewRequest(http.MethodGet, "/api/driver", nil)
	rec := httptest.NewRecorder()
	h.GetAll(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestDriver_GetAll_WithDrivers_Returns200(t *testing.T) {
	drivers := []*models.Driver{
		{ID: uuid.New(), Name: "Alice", LicenseState: "IL", IsActive: true},
		{ID: uuid.New(), Name: "Bob", LicenseState: "WI", IsActive: true},
	}
	h := newDriverHandler(&stubDriverRepo{all: drivers}, &stubBOLRepo{}, &stubHOSRepo{}, &stubAssignRepo{}, &stubDriverHOSSvc{})
	req := httptest.NewRequest(http.MethodGet, "/api/driver", nil)
	rec := httptest.NewRecorder()
	h.GetAll(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

// --- GetRunsheet ---

func TestDriver_GetRunsheet_NotFound_Returns404(t *testing.T) {
	h := newDriverHandler(&stubDriverRepo{err: errNotFound}, &stubBOLRepo{}, &stubHOSRepo{}, &stubAssignRepo{}, &stubDriverHOSSvc{})
	req := withIDParam(httptest.NewRequest(http.MethodGet, "/", nil), uuid.New().String())
	rec := httptest.NewRecorder()
	h.GetRunsheet(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestDriver_GetRunsheet_NoAssignment_Returns200(t *testing.T) {
	driver := &models.Driver{ID: uuid.New(), Name: "Alice"}
	h := newDriverHandler(&stubDriverRepo{driver: driver}, &stubBOLRepo{}, &stubHOSRepo{}, &stubAssignRepo{}, &stubDriverHOSSvc{})
	req := withIDParam(httptest.NewRequest(http.MethodGet, "/", nil), driver.ID.String())
	rec := httptest.NewRecorder()
	h.GetRunsheet(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestDriver_GetRunsheet_WithActiveBOL_Returns200(t *testing.T) {
	driver := &models.Driver{ID: uuid.New(), Name: "Bob"}
	plan := bolWithStatus(models.PlanBOLStatusSubmitted)
	assign := &models.DriverBOLAssignment{ID: uuid.New(), PlanBOLID: plan.ID}
	ar := &activeAssignRepo{active: assign}
	h := NewDriverHandler(&stubDriverRepo{driver: driver}, &stubBOLRepo{bol: plan}, &stubHOSRepo{}, ar, &stubDriverHOSSvc{}, &stubLogisticsClient{}, nil)
	req := withIDParam(httptest.NewRequest(http.MethodGet, "/", nil), driver.ID.String())
	rec := httptest.NewRecorder()
	h.GetRunsheet(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

// --- GetActiveBOL ---

func TestDriver_GetActiveBOL_NoActive_Returns404(t *testing.T) {
	h := newDriverHandler(&stubDriverRepo{}, &stubBOLRepo{}, &stubHOSRepo{}, &stubAssignRepo{}, &stubDriverHOSSvc{})
	req := withIDParam(httptest.NewRequest(http.MethodGet, "/", nil), uuid.New().String())
	rec := httptest.NewRecorder()
	h.GetActiveBOL(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestDriver_GetActiveBOL_Success_Returns200(t *testing.T) {
	plan := bolWithStatus(models.PlanBOLStatusSubmitted)
	assign := &models.DriverBOLAssignment{ID: uuid.New(), PlanBOLID: plan.ID}
	ar := &activeAssignRepo{active: assign}
	h := NewDriverHandler(&stubDriverRepo{}, &stubBOLRepo{bol: plan}, &stubHOSRepo{}, ar, &stubDriverHOSSvc{}, &stubLogisticsClient{}, nil)
	req := withIDParam(httptest.NewRequest(http.MethodGet, "/", nil), uuid.New().String())
	rec := httptest.NewRecorder()
	h.GetActiveBOL(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

// --- GetHOS ---

func TestDriver_GetHOS_NotFound_Returns404(t *testing.T) {
	h := newDriverHandler(&stubDriverRepo{}, &stubBOLRepo{}, &stubHOSRepo{err: errors.New("not found")}, &stubAssignRepo{}, &stubDriverHOSSvc{})
	req := withIDParam(httptest.NewRequest(http.MethodGet, "/", nil), uuid.New().String())
	rec := httptest.NewRecorder()
	h.GetHOS(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestDriver_GetHOS_Success_Returns200(t *testing.T) {
	window := &models.HOSWindow{ID: uuid.New(), DriverID: uuid.New()}
	h := newDriverHandler(&stubDriverRepo{}, &stubBOLRepo{}, &stubHOSRepo{window: window}, &stubAssignRepo{}, &stubDriverHOSSvc{})
	req := withIDParam(httptest.NewRequest(http.MethodGet, "/", nil), uuid.New().String())
	rec := httptest.NewRecorder()
	h.GetHOS(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

// --- LogStop ---

func logStopReq(t *testing.T, driverID, stopID uuid.UUID) *http.Request {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", driverID.String())
	rctx.URLParams.Add("stopId", stopID.String())
	return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
}

func TestDriver_LogStop_StopNotFound_Returns404(t *testing.T) {
	h := newDriverHandler(&stubDriverRepo{}, &stubBOLRepo{stopErr: errNotFound}, &stubHOSRepo{}, &stubAssignRepo{}, &stubDriverHOSSvc{})
	rec := httptest.NewRecorder()
	h.LogStop(rec, logStopReq(t, uuid.New(), uuid.New()))
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestDriver_LogStop_HOSRejection_Returns422(t *testing.T) {
	stop := &models.PlanBOLStop{ID: uuid.New(), PlanBOLID: uuid.New(), LocationID: "ST0001", StopType: models.StopTypeStore}
	h := newDriverHandler(
		&stubDriverRepo{},
		&stubBOLRepo{stop: stop},
		&stubHOSRepo{},
		&stubAssignRepo{},
		&stubDriverHOSSvc{err: errors.New("HOS daily limit exceeded")},
	)
	rec := httptest.NewRecorder()
	h.LogStop(rec, logStopReq(t, uuid.New(), stop.ID))
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestDriver_LogStop_Success_Returns204(t *testing.T) {
	now := time.Now()
	stop := &models.PlanBOLStop{
		ID: uuid.New(), PlanBOLID: uuid.New(),
		LocationID: "ST0001", StopType: models.StopTypeStore,
		IsProcessed: true, ProcessedAt: &now,
	}
	h := newDriverHandler(
		&stubDriverRepo{},
		&stubBOLRepo{stop: stop, stops: []*models.PlanBOLStop{stop}},
		&stubHOSRepo{},
		&stubAssignRepo{},
		&stubDriverHOSSvc{},
	)
	rec := httptest.NewRecorder()
	h.LogStop(rec, logStopReq(t, uuid.New(), stop.ID))
	assert.Equal(t, http.StatusNoContent, rec.Code)
}
