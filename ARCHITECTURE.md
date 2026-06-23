# Switchyard — Architecture v1.3

> **Intended audience:** Claude Code, sprint planning sessions, project handoff.
> **Status:** Pre-sprint architectural design. Approved scope, pending sprint start.
> **Parent system:** [Switchyard](https://github.com/JacobJGalloway/Switchyard) (.NET Core 10 / Go / React / PostgreSQL)
> **Sprint model:** Two-week sprint. Scope not completed by end of sprint moves to v1.4.
> **Sprint goal:** Demo stable — pilot client and user feedback ready.
> **Git history:** Refer to `main` branch merge commit notes for v1.2 and prior landed work.

---

## Table of Contents

1. [Sprint Goal & Definition of Done](#1-sprint-goal--definition-of-done)
2. [v1.3 Scope — In Sprint](#2-v13-scope--in-sprint)
3. [v1.4 Scope — Demo Stable Hardening](#3-v14-scope--demo-stable-hardening)
4. [Feature: Mid-BOL Transfer Stops](#4-feature-mid-bol-transfer-stops)
5. [Feature: Demo Reset / Reseed Script](#5-feature-demo-reset--reseed-script)
6. [Feature: Two-Company Demo Seed](#6-feature-two-company-demo-seed)
7. [Feature: Dispatch Board Dark Mode Nuance + Favicon](#7-feature-dispatch-board-dark-mode-nuance--favicon)
8. [Feature: SKU Unit Price](#8-feature-sku-unit-price)
9. [Human-in-the-Loop Checkpoints](#9-human-in-the-loop-checkpoints)
10. [Out of Scope — v1.3](#10-out-of-scope--v13)

---

## 1. Sprint Goal & Definition of Done

**Demo stable** means: a pilot client can be walked through a live operational demo session without encountering broken UI states, missing data, or incomplete workflows. The board must look and behave like a real working day in dispatch.

A v1.3 release to `main` is considered **complete** when:

- [ ] Mid-BOL driver transfer is functional end-to-end (stop type, custody chain, board state)
- [ ] Demo reset script produces a date-relative board that looks like a live operational day on any run date
- [ ] Two-company seed is in place (Company A and Company B with distinct palettes and complexity levels)
- [ ] Dispatch board dark mode is visually polished with no nuance regressions; favicon is swapped
- [ ] SKU unit price field is in the inventory model, seeded with owner-provided values, and surfaces correctly in revenue vs. profit analytics
- [ ] All v1.3 items above are merged to `main`

---

## 2. v1.3 Scope — In Sprint

| # | Feature | Weight | Notes |
|---|---------|--------|-------|
| 1 | Mid-BOL transfer stops | Heavy | Schema + service + board state changes |
| 2 | Demo reset / reseed script | Medium | Date-relative logic is the core complexity |
| 3 | Two-company demo seed | Medium | Dependent on reset script being in place |
| 4 | Dispatch board dark mode nuance + favicon | Light | UI polish; no backend changes |
| 5 | SKU unit price | Light (code) / Tedious (seed) | Human checkpoint required for seed values |

---

## 3. v1.4 Scope — Demo Stable Hardening

These items are **pilot-required** but deferred from v1.3. They are isolated in scope with no cross-functional dependencies, making them safe to defer without risk to the v1.3 demo.

| # | Feature | Rationale for deferral |
|---|---------|------------------------|
| 1 | ARIA compliance audit — board columns, cards, icon-only buttons, skip-nav | Isolated; no functional dependencies |
| 2 | Color contrast audit (WCAG AA) — light and dark themes | Isolated; no functional dependencies |
| 3 | Rolling refresh tokens for Auth0 sessions | Isolated; no cross-functional dependencies |
| 4 | Mid-week board display state | Refinement on top of two-company seed; not a demo blocker |

> These land in v1.4 as a hardening sprint before the pilot goes live. They are not "someday" items.

---

## 4. Feature: Mid-BOL Transfer Stops

### Problem
A driver may be unable to complete a BOL run due to Hours of Service (HOS) limits or an emergency. Currently there is no mechanism to formally hand off the load to a different driver mid-route while preserving BOL continuity and dispatch board accuracy.

### Current State
- `DriverBOLAssignment` model exists in the Go service
- It supports a single driver assigned to a BOL from origin to destination
- No concept of a mid-route custody transfer exists in the schema or the board

### Target State
Introduce a `transfer` stop type that creates a formal custody checkpoint. The BOL remains continuous; the custody chain records the handoff.

#### Schema Changes (Go / PostgreSQL)

**New: `transfer` stop type**

Add `transfer` as a valid value in the stop type enum alongside existing stop types (pickup, delivery, etc.).

**`DriverBOLAssignment` restructuring**

The existing model assumes one driver per BOL. Restructure to support a custody chain:

```
DriverBOLAssignment
  - id
  - bol_id (FK)
  - driver_id (FK)
  - equipment_id (FK, nullable — equipment may not transfer with driver)
  - segment_start_stop_id (FK → Stop)   -- where this driver's responsibility begins
  - segment_end_stop_id (FK → Stop)     -- where this driver's responsibility ends (transfer or final destination)
  - assigned_at (timestamp)
  - transfer_reason (enum: hos_limit | emergency | planned | other)
  - notes (text, nullable)
```

A BOL with no transfer has one `DriverBOLAssignment` record (origin → destination). A BOL with one mid-route transfer has two records: origin → transfer stop, transfer stop → destination.

#### Migration Path for Existing Records

Existing `DriverBOLAssignment` rows have no stop anchor FKs. When the new columns are added, backfill:
- `segment_start_stop_id` → the BOL's first stop (origin)
- `segment_end_stop_id` → the BOL's last stop (final destination)

This keeps existing single-driver BOLs working without changes to queries or board rendering.

> **Single-driver load is the default and the common case.** Transfer stops exist to handle HOS limits and emergencies — not as a routine pattern. The restructured schema must not change the shape or behavior of single-driver BOL records for any existing query or board view.

**Transfer Stop record**

When a transfer is initiated, create a `Stop` record of type `transfer` at the handoff location. This stop anchors both the outgoing and incoming `DriverBOLAssignment` segments.

#### API Changes (Go Service)

- `POST /bols/{id}/transfer` — initiate a driver transfer; accepts incoming driver, equipment (optional), transfer stop location, and reason; creates the transfer stop and the new `DriverBOLAssignment` segment
- `GET /bols/{id}/assignments` — return the full custody chain for a BOL (ordered by segment)
- Update `GET /bols/{id}` to include the custody chain in the response

#### Dispatch Board State

- A BOL mid-transfer should surface on the board with a visual indicator (e.g., a custody handoff badge on the card)
- The board card should reflect the **current active driver** (the most recent `DriverBOLAssignment` segment)
- The outgoing driver's assignment should be marked complete at the transfer stop, not dropped

#### Constraints & Guardrails for Claude Code

- Do **not** break existing BOL cards that have no transfer — the single-assignment path must remain the default and continue working without changes
- The `transfer_reason` field is required on creation; do not allow null
- Equipment transfer is optional — a driver swap does not always mean an equipment swap
- If unsure whether a schema change affects existing BOL queries, surface the question rather than assuming

---

## 5. Feature: Demo Reset / Reseed Script

### Problem
A static seed produces a board that looks stale relative to today's date. A demo run on any given day should show a board that looks like a live operational morning, not a snapshot frozen at seed time.

### Target State
A script (Go CLI or SQL + shell) that:

1. Truncates or soft-deletes all existing demo data
2. Reseeds with dates calculated **relative to today** at runtime, targeting a **Thursday/Friday operational state** — the end of the work week, where operational complexity is at its peak
3. Can be run at any time before a demo session and produce a consistent, realistic board state
4. Is idempotent — safe to run multiple times without accumulating data

**Thursday/Friday board state requirements:**
- **Driver HOS:** A mix across the driver pool — some drivers at their weekly HOS limit (unavailable for the remainder of the week), others only at their daily limit (available again after a 10-hour reset), and some with hours remaining
- **Equipment:** A mix of statuses — some units in active service, some with maintenance flags from normal operational wear during the week, others ready and available for weekend use
- **BOL state:** Active BOLs in progress, completed BOLs from earlier in the week, and the full range of stop states visible on the board

The goal is that every general-status variant a dispatcher sees in production is visible on the board during the demo — no status type should be absent from the demo view.

### Design Notes
- Date offsets should be defined as named constants or a config block at the top of the script, not scattered through seed logic — makes them easy to tune without hunting through code
- The script should log what it seeded (counts, date ranges) so it's easy to confirm the board state before a demo
- Seed data must be compatible with the two-company demo seed (see Feature 6)
- Implementation: Go CLI (consistent with the rest of the Go backend; testable)

---

## 6. Feature: Two-Company Demo Seed

### Dependency
The demo reset script (Feature 5) must be in place before this seed is built on top of it.

### Target State
Two distinct company datasets within the same demo environment:

**Company A — Monday Morning (Default Brand)**
- Board state: start of week, fresh assignments, clean slate feel
- Brand: default Switchyard palette
- Complexity: baseline — shows core dispatch flow without edge cases

**Company B — Mid-Week Complexity (Client Palette Override)**
- Board state: mid-week operational complexity — BOLs in progress, some transfers, varied stop states
- Brand: client palette override (demonstrates the theming system to the pilot client)
- Complexity: elevated — shows the board handling a real operational load, not just a single truck

### Design Notes
- The two companies should be selectable or togglable during the demo session — a presenter should not have to re-run the full seed to switch contexts
- Company B's mid-week complexity should include at least one BOL that has used the transfer stop feature (Feature 4) so that feature is visible in the demo without requiring a live demo of the transfer flow itself

---

## 7. Feature: Dispatch Board Dark Mode Nuance + Favicon

### Scope
UI polish pass only. No backend changes.

### Dark Mode
- Audit the dispatch board in dark mode for nuance regressions — colors that look acceptable in light mode but lose contrast, hierarchy, or legibility in dark
- Focus areas: board column headers, card status indicators, stop type badges, empty state messaging
- This is a visual polish pass, not a full WCAG audit (that is v1.4)

### Favicon
- Swap the current favicon for the Switchyard branded favicon
- Verify it renders correctly in browser tab, bookmarks, and PWA contexts if applicable

---

## 8. Feature: SKU Unit Price

### Scope
Extend the inventory model to hold a unit price per SKU. This enables revenue vs. profit analytics by giving the system the price side of the margin equation.

### Schema Change

```
SKU (existing model — extend)
  + unit_price (decimal, not null)
  + price_currency (varchar, default 'USD')
  + price_effective_date (date, nullable — for future price history support)
```

### Analytics Surface
Once unit price is in the model, surface revenue vs. profit in the analytics views that currently show cost data. The formula is straightforward — Claude Code can implement the calculations. No estimate work is needed from Claude Code on the pricing values themselves.

### Migration
- The migration script to add the column is Claude Code's work
- **The seed values (actual unit prices per SKU) require owner input** — see Human-in-the-Loop Checkpoints below
- The seed data entry is expected to be tedious (manual per-SKU values) but not architecturally complex

---

## 9. Human-in-the-Loop Checkpoints

These are points in the sprint where Claude Code should **pause and surface** rather than proceed independently.

| Checkpoint | Feature | What to surface |
|------------|---------|-----------------|
| SKU unit prices | SKU Unit Price (#8) | Once the migration script and seed template are ready, provide the seed file/template for owner to fill in unit price values before seeding. Do not guess or estimate prices. |
| Transfer stop edge cases | Mid-BOL Transfer (#4) | If any existing BOL query or board query is ambiguous about how to handle the custody chain, surface the question with a specific example rather than choosing an approach silently. |
| Company B palette | Two-Company Seed (#6) | The client palette override values for Company B are owner-defined. Surface a placeholder config and ask for the values before finalizing the seed. |

---

## 10. Out of Scope — v1.3

The following are explicitly **not** in this sprint. Do not begin work on these without a scope change.

- ARIA compliance audit
- Color contrast audit (WCAG AA)
- Rolling refresh tokens (Auth0)
- Mid-week board display state refinement
- Operating cost / tow rate analytics
- Dead-head run support
- Any features listed as "Backlog" in the README
