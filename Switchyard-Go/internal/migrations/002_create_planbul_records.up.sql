CREATE TABLE plan_bol_record (
    id                UUID        PRIMARY KEY,
    driver_id         UUID        NOT NULL REFERENCES driver(id),
    originating_wh_id TEXT        NOT NULL,
    status            TEXT        NOT NULL,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    submitted_at      TIMESTAMPTZ,
    fulfilled_at      TIMESTAMPTZ
);
