package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/michael/velero-backup-reporter/internal/collector"
	"github.com/michael/velero-backup-reporter/internal/config"
	"github.com/michael/velero-backup-reporter/internal/email"
	"github.com/michael/velero-backup-reporter/internal/report"
	"github.com/michael/velero-backup-reporter/internal/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var version = "dev"

func main() {
	rootCmd := &cobra.Command{
		Use:     "velero-backup-reporter",
		Short:   "Backup report generator and web UI for Velero",
		Version: version,
		RunE:    run,
	}

	rootCmd.Flags().String("config", "", "Path to config file (YAML)")
	rootCmd.Flags().String("kubeconfig", "", "Path to kubeconfig file")
	rootCmd.Flags().String("namespace", "velero", "Namespace to monitor for Velero resources")
	rootCmd.Flags().Int("port", 8080, "HTTP server port")
	rootCmd.Flags().String("collection-interval", "5m", "Interval between backup data collections")

	rootCmd.Flags().String("smtp-host", "", "SMTP server host")
	rootCmd.Flags().Int("smtp-port", 587, "SMTP server port")
	rootCmd.Flags().String("smtp-username", "", "SMTP username")
	rootCmd.Flags().String("smtp-password", "", "SMTP password")
	rootCmd.Flags().String("smtp-from", "", "Email sender address")
	rootCmd.Flags().StringSlice("smtp-to", nil, "Email recipient addresses")
	rootCmd.Flags().Bool("smtp-tls", true, "Enable SMTP TLS")
	rootCmd.Flags().String("email-schedule", "0 8 * * *", "Cron schedule for email reports")
	rootCmd.Flags().String("email-details-window", "24h", "Time window for backups shown in email report details (Go duration, e.g. 24h)")
	rootCmd.Flags().Bool("email-enabled", false, "Enable email notifications")
	rootCmd.Flags().Bool("email-test-enabled", false, "Enable the test email endpoint and UI button")
	rootCmd.Flags().Bool("send-report-now", false, "Generate and send backup report immediately, then exit")

	viper.BindPFlags(rootCmd.Flags())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	// Handle on-demand report generation
	if cfg.Email.SendReportNow {
		return sendReportNow(cfg)
	}

	// Normal operation: start web server and collector
	// Set up context with signal handling for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		log.Println("INFO: shutdown signal received")
		cancel()
	}()

	// Initialize Kubernetes client
	kubeClient, err := collector.NewKubeClient(cfg.Kubeconfig)
	if err != nil {
		return fmt.Errorf("creating kubernetes client: %w", err)
	}

	// Initialize collector
	coll := collector.New(kubeClient, cfg.Namespace, cfg.CollectionInterval)

	// Start collector in background
	go coll.Run(ctx)

	// Initialize email sender if enabled
	serverOpts := []server.Option{server.WithKubeClient(kubeClient)}
	if cfg.Email.Enabled {
		sender, err := email.NewSender(cfg.SMTP, cfg.Email)
		if err != nil {
			return fmt.Errorf("creating email sender: %w", err)
		}

		if cfg.Email.TestEnabled {
			serverOpts = append(serverOpts, server.WithEmailSender(sender))
			log.Println("INFO: email test endpoint enabled")
		}

		scheduler := email.NewScheduler(sender, coll, cfg.Email.Schedule)
		go func() {
			if err := scheduler.Start(ctx); err != nil {
				log.Printf("ERROR: email scheduler: %v", err)
			}
		}()
		log.Printf("INFO: email notifications enabled, schedule: %s, details window: %s", cfg.Email.Schedule, cfg.Email.DetailsWindow)
	}

	// Initialize web server
	srv, err := server.New(coll, serverOpts...)
	if err != nil {
		return fmt.Errorf("creating server: %w", err)
	}

	// Start HTTP server
	addr := fmt.Sprintf(":%d", cfg.Port)
	httpServer := &http.Server{
		Addr:    addr,
		Handler: srv.Handler(),
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			log.Printf("ERROR: HTTP server shutdown: %v", err)
		}
	}()

	log.Printf("INFO: starting server on %s", addr)
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("HTTP server: %w", err)
	}

	log.Println("INFO: server stopped")
	return nil
}

// sendReportNow generates and sends a backup report immediately.
func sendReportNow(cfg *config.Config) error {
	// Validate SMTP configuration
	if cfg.SMTP.Host == "" {
		return fmt.Errorf("SMTP configuration required: smtp-host is missing")
	}
	if cfg.SMTP.From == "" {
		return fmt.Errorf("SMTP configuration required: smtp-from is missing")
	}
	if len(cfg.SMTP.To) == 0 {
		return fmt.Errorf("SMTP configuration required: smtp-to is required (at least one recipient)")
	}

	// Set up context
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Initialize Kubernetes client
	kubeClient, err := collector.NewKubeClient(cfg.Kubeconfig)
	if err != nil {
		return fmt.Errorf("creating kubernetes client: %w", err)
	}

	// Initialize collector
	coll := collector.New(kubeClient, cfg.Namespace, cfg.CollectionInterval)

	// Collect backup and schedule data
	log.Println("INFO: collecting backup data...")
	if err := coll.Collect(ctx); err != nil {
		return fmt.Errorf("collecting backup data: %w", err)
	}

	// Create email sender
	sender, err := email.NewSender(cfg.SMTP, cfg.Email)
	if err != nil {
		return fmt.Errorf("creating email sender: %w", err)
	}

	// Generate report
	log.Println("INFO: generating backup report...")
	rpt := report.Generate(coll.Backups(), coll.Schedules())

	// Send email
	log.Printf("INFO: sending report to %v...", cfg.SMTP.To)
	if err := sender.Send(rpt); err != nil {
		return fmt.Errorf("sending report: %w", err)
	}

	log.Printf("INFO: backup report sent successfully to %v", cfg.SMTP.To)
	return nil
}
