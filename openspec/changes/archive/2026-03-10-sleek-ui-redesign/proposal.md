## Why

The current UI is functional but utilitarian — it relies heavily on default Bootstrap styling with minimal visual hierarchy, no icons, duplicated utility code across views, and a flat card-heavy layout. Improving navigation clarity, data presentation, and overall polish will make the tool more pleasant to use during on-call and NOC monitoring scenarios.

## What Changes

- Extract duplicated utility functions (`statusBadgeClass`, `formatTime`) into shared composables
- Add Bootstrap Icons throughout the UI for better visual cues (navbar, buttons, table headers, status badges)
- Improve the dashboard with better card styling, subtle shadows, and visual hierarchy
- Enhance table presentation with styled sort indicators, better column alignment, and pagination-ready structure
- Redesign the backup detail view to use tabbed sections instead of a long scrolling card list
- Improve empty states with icons and more descriptive messaging
- Add subtle hover effects, transitions, and shadows for a more polished feel
- Improve the navbar with icon-enhanced links and a cleaner toggle button

## Capabilities

### New Capabilities
- `ui-polish`: Shared composables, Bootstrap Icons integration, improved visual styling across all views

### Modified Capabilities
<!-- No existing spec-level behavior changes — this is purely a visual/UX improvement -->

## Impact

- **Frontend (`web/frontend/src/App.vue`)**: Navbar and footer restyling, icon integration
- **Frontend (`web/frontend/src/views/`)**: All three views updated with improved styling, shared composables, tabs in detail view
- **Frontend (`web/frontend/src/composables/`)**: New shared utilities directory
- **Frontend (`web/frontend/package.json`)**: Add `bootstrap-icons` dependency
- **No backend changes required**
