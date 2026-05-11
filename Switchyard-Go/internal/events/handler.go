package events

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/JacobJGalloway/switchyard-go/internal/integrations"
)

// Service interfaces — the event handler defines what it needs from each service.
// Concrete implementations in /internal/services/ satisfy these interfaces implicitly
// (Go interfaces are structural — no explicit declaration required).

type HOSService interface {
	OnStopLogged(ctx context.Context, e StopLoggedPayload) error
}

type WhiteboardService interface {
	OnAssignmentDeparted(ctx context.Context, e AssignmentPayload) error
	OnAssignmentFulfilled(ctx context.Context, e AssignmentPayload) error
	OnDeadheadConfirmed(ctx context.Context, e AssignmentPayload) error
	OnMandatedStop(ctx context.Context, e MandatedStopPayload) error
	OnEquipmentBreakdown(ctx context.Context, e EquipmentBreakdownPayload) error
	OnEquipmentResolved(ctx context.Context, e EquipmentResolvedPayload) error
}

type NotificationService interface {
	OnHOSLimitApproaching(ctx context.Context, e HOSAlertPayload) error
	OnHOSWeeklyLimitReached(ctx context.Context, e HOSAlertPayload) error
	OnBOLWorkflowCompleted(ctx context.Context, e BOLCompletedPayload) error
	OnDeadheadWindowExpiring(ctx context.Context, e DeadheadExpiryPayload) error
	OnRoadsideBreakdownWithLoad(ctx context.Context, e EquipmentBreakdownPayload) error
}

// Handler is the single entry point for all Go backend workflow events.
// It is the sole owner of the Auth0 M2M token — no other package holds it.
type Handler struct {
	hosService          HOSService
	whiteboardService   WhiteboardService
	notificationService NotificationService

	inventoryClient integrations.InventoryClient
	logisticsClient integrations.LogisticsClient

	httpClient *http.Client

	tokenMu     sync.RWMutex
	token       string
	tokenExpiry time.Time

	auth0Domain   string
	auth0ClientID string
	auth0Secret   string
	auth0Audience string
}

type Config struct {
	Auth0Domain   string
	Auth0ClientID string
	Auth0Secret   string
	Auth0Audience string
}

func NewHandler(
	cfg Config,
	hos HOSService,
	wb WhiteboardService,
	notify NotificationService,
	inv integrations.InventoryClient,
	log integrations.LogisticsClient,
) *Handler {
	return &Handler{
		hosService:          hos,
		whiteboardService:   wb,
		notificationService: notify,
		inventoryClient:     inv,
		logisticsClient:     log,
		httpClient:          &http.Client{Timeout: 10 * time.Second},
		auth0Domain:         cfg.Auth0Domain,
		auth0ClientID:       cfg.Auth0ClientID,
		auth0Secret:         cfg.Auth0Secret,
		auth0Audience:       cfg.Auth0Audience,
	}
}

// TokenProvider returns a closure for the integration clients.
// The closure checks expiry on every call and refreshes transparently.
func (h *Handler) TokenProvider() integrations.TokenProvider {
	return func() string {
		h.tokenMu.RLock()
		valid := time.Now().Before(h.tokenExpiry)
		t := h.token
		h.tokenMu.RUnlock()

		if valid {
			return t
		}

		if err := h.refreshToken(context.Background()); err != nil {
			// Return empty — the HTTP call will 401 and the error surfaces there.
			return ""
		}

		h.tokenMu.RLock()
		defer h.tokenMu.RUnlock()
		return h.token
	}
}

// Handle is registered as POST /api/events in cmd/main.go.
func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	var evt Event
	if err := json.NewDecoder(r.Body).Decode(&evt); err != nil {
		http.Error(w, "invalid event payload", http.StatusBadRequest)
		return
	}

	if err := route(r.Context(), h, evt); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

// InventoryClient exposes the inventory adapter to services wired through the handler.
func (h *Handler) InventoryClient() integrations.InventoryClient {
	return h.inventoryClient
}

// LogisticsClient exposes the logistics adapter to services wired through the handler.
func (h *Handler) LogisticsClient() integrations.LogisticsClient {
	return h.logisticsClient
}

type auth0TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

func (h *Handler) refreshToken(ctx context.Context) error {
	body, _ := json.Marshal(map[string]string{
		"client_id":     h.auth0ClientID,
		"client_secret": h.auth0Secret,
		"audience":      h.auth0Audience,
		"grant_type":    "client_credentials",
	})

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		"https://"+h.auth0Domain+"/oauth/token", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("building token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := h.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("fetching Auth0 token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Auth0 token endpoint returned %d", resp.StatusCode)
	}

	var tr auth0TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil {
		return fmt.Errorf("decoding token response: %w", err)
	}

	h.tokenMu.Lock()
	defer h.tokenMu.Unlock()
	h.token = tr.AccessToken
	// Refresh 60 seconds before actual expiry to avoid races at the boundary.
	h.tokenExpiry = time.Now().Add(time.Duration(tr.ExpiresIn-60) * time.Second)

	return nil
}
