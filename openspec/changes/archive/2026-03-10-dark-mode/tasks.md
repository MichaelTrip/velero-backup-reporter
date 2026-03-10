## 1. Theme Initialization

- [x] 1.1 Add theme initialization logic in `main.js` before Vue mounts: read `localStorage('theme')`, fall back to `prefers-color-scheme`, set `data-bs-theme` on `<html>`

## 2. Theme Toggle

- [x] 2.1 Add reactive `theme` ref and `toggleTheme()` function in `App.vue` that updates `localStorage`, `data-bs-theme` attribute, and the ref
- [x] 2.2 Add a theme toggle button to the navbar with moon/sun icons that calls `toggleTheme()`

## 3. Theme-Adaptive Styling

- [x] 3.1 Update footer in `App.vue`: replace hardcoded `bg-light` with a class that adapts to both light and dark modes
- [x] 3.2 Review table headers across views: replace `table-dark` `<thead>` class with a theme-adaptive alternative so headers have proper contrast in both modes

## 4. Build and Verify

- [x] 4.1 Build the frontend and verify both light and dark modes render correctly
