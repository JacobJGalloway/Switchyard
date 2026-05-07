package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/JacobJGalloway/switchyard-go/internal/integrations"
	"github.com/JacobJGalloway/switchyard-go/internal/models"
	"github.com/JacobJGalloway/switchyard-go/internal/repository"
	"github.com/JacobJGalloway/switchyard-go/internal/services"
)

type planBOLService interface {
	PlanRoute(ctx context.Context, in services.PlanRouteInput) (*models.PlanBOLRecord, error)
	ValidatePlan(ctx context.Context, planBOLID uuid.UUID) ([]string, error)
}

// PlanBOLHandler manages the full PlanBOL lifecycle: create, validate, submit, inspect.
type PlanBOLHandler struct {
	svc       planBOLService
	bolRepo   repository.PlanBOLRepository
	logistics integrations.LogisticsClient
}

func NewPlanBOLHandler(
	svc planBOLService,
	bolRepo repository.PlanBOLRepository,
	logistics integrations.LogisticsClient,
) *PlanBOLHandler {
	return &PlanBOLHandler{svc: svc, bolRepo: bolRepo, logistics: logistics}
}

type createPlanBOLRequest struct {
	DriverID               string                `json:"driver_id"`
	OriginWarehouseID      string                `json:"origin_warehouse_id"`
	AdditionalWarehouseIDs []string              `json:"additional_warehouse_ids"`
	StoreStops             []stopRequestJSON     `json:"store_stops"`
	EstimatedRunHours      float64               `json:"estimated_run_hours"`
	StateCode              string                `json:"state_code"`
	CycleLabel             string                `json:"cycle_label"`
}

type stopRequestJSON struct {
	LocationID string         `json:"location_id"`
	Items      map[string]int `json:"items"`
}

type commitRequest struct {
	CustomerFirstName string `json:"customer_first_name"`
	CustomerLastName  string `json:"customer_last_name"`
	City              string `json:"city"`
	State             string `json:"state"`
}

type planBOLResponse struct {
	*models.PlanBOLRecord
	Stops []*models.PlanBOLStop `json:"stops"`
}

// Create handles POST /api/plan-bol
func (h *PlanBOLHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createPlanBOLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	driverID, err := uuid.Parse(req.DriverID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid driver_id")
		return
	}
	if req.OriginWarehouseID == "" {
		writeError(w, http.StatusBadRequest, "origin_warehouse_id is required")
		return
	}
	if len(req.StoreStops) == 0 {
		writeError(w, http.StatusBadRequest, "at least one store stop is required")
		return
	}

	stops := make([]services.StopRequest, len(req.StoreStops))
	for i, s := range req.StoreStops {
		stops[i] = services.StopRequest{
			LocationID: s.LocationID,
			Items:      s.Items,
		}
	}

	plan, err := h.svc.PlanRoute(r.Context(), services.PlanRouteInput{
		DriverID:               driverID,
		OriginWarehouseID:      req.OriginWarehouseID,
		AdditionalWarehouseIDs: req.AdditionalWarehouseIDs,
		StoreStops:             stops,
		EstimatedRunHours:      req.EstimatedRunHours,
		StateCode:              req.StateCode,
		CycleLabel:             req.CycleLabel,
	})
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, plan)
}

// Get handles GET /api/plan-bol/:id
func (h *PlanBOLHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, ok := parseUUID(w, chi.URLParam(r, "id"))
	if !ok {
		return
	}
	plan, err := h.bolRepo.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "plan BOL not found")
		return
	}
	stops, err := h.bolRepo.GetStops(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch stops")
		return
	}
	writeJSON(w, http.StatusOK, planBOLResponse{PlanBOLRecord: plan, Stops: stops})
}

// Validate handles POST /api/plan-bol/:id/validate
// Re-runs truck inventory constraints over the persisted stop sequence.
// Returns violations without advancing plan status — informational only.
func (h *PlanBOLHandler) Validate(w http.ResponseWriter, r *http.Request) {
	id, ok := parseUUID(w, chi.URLParam(r, "id"))
	if !ok {
		return
	}
	violations, err := h.svc.ValidatePlan(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "validation failed: "+err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"violations": violations})
}

// BeginPlanning handles POST /api/plan-bol/:id/begin-planning
// Claims a draft BOL for route planning (draft → plan-progress).
// Restricted to dispatchers and route planners.
func (h *PlanBOLHandler) BeginPlanning(w http.ResponseWriter, r *http.Request) {
	id, ok := parseUUID(w, chi.URLParam(r, "id"))
	if !ok {
		return
	}
	plan, err := h.bolRepo.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "plan BOL not found")
		return
	}
	if plan.Status != models.PlanBOLStatusDraft {
		writeError(w, http.StatusConflict, "plan BOL must be in draft status to begin planning")
		return
	}
	if err := h.bolRepo.UpdateStatus(r.Context(), id, models.PlanBOLStatusPlanProgress); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update plan status")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Commit handles POST /api/plan-bol/:id/commit
