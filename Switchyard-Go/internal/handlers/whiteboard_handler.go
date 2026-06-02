package handlers

import (
	"context"
	"html/template"
	"net/http"

	"github.com/JacobJGalloway/switchyard-go/internal/services"
)

type boardService interface {
	GetBoardState(ctx context.Context) (*services.BoardState, error)
	GetAlerts(ctx context.Context) ([]*services.BoardAlert, error)
}

// WhiteboardHandler serves the dispatch Kanban board and its alert surface.
type WhiteboardHandler struct {
	svc          boardService
	tmpl         *template.Template // "dispatch_board" template, parsed at startup
	reactBaseURL string
}

func NewWhiteboardHandler(svc boardService, tmpl *template.Template, reactBaseURL string) *WhiteboardHandler {
	return &WhiteboardHandler{svc: svc, tmpl: tmpl, reactBaseURL: reactBaseURL}
}

// boardPageData wraps BoardState with view-layer concerns for the HTML template.
type boardPageData struct {
	*services.BoardState
	ReactBaseURL string
}

// GetBoard handles GET /api/dispatch/board
func (h *WhiteboardHandler) GetBoard(w http.ResponseWriter, r *http.Request) {
	board, err := h.svc.GetBoardState(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to assemble board state")
		return
	}
	writeJSON(w, http.StatusOK, board)
}

// GetBoardPage handles GET / — server-rendered HTML Kanban board.
func (h *WhiteboardHandler) GetBoardPage(w http.ResponseWriter, r *http.Request) {
	board, err := h.svc.GetBoardState(r.Context())
	if err != nil {
		http.Error(w, "failed to load board", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.tmpl.ExecuteTemplate(w, "dispatch_board", boardPageData{
		BoardState:   board,
		ReactBaseURL: h.reactBaseURL,
	}); err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
	}
}

// GetAlerts handles GET /api/dispatch/alerts
func (h *WhiteboardHandler) GetAlerts(w http.ResponseWriter, r *http.Request) {
	alerts, err := h.svc.GetAlerts(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch alerts")
		return
	}
	writeJSON(w, http.StatusOK, alerts)
}
