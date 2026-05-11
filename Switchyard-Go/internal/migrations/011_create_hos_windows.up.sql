CREATE TABLE hos_window (
    id                UUID        PRIMARY KEY,
    driver_id         UUID        NOT NULL REFERENCES driver(id),
    window_start      TIMESTAMPTZ NOT NULL,
    daily_hours_used  NUMERIC     NOT NULL DEFAULT 0,
    weekly_hours_used NUMERIC     NOT NULL DEFAULT 0,
    -- last_activity_at is when the driver last departed a stop — used to calculate
    -- driving hours for each leg (leg_hours = stop_logged_at - last_activity_at)
    last_activity_at  TIMESTAMPTZ,
    -- FMCSA 30-minute break required before driving after 8 cumulative on-duty hours
    break_30_taken    BOOLEAN     NOT NULL DEFAULT false,
    break_30_at       TIMESTAMPTZ,
    mandated_stop_at  TIMESTAMPTZ,
    eld_stop_ref      TEXT
);
