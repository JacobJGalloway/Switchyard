package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/JacobJGalloway/switchyard-go/internal/models"
)

type DriverRepo struct{ db *pgxpool.Pool }

func NewDriverRepo(db *pgxpool.Pool) *DriverRepo { return &DriverRepo{db: db} }

func (r *DriverRepo) Create(ctx context.Context, d *models.Driver) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO driver (id, name, home_warehouse_id, auth0_user_id, license_state, is_active, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		d.ID, d.Name, d.HomeWarehouseID, d.Auth0UserID, d.LicenseState, d.IsActive, d.CreatedAt)
	return err
}

func (r *DriverRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Driver, error) {
	d := &models.Driver{}
	err := r.db.QueryRow(ctx,
		`SELECT id, name, home_warehouse_id, auth0_user_id, license_state, is_active, created_at
		 FROM driver WHERE id=$1`, id).
		Scan(&d.ID, &d.Name, &d.HomeWarehouseID, &d.Auth0UserID, &d.LicenseState, &d.IsActive, &d.CreatedAt)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (r *DriverRepo) GetAll(ctx context.Context) ([]*models.Driver, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, name, home_warehouse_id, auth0_user_id, license_state, is_active, created_at
		 FROM driver ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*models.Driver
	for rows.Next() {
		d := &models.Driver{}
		if err := rows.Scan(&d.ID, &d.Name, &d.HomeWarehouseID, &d.Auth0UserID, &d.LicenseState, &d.IsActive, &d.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, rows.Err()
}

func (r *DriverRepo) Update(ctx context.Context, d *models.Driver) error {
	_, err := r.db.Exec(ctx,
		`UPDATE driver SET name=$2, home_warehouse_id=$3, auth0_user_id=$4, license_state=$5, is_active=$6
		 WHERE id=$1`,
		d.ID, d.Name, d.HomeWarehouseID, d.Auth0UserID, d.LicenseState, d.IsActive)
	return err
}
