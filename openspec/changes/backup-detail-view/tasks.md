## 1. Extend Backup Data Collection

- [x] 1.1 Add new fields to `collector.BackupInfo` struct: StorageLocation, TTL, IncludedNamespaces, ExcludedNamespaces, IncludedResources, ExcludedResources, Labels, Annotations, VolumeSnapshotsAttempted, VolumeSnapshotsCompleted, CSIVolumeSnapshotsAttempted, CSIVolumeSnapshotsCompleted, FailureReason, ValidationErrors
- [x] 1.2 Update `extractBackupInfo` to populate the new fields from the Velero Backup CR spec and status
- [x] 1.3 Write unit tests for the extended extraction logic

## 2. Backup Detail Page

- [x] 2.1 Add a `GetBackup(name string) *BackupInfo` method to the Collector that returns a single backup by name
- [x] 2.2 Create `web/templates/detail.html` template showing all backup fields organized into sections (metadata, configuration, status, volume snapshots)
- [x] 2.3 Add a per-page template set for the detail page in `server.New` (clone base + parse detail.html)
- [x] 2.4 Implement `handleBackupDetail` handler that looks up the backup by name from the chi URL param, renders the detail template, or returns 404
- [x] 2.5 Register routes `GET /backups/{name}` and ensure it doesn't conflict with existing `GET /backups`
- [x] 2.6 Write integration tests for the detail handler (found backup, not found)

## 3. Clickable Backup Names

- [x] 3.1 Update `web/templates/backups.html` to wrap backup names in `<a href="/backups/{{.Name}}">` links

## 4. Backup Logs

- [x] 4.1 Add Velero DownloadRequest API type to the kube client scheme in `collector.NewKubeClient`
- [x] 4.2 Implement `logs.FetchBackupLogs(ctx, kubeClient, backupName, namespace string) (string, error)` that creates a DownloadRequest CR, polls until processed (with timeout), fetches the signed URL, decompresses the gzip content, cleans up the CR, and returns the log text
- [x] 4.3 Implement `handleBackupLogs` handler at `GET /backups/{name}/logs` that calls the log fetcher and returns plain text, with error handling for not-found, non-terminal phase, and timeout cases
- [x] 4.4 Pass the kube client to the Server so the logs handler can create DownloadRequests
- [x] 4.5 Add a "View Logs" link/button to the detail page template (disabled/hidden for non-terminal backups)
- [x] 4.6 Update Kubernetes RBAC manifests to include `create`, `get`, `delete` on `downloadrequests.velero.io`
- [x] 4.7 Write unit tests for log retrieval logic (mock kube client scenarios: success, timeout, not found)

## 5. PDF Export

- [x] 5.1 Add `go-pdf/fpdf` dependency (`go get github.com/go-pdf/fpdf`)
- [x] 5.2 Implement `pdf.GenerateBackupReport(backup BackupInfo) ([]byte, error)` function that builds a PDF with backup metadata, configuration, and status sections
- [x] 5.3 Implement `handleBackupPDF` handler that generates the PDF and serves it with appropriate Content-Type and Content-Disposition headers
- [x] 5.4 Register route `GET /backups/{name}/pdf`
- [x] 5.5 Add a "Download PDF" button/link to the detail page template
- [x] 5.6 Write unit tests for PDF generation (verify non-empty output, correct content type)
