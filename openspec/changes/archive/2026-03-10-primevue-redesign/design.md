## Context

The app currently uses Bootstrap 5 + Bootstrap Icons for all UI rendering. While functional, the interface looks generic. The app has 3 views (Dashboard, Backups List, Backup Detail), 2 reusable components (LoadingSpinner, ErrorAlert), and a shared composables file. All data comes from a Go JSON API at `/api/v1/`. The frontend is built with Vite and embedded into the Go binary via `go:embed`.

## Goals / Non-Goals

**Goals:**
- Replace all Bootstrap components with PrimeVue 4 equivalents for a production-grade look
- Use the Aura theme preset for modern, polished styling
- Preserve all existing functionality (dark mode toggle, filtering, sorting, tabs, PDF download, log viewing)
- Use PrimeVue's DataTable with built-in sorting and filtering for a richer data experience
- Maintain the single-binary deployment model (Vite builds, Go embeds)

**Non-Goals:**
- Adding new features, routes, or API endpoints
- Server-side rendering or SSR configuration
- Custom theme token overrides — use Aura defaults
- Auto-import plugin — use explicit imports to keep things clear and simple

## Decisions

### 1. PrimeVue 4 with Aura theme preset

**Choice**: Install `primevue`, `@primeuix/themes` (Aura preset), and `primeicons`. Configure via `app.use(PrimeVue, { theme: { preset: Aura, options: { darkModeSelector: '.app-dark' } } })`.
**Rationale**: PrimeVue is the most mature Vue 3 component library with enterprise-grade components. The Aura theme is modern and clean. Dark mode is controlled by toggling a CSS class on `<html>`.
**Alternatives considered**: Vuetify (Material Design — too opinionated), Naive UI (less mature, smaller community), keeping Bootstrap (doesn't achieve the visual upgrade goal).

### 2. Dark mode via CSS class toggle

**Choice**: Use `darkModeSelector: '.app-dark'` in PrimeVue config. Toggle by adding/removing `.app-dark` on `document.documentElement`. Persist to `localStorage` and detect OS preference on first visit.
**Rationale**: PrimeVue's dark mode is controlled by a CSS class selector, not Bootstrap's `data-bs-theme`. The class-based approach integrates with PrimeVue's theme system. Same localStorage persistence pattern as before.

### 3. PrimeVue DataTable for all tables

**Choice**: Replace Bootstrap `<table>` elements with PrimeVue `DataTable` + `Column` components. Use DataTable's built-in `sortable` prop on columns. Use DataTable's `filterDisplay="row"` for inline filtering on the backups list.
**Rationale**: DataTable provides sorting, filtering, column resize, and responsive behavior out of the box. Eliminates custom sort logic (`toggleSort`, `sortIcon`) in the backups list view. Much more polished than hand-rolled Bootstrap tables.

### 4. PrimeVue Menubar for navigation

**Choice**: Replace Bootstrap's navbar with PrimeVue `Menubar` component. Navigation items use vue-router programmatic navigation.
**Rationale**: Menubar provides responsive collapse, consistent styling with the Aura theme, and supports custom slots for the theme toggle button. Looks cohesive with other PrimeVue components.

### 5. PrimeVue Tabs for detail view

**Choice**: Use PrimeVue `Tabs`, `TabList`, `Tab`, `TabPanels`, `TabPanel` (the v4 API, not deprecated TabView).
**Rationale**: Direct replacement for the Bootstrap nav-tabs currently in the detail view. PrimeVue 4 renamed TabView to the Tabs component system.

### 6. Explicit component imports (no auto-import)

**Choice**: Import each PrimeVue component explicitly in the views that use it.
**Rationale**: Keeps the build simple, avoids adding `unplugin-vue-components` devDependency, makes dependencies obvious. The app only has 3 views — the import overhead is minimal.

### 7. Component mapping

| Current (Bootstrap)       | New (PrimeVue)              |
|---------------------------|-----------------------------|
| `<nav class="navbar">`    | `<Menubar>`                 |
| `<div class="card">`      | `<Card>` or `<Panel>`       |
| `<table class="table">`   | `<DataTable>` + `<Column>`  |
| `<span class="badge">`    | `<Tag>`                     |
| `<div class="alert">`     | `<Message>`                 |
| `<select class="form-select">` | `<Select>`            |
| `<button class="btn">`    | `<Button>`                  |
| Bootstrap nav-tabs        | `<Tabs>` / `<TabList>` / `<Tab>` / `<TabPanels>` / `<TabPanel>` |
| `<div class="spinner">`   | `<Skeleton>` (for better UX)|
| `<i class="bi bi-*">`     | `<i class="pi pi-*">`       |

## Risks / Trade-offs

- **Bundle size increase** → PrimeVue + Aura + PrimeIcons is larger than Bootstrap. Acceptable for an internal monitoring tool. Tree-shaking via explicit imports helps.
- **Learning curve** → PrimeVue API differs from Bootstrap. All views need full rewrites, not incremental changes.
- **DataTable customization** → PrimeVue DataTable's built-in filtering replaces custom filter logic. May need template slots for status badges in columns.
- **Menubar routing** → PrimeVue Menubar uses a model-based API with `command` callbacks; needs integration with vue-router's `router.push()`.
