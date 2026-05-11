CREATE TABLE delivery_confirmation (
    id               UUID        PRIMARY KEY,
    plan_bol_stop_id UUID        NOT NULL REFERENCES plan_bol_stop(id),
    driver_id        UUID        NOT NULL REFERENCES driver(id),
    confirmed_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    -- invoice_id is set after the internal invoice is generated; no FK here
    -- because internal_invoice is created after confirmation (see migration 007)
    invoice_id       UUID
);
