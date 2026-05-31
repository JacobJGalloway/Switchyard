package handlers

import (
	"context"
	"encoding/json"
	"html/template"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/JacobJGalloway/switchyard-go/internal/events"
	"github.com/JacobJGalloway/switchyard-go/internal/integrations"
	"github.com/JacobJGalloway/switchyard-go/internal/models"
	"github.com/JacobJGalloway/switchyard-go/internal/repository"
)

type driverHOSService interface {
	OnStopLogged(ctx context.Context, e events.StopLoggedPayload) error
}

// DriverHandler serves driver run sheets, HOS state, and stop logging.
type DriverHandler struct {
	driverRepo repository.DriverRepository
	bolRepo    repository.PlanBOLRepository
	hosRepo    repository.HOSRepository
	assignRepo repository.AssignmentRepository
	hosSvc     driverHOSService
	logistics  integrations.LogisticsClient
	tmpl       *template.Template // "driver_runsheet" template, parsed at startup
}

func NewDriverHandler(
	driverRepo repository.DriverRepository,
	bolRepo repository.PlanBOLRepository,
	hosRepo repository.HOSRepository,
	assignRepo repository.AssignmentRepository,
	hosSvc driverHOSService,
	logistics integrations.LogisticsClient,
	tmpl *template.Template,
) *DriverHandler {
	return &DriverHandler{
		driverRepo: driverRepo,
		bolRepo:    bolRepo,
		hosRepo:    hosRepo,
		assignRepo: assignRepo,
		hosSvc:     hosSvc,
		logistics:  logistics,
		tmpl:       tmpl,
	}
}

type driverSummary struct {
	*models.Driver
	HOSWindow *models.HOSWindow `json:"hos_window,omitempty"`
}

type runsheetResponse struct {
	Driver     *models.Driver              `json:"driver"`
	Assignment *models.DriverBOLAssignment `json:"assignment,omitempty"`
	PlanBOL    *models.PlanBOLRecord       `json:"plan_bol,omitempty"`
	Stops      []*models.PlanBOLStop       `json:"stops,omitempty"`
	Snapshots  []*models.TruckInventorySnapshot `json:"snapshots,omitempty"`
	HOSWindow  *models.HOSWindow           `json:"hos_window,omitempty"`
}

// GetAll handles GET /api/driver
// Returns all drivers with their current HOS window state.
func (h *DriverHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	drivers, err := h.driverRepo.GetAll(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch drivers")
		return
	}
	summaries := make([]driverSummary, 0, len(drivers))
	for _, d := range drivers {
		window, _ := h.hosRepo.GetWindowByDriver(r.Context(), d.ID)
		summaries = append(summaries, driverSummary{Driver: d, HOSWindow: window})
	}
	writeJSON(w, http.StatusOK, summaries)
}

// GetRunsheet handles GET /api/driver/:id/runsheet
// Returns the driver's current run: active assignment, stop sequence, and truck inventory state.
func (h *DriverHandler) GetRunsheet(w http.ResponseWriter, r *http.Request) {
	driverID, ok := parseUUID(w, chi.URLParam(r, "id"))
	if !ok {
		return
	}
	driver, err := h.driverRepo.GetByID(r.Context(), driverID)
	if err != nil {
		writeError(w, http.StatusNotFound, "driver not found")
		return
	}

	resp := runsheetResponse{Driver: driver}

	assignment, _ := h.assignRepo.GetActiveByDriver(r.Context(), driverID)
	if assignment != nil {
		resp.Assignment = assignment
		bol, _ := h.bolRepo.GetByID(r.Context(), assignment.PlanBOLID)
		if bol != nil {
			resp.PlanBOL = bol
			resp.Stops, _ = h.bolRepo.GetStops(r.Context(), bol.ID)
			resp.Snapshots, _ = h.bolRepo.GetSnapshots(r.Context(), bol.ID)
		}
	}

	resp.HOSWindow, _ = h.hosRepo.GetWindowByDriver(r.Context(), driverID)
	writeJSON(w, http.StatusOK, resp)
}

