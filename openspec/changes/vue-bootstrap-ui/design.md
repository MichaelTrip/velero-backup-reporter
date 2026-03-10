## Context

The Velero Backup Reporter currently serves HTML pages via Go `html/template` with vanilla JavaScript for client-side interactions. All UI assets (4 templates + 1 CSS file) are embedded in the Go binary via `go:embed`. The server package (`internal/server/`) renders templates directly in HTTP handlers. The goal is to replace this with a Vue 3 SPA served from the same Go binary, using Bootstrap 5 for styling.

## Goals / Non-Goals

**Goals:**
- Replace Go templates with a Vue 3 SPA that provides the same views: Dashboard, Backups List, Backup Detail
- Use Bootstrap 5 for responsive layout, components (cards, tables, badges, navbar), and styling
- Expose a JSON REST API from the Go backend for the Vue frontend to consume
- Keep the single-binary deployment model by embedding built frontend assets
- Maintain all existing functionality: filtering, sorting, PDF download, log viewing

**Non-Goals:**
- Real-time data via WebSockets (polling or manual refresh is sufficient)
- Server-side rendering (SSR) - pure client-side SPA is adequate for this admin tool
- Authentication/authorization (not currently present, not adding it now)
- Rewriting the email HTML templates (those stay as Go-rendered HTML)
- State management library (Vuex/Pinia) - the app is simple enough for component-local state

## Decisions

### 1. Vue 3 with Composition API + `<script setup>`

**Choice**: Vue 3 Composition API with `<script setup>` syntax.
**Rationale**: Simpler, more concise component code. Better TypeScript support. The app has few components, so the Composition API keeps things lightweight without needing Options API boilerplate.
**Alternatives considered**: React (heavier ecosystem, user specified Vue), Vue 2 (end-of-life).

### 2. Vite as build tool

**Choice**: Vite for frontend build tooling.
**Rationale**: Fast dev server with HMR, optimized production builds, first-class Vue support, minimal configuration needed.
**Alternatives considered**: Webpack (slower, more config), Vue CLI (deprecated in favor of Vite).

### 3. Bootstrap 5 via npm (no jQuery)

**Choice**: Bootstrap 5 installed via npm, imported in the Vue app. No jQuery dependency.
**Rationale**: Bootstrap 5 dropped jQuery requirement. CSS utilities and components work directly. For interactive Bootstrap components (dropdowns, modals), use `bootstrap` JS directly or lightweight Vue wrappers.
**Alternatives considered**: BootstrapVue (Vue 2 only, incompatible), BootstrapVueNext (still maturing).

### 4. REST API structure

**Choice**: JSON API endpoints under `/api/v1/` prefix.
- `GET /api/v1/dashboard` - summary stats + schedule statistics
- `GET /api/v1/backups` - list of all backups with metadata
- `GET /api/v1/backups/{name}` - single backup detail
- `GET /api/v1/backups/{name}/logs` - backup logs (plain text)
- `GET /api/v1/backups/{name}/pdf` - PDF download (binary)
- `GET /healthz` - health check (unchanged)

**Rationale**: Clean separation. The existing report/collector logic already produces structured data; handlers just need to marshal to JSON instead of passing to templates.
**Alternatives considered**: GraphQL (overkill for 3 views).

### 5. SPA routing with history mode fallback

**Choice**: Vue Router with HTML5 history mode. Go server serves the SPA `index.html` for any non-API, non-static route.
**Rationale**: Clean URLs without hash fragments. The Go server catches unmatched routes and serves `index.html`, letting Vue Router handle client-side routing.

### 6. Frontend project location

**Choice**: `web/frontend/` directory containing the Vue project. Build output goes to `web/dist/` which is embedded by `web/embed.go`.
**Rationale**: Keeps frontend code organized under existing `web/` directory. The `embed.go` file already handles embedding; it just needs to point to `dist/` instead of `templates/` and `static/`.

## Risks / Trade-offs

- **Increased build complexity** → Mitigated by a two-stage Dockerfile (Node build → Go build). Developers need Node.js installed for frontend development.
- **Larger binary size** → Bootstrap CSS + Vue runtime add ~200-300KB gzipped. Acceptable for a server-side tool.
- **SPA initial load** → Single page load fetches the app shell, then API calls for data. For an admin tool with few users, this is acceptable. No SSR needed.
- **Loss of server-side rendering for email** → Email templates are independent and unaffected.
- **Browser compatibility** → Vue 3 and Bootstrap 5 require modern browsers. Acceptable for a Kubernetes admin tool.
