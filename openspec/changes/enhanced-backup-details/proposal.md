## Why

The backup detail view currently shows high-level metadata but lacks granular volume-level information. Operators need to see which specific volumes were backed up, their sizes, progress, and individual statuses to effectively troubleshoot backup issues and verify data protection coverage.

## What Changes

- Collect per-volume backup details by querying PodVolumeBackup CRDs associated with each backup
- Add volume backup information to the backup detail API response: volume name, pod, status, bytes backed up, uploader type
- Display a new "Volume Backups" section in the backup detail view showing each volume with its size and status
- Add additional backup metadata fields: hook execution status, backup format version, async item operation counts
- Show storage size context in the backups list (total bytes backed up)

## Capabilities

### New Capabilities
- `volume-backup-details`: Collect and display per-volume backup information from PodVolumeBackup CRDs, including volume name, pod reference, status, bytes transferred, and uploader type

### Modified Capabilities
- `rest-api`: Add volume backup data to the backup detail endpoint response and additional summary fields
- `vue-spa-frontend`: Add Volume Backups section to the backup detail view and enhanced status information

## Impact

- **Backend (`internal/collector/`)**: New function to list PodVolumeBackup CRDs filtered by backup name label. New struct for per-volume data. Additional RBAC permissions needed for `velero.io/podvolumebackups` (get, list).
- **Backend (`internal/server/`)**: Backup detail API response extended with volume backup array and additional metadata fields.
- **Frontend (`web/frontend/`)**: BackupDetailView updated with new Volume Backups card showing a table of individual volumes.
- **Deployment (`deploy/manifests.yaml`)**: ClusterRole updated to include `podvolumebackups` resource permissions.
- **No new Go dependencies required** — PodVolumeBackup types are already available in the velero v1 API package.
