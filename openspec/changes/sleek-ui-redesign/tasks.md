## 1. Dependencies and Setup

- [x] 1.1 Install `bootstrap-icons` npm package and import its CSS in `main.js`

## 2. Shared Composables

- [x] 2.1 Create `src/composables/useBackupUtils.js` with shared `statusBadgeClass()`, `formatTime()`, and `formatBytes()` functions
- [x] 2.2 Refactor `DashboardView.vue` to use the shared composable instead of local function definitions
- [x] 2.3 Refactor `BackupsListView.vue` to use the shared composable instead of local function definitions
- [x] 2.4 Refactor `BackupDetailView.vue` to use the shared composable instead of local function definitions

## 3. Navbar and App Shell

- [x] 3.1 Update navbar in `App.vue`: add Bootstrap Icons to nav links (speedometer for Dashboard, list-ul for Backups), replace emoji theme toggle with `bi-moon-fill`/`bi-sun-fill` icons
- [x] 3.2 Add `.card-hover` CSS utility class in `App.vue` styles for subtle shadow transition on hover

## 4. Dashboard Improvements

- [x] 4.1 Enhance dashboard status cards: add colored left border accents, subtle shadows, and the `card-hover` class
- [x] 4.2 Improve the empty state for the schedule table with a centered icon and descriptive message

## 5. Backups List Improvements

- [x] 5.1 Replace unicode sort indicators (▲/▼) with Bootstrap Icons (`bi-sort-up`/`bi-sort-down`) and add `cursor: pointer` to sortable headers
- [x] 5.2 Improve the empty state for the backups table with a centered icon and descriptive message

## 6. Backup Detail Tabbed Layout

- [x] 6.1 Reorganize `BackupDetailView.vue` into Bootstrap nav-tabs: "Overview" (metadata, config, status), "Volumes" (snapshots + volume backups table), and "Logs" (log viewer)
- [x] 6.2 Improve empty states within the detail view tabs (no volumes, no logs)

## 7. Build and Verify

- [x] 7.1 Build the frontend and verify all views render correctly in both light and dark modes
