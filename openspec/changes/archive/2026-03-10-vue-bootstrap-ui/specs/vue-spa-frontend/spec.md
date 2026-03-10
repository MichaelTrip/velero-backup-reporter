## ADDED Requirements

### Requirement: SPA application shell
The application SHALL serve a Vue 3 single-page application with a Bootstrap 5 navbar, a router-view content area, and a footer. The navbar SHALL contain the brand name "Velero Backup Reporter" and navigation links to Dashboard and Backups.

#### Scenario: Initial page load
- **WHEN** a user navigates to the application root URL
- **THEN** the browser loads the SPA shell with navbar, content area, and footer, and the Dashboard view is displayed by default

#### Scenario: Navigation between views
- **WHEN** a user clicks a navigation link in the navbar
- **THEN** the view transitions client-side without a full page reload

### Requirement: Dashboard view
The Dashboard view SHALL display backup status summary cards and a schedule statistics table. Summary cards SHALL show counts for: Total, Completed, Failed, Partially Failed, In Progress, and Deleting backups. Each status card SHALL use a contextual Bootstrap color (success, danger, warning, info). The schedule table SHALL show schedule name, total backups, completed count, failed count, last backup time, and success rate.

#### Scenario: Dashboard displays summary cards
- **WHEN** the Dashboard view loads and receives data from the API
- **THEN** six Bootstrap cards are displayed with backup counts by status, each with the appropriate contextual color

#### Scenario: Dashboard displays schedule statistics
- **WHEN** the Dashboard view loads and schedule data is available
- **THEN** a Bootstrap table displays per-schedule statistics with sortable columns

#### Scenario: Dashboard handles empty state
- **WHEN** the Dashboard view loads and no backups exist
- **THEN** the summary cards show zero counts and the schedule table shows an "No schedules found" message

### Requirement: Backups list view
The Backups list view SHALL display all backups in a Bootstrap table with columns: Name, Schedule, Status, Started, Duration, Items, Warnings, Errors. The table SHALL support client-side filtering by status and schedule name, and sorting by any column. Status values SHALL be rendered as Bootstrap badges with contextual colors.

#### Scenario: Backups table renders with data
- **WHEN** the Backups list view loads and receives backup data
- **THEN** a Bootstrap table displays all backups with status badges and all specified columns

#### Scenario: Filter by status
- **WHEN** a user selects a status from the status filter dropdown
- **THEN** the table shows only backups matching that status

#### Scenario: Filter by schedule
- **WHEN** a user selects a schedule from the schedule filter dropdown
- **THEN** the table shows only backups belonging to that schedule

#### Scenario: Sort by column
- **WHEN** a user clicks a column header
- **THEN** the table rows are sorted by that column, toggling between ascending and descending order

#### Scenario: Navigate to backup detail
- **WHEN** a user clicks a backup name in the table
- **THEN** the application navigates to the Backup Detail view for that backup

### Requirement: Backup detail view
The Backup Detail view SHALL display full backup information organized into sections: Metadata, Configuration, Status, Volume Snapshots, and Labels/Annotations. Each section SHALL use Bootstrap cards. Action buttons SHALL include: Back (navigate to list), Download PDF, and View Logs.

#### Scenario: Detail view displays backup metadata
- **WHEN** the Backup Detail view loads for a specific backup
- **THEN** the metadata section shows: name, namespace, status (as badge), schedule, start time, completion time, expiration, storage location

#### Scenario: Detail view displays configuration
- **WHEN** the backup has included/excluded namespaces or resources
- **THEN** the configuration section lists them

#### Scenario: Download PDF
- **WHEN** a user clicks the "Download PDF" button
- **THEN** the browser downloads the PDF report for that backup from the API

#### Scenario: View logs
- **WHEN** a user clicks the "View Logs" button
- **THEN** the backup logs are fetched from the API and displayed in a preformatted code block within a Bootstrap card

### Requirement: Loading and error states
The application SHALL display a loading spinner while API requests are in flight. The application SHALL display a Bootstrap alert with an error message if an API request fails.

#### Scenario: Loading state
- **WHEN** an API request is pending
- **THEN** a Bootstrap spinner is displayed in the content area

#### Scenario: Error state
- **WHEN** an API request returns an error
- **THEN** a Bootstrap danger alert is displayed with the error message

### Requirement: Responsive layout
The application SHALL use Bootstrap's responsive grid system to provide a usable layout on desktop and tablet screen sizes. Dashboard cards SHALL reflow from a multi-column grid to a single column on smaller screens.

#### Scenario: Desktop layout
- **WHEN** the viewport width is 992px or wider
- **THEN** dashboard cards display in a 3-column grid and tables display all columns

#### Scenario: Tablet layout
- **WHEN** the viewport width is between 576px and 991px
- **THEN** dashboard cards display in a 2-column grid and the table remains horizontally scrollable

### Requirement: Embedded SPA assets
The built Vue SPA assets (HTML, JS, CSS) SHALL be embedded into the Go binary using `go:embed` and served as static files. The Go server SHALL serve `index.html` for any route not matching an API endpoint or static file, enabling client-side routing.

#### Scenario: SPA fallback routing
- **WHEN** the browser requests a non-API path (e.g., `/backups/my-backup`)
- **THEN** the Go server responds with the SPA `index.html` and Vue Router handles the route client-side
