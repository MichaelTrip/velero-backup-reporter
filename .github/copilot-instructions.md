# Copilot Instructions: Velero Backup Reporter

## Project Overview

**Velero Backup Reporter** is a lightweight Go application that provides a web UI and email notifications for Kubernetes Velero backup status. It reads Velero Backup and Schedule custom resources (CRDs), aggregates data, and exposes dashboards/reports.

**Key traits:**
- Single Go binary with embedded Vue 3 frontend (no external dependencies to serve)
- Kubernetes-native: runs in-cluster or out-of-cluster with kubeconfig
- Configuration via CLI flags, environment variables, or YAML config file
- Email notifications with cron-based scheduling
- Read-only; does not modify Velero resources

---

## Project Structure

```
cmd/velero-backup-reporter/     # Main entry point (Cobra CLI, signal handling)
internal/
  ├── collector/                # Kubernetes client & backup data collection
  ├── config/                   # Configuration loading & validation (Viper)
  ├── email/                    # SMTP sender & scheduled job runner
  ├── logs/                     # Log file parsing utilities
  ├── pdf/                      # PDF report generation (fpdf)
  ├── report/                   # HTML report templates & helpers
  └── server/                   # HTTP server (Chi router) + API handlers
web/
  ├── embed.go                  # Embeds frontend assets into binary
  └── frontend/                 # Vue 3 + Vite + PrimeVue UI
    ├── src/views/              # Current state-focused UI components
    ├── src/router/
    └── src/composables/        # Shared Vue logic
charts/                         # Helm chart for K8s deployment
deploy/                         # K8s manifests (alternative to Helm)
openspec/                       # Design proposals & specifications
```

---

## Build & Development

### Prerequisites
- Go 1.25.7+
- Node.js 18+ (frontend build)
- kubectl & kubeconfig (for Kubernetes development)

### Build Commands

| Command | Purpose |
|---------|---------|
| `make build` | Build frontend + backend (creates `velero-backup-reporter` binary) |
| `make frontend` | Build Vue/Vite frontend only; outputs to `web/dist/` |
| `make backend` | Build Go backend binary (CGO=0 for portability) |
| `make clean` | Remove built artifacts |
| `npm ci && npm run dev` (in `web/frontend/`) | Hot-reload frontend during development |
| `go build ./cmd/velero-backup-reporter/ && ./velero-backup-reporter ...` | Run backend standalone |

### Common Development Tasks

**Frontend development with hot reload:**
```bash
cd web/frontend
npm ci
npm run dev
# opens http://localhost:5173; backend must run separately
```

**Backend development (out-of-cluster):**
```bash
go build -o velero-backup-reporter ./cmd/velero-backup-reporter/
./velero-backup-reporter --kubeconfig ~/.kube/config --namespace velero --port 8080
# backend on :8080; frontend assets served from embedded web/dist/
```

**Test backend changes:**
```bash
go test ./internal/...
```

**Docker build:**
```bash
docker build -t velero-backup-reporter .
docker run -p 8080:8080 -v ~/.kube/config:/kubeconfig:ro -e KUBECONFIG=/kubeconfig velero-backup-reporter
```

---

## Key Components & Patterns

### Configuration Management
- **File:** [internal/config/config.go](../internal/config/config.go)
- Uses Viper for multi-source config: YAML file → environment variables → CLI flags (priority order reversed)
- Flags bound in [cmd/main.go](../cmd/velero-backup-reporter/main.go); validation in `config.Validate()`
- SMTP & email config are under `cfg.Email` and `cfg.SMTP` structs
- Must call `cfg.Validate()` after loading to check required fields

### Kubernetes Integration
- **Files:** [internal/collector/client.go](../internal/collector/client.go), [internal/collector/collector.go](../internal/collector/collector.go)
- Uses `k8s.io/client-go` + `sigs.k8s.io/controller-runtime`
- Supports in-cluster auth (ServiceAccount) and kubeconfig-based auth
- Watches/lists Velero Backup and Schedule CRDs in configured namespace
- Data collection runs on a ticker (default `5m`); see `collector.Run(ctx)`

### Web Server & API
- **File:** [internal/server/server.go](../internal/server/server.go)
- Chi router (lightweight, composable)
- HTTP handlers serve Vue frontend assets (from `web/embed.go`) + API endpoints
- API returns JSON data from collector for UI consumption
- Email test endpoint at `POST /api/email/test` (if `--email-test-enabled`)

### Email & Scheduling
- **Files:** [internal/email/email.go](../internal/email/email.go), [internal/email/scheduler.go](../internal/email/scheduler.go)
- Email sender wraps `net/smtp` with TLS support
- Scheduler uses `github.com/robfig/cron/v3` to run jobs at cron times (e.g., `0 8 * * *` = 8 AM daily)
- Report generation calls [internal/report/report.go](../internal/report/report.go) (HTML templates) and [internal/pdf/pdf.go](../internal/pdf/pdf.go)

