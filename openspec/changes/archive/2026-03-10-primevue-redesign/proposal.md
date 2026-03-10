## Why

The current Bootstrap-based UI is functional but looks generic. Replacing it with PrimeVue — a production-grade Vue 3 component library with the Aura theme — delivers a dramatically more polished, professional interface with rich data tables, refined form controls, and cohesive dark/light mode theming out of the box.

## What Changes

- **BREAKING**: Remove Bootstrap 5 and Bootstrap Icons as dependencies
- Replace all Bootstrap components (navbar, cards, tables, badges, alerts, spinners, forms) with PrimeVue equivalents
- Switch to PrimeVue's Aura theme with built-in dark/light mode (`prefers-color-scheme` detection and manual toggle)
- Replace the backups list table with PrimeVue DataTable (built-in sorting, filtering, column styling)
- Replace the dashboard schedule table with PrimeVue DataTable
- Replace status badges with PrimeVue Tag components
- Replace cards with PrimeVue Card/Panel components
- Replace the loading spinner with PrimeVue Skeleton screens for better perceived performance
- Replace alerts with PrimeVue Message components
- Replace the detail view tabs with PrimeVue TabView
- Replace form selects with PrimeVue Select components
- Add PrimeIcons for consistent iconography
- Update `main.js` to register PrimeVue plugin and theme configuration
- Rewrite `App.vue` layout with PrimeVue Menubar component

## Capabilities

### New Capabilities
- `primevue-ui`: Complete UI component migration from Bootstrap to PrimeVue with Aura theme

### Modified Capabilities
<!-- No spec-level behavior changes — all existing functionality is preserved with new visual components -->

## Impact

- **Frontend (`web/frontend/package.json`)**: Remove `bootstrap`, `bootstrap-icons`; add `primevue`, `@primevue/themes`, `primeicons`
- **Frontend (`web/frontend/src/main.js`)**: Register PrimeVue plugin with Aura theme configuration, replace CSS imports
- **Frontend (`web/frontend/src/App.vue`)**: Complete rewrite — PrimeVue Menubar, theme toggle, layout
- **Frontend (`web/frontend/src/views/`)**: All three views rewritten with PrimeVue components
- **Frontend (`web/frontend/src/components/`)**: LoadingSpinner and ErrorAlert replaced or removed
- **Frontend (`web/frontend/src/composables/`)**: `useBackupUtils.js` retained (utility functions are framework-agnostic)
- **No backend changes required**
