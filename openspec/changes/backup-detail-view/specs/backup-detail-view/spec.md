## ADDED Requirements

### Requirement: Backup detail page
The system SHALL display a detail page for a single Velero Backup showing comprehensive information from the Backup CR.

#### Scenario: Navigate to backup detail
- **WHEN** a user navigates to `/backups/{name}`
- **THEN** the system SHALL display the full details of the backup with that name

#### Scenario: Backup not found
- **WHEN** a user navigates to `/backups/{name}` and no backup with that name exists
- **THEN** the system SHALL return an HTTP 404 response with a user-friendly message

### Requirement: Display backup metadata
The detail page SHALL display the backup's core metadata.

#### Scenario: Core metadata fields
- **WHEN** the backup detail page is rendered
- **THEN** the page SHALL display: backup name, namespace, status (phase), start time, completion time, duration, storage location, and TTL

#### Scenario: Labels and annotations
- **WHEN** the backup has labels or annotations
- **THEN** the page SHALL display them as key-value pairs

### Requirement: Display backup spec configuration
The detail page SHALL show how the backup was configured.

#### Scenario: Namespace and resource filters
- **WHEN** the backup spec includes included/excluded namespaces or resources
- **THEN** the page SHALL display these filter lists

#### Scenario: Schedule association
- **WHEN** the backup was created by a schedule
- **THEN** the page SHALL display the schedule name

### Requirement: Display backup status details
The detail page SHALL show detailed status information beyond the phase.

#### Scenario: Volume snapshot counts
- **WHEN** the backup has volume snapshot information
- **THEN** the page SHALL display snapshots attempted vs completed, and CSI snapshots attempted vs completed

#### Scenario: Warnings and errors
- **WHEN** the backup has warnings or errors
- **THEN** the page SHALL display the warning count and error count

#### Scenario: Failure reason
- **WHEN** the backup has a failure reason
- **THEN** the page SHALL display the failure reason prominently

#### Scenario: Validation errors
- **WHEN** the backup has validation errors
- **THEN** the page SHALL display the list of validation error messages

#### Scenario: Progress information
- **WHEN** the backup has progress information
- **THEN** the page SHALL display items backed up vs total items
