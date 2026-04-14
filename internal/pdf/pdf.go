package pdf

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/go-pdf/fpdf"
	"github.com/michael/velero-backup-reporter/internal/collector"
	reportpkg "github.com/michael/velero-backup-reporter/internal/report"
)

// GenerateWindowReport creates a consolidated PDF report for a selected time window.
func GenerateWindowReport(rpt reportpkg.BackupReport, windowLabel string) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetAutoPageBreak(true, 15)
	pdf.AddPage()
	pdf.SetMargins(12, 12, 12)

	drawReportHeader(pdf, "Velero Backup Report", windowLabel, rpt.GeneratedAt)

	sectionHeader(pdf, "Summary")
	tableRow(pdf, "Total Backups", fmt.Sprintf("%d", rpt.Summary.TotalBackups))
	tableRowWithColor(pdf, "Completed", fmt.Sprintf("%d", rpt.Summary.Completed), 22, 101, 52)
	tableRowWithColor(pdf, "Failed", fmt.Sprintf("%d", rpt.Summary.Failed), 153, 27, 27)
	tableRowWithColor(pdf, "Partially Failed", fmt.Sprintf("%d", rpt.Summary.PartiallyFailed), 146, 64, 14)
	tableRowWithColor(pdf, "Missing / Not Started", fmt.Sprintf("%d", rpt.Summary.NotStarted), 153, 27, 27)
	tableRow(pdf, "Last Successful", formatTimePtr(rpt.Summary.LastSuccessful))
	tableRow(pdf, "Last Failed", formatTimePtr(rpt.Summary.LastFailed))
	pdf.Ln(4)

	if len(rpt.ScheduleSummaries) > 0 {
		sectionHeader(pdf, "Schedules")
		for _, s := range rpt.ScheduleSummaries {
			statusColorR, statusColorG, statusColorB := statusTextColor(s.LastBackupStatus)
			tableRow(pdf, "Schedule", s.ScheduleName)
			tableRowWithColor(pdf, "Last Status", valueOrDash(s.LastBackupStatus), statusColorR, statusColorG, statusColorB)
			tableRow(pdf, "Total", fmt.Sprintf("%d", s.TotalBackups))
			tableRow(pdf, "Success Rate", fmt.Sprintf("%.1f%%", s.SuccessRate))
			pdf.Ln(1)
		}
		pdf.Ln(3)
	}

	sectionHeader(pdf, "Backup Details")
	for _, b := range rpt.Backups {
		pdf.SetFillColor(255, 247, 237)
		pdf.SetDrawColor(254, 186, 116)
		pdf.MultiCell(0, 7, "", "1", "L", true)
		currentY := pdf.GetY()
		pdf.SetXY(12, currentY-7)

		pdf.SetFont("Helvetica", "B", 10)
		pdf.CellFormat(120, 7, b.Name, "", 0, "L", false, 0, "")
		pdf.SetFont("Helvetica", "B", 9)
		r, g, bb := statusTextColor(b.Status)
		pdf.SetTextColor(r, g, bb)
		pdf.CellFormat(0, 7, b.Status, "", 1, "R", false, 0, "")
		pdf.SetTextColor(0, 0, 0)

		tableRow(pdf, "Schedule", valueOrDash(b.ScheduleName))
		tableRow(pdf, "Started", formatTimePtr(b.StartTime))
		tableRow(pdf, "Duration", formatDurationHuman(b.Duration))
		tableRow(pdf, "Items Backed Up", fmt.Sprintf("%d / %d", b.ItemsBackedUp, b.TotalItems))
		tableRow(pdf, "Warnings / Errors", fmt.Sprintf("%d / %d", b.Warnings, b.Errors))
		tableRow(pdf, "Failure Reason", valueOrDash(b.FailureReason))
		if len(b.ValidationErrors) > 0 {
			tableRow(pdf, "Validation Errors", strings.Join(b.ValidationErrors, "; "))
		} else {
			tableRow(pdf, "Validation Errors", "-")
		}
		pdf.Ln(2)
	}

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("generating PDF: %w", err)
	}
	return buf.Bytes(), nil
}

