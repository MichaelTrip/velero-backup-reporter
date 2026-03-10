## Context

The Vue + Bootstrap 5 SPA currently uses Bootstrap's default light theme. Bootstrap 5.3+ ships with built-in dark mode support via the `data-bs-theme` attribute on the root `<html>` element. Setting `data-bs-theme="dark"` inverts all Bootstrap component colors automatically — cards, tables, badges, navbar, alerts, etc. — without needing a separate CSS file or theme library.

The app's custom styles are minimal (flexbox layout on body/app, no hardcoded colors beyond Bootstrap classes), so the switch is mostly about toggling the attribute and persisting the preference.

## Goals / Non-Goals

**Goals:**
- Allow users to toggle between light and dark themes via a button in the navbar
- Persist the chosen theme in `localStorage`
- Default to the user's OS preference (`prefers-color-scheme: dark`) when no saved preference exists
- Apply the theme before Vue mounts to avoid a flash of the wrong theme

**Non-Goals:**
- Custom color palettes or branding overrides beyond Bootstrap's built-in dark mode
- Per-user server-side theme preferences (no backend changes)
- Auto-switching on OS preference change while the app is open (toggle is manual)

## Decisions

### 1. Bootstrap 5.3 `data-bs-theme` attribute

**Choice**: Set `data-bs-theme="light"` or `data-bs-theme="dark"` on `<html>`.
**Rationale**: This is Bootstrap's official dark mode mechanism. All Bootstrap components (cards, tables, badges, alerts, forms, navbar) automatically adapt. No additional CSS framework or dependency needed.
**Alternatives considered**: CSS custom properties with manual overrides (much more work), separate dark stylesheet (duplication), third-party theme libraries (unnecessary dependency).

### 2. Theme initialization in `main.js` before mount

**Choice**: Read the theme from `localStorage` (key: `theme`) before calling `createApp().mount()`. If no stored value, check `window.matchMedia('(prefers-color-scheme: dark)')`. Apply `data-bs-theme` on `<html>` immediately.
**Rationale**: Prevents a flash of light theme when the user prefers dark. The attribute is set synchronously before rendering.

### 3. Toggle button with icon in navbar

**Choice**: A button in the navbar that switches between a sun icon (light mode active) and moon icon (dark mode active). Uses simple Unicode characters (or Bootstrap Icons if available).
**Rationale**: Simple, universally understood UI pattern. Placed in the navbar for discoverability.

### 4. Reactive theme state via Vue `ref`

**Choice**: A `ref` in `App.vue` holding the current theme. The toggle function updates `localStorage`, the `data-bs-theme` attribute, and the ref.
**Rationale**: Keeps the toggle reactive so the button icon updates. No need for Pinia/Vuex since it's a single component managing the state.

## Risks / Trade-offs

- **Footer hardcoded `bg-light`** → Must change to a theme-adaptive class or remove the hardcoded background so it adapts with Bootstrap's dark mode.
- **Custom badge colors on cards** → The status cards use `text-white` which works on colored backgrounds in both modes; no issue expected.
- **`table-dark` header class** → Currently uses `table-dark` for `<thead>` which looks the same in both modes. May need adjustment for contrast in dark mode.
