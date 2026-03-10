## ADDED Requirements

### Requirement: Volume backups in backup detail response
The backup detail API endpoint (`GET /api/v1/backups/{name}`) SHALL include a `volumeBackups` array in the JSON response. Each entry SHALL contain: `volumeName`, `podName`, `podNamespace`, `nodeName`, `uploaderType`, `phase`, `startTimestamp`, `completionTimestamp`, `totalBytes`, `bytesDone`, and `snapshotId`.

#### Scenario: Backup detail with volume backups
- **WHEN** a client requests `GET /api/v1/backups/{name}` for a backup with PodVolumeBackup CRs
- **THEN** the response includes a `volumeBackups` array with one entry per PodVolumeBackup

#### Scenario: Backup detail without volume backups
- **WHEN** a client requests `GET /api/v1/backups/{name}` for a backup with no PodVolumeBackup CRs
- **THEN** the response includes an empty `volumeBackups` array `[]`

### Requirement: Additional status fields in backup detail response
The backup detail API endpoint SHALL include the following additional fields: `hooksAttempted` (int), `hooksFailed` (int), `backupItemOperationsAttempted` (int), `backupItemOperationsCompleted` (int), `backupItemOperationsFailed` (int), and `formatVersion` (string).

#### Scenario: Response includes hook and operation counts
- **WHEN** a client requests `GET /api/v1/backups/{name}`
- **THEN** the response JSON includes `hooksAttempted`, `hooksFailed`, `backupItemOperationsAttempted`, `backupItemOperationsCompleted`, `backupItemOperationsFailed`, and `formatVersion` fields
