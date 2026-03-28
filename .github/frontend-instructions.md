---
applyTo:
  - web/frontend/**/*.{js,vue,json}
  - web/embed.go
description: Frontend development guide for Vue 3 + PrimeVue UI
---

# Frontend Development Guide

This guide applies to Vue 3 frontend development in `web/frontend/`.

## Quick Start

```bash
# Install dependencies
cd web/frontend
npm ci

# Start dev server (hot reload on localhost:5173)
npm run dev

# Build for production (outputs to web/dist/)
npm run build
```

## Architecture

### Views (Page Components)

Located in `src/views/`:

- **DashboardView.vue** — Overview of backup status: counts, last success/failure, per-schedule stats
- **BackupsListView.vue** — Filterable, sortable table of all backups with inline actions
- **BackupDetailView.vue** — Single backup detailed view with logs, timeline, PDF export

### Composables (Shared Logic)

Located in `src/composables/`:

- **useBackupUtils.js** — Utility functions for backup filtering, status formatting, date handling

Composables should:
- Return refs/computed for reactive data
- Encapsulate complex logic separate from components
- Use descriptive names (`useBackupFilter`, `useStatusFormatter`, etc.)

### Router

Located in `src/router/index.js`:

- Defines route paths: `/`, `/backups`, `/backups/:id`
- Guards/redirects as needed
- Lazy-loads views when possible

### Component Patterns

**Use PrimeVue components** for UI consistency:
- `DataTable` for lists and tables
- `Card` for content containers
- `Button` for actions
- `Dialog` for modals
- `Toast` for notifications (inject `useToast()`)
- `Loader`/`Skeleton` for loading states

**Example component structure:**
```vue
<template>
  <div class="backup-list">
    <PrimeDataTable :value="backups" />
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useBackupUtils } from '@/composables/useBackupUtils'

const { formatStatus, filterByName } = useBackupUtils()
const backups = ref([])

const filteredBackups = computed(() => {
  return filterByName(backups.value, searchTerm.value)
})
</script>

<style scoped>
/* BEM-style class naming */
.backup-list { }
.backup-list__header { }
.backup-list__item { }
</style>
```

## Styling & Theming

- **Framework:** PrimeVue + PrimeIcons (theme + icons included)
- **CSS:** Scoped styles in Vue components (SFC style block)
- **Class naming:** BEM convention (`.component__element--modifier`)
- **Design tokens:** Use PrimeVue's CSS variables (e.g., `var(--primary-color)`, `var(--surface-border)`)
- **Responsiveness:** PrimeVue components are mobile-first; use Media queries for custom layouts

**Theme colors from PrimeVue:**
```css
/* Semantic colors */
var(--primary-color)     /* Main action color */
var(--success-color)     /* Success states */
var(--warning-color)     /* Warning states */
var(--danger-color)      /* Error states */
var(--info-color)        /* Info states */

/* Surface colors */
var(--surface-ground)    /* Page background */
var(--surface-section)   /* Card backgrounds */
var(--surface-border)    /* Dividers/borders */
```

## State Management

**Strategy:** Composables + component-level refs

- Avoid global state library (Pinia not in dependencies)
- Fetch data from backend API endpoints via `fetch()` or custom `useApiClient()` composable
- Handle loading/error states in components
- Use `ref()` for mutable state, `computed()` for derived state

**Example API call in composable:**
```javascript
export function useBackupData() {
  const backups = ref([])
  const loading = ref(false)
  const error = ref(null)

  async function fetchBackups() {
    loading.value = true
    error.value = null
    try {
      const res = await fetch('/api/backups')
      backups.value = await res.json()
    } catch (err) {
      error.value = err.message
    } finally {
      loading.value = false
    }
  }

  onMounted(fetchBackups)
  return { backups, loading, error, fetchBackups }
}
```

## API Integration

Backend serves HTTP API at `/api/` endpoints. Common patterns:

