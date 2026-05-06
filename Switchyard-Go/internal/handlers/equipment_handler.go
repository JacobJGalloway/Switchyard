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

type equipNotificationService interface {
	OnRoadsideBreakdownWithLoad(ctx context.Context, e events.EquipmentBreakdownPayload) error
}

// EquipmentHandler manages truck and tractor lifecycle, maintenance, and breakdowns.
type EquipmentHandler struct {
	equipRepo repository.EquipmentRepository
	notifySvc equipNotificationService
}

func NewEquipmentHandler(
	equipRepo repository.EquipmentRepository,
	notifySvc equipNotificationService,
) *EquipmentHandler {
	return &EquipmentHandler{equipRepo: equipRepo, notifySvc: notifySvc}
}

type createEquipmentRequest struct {
	UnitID          string                `json:"unit_id"`
	EquipmentType   models.EquipmentType  `json:"equipment_type"`
	HomeWarehouseID string                `json:"home_warehouse_id"`
}

type reportMaintenanceRequest struct {
	Description     string     `json:"description"`
	ScheduledAt     time.Time  `json:"scheduled_at"`
	EstimatedReturn *time.Time `json:"estimated_return"`
}

type reportBreakdownRequest struct {
	BreakdownType models.BreakdownType `json:"breakdown_type"`
	LocationDesc  *string              `json:"location_desc"`
	DriverID      *string              `json:"driver_id"`
	LoadAttached  bool                 `json:"load_attached"`
}

// GetAll handles GET /api/equipment
func (h *EquipmentHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	equipment, err := h.equipRepo.GetAll(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch equipment")
		return
	}
	if equipment == nil {
		equipment = []*models.Equipment{}
	}
	writeJSON(w, http.StatusOK, equipment)
}

// Create handles POST /api/equipment
func (h *EquipmentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createEquipmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.UnitID == "" || req.HomeWarehouseID == "" {
		writeError(w, http.StatusBadRequest, "unit_id and home_warehouse_id are required")
		return
	}
	if req.EquipmentType != models.EquipmentTypeTruck && req.EquipmentType != models.EquipmentTypeTractor {
		writeError(w, http.StatusBadRequest, "equipment_type must be truck or tractor")
		return
	}

	e := &models.Equipment{
		ID:              uuid.New(),
		UnitID:          req.UnitID,
		EquipmentType:   req.EquipmentType,
		HomeWarehouseID: req.HomeWarehouseID,
		Status:          models.EquipmentStatusAvailable,
		CreatedAt:       time.Now().UTC(),
	}
	if err := h.equipRepo.Create(r.Context(), e); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create equipment")
		return
	}
	writeJSON(w, http.StatusCreated, e)
}

// ReportMaintenance handles PATCH /api/equipment/:id/maintenance
func (h *EquipmentHandler) ReportMaintenance(w http.ResponseWriter, r *http.Request) {
	id, ok := parseUUID(w, chi.URLParam(r, "id"))
	if !ok {
		return
	}
	var req reportMaintenanceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Description == "" {
		writeError(w, http.StatusBadRequest, "description is required")
		return
	}

	equip, err := h.equipRepo.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "equipment not found")
		return
	}
	if equip.Status == models.EquipmentStatusMaintenance || equip.Status == models.EquipmentStatusBreakdown {
		writeError(w, http.StatusConflict, "equipment is already in maintenance or breakdown")
		return
	}

	rec := &models.MaintenanceRecord{
		ID:              uuid.New(),
		EquipmentID:     id,
		Description:     req.Description,
		ScheduledAt:     req.ScheduledAt,
		EstimatedReturn: req.EstimatedReturn,
	}
	if err := h.equipRepo.CreateMaintenanceRecord(r.Context(), rec); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create maintenance record")
		return
	}
	if err := h.equipRepo.UpdateStatus(r.Context(), id, models.EquipmentStatusMaintenance); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update equipment status")
		return
	}
	writeJSON(w, http.StatusOK, rec)
}

// ReportBreakdown handles PATCH /api/equipment/:id/breakdown
func (h *EquipmentHandler) ReportBreakdown(w http.ResponseWriter, r *http.Request) {
	id, ok := parseUUID(w, chi.URLParam(r, "id"))
	if !ok {
		return
	}
	var req reportBreakdownRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.BreakdownType != models.BreakdownTypeDepot && req.BreakdownType != models.BreakdownTypeRoadside {
		writeError(w, http.StatusBadRequest, "breakdown_type must be depot or roadside")
		return
	}
	if req.BreakdownType == models.BreakdownTypeRoadside && req.LocationDesc == nil {
		writeError(w, http.StatusBadRequest, "location_desc is required for roadside breakdowns")
		return
	}

	var driverID *uuid.UUID
	if req.DriverID != nil {
		parsed, err := uuid.Parse(*req.DriverID)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid driver_id")
			return
		}
		driverID = &parsed
	}

	rec := &models.BreakdownRecord{
		ID:            uuid.New(),
		EquipmentID:   id,
		BreakdownType: req.BreakdownType,
		LocationDesc:  req.LocationDesc,
		DriverID:      driverID,
		LoadAttached:  req.LoadAttached,
		ReportedAt:    time.Now().UTC(),
	}
	if err := h.equipRepo.CreateBreakdownRecord(r.Context(), rec); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create breakdown record")
		return
	}
	if err := h.equipRepo.UpdateStatus(r.Context(), id, models.EquipmentStatusBreakdown); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update equipment status")
		return
	}

	// Roadside breakdown with load attached — critical, notify dispatcher immediately.
	if req.BreakdownType == models.BreakdownTypeRoadside && req.LoadAttached {
		_ = h.notifySvc.OnRoadsideBreakdownWithLoad(r.Context(), events.EquipmentBreakdownPayload{
			EquipmentID:   id,
			BreakdownID:   rec.ID,
			BreakdownType: string(req.BreakdownType),
			LocationDesc:  req.LocationDesc,
			DriverID:      driverID,
			LoadAttached:  req.LoadAttached,
			ReportedAt:    rec.ReportedAt,
		})
	}

	writeJSON(w, http.StatusOK, rec)
}

// Resolve handles PATCH /api/equipment/:id/resolve
// Resolves the active maintenance or breakdown record and returns equipment to available.
func (h *EquipmentHandler) Resolve(w http.ResponseWriter, r *http.Request) {
	id, ok := parseUUID(w, chi.URLParam(r, "id"))
	if !ok {
		return
	}
	equip, err := h.equipRepo.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "equipment not found")
		return
	}

	now := time.Now().UTC()
	switch equip.Status {
	case models.EquipmentStatusMaintenance:
		rec, err := h.equipRepo.GetActiveMaintenanceByEquipment(r.Context(), id)
		if err != nil || rec == nil {
			writeError(w, http.StatusNotFound, "no active maintenance record found")
			return
		}
		if err := h.equipRepo.ResolveMaintenanceRecord(r.Context(), rec.ID, now); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to resolve maintenance record")
			return
		}
	case models.EquipmentStatusBreakdown:
		rec, err := h.equipRepo.GetActiveBreakdownByEquipment(r.Context(), id)
		if err != nil || rec == nil {
			writeError(w, http.StatusNotFound, "no active breakdown record found")
			return
		}
		if err := h.equipRepo.ResolveBreakdownRecord(r.Context(), rec.ID, now); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to resolve breakdown record")
			return
		}
	default:
		writeError(w, http.StatusConflict, "equipment is not currently in maintenance or breakdown")
		return
	}

	if err := h.equipRepo.UpdateStatus(r.Context(), id, models.EquipmentStatusAvailable); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update equipment status")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
