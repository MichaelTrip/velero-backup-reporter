## ADDED Requirements

### Requirement: Collect PodVolumeBackup data per backup
The system SHALL provide a function to list PodVolumeBackup CRs associated with a given backup name, filtered by the `velero.io/backup-name` label. For each PodVolumeBackup, the system SHALL extract: volume name, pod name, pod namespace, node name, uploader type, phase, start timestamp, completion timestamp, total bytes, bytes done, and snapshot ID.

#### Scenario: Backup has file-system volume backups
- **WHEN** a backup detail is requested and PodVolumeBackup CRs exist with the label `velero.io/backup-name` matching the backup name
- **THEN** the system returns a list of volume backup records with volume name, pod, status, and byte counts

#### Scenario: Backup has no file-system volume backups
- **WHEN** a backup detail is requested and no PodVolumeBackup CRs exist for that backup name
- **THEN** the system returns an empty list of volume backups

#### Scenario: PodVolumeBackup query fails
- **WHEN** the Kubernetes API call to list PodVolumeBackup CRs fails
- **THEN** the system logs the error and returns an empty volume backup list (non-fatal)

### Requirement: Collect additional backup status fields
The system SHALL extract the following additional fields from the Backup CR status: hook status (hooks attempted, hooks failed), backup item operations (attempted, completed, failed), and format version.

#### Scenario: Backup has hook status
- **WHEN** a backup has hook execution data in its status
- **THEN** the collector extracts hooksAttempted and hooksFailed counts

#### Scenario: Backup has no hook status
- **WHEN** a backup has no hook status data (nil HookStatus)
- **THEN** the hook counts default to zero
