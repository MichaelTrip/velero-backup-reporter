## Why

The application currently relies on PrimeVue's default font stack, which falls back to generic browser fonts. This makes the UI look unpolished and hard to read, especially in data-dense views like tables and detail pages. A proper font with good readability at small sizes will significantly improve the professional appearance.

## What Changes

- Add Inter as the primary UI font (self-hosted via npm, no external CDN dependency)
- Add JetBrains Mono as the monospace font for log output and code-like content
- Override PrimeVue's default font-family via CSS custom properties
- Apply proper font sizing, line height, and weight tuning across the app for better readability
- Use the monospace font for log pre blocks and label/annotation badges

## Capabilities

### New Capabilities
- `custom-typography`: Custom font configuration with Inter and JetBrains Mono, typography tuning

### Modified Capabilities
<!-- No spec-level behavior changes -->

## Impact

- **Frontend (`web/frontend/package.json`)**: Add `@fontsource-variable/inter` and `@fontsource-variable/jetbrains-mono` dependencies
- **Frontend (`web/frontend/src/main.js`)**: Import font CSS files
- **Frontend (`web/frontend/src/App.vue`)**: Override font-family CSS custom properties, add typography styles
- **Frontend (`web/frontend/src/views/BackupDetailView.vue`)**: Apply monospace font to log output and label/annotation badges
- **No backend changes required**
