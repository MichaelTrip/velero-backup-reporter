## ADDED Requirements

### Requirement: PrimeVue component library setup
The application SHALL use PrimeVue 4 with the Aura theme preset as its UI component library. All UI elements SHALL be rendered using PrimeVue components instead of Bootstrap.

#### Scenario: PrimeVue is loaded on application start
- **WHEN** the application loads
- **THEN** all components render using PrimeVue's Aura theme styling

### Requirement: Dark mode with PrimeVue theming
The application SHALL support dark and light mode via the `.app-dark` CSS class on `<html>`. The theme toggle SHALL persist the preference to `localStorage` and detect the OS color scheme preference as default.

#### Scenario: Toggle to dark mode
- **WHEN** the user clicks the theme toggle while in light mode
- **THEN** the `.app-dark` class is added to `<html>`, all PrimeVue components render in dark mode, and the preference is saved to `localStorage`

#### Scenario: Toggle to light mode
- **WHEN** the user clicks the theme toggle while in dark mode
- **THEN** the `.app-dark` class is removed from `<html>`, all PrimeVue components render in light mode, and the preference is saved to `localStorage`

#### Scenario: OS preference detection
- **WHEN** no `localStorage` theme preference exists and the OS prefers dark mode
- **THEN** the application defaults to dark mode

#### Scenario: No flash of wrong theme
- **WHEN** a user with dark mode saved in `localStorage` loads the application
- **THEN** the page renders in dark mode from the first paint with no flash of light mode

### Requirement: PrimeVue Menubar navigation
The application SHALL use a PrimeVue Menubar component for navigation with items for Dashboard and Backups. The Menubar SHALL include a theme toggle button in its end slot.

#### Scenario: Navigation via Menubar
- **WHEN** the user clicks "Dashboard" or "Backups" in the Menubar
- **THEN** the application navigates to the corresponding route

#### Scenario: Active route indication
- **WHEN** the user is on a particular route
- **THEN** the corresponding Menubar item is visually highlighted

### Requirement: Dashboard with PrimeVue components
The dashboard SHALL display backup summary statistics using PrimeVue Card components and schedule statistics using a PrimeVue DataTable.

#### Scenario: Summary cards display
- **WHEN** the dashboard loads
- **THEN** status summary cards render using PrimeVue Card components with status-colored Tag indicators

#### Scenario: Schedule statistics table
- **WHEN** the dashboard displays schedule data
- **THEN** the schedule statistics render in a PrimeVue DataTable with sortable columns

### Requirement: Backups list with PrimeVue DataTable
The backups list SHALL use a PrimeVue DataTable with built-in sortable columns. Filtering SHALL use PrimeVue Select components for status and schedule.

#### Scenario: Column sorting
- **WHEN** the user clicks a sortable column header
- **THEN** the DataTable sorts by that column (ascending first, then descending on subsequent clicks)

#### Scenario: Status filtering
- **WHEN** the user selects a status from the filter Select
- **THEN** the DataTable displays only backups matching that status

#### Scenario: Empty state
- **WHEN** no backups match the current filters
- **THEN** the DataTable displays a styled empty message

### Requirement: Backup detail with PrimeVue Tabs
The backup detail view SHALL organize content into PrimeVue Tabs: "Overview", "Volumes", and "Logs".

#### Scenario: Default tab
- **WHEN** the user navigates to a backup detail page
- **THEN** the "Overview" tab is active displaying metadata, configuration, and status

#### Scenario: Volumes tab
- **WHEN** the user clicks the "Volumes" tab
- **THEN** the volume snapshots and file system volume backups are displayed

#### Scenario: Logs tab with lazy loading
- **WHEN** the user clicks the "Logs" tab
- **THEN** a "Load Logs" button is shown; clicking it fetches and displays the log output

### Requirement: Status badges as PrimeVue Tags
All backup status indicators SHALL be rendered as PrimeVue Tag components with appropriate severity colors.

#### Scenario: Completed status
- **WHEN** a backup has status "Completed"
- **THEN** it displays a Tag with success severity

#### Scenario: Failed status
- **WHEN** a backup has status "Failed"
- **THEN** it displays a Tag with danger severity

### Requirement: Loading states with Skeleton screens
The application SHALL display PrimeVue Skeleton components during data loading instead of a simple spinner.

#### Scenario: Dashboard loading
- **WHEN** the dashboard data is being fetched
- **THEN** skeleton placeholders render in place of cards and tables

### Requirement: Error display with PrimeVue Message
Errors and alerts SHALL be displayed using PrimeVue Message components with appropriate severity levels.

#### Scenario: API error
- **WHEN** an API call fails
- **THEN** an error Message component is displayed with the error details

#### Scenario: Backup failure reason
- **WHEN** a backup has a failure reason
- **THEN** it is displayed in a Message component with error severity
