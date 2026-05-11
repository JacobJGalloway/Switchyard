package models

import (
	"time"

	"github.com/google/uuid"
)

// HOSLimit stores state-level HOS rules. New states are added by inserting rows —
// no code change required (see ARCHITECTURE.md §4.3).
// DailyDrivingLimitHours and DailyPeriodHours are distinct: a driver may hit
// the driving cap before the on-duty window closes, or vice versa.
type HOSLimit struct {
	ID                            uuid.UUID `json:"id"`
	StateCode                     string    `json:"state_code"`
	DailyDrivingLimitHours        float64   `json:"daily_driving_limit_hours"`
	DailyPeriodHours              float64   `json:"daily_period_hours"`
	RestPeriodHours               float64   `json:"rest_period_hours"`
	WeeklyLimitHours              float64   `json:"weekly_limit_hours"`
	WeeklyPeriodDays              int       `json:"weekly_period_days"`
	WeeklyResetHours              float64   `json:"weekly_reset_hours"`
	SleeperCabMinHours            *float64  `json:"sleeper_cab_min_hours,omitempty"`
	ShortHaulRadiusMiles          *int      `json:"short_haul_radius_miles,omitempty"`
	AdverseWeatherExtensionHours  *float64  `json:"adverse_weather_extension_hours,omitempty"`
	BreakRequiredAfterHours       float64   `json:"break_required_after_hours"`
	// SleeperSplitAllowed indicates the split sleeper berth provision is available,
	// which pauses the daily period clock during the sleeper portion of the rest.
	SleeperSplitAllowed           bool      `json:"sleeper_split_allowed"`
	SleeperSplitOptions           *string   `json:"sleeper_split_options,omitempty"` // e.g. "7/3,8/2"
	// CycleLabel is "60/7" or "70/8". States that permit both have one row per cycle.
	// The HOS service selects the row matching the driver's assigned operating cycle.
	CycleLabel                    string    `json:"cycle_label"`
	EffectiveFrom                 time.Time `json:"effective_from"`
	Notes                         *string   `json:"notes,omitempty"`
}

// HOSWindow tracks accumulated hours for one driver across a run window.
// Break30Taken / Break30At enforce the FMCSA rule: a 30-minute break is required
// before driving after 8 cumulative on-duty hours within the window.
type HOSWindow struct {
	ID              uuid.UUID  `json:"id"`
	DriverID        uuid.UUID  `json:"driver_id"`
	WindowStart     time.Time  `json:"window_start"`
	DailyHoursUsed  float64    `json:"daily_hours_used"`
	WeeklyHoursUsed float64    `json:"weekly_hours_used"`
	LastActivityAt  *time.Time `json:"last_activity_at,omitempty"`
	Break30Taken    bool       `json:"break_30_taken"`
	Break30At       *time.Time `json:"break_30_at,omitempty"`
	MandatedStopAt  *time.Time `json:"mandated_stop_at,omitempty"`
	ELDStopRef      *string    `json:"eld_stop_ref,omitempty"`
}
