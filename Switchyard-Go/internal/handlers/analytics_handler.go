package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// AnalyticsHandler serves aggregate summary data for the dispatch whiteboard charts.
// All queries are read-only against existing tables — no dedicated analytics schema.
type AnalyticsHandler struct {
	db *pgxpool.Pool
}

func NewAnalyticsHandler(db *pgxpool.Pool) *AnalyticsHandler {
	return &AnalyticsHandler{db: db}
}

type bolStatusCount struct {
	Status string `json:"status"`
	Count  int    `json:"count"`
}

type analyticsSummary struct {
	BOLsByStatus      []bolStatusCount `json:"bols_by_status"`
	StopCompletionPct float64          `json:"stop_completion_pct"`
	FulfilledLast7d   int              `json:"fulfilled_last_7d"`
	FulfilledLast30d  int              `json:"fulfilled_last_30d"`
}

// GetSummary handles GET /api/analytics/summary
func (h *AnalyticsHandler) GetSummary(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	counts, err := h.bolsByStatus(ctx)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to query BOL counts")
		return
	}

	completionPct, err := h.stopCompletionRate(ctx)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to query stop completion rate")
		return
	}

	last7d, err := h.fulfilledSince(ctx, 7*24*time.Hour)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to query 7-day throughput")
		return
	}

	last30d, err := h.fulfilledSince(ctx, 30*24*time.Hour)
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

func (h *AnalyticsHandler) bolsByStatus(ctx context.Context) ([]bolStatusCount, error) {
	rows, err := h.db.Query(ctx,
		`SELECT status, COUNT(*) FROM plan_bol_record GROUP BY status ORDER BY status`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []bolStatusCount
	for rows.Next() {
		var c bolStatusCount
		if err := rows.Scan(&c.Status, &c.Count); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	if out == nil {
		out = []bolStatusCount{}
	}
	return out, rows.Err()
}

func (h *AnalyticsHandler) stopCompletionRate(ctx context.Context) (float64, error) {
	var total, processed int
	err := h.db.QueryRow(ctx,
		`SELECT COUNT(*), COUNT(*) FILTER (WHERE is_processed) FROM plan_bol_stop`).
		Scan(&total, &processed)
	if err != nil {
		return 0, err
	}
	if total == 0 {
		return 0, nil
	}
	return float64(processed) / float64(total) * 100, nil
}

func (h *AnalyticsHandler) fulfilledSince(ctx context.Context, window time.Duration) (int, error) {
	since := time.Now().Add(-window)
	var count int
	err := h.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM plan_bol_record WHERE status='fulfilled' AND fulfilled_at >= $1`, since).
		Scan(&count)
	return count, err
}
