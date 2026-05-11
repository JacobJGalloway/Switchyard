package integrations

import (
	"net/http"
	"time"
)

// TokenProvider returns the current M2M bearer token. The event handler holds
// the single Auth0 M2M token and provides a closure here — no other package
// ever sees or stores the token directly.
type TokenProvider func() string

func newHTTPClient() *http.Client {
	return &http.Client{Timeout: 30 * time.Second}
}
