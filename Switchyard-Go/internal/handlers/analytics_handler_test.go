package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/JacobJGalloway/switchyard-go/internal/models"
)

// stubAnalyticsQuerier implements repository.AnalyticsQuerier for testing.
type stubAnalyticsQuerier struct {
	statusErr        error
	completionErr    error
	windowErr        error
	costErr          error
	costDriverErr    error
	costWarehouseErr error
}

func (s *stubAnalyticsQuerier) BOLsByStatus(_ context.Context) ([]models.BOLStatusCount, error) {
	if s.statusErr != nil {
		return nil, s.statusErr
	}
	return []models.BOLStatusCount{{Status: "fulfilled", Count: 3}}, nil
}

func (s *stubAnalyticsQuerier) StopCompletionRate(_ context.Context) (float64, error) {
	if s.completionErr != nil {
		return 0, s.completionErr
	}
	return 75.0, nil
}

func (s *stubAnalyticsQuerier) FulfilledInWindow(_ context.Context, _ time.Time) (int, error) {
	if s.windowErr != nil {
		return 0, s.windowErr
	}
	return 5, nil
}

func (s *stubAnalyticsQuerier) OperatingCostByBOL(_ context.Context) ([]models.BOLOperatingCost, error) {
	if s.costErr != nil {
		return nil, s.costErr
	}
	return []models.BOLOperatingCost{}, nil
}

func (s *stubAnalyticsQuerier) OperatingCostByDriver(_ context.Context) ([]models.DriverOperatingCost, error) {
	if s.costDriverErr != nil {
		return nil, s.costDriverErr
	}
	if s.costErr != nil {
		return nil, s.costErr
	}
	return []models.DriverOperatingCost{}, nil
}

func (s *stubAnalyticsQuerier) OperatingCostByWarehouse(_ context.Context) ([]models.WarehouseOperatingCost, error) {
	if s.costWarehouseErr != nil {
		return nil, s.costWarehouseErr
	}
	if s.costErr != nil {
		return nil, s.costErr
	}
	return []models.WarehouseOperatingCost{}, nil
}

// --- GetSummary ---

func TestAnalytics_GetSummary_Returns200(t *testing.T) {
	h := NewAnalyticsHandler(&stubAnalyticsQuerier{})
	req := httptest.NewRequest(http.MethodGet, "/api/analytics/summary", nil)
	rec := httptest.NewRecorder()
	h.GetSummary(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	var body map[string]any
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	assert.Contains(t, body, "bols_by_status")
	assert.Contains(t, body, "stop_completion_pct")
	assert.Contains(t, body, "fulfilled_last_7d")
	assert.Contains(t, body, "fulfilled_last_30d")
}

func TestAnalytics_GetSummary_BOLStatusError_Returns500(t *testing.T) {
	h := NewAnalyticsHandler(&stubAnalyticsQuerier{statusErr: errors.New("db error")})
	req := httptest.NewRequest(http.MethodGet, "/api/analytics/summary", nil)
	rec := httptest.NewRecorder()
	h.GetSummary(rec, req)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestAnalytics_GetSummary_CompletionRateError_Returns500(t *testing.T) {
	h := NewAnalyticsHandler(&stubAnalyticsQuerier{completionErr: errors.New("db error")})
	req := httptest.NewRequest(http.MethodGet, "/api/analytics/summary", nil)
	rec := httptest.NewRecorder()
	h.GetSummary(rec, req)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestAnalytics_GetSummary_WindowError_Returns500(t *testing.T) {
	h := NewAnalyticsHandler(&stubAnalyticsQuerier{windowErr: errors.New("db error")})
	req := httptest.NewRequest(http.MethodGet, "/api/analytics/summary", nil)
	rec := httptest.NewRecorder()
	h.GetSummary(rec, req)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

// --- GetOperatingCost ---

func TestAnalytics_GetOperatingCost_Returns200(t *testing.T) {
	h := NewAnalyticsHandler(&stubAnalyticsQuerier{})
	req := httptest.NewRequest(http.MethodGet, "/api/analytics/operating-cost", nil)
	rec := httptest.NewRecorder()
	h.GetOperatingCost(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	var body map[string]any
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	assert.Contains(t, body, "by_bol")
	assert.Contains(t, body, "by_driver")
	assert.Contains(t, body, "by_warehouse")
}

func TestAnalytics_GetOperatingCost_BOLError_Returns500(t *testing.T) {
	h := NewAnalyticsHandler(&stubAnalyticsQuerier{costErr: errors.New("db error")})
	req := httptest.NewRequest(http.MethodGet, "/api/analytics/operating-cost", nil)
	rec := httptest.NewRecorder()
	h.GetOperatingCost(rec, req)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestAnalytics_GetOperatingCost_DriverError_Returns500(t *testing.T) {
	h := NewAnalyticsHandler(&stubAnalyticsQuerier{costDriverErr: errors.New("db error")})
	req := httptest.NewRequest(http.MethodGet, "/api/analytics/operating-cost", nil)
	rec := httptest.NewRecorder()
	h.GetOperatingCost(rec, req)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestAnalytics_GetOperatingCost_WarehouseError_Returns500(t *testing.T) {
	h := NewAnalyticsHandler(&stubAnalyticsQuerier{costWarehouseErr: errors.New("db error")})
	req := httptest.NewRequest(http.MethodGet, "/api/analytics/operating-cost", nil)
	rec := httptest.NewRecorder()
	h.GetOperatingCost(rec, req)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
