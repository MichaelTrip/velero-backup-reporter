## ADDED Requirements

### Requirement: Backup log retrieval
The system SHALL retrieve backup logs from Velero's object storage via the DownloadRequest CR mechanism.

#### Scenario: Fetch logs for a completed backup
- **WHEN** a user requests logs for a backup in a terminal phase (Completed, PartiallyFailed, Failed)
- **THEN** the system SHALL create a DownloadRequest CR, poll until processed, fetch the gzip-compressed logs from the signed URL, decompress them, and return the log content

#### Scenario: Backup not in terminal phase
- **WHEN** a user requests logs for a backup that is not in a terminal phase (e.g., InProgress, New)
- **THEN** the system SHALL return an error indicating logs are not yet available

#### Scenario: DownloadRequest timeout
- **WHEN** the Velero controller does not process the DownloadRequest within the timeout period
- **THEN** the system SHALL return an error indicating the log retrieval timed out

#### Scenario: Backup not found for logs
- **WHEN** a user requests logs for a backup name that does not exist
- **THEN** the system SHALL return an HTTP 404 response

### Requirement: Backup log endpoint
The system SHALL expose an HTTP endpoint for retrieving backup logs.

#### Scenario: Log endpoint
- **WHEN** a GET request is made to `/backups/{name}/logs`
- **THEN** the system SHALL return the backup log content as plain text

### Requirement: Backup log access from detail page
The backup detail page SHALL provide access to the backup's logs.

#### Scenario: View logs link
- **WHEN** the backup detail page is rendered for a backup in a terminal phase
- **THEN** the page SHALL include a "View Logs" link or button that navigates to the log content

#### Scenario: Logs unavailable
- **WHEN** the backup detail page is rendered for a backup not in a terminal phase
- **THEN** the "View Logs" link SHALL be disabled or hidden

### Requirement: DownloadRequest cleanup
The system SHALL clean up DownloadRequest CRs after use.

#### Scenario: Delete DownloadRequest after retrieval
- **WHEN** backup logs have been successfully retrieved or the request has timed out
- **THEN** the system SHALL delete the DownloadRequest CR to avoid accumulating stale resources

### Requirement: RBAC for log retrieval
The deployment manifests SHALL include RBAC permissions for DownloadRequest CRs.

#### Scenario: RBAC permissions
- **WHEN** the application is deployed with the provided Kubernetes manifests
- **THEN** the ServiceAccount SHALL have `create`, `get`, and `delete` permissions on `downloadrequests.velero.io`
