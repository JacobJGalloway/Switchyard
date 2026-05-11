package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/JacobJGalloway/switchyard-go/internal/models"
)

type PairingRepo struct{ db *pgxpool.Pool }

func NewPairingRepo(db *pgxpool.Pool) *PairingRepo { return &PairingRepo{db: db} }

const pairingCols = `id, active_bol_id, deadhead_bol_id, paired_at, earliest_valid_at, origin_warehouse, status`

func scanPairing(row interface{ Scan(...any) error }) (*models.PlanBOLPairing, error) {
	p := &models.PlanBOLPairing{}
	var status string
	err := row.Scan(&p.ID, &p.ActiveBOLID, &p.DeadheadBOLID, &p.PairedAt, &p.EarliestValidAt, &p.OriginWarehouse, &status)
	if err != nil {
		return nil, err
	}
	p.Status = models.PairingStatus(status)
	return p, nil
}

func (r *PairingRepo) Create(ctx context.Context, p *models.PlanBOLPairing) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO plan_bol_pairing (`+pairingCols+`)
		 VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		p.ID, p.ActiveBOLID, p.DeadheadBOLID, p.PairedAt, p.EarliestValidAt, p.OriginWarehouse, string(p.Status))
	return err
}

func (r *PairingRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.PlanBOLPairing, error) {
	return scanPairing(r.db.QueryRow(ctx,
		`SELECT `+pairingCols+` FROM plan_bol_pairing WHERE id=$1`, id))
}

func (r *PairingRepo) GetByActiveBOL(ctx context.Context, activeBOLID uuid.UUID) (*models.PlanBOLPairing, error) {
	return scanPairing(r.db.QueryRow(ctx,
		`SELECT `+pairingCols+` FROM plan_bol_pairing WHERE active_bol_id=$1
		 ORDER BY paired_at DESC LIMIT 1`, activeBOLID))
}

// GetEligible returns proposed pairings whose origin warehouse matches the given location
// and whose earliest_valid_at is within the estimated completion window.
// Geographic nearest-warehouse resolution is deferred — location is matched exactly for v1.1.
func (r *PairingRepo) GetEligible(ctx context.Context, location string, estimatedCompletion time.Time) ([]*models.PlanBOLPairing, error) {
	rows, err := r.db.Query(ctx,
		`SELECT `+pairingCols+` FROM plan_bol_pairing
		 WHERE origin_warehouse=$1 AND earliest_valid_at <= $2 AND status='proposed'
		 ORDER BY earliest_valid_at`, location, estimatedCompletion)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*models.PlanBOLPairing
	for rows.Next() {
		p, err := scanPairing(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func (r *PairingRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status models.PairingStatus) error {
	_, err := r.db.Exec(ctx, `UPDATE plan_bol_pairing SET status=$2 WHERE id=$1`, id, string(status))
	return err
}
