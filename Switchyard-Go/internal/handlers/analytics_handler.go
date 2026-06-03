package handlers

import (
	"net/http"
	"time"

	"github.com/JacobJGalloway/switchyard-go/internal/repository"
)

// AnalyticsHandler serves aggregate summary and operating cost data for the dispatch whiteboard.
type AnalyticsHandler struct {
	querier repository.AnalyticsQuerier
}

func NewAnalyticsHandler(q repository.AnalyticsQuerier) *AnalyticsHandler {
	return &AnalyticsHandler{querier: q}
}

type analyticsSummary struct {
	BOLsByStatus      any     `json:"bols_by_status"`
	StopCompletionPct float64 `json:"stop_completion_pct"`
	FulfilledLast7d   int     `json:"fulfilled_last_7d"`
	FulfilledLast30d  int     `json:"fulfilled_last_30d"`
}

type operatingCostResponse struct {
	ByBOL      any `json:"by_bol"`
	ByDriver   any `json:"by_driver"`
	ByWarehouse any `json:"by_warehouse"`
}

// GetSummary handles GET /api/analytics/summary
func (h *AnalyticsHandler) GetSummary(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	counts, err := h.querier.BOLsByStatus(ctx)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to query BOL counts")
		return
	}

	completionPct, err := h.querier.StopCompletionRate(ctx)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to query stop completion rate")
		return
	}

	last7d, err := h.querier.FulfilledInWindow(ctx, time.Now().Add(-7*24*time.Hour))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to query 7-day throughput")
		return
	}

	last30d, err := h.querier.FulfilledInWindow(ctx, time.Now().Add(-30*24*time.Hour))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to query 30-day throughput")
		return
	}

	writeJSON(w, http.StatusOK, analyticsSummary{
		BOLsByStatus:      counts,
		StopCompletionPct: completionPct,
		FulfilledLast7d:   last7d,
		FulfilledLast30d:  last30d,
	})
}

// GetOperatingCost handles GET /api/analytics/operating-cost
func (h *AnalyticsHandler) GetOperatingCost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	byBOL, err := h.querier.OperatingCostByBOL(ctx)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to query BOL operating costs")
		return
	}

	byDriver, err := h.querier.OperatingCostByDriver(ctx)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to query driver operating costs")
		return
	}

	byWarehouse, err := h.querier.OperatingCostByWarehouse(ctx)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to query warehouse operating costs")
		return
	}

	writeJSON(w, http.StatusOK, operatingCostResponse{
		ByBOL:       byBOL,
		ByDriver:    byDriver,
		ByWarehouse: byWarehouse,
	})
}
