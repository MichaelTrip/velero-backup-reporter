package email

import (
	"context"
	"log"

	"github.com/michael/velero-backup-reporter/internal/collector"
	"github.com/michael/velero-backup-reporter/internal/report"
	"github.com/robfig/cron/v3"
)

// Scheduler manages scheduled email report delivery.
type Scheduler struct {
	sender    *Sender
	collector *collector.Collector
	cron      *cron.Cron
	schedule  string
}

// NewScheduler creates a new email Scheduler.
func NewScheduler(sender *Sender, c *collector.Collector, schedule string) *Scheduler {
	return &Scheduler{
		sender:    sender,
		collector: c,
		cron:      cron.New(),
		schedule:  schedule,
	}
}

// Start begins the email schedule. It blocks until ctx is cancelled.
func (s *Scheduler) Start(ctx context.Context) error {
	_, err := s.cron.AddFunc(s.schedule, func() {
		rpt := report.Generate(s.collector.Backups(), s.collector.Schedules())
		if err := s.sender.Send(rpt); err != nil {
			log.Printf("ERROR: failed to send email report: %v", err)
		}
	})
	if err != nil {
		return err
	}

	s.cron.Start()
	<-ctx.Done()
	s.cron.Stop()
	return nil
}
