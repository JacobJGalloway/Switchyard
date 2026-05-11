package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/JacobJGalloway/switchyard-go/internal/services"
)

type stubBoardSvc struct{ err error }

func (s *stubBoardSvc) GetBoardState(_ context.Context) (*services.BoardState, error) {
	if s.err != nil {
		return nil, s.err
	}
	return &services.BoardState{}, nil
}

func (s *stubBoardSvc) GetAlerts(_ context.Context) ([]*services.BoardAlert, error) {
	if s.err != nil {
		return nil, s.err
	}
	return []*services.BoardAlert{}, nil
}

// --- GetBoard ---

func TestWhiteboard_GetBoard_ServiceError_Returns500(t *testing.T) {
	h := NewWhiteboardHandler(&stubBoardSvc{err: errors.New("db error")}, nil)
	req := httptest.NewRequest(http.MethodGet, "/api/dispatch/board", nil)
	rec := httptest.NewRecorder()
	h.GetBoard(rec, req)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestWhiteboard_GetBoard_Success_Returns200(t *testing.T) {
	h := NewWhiteboardHandler(&stubBoardSvc{}, nil)
	req := httptest.NewRequest(http.MethodGet, "/api/dispatch/board", nil)
	rec := httptest.NewRecorder()
	h.GetBoard(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

// --- GetAlerts ---

func TestWhiteboard_GetAlerts_ServiceError_Returns500(t *testing.T) {
	h := NewWhiteboardHandler(&stubBoardSvc{err: errors.New("db error")}, nil)
	req := httptest.NewRequest(http.MethodGet, "/api/dispatch/alerts", nil)
	rec := httptest.NewRecorder()
	h.GetAlerts(rec, req)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestWhiteboard_GetAlerts_Success_Returns200(t *testing.T) {
	h := NewWhiteboardHandler(&stubBoardSvc{}, nil)
	req := httptest.NewRequest(http.MethodGet, "/api/dispatch/alerts", nil)
	rec := httptest.NewRecorder()
	h.GetAlerts(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}
