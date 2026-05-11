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
	"github.com/stretchr/testify/require"

	"github.com/JacobJGalloway/switchyard-go/internal/integrations"
	"github.com/JacobJGalloway/switchyard-go/internal/models"
	"github.com/JacobJGalloway/switchyard-go/internal/services"
)

// --- stubs ---

type stubPlanBOLSvc struct {
	plan       *models.PlanBOLRecord
	violations []string
	err        error
}

func (s *stubPlanBOLSvc) PlanRoute(_ context.Context, _ services.PlanRouteInput) (*models.PlanBOLRecord, error) {
	return s.plan, s.err
}
func (s *stubPlanBOLSvc) ValidatePlan(_ context.Context, _ uuid.UUID) ([]string, error) {
	return s.violations, s.err
}

type stubLogisticsClient struct {
	txID string
	err  error
}

func (s *stubLogisticsClient) CreateBOL(_ context.Context, _ *integrations.CreateBOLRequest) (string, error) {
	return s.txID, s.err
}
func (s *stubLogisticsClient) ProcessStop(_ context.Context, _, _ string) error { return nil }
func (s *stubLogisticsClient) ReplaceStop(_ context.Context, _ string, _ *integrations.ReplaceStopRequest) error {
	return nil
}

func newPlanBOLHandler(svc *stubPlanBOLSvc, bol *stubBOLRepo, log *stubLogisticsClient) *PlanBOLHandler {
	return NewPlanBOLHandler(svc, bol, log)
}

func bolWithStatus(status models.PlanBOLStatus) *models.PlanBOLRecord {
	return &models.PlanBOLRecord{
		ID:              uuid.New(),
		DriverID:        uuid.New(),
		OriginatingWhID: "WH001",
		Status:          status,
		CreatedAt:       time.Now(),
	}
}

// --- Create ---

