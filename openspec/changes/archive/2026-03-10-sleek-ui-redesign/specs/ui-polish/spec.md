## ADDED Requirements

### Requirement: Bootstrap Icons integration
The application SHALL use Bootstrap Icons (`bi bi-*` classes) for all iconography instead of emoji characters. Icons SHALL render consistently across all platforms and adapt to the active theme.

#### Scenario: Navbar displays icons
- **WHEN** the application loads
- **THEN** the navbar links display appropriate icons (e.g., speedometer for Dashboard, list for Backups) alongside their text labels

#### Scenario: Theme toggle uses proper icon
- **WHEN** the user views the theme toggle button
- **THEN** it displays a `bi-moon-fill` icon in light mode and a `bi-sun-fill` icon in dark mode instead of emoji characters

### Requirement: Shared composables for utility functions
The application SHALL provide shared composable functions so that `statusBadgeClass()`, `formatTime()`, and `formatBytes()` are defined once and reused across all views.

#### Scenario: Status badge class consistency
- **WHEN** a backup status is displayed in any view (Dashboard, Backups List, Backup Detail)
- **THEN** the badge class is computed using the same shared `statusBadgeClass()` function

#### Scenario: Time formatting consistency
- **WHEN** a timestamp is displayed in any view
- **THEN** it is formatted using the same shared `formatTime()` function

### Requirement: Dashboard card visual improvements
The dashboard status cards SHALL have subtle shadows, colored left borders matching their status color, and a hover effect that slightly elevates the card.

#### Scenario: Card hover effect
- **WHEN** the user hovers over a dashboard status card
- **THEN** the card displays a slightly elevated shadow via CSS transition

#### Scenario: Card border accent
- **WHEN** the dashboard loads
- **THEN** each status card displays a colored left border matching its status color (e.g., green for Completed, red for Failed)

### Requirement: Improved table sort indicators
The backups list table SHALL display styled sort indicator icons (`bi-sort-up` / `bi-sort-down`) instead of plain unicode triangles. Sortable column headers SHALL have a visible cursor pointer.

#### Scenario: Sort icon display
- **WHEN** a column is sorted ascending
- **THEN** the column header displays a `bi-sort-up` icon

#### Scenario: Sort icon descending
- **WHEN** a column is sorted descending
- **THEN** the column header displays a `bi-sort-down` icon

### Requirement: Tabbed backup detail view
The backup detail view SHALL organize content into Bootstrap nav-tabs: "Overview" (metadata, configuration, status cards), "Volumes" (volume snapshots and volume backups), and "Logs" (log output). The Overview tab SHALL be active by default.

#### Scenario: Default tab on load
- **WHEN** the user navigates to a backup detail page
- **THEN** the "Overview" tab is active and displays metadata, configuration, and status information

#### Scenario: Volumes tab
- **WHEN** the user clicks the "Volumes" tab
- **THEN** the view displays volume snapshot counts and the file system volume backups table

#### Scenario: Logs tab
- **WHEN** the user clicks the "Logs" tab and logs have not been loaded
- **THEN** the view displays a "Load Logs" button; clicking it fetches and displays the logs

### Requirement: Improved empty states
Empty states (no backups, no schedules, no volume backups) SHALL display a relevant Bootstrap Icon and descriptive text instead of plain muted text.

#### Scenario: No backups empty state
- **WHEN** the backups list has no data
- **THEN** the view displays a centered icon and message such as "No backups found"

#### Scenario: No schedules empty state
- **WHEN** the dashboard has no schedule statistics
- **THEN** the view displays a centered icon and message such as "No schedules found"

### Requirement: Navbar visual improvements
The navbar links SHALL include icons alongside text labels. The active link SHALL have clear visual distinction.

#### Scenario: Navbar link icons
- **WHEN** the application renders the navbar
- **THEN** each navigation link displays an icon before its text label