// GenerateBackupReport creates a PDF report for a single backup.
func GenerateBackupReport(backup *collector.BackupInfo) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetAutoPageBreak(true, 15)
	pdf.AddPage()
	pdf.SetMargins(12, 12, 12)

	// Card container background
	pdf.SetFillColor(255, 255, 255)
	pdf.Rect(10, 10, 190, 277, "F")

	// Header band approximates the email gradient using layered fills.
	pdf.SetFillColor(15, 23, 42)
	pdf.Rect(10, 10, 190, 18, "F")
	pdf.SetFillColor(29, 78, 216)
	pdf.Rect(95, 10, 105, 18, "F")
	pdf.SetFillColor(14, 165, 233)
	pdf.Rect(145, 10, 55, 18, "F")
	pdf.SetFillColor(245, 158, 11)
	pdf.Rect(10, 28, 190, 2, "F")

	pdf.SetY(14)
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Helvetica", "B", 17)
	pdf.CellFormat(0, 7, "Velero Backup Report", "", 1, "L", false, 0, "")
	pdf.SetFont("Helvetica", "", 9)
	pdf.SetTextColor(219, 234, 254)
	pdf.CellFormat(0, 5, fmt.Sprintf("Generated at %s", time.Now().Format("2006-01-02 15:04:05 UTC")), "", 1, "L", false, 0, "")
	pdf.SetTextColor(0, 0, 0)
	pdf.Ln(8)

	// Failure reason alert
	if backup.FailureReason != "" {
		pdf.SetFillColor(254, 226, 226)
		pdf.SetFont("Helvetica", "B", 10)
		pdf.SetTextColor(153, 27, 27)
		pdf.CellFormat(0, 8, "Failure Reason: "+backup.FailureReason, "", 1, "L", true, 0, "")
		pdf.SetTextColor(0, 0, 0)
		pdf.Ln(4)
	}

	// Validation errors
	if len(backup.ValidationErrors) > 0 {
		pdf.SetFillColor(254, 243, 199)
		pdf.SetFont("Helvetica", "B", 10)
		pdf.SetTextColor(146, 64, 14)
		pdf.CellFormat(0, 8, "Validation Errors:", "", 1, "L", true, 0, "")
		pdf.SetTextColor(0, 0, 0)
		pdf.SetFont("Helvetica", "", 9)
		for _, e := range backup.ValidationErrors {
			pdf.CellFormat(0, 6, "  - "+e, "", 1, "L", false, 0, "")
		}
		pdf.Ln(4)
	}

	// Metadata section
	sectionHeader(pdf, "Metadata")
	var duration string
	if backup.StartTimestamp != nil && backup.CompletionTimestamp != nil {
		duration = backup.CompletionTimestamp.Sub(*backup.StartTimestamp).String()
	} else {
		duration = "-"
	}
	tableRow(pdf, "Name", backup.Name)
	tableRow(pdf, "Namespace", backup.Namespace)
	tableRowStatus(pdf, "Status", backup.Phase)
	tableRow(pdf, "Schedule", valueOrDash(backup.ScheduleName))
	tableRow(pdf, "Start Time", formatTimePtr(backup.StartTimestamp))
	tableRow(pdf, "Completion Time", formatTimePtr(backup.CompletionTimestamp))
	tableRow(pdf, "Duration", duration)
	tableRow(pdf, "Expiration", formatTimePtr(backup.Expiration))
	tableRow(pdf, "Storage Location", valueOrDash(backup.StorageLocation))
	tableRow(pdf, "TTL", valueOrDash(backup.TTL))
	pdf.Ln(4)

	// Configuration section
	sectionHeader(pdf, "Configuration")
	tableRow(pdf, "Included Namespaces", sliceOrAll(backup.IncludedNamespaces))
	tableRow(pdf, "Excluded Namespaces", sliceOrDash(backup.ExcludedNamespaces))
	tableRow(pdf, "Included Resources", sliceOrAll(backup.IncludedResources))
	tableRow(pdf, "Excluded Resources", sliceOrDash(backup.ExcludedResources))
	pdf.Ln(4)

	// Status section
	sectionHeader(pdf, "Status")
	tableRow(pdf, "Items Backed Up", fmt.Sprintf("%d / %d", backup.ItemsBackedUp, backup.TotalItems))
	tableRow(pdf, "Warnings", fmt.Sprintf("%d", backup.Warnings))
	tableRow(pdf, "Errors", fmt.Sprintf("%d", backup.Errors))
	pdf.Ln(4)

	// Volume Snapshots section
	sectionHeader(pdf, "Volume Snapshots")
	tableRow(pdf, "Snapshots Attempted", fmt.Sprintf("%d", backup.VolumeSnapshotsAttempted))
	tableRow(pdf, "Snapshots Completed", fmt.Sprintf("%d", backup.VolumeSnapshotsCompleted))
	tableRow(pdf, "CSI Snapshots Attempted", fmt.Sprintf("%d", backup.CSIVolumeSnapshotsAttempted))
	tableRow(pdf, "CSI Snapshots Completed", fmt.Sprintf("%d", backup.CSIVolumeSnapshotsCompleted))
	pdf.Ln(4)

	// Labels
	if len(backup.Labels) > 0 {
		sectionHeader(pdf, "Labels")
		for k, v := range backup.Labels {
			tableRow(pdf, k, v)
		}
		pdf.Ln(4)
	}

	// Annotations
	if len(backup.Annotations) > 0 {
		sectionHeader(pdf, "Annotations")
		for k, v := range backup.Annotations {
			tableRow(pdf, k, v)
		}
		pdf.Ln(4)
	}

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("generating PDF: %w", err)
	}
	return buf.Bytes(), nil
}

