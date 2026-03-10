package report

import (
	"time"

	"github.com/michael/velero-backup-reporter/internal/collector"
)

// BackupReport is the full report containing summary and details.
type BackupReport struct {
	GeneratedAt      time.Time
	Summary          BackupSummary
	ScheduleSummaries []ScheduleSummary
	Backups          []BackupDetail
}

// BackupSummary contains aggregate stats across all backups.
type BackupSummary struct {
	TotalBackups      int
	Completed         int
	Failed            int
	PartiallyFailed   int
	InProgress        int
	Deleting          int
	Other             int
	LastSuccessful    *time.Time
	LastFailed        *time.Time
}

// ScheduleSummary contains per-schedule statistics.
type ScheduleSummary struct {
	ScheduleName     string
	LastBackupStatus string
	LastBackupTime   *time.Time
	TotalBackups     int
	SuccessfulBackups int
	SuccessRate      float64
}

// BackupDetail contains information about a single backup.
type BackupDetail struct {
	Name           string
	ScheduleName   string
	Status         string
	StartTime      *time.Time
	CompletionTime *time.Time
	Duration       time.Duration
	ItemsBackedUp  int
	TotalItems     int
	Warnings       int
	Errors         int
}

// Generate creates a BackupReport from collected backup and schedule data.
func Generate(backups []collector.BackupInfo, schedules []collector.ScheduleInfo) BackupReport {
	report := BackupReport{
		GeneratedAt: time.Now(),
	}

	report.Summary = generateSummary(backups)
	report.Backups = generateDetails(backups)
	report.ScheduleSummaries = generateScheduleSummaries(backups, schedules)

	return report
}

func generateSummary(backups []collector.BackupInfo) BackupSummary {
	summary := BackupSummary{
		TotalBackups: len(backups),
	}

	for _, b := range backups {
		switch b.Phase {
		case "Completed":
			summary.Completed++
			if summary.LastSuccessful == nil || (b.CompletionTimestamp != nil && b.CompletionTimestamp.After(*summary.LastSuccessful)) {
				summary.LastSuccessful = b.CompletionTimestamp
			}
		case "Failed":
			summary.Failed++
			if summary.LastFailed == nil || (b.CompletionTimestamp != nil && b.CompletionTimestamp.After(*summary.LastFailed)) {
				summary.LastFailed = b.CompletionTimestamp
			}
		case "PartiallyFailed":
			summary.PartiallyFailed++
		case "InProgress":
			summary.InProgress++
		case "Deleting":
			summary.Deleting++
		default:
			summary.Other++
		}
	}

	return summary
}

func generateDetails(backups []collector.BackupInfo) []BackupDetail {
	details := make([]BackupDetail, 0, len(backups))

	for _, b := range backups {
		detail := BackupDetail{
			Name:           b.Name,
			ScheduleName:   b.ScheduleName,
			Status:         b.Phase,
			StartTime:      b.StartTimestamp,
			CompletionTime: b.CompletionTimestamp,
			ItemsBackedUp:  b.ItemsBackedUp,
			TotalItems:     b.TotalItems,
			Warnings:       b.Warnings,
			Errors:         b.Errors,
		}

		if b.StartTimestamp != nil && b.CompletionTimestamp != nil {
			detail.Duration = b.CompletionTimestamp.Sub(*b.StartTimestamp)
		}

		details = append(details, detail)
	}

	return details
}

func generateScheduleSummaries(backups []collector.BackupInfo, schedules []collector.ScheduleInfo) []ScheduleSummary {
	// Group backups by schedule
	bySchedule := make(map[string][]collector.BackupInfo)
	for _, b := range backups {
		name := b.ScheduleName
		if name == "" {
			name = "Unscheduled"
		}
		bySchedule[name] = append(bySchedule[name], b)
	}

	// Build schedule name set
	scheduleNames := make(map[string]bool)
	for _, s := range schedules {
		scheduleNames[s.Name] = true
	}

	// Ensure all known schedules appear even if they have no backups
	for _, s := range schedules {
		if _, exists := bySchedule[s.Name]; !exists {
			bySchedule[s.Name] = nil
		}
	}

	summaries := make([]ScheduleSummary, 0, len(bySchedule))
	for name, sBackups := range bySchedule {
		s := ScheduleSummary{
			ScheduleName:  name,
			TotalBackups:  len(sBackups),
		}

		var lastTime *time.Time
		for _, b := range sBackups {
			if b.Phase == "Completed" {
				s.SuccessfulBackups++
			}
			if b.StartTimestamp != nil && (lastTime == nil || b.StartTimestamp.After(*lastTime)) {
				lastTime = b.StartTimestamp
				s.LastBackupStatus = b.Phase
				s.LastBackupTime = b.StartTimestamp
			}
		}

		if s.TotalBackups > 0 {
			s.SuccessRate = float64(s.SuccessfulBackups) / float64(s.TotalBackups) * 100
		}

		summaries = append(summaries, s)
	}

	return summaries
}
