package report

import (
	"testing"
	"time"

	"github.com/michael/velero-backup-reporter/internal/collector"
)

func timePtr(t time.Time) *time.Time {
	return &t
}

func TestGenerateSummary(t *testing.T) {
	now := time.Now()
	backups := []collector.BackupInfo{
		{Name: "b1", Phase: "Completed", CompletionTimestamp: timePtr(now.Add(-1 * time.Hour))},
		{Name: "b2", Phase: "Completed", CompletionTimestamp: timePtr(now.Add(-30 * time.Minute))},
		{Name: "b3", Phase: "Failed", CompletionTimestamp: timePtr(now.Add(-2 * time.Hour))},
		{Name: "b4", Phase: "PartiallyFailed"},
		{Name: "b5", Phase: "InProgress"},
		{Name: "b6", Phase: "Deleting"},
		{Name: "b7", Phase: "New"},
	}

	summary := generateSummary(backups)

	if summary.TotalBackups != 7 {
		t.Errorf("expected 7 total, got %d", summary.TotalBackups)
	}
	if summary.Completed != 2 {
		t.Errorf("expected 2 completed, got %d", summary.Completed)
	}
	if summary.Failed != 1 {
		t.Errorf("expected 1 failed, got %d", summary.Failed)
	}
	if summary.PartiallyFailed != 1 {
		t.Errorf("expected 1 partially failed, got %d", summary.PartiallyFailed)
	}
	if summary.InProgress != 1 {
		t.Errorf("expected 1 in progress, got %d", summary.InProgress)
	}
	if summary.Deleting != 1 {
		t.Errorf("expected 1 deleting, got %d", summary.Deleting)
	}
	if summary.Other != 1 {
		t.Errorf("expected 1 other, got %d", summary.Other)
	}
	if summary.LastSuccessful == nil {
		t.Fatal("expected last successful timestamp")
	}
	if summary.LastFailed == nil {
		t.Fatal("expected last failed timestamp")
	}
}

func TestGenerateDetails_Duration(t *testing.T) {
	start := time.Now().Add(-10 * time.Minute)
	end := time.Now()

	backups := []collector.BackupInfo{
		{
			Name:                "b1",
			Phase:               "Completed",
			StartTimestamp:      timePtr(start),
			CompletionTimestamp: timePtr(end),
			ItemsBackedUp:      50,
			TotalItems:         100,
		},
	}

	details := generateDetails(backups)

	if len(details) != 1 {
		t.Fatalf("expected 1 detail, got %d", len(details))
	}
	if details[0].Duration < 9*time.Minute || details[0].Duration > 11*time.Minute {
		t.Errorf("expected ~10m duration, got %v", details[0].Duration)
	}
}

func TestGenerateScheduleSummaries(t *testing.T) {
	now := time.Now()
	backups := []collector.BackupInfo{
		{Name: "s1-1", ScheduleName: "daily", Phase: "Completed", StartTimestamp: timePtr(now.Add(-2 * time.Hour))},
		{Name: "s1-2", ScheduleName: "daily", Phase: "Completed", StartTimestamp: timePtr(now.Add(-1 * time.Hour))},
		{Name: "s1-3", ScheduleName: "daily", Phase: "Failed", StartTimestamp: timePtr(now)},
		{Name: "adhoc", ScheduleName: "", Phase: "Completed", StartTimestamp: timePtr(now)},
	}

	schedules := []collector.ScheduleInfo{
		{Name: "daily"},
		{Name: "weekly"}, // no backups
	}

	summaries := generateScheduleSummaries(backups, schedules)

	byName := make(map[string]ScheduleSummary)
	for _, s := range summaries {
		byName[s.ScheduleName] = s
	}

	daily, ok := byName["daily"]
	if !ok {
		t.Fatal("expected 'daily' schedule summary")
	}
	if daily.TotalBackups != 3 {
		t.Errorf("expected 3 backups for daily, got %d", daily.TotalBackups)
	}
	if daily.SuccessfulBackups != 2 {
		t.Errorf("expected 2 successful for daily, got %d", daily.SuccessfulBackups)
	}
	expectedRate := float64(2) / float64(3) * 100
	if daily.SuccessRate < expectedRate-0.1 || daily.SuccessRate > expectedRate+0.1 {
		t.Errorf("expected ~%.1f%% success rate, got %.1f%%", expectedRate, daily.SuccessRate)
	}
	if daily.LastBackupStatus != "Failed" {
		t.Errorf("expected last status 'Failed', got '%s'", daily.LastBackupStatus)
	}

	weekly, ok := byName["weekly"]
	if !ok {
		t.Fatal("expected 'weekly' schedule summary")
	}
	if weekly.TotalBackups != 0 {
		t.Errorf("expected 0 backups for weekly, got %d", weekly.TotalBackups)
	}

	unsched, ok := byName["Unscheduled"]
	if !ok {
		t.Fatal("expected 'Unscheduled' summary")
	}
	if unsched.TotalBackups != 1 {
		t.Errorf("expected 1 unscheduled backup, got %d", unsched.TotalBackups)
	}
}

func TestGenerate_ReportTimestamp(t *testing.T) {
	before := time.Now()
	report := Generate(nil, nil)
	after := time.Now()

	if report.GeneratedAt.Before(before) || report.GeneratedAt.After(after) {
		t.Error("expected report timestamp to be between before and after")
	}
}
