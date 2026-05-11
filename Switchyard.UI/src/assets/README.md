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

---

## Adding a new client

1. **Create the client folder** — `clients/{client-slug}/` where the slug is lowercase kebab-case (e.g. `clients/acme-logistics/`)
2. **Add asset files** — place any canonical-name assets the client provides into that folder. Only include what the client actually has — omitted assets fall back to Switchyard defaults automatically.
3. **Add imports to `src/utils/assetResolver.ts`** — one import per file at the top of the file
4. **Add a client block to `assetMap`** — only include the assets the client provides; any missing canonical name falls through to the Switchyard entry
5. **Add CSS variable overrides to `src/index.css`** — add a `:root[data-client="{client-slug}"]` block with the client's brand colors (mirrors the existing `[data-theme="dark"]` pattern)
6. **Update `.env.example`** — document the new client slug as a valid `VITE_CLIENT` value
7. **Favicon** — `public/favicon.svg` is outside the asset pipeline; replace manually for client deployments using the client's `logo-simple` source

### Notes
- `logo-name-only` is the key differentiator: if absent (null), the header renders a single integrated `logo-full-name` image. If present, it renders `logo-detail` + `logo-name-only` side by side.
- `logo-powered-by` is Switchyard-only — it is never overridden by a client. It is the attribution mark shown when a client override is active.
- `logo-simple` falls back to Switchyard by default; a client may provide their own version if they have a simplified mark.
- Set `VITE_CLIENT={client-slug}` and `VITE_CLIENT_NAME={Display Name}` in the deployment environment to activate the override.
