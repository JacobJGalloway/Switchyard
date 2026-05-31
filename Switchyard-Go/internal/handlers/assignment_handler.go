package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/JacobJGalloway/switchyard-go/internal/events"
	"github.com/JacobJGalloway/switchyard-go/internal/models"
	"github.com/JacobJGalloway/switchyard-go/internal/repository"
)

type assignHOSService interface {
	CanAssign(ctx context.Context, driverID uuid.UUID, estimatedRunHours float64, stateCode, cycleLabel string) error
}

type assignWhiteboardService interface {
	OnAssignmentDeparted(ctx context.Context, e events.AssignmentPayload) error
	OnAssignmentFulfilled(ctx context.Context, e events.AssignmentPayload) error
	OnDeadheadConfirmed(ctx context.Context, e events.AssignmentPayload) error
}

type assignNotificationService interface {
	OnBOLWorkflowCompleted(ctx context.Context, e events.BOLCompletedPayload) error
}

// AssignmentHandler manages driver-BOL-equipment assignment lifecycle.
type AssignmentHandler struct {
	assignRepo repository.AssignmentRepository
	driverRepo repository.DriverRepository
	bolRepo    repository.PlanBOLRepository
	equipRepo  repository.EquipmentRepository
	hosSvc     assignHOSService
	wbSvc      assignWhiteboardService
	notifySvc  assignNotificationService
}

func NewAssignmentHandler(
	assignRepo repository.AssignmentRepository,
	driverRepo repository.DriverRepository,
	bolRepo repository.PlanBOLRepository,
	equipRepo repository.EquipmentRepository,
	hosSvc assignHOSService,
	wbSvc assignWhiteboardService,
	notifySvc assignNotificationService,
) *AssignmentHandler {
	return &AssignmentHandler{
		assignRepo: assignRepo,
		driverRepo: driverRepo,
		bolRepo:    bolRepo,
		equipRepo:  equipRepo,
		hosSvc:     hosSvc,
		wbSvc:      wbSvc,
		notifySvc:  notifySvc,
	}
}

type createAssignmentRequest struct {
	DriverID          string   `json:"driver_id"`
	PlanBOLID         string   `json:"plan_bol_id"`
	EquipmentID       string   `json:"equipment_id"`
	EstimatedRunHours float64  `json:"estimated_run_hours"`
	StateCode         string   `json:"state_code"`
	CycleLabel        string   `json:"cycle_label"`
	BaseRatePerMile   *float64 `json:"base_rate_per_mile"`
}

type assignmentResponse struct {
	*models.DriverBOLAssignment
	Driver    *models.Driver        `json:"driver"`
	PlanBOL   *models.PlanBOLRecord `json:"plan_bol"`
	Equipment *models.Equipment     `json:"equipment"`
}

// Create handles POST /api/assignment
// Validates HOS eligibility, equipment availability, and BOL status before creating.
func (h *AssignmentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createAssignmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	driverID, err := uuid.Parse(req.DriverID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid driver_id")
		return
	}
	planBOLID, err := uuid.Parse(req.PlanBOLID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid plan_bol_id")
		return
	}
	equipmentID, err := uuid.Parse(req.EquipmentID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid equipment_id")
		return
	}
	if req.StateCode == "" || req.CycleLabel == "" {
		writeError(w, http.StatusBadRequest, "state_code and cycle_label are required")
		return
	}

	bol, err := h.bolRepo.GetByID(r.Context(), planBOLID)
	if err != nil {
		writeError(w, http.StatusNotFound, "plan BOL not found")
		return
	}
	if bol.Status != models.PlanBOLStatusValidated {
		writeError(w, http.StatusConflict, "plan BOL must be in validated status to assign a driver")
		return
	}

	equip, err := h.equipRepo.GetByID(r.Context(), equipmentID)
	if err != nil {
		writeError(w, http.StatusNotFound, "equipment not found")
		return
	}
	if equip.Status != models.EquipmentStatusAvailable {
		writeError(w, http.StatusConflict, "equipment is not available")
		return
	}

	// Hard HOS constraint — assignment is rejected if the driver cannot legally complete the run.
	if err := h.hosSvc.CanAssign(r.Context(), driverID, req.EstimatedRunHours, req.StateCode, req.CycleLabel); err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	baseRate := 3.20
	if req.BaseRatePerMile != nil {
		baseRate = *req.BaseRatePerMile
	}
	assignment := &models.DriverBOLAssignment{
		ID:              uuid.New(),
		DriverID:        driverID,
		PlanBOLID:       planBOLID,
		EquipmentID:     equipmentID,
		BaseRatePerMile: baseRate,
		AssignedAt:      time.Now().UTC(),
	}
	if err := h.assignRepo.Create(r.Context(), assignment); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create assignment")
		return
	}
	if err := h.equipRepo.UpdateStatus(r.Context(), equipmentID, models.EquipmentStatusAssigned); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update equipment status")
		return
	}

	writeJSON(w, http.StatusCreated, assignment)
}

