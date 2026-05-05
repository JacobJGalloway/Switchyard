-- HOSLimit seed data — Midwest operating states.
-- IL: 60/7 only. All others: 60/7 and 70/8 except IA (70/8 confirmed; 60/7 pending verification).
-- WI in-state extended limits (12hr/16hr) apply only when route stays entirely within WI.
-- IA 16-hr short-haul exception (once per 7 days) is a v1.2 candidate — not enforced at runtime.
-- Agricultural exemptions for IA and MI need further verification before modeling.

INSERT INTO hos_limit (
    id, state_code, cycle_label,
    daily_driving_limit_hours, daily_period_hours, rest_period_hours,
    weekly_limit_hours, weekly_period_days, weekly_reset_hours,
    sleeper_cab_min_hours, short_haul_radius_miles, adverse_weather_extension_hours,
    break_required_after_hours, sleeper_split_allowed, sleeper_split_options,
    effective_from, notes
) VALUES

-- Illinois — 60/7 only
(gen_random_uuid(), 'IL', '60/7',
 11, 14, 10, 60, 7, 34,
 7.0, 150, 2, 8, false, null, '2024-01-01',
 'Sleeper cab: 7-8 hrs min during 10 hr rest. Short-haul (150 air mi): no ELD required if within 11 hr driving limit. Adverse weather: +2 hr driving within 14 hr period.'),

-- Wisconsin — 60/7
(gen_random_uuid(), 'WI', '60/7',
 11, 14, 10, 60, 7, 34,
 7.0, 150, 2, 8, false, null, '2024-01-01',
 'Federal defaults. In-state routes only (all stops within WI): 12 hr driving / 16 hr period permitted. Short-haul (150 air mi): logging exempt if within 14 hr period. Service must verify all stops within WI before applying extended limits.'),

-- Wisconsin — 70/8
(gen_random_uuid(), 'WI', '70/8',
 11, 14, 10, 70, 8, 34,
 7.0, 150, 2, 8, false, null, '2024-01-01',
 'Federal defaults. In-state routes only (all stops within WI): 12 hr driving / 16 hr period permitted. Short-haul (150 air mi): logging exempt if within 14 hr period. Service must verify all stops within WI before applying extended limits.'),

-- Indiana — 60/7
(gen_random_uuid(), 'IN', '60/7',
 11, 14, 10, 60, 7, 34,
 6.5, 150, 2, 8, false, null, '2024-01-01',
 'Sleeper cab: 6.5 hrs min during 10 hr rest (less than federal 7 hr minimum). Short-haul (150 air mi): logging exempt if within 14 hr period.'),

-- Indiana — 70/8
(gen_random_uuid(), 'IN', '70/8',
 11, 14, 10, 70, 8, 34,
 6.5, 150, 2, 8, false, null, '2024-01-01',
 'Sleeper cab: 6.5 hrs min during 10 hr rest (less than federal 7 hr minimum). Short-haul (150 air mi): logging exempt if within 14 hr period.'),

-- Arkansas — 60/7
(gen_random_uuid(), 'AR', '60/7',
 11, 14, 10, 60, 7, 34,
 null, null, 2, 8, false, null, '2024-01-01',
 'Intrastate min. age: 18; interstate min. age: 21. Agricultural HOS exceptions do not apply during seeding or harvesting seasons.'),

-- Arkansas — 70/8
(gen_random_uuid(), 'AR', '70/8',
 11, 14, 10, 70, 8, 34,
 null, null, 2, 8, false, null, '2024-01-01',
 'Intrastate min. age: 18; interstate min. age: 21. Agricultural HOS exceptions do not apply during seeding or harvesting seasons.'),

-- Tennessee — 60/7
(gen_random_uuid(), 'TN', '60/7',
 11, 14, 10, 60, 7, 34,
 7.0, 150, 2, 8, true, '7/3,8/2', '2024-01-01',
 'Short-haul (150 air mi): logbook exempt but must comply with 14 hr period and 11 hr driving limits. Sleeper split (7/3 or 8/2) pauses the 14 hr period clock.'),

-- Tennessee — 70/8
(gen_random_uuid(), 'TN', '70/8',
 11, 14, 10, 70, 8, 34,
 7.0, 150, 2, 8, true, '7/3,8/2', '2024-01-01',
 'Short-haul (150 air mi): logbook exempt but must comply with 14 hr period and 11 hr driving limits. Sleeper split (7/3 or 8/2) pauses the 14 hr period clock.'),

-- Kentucky — 60/7
(gen_random_uuid(), 'KY', '60/7',
 11, 14, 10, 60, 7, 34,
 7.0, null, null, 8, true, '7/3,8/2', '2024-01-01',
 '14 hr window is continuous — does not pause for breaks, fuel, or traffic unless using split sleeper (7/3 or 8/2). As of 2026: 30-min break may be satisfied by off-duty, sleeper berth, or on-duty-not-driving time.'),

-- Kentucky — 70/8
(gen_random_uuid(), 'KY', '70/8',
 11, 14, 10, 70, 8, 34,
 7.0, null, null, 8, true, '7/3,8/2', '2024-01-01',
 '14 hr window is continuous — does not pause for breaks, fuel, or traffic unless using split sleeper (7/3 or 8/2). As of 2026: 30-min break may be satisfied by off-duty, sleeper berth, or on-duty-not-driving time.'),

-- Michigan — 60/7
(gen_random_uuid(), 'MI', '60/7',
 11, 14, 10, 60, 7, 34,
 null, null, 2, 8, false, null, '2024-01-01',
 'Intrastate agricultural driving exceptions exist but specific regulations are pending verification — do not apply to agricultural routes until confirmed.'),

-- Michigan — 70/8
(gen_random_uuid(), 'MI', '70/8',
 11, 14, 10, 70, 8, 34,
 null, null, 2, 8, false, null, '2024-01-01',
 'Intrastate agricultural driving exceptions exist but specific regulations are pending verification — do not apply to agricultural routes until confirmed.'),

-- Iowa — 70/8 only (60/7 pending driver contact verification)
(gen_random_uuid(), 'IA', '70/8',
 11, 14, 10, 70, 8, 34,
 7.0, 150, null, 8, true, '7/3,8/2', '2024-01-01',
 'Short-haul (150 air mi): 16-hr work period exception permitted once per 7 days — v1.2 candidate, not enforced at runtime; flag for review if active routes reach October. Sleeper splits available to pause 14 hr clock. 60/7 cycle pending verification.'),

-- Missouri — 60/7
(gen_random_uuid(), 'MO', '60/7',
 11, 14, 10, 60, 7, 34,
 null, null, null, 8, false, null, '2024-01-01',
 'Intrastate rules mirror interstate. No additional state-specific exemptions identified.'),

-- Missouri — 70/8
(gen_random_uuid(), 'MO', '70/8',
 11, 14, 10, 70, 8, 34,
 null, null, null, 8, false, null, '2024-01-01',
 'Intrastate rules mirror interstate. No additional state-specific exemptions identified.');
