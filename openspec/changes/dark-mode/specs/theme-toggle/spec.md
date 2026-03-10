## ADDED Requirements

### Requirement: Theme toggle button in navbar
The application SHALL display a theme toggle button in the navbar. The button SHALL show a moon icon when light mode is active (indicating "switch to dark") and a sun icon when dark mode is active (indicating "switch to light").

#### Scenario: Toggle from light to dark
- **WHEN** the user clicks the theme toggle button while in light mode
- **THEN** the application switches to dark mode, the `data-bs-theme` attribute on `<html>` is set to `"dark"`, and the button icon changes to a sun

#### Scenario: Toggle from dark to light
- **WHEN** the user clicks the theme toggle button while in dark mode
- **THEN** the application switches to light mode, the `data-bs-theme` attribute on `<html>` is set to `"light"`, and the button icon changes to a moon

### Requirement: Theme persistence in localStorage
The application SHALL persist the user's theme choice in `localStorage` under the key `theme`. On subsequent page loads, the application SHALL read and apply the stored theme preference.

#### Scenario: Theme persists across reload
- **WHEN** the user selects dark mode and reloads the page
- **THEN** the application loads in dark mode without a flash of light mode

#### Scenario: No stored preference
- **WHEN** `localStorage` has no `theme` key
- **THEN** the application uses the OS-level color scheme preference as the default

### Requirement: OS color scheme preference detection
The application SHALL detect the user's OS-level color scheme preference using `prefers-color-scheme` media query. This preference SHALL be used as the default theme when no `localStorage` preference exists.

#### Scenario: OS prefers dark
- **WHEN** the user's OS is set to dark mode and no `localStorage` preference exists
- **THEN** the application defaults to dark mode

#### Scenario: OS prefers light
- **WHEN** the user's OS is set to light mode and no `localStorage` preference exists
- **THEN** the application defaults to light mode

### Requirement: No flash of wrong theme on load
The application SHALL apply the correct theme to the `<html>` element before Vue mounts, preventing a visible flash of the incorrect theme during page load.

#### Scenario: Dark mode user loads page
- **WHEN** a user with dark mode saved in `localStorage` loads the application
- **THEN** the page renders in dark mode from the first paint with no flash of light mode

### Requirement: All UI elements adapt to dark mode
All Bootstrap components (cards, tables, badges, alerts, navbar, footer, forms) SHALL render correctly in both light and dark modes. No hardcoded light-only background colors SHALL remain on theme-aware elements.

#### Scenario: Footer adapts to dark mode
- **WHEN** dark mode is active
- **THEN** the footer background and text colors adapt to the dark theme instead of remaining light

#### Scenario: Tables adapt to dark mode
- **WHEN** dark mode is active
- **THEN** table backgrounds, borders, and text adapt to the dark theme
