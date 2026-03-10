## MODIFIED Requirements

### Requirement: Generate per-backup detail
The system SHALL include detailed information for each individual backup in the report.

#### Scenario: Backup detail content
- **WHEN** a report includes backup details
- **THEN** each backup entry SHALL show: backup name, schedule name (if any), status, start time, completion time, duration, items backed up, warnings count, errors count, storage location, TTL, volume snapshots attempted/completed, CSI volume snapshots attempted/completed, failure reason, validation errors, included/excluded namespaces, included/excluded resources, labels, and annotations