// Get handles GET /api/assignment/:id
func (h *AssignmentHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, ok := parseUUID(w, chi.URLParam(r, "id"))
	if !ok {
		return
	}
	assignment, err := h.assignRepo.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "assignment not found")
		return
	}
	driver, _ := h.driverRepo.GetByID(r.Context(), assignment.DriverID)
	bol, _ := h.bolRepo.GetByID(r.Context(), assignment.PlanBOLID)
	equip, _ := h.equipRepo.GetByID(r.Context(), assignment.EquipmentID)
	writeJSON(w, http.StatusOK, assignmentResponse{
		DriverBOLAssignment: assignment,
		Driver:              driver,
		PlanBOL:             bol,
		Equipment:           equip,
	})
}

// Depart handles PATCH /api/assignment/:id/depart
// Marks the assignment as departed, moving the driver card to In Transit on the board.
func (h *AssignmentHandler) Depart(w http.ResponseWriter, r *http.Request) {
	id, ok := parseUUID(w, chi.URLParam(r, "id"))
	if !ok {
		return
	}
	assignment, err := h.assignRepo.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "assignment not found")
		return
	}
	if assignment.DepartedAt != nil {
		writeError(w, http.StatusConflict, "assignment has already departed")
		return
	}

	now := time.Now().UTC()
	if err := h.assignRepo.MarkDeparted(r.Context(), id, now); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to mark departure")
		return
	}
	if err := h.bolRepo.UpdateStatus(r.Context(), assignment.PlanBOLID, models.PlanBOLStatusSubmitted); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update BOL status")
		return
	}
	_ = h.wbSvc.OnAssignmentDeparted(r.Context(), events.AssignmentPayload{
		AssignmentID: id,
		DriverID:     assignment.DriverID,
		PlanBOLID:    assignment.PlanBOLID,
		EquipmentID:  assignment.EquipmentID,
	})
	w.WriteHeader(http.StatusNoContent)
}

// Fulfill handles PATCH /api/assignment/:id/fulfill
// Marks all stops confirmed, starts the dead-head timer, notifies dispatch.
func (h *AssignmentHandler) Fulfill(w http.ResponseWriter, r *http.Request) {
	id, ok := parseUUID(w, chi.URLParam(r, "id"))
	if !ok {
		return
	}
	assignment, err := h.assignRepo.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "assignment not found")
		return
	}
	if assignment.FulfilledAt != nil {
		writeError(w, http.StatusConflict, "assignment is already fulfilled")
		return
	}

	now := time.Now().UTC()
	if err := h.assignRepo.MarkFulfilled(r.Context(), id, now); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to mark fulfilled")
		return
	}
	if err := h.bolRepo.UpdateStatus(r.Context(), assignment.PlanBOLID, models.PlanBOLStatusFulfilled); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update BOL status")
		return
	}

	payload := events.AssignmentPayload{
		AssignmentID: id,
		DriverID:     assignment.DriverID,
		PlanBOLID:    assignment.PlanBOLID,
		EquipmentID:  assignment.EquipmentID,
	}
	_ = h.wbSvc.OnAssignmentFulfilled(r.Context(), payload)
	_ = h.notifySvc.OnBOLWorkflowCompleted(r.Context(), events.BOLCompletedPayload{
		AssignmentID: id,
		DriverID:     assignment.DriverID,
		PlanBOLID:    assignment.PlanBOLID,
	})
	w.WriteHeader(http.StatusNoContent)
}

// ConfirmDeadhead handles PATCH /api/assignment/:id/deadhead
// Confirms the dead-head return run, clearing the driver card from the board.
func (h *AssignmentHandler) ConfirmDeadhead(w http.ResponseWriter, r *http.Request) {
	id, ok := parseUUID(w, chi.URLParam(r, "id"))
	if !ok {
		return
	}
	assignment, err := h.assignRepo.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "assignment not found")
		return
	}
	if assignment.FulfilledAt == nil {
		writeError(w, http.StatusConflict, "assignment must be fulfilled before confirming dead-head")
		return
	}
	if assignment.DeadheadConfirmedAt != nil {
		writeError(w, http.StatusConflict, "dead-head is already confirmed")
		return
	}

	now := time.Now().UTC()
	if err := h.assignRepo.ConfirmDeadhead(r.Context(), id, now); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to confirm dead-head")
		return
	}
	_ = h.wbSvc.OnDeadheadConfirmed(r.Context(), events.AssignmentPayload{
		AssignmentID: id,
		DriverID:     assignment.DriverID,
		PlanBOLID:    assignment.PlanBOLID,
		EquipmentID:  assignment.EquipmentID,
	})
	w.WriteHeader(http.StatusNoContent)
}
