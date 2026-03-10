## 1. Backend REST API

- [x] 1.1 Add JSON API route group under `/api/v1/` in `internal/server/server.go` with `Content-Type: application/json` middleware
- [x] 1.2 Implement `GET /api/v1/dashboard` handler returning summary stats and schedule statistics as JSON
- [x] 1.3 Implement `GET /api/v1/backups` handler returning all backups as a JSON array
- [x] 1.4 Implement `GET /api/v1/backups/{name}` handler returning full backup detail as JSON (with 404 for missing backups)
- [x] 1.5 Implement `GET /api/v1/backups/{name}/logs` handler returning plain text logs via the API route
- [x] 1.6 Implement `GET /api/v1/backups/{name}/pdf` handler returning PDF binary via the API route
- [x] 1.7 Add JSON error response helper returning `{"error": "<message>"}` with appropriate status codes
- [x] 1.8 Write unit tests for all API endpoints

## 2. Vue Frontend Setup

- [x] 2.1 Initialize Vue 3 project with Vite in `web/frontend/` (`npm create vite@latest` with Vue template)
- [x] 2.2 Install dependencies: `vue-router`, `bootstrap` (npm packages)
- [x] 2.3 Configure Vite: set `base` for production, configure dev server proxy to Go backend at `:8080`
- [x] 2.4 Set up main app entry point (`main.js`): import Bootstrap CSS, create Vue app with router
- [x] 2.5 Create `App.vue` with Bootstrap navbar (brand + Dashboard/Backups links), `<router-view>`, and footer

## 3. Vue Views and Components

- [x] 3.1 Create `DashboardView.vue`: fetch `/api/v1/dashboard`, render summary cards (Bootstrap cards with contextual colors) and schedule statistics table
- [x] 3.2 Create `BackupsListView.vue`: fetch `/api/v1/backups`, render Bootstrap table with status badges, implement reactive filtering (status + schedule dropdowns) and column sorting
- [x] 3.3 Create `BackupDetailView.vue`: fetch `/api/v1/backups/{name}`, render metadata/configuration/status/snapshots/labels sections in Bootstrap cards, add Back/PDF/Logs action buttons
- [x] 3.4 Implement log viewing: fetch logs on button click, display in `<pre>` block within a Bootstrap card
- [x] 3.5 Implement PDF download: trigger browser download from `/api/v1/backups/{name}/pdf`
- [x] 3.6 Create reusable `LoadingSpinner.vue` component (Bootstrap spinner) and `ErrorAlert.vue` component (Bootstrap danger alert)

## 4. Vue Router Configuration

- [x] 4.1 Configure Vue Router with routes: `/` → DashboardView, `/backups` → BackupsListView, `/backups/:name` → BackupDetailView
- [x] 4.2 Set router to HTML5 history mode

## 5. Embed and Serve SPA

- [x] 5.1 Update `web/embed.go` to embed `dist/` directory instead of `templates/` and `static/`
- [x] 5.2 Update Go server to serve embedded SPA static files (JS, CSS, assets) at root
- [x] 5.3 Add SPA fallback: serve `index.html` for any non-API, non-static route
- [x] 5.4 Remove old template rendering handlers and template-related code from `internal/server/server.go`
- [x] 5.5 Remove old template files (`web/templates/`) and stylesheet (`web/static/style.css`)

## 6. Build Pipeline

- [x] 6.1 Update Dockerfile with Node.js build stage: install deps, run `npm run build`, copy `dist/` to Go build stage
- [x] 6.2 Add `Makefile` or build script target to build frontend before Go binary (`cd web/frontend && npm ci && npm run build`)
- [x] 6.3 Verify the complete build produces a working single binary with embedded SPA
