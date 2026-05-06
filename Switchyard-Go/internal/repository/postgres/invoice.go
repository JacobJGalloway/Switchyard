package postgres

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/JacobJGalloway/switchyard-go/internal/models"
)

type InvoiceRepo struct{ db *pgxpool.Pool }

func NewInvoiceRepo(db *pgxpool.Pool) *InvoiceRepo { return &InvoiceRepo{db: db} }

func (r *InvoiceRepo) CreateConfirmation(ctx context.Context, c *models.DeliveryConfirmation) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO delivery_confirmation (id, plan_bol_stop_id, driver_id, confirmed_at, invoice_id)
		 VALUES ($1,$2,$3,$4,$5)`,
		c.ID, c.PlanBOLStopID, c.DriverID, c.ConfirmedAt, c.InvoiceID)
	return err
}

func (r *InvoiceRepo) GetConfirmation(ctx context.Context, id uuid.UUID) (*models.DeliveryConfirmation, error) {
	c := &models.DeliveryConfirmation{}
	err := r.db.QueryRow(ctx,
		`SELECT id, plan_bol_stop_id, driver_id, confirmed_at, invoice_id
		 FROM delivery_confirmation WHERE id=$1`, id).
		Scan(&c.ID, &c.PlanBOLStopID, &c.DriverID, &c.ConfirmedAt, &c.InvoiceID)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *InvoiceRepo) GetConfirmationsByBOL(ctx context.Context, planBOLID uuid.UUID) ([]*models.DeliveryConfirmation, error) {
	rows, err := r.db.Query(ctx,
		`SELECT dc.id, dc.plan_bol_stop_id, dc.driver_id, dc.confirmed_at, dc.invoice_id
		 FROM delivery_confirmation dc
		 JOIN plan_bol_stop s ON s.id = dc.plan_bol_stop_id
		 WHERE s.plan_bol_id=$1 ORDER BY dc.confirmed_at`, planBOLID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*models.DeliveryConfirmation
	for rows.Next() {
		c := &models.DeliveryConfirmation{}
		if err := rows.Scan(&c.ID, &c.PlanBOLStopID, &c.DriverID, &c.ConfirmedAt, &c.InvoiceID); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func (r *InvoiceRepo) SetConfirmationInvoice(ctx context.Context, confirmationID, invoiceID uuid.UUID) error {
	_, err := r.db.Exec(ctx,
		`UPDATE delivery_confirmation SET invoice_id=$2 WHERE id=$1`, confirmationID, invoiceID)
	return err
}

func (r *InvoiceRepo) CreateInvoice(ctx context.Context, i *models.InternalInvoice) error {
	lineItemsJSON, _ := json.Marshal(i.LineItems)
	_, err := r.db.Exec(ctx,
		`INSERT INTO internal_invoice (id, store_id, plan_bol_id, line_items, output_path, generated_at)
		 VALUES ($1,$2,$3,$4,$5,$6)`,
		i.ID, i.StoreID, i.PlanBOLID, lineItemsJSON, i.OutputPath, i.GeneratedAt)
	return err
}

func (r *InvoiceRepo) GetInvoice(ctx context.Context, id uuid.UUID) (*models.InternalInvoice, error) {
	i := &models.InternalInvoice{}
	var lineItemsRaw []byte
	err := r.db.QueryRow(ctx,
		`SELECT id, store_id, plan_bol_id, line_items, output_path, generated_at
		 FROM internal_invoice WHERE id=$1`, id).
		Scan(&i.ID, &i.StoreID, &i.PlanBOLID, &lineItemsRaw, &i.OutputPath, &i.GeneratedAt)
	if err != nil {
		return nil, err
	}
	_ = json.Unmarshal(lineItemsRaw, &i.LineItems)
	return i, nil
}

func (r *InvoiceRepo) GetInvoicesByStore(ctx context.Context, storeID string) ([]*models.InternalInvoice, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, store_id, plan_bol_id, line_items, output_path, generated_at
		 FROM internal_invoice WHERE store_id=$1 ORDER BY generated_at DESC`, storeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*models.InternalInvoice
	for rows.Next() {
		i := &models.InternalInvoice{}
		var lineItemsRaw []byte
		if err := rows.Scan(&i.ID, &i.StoreID, &i.PlanBOLID, &lineItemsRaw, &i.OutputPath, &i.GeneratedAt); err != nil {
			return nil, err
		}
		_ = json.Unmarshal(lineItemsRaw, &i.LineItems)
		out = append(out, i)
	}
	return out, rows.Err()
}
