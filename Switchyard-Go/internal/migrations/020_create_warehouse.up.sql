CREATE TABLE warehouse (
    id         TEXT        PRIMARY KEY,
    region     TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
