package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/JacobJGalloway/switchyard-go/internal/models"
)

type WarehouseRepo struct{ db *pgxpool.Pool }

func NewWarehouseRepo(db *pgxpool.Pool) *WarehouseRepo { return &WarehouseRepo{db: db} }

func (r *WarehouseRepo) GetAll(ctx context.Context) ([]*models.Warehouse, error) {
	rows, err := r.db.Query(ctx, `SELECT id, region, created_at FROM warehouse ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanWarehouses(rows)
}

func (r *WarehouseRepo) GetByRegion(ctx context.Context, region string) ([]*models.Warehouse, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, region, created_at FROM warehouse WHERE region = $1 ORDER BY id`, region)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanWarehouses(rows)
}

func (r *WarehouseRepo) Create(ctx context.Context, w *models.Warehouse) error {
	if w.CreatedAt.IsZero() {
		w.CreatedAt = time.Now().UTC()
	}
	_, err := r.db.Exec(ctx,
		`INSERT INTO warehouse (id, region, created_at) VALUES ($1, $2, $3)`,
		w.ID, w.Region, w.CreatedAt)
	return err
}

func scanWarehouses(rows interface {
	Next() bool
	Scan(...any) error
	Err() error
}) ([]*models.Warehouse, error) {
	var out []*models.Warehouse
	for rows.Next() {
		w := &models.Warehouse{}
		if err := rows.Scan(&w.ID, &w.Region, &w.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, w)
	}
	if out == nil {
		out = []*models.Warehouse{}
	}
	return out, rows.Err()
}
