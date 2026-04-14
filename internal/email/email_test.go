package email

import (
	"strings"
	"testing"
	"time"

	"github.com/michael/velero-backup-reporter/internal/config"
	"github.com/michael/velero-backup-reporter/internal/report"
)

func timePtr(t time.Time) *time.Time {
	return &t
}

func TestNewSender(t *testing.T) {
	cfg := config.SMTPConfig{
		Host: "smtp.example.com",
		Port: 587,
		From: "test@example.com",
		To:   []string{"admin@example.com"},
	}

	sender, err := NewSender(cfg, config.EmailConfig{DetailsWindow: 24 * time.Hour})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sender == nil {
		t.Fatal("expected non-nil sender")
	}
}

func TestBuildMessage(t *testing.T) {
	msg := buildMessage(
		"from@example.com",
		[]string{"to1@example.com", "to2@example.com"},
		"Test Subject",
		"<html><body>Hello</body></html>",
	)

	if !strings.Contains(msg, "From: from@example.com") {
		t.Error("expected From header")
	}
	if !strings.Contains(msg, "To: to1@example.com, to2@example.com") {
		t.Error("expected To header")
	}
	if !strings.Contains(msg, "Subject: Test Subject") {
		t.Error("expected Subject header")
	}
	if !strings.Contains(msg, "Content-Type: text/html") {
		t.Error("expected HTML content type")
	}
	if !strings.Contains(msg, "<html>") {
		t.Error("expected HTML body")
	}
}

func TestEmailTemplateRendering(t *testing.T) {
	cfg := config.SMTPConfig{
		Host: "smtp.example.com",
		Port: 587,
		From: "test@example.com",
		To:   []string{"admin@example.com"},
	}

	sender, err := NewSender(cfg, config.EmailConfig{DetailsWindow: 24 * time.Hour})
	if err != nil {
		t.Fatalf("creating sender: %v", err)
	}

	rpt := report.BackupReport{
		GeneratedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		Summary: report.BackupSummary{
			TotalBackups:    5,
			Completed:       3,
			Failed:          1,
			PartiallyFailed: 1,
			NotStarted:      1,
			LastSuccessful:  timePtr(time.Date(2024, 1, 15, 9, 0, 0, 0, time.UTC)),
			LastFailed:      timePtr(time.Date(2024, 1, 14, 0, 0, 0, 0, time.UTC)),
		},
		ScheduleSummaries: []report.ScheduleSummary{
			{
				ScheduleName:      "daily",
				LastBackupStatus:  "Completed",
				TotalBackups:      3,
				SuccessfulBackups: 2,
				SuccessRate:       66.7,
			},
		},
		Backups: []report.BackupDetail{
			{
				Name:             "backup-1",
				Status:           "Failed",
				StartTime:        timePtr(time.Date(2024, 1, 15, 9, 0, 0, 0, time.UTC)),
				Duration:         5 * time.Minute,
				ItemsBackedUp:    100,
				TotalItems:       100,
				FailureReason:    "plugin error",
				ValidationErrors: []string{"invalid include resource"},
			},
		},
	}

	var buf strings.Builder
	err = sender.template.Execute(&buf, rpt)
	if err != nil {
		t.Fatalf("template execution failed: %v", err)
	}

	html := buf.String()
	if !strings.Contains(html, "Velero Backup Report") {
		t.Error("expected report title in output")
	}
	if !strings.Contains(html, "2024-01-15") {
		t.Error("expected date in output")
	}
	if !strings.Contains(html, "backup-1") {
		t.Error("expected backup name in output")
	}
	if !strings.Contains(html, "plugin error") {
		t.Error("expected failure reason in output")
	}
	if !strings.Contains(html, "invalid include resource") {
		t.Error("expected validation errors in output")
	}
	if !strings.Contains(html, "Not Started") {
		t.Error("expected not started summary label in output")
	}
}

func TestNewSender_DefaultDetailsWindow(t *testing.T) {
	cfg := config.SMTPConfig{
		Host: "smtp.example.com",
		Port: 587,
		From: "test@example.com",
		To:   []string{"admin@example.com"},
	}

	sender, err := NewSender(cfg, config.EmailConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if sender.detailsWindow != 24*time.Hour {
		t.Fatalf("expected default details window 24h, got %v", sender.detailsWindow)
	}
}

func TestFilterBackupDetailsWithinWindow(t *testing.T) {
	now := time.Date(2026, 3, 28, 12, 0, 0, 0, time.UTC)
	inside := now.Add(-23 * time.Hour)
	outside := now.Add(-25 * time.Hour)

	backups := []report.BackupDetail{
		{Name: "inside", StartTime: timePtr(inside)},
		{Name: "outside", StartTime: timePtr(outside)},
	}

	filtered := filterBackupDetailsWithinWindow(backups, now, 24*time.Hour)

	if len(filtered) != 1 {
		t.Fatalf("expected 1 backup in window, got %d", len(filtered))
	}
	if filtered[0].Name != "inside" {
		t.Fatalf("expected backup 'inside', got %q", filtered[0].Name)
	}
}

func TestFilterBackupDetailsWithinWindow_UsesCompletionWhenStartMissing(t *testing.T) {
	now := time.Date(2026, 3, 28, 12, 0, 0, 0, time.UTC)
	completion := now.Add(-2 * time.Hour)

	backups := []report.BackupDetail{
		{Name: "completion-only", CompletionTime: timePtr(completion)},
		{Name: "missing-times"},
	}

	filtered := filterBackupDetailsWithinWindow(backups, now, 24*time.Hour)

	if len(filtered) != 1 {
		t.Fatalf("expected 1 backup in window, got %d", len(filtered))
	}
	if filtered[0].Name != "completion-only" {
		t.Fatalf("expected backup 'completion-only', got %q", filtered[0].Name)
	}
}

func TestFilterBackupDetailsWithinWindow_IncludesNotStartedWithoutTimestamps(t *testing.T) {
	now := time.Date(2026, 3, 28, 12, 0, 0, 0, time.UTC)

	backups := []report.BackupDetail{
		{Name: "not-started", Status: "FailedValidation"},
		{Name: "missing-times", Status: "InProgress"},
	}

	filtered := filterBackupDetailsWithinWindow(backups, now, 24*time.Hour)

	if len(filtered) != 1 {
		t.Fatalf("expected 1 backup in window, got %d", len(filtered))
	}
	if filtered[0].Name != "not-started" {
		t.Fatalf("expected backup 'not-started', got %q", filtered[0].Name)
	}
}
