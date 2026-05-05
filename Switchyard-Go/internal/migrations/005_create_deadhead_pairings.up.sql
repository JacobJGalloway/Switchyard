CREATE TABLE plan_bol_pairing (
    id                UUID        PRIMARY KEY,
    active_bol_id     UUID        NOT NULL REFERENCES plan_bol_record(id),
    deadhead_bol_id   UUID        NOT NULL REFERENCES plan_bol_record(id),
    paired_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    earliest_valid_at TIMESTAMPTZ NOT NULL,
    origin_warehouse  TEXT        NOT NULL,
    status            TEXT        NOT NULL
);
