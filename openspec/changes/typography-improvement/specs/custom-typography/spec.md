## ADDED Requirements

### Requirement: Inter as primary UI font
The application SHALL use Inter as the primary sans-serif font for all UI text, including PrimeVue components, headings, body text, and table content.

#### Scenario: All text renders in Inter
- **WHEN** the application loads
- **THEN** all UI text (headings, body, tables, buttons, menus) renders in the Inter font family

### Requirement: JetBrains Mono as monospace font
The application SHALL use JetBrains Mono as the monospace font for log output and code-like content such as label/annotation key-value badges.

#### Scenario: Log output uses monospace font
- **WHEN** the user views backup logs in the Logs tab
- **THEN** the log output renders in JetBrains Mono

#### Scenario: Labels and annotations use monospace font
- **WHEN** the user views labels or annotations in the backup detail
- **THEN** the key=value badges render in JetBrains Mono for clear readability

### Requirement: Self-hosted fonts
Fonts SHALL be self-hosted via npm packages and bundled by Vite. The application SHALL NOT load fonts from external CDNs.

#### Scenario: No external font requests
- **WHEN** the application loads in an air-gapped environment
- **THEN** all fonts load successfully from the local bundle

### Requirement: Readable typography tuning
The application SHALL apply appropriate font sizes, line heights, and font weights to ensure good readability across all views, particularly in data-dense tables and detail pages.

#### Scenario: Table text readability
- **WHEN** the user views a data table with many rows
- **THEN** the text is legible with adequate line height and font size
