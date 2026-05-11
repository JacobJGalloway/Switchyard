package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/JacobJGalloway/switchyard-go/internal/models"
)

type stubInvoiceRepo struct {
	invoice  *models.InternalInvoice
	invoices []*models.InternalInvoice
	err      error
}

func (r *stubInvoiceRepo) CreateConfirmation(_ context.Context, _ *models.DeliveryConfirmation) error {
	return nil
}
func (r *stubInvoiceRepo) GetConfirmation(_ context.Context, _ uuid.UUID) (*models.DeliveryConfirmation, error) {
	return nil, nil
}
func (r *stubInvoiceRepo) GetConfirmationsByBOL(_ context.Context, _ uuid.UUID) ([]*models.DeliveryConfirmation, error) {
	return nil, nil
}
func (r *stubInvoiceRepo) SetConfirmationInvoice(_ context.Context, _, _ uuid.UUID) error { return nil }
func (r *stubInvoiceRepo) CreateInvoice(_ context.Context, _ *models.InternalInvoice) error {
	return nil
}
func (r *stubInvoiceRepo) GetInvoice(_ context.Context, _ uuid.UUID) (*models.InternalInvoice, error) {
	if r.err != nil {
		return nil, r.err
	}
	return r.invoice, nil
}
func (r *stubInvoiceRepo) GetInvoicesByStore(_ context.Context, _ string) ([]*models.InternalInvoice, error) {
	if r.err != nil {
		return nil, r.err
	}
	return r.invoices, nil
}

// --- GetInvoice ---

func TestInvoice_GetInvoice_BadUUID_Returns400(t *testing.T) {
	h := NewInvoiceHandler(&stubInvoiceRepo{})
	req := withIDParam(httptest.NewRequest(http.MethodGet, "/", nil), "not-a-uuid")
	rec := httptest.NewRecorder()
	h.GetInvoice(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestInvoice_GetInvoice_NotFound_Returns404(t *testing.T) {
	h := NewInvoiceHandler(&stubInvoiceRepo{err: errNotFound})
	req := withIDParam(httptest.NewRequest(http.MethodGet, "/", nil), uuid.New().String())
	rec := httptest.NewRecorder()
	h.GetInvoice(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestInvoice_GetInvoice_Success_Returns200(t *testing.T) {
	inv := &models.InternalInvoice{ID: uuid.New(), StoreID: "ST0001"}
	h := NewInvoiceHandler(&stubInvoiceRepo{invoice: inv})
	req := withIDParam(httptest.NewRequest(http.MethodGet, "/", nil), inv.ID.String())
	rec := httptest.NewRecorder()
	h.GetInvoice(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

// --- GetInvoicesByStore ---

func TestInvoice_GetInvoicesByStore_EmptyList_Returns200(t *testing.T) {
	h := NewInvoiceHandler(&stubInvoiceRepo{})
	req := httptest.NewRequest(http.MethodGet, "/api/invoice/store/ST0001", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("storeId", "ST0001")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	rec := httptest.NewRecorder()
	h.GetInvoicesByStore(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestInvoice_GetInvoicesByStore_RepoError_Returns500(t *testing.T) {
	h := NewInvoiceHandler(&stubInvoiceRepo{err: errNotFound})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("storeId", "ST0001")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	rec := httptest.NewRecorder()
	h.GetInvoicesByStore(rec, req)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestInvoice_GetInvoicesByStore_WithResults_Returns200(t *testing.T) {
	invs := []*models.InternalInvoice{{ID: uuid.New(), StoreID: "ST0002"}}
	h := NewInvoiceHandler(&stubInvoiceRepo{invoices: invs})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("storeId", "ST0002")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	rec := httptest.NewRecorder()
	h.GetInvoicesByStore(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}
