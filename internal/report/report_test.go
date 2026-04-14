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
		{Name: "b8", Phase: "FailedValidation", StartTimestamp: timePtr(now.Add(-20 * time.Minute))},
		{Name: "b9", Phase: "FinalizingPartiallyFailed", CompletionTimestamp: timePtr(now.Add(-10 * time.Minute))},
	}

	summary := generateSummary(backups)

	if summary.TotalBackups != 9 {
		t.Errorf("expected 9 total, got %d", summary.TotalBackups)
	}
	if summary.Completed != 2 {
		t.Errorf("expected 2 completed, got %d", summary.Completed)
	}
	if summary.Failed != 2 {
		t.Errorf("expected 2 failed, got %d", summary.Failed)
	}
	if summary.PartiallyFailed != 2 {
		t.Errorf("expected 2 partially failed, got %d", summary.PartiallyFailed)
	}
	if summary.NotStarted != 1 {
		t.Errorf("expected 1 not started, got %d", summary.NotStarted)
	}
	if summary.InProgress != 2 {
		t.Errorf("expected 2 in progress, got %d", summary.InProgress)
	}
	if summary.Deleting != 1 {
		t.Errorf("expected 1 deleting, got %d", summary.Deleting)
	}
	if summary.Other != 0 {
		t.Errorf("expected 0 other, got %d", summary.Other)
	}
	if summary.LastSuccessful == nil {
		t.Fatal("expected last successful timestamp")
	}
	if summary.LastFailed == nil {
		t.Fatal("expected last failed timestamp")
	}
	if !summary.LastFailed.Equal(now.Add(-20 * time.Minute)) {
		t.Errorf("expected last failed to use latest failed-like timestamp, got %v", summary.LastFailed)
	}
}

func TestGeneratePeriodSummaries_MapsFailureVariants(t *testing.T) {
	now := time.Now()
	backups := []collector.BackupInfo{
		{Name: "b1", Phase: "Completed", CompletionTimestamp: timePtr(now.Add(-2 * time.Hour))},
		{Name: "b2", Phase: "FailedValidation", CompletionTimestamp: timePtr(now.Add(-3 * time.Hour))},
		{Name: "b3", Phase: "FinalizingPartiallyFailed", CompletionTimestamp: timePtr(now.Add(-4 * time.Hour))},
	}

	periods := generatePeriodSummaries(backups)
	summary := periods["Last 24 Hours"]

	if summary.TotalBackups != 3 {
		t.Fatalf("expected 3 total backups, got %d", summary.TotalBackups)
	}
	if summary.Completed != 1 {
		t.Fatalf("expected 1 completed, got %d", summary.Completed)
	}
	if summary.Failed != 1 {
		t.Fatalf("expected 1 failed, got %d", summary.Failed)
	}
	if summary.PartiallyFailed != 1 {
		t.Fatalf("expected 1 partially failed, got %d", summary.PartiallyFailed)
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
			ItemsBackedUp:       50,
			TotalItems:          100,
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

func TestGenerateDetails_FailedFirst(t *testing.T) {
	now := time.Now()
	backups := []collector.BackupInfo{
		{Name: "completed-new", Phase: "Completed", StartTimestamp: timePtr(now.Add(-1 * time.Hour))},
		{Name: "failed-old", Phase: "Failed", StartTimestamp: timePtr(now.Add(-5 * time.Hour))},
		{Name: "partial-old", Phase: "PartiallyFailed", StartTimestamp: timePtr(now.Add(-4 * time.Hour))},
	}

	details := generateDetails(backups)
	if len(details) != 3 {
		t.Fatalf("expected 3 details, got %d", len(details))
	}

	if details[0].Status != "Failed" {
		t.Fatalf("expected first detail to be failed, got %s", details[0].Status)
	}
	if details[1].Status != "PartiallyFailed" {
		t.Fatalf("expected second detail to be partially failed, got %s", details[1].Status)
	}
	if details[2].Status != "Completed" {
		t.Fatalf("expected completed detail after failures, got %s", details[2].Status)
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

	summaries := generateScheduleSummaries(backups, schedules, map[string]bool{})

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

func TestGenerateMissedScheduleDetails(t *testing.T) {
	now := time.Date(2026, 4, 14, 10, 58, 44, 0, time.UTC)
	lastRun := time.Date(2026, 4, 12, 5, 0, 0, 0, time.UTC)

	schedules := []collector.ScheduleInfo{
		{Name: "immich", Schedule: "0 5 * * *", LastBackupTime: &lastRun},
	}

	missed, bySchedule := generateMissedScheduleDetails(nil, schedules, now)

	if !bySchedule["immich"] {
		t.Fatal("expected immich schedule to be flagged as missed")
	}
	if len(missed) != 1 {
		t.Fatalf("expected 1 recent missed run, got %d", len(missed))
	}
	if missed[0].Status != "Missed" {
		t.Fatalf("expected status Missed, got %s", missed[0].Status)
	}
	if missed[0].FailureReason == "" {
		t.Fatal("expected failure reason for missed run")
	}
}

func TestGenerateScheduleSummaries_MarkedMissed(t *testing.T) {
	now := time.Now()
	backups := []collector.BackupInfo{
		{Name: "s1-1", ScheduleName: "daily", Phase: "Completed", StartTimestamp: timePtr(now.Add(-2 * time.Hour))},
	}
	schedules := []collector.ScheduleInfo{{Name: "daily"}}

	summaries := generateScheduleSummaries(backups, schedules, map[string]bool{"daily": true})
	if len(summaries) == 0 {
		t.Fatal("expected at least one summary")
	}
	if summaries[0].LastBackupStatus != "Missed" {
		t.Fatalf("expected last backup status Missed, got %s", summaries[0].LastBackupStatus)
	}
}
