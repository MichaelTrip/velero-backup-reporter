## 1. Dependencies and Setup

- [x] 1.1 Remove `bootstrap` and `bootstrap-icons` npm packages; install `primevue`, `@primeuix/themes`, and `primeicons`
- [x] 1.2 Rewrite `main.js`: remove Bootstrap CSS/JS imports, configure PrimeVue plugin with Aura theme and `darkModeSelector: '.app-dark'`, import PrimeIcons CSS, apply theme from localStorage/OS preference before mount

## 2. App Shell

- [x] 2.1 Rewrite `App.vue`: replace Bootstrap navbar with PrimeVue Menubar (Dashboard, Backups items with router navigation), add theme toggle Button in end slot, update dark mode toggle to add/remove `.app-dark` class on `<html>`, update layout styles
- [x] 2.2 Remove `LoadingSpinner.vue` and `ErrorAlert.vue` components (will be replaced by PrimeVue Skeleton and Message inline)

## 3. Shared Composables

- [x] 3.1 Update `useBackupUtils.js`: change `statusBadgeClass()` to return PrimeVue Tag severity strings (`success`, `danger`, `warn`, `info`, `secondary`) instead of Bootstrap CSS classes; rename to `statusSeverity()`

## 4. Dashboard View

- [x] 4.1 Rewrite `DashboardView.vue`: replace Bootstrap cards with PrimeVue Card components showing status icon + count + Tag, replace loading spinner with Skeleton, replace error alert with Message
- [x] 4.2 Replace Bootstrap schedule table with PrimeVue DataTable + Column components with sortable columns and Tag for status badges

## 5. Backups List View

- [x] 5.1 Rewrite `BackupsListView.vue`: replace filter `<select>` elements with PrimeVue Select components, replace Bootstrap table with PrimeVue DataTable + Column with `sortable` prop, use Tag for status column, use template slots for name links and warning/error highlighting
- [x] 5.2 Add empty state template to the DataTable

## 6. Backup Detail View

- [x] 6.1 Rewrite `BackupDetailView.vue` header: PrimeVue Button components for Back and Download PDF
- [x] 6.2 Rewrite Overview tab content: replace Bootstrap cards with PrimeVue Card/Panel, replace alerts with Message, use Tag for status badges
- [x] 6.3 Rewrite Volumes tab: replace Bootstrap table with PrimeVue DataTable for volume backups, use Tag for phase badges
- [x] 6.4 Rewrite Logs tab: PrimeVue Button for Load Logs, Skeleton while loading, styled pre for log output
- [x] 6.5 Wire up PrimeVue Tabs (Tabs, TabList, Tab, TabPanels, TabPanel) to contain the three tab sections

## 7. Build and Verify

- [x] 7.1 Build the frontend and verify all views render correctly
- [x] 7.2 Build the Go binary and run tests to ensure everything compiles with the new embedded assets
