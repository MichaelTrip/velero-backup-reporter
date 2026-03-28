# Velero Backup Reporter

A lightweight Go application that generates backup reports from Velero and serves them in a web UI, with optional email notifications.

## Features

- **Dashboard** - Overview of backup status counts, last success/failure, and per-schedule statistics
- **Backup List** - Filterable, sortable table of all backups with details
- **Email Reports** - Scheduled HTML email reports via SMTP
- **Single Binary** - All web assets embedded, no external dependencies
- **Kubernetes Native** - Runs in-cluster with service account auth, or out-of-cluster with kubeconfig

## Quick Start

### In-Cluster

```bash
kubectl apply -f deploy/manifests.yaml
```

### Out-of-Cluster

```bash
go build -o velero-backup-reporter ./cmd/velero-backup-reporter/
./velero-backup-reporter --kubeconfig ~/.kube/config --namespace velero
```

Then open http://localhost:8080.

### Docker

```bash
docker build -t velero-backup-reporter .
docker run -p 8080:8080 -v ~/.kube/config:/kubeconfig:ro \
  -e KUBECONFIG=/kubeconfig velero-backup-reporter
```

## Configuration

Configuration can be provided via CLI flags, environment variables, or a YAML config file. CLI flags take highest precedence, then environment variables, then config file.

| Flag | Env Var | Default | Description |
|------|---------|---------|-------------|
| `--config` | - | - | Path to YAML config file |
| `--kubeconfig` | `KUBECONFIG` | - | Path to kubeconfig (uses in-cluster if not set) |
| `--namespace` | `NAMESPACE` | `velero` | Namespace to monitor for Velero resources |
| `--port` | `PORT` | `8080` | HTTP server port |
| `--collection-interval` | `COLLECTION_INTERVAL` | `5m` | Data collection interval |
| `--email-enabled` | `EMAIL_ENABLED` | `false` | Enable email notifications |
| `--email-schedule` | `EMAIL_SCHEDULE` | `0 8 * * *` | Cron schedule for email reports |
| `--email-details-window` | `EMAIL_DETAILS_WINDOW` | `24h` | Time window for backups shown in email report details |
| `--smtp-host` | `SMTP_HOST` | - | SMTP server host |
| `--smtp-port` | `SMTP_PORT` | `587` | SMTP server port |
| `--smtp-username` | `SMTP_USERNAME` | - | SMTP username |
| `--smtp-password` | `SMTP_PASSWORD` | - | SMTP password |
| `--smtp-from` | `SMTP_FROM` | - | Sender email address |
| `--smtp-to` | `SMTP_TO` | - | Recipient email addresses (comma-separated) |
| `--smtp-tls` | `SMTP_TLS` | `true` | Enable SMTP TLS |

### Example Config File

```yaml
namespace: velero
port: 8080
collection-interval: 10m
email-enabled: true
email-schedule: "0 8 * * 1-5"
email-details-window: 24h
smtp-host: smtp.example.com
smtp-port: 587
smtp-username: user@example.com
smtp-password: secret
smtp-from: velero-reports@example.com
smtp-to:
  - ops-team@example.com
```

## RBAC

The application requires read access to Velero Backup and Schedule custom resources. The included manifests create a ClusterRole with minimal permissions:

```yaml
rules:
  - apiGroups: ["velero.io"]
    resources: ["backups", "schedules"]
    verbs: ["get", "list", "watch"]
```

## Velero Compatibility

Tested with Velero v1.x API. The application is read-only and does not modify any Velero resources.