func sectionHeader(pdf *fpdf.Fpdf, title string) {
	pdf.SetFont("Helvetica", "B", 13)
	pdf.SetTextColor(15, 23, 42)
	switch title {
	case "Metadata":
		pdf.SetFillColor(224, 231, 255)
		pdf.SetTextColor(29, 78, 216)
	case "Configuration":
		pdf.SetFillColor(240, 253, 250)
		pdf.SetTextColor(15, 118, 110)
	case "Status":
		pdf.SetFillColor(255, 237, 213)
		pdf.SetTextColor(124, 45, 18)
	case "Volume Snapshots":
		pdf.SetFillColor(236, 254, 255)
		pdf.SetTextColor(14, 116, 144)
	default:
		pdf.SetFillColor(240, 242, 245)
	}
	pdf.CellFormat(0, 9, title, "", 1, "L", true, 0, "")
	pdf.Ln(2)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Helvetica", "", 10)
}

func tableRow(pdf *fpdf.Fpdf, label, value string) {
	pdf.SetFont("Helvetica", "B", 9)
	pdf.SetFillColor(255, 237, 213)
	pdf.SetDrawColor(254, 215, 170)
	pdf.CellFormat(55, 7, label, "1", 0, "L", true, 0, "")
	pdf.SetFont("Helvetica", "", 9)
	pdf.CellFormat(0, 7, value, "1", 1, "L", false, 0, "")
}

func tableRowStatus(pdf *fpdf.Fpdf, label, value string) {
	r, g, b := statusTextColor(value)
	tableRowWithColor(pdf, label, value, r, g, b)
}

func tableRowWithColor(pdf *fpdf.Fpdf, label, value string, r, g, b int) {
	pdf.SetFont("Helvetica", "B", 9)
	pdf.SetFillColor(255, 237, 213)
	pdf.SetDrawColor(254, 215, 170)
	pdf.CellFormat(55, 7, label, "1", 0, "L", true, 0, "")
	pdf.SetFont("Helvetica", "", 9)
	pdf.SetTextColor(r, g, b)
	pdf.CellFormat(0, 7, value, "1", 1, "L", false, 0, "")
	pdf.SetTextColor(0, 0, 0)
}

func statusTextColor(status string) (int, int, int) {
	switch {
	case status == "Completed":
		return 22, 101, 52
	case status == "Missed":
		return 153, 27, 27
	case strings.Contains(status, "PartiallyFailed"):
		return 146, 64, 14
	case strings.Contains(status, "Failed"):
		return 153, 27, 27
	case status == "InProgress" || status == "New" || status == "Queued" || status == "ReadyToStart" || status == "WaitingForPluginOperations" || status == "Finalizing":
		return 30, 64, 175
	default:
		return 75, 85, 99
	}
}

func drawReportHeader(pdf *fpdf.Fpdf, title, subtitle string, generatedAt time.Time) {
	pdf.SetFillColor(15, 23, 42)
	pdf.Rect(10, 10, 190, 18, "F")
	pdf.SetFillColor(29, 78, 216)
	pdf.Rect(95, 10, 105, 18, "F")
	pdf.SetFillColor(14, 165, 233)
	pdf.Rect(145, 10, 55, 18, "F")
	pdf.SetFillColor(245, 158, 11)
	pdf.Rect(10, 28, 190, 2, "F")

	pdf.SetY(14)
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Helvetica", "B", 17)
	pdf.CellFormat(0, 7, title, "", 1, "L", false, 0, "")
	pdf.SetFont("Helvetica", "", 9)
	pdf.SetTextColor(219, 234, 254)
	pdf.CellFormat(0, 5, fmt.Sprintf("Generated at %s", generatedAt.UTC().Format("2006-01-02 15:04:05 UTC")), "", 1, "L", false, 0, "")
	if subtitle != "" {
		pdf.CellFormat(0, 5, subtitle, "", 1, "L", false, 0, "")
	}
	pdf.SetTextColor(0, 0, 0)
	pdf.Ln(8)
}

func formatDurationHuman(d time.Duration) string {
	if d <= 0 {
		return "-"
	}
	return d.Round(time.Second).String()
}

func formatTimePtr(t *time.Time) string {
	if t == nil {
		return "-"
	}
	return t.Format("2006-01-02 15:04:05 UTC")
}

func valueOrDash(s string) string {
	if s == "" {
		return "-"
	}
	return s
}

func sliceOrDash(s []string) string {
	if len(s) == 0 {
		return "-"
	}
	return strings.Join(s, ", ")
}

func sliceOrAll(s []string) string {
	if len(s) == 0 {
		return "All"
	}
	return strings.Join(s, ", ")
}