// Commits the finalized route plan to the .NET Logistics API and begins dock loading (plan-progress → loading).
// Stores the returned transaction ID for stop-processing calls.
// Restricted to dispatchers and route planners.
//
// The request body carries store contact details because the Go backend does not
// yet have a store master-data lookup. A future store integration adapter replaces this.
func (h *PlanBOLHandler) Commit(w http.ResponseWriter, r *http.Request) {
	id, ok := parseUUID(w, chi.URLParam(r, "id"))
	if !ok {
		return
	}
	var req commitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.CustomerFirstName == "" || req.CustomerLastName == "" || req.City == "" || req.State == "" {
		writeError(w, http.StatusBadRequest, "customer_first_name, customer_last_name, city, and state are required")
		return
	}

	plan, err := h.bolRepo.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "plan BOL not found")
		return
	}
	if plan.Status != models.PlanBOLStatusPlanProgress {
		writeError(w, http.StatusConflict, "plan BOL must be in plan-progress status to commit")
		return
	}

	stops, err := h.bolRepo.GetStops(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch stops")
		return
	}

	txID, err := h.logistics.CreateBOL(r.Context(), &integrations.CreateBOLRequest{
		CustomerFirstName: req.CustomerFirstName,
		CustomerLastName:  req.CustomerLastName,
		City:              req.City,
		State:             req.State,
		LineEntries:       buildLineEntries(stops),
	})
	if err != nil {
		writeError(w, http.StatusBadGateway, "logistics API error: "+err.Error())
		return
	}

	if err := h.bolRepo.SetSubmittedTransactionID(r.Context(), id, txID); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to store transaction ID")
		return
	}
	if err := h.bolRepo.UpdateStatus(r.Context(), id, models.PlanBOLStatusLoading); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update plan status")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"transaction_id": txID})
}

// MarkLoaded handles PATCH /api/plan-bol/:id/mark-loaded
// Confirms the trailer is loaded and ready for driver assignment (loading → validated).
func (h *PlanBOLHandler) MarkLoaded(w http.ResponseWriter, r *http.Request) {
	id, ok := parseUUID(w, chi.URLParam(r, "id"))
	if !ok {
		return
	}
	plan, err := h.bolRepo.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "plan BOL not found")
		return
	}
	if plan.Status != models.PlanBOLStatusLoading {
		writeError(w, http.StatusConflict, "plan BOL must be in loading status to mark as loaded")
		return
	}
	if err := h.bolRepo.UpdateStatus(r.Context(), id, models.PlanBOLStatusValidated); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update plan status")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// GetStatusHistory handles GET /api/plan-bol/:id/history
func (h *PlanBOLHandler) GetStatusHistory(w http.ResponseWriter, r *http.Request) {
	id, ok := parseUUID(w, chi.URLParam(r, "id"))
	if !ok {
		return
	}
	history, err := h.bolRepo.GetStatusHistory(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch status history")
		return
	}
	if history == nil {
		history = []*models.BOLStatusHistory{}
	}
	writeJSON(w, http.StatusOK, history)
}

// GetTruckState handles GET /api/plan-bol/:id/truck-state
func (h *PlanBOLHandler) GetTruckState(w http.ResponseWriter, r *http.Request) {
	id, ok := parseUUID(w, chi.URLParam(r, "id"))
	if !ok {
		return
	}
	snapshots, err := h.bolRepo.GetSnapshots(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch truck state")
		return
	}
	if snapshots == nil {
		snapshots = []*models.TruckInventorySnapshot{}
	}
	writeJSON(w, http.StatusOK, snapshots)
}

// buildLineEntries converts plan BOL stops into the .NET line entry format.
// Warehouse stops produce positive quantities (loading onto truck);
// store stops produce negative quantities (delivering off truck).
func buildLineEntries(stops []*models.PlanBOLStop) []integrations.LineEntryRequest {
	var entries []integrations.LineEntryRequest
	for _, stop := range stops {
		sign := 1
		if stop.StopType == models.StopTypeStore {
			sign = -1
		}
		for sku, qty := range stop.DeliveryItems {
			entries = append(entries, integrations.LineEntryRequest{
				LocationID: stop.LocationID,
				SKUMarker:  sku,
				Quantity:   sign * qty,
			})
		}
	}
	return entries
}
