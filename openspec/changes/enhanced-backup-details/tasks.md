## 1. Collector: Volume Backup Data

- [x] 1.1 Add `VolumeBackupInfo` struct to `internal/collector/collector.go` with fields: VolumeName, PodName, PodNamespace, NodeName, UploaderType, Phase, StartTimestamp, CompletionTimestamp, TotalBytes, BytesDone, SnapshotID
- [x] 1.2 Add `ListVolumeBackups(ctx, client, backupName, namespace)` function that queries PodVolumeBackup CRs using label selector `velero.io/backup-name=<name>` and returns `[]VolumeBackupInfo`
- [x] 1.3 Add additional fields to `BackupInfo` struct: HooksAttempted, HooksFailed, BackupItemOperationsAttempted, BackupItemOperationsCompleted, BackupItemOperationsFailed, FormatVersion
- [x] 1.4 Update `extractBackupInfo()` to populate the new BackupInfo fields from Backup.Status
- [x] 1.5 Write unit tests for `ListVolumeBackups` and the new `extractBackupInfo` fields

## 2. API: Enhanced Backup Detail Response

- [x] 2.1 Add `volumeBackupJSON` struct and `volumeBackups` field to `backupDetailJSON` in `internal/server/server.go`
- [x] 2.2 Add new BackupInfo fields (hooks, operations, formatVersion) to `backupDetailJSON`
- [x] 2.3 Update `handleAPIBackupDetail` to call `ListVolumeBackups` and include results in the response
- [x] 2.4 Write unit tests for the enhanced backup detail API response

## 3. Frontend: Volume Backups Display

- [x] 3.1 Add `formatBytes` utility function in `BackupDetailView.vue` to format byte values to human-readable units (B, KB, MB, GB, TB)
- [x] 3.2 Add "File System Volume Backups" card to `BackupDetailView.vue` with table showing Volume, Pod, Node, Status, Uploader, Size, Progress columns (only rendered when volumeBackups is non-empty)
- [x] 3.3 Add hook and async operations display to the Status card in `BackupDetailView.vue` (shown when hooksAttempted > 0 or operations attempted > 0)
- [x] 3.4 Build frontend and verify the updated detail view renders correctly

## 4. Deployment

- [x] 4.1 Update `deploy/manifests.yaml` ClusterRole to include `podvolumebackups` in the velero.io API group resources (get, list)
