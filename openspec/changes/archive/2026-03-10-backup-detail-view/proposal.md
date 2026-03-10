## Why

Users can see a list of backups on the backups page, but clicking on a backup does nothing. There is no way to drill into a single backup to see its full details — resource counts, error messages, volume snapshots, labels, and other metadata stored in the Velero Backup CR. Additionally, there is no way to view backup logs (stored in object storage) through the web UI, and no way to export a backup report as a PDF for offline sharing, auditing, or archival purposes.

## What Changes

- Add a backup detail page (`/backups/{name}`) that displays comprehensive information for a single Velero Backup CR, including status, timing, resource details, errors, warnings, volume snapshot info, and labels/annotations
- Make backup names in the backup list page clickable, linking to the detail view
- Enrich the backup data collected from Velero CRs to include additional fields (resource list, volume snapshots, labels, annotations, storage location, TTL, included/excluded resources and namespaces)
- Add backup log retrieval via Velero's DownloadRequest CR mechanism, displaying logs on the detail page
- Add PDF export functionality allowing users to download a backup detail report as a PDF file

## Capabilities

### New Capabilities
- `backup-detail-view`: Displaying a detailed view of a single Velero Backup CR with all available metadata, resource information, and error/warning details
- `pdf-export`: Generating and downloading a PDF report for a single backup's details
- `backup-logs`: Retrieving and displaying Velero backup logs via the DownloadRequest CR mechanism

### Modified Capabilities
- `backup-collection`: Collecting additional fields from Velero Backup CRs (resource list, volume snapshots, labels, annotations, spec fields)
- `web-ui`: Adding the detail page route, clickable backup names in the list, and PDF download button

## Impact

- **Code**: Changes to collector (additional CR fields), report models (richer detail struct), server (new route + handler), new HTML template for detail page, new PDF generation logic
- **Dependencies**: Requires a Go PDF generation library (e.g., `go-pdf/fpdf` or `jung-kurt/gofpdf`)
- **APIs**: New routes `GET /backups/{name}` (detail page), `GET /backups/{name}/pdf` (PDF download), `GET /backups/{name}/logs` (backup logs)
- **Systems**: Requires write access to create DownloadRequest CRs (RBAC change). Requires the Velero server controller to be running to process download requests. Requires network access from the reporter pod to the backup storage location's signed URLs.
- **RBAC**: ServiceAccount needs `create` and `get` permissions on `downloadrequests.velero.io` in addition to existing `list` permissions on backups and schedules
