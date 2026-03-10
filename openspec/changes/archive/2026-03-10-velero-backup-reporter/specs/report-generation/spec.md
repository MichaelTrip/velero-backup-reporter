## ADDED Requirements

### Requirement: Generate backup summary report
The system SHALL generate a summary report containing an overview of all collected backups.

#### Scenario: Summary report content
- **WHEN** a report is generated
- **THEN** the report SHALL include: total number of backups, count by status (Completed, Failed, PartiallyFailed, InProgress, Deleting), last successful backup timestamp, and last failed backup timestamp

### Requirement: Generate per-backup detail
The system SHALL include detailed information for each individual backup in the report.

#### Scenario: Backup detail content
- **WHEN** a report includes backup details
- **THEN** each backup entry SHALL show: backup name, schedule name (if any), status, start time, completion time, duration, items backed up, warnings count, and errors count

### Requirement: Generate per-schedule summary
The system SHALL group backups by schedule and provide per-schedule statistics.

#### Scenario: Schedule summary
- **WHEN** backups belong to a schedule
- **THEN** the report SHALL show per-schedule: schedule name, last backup status, last backup time, total backups, success rate percentage

#### Scenario: Unscheduled backups
- **WHEN** backups do not belong to any schedule
- **THEN** they SHALL be grouped under an "Unscheduled" category

### Requirement: Report timestamp
The system SHALL include a generation timestamp on every report.

#### Scenario: Report timestamp present
- **WHEN** a report is generated
- **THEN** the report SHALL display the date and time it was generated
