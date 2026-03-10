## Context

The Vue 3 + Bootstrap 5 SPA serves as a monitoring dashboard for Velero backups. The current UI is functional but visually plain — it uses default Bootstrap cards and tables without icons, has duplicated utility functions across views, and presents the backup detail as a long scrolling page of cards. The app already has dark mode support via `data-bs-theme`.

## Goals / Non-Goals

**Goals:**
- Add Bootstrap Icons for visual cues across the entire UI (navbar, buttons, status badges, empty states)
- Extract duplicated code (`statusBadgeClass`, `formatTime`, `formatBytes`) into shared composables
- Improve dashboard cards with subtle shadows, border accents, and better visual hierarchy
- Enhance tables with styled sort indicators and better alignment
- Reorganize the backup detail view using Bootstrap tabs to reduce scrolling
- Improve empty states with icons and descriptive messaging
- Add subtle CSS transitions for hover effects on cards and buttons

**Non-Goals:**
- Custom color palettes or branding beyond Bootstrap defaults
- Adding pagination (data volumes are small — typically < 100 backups)
- Changing any backend APIs or data structures
- Adding new features or data — this is purely visual/UX
- Custom fonts or typography beyond Bootstrap's defaults

## Decisions

### 1. Bootstrap Icons via npm

**Choice**: Install `bootstrap-icons` npm package and import its CSS in `main.js`.
**Rationale**: Official Bootstrap icon library, integrates seamlessly, provides 2000+ icons via `<i class="bi bi-*">` elements. No external CDN needed — bundled by Vite.
**Alternatives considered**: Font Awesome (heavier, separate ecosystem), inline SVGs (more work, harder to maintain), emoji (current approach — inconsistent across platforms).

### 2. Shared composables in `src/composables/`

**Choice**: Create `useBackupUtils.js` composable exporting `statusBadgeClass()`, `formatTime()`, and `formatBytes()`.
**Rationale**: These functions are duplicated across 3 views. A composable is the Vue 3 idiomatic pattern for sharing reactive/non-reactive utilities. Reduces duplication and ensures consistency.

### 3. Bootstrap Nav Tabs for backup detail

**Choice**: Group the backup detail cards into tab panes: "Overview" (metadata, status, config), "Volumes" (snapshots + volume backups), and "Logs".
**Rationale**: The current detail page is a long vertical scroll of 6+ cards. Tabs reduce cognitive load and let users jump directly to the section they need. Bootstrap's nav-tabs component works natively with dark mode.

### 4. CSS utility class for card hover effects

**Choice**: Add a small custom CSS class (`.card-hover`) in `App.vue` styles that adds a subtle shadow transition on hover.
**Rationale**: Minimal CSS addition that significantly improves perceived interactivity. Keeps custom styles centralized and minimal.

### 5. Navbar icon integration

**Choice**: Add Bootstrap Icons to navbar links (dashboard icon, list icon) and replace the emoji theme toggle with a proper `bi-moon-fill`/`bi-sun-fill` icon.
**Rationale**: Icons improve scannability and visual consistency. The emoji toggle renders differently across platforms — a proper icon is more reliable.

## Risks / Trade-offs

- **Bundle size increase** → `bootstrap-icons` CSS adds ~80KB (gzipped ~15KB). Acceptable for an internal monitoring tool. Could tree-shake later if needed.
- **Tab state not in URL** → Switching tabs on the detail view won't update the URL hash. Acceptable since tabs are a UI convenience, not distinct pages.
- **Existing tests** → Server-side tests don't test frontend rendering, so no test breakage expected. Frontend is not unit-tested.
