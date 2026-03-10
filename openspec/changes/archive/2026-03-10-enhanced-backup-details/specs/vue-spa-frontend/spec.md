## ADDED Requirements

### Requirement: Volume Backups section in backup detail view
The backup detail view SHALL display a "File System Volume Backups" card listing each volume backup from the `volumeBackups` array. The table SHALL show columns: Volume, Pod, Node, Status (as badge), Uploader, Size (formatted bytes), and Progress (bytes done / total bytes).

#### Scenario: Backup has volume backups
- **WHEN** the backup detail view loads and the `volumeBackups` array is non-empty
- **THEN** a "File System Volume Backups" card is displayed with a table row for each volume backup, showing volume name, pod name, status badge, uploader type, and formatted byte sizes

#### Scenario: Backup has no volume backups
- **WHEN** the backup detail view loads and the `volumeBackups` array is empty
- **THEN** the "File System Volume Backups" card is not displayed

#### Scenario: Volume backup status badge colors
- **WHEN** a volume backup has a status of Completed, Failed, InProgress, or other
- **THEN** the status badge uses the same contextual color scheme as backup status badges (green for Completed, red for Failed, blue for InProgress, grey for other)

### Requirement: Hook and operations status in backup detail view
The backup detail view SHALL display hook execution and async operation counts in the Status card when they are non-zero. Hook status SHALL show "Hooks: X attempted, Y failed". Operations status SHALL show "Async Operations: X attempted, Y completed, Z failed".

#### Scenario: Backup has hook data
- **WHEN** the backup detail view loads and `hooksAttempted` is greater than zero
- **THEN** the Status card displays hook execution counts

#### Scenario: Backup has no hook data
- **WHEN** the backup detail view loads and `hooksAttempted` is zero
- **THEN** no hook information is displayed in the Status card

### Requirement: Human-readable byte formatting
The frontend SHALL format byte values into human-readable units (B, KB, MB, GB, TB) for volume backup sizes.

#### Scenario: Format bytes
- **WHEN** a byte value is displayed (e.g., totalBytes or bytesDone)
- **THEN** the value is shown in the most appropriate unit (e.g., 1073741824 → "1.0 GB")
