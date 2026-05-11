CREATE TABLE driver (
    id                UUID        PRIMARY KEY,
    name              TEXT        NOT NULL,
    home_warehouse_id TEXT        NOT NULL,
    auth0_user_id     TEXT        NOT NULL,
    license_state     TEXT        NOT NULL,
    is_active         BOOLEAN     NOT NULL DEFAULT true,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now()
);
