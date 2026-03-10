## MODIFIED Requirements

### Requirement: Backup list page
The system SHALL display a page listing all individual backups with their details.

#### Scenario: Backup list with details
- **WHEN** a user navigates to the backups page
- **THEN** the system SHALL display a table of all backups with columns: name, schedule, status, start time, duration, items backed up, warnings, errors

#### Scenario: Clickable backup names
- **WHEN** a user views the backup list
- **THEN** each backup name SHALL be a hyperlink to the backup detail page at `/backups/{name}`

#### Scenario: Filter by status
- **WHEN** a user selects a status filter
- **THEN** the backup list SHALL show only backups matching the selected status

#### Scenario: Filter by schedule
- **WHEN** a user selects a schedule filter
- **THEN** the backup list SHALL show only backups belonging to the selected schedule

#### Scenario: Sort backups
- **WHEN** a user clicks a column header
- **THEN** the backup list SHALL sort by that column in ascending or descending order
