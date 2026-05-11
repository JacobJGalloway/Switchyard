package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/JacobJGalloway/switchyard-go/internal/models"
)

// --- stub ---

type stubPairingRepo struct {
	pairings []*models.PlanBOLPairing
	created  *models.PlanBOLPairing
}

func (r *stubPairingRepo) GetEligible(_ context.Context, _ string, _ time.Time) ([]*models.PlanBOLPairing, error) {
	return r.pairings, nil
}
func (r *stubPairingRepo) Create(_ context.Context, p *models.PlanBOLPairing) error {
	r.created = p
	return nil
}
func (r *stubPairingRepo) GetByID(_ context.Context, _ uuid.UUID) (*models.PlanBOLPairing, error) {
	return nil, nil
}
func (r *stubPairingRepo) GetByActiveBOL(_ context.Context, _ uuid.UUID) (*models.PlanBOLPairing, error) {
	return nil, nil
}
func (r *stubPairingRepo) UpdateStatus(_ context.Context, _ uuid.UUID, _ models.PairingStatus) error {
	return nil
}

// --- GetEligible ---

func TestDeadheadGetEligible_MissingLocation_Returns400(t *testing.T) {
	h := NewDeadheadHandler(&stubPairingRepo{}, 4.0)
	q := url.Values{"estimated_completion": {time.Now().Format(time.RFC3339)}}
	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	h.GetEligible(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestDeadheadGetEligible_MissingEstimatedCompletion_Returns400(t *testing.T) {
	h := NewDeadheadHandler(&stubPairingRepo{}, 4.0)
	q := url.Values{"location": {"wh-1"}}
	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	h.GetEligible(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestDeadheadGetEligible_InvalidTimeFormat_Returns400(t *testing.T) {
	h := NewDeadheadHandler(&stubPairingRepo{}, 4.0)
	q := url.Values{"location": {"wh-1"}, "estimated_completion": {"not-a-time"}}
	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	h.GetEligible(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestDeadheadGetEligible_ValidParams_Returns200(t *testing.T) {
	h := NewDeadheadHandler(&stubPairingRepo{pairings: []*models.PlanBOLPairing{}}, 4.0)
	q := url.Values{
		"location":             {"wh-1"},
		"estimated_completion": {time.Now().Add(2 * time.Hour).Format(time.RFC3339)},
	}
	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	h.GetEligible(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

// --- Pair (4-hour window enforcement — architecture §4.2 hard constraint) ---

func TestDeadheadPair_WindowAlreadyClosed_Returns409(t *testing.T) {
	h := NewDeadheadHandler(&stubPairingRepo{}, 4.0)
	// Fulfillment in 1 hour — less than the 4-hour required lead time.
	body := map[string]any{
		"active_bol_id":            uuid.New().String(),
		"deadhead_bol_id":          uuid.New().String(),
		"estimated_fulfillment_at": time.Now().Add(1 * time.Hour),
		"origin_warehouse":         "wh-1",
	}
	req := httptest.NewRequest(http.MethodPost, "/api/deadhead/pair", postBody(t, body))
	rec := httptest.NewRecorder()
	h.Pair(rec, req)
	assert.Equal(t, http.StatusConflict, rec.Code)
}

func TestDeadheadPair_WindowOpen_Returns201(t *testing.T) {
	repo := &stubPairingRepo{}
	h := NewDeadheadHandler(repo, 4.0)
	// Fulfillment in 6 hours — well within the 4-hour lead-time window.
	body := map[string]any{
		"active_bol_id":            uuid.New().String(),
		"deadhead_bol_id":          uuid.New().String(),
		"estimated_fulfillment_at": time.Now().Add(6 * time.Hour),
		"origin_warehouse":         "wh-1",
	}
	req := httptest.NewRequest(http.MethodPost, "/api/deadhead/pair", postBody(t, body))
	rec := httptest.NewRecorder()
	h.Pair(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)
	assert.NotNil(t, repo.created, "pairing record must be persisted")
	assert.Equal(t, models.PairingStatusProposed, repo.created.Status)
}

func TestDeadheadPair_ExactlyAtWindow_Returns409(t *testing.T) {
	// Fulfillment in exactly 4 hours: earliestValidAt == now, window is closed (< not <=).
	h := NewDeadheadHandler(&stubPairingRepo{}, 4.0)
	body := map[string]any{
		"active_bol_id":            uuid.New().String(),
		"deadhead_bol_id":          uuid.New().String(),
		"estimated_fulfillment_at": time.Now().Add(4 * time.Hour),
		"origin_warehouse":         "wh-1",
	}
	req := httptest.NewRequest(http.MethodPost, "/api/deadhead/pair", postBody(t, body))
	rec := httptest.NewRecorder()
	h.Pair(rec, req)
	assert.Equal(t, http.StatusConflict, rec.Code)
}

func TestDeadheadPair_MissingOriginWarehouse_Returns400(t *testing.T) {
	h := NewDeadheadHandler(&stubPairingRepo{}, 4.0)
	body := map[string]any{
		"active_bol_id":            uuid.New().String(),
		"deadhead_bol_id":          uuid.New().String(),
		"estimated_fulfillment_at": time.Now().Add(6 * time.Hour),
	}
	req := httptest.NewRequest(http.MethodPost, "/api/deadhead/pair", postBody(t, body))
	rec := httptest.NewRecorder()
	h.Pair(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestDeadheadPair_ZeroEstimatedFulfillment_Returns400(t *testing.T) {
	h := NewDeadheadHandler(&stubPairingRepo{}, 4.0)
	body := map[string]any{
		"active_bol_id":   uuid.New().String(),
		"deadhead_bol_id": uuid.New().String(),
		"origin_warehouse": "wh-1",
		// estimated_fulfillment_at intentionally omitted → zero time
	}
	req := httptest.NewRequest(http.MethodPost, "/api/deadhead/pair", postBody(t, body))
	rec := httptest.NewRecorder()
	h.Pair(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestDeadheadPair_InvalidBOLID_Returns400(t *testing.T) {
	h := NewDeadheadHandler(&stubPairingRepo{}, 4.0)
	body := map[string]any{
		"active_bol_id":            "not-a-uuid",
		"deadhead_bol_id":          uuid.New().String(),
		"estimated_fulfillment_at": time.Now().Add(6 * time.Hour),
		"origin_warehouse":         "wh-1",
	}
	req := httptest.NewRequest(http.MethodPost, "/api/deadhead/pair", postBody(t, body))
	rec := httptest.NewRecorder()
	h.Pair(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// --- Cancel ---

func TestDeadheadCancel_ValidID_Returns204(t *testing.T) {
	h := NewDeadheadHandler(&stubPairingRepo{}, 4.0)
	rctx := withPairingIDParam(httptest.NewRequest(http.MethodDelete, "/", nil), uuid.New().String())
	rec := httptest.NewRecorder()
	h.Cancel(rec, rctx)
	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestDeadheadCancel_InvalidID_Returns400(t *testing.T) {
	h := NewDeadheadHandler(&stubPairingRepo{}, 4.0)
	rctx := withPairingIDParam(httptest.NewRequest(http.MethodDelete, "/", nil), "bad-id")
	rec := httptest.NewRecorder()
	h.Cancel(rec, rctx)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func withPairingIDParam(r *http.Request, id string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("pairingId", id)
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
}
