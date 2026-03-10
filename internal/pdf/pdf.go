package pdf

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/go-pdf/fpdf"
	"github.com/michael/velero-backup-reporter/internal/collector"
)

// GenerateBackupReport creates a PDF report for a single backup.
func GenerateBackupReport(backup *collector.BackupInfo) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetAutoPageBreak(true, 15)
	pdf.AddPage()

	// Title
	pdf.SetFont("Helvetica", "B", 18)
	pdf.CellFormat(0, 12, "Velero Backup Report", "", 1, "C", false, 0, "")
	pdf.Ln(4)
	pdf.SetFont("Helvetica", "", 10)
	pdf.SetTextColor(128, 128, 128)
	pdf.CellFormat(0, 6, fmt.Sprintf("Generated: %s", time.Now().Format("2006-01-02 15:04:05 UTC")), "", 1, "C", false, 0, "")
	pdf.SetTextColor(0, 0, 0)
	pdf.Ln(8)

	// Failure reason alert
	if backup.FailureReason != "" {
		pdf.SetFillColor(254, 226, 226)
		pdf.SetFont("Helvetica", "B", 10)
		pdf.CellFormat(0, 8, "Failure Reason: "+backup.FailureReason, "", 1, "L", true, 0, "")
		pdf.Ln(4)
	}

	// Validation errors
	if len(backup.ValidationErrors) > 0 {
		pdf.SetFillColor(254, 243, 199)
		pdf.SetFont("Helvetica", "B", 10)
		pdf.CellFormat(0, 8, "Validation Errors:", "", 1, "L", true, 0, "")
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
	tableRow(pdf, "Status", backup.Phase)
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
	pdf.SetFillColor(240, 242, 245)
	pdf.CellFormat(0, 9, title, "", 1, "L", true, 0, "")
	pdf.Ln(2)
	pdf.SetFont("Helvetica", "", 10)
}

func tableRow(pdf *fpdf.Fpdf, label, value string) {
	pdf.SetFont("Helvetica", "B", 9)
	pdf.SetFillColor(248, 249, 251)
	pdf.CellFormat(55, 7, label, "1", 0, "L", true, 0, "")
	pdf.SetFont("Helvetica", "", 9)
	pdf.CellFormat(0, 7, value, "1", 1, "L", false, 0, "")
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
