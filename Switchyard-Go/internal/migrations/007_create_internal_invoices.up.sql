CREATE TABLE internal_invoice (
    id           UUID        PRIMARY KEY,
    store_id     TEXT        NOT NULL,
    plan_bol_id  UUID        NOT NULL REFERENCES plan_bol_record(id),
    line_items   JSONB       NOT NULL,
    output_path  TEXT        NOT NULL,
    generated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Back-fill the FK now that both tables exist
ALTER TABLE delivery_confirmation
    ADD CONSTRAINT fk_delivery_confirmation_invoice
    FOREIGN KEY (invoice_id) REFERENCES internal_invoice(id);
