CREATE TABLE truck_inventory_snapshot (
    id                 UUID        PRIMARY KEY,
    plan_bol_id        UUID        NOT NULL REFERENCES plan_bol_record(id),
    plan_bol_stop_id   UUID        NOT NULL REFERENCES plan_bol_stop(id),
    sku_id             TEXT        NOT NULL,
    quantity_loaded    INTEGER     NOT NULL,
    quantity_remaining INTEGER     NOT NULL,
    snapshot_at        TIMESTAMPTZ NOT NULL DEFAULT now()
);
