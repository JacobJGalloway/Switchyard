CREATE TABLE hos_limit (
    id                              UUID    PRIMARY KEY,
    state_code                      TEXT    NOT NULL,
    daily_driving_limit_hours       NUMERIC NOT NULL,  -- max driving hours before mandated rest
    daily_period_hours              NUMERIC NOT NULL,  -- total on-duty window (driving + other duty)
    rest_period_hours               NUMERIC NOT NULL,  -- mandated rest after daily limit is hit
    weekly_limit_hours              NUMERIC NOT NULL,  -- rolling weekly driving cap
    weekly_period_days              INTEGER NOT NULL,  -- days in the weekly window
    weekly_reset_hours              NUMERIC NOT NULL,  -- rest hours required to reset weekly clock
    sleeper_cab_min_hours           NUMERIC,           -- min hours in sleeper during rest (null = no cab rule)
    short_haul_radius_miles         INTEGER,           -- exemption radius in air miles (null = no exemption)
    adverse_weather_extension_hours NUMERIC,           -- additional driving hours permitted in adverse weather
    break_required_after_hours      NUMERIC NOT NULL DEFAULT 8, -- 30-min break trigger threshold
    effective_from                  DATE    NOT NULL,
    notes                           TEXT
);

CREATE UNIQUE INDEX hos_limit_state_effective ON hos_limit (state_code, effective_from);
