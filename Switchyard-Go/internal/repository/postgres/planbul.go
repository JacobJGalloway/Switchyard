package postgres

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/JacobJGalloway/switchyard-go/internal/models"
)

type PlanBOLRepo struct{ db *pgxpool.Pool }

func NewPlanBOLRepo(db *pgxpool.Pool) *PlanBOLRepo { return &PlanBOLRepo{db: db} }

const bolCols = `id, driver_id, originating_wh_id, status, created_at, submitted_at, fulfilled_at, miles_driven, submitted_transaction_id`

func scanBOL(row interface{ Scan(...any) error }) (*models.PlanBOLRecord, error) {
	p := &models.PlanBOLRecord{}
	var status string
	err := row.Scan(&p.ID, &p.DriverID, &p.OriginatingWhID, &status,
		&p.CreatedAt, &p.SubmittedAt, &p.FulfilledAt, &p.MilesDriven, &p.SubmittedTransactionID)
	if err != nil {
		return nil, err
	}
	p.Status = models.PlanBOLStatus(status)
	return p, nil
}

func (r *PlanBOLRepo) Create(ctx context.Context, p *models.PlanBOLRecord) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO plan_bol_record (`+bolCols+`) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		p.ID, p.DriverID, p.OriginatingWhID, string(p.Status),
		p.CreatedAt, p.SubmittedAt, p.FulfilledAt, p.MilesDriven, p.SubmittedTransactionID)
	return err
}

func (r *PlanBOLRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.PlanBOLRecord, error) {
	return scanBOL(r.db.QueryRow(ctx,
		`SELECT `+bolCols+` FROM plan_bol_record WHERE id=$1`, id))
}

func (r *PlanBOLRepo) GetByStatus(ctx context.Context, status models.PlanBOLStatus) ([]*models.PlanBOLRecord, error) {
	rows, err := r.db.Query(ctx,
		`SELECT `+bolCols+` FROM plan_bol_record WHERE status=$1 ORDER BY created_at`, string(status))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*models.PlanBOLRecord
	for rows.Next() {
		p, err := scanBOL(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func (r *PlanBOLRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status models.PlanBOLStatus) error {
	_, err := r.db.Exec(ctx, `UPDATE plan_bol_record SET status=$2 WHERE id=$1`, id, string(status))
	return err
}

func (r *PlanBOLRepo) SetSubmittedTransactionID(ctx context.Context, id uuid.UUID, txID string) error {
	_, err := r.db.Exec(ctx, `UPDATE plan_bol_record SET submitted_transaction_id=$2 WHERE id=$1`, id, txID)
	return err
}

func (r *PlanBOLRepo) CreateStop(ctx context.Context, s *models.PlanBOLStop) error {
	var itemsJSON []byte
	if s.DeliveryItems != nil {
		itemsJSON, _ = json.Marshal(s.DeliveryItems)
	}
	_, err := r.db.Exec(ctx,
		`INSERT INTO plan_bol_stop (id, plan_bol_id, sequence, location_id, stop_type, delivery_items, is_processed, processed_at, miles_leg, driver_log_ref)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		s.ID, s.PlanBOLID, s.Sequence, s.LocationID, string(s.StopType),
		itemsJSON, s.IsProcessed, s.ProcessedAt, s.MilesLeg, s.DriverLogRef)
	return err
}

func (r *PlanBOLRepo) GetStopByID(ctx context.Context, stopID uuid.UUID) (*models.PlanBOLStop, error) {
	return scanStop(r.db.QueryRow(ctx,
		`SELECT id, plan_bol_id, sequence, location_id, stop_type, delivery_items, is_processed, processed_at, miles_leg, driver_log_ref
		 FROM plan_bol_stop WHERE id=$1`, stopID))
}

func (r *PlanBOLRepo) GetStops(ctx context.Context, planBOLID uuid.UUID) ([]*models.PlanBOLStop, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, plan_bol_id, sequence, location_id, stop_type, delivery_items, is_processed, processed_at, miles_leg, driver_log_ref
		 FROM plan_bol_stop WHERE plan_bol_id=$1 ORDER BY sequence`, planBOLID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*models.PlanBOLStop
	for rows.Next() {
		s, err := scanStop(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

func (r *PlanBOLRepo) MarkStopProcessed(ctx context.Context, stopID uuid.UUID, processedAt time.Time, milesLeg *float64) error {
	_, err := r.db.Exec(ctx,
		`UPDATE plan_bol_stop SET is_processed=true, processed_at=$2, miles_leg=$3 WHERE id=$1`,
		stopID, processedAt, milesLeg)
	return err
}

func (r *PlanBOLRepo) SetMilesDriven(ctx context.Context, id uuid.UUID, miles float64) error {
	_, err := r.db.Exec(ctx, `UPDATE plan_bol_record SET miles_driven=$2 WHERE id=$1`, id, miles)
	return err
}

func (r *PlanBOLRepo) CreateSnapshot(ctx context.Context, s *models.TruckInventorySnapshot) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO truck_inventory_snapshot (id, plan_bol_id, plan_bol_stop_id, sku_id, quantity_loaded, quantity_remaining, snapshot_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		s.ID, s.PlanBOLID, s.PlanBOLStopID, s.SKUID, s.QuantityLoaded, s.QuantityRemaining, s.SnapshotAt)
	return err
}

func (r *PlanBOLRepo) GetSnapshots(ctx context.Context, planBOLID uuid.UUID) ([]*models.TruckInventorySnapshot, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, plan_bol_id, plan_bol_stop_id, sku_id, quantity_loaded, quantity_remaining, snapshot_at
		 FROM truck_inventory_snapshot WHERE plan_bol_id=$1 ORDER BY snapshot_at`, planBOLID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*models.TruckInventorySnapshot
	for rows.Next() {
		s := &models.TruckInventorySnapshot{}
		if err := rows.Scan(&s.ID, &s.PlanBOLID, &s.PlanBOLStopID, &s.SKUID, &s.QuantityLoaded, &s.QuantityRemaining, &s.SnapshotAt); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

func (r *PlanBOLRepo) GetStatusHistory(ctx context.Context, planBOLID uuid.UUID) ([]*models.BOLStatusHistory, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, plan_bol_id, from_status, to_status, changed_at
		 FROM plan_bol_status_history WHERE plan_bol_id=$1 ORDER BY changed_at`, planBOLID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*models.BOLStatusHistory
	for rows.Next() {
		h := &models.BOLStatusHistory{}
		var fromStatus *string
		var toStatus string
		if err := rows.Scan(&h.ID, &h.PlanBOLID, &fromStatus, &toStatus, &h.ChangedAt); err != nil {
			return nil, err
		}
		if fromStatus != nil {
			s := models.PlanBOLStatus(*fromStatus)
			h.FromStatus = &s
		}
		h.ToStatus = models.PlanBOLStatus(toStatus)
		out = append(out, h)
	}
	return out, rows.Err()
}

func scanStop(row interface{ Scan(...any) error }) (*models.PlanBOLStop, error) {
	s := &models.PlanBOLStop{}
	var stopType string
	var itemsRaw []byte
	err := row.Scan(&s.ID, &s.PlanBOLID, &s.Sequence, &s.LocationID, &stopType,
		&itemsRaw, &s.IsProcessed, &s.ProcessedAt, &s.MilesLeg, &s.DriverLogRef)
	if err != nil {
		return nil, err
	}
	s.StopType = models.StopType(stopType)
	if len(itemsRaw) > 0 {
		_ = json.Unmarshal(itemsRaw, &s.DeliveryItems)
	}
	return s, nil
}
