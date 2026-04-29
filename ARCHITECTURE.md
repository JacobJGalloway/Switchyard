# Switchyard — Architecture v1.1

> **Intended audience:** Claude Code, onboarding developers, project planning sessions.
> **Status:** Pre-implementation architectural design. Approved concept, pending sprint start (late April / early May 2026).
> **Parent system:** [Switchyard (.NET)](https://github.com/JacobJGalloway/Switchyard) (.NET Core / React / SQLite)
> **Sprint model:** Two-week sprints. Scope not completed by end of sprint moves to v1.2.

---

## Table of Contents

1. [Problem Statement](#1-problem-statement)
2. [Scope — What This Project Owns](#2-scope--what-this-project-owns)
3. [Event Handler / Listener — Go Entry Point](#3-event-handler--listener--go-entry-point)
4. [Core Business Rules](#4-core-business-rules)
5. [Digital Dispatch Whiteboard — Kanban Board Design](#5-digital-dispatch-whiteboard--kanban-board-design)
6. [Project Folder Structure](#6-project-folder-structure)
7. [Data Model — Go Backend Entities (PostgreSQL)](#7-data-model--go-backend-entities-postgresql)
8. [API Endpoints (Go Backend)](#8-api-endpoints-go-backend)
9. [Auth0 — Permissions](#9-auth0--permissions)
10. [Deployment — Island Architecture](#10-deployment--island-architecture)
11. [Environment Variables](#11-environment-variables)
12. [Technology Stack](#12-technology-stack)
13. [Integration Points with Switchyard (.NET)](#13-integration-points-with-switchyard-net)
14. [Wanted Features Carried Forward](#14-wanted-features-carried-forward)
15. [Handoff Instructions for Claude Code](#15-handoff-instructions-for-claude-code)
16. [v1.2 Candidate Features](#16-v12-candidate-features)

---

## 1. Problem Statement

The existing Switchyard (.NET) system generates Bills of Lading (BOL) and tracks inventory
transfers between warehouses and stores. Three unsolved problems remain from MVP:

1. **Inventory shortfall detection** is implicit — no system catches the case where the originating
   warehouse cannot cover all store stops before the BOL is committed.
2. **Route sequencing** is manual — dispatchers cannot verify that a driver's stop order is
   inventory-safe before the truck rolls.
3. **Dead-head run pairing** is informal — dispatchers match return runs mentally, with no system
   enforcing the 4-hour pre-arrangement rule or surfacing eligible BOLs.

Additionally, v1.1 introduces three new operational surfaces not present in MVP:

4. **Driver HOS (Hours of Service) tracking** — state-level daily and weekly limits must be
   enforced and surfaced to dispatchers before assignment decisions are made.
5. **Digital Dispatch Whiteboard** — a Kanban-style real-time board replacing the physical
   whiteboard. Tracks drivers, equipment, and BOLs through their full lifecycle.
6. **Equipment management** — trucks and tractors must be tracked independently of drivers,
   including scheduled maintenance windows and roadside or depot breakdown states.

---

## 2. Scope — What This Project Owns

```
Switchyard (.NET) (existing — do not modify core logic)
  └── BOL lifecycle, stop CUD (stops are children of BOL in the .NET domain),
      inventory CRUD, Auth0, CQRS read replica, SQLite

Switchyard (Go) (this project — extends upstream and downstream)
  ├── Event handler / listener — single Go entry point, routes all workflow events
  ├── PlanBOL route planning and constraint resolution (planning phase only)
  ├── Truck inventory state management across a planned run
  ├── Dead-head BOL pairing (4-hour pre-arrangement window)
  ├── Driver entity, HOS state, and run sheet generation
  ├── Driver-to-BOL-to-equipment assignment record
  ├── Digital Dispatch Whiteboard (Kanban board — real-time)
  ├── Equipment (truck/tractor) management — maintenance and breakdown states
  ├── Delivery confirmation per store stop
  ├── Internal store invoice output (env-configurable destination)
  └── Email notification service — workflow completion and alert notifications
```

The Go backend calls the existing Switchyard (.NET) Logistics API. It does **not** replace it.
The existing `POST /BillOfLading/{transactionId}/replace-stop` endpoint becomes an emergency
manual override for dispatchers after BOL commitment — not the primary workflow.

**CUD authority boundary:**

```
Go backend owns CUD for:
  PlanBOLRecord, PlanBOLStop, PlanBOLPairing   (planning phase only)
  Driver entity, HOS state, and run state
  DriverBOLAssignment                           (driver + BOL + equipment linkage)
  Equipment entity (truck/tractor), MaintenanceRecord, BreakdownRecord
  Truck inventory snapshots
  Delivery confirmations
  Internal invoices

.NET backend owns CUD for (unchanged):
  Committed BOLs and their Stop children
  Inventory quantities (Clothing, PPE, Tool)
  Warehouse and store master data
  Auth0 users
```

The seam between systems is the `submit` action on a validated PlanBOLRecord. At that moment
the Go backend calls `POST /BillOfLading` on the Switchyard (.NET) Logistics API, hands off
authority, and the PlanBOL entities become read-only archived records.

---

## 3. Event Handler / Listener — Go Entry Point

The Go backend exposes a single event handler as its primary entry point. This design:

- Holds the single Auth0 M2M token used to authenticate all calls to the .NET APIs
- Receives all incoming workflow events (HTTP, internal triggers, completion callbacks)
- Routes each event to the appropriate internal Go workflow service
- Fires the email notification service on workflow completion or alert threshold

This pattern keeps the Auth0 free tier M2M token count to one, regardless of how many
internal Go workflow services are added in future versions.

```
Incoming event
      │
      ▼
┌─────────────────────┐
│   Event Handler     │  ← Single M2M token lives here
│   (entry point)     │
└──────────┬──────────┘
           │  routes to
   ┌───────┴────────┬──────────────┬──────────────┬──────────────┐
   ▼                ▼              ▼               ▼              ▼
Route          HOS Service    Whiteboard      Equipment      Email
Planner                       Service         Service        Notifier
```

**Email notification service** is a lightweight internal module triggered by the event handler
on workflow completion or alert condition. Uses Go's `net/smtp` or `gomail`. Does not require
its own M2M token — internal to the Go backend.

Notification triggers (initial scope):
- Driver HOS limit approaching (threshold configurable via env var)
- Driver HOS weekly limit reached
- BOL workflow completed (all stops confirmed)
- Dead-head pairing window expiring (< 1 hour remaining)
- Equipment breakdown reported

---

## 4. Core Business Rules (Encode as Service-Layer Constraints)

Read this section fully before writing any code. Ask clarifying questions if any rule
is ambiguous. These rules are hard constraints, not soft warnings.

### 4.1 Truck Inventory State

- Driver always starts at the originating (home) warehouse — this is stop 1, always.
- Stop 1 is auto-processed at BOL creation time. It is an invariant, not a user action.
- Truck maintains a running inventory balance as stops are evaluated in sequence.
- A store stop is only valid if the truck currently holds sufficient quantity of all
  required products for that stop at the moment of arrival.
- **Preferred resolution strategy (Problem 1.1):** When the originating warehouse cannot
  cover all store stops, the constraint solver first attempts warehouse-to-warehouse
  loading — all warehouse pickups are sequenced before any store deliveries begin.
  This is the cleaner operational path and easier for the driver to plan around.
- **Fallback resolution strategy:** If warehouse-to-warehouse loading cannot satisfy the
  full stop sequence, mid-route warehouse insertion is attempted. A warehouse stop may
  be inserted at any point in the sequence where a shortfall is detected, provided that
  warehouse has the required inventory available. Store stops for which the truck already
  holds sufficient inventory may proceed before the additional warehouse stop.
- The truck **must be empty** at the completion of the final stop. A PlanBOLRecord that
  cannot be constructed to result in an empty truck is invalid and must not be submitted.

### 4.2 Dead-Head Run Pairing

- A dead-head (return) BOL must be arranged **at least 4 hours** before the active BOL
  is fulfilled.
- The dead-head BOL must originate from the warehouse geographically closest to the
  active BOL's projected final stop.
- Dead-head runs work the driver back toward his home warehouse without requiring an
  explicit empty return trip.
- The system must expose a query surface for eligible dead-head BOLs given a projected
  final stop location and estimated completion time.
- Dead-head BOL creation is a first-class workflow — the system generates the dead-head
  BOL record, not just pairs existing ones.
- A PlanBOLPairing that does not satisfy the 4-hour window must be rejected at the
  service layer — hard constraint, not a warning.

### 4.3 Driver HOS (Hours of Service)

- Daily and weekly HOS limits are enforced at the state level for the states in which
  the company's warehouses are located (initial scope — not all 50 states).
- HOS limits are stored as configuration per state. New states are added by configuration,
  not code change.
- The event handler evaluates HOS state on every driver stop log event.
- A driver who cannot legally complete a planned run without exceeding their daily or
  weekly limit must not be assigned to that BOL at the planning phase.
- A driver who reaches a mandated stop (HOS-triggered rest) mid-run transitions their
  Whiteboard card to the **Mandated Stop** sub-column. If ELD data is accessible,
  the stop timestamp is pulled from the ELD. Otherwise the dispatcher logs it manually.
- The dead-head timer is paused during a mandated stop.

### 4.4 Driver Stop Logging

- Each stop must be logged individually (regulatory requirement — digital logbook compliance).
- Stop completion triggers a dispatch notification. Dispatcher acknowledges; driver confirms.
- Driver and dispatcher views are separate concerns.

### 4.5 Delivery Confirmation and Internal Invoice

- Each store stop that is successfully processed generates a delivery confirmation record.
- A delivery confirmation triggers an internal invoice for that store (cost center recipient).
- Invoice output destination is controlled by `INVOICE_OUTPUT_PATH` env var.
  Default: `./output/invoices/`. Future: Azure Blob or equivalent — no code change needed.
- Invoice format: structured `.json` matching the existing BOL `.txt` output pattern.

### 4.6 Dead-Head Timer (Post-Delivery)

- The dead-head timer starts when the final stop on the active BOL is confirmed (fulfilled_at).
- The timer gives the dispatcher a window to find and confirm a dead-head return run
  before the driver goes idle or their HOS window closes.
- Timer duration is configurable via `DEADHEAD_SEARCH_WINDOW_HOURS` env var.
- Timer expiry triggers an email notification to the dispatcher.
- The Whiteboard **Delivered** column displays the timer countdown on the driver card.

---

## 5. Digital Dispatch Whiteboard — Kanban Board Design

The whiteboard is a real-time Kanban board served at `GET /` (HTML) and `GET /api/dispatch/board`
(JSON). It replaces the physical dispatch whiteboard.

### Column Structure

```
┌──────────────┬──────────────┬─────────────────────────────┬──────────────┬──────────────┬──────────────┬──────────────┐
│   AVAILABLE  │   PENDING    │         IN TRANSIT           │  DELIVERED   │ HOS LIMITED  │ MAINTENANCE  │  BREAKDOWN   │
│              │  DISPATCH    │                              │              │              │              │              │
│  BOL card    │  BOL card    │  Driver card (primary)       │  Driver card │  Driver card │  Equipment   │  Equipment   │
│  (loaded,    │  (skinny —   │  └─ BOL sub-card             │  + countdown │  (at/near    │  card +      │  card +      │
│  ready for   │  needs to    │                              │  timer for   │  weekly      │  timeframe   │  location    │
│  assignment) │  move now)   │  ┌─────────┬──────────────┐  │  dead-head   │  HOS limit)  │  estimate    │  indicator   │
│              │              │  │IN DELIVERY│MANDATED STOP│  │  search      │              │              │              │
│              │              │  │          │+ ELD time   │  │              │              │              │              │
│              │              │  └─────────┴──────────────┘  │              │              │              │              │
└──────────────┴──────────────┴─────────────────────────────┴──────────────┴──────────────┴──────────────┴──────────────┘
```

### Column Definitions

**AVAILABLE**
- Primary card: BOL
- State: BOL created, inventory loaded onto trailer, no driver/equipment assigned yet
- Dispatcher action: assign driver and equipment to move card to Pending Dispatch

**PENDING DISPATCH** *(skinny column — urgent)*
- Primary card: BOL
- State: driver and equipment assigned, BOL has not yet departed
- Visual treatment: narrow column width signals urgency — this BOL needs to roll
- Dispatcher action: confirm departure to move card to In Transit

**IN TRANSIT**
- Primary card: Driver
- Sub-card: BOL (beneath driver card)
- Two sub-columns:
  - **In Delivery** — driver actively moving between stops
  - **Mandated Stop** — HOS or legally required rest stop. Displays ELD timestamp
    if accessible; dispatcher manual entry if not. Dead-head timer paused.
- Driver card displays: driver name, equipment ID, current stop, HOS status indicator
  (green / yellow approaching / red at limit)

**DELIVERED**
- Primary card: Driver
- State: all BOL stops confirmed, driver away from originating warehouse
- Countdown timer displayed on card — dispatcher window to arrange dead-head return run
- Timer expiry fires email notification to dispatcher
- Dispatcher action: confirm dead-head pairing to clear card from board

**HOS LIMITED**
- Primary card: Driver
- State: driver has reached daily or weekly HOS limit, unavailable for new runs
- Informational column — dispatcher knows not to assign this driver
- Card clears automatically when HOS window resets (configurable per state)

**MAINTENANCE** *(equipment columns, right side of board)*
- Primary card: Equipment (truck or tractor)
- State: scheduled planned maintenance, equipment temporarily unavailable
- Card displays: equipment ID, maintenance description, estimated return timeframe
- Planned maintenance only — unplanned goes to Breakdown

**BREAKDOWN**
- Primary card: Equipment
- Two scenarios on one card, distinguished by location indicator:
  - **Depot breakdown** — equipment made it back to originating warehouse before failure
  - **Roadside breakdown** — equipment failed in the field, driver and load still attached
    (urgent — dispatcher must arrange rescue dispatch)
- Roadside breakdown with load triggers immediate email notification to dispatcher
- Card displays: equipment ID, driver reference (if load attached), breakdown location,
  reported timestamp

### Driver-BOL-Equipment Assignment Record

The assignment record is the linking entity that ties all three together:

```
Driver ──assigned to──► BOL ──on──► Equipment (truck/tractor)
```

One assignment record answers the dispatcher's core question:
"Who is taking what, in which vehicle, to where?"

Assignment is created at Pending Dispatch, transitions with the driver card through
In Transit → Delivered, and is archived when the dead-head return is confirmed.

---

## 6. Project Folder Structure

```
Switchyard-Go/
├── cmd/
│   └── main.go                        # Entry point, DI wiring, config bootstrap
├── internal/
│   ├── events/                        # Event handler — single entry point
│   │   ├── handler.go                 # Receives all events, routes to services
│   │   └── router.go                  # Event type → service mapping
│   ├── handlers/                      # HTTP controllers — bind input, return JSON or HTML
│   │   ├── planbul_handler.go         # PlanBOL create, validate, submit
│   │   ├── driver_handler.go          # Run sheet, stop logging, active BOL
│   │   ├── deadhead_handler.go        # Pairing query, pair, cancel
│   │   ├── whiteboard_handler.go      # Dispatch board, alerts
│   │   └── equipment_handler.go       # Truck/tractor CRUD, maintenance, breakdown
│   ├── services/                      # Business logic — never imports /web
│   │   ├── route_planner.go           # Constraint-satisfaction pass over stop sequence
│   │   ├── truck_inventory.go         # Running inventory state across a planned run
│   │   ├── deadhead_service.go        # 4-hour window pairing logic + timer
│   │   ├── hos_service.go             # State-level HOS daily/weekly limit enforcement
│   │   ├── whiteboard_service.go      # Kanban board state assembly and transitions
│   │   ├── equipment_service.go       # Truck/tractor lifecycle, maintenance, breakdown
│   │   ├── invoice_service.go         # Delivery confirmation + invoice output
│   │   └── notification_service.go    # Email notifications — workflow completion + alerts
│   ├── models/                        # Domain structs
│   │   ├── planbul.go                 # PlanBOLRecord, PlanBOLStop, StopType enum
│   │   ├── planpairing.go             # PlanBOLPairing, PairingStatus enum
│   │   ├── driver.go                  # Driver, HomeWarehouse, HOSState, ActiveBOL reference
│   │   ├── assignment.go              # DriverBOLAssignment — driver + BOL + equipment link
│   │   ├── equipment.go               # Equipment, MaintenanceRecord, BreakdownRecord
│   │   ├── truck.go                   # TruckInventoryState, Capacity constraints
│   │   ├── hos.go                     # HOSLimit (per state), HOSWindow
│   │   └── invoice.go                 # DeliveryConfirmation, InternalInvoice
│   ├── repository/                    # DB access — interface-driven for testability
│   │   ├── interfaces.go              # Repository contracts defined before implementations
│   │   ├── planbul_repo.go
│   │   ├── driver_repo.go
│   │   ├── assignment_repo.go
│   │   └── equipment_repo.go
│   ├── integrations/                  # Adapter to Switchyard (.NET) APIs only
│   │   ├── logistics_client.go        # POST /BillOfLading, process, replace-stop
│   │   └── inventory_client.go        # GET /Clothing, /PPE, /Tool by location
│   └── migrations/                    # golang-migrate SQL files
│       ├── 001_create_drivers.up.sql
│       ├── 001_create_drivers.down.sql
│       ├── 002_create_planbul_records.up.sql
│       ├── 002_create_planbul_records.down.sql
│       ├── 003_create_planbul_stops.up.sql
│       ├── 003_create_planbul_stops.down.sql
│       ├── 004_create_truck_snapshots.up.sql
│       ├── 004_create_truck_snapshots.down.sql
│       ├── 005_create_deadhead_pairings.up.sql
│       ├── 005_create_deadhead_pairings.down.sql
│       ├── 006_create_delivery_confirmations.up.sql
│       ├── 006_create_delivery_confirmations.down.sql
│       ├── 007_create_internal_invoices.up.sql
│       ├── 007_create_internal_invoices.down.sql
│       ├── 008_create_equipment.up.sql
│       ├── 008_create_equipment.down.sql
│       ├── 009_create_driver_bol_assignments.up.sql
│       ├── 009_create_driver_bol_assignments.down.sql
│       ├── 010_create_hos_limits.up.sql
│       └── 010_create_hos_limits.down.sql
├── web/
│   ├── templates/
│   │   ├── dispatch_board.html        # Dispatcher Kanban whiteboard view
│   │   └── driver_runsheet.html       # Driver view — stop list, current inventory state
│   └── static/                        # CSS, JS — mobile-responsive from day one
├── deploy/
│   ├── Dockerfile                     # Multi-stage: golang:alpine builder → alpine runtime
│   ├── docker-compose.yml             # Island network config — see Section 10
│   └── .github/
│       └── workflows/
│           └── ci.yml                 # go vet → go test → docker build → push → deploy
├── output/
│   └── invoices/                      # Default INVOICE_OUTPUT_PATH — gitignore contents
├── .env.example
├── go.mod
├── go.sum
└── CLAUDE.md                          # Claude Code instructions — append, do not replace
```

---

## 7. Data Model — Go Backend Entities (PostgreSQL)

All IDs are UUIDs. Foreign keys to the existing Switchyard (.NET) system store those IDs as
strings. The Go backend never joins across system boundaries.

### Driver
```sql
id                UUID        PRIMARY KEY
name              TEXT        NOT NULL
home_warehouse_id TEXT        NOT NULL  -- warehouse ID from .NET system
auth0_user_id     TEXT        NOT NULL  -- links to existing Auth0 user
license_state     TEXT        NOT NULL  -- primary state for HOS limit lookup
is_active         BOOLEAN     NOT NULL DEFAULT true
created_at        TIMESTAMPTZ NOT NULL DEFAULT now()
```

### HOSLimit (configuration per state)
```sql
id                UUID        PRIMARY KEY
state_code        TEXT        NOT NULL  -- e.g. 'IL', 'IN', 'WI'
daily_limit_hours NUMERIC     NOT NULL
weekly_limit_hours NUMERIC    NOT NULL
effective_from    DATE        NOT NULL
notes             TEXT                  -- state-specific exceptions or exemptions
```

### HOSWindow (per driver, per run)
```sql
id                UUID        PRIMARY KEY
driver_id         UUID        NOT NULL REFERENCES driver(id)
window_start      TIMESTAMPTZ NOT NULL
daily_hours_used  NUMERIC     NOT NULL DEFAULT 0
weekly_hours_used NUMERIC     NOT NULL DEFAULT 0
mandated_stop_at  TIMESTAMPTZ           -- set if driver hits mandated rest
eld_stop_ref      TEXT                  -- ELD logbook reference if accessible
```

### Equipment
```sql
id                UUID        PRIMARY KEY
unit_id           TEXT        NOT NULL UNIQUE  -- dispatcher-facing identifier (e.g. TRUCK-04)
equipment_type    TEXT        NOT NULL  -- truck | tractor
home_warehouse_id TEXT        NOT NULL
status            TEXT        NOT NULL  -- available | assigned | maintenance | breakdown
created_at        TIMESTAMPTZ NOT NULL DEFAULT now()
```

### MaintenanceRecord
```sql
id                UUID        PRIMARY KEY
equipment_id      UUID        NOT NULL REFERENCES equipment(id)
description       TEXT        NOT NULL
scheduled_at      TIMESTAMPTZ NOT NULL
estimated_return  TIMESTAMPTZ           -- dispatcher estimate for whiteboard display
completed_at      TIMESTAMPTZ
```

### BreakdownRecord
```sql
id                UUID        PRIMARY KEY
equipment_id      UUID        NOT NULL REFERENCES equipment(id)
breakdown_type    TEXT        NOT NULL  -- depot | roadside
location_desc     TEXT                  -- free text location (roadside only)
driver_id         UUID        REFERENCES driver(id)  -- set if load still attached
load_attached     BOOLEAN     NOT NULL DEFAULT false
reported_at       TIMESTAMPTZ NOT NULL DEFAULT now()
resolved_at       TIMESTAMPTZ
```

### DriverBOLAssignment
```sql
id                    UUID        PRIMARY KEY
driver_id             UUID        NOT NULL REFERENCES driver(id)
plan_bol_id           UUID        NOT NULL REFERENCES plan_bol_record(id)
equipment_id          UUID        NOT NULL REFERENCES equipment(id)
assigned_at           TIMESTAMPTZ NOT NULL DEFAULT now()
departed_at           TIMESTAMPTZ
fulfilled_at          TIMESTAMPTZ
deadhead_confirmed_at TIMESTAMPTZ
```

### PlanBOLRecord
```sql
id                UUID        PRIMARY KEY
driver_id         UUID        NOT NULL REFERENCES driver(id)
originating_wh_id TEXT        NOT NULL
status            TEXT        NOT NULL  -- draft | validated | submitted | fulfilled
created_at        TIMESTAMPTZ NOT NULL DEFAULT now()
submitted_at      TIMESTAMPTZ
fulfilled_at      TIMESTAMPTZ
```

### PlanBOLStop
```sql
id                UUID        PRIMARY KEY
plan_bol_id       UUID        NOT NULL REFERENCES plan_bol_record(id)
sequence          INTEGER     NOT NULL
location_id       TEXT        NOT NULL
stop_type         TEXT        NOT NULL  -- warehouse | store | return_depot (future)
is_processed      BOOLEAN     NOT NULL DEFAULT false
processed_at      TIMESTAMPTZ
driver_log_ref    TEXT                  -- digital logbook entry reference
```

### TruckInventorySnapshot
```sql
id                  UUID        PRIMARY KEY
plan_bol_id         UUID        NOT NULL REFERENCES plan_bol_record(id)
plan_bol_stop_id    UUID        NOT NULL REFERENCES plan_bol_stop(id)
sku_id              TEXT        NOT NULL
quantity_loaded     INTEGER     NOT NULL
quantity_remaining  INTEGER     NOT NULL
snapshot_at         TIMESTAMPTZ NOT NULL DEFAULT now()
```

### PlanBOLPairing
```sql
id                  UUID        PRIMARY KEY
active_bol_id       UUID        NOT NULL REFERENCES plan_bol_record(id)
deadhead_bol_id     UUID        NOT NULL REFERENCES plan_bol_record(id)
paired_at           TIMESTAMPTZ NOT NULL DEFAULT now()
earliest_valid_at   TIMESTAMPTZ NOT NULL  -- active BOL estimated fulfillment minus 4 hours
origin_warehouse    TEXT        NOT NULL
status              TEXT        NOT NULL  -- proposed | confirmed | cancelled
```

### DeliveryConfirmation
```sql
id                  UUID        PRIMARY KEY
plan_bol_stop_id    UUID        NOT NULL REFERENCES plan_bol_stop(id)
driver_id           UUID        NOT NULL REFERENCES driver(id)
confirmed_at        TIMESTAMPTZ NOT NULL DEFAULT now()
invoice_id          UUID        REFERENCES internal_invoice(id)
```

### InternalInvoice
```sql
id                  UUID        PRIMARY KEY
store_id            TEXT        NOT NULL
plan_bol_id         UUID        NOT NULL REFERENCES plan_bol_record(id)
line_items          JSONB       NOT NULL  -- [{sku_id, qty_delivered, unit_ref}]
output_path         TEXT        NOT NULL
generated_at        TIMESTAMPTZ NOT NULL DEFAULT now()
```

---

## 8. API Endpoints (Go Backend)

All `/api/*` routes return JSON. Routes under `/` return server-rendered HTML.

### Event Handler
```
POST   /api/events                       Receive workflow event, route to service
```

### PlanBOL
```
POST   /api/plan-bol                     Create and begin constraint resolution
GET    /api/plan-bol/:id                 Get PlanBOLRecord with full stop sequence
POST   /api/plan-bol/:id/validate        Run full constraint pass, return violations
POST   /api/plan-bol/:id/submit          Submit validated plan to Switchyard (.NET)
GET    /api/plan-bol/:id/truck-state     Truck inventory state snapshot at each stop
```

### Driver
```
GET    /api/driver                       All drivers with current HOS state
GET    /api/driver/:id/runsheet          Current run — stops + live inventory state
POST   /api/driver/:id/stop/:stopId/log  Log stop completion, notify dispatch
GET    /api/driver/:id/active-bol        Current active PlanBOLRecord for this driver
GET    /api/driver/:id/hos               Current HOS window state
```

### Assignment
```
POST   /api/assignment                   Create driver-BOL-equipment assignment
GET    /api/assignment/:id               Get assignment with all linked entities
PATCH  /api/assignment/:id/depart        Mark departed (moves to In Transit)
PATCH  /api/assignment/:id/fulfill       Mark fulfilled (starts dead-head timer)
PATCH  /api/assignment/:id/deadhead      Confirm dead-head return (clears from board)
```

### Equipment
```
GET    /api/equipment                    All equipment with current status
POST   /api/equipment                    Register new truck or tractor
PATCH  /api/equipment/:id/maintenance    Report scheduled maintenance
PATCH  /api/equipment/:id/breakdown      Report breakdown (depot or roadside)
PATCH  /api/equipment/:id/resolve        Resolve maintenance or breakdown
```

### Dead-Head
```
GET    /api/deadhead/eligible            Eligible return BOLs given location + time
POST   /api/deadhead/pair                Pair two BOLs (enforces 4-hour rule)
DELETE /api/deadhead/:pairingId          Cancel a confirmed pairing
```

### Dispatch
```
GET    /api/dispatch/board               Full Kanban board state — all columns
GET    /api/dispatch/alerts              HOS warnings, breakdown alerts, expiring timers
```

### Internal Invoice
```
GET    /api/invoice/:id                  Get invoice record
GET    /api/invoice/store/:storeId       All invoices for a given store
```

### HTML Views
```
GET    /                                 Dispatch Kanban whiteboard
GET    /driver/:id                       Driver run sheet
```

---

## 9. Auth0 — Permissions

Uses the existing Auth0 tenant and M2M application. No new tenant required.
The Go event handler holds the single M2M token for all .NET API calls.

```
Existing permissions (do not modify):
  read:inventory | read:bol | create:bol | modify:bol | manage:users

Add for v1.1:
  fulfill:bol        — stop logging, truck state updates, delivery confirmation
  read:runsheet      — driver's own run sheet (scoped to active PlanBOLRecord only)
  manage:equipment   — equipment CRUD, maintenance and breakdown reporting
  manage:drivers     — driver HOS management, assignment creation
```

Auth0 M2M token limit on free tier: confirm available slots before adding new
applications. Current usage: one M2M for Scalar UI (existing). Go event handler
requires one additional M2M — confirm free tier allows this before implementation.

---

## 10. Deployment — Island Architecture

```
                    ┌──────────────────────┐
                    │    Load Balancer      │
                    │   (nginx / Traefik)   │
                    └───────────┬──────────┘
                                │
          ┌─────────────────────┼─────────────────────┐
          │                     │                     │
   ┌──────▼──────┐       ┌──────▼──────┐      ┌──────▼──────┐
   │  Go Backend │       │  Go Backend │      │  Go Backend │
   │  Island 1   │       │  Island 2   │      │  Islands 3-6│
   └──────┬──────┘       └──────┬──────┘      └──────┬──────┘
          │                     │                     │
          └─────────────────────┼─────────────────────┘
                                │
                    ┌───────────▼──────────┐
                    │     PostgreSQL        │
                    │   (Go backend DB)     │
                    └───────────┬──────────┘
                                │
                    ┌───────────▼──────────┐
                    │  Switchyard (.NET)    │
                    │  APIs                 │
                    │  Inventory  port 7000 │
                    │  Logistics  port 7001 │
                    └──────────────────────┘

Frontend (shared across all backend islands):
                    ┌──────────────────────┐
                    │   WarehouseSalesUI    │
                    │  (existing React/TS)  │
                    │   Static host / CDN   │
                    └──────────────────────┘
```

---

## 11. Environment Variables

```bash
# Required
DATABASE_URL=postgres://user:pass@host:5432/switchyard_go
LOGISTICS_API_URL=https://localhost:7001
INVENTORY_API_URL=https://localhost:7000
AUTH0_AUTHORITY=https://your-tenant.auth0.com/
AUTH0_AUDIENCE=your-api-audience

# Invoice output — swap value only, no code change needed
INVOICE_OUTPUT_PATH=./output/invoices

# HOS configuration
HOS_WARNING_THRESHOLD_HOURS=1.5        # Hours remaining before whiteboard warning color
DEADHEAD_WINDOW_HOURS=4                # Minimum dead-head pairing lead time (hard constraint)
DEADHEAD_SEARCH_WINDOW_HOURS=2         # Timer duration after delivery before dispatcher alert

# Email notifications
SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_USER=notifications@switchyard.com
SMTP_PASS=your-smtp-password
NOTIFY_DISPATCH_EMAIL=dispatch@switchyard.com

# Optional
PORT=8080
LOG_LEVEL=info
```

---

## 12. Technology Stack

| Concern | Choice | Rationale |
|---|---|---|
| Language | Go (open-source, BSD license) | Concurrency, low memory, compiled binary, event handler pattern |
| Router | `chi` | Lightweight, idiomatic, composable middleware |
| DB queries | `sqlc` | Type-safe, generated from SQL, no ORM magic |
| Migrations | `golang-migrate` | Versioned schema |
| Email | `gomail` or `net/smtp` | Lightweight, no external service dependency |
| Templates | `html/template` (stdlib) | No dependency, peelable when frontend separates |
| Config | `viper` + env vars | 12-factor app |
| Testing | `go test` + `testify` | Standard library, table-driven tests |
| Container | Docker multi-stage | Alpine runtime image, minimal footprint |
| CI/CD | GitHub Actions | Consistent with existing repo |

---

## 13. Integration Points with Switchyard (.NET)

The `/internal/integrations` adapter layer is the **only** place in the Go backend that
knows about the existing system's API contracts. No other package may import or call
the existing system directly. This boundary is enforced by code review.

```
Go backend calls existing system for:

  Inventory reads (before constraint resolution):
    GET  :7000/api/Clothing/location/:locationId
    GET  :7000/api/PPE/location/:locationId
    GET  :7000/api/Tool/location/:locationId

  BOL commitment (on PlanBOLRecord submit):
    POST :7001/api/BillOfLading

  Stop fulfillment tracking (after commitment):
    POST :7001/api/BillOfLading/:id/process/:locationId

  Emergency manual override (dispatch only, after commitment):
    POST :7001/api/BillOfLading/:id/replace-stop
```

---

## 14. Wanted Features Carried Forward

Features from the existing Switchyard (.NET) README that have architectural impact on
the Go backend and should be kept visible during sprint planning.

- **BOL status history** — `PlanBOLRecord.status` transitions provide this for the planning
  phase. Pairs with the existing system's audit trail when that feature is built.
- **Analytics endpoints** — `TruckInventorySnapshot` and `DeliveryConfirmation` records
  provide the raw data foundation. Analytics endpoints and reporting charts are scoped
  to v1.2 (see Section 16).
- **Returns** — not in scope for 1.1. Flagged for future sprint. Likely implemented as
  `stop_type: return_depot` on `PlanBOLStop` — the constraint solver already accommodates
  an additional stop type without structural changes.
- **Extract Data/ to class library** — the Go backend establishes clean domain separation
  from day one and can serve as the reference pattern during the existing system's refactor.

---

## 15. Handoff Instructions for Claude Code

Append the following block to the existing `CLAUDE.md` in the project root.
Do not replace the existing file.

```markdown
## Switchyard Go Architecture (v1.1)

Read `ARCHITECTURE.md` before writing any code for this project.

Priority reading order:
1. Section 2 — CUD authority boundary between Go and .NET. PlanBOL entities are
   Go's domain. Committed BOL stops are .NET's domain. These never cross.
2. Section 3 — The event handler is the single Go entry point. All workflow events
   route through it. The M2M token lives here and nowhere else.
3. Section 4 — Business rules are hard constraints. The empty truck rule, the 4-hour
   dead-head window, and state-level HOS limits are enforced at the service layer.
   They are not warnings. They are rejections.
4. Section 5 — Read the full Kanban board design before touching any whiteboard code.
   Column transitions are driven by assignment state, not by manual dispatcher input.
5. Section 13 — The integrations adapter is the only place that calls the .NET system.
   No exceptions.

Suggested implementation order:
  1. go.mod and project scaffold (Section 6)
  2. Domain models /internal/models — no dependencies, start here
  3. Repository interfaces /internal/repository/interfaces.go
  4. Migration files /internal/migrations — schema before any DB code
  5. HOSLimit seed data — state configurations before HOS service
  6. Integration adapters /internal/integrations — confirm existing API contracts first
  7. Event handler /internal/events — wire routing before any service code
  8. Service layer — route_planner and hos_service first (core constraints)
  9. Whiteboard service — depends on driver, equipment, and assignment services
  10. Notification service — last service, depends on all others for trigger points
  11. Handlers /internal/handlers — wire services to HTTP after services pass tests
  12. Web templates /web/templates — last, only if HTML views are in current sprint
```

---

## 16. v1.2 Candidate Features

Items moved from 1.1 scope due to sprint velocity or late addition after initial design.
Review at the start of the 1.2 planning session. Items may be re-prioritized or deferred
further depending on 1.1 completion and team direction.

### Operating Cost Tracking
Late request to establish data foundations for revenue vs. profit analytics.

- **Base operating cost:** $2.55/mile — applied to all active BOL runs
- **Tow rate (no full trailer):** $3.20/mile — roadside breakdown, tractor/truck without trailer
- **Flat-bed tow rate:** $3.80/mile — roadside breakdown requiring flat-bed recovery

Implementation notes:
- Cost records attach to `DriverBOLAssignment` or a new `RunCostRecord` entity
- Mileage source: derive from stop sequence coordinates or dispatcher-entered odometer
- Tow cost triggers from a resolved `BreakdownRecord` with `breakdown_type: roadside`
- Each facility is assumed to have an on-site maintenance/repair shop — tow cost covers
  the roadside-to-warehouse leg only

### Analytics Endpoints and Reporting
Depends on operating cost tracking being in place first.

- Revenue vs. profit charts per BOL, per driver, per warehouse, per time period
- Shipping figures aggregated from `DeliveryConfirmation` and `TruckInventorySnapshot`
- Operating cost vs. revenue overlay — requires cost records from above
- Data is already being captured in 1.1; endpoints and chart surfaces are the 1.2 addition

### Returns
- Implement `stop_type: return_depot` on `PlanBOLStop`
- Constraint solver accommodates this without structural changes — stop type enum is the
  only addition required
- Scoped out of 1.1 to keep sprint velocity focused on core dispatch workflow

### Theming and White-Label Client Overrides
- Add a "Theming" section which covers the default base color themes ("Industrial Cool" light
  and dark) and either instructions on how to provide the client configuration variables for
  theme overrides, or a link to the section README.md on what those instructions are. Will
  require an update to the section README.md if the link workflow is implemented.
- Client override scheme: `{client_dns_name}.switchyard.com` with SCSS-configurable color
  variables. Default operational domain is `@switchyard.com`.

### Switchyard Brand Assets
- Complete asset creation for the Switchyard logo and naming by this point. This work can be
  started in 1.1, but only as an additional feature if time allows. This includes:
  - Logo only
  - Name only
  - Logo and name together
  - "Powered by Switchyard" lockup
- All assets are required in both Light and Dark mode variants.
- Two-tone treatment on the "S" mark. Two modes only — no additional palette variants.

### README Refresh
- Rewrite the main `README.md` to reflect Switchyard naming and current 1.1 scope
- Add a `Future Features` section aligned to this document's v1.2 candidate list
- Scheduled for early in the week following 1.1 sprint kickoff
