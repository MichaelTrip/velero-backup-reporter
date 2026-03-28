package report

import (
	"sort"
	"time"

	"github.com/michael/velero-backup-reporter/internal/collector"
)

// BackupReport is the full report containing summary and details.
type BackupReport struct {
	GeneratedAt       time.Time
	Summary           BackupSummary
	ScheduleSummaries []ScheduleSummary
	Backups           []BackupDetail
	PeriodSummaries   map[string]BackupPeriodSummary
}

// BackupSummary contains aggregate stats across all backups.
type BackupSummary struct {
	TotalBackups    int
	Completed       int
	Failed          int
	PartiallyFailed int
	InProgress      int
	Deleting        int
	Other           int
	LastSuccessful  *time.Time
	LastFailed      *time.Time
}

// BackupPeriodSummary contains statistics for a time period.
type BackupPeriodSummary struct {
	Period          string
	TotalBackups    int
	Completed       int
	Failed          int
	PartiallyFailed int
	AverageDuration time.Duration
	TotalItems      int
}

// ScheduleSummary contains per-schedule statistics.
type ScheduleSummary struct {
	ScheduleName      string
	LastBackupStatus  string
	LastBackupTime    *time.Time
	TotalBackups      int
	SuccessfulBackups int
	SuccessRate       float64
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
	report.PeriodSummaries = generatePeriodSummaries(backups)

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

	// Sort by timestamp (newest first)
	sort.Slice(details, func(i, j int) bool {
		// Handle nil timestamps
		if details[i].StartTime == nil {
			return false
		}
		if details[j].StartTime == nil {
			return true
		}
		return details[i].StartTime.After(*details[j].StartTime)
	})

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
			ScheduleName: name,
			TotalBackups: len(sBackups),
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

// generatePeriodSummaries creates summary statistics for different time periods.
func generatePeriodSummaries(backups []collector.BackupInfo) map[string]BackupPeriodSummary {
	now := time.Now()
	periods := map[string]time.Time{
		"Last 24 Hours": now.Add(-24 * time.Hour),
		"Last 7 Days":   now.Add(-7 * 24 * time.Hour),
		"Last 30 Days":  now.Add(-30 * 24 * time.Hour),
	}

	summaries := make(map[string]BackupPeriodSummary)

	for periodName, startTime := range periods {
		period := BackupPeriodSummary{
			Period: periodName,
		}

		durations := []time.Duration{}

		for _, b := range backups {
			// Check if backup completed within the period
			completionTime := b.CompletionTimestamp
			if completionTime == nil {
				completionTime = b.StartTimestamp
			}

			if completionTime != nil && completionTime.After(startTime) {
				period.TotalBackups++
				period.TotalItems += b.TotalItems

				switch b.Phase {
				case "Completed":
					period.Completed++
				case "Failed":
					period.Failed++
				case "PartiallyFailed":
					period.PartiallyFailed++
				}

				// Track duration for average calculation
				if b.StartTimestamp != nil && b.CompletionTimestamp != nil {
					durations = append(durations, b.CompletionTimestamp.Sub(*b.StartTimestamp))
				}
			}
		}

		// Calculate average duration
		if len(durations) > 0 {
			totalDuration := time.Duration(0)
			for _, d := range durations {
				totalDuration += d
			}
			period.AverageDuration = totalDuration / time.Duration(len(durations))
		}

		summaries[periodName] = period
	}

	return summaries
}
