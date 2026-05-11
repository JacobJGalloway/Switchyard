package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/JacobJGalloway/switchyard-go/internal/models"
	"github.com/JacobJGalloway/switchyard-go/internal/repository"
)

// DeadheadHandler manages dead-head BOL pairing and the 4-hour window constraint.
type DeadheadHandler struct {
	pairingRepo         repository.PairingRepository
	deadheadWindowHours float64 // DEADHEAD_WINDOW_HOURS env var (hard minimum lead time)
}

func NewDeadheadHandler(
	pairingRepo repository.PairingRepository,
	deadheadWindowHours float64,
) *DeadheadHandler {
	return &DeadheadHandler{
		pairingRepo:         pairingRepo,
		deadheadWindowHours: deadheadWindowHours,
	}
}

type createPairingRequest struct {
	ActiveBOLID              string    `json:"active_bol_id"`
	DeadheadBOLID            string    `json:"deadhead_bol_id"`
	EstimatedFulfillmentAt   time.Time `json:"estimated_fulfillment_at"`
	OriginWarehouse          string    `json:"origin_warehouse"`
}

// GetEligible handles GET /api/deadhead/eligible
// Query params: location (string), estimated_completion (RFC3339 time)
func (h *DeadheadHandler) GetEligible(w http.ResponseWriter, r *http.Request) {
	location := r.URL.Query().Get("location")
	if location == "" {
		writeError(w, http.StatusBadRequest, "location query param is required")
		return
	}
	rawTime := r.URL.Query().Get("estimated_completion")
	if rawTime == "" {
		writeError(w, http.StatusBadRequest, "estimated_completion query param is required")
		return
	}
	estCompletion, err := time.Parse(time.RFC3339, rawTime)
	if err != nil {
		writeError(w, http.StatusBadRequest, "estimated_completion must be RFC3339 format")
		return
	}

	pairings, err := h.pairingRepo.GetEligible(r.Context(), location, estCompletion)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch eligible pairings")
		return
	}
	if pairings == nil {
		pairings = []*models.PlanBOLPairing{}
	}
	writeJSON(w, http.StatusOK, pairings)
}

// Pair handles POST /api/deadhead/pair
// Enforces the 4-hour minimum lead time (ARCHITECTURE.md §4.2 hard constraint).
func (h *DeadheadHandler) Pair(w http.ResponseWriter, r *http.Request) {
	var req createPairingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	activeBOLID, err := uuid.Parse(req.ActiveBOLID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid active_bol_id")
		return
	}
	deadheadBOLID, err := uuid.Parse(req.DeadheadBOLID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid deadhead_bol_id")
		return
	}
	if req.OriginWarehouse == "" {
		writeError(w, http.StatusBadRequest, "origin_warehouse is required")
		return
	}
	if req.EstimatedFulfillmentAt.IsZero() {
		writeError(w, http.StatusBadRequest, "estimated_fulfillment_at is required")
		return
	}

	// Hard constraint: pairing must be arranged at least deadheadWindowHours before
	// the active BOL is projected to be fulfilled (ARCHITECTURE.md §4.2).
	window := time.Duration(h.deadheadWindowHours * float64(time.Hour))
	earliestValidAt := req.EstimatedFulfillmentAt.Add(-window)
	if !time.Now().Before(earliestValidAt) {
		writeError(w, http.StatusConflict,
			"dead-head pairing window has closed — must be arranged at least "+
				fmt.Sprintf("%.0f hours", h.deadheadWindowHours)+
				" before estimated BOL fulfillment")
		return
	}

	pairing := &models.PlanBOLPairing{
		ID:              uuid.New(),
		ActiveBOLID:     activeBOLID,
		DeadheadBOLID:   deadheadBOLID,
		PairedAt:        time.Now().UTC(),
		EarliestValidAt: earliestValidAt,
		OriginWarehouse: req.OriginWarehouse,
		Status:          models.PairingStatusProposed,
	}
	if err := h.pairingRepo.Create(r.Context(), pairing); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create pairing")
		return
	}
	writeJSON(w, http.StatusCreated, pairing)
}

// Cancel handles DELETE /api/deadhead/:pairingId
func (h *DeadheadHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	id, ok := parseUUID(w, chi.URLParam(r, "pairingId"))
	if !ok {
		return
	}
	if err := h.pairingRepo.UpdateStatus(r.Context(), id, models.PairingStatusCancelled); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to cancel pairing")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
