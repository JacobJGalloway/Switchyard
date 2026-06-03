package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/JacobJGalloway/switchyard-go/internal/models"
)

type AnalyticsRepo struct{ db *pgxpool.Pool }

func NewAnalyticsRepo(db *pgxpool.Pool) *AnalyticsRepo { return &AnalyticsRepo{db: db} }

func (r *AnalyticsRepo) BOLsByStatus(ctx context.Context) ([]models.BOLStatusCount, error) {
	rows, err := r.db.Query(ctx,
		`SELECT status, COUNT(*) FROM plan_bol_record GROUP BY status ORDER BY status`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.BOLStatusCount
	for rows.Next() {
		var c models.BOLStatusCount
		if err := rows.Scan(&c.Status, &c.Count); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	if out == nil {
		out = []models.BOLStatusCount{}
	}
	return out, rows.Err()
}

func (r *AnalyticsRepo) StopCompletionRate(ctx context.Context) (float64, error) {
	var total, processed int
	err := r.db.QueryRow(ctx,
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

func (r *AnalyticsRepo) FulfilledInWindow(ctx context.Context, since time.Time) (int, error) {
	var count int
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM plan_bol_record WHERE status='fulfilled' AND fulfilled_at >= $1`, since).
		Scan(&count)
	return count, err
}

func (r *AnalyticsRepo) OperatingCostByBOL(ctx context.Context) ([]models.BOLOperatingCost, error) {
	rows, err := r.db.Query(ctx, `
		SELECT
			p.id,
			p.driver_id,
			p.originating_wh_id,
			COALESCE(p.miles_driven, 0),
			COALESCE(a.base_rate_per_mile, 3.20),
			COALESCE(SUM(b.tow_cost), 0),
			COALESCE(p.miles_driven, 0) * COALESCE(a.base_rate_per_mile, 3.20) + COALESCE(SUM(b.tow_cost), 0)
		FROM plan_bol_record p
		LEFT JOIN driver_bol_assignment a ON a.plan_bol_id = p.id
		LEFT JOIN breakdown_record b ON b.driver_id = a.driver_id AND b.resolved_at IS NOT NULL
			AND b.reported_at >= COALESCE(a.departed_at, a.assigned_at)
			AND b.reported_at <= COALESCE(a.fulfilled_at, now())
		WHERE p.status = 'fulfilled'
		GROUP BY p.id, p.driver_id, p.originating_wh_id, p.miles_driven, a.base_rate_per_mile
		ORDER BY p.originating_wh_id, p.driver_id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.BOLOperatingCost
	for rows.Next() {
		var c models.BOLOperatingCost
		if err := rows.Scan(&c.PlanBOLID, &c.DriverID, &c.WarehouseID,
			&c.MilesDriven, &c.BaseRate, &c.TowCost, &c.TotalCost); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	if out == nil {
		out = []models.BOLOperatingCost{}
	}
	return out, rows.Err()
}

func (r *AnalyticsRepo) OperatingCostByDriver(ctx context.Context) ([]models.DriverOperatingCost, error) {
	rows, err := r.db.Query(ctx, `
		SELECT
			p.driver_id,
			COALESCE(SUM(p.miles_driven), 0),
			COALESCE(SUM(p.miles_driven * COALESCE(a.base_rate_per_mile, 3.20)), 0) + COALESCE(SUM(b.tow_cost), 0),
			COUNT(DISTINCT p.id)
		FROM plan_bol_record p
		LEFT JOIN driver_bol_assignment a ON a.plan_bol_id = p.id
		LEFT JOIN breakdown_record b ON b.driver_id = a.driver_id AND b.resolved_at IS NOT NULL
			AND b.reported_at >= COALESCE(a.departed_at, a.assigned_at)
			AND b.reported_at <= COALESCE(a.fulfilled_at, now())
		WHERE p.status = 'fulfilled'
		GROUP BY p.driver_id
		ORDER BY p.driver_id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.DriverOperatingCost
	for rows.Next() {
		var c models.DriverOperatingCost
		if err := rows.Scan(&c.DriverID, &c.TotalMiles, &c.TotalCost, &c.BOLCount); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	if out == nil {
		out = []models.DriverOperatingCost{}
	}
	return out, rows.Err()
}

func (r *AnalyticsRepo) OperatingCostByWarehouse(ctx context.Context) ([]models.WarehouseOperatingCost, error) {
	rows, err := r.db.Query(ctx, `
		SELECT
			p.originating_wh_id,
			COALESCE(SUM(p.miles_driven), 0),
			COALESCE(SUM(p.miles_driven * COALESCE(a.base_rate_per_mile, 3.20)), 0) + COALESCE(SUM(b.tow_cost), 0),
			COUNT(DISTINCT p.id)
		FROM plan_bol_record p
		LEFT JOIN driver_bol_assignment a ON a.plan_bol_id = p.id
		LEFT JOIN breakdown_record b ON b.driver_id = a.driver_id AND b.resolved_at IS NOT NULL
			AND b.reported_at >= COALESCE(a.departed_at, a.assigned_at)
			AND b.reported_at <= COALESCE(a.fulfilled_at, now())
		WHERE p.status = 'fulfilled'
		GROUP BY p.originating_wh_id
		ORDER BY p.originating_wh_id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.WarehouseOperatingCost
	for rows.Next() {
		var c models.WarehouseOperatingCost
		if err := rows.Scan(&c.WarehouseID, &c.TotalMiles, &c.TotalCost, &c.BOLCount); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	if out == nil {
		out = []models.WarehouseOperatingCost{}
	}
	return out, rows.Err()
}
