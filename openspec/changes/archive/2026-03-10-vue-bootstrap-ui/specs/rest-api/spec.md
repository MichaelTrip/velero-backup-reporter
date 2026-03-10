## ADDED Requirements

### Requirement: Dashboard API endpoint
The server SHALL expose `GET /api/v1/dashboard` returning a JSON object with backup summary statistics and per-schedule statistics. The summary SHALL include total count and counts by status (Completed, Failed, PartiallyFailed, InProgress, Deleting). Schedule statistics SHALL include schedule name, total count, completed count, failed count, last backup time, and success rate.

#### Scenario: Successful dashboard response
- **WHEN** a client sends `GET /api/v1/dashboard`
- **THEN** the server responds with HTTP 200 and a JSON body containing `summary` and `schedules` fields

#### Scenario: Dashboard response with no data
- **WHEN** a client sends `GET /api/v1/dashboard` and no backups have been collected
- **THEN** the server responds with HTTP 200 and a JSON body with zero counts and an empty schedules array

### Requirement: Backups list API endpoint
The server SHALL expose `GET /api/v1/backups` returning a JSON array of backup objects. Each backup object SHALL include: name, namespace, phase (status), scheduleName, startTimestamp, completionTimestamp, duration, itemsBackedUp, warnings, errors.

#### Scenario: Successful backups list response
- **WHEN** a client sends `GET /api/v1/backups`
- **THEN** the server responds with HTTP 200 and a JSON array of backup objects

#### Scenario: Empty backups list
- **WHEN** a client sends `GET /api/v1/backups` and no backups exist
- **THEN** the server responds with HTTP 200 and an empty JSON array `[]`

### Requirement: Backup detail API endpoint
The server SHALL expose `GET /api/v1/backups/{name}` returning a JSON object with full backup details including: metadata (name, namespace, status, schedule, timestamps, TTL, storage location), configuration (included/excluded namespaces and resources), status (items, warnings, errors, validation errors), volume snapshots, labels, and annotations.

#### Scenario: Successful backup detail response
- **WHEN** a client sends `GET /api/v1/backups/{name}` with a valid backup name
- **THEN** the server responds with HTTP 200 and a JSON object with full backup details

#### Scenario: Backup not found
- **WHEN** a client sends `GET /api/v1/backups/{name}` with a name that does not exist
- **THEN** the server responds with HTTP 404 and a JSON error object `{"error": "backup not found"}`

### Requirement: Backup logs API endpoint
The server SHALL expose `GET /api/v1/backups/{name}/logs` returning the backup log content as plain text (`Content-Type: text/plain`).

#### Scenario: Successful logs response
- **WHEN** a client sends `GET /api/v1/backups/{name}/logs` for a backup with available logs
- **THEN** the server responds with HTTP 200 and the log content as plain text

#### Scenario: Logs not available
- **WHEN** a client sends `GET /api/v1/backups/{name}/logs` and logs cannot be retrieved
- **THEN** the server responds with HTTP 500 and a JSON error object with a descriptive message

### Requirement: Backup PDF API endpoint
The server SHALL expose `GET /api/v1/backups/{name}/pdf` returning the PDF report as a binary download (`Content-Type: application/pdf`).

#### Scenario: Successful PDF download
- **WHEN** a client sends `GET /api/v1/backups/{name}/pdf` for a valid backup
- **THEN** the server responds with HTTP 200, `Content-Type: application/pdf`, and the PDF binary content

### Requirement: JSON error responses
All API endpoints SHALL return errors as JSON objects with the format `{"error": "<message>"}` and appropriate HTTP status codes. API endpoints SHALL set the `Content-Type: application/json` header on all JSON responses.

#### Scenario: Internal server error
- **WHEN** an API endpoint encounters an unexpected error
- **THEN** the server responds with HTTP 500 and `{"error": "<descriptive message>"}`

#### Scenario: Content-Type header
- **WHEN** any API endpoint returns a JSON response
- **THEN** the response includes `Content-Type: application/json` header
