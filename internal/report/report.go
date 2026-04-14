package report

import (
	"sort"
	"strings"
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
	NotStarted      int
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
	Name             string
	ScheduleName     string
	Status           string
	StartTime        *time.Time
	CompletionTime   *time.Time
	Duration         time.Duration
	ItemsBackedUp    int
	TotalItems       int
	Warnings         int
	Errors           int
	FailureReason    string
	ValidationErrors []string
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
		if backupDidNotStart(b) {
			summary.NotStarted++
		}

		switch classifyPhase(b.Phase) {
		case phaseClassCompleted:
			summary.Completed++
			if ts := backupRelevantTimestamp(b); ts != nil {
				if summary.LastSuccessful == nil || ts.After(*summary.LastSuccessful) {
					summary.LastSuccessful = ts
				}
			}
		case phaseClassFailed:
			summary.Failed++
			if ts := backupRelevantTimestamp(b); ts != nil {
				if summary.LastFailed == nil || ts.After(*summary.LastFailed) {
					summary.LastFailed = ts
				}
			}
		case phaseClassPartiallyFailed:
			summary.PartiallyFailed++
		case phaseClassInProgress:
			summary.InProgress++
		case phaseClassDeleting:
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
			Name:             b.Name,
			ScheduleName:     b.ScheduleName,
			Status:           b.Phase,
			StartTime:        b.StartTimestamp,
			CompletionTime:   b.CompletionTimestamp,
			ItemsBackedUp:    b.ItemsBackedUp,
			TotalItems:       b.TotalItems,
			Warnings:         b.Warnings,
			Errors:           b.Errors,
			FailureReason:    b.FailureReason,
			ValidationErrors: b.ValidationErrors,
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

				switch classifyPhase(b.Phase) {
				case phaseClassCompleted:
					period.Completed++
				case phaseClassFailed:
					period.Failed++
				case phaseClassPartiallyFailed:
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

type backupPhaseClass int

const (
	phaseClassOther backupPhaseClass = iota
	phaseClassCompleted
	phaseClassFailed
	phaseClassPartiallyFailed
	phaseClassInProgress
	phaseClassDeleting
)

func classifyPhase(phase string) backupPhaseClass {
	switch {
	case strings.Contains(phase, "PartiallyFailed"):
		return phaseClassPartiallyFailed
	case strings.Contains(phase, "Failed"):
		return phaseClassFailed
	case phase == "Completed":
		return phaseClassCompleted
	case phase == "Deleting":
		return phaseClassDeleting
	case phase == "InProgress" || phase == "New" || phase == "Queued" || phase == "ReadyToStart" || phase == "Finalizing" || phase == "WaitingForPluginOperations":
		return phaseClassInProgress
	default:
		return phaseClassOther
	}
}

func backupRelevantTimestamp(b collector.BackupInfo) *time.Time {
	if b.CompletionTimestamp != nil {
		return b.CompletionTimestamp
	}
	return b.StartTimestamp
}

func backupDidNotStart(b collector.BackupInfo) bool {
	if b.StartTimestamp != nil || b.CompletionTimestamp != nil {
		return false
	}

	switch b.Phase {
	case "New", "Queued", "ReadyToStart", "FailedValidation":
		return true
	default:
		return false
	}
}
