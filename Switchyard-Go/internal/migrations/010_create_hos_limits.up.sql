CREATE TABLE hos_limit (
    id                 UUID    PRIMARY KEY,
    state_code         TEXT    NOT NULL,
    daily_limit_hours  NUMERIC NOT NULL,
    weekly_limit_hours NUMERIC NOT NULL,
    effective_from     DATE    NOT NULL,
    notes              TEXT
);

CREATE UNIQUE INDEX hos_limit_state_effective ON hos_limit (state_code, effective_from);
