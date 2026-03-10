## Context

Velero is the de facto standard for Kubernetes cluster backups. It stores backup metadata as Custom Resources (Backup, Schedule, BackupStorageLocation). Currently, operators must use `velero backup get` or query CRs directly to understand backup health. There is no consolidated reporting view or proactive notification system.

This application will run as a Deployment inside a Kubernetes cluster, using `client-go` with in-cluster or kubeconfig-based authentication to read Velero CRs. It will generate reports and serve them via a web UI, with optional email delivery.

## Goals / Non-Goals

**Goals:**
- Provide a single-binary Go application that collects Velero backup data and presents it in a web UI
- Support email delivery of reports via SMTP
- Make deployment simple via container image with Kubernetes manifests
- Keep the application read-only — it SHALL NOT modify any Velero resources
- Support both in-cluster and out-of-cluster (kubeconfig) operation

**Non-Goals:**
- Managing or triggering Velero backups (this is strictly a reporting tool)
- Supporting backup solutions other than Velero
- User authentication/authorization for the web UI (can be added later or handled by an ingress proxy)
- Persistent storage for historical report data (reports are generated from live CR state)
- Real-time streaming or WebSocket updates

## Decisions

### 1. Go with standard library HTTP server

**Decision**: Use Go's `net/http` with a lightweight router (e.g., `chi` or `gorilla/mux`) rather than a full framework.

**Rationale**: The application has simple routing needs. The standard library is well-tested, has no external dependencies, and Go developers are familiar with it. Chi adds minimal overhead while providing middleware support.

**Alternatives considered**:
- Gin/Echo: More opinionated, heavier dependencies for features we don't need
- Pure `net/http`: Lacks convenient routing patterns, middleware chaining

### 2. Embedded web assets with `embed` package

**Decision**: Use Go's `embed` package to bundle HTML templates and static assets into the binary.

**Rationale**: Single binary deployment is a key goal. Embedding assets avoids the need for a separate file system mount or init container. Templates can use Go's `html/template` package.

**Alternatives considered**:
- Separate frontend SPA (React/Vue): Over-engineered for a reporting dashboard; adds build complexity
- External static file serving: Requires volume mounts, complicates deployment

### 3. Kubernetes client-go for CR access

**Decision**: Use `client-go` with dynamic client or typed clients generated from Velero's API types.

**Rationale**: Velero publishes Go API types (`github.com/vmware-tanzu/velero/pkg/apis`). Using typed clients provides compile-time safety. The application can use in-cluster config when running as a pod or kubeconfig when running locally.

**Alternatives considered**:
- Direct REST API calls: Loses type safety, more boilerplate
- Controller-runtime: Designed for controllers with reconciliation loops; too heavy for read-only reporting

### 4. Configuration via environment variables and config file

**Decision**: Support configuration through environment variables (12-factor), an optional YAML config file, and CLI flags. Use a library like `spf13/viper` + `cobra` for configuration management.

**Rationale**: Environment variables work well in Kubernetes (ConfigMaps/Secrets). A config file is convenient for local development. CLI flags are useful for one-off overrides.

### 5. Go `net/smtp` or `jordan-wright/email` for email delivery

**Decision**: Use a lightweight SMTP library for sending HTML email reports.

**Rationale**: The email needs are simple — send formatted HTML reports to a list of recipients. No need for a full email service integration. Standard SMTP is universally supported.

### 6. Periodic report generation via internal scheduler

**Decision**: Use a simple ticker or cron-like scheduler (e.g., `robfig/cron`) within the application to periodically collect backup data and generate reports.

**Rationale**: Reports need periodic refresh to reflect current backup state. An internal scheduler avoids external dependencies (CronJobs). The refresh interval should be configurable.

## Risks / Trade-offs

- **[Velero API version changes]** → Pin to a specific Velero API version and document compatibility. The typed client approach means API changes require a rebuild. Mitigation: support Velero v1.x API which is stable.

- **[Large number of backups]** → If a cluster has thousands of historical backups, listing all CRs could be slow. Mitigation: support label selectors and time-based filtering; paginate API requests.

- **[SMTP misconfiguration]** → Email failures should not affect the web UI or report generation. Mitigation: email sending is async and errors are logged but do not block the main application.

- **[No persistence]** → Reports are generated from live state. If Velero CRs are garbage collected, historical data is lost. Mitigation: document this limitation; persistence can be added as a future enhancement.

- **[No auth on web UI]** → The web UI is unauthenticated. Mitigation: document that users should place it behind an authenticating proxy (OAuth2 Proxy, ingress auth) in production.
