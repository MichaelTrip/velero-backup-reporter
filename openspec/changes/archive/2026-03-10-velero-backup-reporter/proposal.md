## Why

There is no lightweight, self-hosted tool for generating and viewing backup reports from Velero, the popular Kubernetes backup solution. Cluster operators need visibility into backup status, success/failure trends, and schedule adherence without manually querying the Velero API or reading raw CR status fields. Additionally, teams need automated email notifications for backup reports to stay informed without checking a dashboard.

## What Changes

- New Go application that connects to a Kubernetes cluster and reads Velero Backup and Schedule custom resources
- Generates structured backup reports summarizing status, duration, items backed up, errors, and warnings
- Serves reports through a web UI with filtering and sorting capabilities
- Configurable SMTP-based email delivery of backup reports on a schedule or on-demand
- Helm chart or Kubernetes manifests for easy in-cluster deployment

## Capabilities

### New Capabilities
- `backup-collection`: Collecting and parsing Velero Backup and Schedule CRs from the Kubernetes API
- `report-generation`: Generating structured backup reports from collected data (status summaries, trends, per-backup details)
- `web-ui`: HTTP server serving a web interface to browse and filter backup reports
- `email-notifications`: SMTP-based email delivery of backup reports with configurable recipients and schedule
- `configuration`: Application configuration via environment variables, config file, and CLI flags

### Modified Capabilities
<!-- No existing capabilities to modify -->

## Impact

- **Dependencies**: Requires `client-go` for Kubernetes API access, Velero API types for CR parsing, an SMTP library for email, and a web framework or standard library HTTP server
- **APIs**: Exposes an HTTP API for the web UI and potentially a REST API for programmatic access to reports
- **Systems**: Requires RBAC permissions to read Velero CRs in the target cluster; optional SMTP server access for email notifications
- **Deployment**: Designed to run as a Deployment inside a Kubernetes cluster with a ServiceAccount that has read access to Velero resources
