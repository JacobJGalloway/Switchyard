package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/auth0/go-jwt-middleware/v3/core"
	"github.com/auth0/go-jwt-middleware/v3/validator"
	"github.com/stretchr/testify/assert"
)

// injectClaims is a test helper that sets ValidatedClaims in the request context,
// simulating what CheckJWT does after successful token verification.
func injectClaims(permissions []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := &validator.ValidatedClaims{
				CustomClaims: &CustomClaims{Permissions: permissions},
			}
			ctx := core.SetClaims(r.Context(), claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func okHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}

func TestRequirePermission_MissingClaims_Returns403(t *testing.T) {
	// No JWT middleware ran — no claims in context.
	handler := RequirePermission("manage:drivers")(okHandler())
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestRequirePermission_PermissionPresent_PassesThrough(t *testing.T) {
	handler := injectClaims([]string{"read:bol", "manage:drivers"})(
		RequirePermission("manage:drivers")(okHandler()),
	)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestRequirePermission_PermissionAbsent_Returns403(t *testing.T) {
	handler := injectClaims([]string{"read:bol"})(
		RequirePermission("manage:drivers")(okHandler()),
	)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestRequirePermission_EmptyPermissionList_Returns403(t *testing.T) {
	handler := injectClaims([]string{})(
		RequirePermission("manage:drivers")(okHandler()),
	)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestRequirePermission_ExactMatchOnly(t *testing.T) {
	// "manage:driver" (no 's') must not satisfy "manage:drivers".
	handler := injectClaims([]string{"manage:driver"})(
		RequirePermission("manage:drivers")(okHandler()),
	)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestCustomClaims_Validate_ReturnsNil(t *testing.T) {
	c := &CustomClaims{Permissions: []string{"read:bol"}}
	assert.NoError(t, c.Validate(nil))
}

// wrongClaims implements validator.CustomClaims but is NOT *CustomClaims.
// Injecting it simulates a misconfigured middleware chain.
type wrongClaims struct{}

func (w *wrongClaims) Validate(_ context.Context) error { return nil }

func TestRequirePermission_WrongClaimsType_Returns403(t *testing.T) {
	inject := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := &validator.ValidatedClaims{CustomClaims: &wrongClaims{}}
			ctx := core.SetClaims(r.Context(), claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
	handler := inject(RequirePermission("manage:drivers")(okHandler()))
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusForbidden, rec.Code)
}
