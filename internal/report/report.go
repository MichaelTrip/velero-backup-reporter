package report

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/michael/velero-backup-reporter/internal/collector"
	"github.com/robfig/cron/v3"
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

	missedDetails, missedBySchedule := generateMissedScheduleDetails(backups, schedules, report.GeneratedAt)
	report.Backups = append(report.Backups, missedDetails...)
	sortBackupDetails(report.Backups)
	report.Summary.NotStarted += len(missedDetails)

	report.ScheduleSummaries = generateScheduleSummaries(backups, schedules, missedBySchedule)
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

	sortBackupDetails(details)

	return details
}

func sortBackupDetails(details []BackupDetail) {
	// Sort by timestamp (newest first)
	sort.Slice(details, func(i, j int) bool {
		iTime := details[i].StartTime
		if iTime == nil {
			iTime = details[i].CompletionTime
		}
		jTime := details[j].StartTime
		if jTime == nil {
			jTime = details[j].CompletionTime
		}

		if iTime == nil {
			return false
		}
		if jTime == nil {
			return true
		}

		return iTime.After(*jTime)
	})
}

func generateScheduleSummaries(backups []collector.BackupInfo, schedules []collector.ScheduleInfo, missedBySchedule map[string]bool) []ScheduleSummary {
	// Group backups by schedule
	bySchedule := make(map[string][]collector.BackupInfo)
	for _, b := range backups {
		name := b.ScheduleName
		if name == "" {
			name = "Unscheduled"
		}
		bySchedule[name] = append(bySchedule[name], b)
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

		if missedBySchedule[name] {
			s.LastBackupStatus = "Missed"
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

const missedRunsLookback = 24 * time.Hour

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

func generateMissedScheduleDetails(backups []collector.BackupInfo, schedules []collector.ScheduleInfo, now time.Time) ([]BackupDetail, map[string]bool) {
	missedDetails := []BackupDetail{}
	missedBySchedule := make(map[string]bool)

	backupsBySchedule := make(map[string][]collector.BackupInfo)
	for _, b := range backups {
		if b.ScheduleName == "" {
			continue
		}
		backupsBySchedule[b.ScheduleName] = append(backupsBySchedule[b.ScheduleName], b)
	}

	for _, s := range schedules {
		if s.Paused || s.Schedule == "" {
			continue
		}

		anchor := scheduleAnchorTime(s, backupsBySchedule[s.Name])
		if anchor == nil {
			continue
		}

		missedRuns, err := missedRunTimes(s.Schedule, *anchor, now)
		if err != nil || len(missedRuns) == 0 {
			continue
		}

		hasRecentMissed := false
		for _, run := range missedRuns {
			if run.Before(now.Add(-missedRunsLookback)) {
				continue
			}
			hasRecentMissed = true
			runTime := run
			missedDetails = append(missedDetails, BackupDetail{
				Name:           fmt.Sprintf("%s (missed)", s.Name),
				ScheduleName:   s.Name,
				Status:         "Missed",
				StartTime:      &runTime,
				FailureReason:  "Scheduled run did not create a backup",
				ItemsBackedUp:  0,
				TotalItems:     0,
				Warnings:       0,
				Errors:         0,
				CompletionTime: nil,
			})
		}
		if hasRecentMissed {
			missedBySchedule[s.Name] = true
		}
	}

	return missedDetails, missedBySchedule
}

func scheduleAnchorTime(s collector.ScheduleInfo, backups []collector.BackupInfo) *time.Time {
	if s.LastBackupTime != nil {
		return s.LastBackupTime
	}

	var latest *time.Time
	for _, b := range backups {
		ts := backupRelevantTimestamp(b)
		if ts == nil {
			continue
		}
		if latest == nil || ts.After(*latest) {
			t := *ts
			latest = &t
		}
	}

	return latest
}

func missedRunTimes(scheduleExpr string, anchor, now time.Time) ([]time.Time, error) {
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	spec, err := parser.Parse(scheduleExpr)
	if err != nil {
		return nil, err
	}

	next := spec.Next(anchor)
	if next.After(now) {
		return nil, nil
	}

	runs := make([]time.Time, 0, 4)
	for !next.After(now) {
		runs = append(runs, next)
		if len(runs) >= 32 {
			break
		}
		next = spec.Next(next)
	}

	return runs, nil
}
