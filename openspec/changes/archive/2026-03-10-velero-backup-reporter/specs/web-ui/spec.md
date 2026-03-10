## ADDED Requirements

### Requirement: Serve web UI over HTTP
The system SHALL serve a web-based user interface over HTTP on a configurable port.

#### Scenario: Default port
- **WHEN** no port is configured
- **THEN** the web UI SHALL be served on port 8080

#### Scenario: Custom port
- **WHEN** a port is configured
- **THEN** the web UI SHALL be served on the specified port

### Requirement: Dashboard page
The system SHALL display a dashboard page as the default landing page showing the backup summary report.

#### Scenario: Dashboard displays summary
- **WHEN** a user navigates to the root URL
- **THEN** the dashboard SHALL display the backup summary with status counts, last successful backup, and last failed backup

#### Scenario: Dashboard displays schedule overview
- **WHEN** a user navigates to the root URL
- **THEN** the dashboard SHALL display a table of schedules with their latest backup status and success rate

### Requirement: Backup list page
The system SHALL display a page listing all individual backups with their details.

#### Scenario: Backup list with details
- **WHEN** a user navigates to the backups page
- **THEN** the system SHALL display a table of all backups with columns: name, schedule, status, start time, duration, items backed up, warnings, errors

#### Scenario: Filter by status
- **WHEN** a user selects a status filter
- **THEN** the backup list SHALL show only backups matching the selected status

#### Scenario: Filter by schedule
- **WHEN** a user selects a schedule filter
- **THEN** the backup list SHALL show only backups belonging to the selected schedule

#### Scenario: Sort backups
- **WHEN** a user clicks a column header
- **THEN** the backup list SHALL sort by that column in ascending or descending order

### Requirement: Health endpoint
The system SHALL expose a `/healthz` endpoint for Kubernetes liveness/readiness probes.

#### Scenario: Health check response
- **WHEN** a GET request is made to `/healthz`
- **THEN** the system SHALL return HTTP 200 with a JSON body indicating health status

### Requirement: Static assets embedded
The system SHALL embed all static assets (HTML, CSS, JS) in the Go binary.

#### Scenario: No external file dependencies
- **WHEN** the application binary is deployed
- **THEN** all web UI assets SHALL be served from the embedded filesystem without requiring external files