- `GET /api/backups` — List all backups (with filters as query params)
- `GET /api/backups/:id` — Single backup details + logs
- `GET /api/schedules` — List schedules
- `GET /api/dashboard` — Dashboard metrics
- `POST /api/email/test` — Send test email (if `--email-test-enabled`)

**Error handling:**
- Check `response.ok` after `fetch()`
- Display error toasts via `useToast()`
- Fall back to sensible defaults (empty arrays, "-" for missing fields)

## Development Tips

### Hot Reload in Dev Mode
- Vite watches files; changes reflect instantly on `localhost:5173`
- Backend must run separately with embedded frontend assets
- Set backend `--port 8080` to avoid conflict

### Debugging
- Vue DevTools browser extension (for Vue 3)
- Browser DevTools Network tab to inspect `/api/*` calls
- `console.log()` in components or composables

### Building for Production
- Frontend: `npm run build` → outputs optimized bundle to `web/dist/`
- Backend: `make backend` → embeds `web/dist/` into binary
- Always rebuild frontend before rebuilding backend for changes to appear in binary

## Common Patterns

### Loading State
```vue
<template>
  <div v-if="loading" class="loader">
    <ProgressSpinner />
  </div>
  <div v-else>{{ data }}</div>
</template>
```

### Error Handling
```javascript
const { $toast } = useContext()
try {
  await fetchBackups()
} catch (err) {
  $toast.add({ severity: 'error', summary: 'Error', detail: err.message, life: 5000 })
}
```

### Pagination
```vue
<DataTable
  :value="backups"
  :paginator="true"
  :rows="20"
  :totalRecords="totalBackups"
  @page="onPageChange"
/>
```

### Search/Filter
```javascript
const searchTerm = ref('')
const filtered = computed(() => {
  return data.value.filter(item =>
    item.name.toLowerCase().includes(searchTerm.value.toLowerCase())
  )
})
```

## File Organization

```
web/frontend/
├── src/
│   ├── App.vue              # Root layout + navigation
│   ├── main.js              # Vue app initialization
│   ├── composables/
│   │   └── useBackupUtils.js
│   ├── views/
│   │   ├── DashboardView.vue
│   │   ├── BackupsListView.vue
│   │   └── BackupDetailView.vue
│   ├── components/          # Shared sub-components (if needed)
│   └── router/
│       └── index.js
├── index.html               # Entry HTML
├── vite.config.js           # Vite build config
└── package.json
```

## Testing (Future)

Currently no frontend tests; consider:
- Vue Test Utils + Vitest for component tests
- Playwright for E2E tests (against backend at `localhost:8080`)

## Common Issues

| Issue | Solution |
|-------|----------|
| Changes not appearing in browser | Hard reload (Cmd+Shift+R) or rebuild frontend with `make frontend` |
| `npm ci` fails on M1/ARM | Clear cache: `npm cache clean --force` |
| Vite port conflicts | Change `vite.config.js` server port or kill existing process |
| PrimeVue components not styled | Ensure `main.js` imports PrimeVue CSS:  `import 'primevue/resources/themes/lara-light-blue/theme.css'` |
| API calls timeout | Verify backend is running on `localhost:8080` (or configured proxy) |

## Dependencies

| Package | Purpose | Notes |
|---------|---------|-------|
| Vue 3 | UI framework | Latest compatible version |
| Vue Router 4 | Client-side routing | Handles `/`, `/backups`, `/backups/:id` |
| PrimeVue 4 | Component library | Pre-styled, accessible UI components |
| PrimeIcons 7 | Icon library | Used inline in templates |
| Vite | Build tool | Fast dev server + optimized production builds |
| @vitejs/plugin-vue | Vite plugin | Support for `.vue` single-file components |

## Useful Commands

```bash
# Dev
npm run dev              # Start hot-reload server

# Build
npm run build            # Production build → web/dist/

# Formatting (if linters added)
npm run lint            # ESLint (if configured)
npm run format          # Prettier (if configured)

# Integration with backend
make build              # Build frontend AND backend together
```
