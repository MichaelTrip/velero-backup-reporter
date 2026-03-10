## Why

The current UI is built with server-side Go templates and vanilla JavaScript. While functional, it lacks interactivity, modern component architecture, and a polished visual design. Migrating to Vue.js with Bootstrap will provide a responsive, component-driven SPA experience with richer user interactions (live filtering, dynamic updates, better navigation) and a professional look-and-feel out of the box.

## What Changes

- **BREAKING**: Replace server-side Go HTML templates with a Vue.js single-page application (SPA)
- Add a Vue 3 frontend application with Bootstrap 5 for styling and layout
- Convert the Go backend from serving HTML to serving a JSON REST API
- Embed the built Vue SPA static assets in the Go binary (replacing embedded templates)
- Redesign all three views (Dashboard, Backups List, Backup Detail) with Bootstrap components
- Add client-side routing via Vue Router for SPA navigation
- Improve table interactivity with reactive filtering, sorting, and search
- Add loading states, error handling, and responsive design via Bootstrap's grid system

## Capabilities

### New Capabilities
- `vue-spa-frontend`: Vue 3 SPA application with Bootstrap 5 styling, component architecture, and client-side routing
- `rest-api`: JSON REST API endpoints replacing the template-rendering handlers to serve data to the Vue frontend

### Modified Capabilities
<!-- No existing capabilities to modify - this is a greenfield frontend rewrite -->

## Impact

- **Backend (`internal/server/`)**: Handlers refactored from template rendering to JSON responses. Routes prefixed with `/api/v1/`. Static file serving updated for SPA assets with fallback routing.
- **Frontend (`web/`)**: Templates and CSS replaced by a Vue 3 project (`web/frontend/`) with build tooling (Vite). Built assets embedded via `go:embed`.
- **Build process**: Dockerfile updated to include Node.js build step for frontend. New `web/frontend/` directory with `package.json`, Vite config, Vue components.
- **Dependencies**: New npm dependencies (vue, vue-router, bootstrap, vite). No new Go dependencies required.
- **Deployment**: No changes to Kubernetes manifests. The binary remains a single self-contained artifact.
- **Email templates**: Unaffected (remain as Go-rendered HTML within the email package).
