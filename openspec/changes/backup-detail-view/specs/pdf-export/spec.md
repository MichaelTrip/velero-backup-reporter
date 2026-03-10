## ADDED Requirements

### Requirement: PDF download endpoint
The system SHALL provide an HTTP endpoint to download a PDF report for a single backup.

#### Scenario: Download PDF
- **WHEN** a GET request is made to `/backups/{name}/pdf`
- **THEN** the system SHALL return a PDF file with Content-Type `application/pdf` and a Content-Disposition header for download

#### Scenario: PDF for nonexistent backup
- **WHEN** a GET request is made to `/backups/{name}/pdf` and no backup with that name exists
- **THEN** the system SHALL return an HTTP 404 response

### Requirement: PDF content
The generated PDF SHALL contain the same information displayed on the backup detail page.

#### Scenario: PDF includes metadata
- **WHEN** a PDF is generated for a backup
- **THEN** the PDF SHALL include: backup name, namespace, status, start time, completion time, duration, storage location, and TTL

#### Scenario: PDF includes spec configuration
- **WHEN** a PDF is generated for a backup with namespace/resource filters
- **THEN** the PDF SHALL include the included/excluded namespaces and resources

#### Scenario: PDF includes status details
- **WHEN** a PDF is generated for a backup
- **THEN** the PDF SHALL include volume snapshot counts, warning/error counts, failure reason (if any), and validation errors (if any)

### Requirement: PDF download from detail page
The backup detail page SHALL provide a button or link to download the PDF.

#### Scenario: Download button present
- **WHEN** the backup detail page is rendered
- **THEN** the page SHALL include a download button/link that triggers the PDF download
