## Why

The application currently only supports a light theme. Many operators use monitoring tools in dark environments (NOCs, on-call) where a light UI causes eye strain. Adding a dark mode toggle provides a more comfortable viewing experience and follows modern UI conventions.

## What Changes

- Add a light/dark theme toggle button in the navbar
- Use Bootstrap 5.3's built-in color mode system (`data-bs-theme` attribute on `<html>`)
- Persist the user's theme preference in `localStorage` so it survives page reloads
- Respect the user's OS-level color scheme preference (`prefers-color-scheme`) as the default when no preference is saved
- Adjust any custom styles (footer, status cards) to work correctly in both modes

## Capabilities

### New Capabilities
- `theme-toggle`: Light/dark theme switching with persistence and OS preference detection

### Modified Capabilities
<!-- No existing spec-level behavior changes required -->

## Impact

- **Frontend (`web/frontend/src/App.vue`)**: Theme toggle button added to navbar, theme initialization logic added
- **Frontend (`web/frontend/src/main.js`)**: Theme applied before mount to prevent flash of wrong theme
- **Frontend (views)**: Minor adjustments if any hardcoded colors don't adapt to Bootstrap's dark mode automatically
- **No backend changes required**
- **No new dependencies** — Bootstrap 5.3+ includes dark mode support natively