// LogStop handles POST /api/driver/:id/stop/:stopId/log
// Records stop completion, updates HOS, syncs to .NET if BOL is committed,
// and automatically marks the assignment fulfilled when the last stop is logged.
func (h *DriverHandler) LogStop(w http.ResponseWriter, r *http.Request) {
	driverID, ok := parseUUID(w, chi.URLParam(r, "id"))
	if !ok {
		return
	}
	stopID, ok := parseUUID(w, chi.URLParam(r, "stopId"))
	if !ok {
		return
	}

	var req struct {
		LoggedAt *time.Time `json:"logged_at"`
		MilesLeg *float64   `json:"miles_leg"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	loggedAt := time.Now().UTC()
	if req.LoggedAt != nil {
		loggedAt = req.LoggedAt.UTC()
	}

	stop, err := h.bolRepo.GetStopByID(r.Context(), stopID)
	if err != nil {
		writeError(w, http.StatusNotFound, "stop not found")
		return
	}

	// HOS enforcement — hard constraint. Stop log is rejected if a rule is violated.
	hosPayload := events.StopLoggedPayload{
		DriverID:      driverID,
		PlanBOLStopID: stopID,
		PlanBOLID:     stop.PlanBOLID,
		LocationID:    stop.LocationID,
		StopType:      string(stop.StopType),
		LoggedAt:      loggedAt,
	}
	if err := h.hosSvc.OnStopLogged(r.Context(), hosPayload); err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if err := h.bolRepo.MarkStopProcessed(r.Context(), stopID, loggedAt, req.MilesLeg); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to mark stop processed")
		return
	}

	// Forward to .NET if the BOL has been committed and this is a store delivery.
	bol, err := h.bolRepo.GetByID(r.Context(), stop.PlanBOLID)
	if err == nil && bol.SubmittedTransactionID != nil && stop.StopType == models.StopTypeStore {
		// Non-fatal: .NET sync failure does not block the driver logbook.
		_ = h.logistics.ProcessStop(r.Context(), *bol.SubmittedTransactionID, stop.LocationID)
	}

	// If all stops are now processed, total miles and auto-fulfill the assignment.
	allStops, err := h.bolRepo.GetStops(r.Context(), stop.PlanBOLID)
	if err == nil {
		allProcessed := true
		var totalMiles float64
		for _, s := range allStops {
			if !s.IsProcessed {
				allProcessed = false
				break
			}
			if s.MilesLeg != nil {
				totalMiles += *s.MilesLeg
			}
		}
		if allProcessed {
			if totalMiles > 0 {
				_ = h.bolRepo.SetMilesDriven(r.Context(), stop.PlanBOLID, totalMiles)
			}
			assignment, err := h.assignRepo.GetActiveByDriver(r.Context(), driverID)
			if err == nil && assignment != nil {
				now := time.Now().UTC()
				_ = h.assignRepo.MarkFulfilled(r.Context(), assignment.ID, now)
				_ = h.bolRepo.UpdateStatus(r.Context(), stop.PlanBOLID, models.PlanBOLStatusFulfilled)
			}
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetActiveBOL handles GET /api/driver/:id/active-bol
func (h *DriverHandler) GetActiveBOL(w http.ResponseWriter, r *http.Request) {
	driverID, ok := parseUUID(w, chi.URLParam(r, "id"))
	if !ok {
		return
	}
	assignment, err := h.assignRepo.GetActiveByDriver(r.Context(), driverID)
	if err != nil || assignment == nil {
		writeError(w, http.StatusNotFound, "no active BOL for this driver")
		return
	}
	bol, err := h.bolRepo.GetByID(r.Context(), assignment.PlanBOLID)
	if err != nil {
		writeError(w, http.StatusNotFound, "plan BOL not found")
		return
	}
	stops, _ := h.bolRepo.GetStops(r.Context(), bol.ID)
	writeJSON(w, http.StatusOK, planBOLResponse{PlanBOLRecord: bol, Stops: stops})
}

// GetRunsheetPage handles GET /driver/:id — server-rendered HTML run sheet.
func (h *DriverHandler) GetRunsheetPage(w http.ResponseWriter, r *http.Request) {
	driverID, ok := parseUUID(w, chi.URLParam(r, "id"))
	if !ok {
		return
	}
	driver, err := h.driverRepo.GetByID(r.Context(), driverID)
	if err != nil {
		http.Error(w, "driver not found", http.StatusNotFound)
		return
	}
	resp := runsheetResponse{Driver: driver}
	assignment, _ := h.assignRepo.GetActiveByDriver(r.Context(), driverID)
	if assignment != nil {
		resp.Assignment = assignment
		bol, _ := h.bolRepo.GetByID(r.Context(), assignment.PlanBOLID)
		if bol != nil {
			resp.PlanBOL = bol
			resp.Stops, _ = h.bolRepo.GetStops(r.Context(), bol.ID)
			resp.Snapshots, _ = h.bolRepo.GetSnapshots(r.Context(), bol.ID)
		}
	}
	resp.HOSWindow, _ = h.hosRepo.GetWindowByDriver(r.Context(), driverID)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.tmpl.ExecuteTemplate(w, "driver_runsheet", resp); err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
	}
}

// GetHOS handles GET /api/driver/:id/hos
func (h *DriverHandler) GetHOS(w http.ResponseWriter, r *http.Request) {
	driverID, ok := parseUUID(w, chi.URLParam(r, "id"))
	if !ok {
		return
	}
	window, err := h.hosRepo.GetWindowByDriver(r.Context(), driverID)
	if err != nil || window == nil {
		writeError(w, http.StatusNotFound, "no HOS window found for this driver")
		return
	}
	writeJSON(w, http.StatusOK, window)
}