func TestPlanBOL_Create_BadBody_Returns400(t *testing.T) {
	h := newPlanBOLHandler(&stubPlanBOLSvc{}, &stubBOLRepo{}, &stubLogisticsClient{})
	req := httptest.NewRequest(http.MethodPost, "/api/plan-bol", nil)
	rec := httptest.NewRecorder()
	h.Create(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestPlanBOL_Create_MissingDriverID_Returns400(t *testing.T) {
	h := newPlanBOLHandler(&stubPlanBOLSvc{}, &stubBOLRepo{}, &stubLogisticsClient{})
	body := map[string]any{"origin_warehouse_id": "WH001", "store_stops": []any{map[string]any{"location_id": "ST0001"}}}
	req := httptest.NewRequest(http.MethodPost, "/api/plan-bol", postBody(t, body))
	rec := httptest.NewRecorder()
	h.Create(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestPlanBOL_Create_MissingOrigin_Returns400(t *testing.T) {
	h := newPlanBOLHandler(&stubPlanBOLSvc{}, &stubBOLRepo{}, &stubLogisticsClient{})
	body := map[string]any{"driver_id": uuid.New().String(), "store_stops": []any{map[string]any{"location_id": "ST0001"}}}
	req := httptest.NewRequest(http.MethodPost, "/api/plan-bol", postBody(t, body))
	rec := httptest.NewRecorder()
	h.Create(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestPlanBOL_Create_NoStops_Returns400(t *testing.T) {
	h := newPlanBOLHandler(&stubPlanBOLSvc{}, &stubBOLRepo{}, &stubLogisticsClient{})
	body := map[string]any{"driver_id": uuid.New().String(), "origin_warehouse_id": "WH001", "store_stops": []any{}}
	req := httptest.NewRequest(http.MethodPost, "/api/plan-bol", postBody(t, body))
	rec := httptest.NewRecorder()
	h.Create(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestPlanBOL_Create_ServiceError_Returns422(t *testing.T) {
	h := newPlanBOLHandler(&stubPlanBOLSvc{err: errors.New("inventory shortfall")}, &stubBOLRepo{}, &stubLogisticsClient{})
	body := map[string]any{
		"driver_id": uuid.New().String(), "origin_warehouse_id": "WH001",
		"store_stops": []any{map[string]any{"location_id": "ST0001", "items": map[string]int{"SKU-A": 2}}},
	}
	req := httptest.NewRequest(http.MethodPost, "/api/plan-bol", postBody(t, body))
	rec := httptest.NewRecorder()
	h.Create(rec, req)
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestPlanBOL_Create_Success_Returns201(t *testing.T) {
	plan := bolWithStatus(models.PlanBOLStatusDraft)
	h := newPlanBOLHandler(&stubPlanBOLSvc{plan: plan}, &stubBOLRepo{}, &stubLogisticsClient{})
	body := map[string]any{
		"driver_id": uuid.New().String(), "origin_warehouse_id": "WH001",
		"store_stops": []any{map[string]any{"location_id": "ST0001", "items": map[string]int{"SKU-A": 2}}},
	}
	req := httptest.NewRequest(http.MethodPost, "/api/plan-bol", postBody(t, body))
	rec := httptest.NewRecorder()
	h.Create(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)
}

// --- Get ---

func TestPlanBOL_Get_BadUUID_Returns400(t *testing.T) {
	h := newPlanBOLHandler(&stubPlanBOLSvc{}, &stubBOLRepo{}, &stubLogisticsClient{})
	req := withIDParam(httptest.NewRequest(http.MethodGet, "/", nil), "not-a-uuid")
	rec := httptest.NewRecorder()
	h.Get(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestPlanBOL_Get_NotFound_Returns404(t *testing.T) {
	h := newPlanBOLHandler(&stubPlanBOLSvc{}, &stubBOLRepo{}, &stubLogisticsClient{})
	req := httptest.NewRequest(http.MethodGet, "/api/plan-bol/"+uuid.New().String(), nil)
	req = withIDParam(req, uuid.New().String())
	rec := httptest.NewRecorder()
	h.Get(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestPlanBOL_Get_Success_Returns200(t *testing.T) {
	plan := bolWithStatus(models.PlanBOLStatusDraft)
	h := newPlanBOLHandler(&stubPlanBOLSvc{}, &stubBOLRepo{bol: plan}, &stubLogisticsClient{})
	req := withIDParam(httptest.NewRequest(http.MethodGet, "/", nil), plan.ID.String())
	rec := httptest.NewRecorder()
	h.Get(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

// --- BeginPlanning ---

func TestPlanBOL_BeginPlanning_BadUUID_Returns400(t *testing.T) {
	h := newPlanBOLHandler(&stubPlanBOLSvc{}, &stubBOLRepo{}, &stubLogisticsClient{})
	req := withIDParam(httptest.NewRequest(http.MethodPost, "/", nil), "not-a-uuid")
	rec := httptest.NewRecorder()
	h.BeginPlanning(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestPlanBOL_BeginPlanning_NotFound_Returns404(t *testing.T) {
	h := newPlanBOLHandler(&stubPlanBOLSvc{}, &stubBOLRepo{}, &stubLogisticsClient{})
	req := withIDParam(httptest.NewRequest(http.MethodPost, "/", nil), uuid.New().String())
	rec := httptest.NewRecorder()
	h.BeginPlanning(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestPlanBOL_BeginPlanning_WrongStatus_Returns409(t *testing.T) {
	plan := bolWithStatus(models.PlanBOLStatusLoading)
	h := newPlanBOLHandler(&stubPlanBOLSvc{}, &stubBOLRepo{bol: plan}, &stubLogisticsClient{})
	req := withIDParam(httptest.NewRequest(http.MethodPost, "/", nil), plan.ID.String())
	rec := httptest.NewRecorder()
	h.BeginPlanning(rec, req)
	assert.Equal(t, http.StatusConflict, rec.Code)
}

func TestPlanBOL_BeginPlanning_Success_Returns204(t *testing.T) {
	plan := bolWithStatus(models.PlanBOLStatusDraft)
	h := newPlanBOLHandler(&stubPlanBOLSvc{}, &stubBOLRepo{bol: plan}, &stubLogisticsClient{})
	req := withIDParam(httptest.NewRequest(http.MethodPost, "/", nil), plan.ID.String())
	rec := httptest.NewRecorder()
	h.BeginPlanning(rec, req)
	assert.Equal(t, http.StatusNoContent, rec.Code)
}

// --- MarkLoaded ---

func TestPlanBOL_MarkLoaded_BadUUID_Returns400(t *testing.T) {
	h := newPlanBOLHandler(&stubPlanBOLSvc{}, &stubBOLRepo{}, &stubLogisticsClient{})
	req := withIDParam(httptest.NewRequest(http.MethodPatch, "/", nil), "not-a-uuid")
	rec := httptest.NewRecorder()
	h.MarkLoaded(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestPlanBOL_MarkLoaded_NotFound_Returns404(t *testing.T) {
	h := newPlanBOLHandler(&stubPlanBOLSvc{}, &stubBOLRepo{}, &stubLogisticsClient{})
	req := withIDParam(httptest.NewRequest(http.MethodPatch, "/", nil), uuid.New().String())
	rec := httptest.NewRecorder()
	h.MarkLoaded(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestPlanBOL_MarkLoaded_WrongStatus_Returns409(t *testing.T) {
	plan := bolWithStatus(models.PlanBOLStatusDraft)
	h := newPlanBOLHandler(&stubPlanBOLSvc{}, &stubBOLRepo{bol: plan}, &stubLogisticsClient{})
	req := withIDParam(httptest.NewRequest(http.MethodPatch, "/", nil), plan.ID.String())
	rec := httptest.NewRecorder()
	h.MarkLoaded(rec, req)
	assert.Equal(t, http.StatusConflict, rec.Code)
}

func TestPlanBOL_MarkLoaded_Success_Returns204(t *testing.T) {
	plan := bolWithStatus(models.PlanBOLStatusLoading)
	h := newPlanBOLHandler(&stubPlanBOLSvc{}, &stubBOLRepo{bol: plan}, &stubLogisticsClient{})
	req := withIDParam(httptest.NewRequest(http.MethodPatch, "/", nil), plan.ID.String())
	rec := httptest.NewRecorder()
	h.MarkLoaded(rec, req)
	assert.Equal(t, http.StatusNoContent, rec.Code)
}

// --- Commit ---

func TestPlanBOL_Commit_BadUUID_Returns400(t *testing.T) {
	h := newPlanBOLHandler(&stubPlanBOLSvc{}, &stubBOLRepo{}, &stubLogisticsClient{})
	req := withIDParam(httptest.NewRequest(http.MethodPost, "/", nil), "not-a-uuid")
	rec := httptest.NewRecorder()
	h.Commit(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestPlanBOL_Commit_NotFound_Returns404(t *testing.T) {
	h := newPlanBOLHandler(&stubPlanBOLSvc{}, &stubBOLRepo{}, &stubLogisticsClient{})
	body := map[string]any{"customer_first_name": "Jane", "customer_last_name": "Doe", "city": "Chicago", "state": "IL"}
	req := withIDParam(httptest.NewRequest(http.MethodPost, "/", postBody(t, body)), uuid.New().String())
	rec := httptest.NewRecorder()
	h.Commit(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestPlanBOL_Commit_MissingFields_Returns400(t *testing.T) {
	plan := bolWithStatus(models.PlanBOLStatusPlanProgress)
	h := newPlanBOLHandler(&stubPlanBOLSvc{}, &stubBOLRepo{bol: plan}, &stubLogisticsClient{})
	body := map[string]any{"customer_first_name": "Jane"}
	req := withIDParam(httptest.NewRequest(http.MethodPost, "/", postBody(t, body)), plan.ID.String())
	rec := httptest.NewRecorder()
	h.Commit(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestPlanBOL_Commit_WrongStatus_Returns409(t *testing.T) {
	plan := bolWithStatus(models.PlanBOLStatusDraft)
	h := newPlanBOLHandler(&stubPlanBOLSvc{}, &stubBOLRepo{bol: plan}, &stubLogisticsClient{})
	body := map[string]any{"customer_first_name": "Jane", "customer_last_name": "Doe", "city": "Chicago", "state": "IL"}
	req := withIDParam(httptest.NewRequest(http.MethodPost, "/", postBody(t, body)), plan.ID.String())
	rec := httptest.NewRecorder()
	h.Commit(rec, req)
	assert.Equal(t, http.StatusConflict, rec.Code)
}

func TestPlanBOL_Commit_LogisticsError_Returns502(t *testing.T) {
	plan := bolWithStatus(models.PlanBOLStatusPlanProgress)
	lc := &stubLogisticsClient{err: errors.New("downstream unavailable")}
	h := newPlanBOLHandler(&stubPlanBOLSvc{}, &stubBOLRepo{bol: plan}, lc)
	body := map[string]any{"customer_first_name": "Jane", "customer_last_name": "Doe", "city": "Chicago", "state": "IL"}
	req := withIDParam(httptest.NewRequest(http.MethodPost, "/", postBody(t, body)), plan.ID.String())
	rec := httptest.NewRecorder()
	h.Commit(rec, req)
	assert.Equal(t, http.StatusBadGateway, rec.Code)
}

func TestPlanBOL_Commit_Success_Returns200(t *testing.T) {
	plan := bolWithStatus(models.PlanBOLStatusPlanProgress)
	lc := &stubLogisticsClient{txID: "TX-001"}
	h := newPlanBOLHandler(&stubPlanBOLSvc{}, &stubBOLRepo{bol: plan}, lc)
	body := map[string]any{"customer_first_name": "Jane", "customer_last_name": "Doe", "city": "Chicago", "state": "IL"}
	req := withIDParam(httptest.NewRequest(http.MethodPost, "/", postBody(t, body)), plan.ID.String())
	rec := httptest.NewRecorder()
	h.Commit(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

// --- Validate, GetStatusHistory, GetTruckState ---

func TestPlanBOL_GetStatusHistory_BadUUID_Returns400(t *testing.T) {
	h := newPlanBOLHandler(&stubPlanBOLSvc{}, &stubBOLRepo{}, &stubLogisticsClient{})
	req := withIDParam(httptest.NewRequest(http.MethodGet, "/", nil), "not-a-uuid")
	rec := httptest.NewRecorder()
	h.GetStatusHistory(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestPlanBOL_GetTruckState_BadUUID_Returns400(t *testing.T) {
	h := newPlanBOLHandler(&stubPlanBOLSvc{}, &stubBOLRepo{}, &stubLogisticsClient{})
	req := withIDParam(httptest.NewRequest(http.MethodGet, "/", nil), "not-a-uuid")
	rec := httptest.NewRecorder()
	h.GetTruckState(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestPlanBOL_Validate_BadUUID_Returns400(t *testing.T) {
	h := newPlanBOLHandler(&stubPlanBOLSvc{}, &stubBOLRepo{}, &stubLogisticsClient{})
	req := withIDParam(httptest.NewRequest(http.MethodPost, "/", nil), "not-a-uuid")
	rec := httptest.NewRecorder()
	h.Validate(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestPlanBOL_Validate_ServiceError_Returns500(t *testing.T) {
	h := newPlanBOLHandler(&stubPlanBOLSvc{err: errors.New("db error")}, &stubBOLRepo{}, &stubLogisticsClient{})
	req := withIDParam(httptest.NewRequest(http.MethodPost, "/", nil), uuid.New().String())
	rec := httptest.NewRecorder()
	h.Validate(rec, req)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestPlanBOL_Validate_Success_Returns200(t *testing.T) {
	plan := bolWithStatus(models.PlanBOLStatusPlanProgress)
	h := newPlanBOLHandler(&stubPlanBOLSvc{violations: []string{}}, &stubBOLRepo{bol: plan}, &stubLogisticsClient{})
	req := withIDParam(httptest.NewRequest(http.MethodPost, "/", nil), plan.ID.String())
	rec := httptest.NewRecorder()
	h.Validate(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestPlanBOL_GetStatusHistory_Returns200(t *testing.T) {
	plan := bolWithStatus(models.PlanBOLStatusDraft)
	h := newPlanBOLHandler(&stubPlanBOLSvc{}, &stubBOLRepo{bol: plan}, &stubLogisticsClient{})
	req := withIDParam(httptest.NewRequest(http.MethodGet, "/", nil), plan.ID.String())
	rec := httptest.NewRecorder()
	h.GetStatusHistory(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

// --- buildLineEntries ---

func TestBuildLineEntries_MixedStops(t *testing.T) {
	stops := []*models.PlanBOLStop{
		{
			LocationID:    "WH001",
			StopType:      models.StopTypeWarehouse,
			DeliveryItems: map[string]int{"SKU-A": 10, "SKU-B": 5},
		},
		{
			LocationID:    "ST0001",
			StopType:      models.StopTypeStore,
			DeliveryItems: map[string]int{"SKU-A": 3},
		},
	}
	entries := buildLineEntries(stops)
	// Warehouse items → positive qty; store items → negative qty.
	require.Len(t, entries, 3)
	skuAWarehouse, skuBWarehouse, skuAStore := 0, 0, 0
	for _, e := range entries {
		switch {
		case e.LocationID == "WH001" && e.SKUMarker == "SKU-A":
			skuAWarehouse = e.Quantity
		case e.LocationID == "WH001" && e.SKUMarker == "SKU-B":
			skuBWarehouse = e.Quantity
		case e.LocationID == "ST0001" && e.SKUMarker == "SKU-A":
			skuAStore = e.Quantity
		}
	}
	assert.Equal(t, 10, skuAWarehouse)
	assert.Equal(t, 5, skuBWarehouse)
	assert.Equal(t, -3, skuAStore)
}

func TestBuildLineEntries_EmptyStops(t *testing.T) {
	assert.Empty(t, buildLineEntries(nil))
	assert.Empty(t, buildLineEntries([]*models.PlanBOLStop{}))
}

func TestPlanBOL_GetTruckState_Returns200(t *testing.T) {
	h := newPlanBOLHandler(&stubPlanBOLSvc{}, &stubBOLRepo{}, &stubLogisticsClient{})
	req := withIDParam(httptest.NewRequest(http.MethodGet, "/", nil), uuid.New().String())

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", uuid.New().String())
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rec := httptest.NewRecorder()
	h.GetTruckState(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}
