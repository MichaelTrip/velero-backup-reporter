## 1. Project Setup

- [x] 1.1 Initialize Go module (`go mod init`) and set up project directory structure (cmd/, internal/, web/)
- [x] 1.2 Add core dependencies: client-go, Velero API types, cobra, viper, chi router
- [x] 1.3 Create main.go with cobra root command and basic CLI flag/env var setup

## 2. Configuration

- [x] 2.1 Implement config struct and loading from YAML file, env vars, and CLI flags with viper
- [x] 2.2 Add config validation (required fields, port range, SMTP completeness check)
- [x] 2.3 Write unit tests for config loading and validation

## 3. Kubernetes Client & Backup Collection

- [x] 3.1 Implement Kubernetes client initialization (in-cluster and kubeconfig modes)
- [x] 3.2 Implement Velero Backup CR listing and metadata extraction
- [x] 3.3 Implement Velero Schedule CR listing and backup-to-schedule association
- [x] 3.4 Implement periodic data collection with configurable interval using a ticker
- [x] 3.5 Write unit tests for backup metadata extraction and schedule association

## 4. Report Generation

- [x] 4.1 Define report data models (BackupReport, BackupSummary, ScheduleSummary, BackupDetail)
- [x] 4.2 Implement summary report generation (status counts, last success/failure timestamps)
- [x] 4.3 Implement per-schedule summary generation (success rate, last backup status)
- [x] 4.4 Write unit tests for report generation logic

## 5. Web UI

- [x] 5.1 Set up chi router with health endpoint (`/healthz`)
- [x] 5.2 Create HTML templates for dashboard page (summary stats, schedule table)
- [x] 5.3 Create HTML templates for backup list page (table with sorting and filtering)
- [x] 5.4 Add static CSS for styling the web UI
- [x] 5.5 Embed templates and static assets using Go `embed` package
- [x] 5.6 Implement HTTP handlers for dashboard and backup list pages
- [x] 5.7 Add client-side filtering by status and schedule, and column sorting
- [x] 5.8 Write integration tests for HTTP handlers

## 6. Email Notifications

- [x] 6.1 Implement SMTP client wrapper for sending HTML emails
- [x] 6.2 Create HTML email template for backup report
- [x] 6.3 Implement email scheduler using cron library with configurable schedule
- [x] 6.4 Add graceful error handling for SMTP failures (log and continue)
- [x] 6.5 Write unit tests for email formatting and scheduler

## 7. Integration & Deployment

- [x] 7.1 Wire all components together in main.go (config → client → collector → reporter → web server + email)
- [x] 7.2 Add graceful shutdown handling (context cancellation, HTTP server shutdown)
- [x] 7.3 Create Dockerfile for building the application container image
- [x] 7.4 Create Kubernetes manifests (Deployment, Service, ServiceAccount, RBAC) or Helm chart
- [x] 7.5 Add a README with configuration reference and deployment instructions
