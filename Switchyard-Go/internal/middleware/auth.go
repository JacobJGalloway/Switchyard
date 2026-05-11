package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v3"
	"github.com/auth0/go-jwt-middleware/v3/jwks"
	"github.com/auth0/go-jwt-middleware/v3/validator"
)

// CustomClaims extends the validated token with Auth0 API permissions.
// Permissions are assigned to roles in the Auth0 dashboard — no code change needed
// when role membership changes.
type CustomClaims struct {
	Permissions []string `json:"permissions"`
}

func (c *CustomClaims) Validate(_ context.Context) error { return nil }

// NewJWTMiddleware creates a CheckJWT middleware for the given Auth0 tenant.
// It verifies RS256 tokens against the tenant JWKS endpoint, validates issuer
// and audience, and populates CustomClaims (including permissions) into context.
// The default error handler returns RFC 6750 compliant JSON on 401.
func NewJWTMiddleware(domain, audience string) (func(http.Handler) http.Handler, error) {
	issuerURL, err := url.Parse("https://" + domain + "/")
	if err != nil {
		return nil, err
	}

	provider, err := jwks.NewCachingProvider(jwks.WithIssuerURL(issuerURL))
	if err != nil {
		return nil, err
	}

	jwtValidator, err := validator.New(
		validator.WithKeyFunc(provider.KeyFunc),
		validator.WithAlgorithm(validator.RS256),
		validator.WithIssuer(issuerURL.String()),
		validator.WithAudience(audience),
		validator.WithCustomClaims(func() *CustomClaims {
			return &CustomClaims{}
		}),
	)
	if err != nil {
		return nil, err
	}

	m, err := jwtmiddleware.New(jwtmiddleware.WithValidator(jwtValidator))
	if err != nil {
		return nil, err
	}

	return m.CheckJWT, nil
}

// RequirePermission returns a middleware that 403s if the validated JWT does not
// include the named permission in its permissions claim. Must sit downstream of
// NewJWTMiddleware — it assumes CheckJWT has already run and populated claims.
func RequirePermission(permission string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, err := jwtmiddleware.GetClaims[*validator.ValidatedClaims](r.Context())
			if err != nil {
				writeErr(w, http.StatusForbidden, "missing claims")
				return
			}
			custom, ok := claims.CustomClaims.(*CustomClaims)
			if !ok {
				writeErr(w, http.StatusForbidden, "invalid claims")
				return
			}
			for _, p := range custom.Permissions {
				if p == permission {
					next.ServeHTTP(w, r)
					return
				}
			}
			writeErr(w, http.StatusForbidden, "insufficient permissions")
		})
	}
}

func writeErr(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
