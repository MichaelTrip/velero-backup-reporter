## Context

The velero-backup-reporter currently shows a summary dashboard and a flat backup list. Users can see backup names, statuses, and basic metrics but cannot drill into individual backups. Velero Backup CRs contain rich metadata (spec filters, volume snapshot counts, validation errors, failure reasons, TTL, storage location, hooks status) that is not currently collected or displayed.

The existing architecture follows a simple pipeline: collector fetches CRs periodically → report package generates summary models → server renders HTML templates. The collector currently extracts a minimal `BackupInfo` struct. The templates use Go's `html/template` with embedded assets.

## Goals / Non-Goals

**Goals:**
- Allow users to click a backup name and see its full details on a dedicated page
- Collect and display additional Velero Backup CR fields (spec configuration, volume snapshots, failure reasons, validation errors, labels, annotations, storage location, TTL)
- Retrieve and display backup logs on the detail page via Velero's DownloadRequest mechanism
- Generate a downloadable PDF report for a single backup
- Keep the existing dashboard and backup list pages working as-is

**Non-Goals:**
- Aggregated PDF reports across multiple backups (single backup only for now)
- Editable backup details (the only write operation is creating DownloadRequest CRs for log retrieval)
- Real-time updates to the detail page (uses the same periodic collection model)
- Caching or persisting logs (fetched on-demand each time)
- Styled/branded PDF templates (functional layout is sufficient)

## Decisions

### 1. Extend BackupInfo rather than create a separate detail model

**Decision**: Add the additional fields directly to the existing `collector.BackupInfo` struct rather than creating a separate "detailed" fetch path.

**Rationale**: The collector already lists all backups. Extracting a few more fields from the same CR costs negligible overhead. A single model keeps the code simple — the list page uses a subset, the detail page uses all fields.

**Alternatives considered**:
- Lazy-load detail on demand (per-backup API call): Adds complexity, requires the detail handler to have direct kube client access, breaks the collector-as-data-source pattern
- Separate DetailedBackupInfo struct: Duplicates fields, requires maintaining two extraction functions

### 2. Server-side PDF generation with go-pdf/fpdf

**Decision**: Use `go-pdf/fpdf` (the maintained fork of `jung-kurt/gofpdf`) to generate PDFs server-side, served from a `GET /backups/{name}/pdf` endpoint.

**Rationale**: Server-side generation keeps the frontend simple (just a download link). `fpdf` is a mature, dependency-free PDF library for Go with no CGO requirements, making it easy to build and deploy. The PDF content mirrors the detail page.

**Alternatives considered**:
- Browser-side PDF via `window.print()` or jsPDF: Depends on browser rendering, inconsistent results, requires JS dependencies
- `chromedp`/headless Chrome: Heavy dependency, complex to deploy in a minimal container
- `wkhtmltopdf`: Requires external binary, complicates the single-binary deployment goal

### 3. Chi URL parameters for backup name routing

**Decision**: Use chi's URL parameter `{name}` for the backup detail route: `/backups/{name}` and `/backups/{name}/pdf`.

**Rationale**: Chi is already the router in use. URL parameters are idiomatic and make the routes bookmarkable. Backup names in Velero are unique within a namespace.

### 5. On-demand backup log retrieval via DownloadRequest CRs

**Decision**: Fetch backup logs on-demand when a user requests them (not pre-fetched by the collector). The handler creates a `DownloadRequest` CR with `Target.Kind = BackupLog`, polls until the Velero controller provides a signed URL, then fetches and decompresses the gzip log content.

**Rationale**: Backup logs are large (potentially megabytes) and stored in object storage, not in CRs. Pre-fetching all logs during collection would be expensive and wasteful. On-demand retrieval matches the Velero CLI's own approach. The DownloadRequest CR is the official Velero API for this — it works with any storage provider.

**Alternatives considered**:
- Direct object store access (S3/GCS SDK): Would need to support every storage backend Velero supports, duplicating Velero's plugin system. The DownloadRequest mechanism abstracts this away.
- Pre-fetch and cache logs: Memory-intensive for large clusters, logs may be stale, adds complexity for questionable benefit.
- Link directly to signed URL: Would expose object store URLs to the browser, may have CORS issues, and signed URLs expire.

**Flow**:
```
User clicks "View Logs"
  → GET /backups/{name}/logs
  → Handler creates DownloadRequest CR (Kind=BackupLog, Name={name})
  → Poll DownloadRequest status until Phase=Processed (with timeout)
  → HTTP GET to status.downloadURL
  → Decompress gzip response
  → Return plain text logs to browser
  → Clean up DownloadRequest CR
```

### 6. Per-page template sets (already fixed)

**Decision**: Continue using the Clone()-based per-page template pattern established in the recent template fix. The new detail page will get its own template set.

**Rationale**: This avoids the `{{define "content"}}` collision issue and is the idiomatic Go approach for multi-page apps with a shared layout.

## Risks / Trade-offs

- **[Large number of fields]** → The BackupInfo struct grows significantly. Mitigation: fields are simple value types, memory impact is negligible even with thousands of backups.

- **[PDF library size]** → `fpdf` adds to binary size. Mitigation: fpdf is pure Go with no CGO, impact is modest (~2-3MB).

- **[Backup name URL encoding]** → Velero backup names follow Kubernetes naming conventions (lowercase, alphanumeric, hyphens) so URL encoding is not a concern.

- **[Missing backup in detail view]** → A user could request a backup name that doesn't exist (deleted between page load and click). Mitigation: return 404 with a friendly message.

- **[DownloadRequest timeout]** → The Velero controller may be slow or unavailable. Mitigation: use a configurable timeout (default 30s) when polling the DownloadRequest status. Return a clear error message if it times out.

- **[RBAC escalation]** → The app now needs `create` and `get` on `downloadrequests.velero.io`, changing it from purely read-only. Mitigation: document the additional RBAC requirement. Log retrieval is optional — if RBAC isn't granted, the logs button can show an appropriate error.

- **[Log size]** → Backup logs can be very large. Mitigation: stream the response rather than buffering entirely in memory. Consider a reasonable size limit or truncation for the web UI display, with an option to download the full log file.
