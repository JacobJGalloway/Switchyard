# Asset Naming Convention

All brand assets follow canonical names. Switchyard defaults live at this root.
Client overrides go in `clients/{client-slug}/` using the same names.

| Canonical name | Purpose |
|---|---|
| `logo-full-name-{light\|dark}` | Integrated lockup — nav header + login screen |
| `logo-name-only-{light\|dark}` | Text/wordmark only — used when client has separate icon and name assets; absent for Switchyard |
| `logo-detail-{light\|dark}` | Detailed standalone graphic — large display contexts |
| `logo-simple-{light\|dark}` | Simplified mark — small contexts, favicon default |
| `logo-powered-by-{light\|dark}` | Switchyard attribution — shown in footer when a client override is active |

Each name ships in three formats: `.svg` (UI), `.png`, `.jpg` (cross-platform/backend use).
