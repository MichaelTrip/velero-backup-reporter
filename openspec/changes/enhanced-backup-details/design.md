## Context

The Velero Backup Reporter currently collects high-level backup metadata (phase, timestamps, item counts, volume snapshot counts) but does not provide per-volume granularity. Velero creates a `PodVolumeBackup` CR for each volume backed up via file-system backup (Kopia/Restic), containing the volume name, pod, status, and byte-level progress. This data is already available in the cluster but not surfaced in the reporter.

The collector (`internal/collector/`) queries Velero Backup and Schedule CRDs on a periodic interval. The server (`internal/server/`) exposes a JSON REST API consumed by the Vue SPA frontend.

## Goals / Non-Goals

**Goals:**
- Surface per-volume backup information (volume name, pod, node, status, bytes, uploader type) in the backup detail view
- Add hook execution status and backup item operation counts to backup details
- Maintain the existing periodic collection model — volume data is fetched on-demand per backup detail request, not cached globally

**Non-Goals:**
- Real-time progress tracking for in-progress volume backups (periodic refresh is sufficient)
- Collecting VolumeSnapshotLocation or BackupStorageLocation details (infrastructure-level, not per-backup)
- Modifying the backups list view with volume data (too much data for a list table)
- Adding backup size to the list view (PodVolumeBackup data requires per-backup queries, not efficient for list)

## Decisions

### 1. On-demand PodVolumeBackup fetching vs. periodic caching

**Choice**: Fetch PodVolumeBackup CRs on-demand when a backup detail is requested, not during periodic collection.
**Rationale**: PodVolumeBackups are only needed for the detail view. Caching all PodVolumeBackups globally would increase memory usage and API calls proportionally with backup count. The detail view is accessed infrequently (one backup at a time), making on-demand fetching acceptable.
**Alternatives considered**: Adding to periodic collection cache (unnecessary memory overhead for rarely-accessed data).

### 2. Label-based filtering for PodVolumeBackups

**Choice**: Filter PodVolumeBackups using the `velero.io/backup-name` label selector.
**Rationale**: Velero applies this label to all PodVolumeBackup CRs associated with a backup. This is the standard Velero convention and avoids listing all PodVolumeBackups then filtering client-side.

### 3. Additional Backup.Status fields

**Choice**: Extract `HookStatus` (attempted/failed counts) and `BackupItemOperations*` (attempted/completed/failed) from Backup.Status, plus `FormatVersion`.
**Rationale**: These are lightweight fields already present on the Backup CR. They provide useful operational insight (were hooks healthy? did async operations complete?) without requiring additional API calls.

### 4. VolumeBackupInfo struct in collector

**Choice**: Add a new `VolumeBackupInfo` struct in the collector package, and a new `ListVolumeBackups(ctx, backupName, namespace)` function that returns `[]VolumeBackupInfo`.
**Rationale**: Keeps volume data separate from the cached BackupInfo. The server calls this directly when handling detail requests, passing the existing kubeClient.

## Risks / Trade-offs

- **Additional RBAC permissions** → ClusterRole must add `podvolumebackups` to the list of permitted resources. This is a deployment change.
- **Detail view latency** → On-demand PodVolumeBackup queries add a small API call per detail view. Mitigated by the low frequency of detail view access and small data size.
- **PodVolumeBackups only cover file-system backups** → Native volume snapshots (via cloud provider) don't create PodVolumeBackup CRs. The existing `VolumeSnapshotsAttempted/Completed` counts cover those; the new Volume Backups table is for file-system (Kopia/Restic) backups specifically. The UI should label this clearly.
