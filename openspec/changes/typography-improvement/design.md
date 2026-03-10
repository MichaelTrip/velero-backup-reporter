## Context

The app uses PrimeVue 4 with the Aura theme. PrimeVue sets its body font via the CSS custom property `--p-font-family`, which defaults to a generic system font stack. There is no custom font currently loaded. Log output uses a bare `<pre>` tag with no explicit monospace font. The app is a single-binary Go deployment where the Vue frontend is embedded — fonts must be self-hosted (no external CDN).

## Goals / Non-Goals

**Goals:**
- Set Inter as the primary sans-serif font for all UI text
- Set JetBrains Mono as the monospace font for logs and code-like content
- Self-host fonts via `@fontsource-variable` npm packages (bundled by Vite, embedded in Go binary)
- Tune font sizes and line heights for better readability in data tables and detail views

**Non-Goals:**
- Changing PrimeVue component structure or layout
- Custom font weight subsets or performance optimization (the app is internal)
- Changing colors or theme tokens beyond font-related properties

## Decisions

### 1. Inter as primary font

**Choice**: Use Inter via `@fontsource-variable/inter` npm package.
**Rationale**: Inter is designed specifically for screen UI readability. It has excellent legibility at small sizes (important for tables), clear number forms, and a professional appearance. It's open source (SIL OFL), widely used in production tools (GitHub, Figma, Linear).
**Alternatives considered**: IBM Plex Sans (good but less common), Geist (newer, less proven), system font stack (current — the problem we're solving).

### 2. JetBrains Mono as monospace font

**Choice**: Use JetBrains Mono via `@fontsource-variable/jetbrains-mono` npm package.
**Rationale**: Specifically designed for code and terminal output readability. Clear distinction between similar characters (0/O, 1/l/I). Works well at the small sizes used in log viewers. Free and open source.
**Alternatives considered**: Fira Code (ligatures not needed), Source Code Pro (less distinctive), system monospace (inconsistent across platforms).

### 3. @fontsource-variable for self-hosting

**Choice**: Use `@fontsource-variable` packages which provide variable font files + CSS.
**Rationale**: Vite bundles the font files into `dist/assets/`, which gets embedded by Go's `go:embed`. No external CDN requests. Variable fonts mean a single file covers all weights (lighter than multiple static weight files).

### 4. CSS custom property override

**Choice**: Override PrimeVue's `--p-font-family` on `:root` to set Inter as the base font. Set a custom `--font-mono` property for monospace usage.
**Rationale**: PrimeVue reads `--p-font-family` for all component text. Overriding at `:root` level ensures all PrimeVue components use Inter without modifying theme tokens at the plugin level.

## Risks / Trade-offs

- **Bundle size** → Inter variable font ~300KB, JetBrains Mono variable font ~200KB. Acceptable for an embedded internal tool. The fonts are cached by the browser after first load.
- **Font flash** → Variable fonts load fast from local bundle. No FOUT expected since fonts are served from the same origin as the app.