### Frontend (Vue 3 + PrimeVue)
- **Entry:** [web/frontend/src/main.js](../web/frontend/src/main.js)
- Router-based layout with views: [DashboardView](../web/frontend/src/views/DashboardView.vue), [BackupsListView](../web/frontend/src/views/BackupsListView.vue), [BackupDetailView](../web/frontend/src/views/BackupDetailView.vue)
- Composables in [src/composables/](../web/frontend/src/composables/) for shared logic (e.g., backup utilities)
- UI components from PrimeVue library (PrimeVue components + PrimeIcons)
- Vite bundler outputs to `web/dist/` → embedded via `web/embed.go`

---

## Conventions & Patterns

1. **Error Handling:** Wrap errors with context: `fmt.Errorf("action: %w", err)`
2. **Logging:** Use `log.Println` for INFO, `log.Printf("ERROR: ...")` for errors (see [cmd/main.go](../cmd/velero-backup-reporter/main.go))
3. **Testing:** `_test.go` files in same package; use `go test ./internal/...`
4. **Config Validation:** Always validate after loading; return descriptive errors
5. **Graceful Shutdown:** Use context cancellation + signal handling (see [cmd/main.go](../cmd/velero-backup-reporter/main.go))
6. **Vue Composables:** Encapsulate reusable logic in functions returning refs/computed
7. **Embedded Assets:** Frontend dist files are embedded; rebuild frontend for changes to show up in binary

---

## Common Pitfalls

1. **Frontend changes not appearing:** Frontend must be rebuilt with `make frontend` or `npm run build` before rebuilding the backend. Dev mode requires running Vite separately on `localhost:5173`.

2. **SMTP TLS flag logic:** Verify `--smtp-tls` handling in [internal/email/email.go](../internal/email/email.go); TLS is enabled by default. Check related issue in `openspec/` for context.

3. **Kubernetes auth:** Out-of-cluster requires valid kubeconfig; in-cluster uses ServiceAccount token from `/var/run/secrets/kubernetes.io/serviceaccount/`.

4. **Collector restart:** Changing collection interval requires app restart; collector runs as a background goroutine.

5. **Email test endpoint:** Only available if `--email-test-enabled=true`; useful for development but should be disabled in production.

---

## Documentation & Proposals

- **README:** See [README.md](../README.md) for feature overview, configuration reference, and Quick Start
- **Design Specs:** [openspec/](../openspec/) contains structured proposals for features (e.g., UI redesigns, dark mode, CICD pipeline)
- **Helm Chart:** [charts/velero-backup-reporter/](../charts/velero-backup-reporter/) for production K8s deployment
- **Deploy Manifests:** [deploy/manifests.yaml](../deploy/manifests.yaml) for quick testing

---

## Useful Commands for AI Agents

```bash
# Build & test
make build
go test ./internal/... -v
go fmt ./...

# Quick dev loop (requires 2 terminals)
# Terminal 1 (frontend):
cd web/frontend && npm run dev

# Terminal 2 (backend):
go run ./cmd/velero-backup-reporter/ --kubeconfig ~/.kube/config --namespace velero

# Linting/formatting
go vet ./...
gofmt -w .

# Docker quick test
docker build -t velero-backup-reporter . && docker run -p 8080:8080 velero-backup-reporter

# Check config validation
go run ./cmd/velero-backup-reporter/ --help | grep -A 100 "Flags:"
```

---

## Glossary

| Term | Meaning |
|------|---------|
| **Velero** | Open-source Kubernetes backup/restore tool; provides Backup & Schedule CRDs |
| **CRD** | Custom Resource Definition; Kubernetes extension for custom object types |
| **Collector** | Background goroutine that fetches Velero Backup/Schedule resources and aggregates data |
| **Chi** | HTTP router library used for request multiplexing |
| **Viper** | Config management library supporting file, env, and flag sources |
| **Cobra** | CLI framework for argument parsing and command structure |
| **PrimeVue** | Vue 3 UI component library with built-in styling and icons |
| **Vite** | Modern frontend build tool (fast dev server + optimized production build) |
| **cron** | Time-based job scheduling syntax (e.g., `0 8 * * *` = 8 AM daily) |
| **CRON expression** | String defining when a task runs (minute, hour, day, month, weekday) |
| **ServiceAccount** | Kubernetes identity for in-cluster applications (provides auth token) |

---

## Feedback & Improvements

This file is maintained as the authoritative source for Copilot guidance in this workspace. If you discover missing information, conflicting details, or outdated sections, please update this file directly.
