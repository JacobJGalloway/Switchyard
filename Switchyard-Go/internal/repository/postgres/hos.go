package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/JacobJGalloway/switchyard-go/internal/models"
)

type HOSRepo struct{ db *pgxpool.Pool }

func NewHOSRepo(db *pgxpool.Pool) *HOSRepo { return &HOSRepo{db: db} }

func (r *HOSRepo) CreateLimit(ctx context.Context, l *models.HOSLimit) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO hos_limit (
			id, state_code, daily_driving_limit_hours, daily_period_hours, rest_period_hours,
			weekly_limit_hours, weekly_period_days, weekly_reset_hours, sleeper_cab_min_hours,
			short_haul_radius_miles, adverse_weather_extension_hours, break_required_after_hours,
			sleeper_split_allowed, sleeper_split_options, cycle_label, effective_from, notes
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)`,
		l.ID, l.StateCode, l.DailyDrivingLimitHours, l.DailyPeriodHours, l.RestPeriodHours,
		l.WeeklyLimitHours, l.WeeklyPeriodDays, l.WeeklyResetHours, l.SleeperCabMinHours,
		l.ShortHaulRadiusMiles, l.AdverseWeatherExtensionHours, l.BreakRequiredAfterHours,
		l.SleeperSplitAllowed, l.SleeperSplitOptions, l.CycleLabel, l.EffectiveFrom, l.Notes)
	return err
}

func (r *HOSRepo) GetLimitByStateAndCycle(ctx context.Context, stateCode, cycleLabel string) (*models.HOSLimit, error) {
	l := &models.HOSLimit{}
	err := r.db.QueryRow(ctx, `
		SELECT id, state_code, daily_driving_limit_hours, daily_period_hours, rest_period_hours,
		       weekly_limit_hours, weekly_period_days, weekly_reset_hours, sleeper_cab_min_hours,
		       short_haul_radius_miles, adverse_weather_extension_hours, break_required_after_hours,
		       sleeper_split_allowed, sleeper_split_options, cycle_label, effective_from, notes
		FROM hos_limit
		WHERE state_code=$1 AND cycle_label=$2
		ORDER BY effective_from DESC LIMIT 1`, stateCode, cycleLabel).
		Scan(
			&l.ID, &l.StateCode, &l.DailyDrivingLimitHours, &l.DailyPeriodHours, &l.RestPeriodHours,
			&l.WeeklyLimitHours, &l.WeeklyPeriodDays, &l.WeeklyResetHours, &l.SleeperCabMinHours,
			&l.ShortHaulRadiusMiles, &l.AdverseWeatherExtensionHours, &l.BreakRequiredAfterHours,
			&l.SleeperSplitAllowed, &l.SleeperSplitOptions, &l.CycleLabel, &l.EffectiveFrom, &l.Notes)
	if err != nil {
		return nil, err
	}
	return l, nil
}

func (r *HOSRepo) CreateWindow(ctx context.Context, w *models.HOSWindow) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO hos_window (id, driver_id, window_start, daily_hours_used, weekly_hours_used,
		    last_activity_at, break_30_taken, break_30_at, mandated_stop_at, eld_stop_ref)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		w.ID, w.DriverID, w.WindowStart, w.DailyHoursUsed, w.WeeklyHoursUsed,
		w.LastActivityAt, w.Break30Taken, w.Break30At, w.MandatedStopAt, w.ELDStopRef)
	return err
}

func (r *HOSRepo) GetWindowByDriver(ctx context.Context, driverID uuid.UUID) (*models.HOSWindow, error) {
	w := &models.HOSWindow{}
	err := r.db.QueryRow(ctx, `
		SELECT id, driver_id, window_start, daily_hours_used, weekly_hours_used,
		       last_activity_at, break_30_taken, break_30_at, mandated_stop_at, eld_stop_ref
		FROM hos_window WHERE driver_id=$1
		ORDER BY window_start DESC LIMIT 1`, driverID).
		Scan(&w.ID, &w.DriverID, &w.WindowStart, &w.DailyHoursUsed, &w.WeeklyHoursUsed,
			&w.LastActivityAt, &w.Break30Taken, &w.Break30At, &w.MandatedStopAt, &w.ELDStopRef)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (r *HOSRepo) UpdateWindow(ctx context.Context, w *models.HOSWindow) error {
	_, err := r.db.Exec(ctx, `
		UPDATE hos_window SET
		    daily_hours_used=$2, weekly_hours_used=$3, last_activity_at=$4,
		    break_30_taken=$5, break_30_at=$6, mandated_stop_at=$7, eld_stop_ref=$8
		WHERE id=$1`,
		w.ID, w.DailyHoursUsed, w.WeeklyHoursUsed, w.LastActivityAt,
		w.Break30Taken, w.Break30At, w.MandatedStopAt, w.ELDStopRef)
	return err
}
