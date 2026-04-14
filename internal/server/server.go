package server

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/michael/velero-backup-reporter/internal/collector"
	"github.com/michael/velero-backup-reporter/internal/email"
	"github.com/michael/velero-backup-reporter/internal/logs"
	"github.com/michael/velero-backup-reporter/internal/pdf"
	"github.com/michael/velero-backup-reporter/internal/report"
	"github.com/michael/velero-backup-reporter/web"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Server struct {
	collector   *collector.Collector
	kubeClient  client.Client
	emailSender *email.Sender
	router      chi.Router
}

func New(c *collector.Collector, opts ...Option) (*Server, error) {
	s := &Server{
		collector: c,
	}
	for _, opt := range opts {
		opt(s)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/healthz", s.handleHealthz)

	// JSON API routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(jsonContentType)
		r.Get("/dashboard", s.handleAPIDashboard)
		r.Get("/report", s.handleAPIReport)
		r.Get("/backups", s.handleAPIBackups)
		r.Get("/backups/{name}", s.handleAPIBackupDetail)
		r.Get("/backups/{name}/logs", s.handleAPIBackupLogs)
		r.Get("/backups/{name}/pdf", s.handleAPIBackupPDF)
		r.Post("/email/test", s.handleAPIEmailTest)
	})

	// Serve SPA static files and fallback
	s.setupSPA(r)

	s.router = r
	return s, nil
}

func jsonContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func (s *Server) setupSPA(r chi.Router) {
	distFS, err := fs.Sub(web.DistFS, ".")
	if err != nil {
		log.Printf("WARNING: could not set up SPA filesystem: %v", err)
		return
	}

	fileServer := http.FileServer(http.FS(distFS))

	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		// Try to serve the file directly first
		path := r.URL.Path
		if path == "/" {
			path = "/index.html"
		}

		// Check if the file exists in the embedded FS
		f, err := distFS.(fs.ReadFileFS).ReadFile(path[1:]) // strip leading /
		if err == nil {
			// Serve the actual file
			_ = f
			fileServer.ServeHTTP(w, r)
			return
		}

		// Fallback: serve index.html for SPA routing
		r.URL.Path = "/"
		fileServer.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("ERROR: encoding JSON response: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func (s *Server) Handler() http.Handler {
	return s.router
}

// Option configures the Server.
type Option func(*Server)

// WithKubeClient sets the Kubernetes client for operations like log retrieval.
func WithKubeClient(c client.Client) Option {
	return func(s *Server) {
		s.kubeClient = c
	}
}

// WithEmailSender sets the email sender for on-demand test emails.
func WithEmailSender(sender *email.Sender) Option {
	return func(s *Server) {
		s.emailSender = sender
	}
}

func (s *Server) handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// API response types

type dashboardResponse struct {
	Summary      backupSummaryJSON     `json:"summary"`
	Schedules    []scheduleSummaryJSON `json:"schedules"`
	EmailEnabled bool                  `json:"emailEnabled"`
}

type backupSummaryJSON struct {
	Total           int        `json:"total"`
	Completed       int        `json:"completed"`
	Failed          int        `json:"failed"`
	PartiallyFailed int        `json:"partiallyFailed"`
	NotStarted      int        `json:"notStarted"`
	InProgress      int        `json:"inProgress"`
	Deleting        int        `json:"deleting"`
	LastSuccessful  *time.Time `json:"lastSuccessful"`
	LastFailed      *time.Time `json:"lastFailed"`
}

type scheduleSummaryJSON struct {
	ScheduleName      string     `json:"scheduleName"`
	TotalBackups      int        `json:"totalBackups"`
	SuccessfulBackups int        `json:"successfulBackups"`
	FailedBackups     int        `json:"failedBackups"`
	LastBackupTime    *time.Time `json:"lastBackupTime"`
	LastBackupStatus  string     `json:"lastBackupStatus"`
	SuccessRate       float64    `json:"successRate"`
}

type backupListItemJSON struct {
	Name                string     `json:"name"`
	ScheduleName        string     `json:"scheduleName"`
	Status              string     `json:"status"`
	StartTimestamp      *time.Time `json:"startTimestamp"`
	CompletionTimestamp *time.Time `json:"completionTimestamp"`
	Duration            string     `json:"duration"`
	ItemsBackedUp       int        `json:"itemsBackedUp"`
	Warnings            int        `json:"warnings"`
	Errors              int        `json:"errors"`
}

type backupDetailJSON struct {
	Name                        string            `json:"name"`
	Namespace                   string            `json:"namespace"`
	Status                      string            `json:"status"`
	ScheduleName                string            `json:"scheduleName"`
	StartTimestamp              *time.Time        `json:"startTimestamp"`
	CompletionTimestamp         *time.Time        `json:"completionTimestamp"`
	Expiration                  *time.Time        `json:"expiration"`
	Duration                    string            `json:"duration"`
	TTL                         string            `json:"ttl"`
	StorageLocation             string            `json:"storageLocation"`
	ItemsBackedUp               int               `json:"itemsBackedUp"`
	TotalItems                  int               `json:"totalItems"`
	Warnings                    int               `json:"warnings"`
	Errors                      int               `json:"errors"`
	IncludedNamespaces          []string          `json:"includedNamespaces"`
	ExcludedNamespaces          []string          `json:"excludedNamespaces"`
	IncludedResources           []string          `json:"includedResources"`
	ExcludedResources           []string          `json:"excludedResources"`
	Labels                      map[string]string `json:"labels"`
	Annotations                 map[string]string `json:"annotations"`
	VolumeSnapshotsAttempted    int               `json:"volumeSnapshotsAttempted"`
	VolumeSnapshotsCompleted    int               `json:"volumeSnapshotsCompleted"`
	CSIVolumeSnapshotsAttempted int               `json:"csiVolumeSnapshotsAttempted"`
	CSIVolumeSnapshotsCompleted int               `json:"csiVolumeSnapshotsCompleted"`
	FailureReason               string            `json:"failureReason"`
	ValidationErrors            []string          `json:"validationErrors"`
	IsTerminal                  bool              `json:"isTerminal"`

	// Hook and operation status
	HooksAttempted                int    `json:"hooksAttempted"`
	HooksFailed                   int    `json:"hooksFailed"`
	BackupItemOperationsAttempted int    `json:"backupItemOperationsAttempted"`
	BackupItemOperationsCompleted int    `json:"backupItemOperationsCompleted"`
	BackupItemOperationsFailed    int    `json:"backupItemOperationsFailed"`
	FormatVersion                 string `json:"formatVersion"`

	// Volume backups
	VolumeBackups []volumeBackupJSON `json:"volumeBackups"`
}

type volumeBackupJSON struct {
	VolumeName          string     `json:"volumeName"`
	PodName             string     `json:"podName"`
	PodNamespace        string     `json:"podNamespace"`
	NodeName            string     `json:"nodeName"`
	UploaderType        string     `json:"uploaderType"`
	Phase               string     `json:"phase"`
	StartTimestamp      *time.Time `json:"startTimestamp"`
	CompletionTimestamp *time.Time `json:"completionTimestamp"`
	TotalBytes          int64      `json:"totalBytes"`
	BytesDone           int64      `json:"bytesDone"`
	SnapshotID          string     `json:"snapshotId"`
}

func (s *Server) handleAPIDashboard(w http.ResponseWriter, r *http.Request) {
	rpt := report.Generate(s.collector.Backups(), s.collector.Schedules())

	sort.Slice(rpt.ScheduleSummaries, func(i, j int) bool {
		return rpt.ScheduleSummaries[i].ScheduleName < rpt.ScheduleSummaries[j].ScheduleName
	})

	schedules := make([]scheduleSummaryJSON, 0, len(rpt.ScheduleSummaries))
	for _, ss := range rpt.ScheduleSummaries {
		failed := ss.TotalBackups - ss.SuccessfulBackups
		schedules = append(schedules, scheduleSummaryJSON{
			ScheduleName:      ss.ScheduleName,
			TotalBackups:      ss.TotalBackups,
			SuccessfulBackups: ss.SuccessfulBackups,
			FailedBackups:     failed,
			LastBackupTime:    ss.LastBackupTime,
			LastBackupStatus:  ss.LastBackupStatus,
			SuccessRate:       ss.SuccessRate,
		})
	}

	resp := dashboardResponse{
		Summary: backupSummaryJSON{
			Total:           rpt.Summary.TotalBackups,
			Completed:       rpt.Summary.Completed,
			Failed:          rpt.Summary.Failed,
			PartiallyFailed: rpt.Summary.PartiallyFailed,
			NotStarted:      rpt.Summary.NotStarted,
			InProgress:      rpt.Summary.InProgress,
			Deleting:        rpt.Summary.Deleting,
			LastSuccessful:  rpt.Summary.LastSuccessful,
			LastFailed:      rpt.Summary.LastFailed,
		},
		Schedules:    schedules,
		EmailEnabled: s.emailSender != nil,
	}

	writeJSON(w, http.StatusOK, resp)
}

// reportResponse represents the complete backup report
type reportResponse struct {
	GeneratedAt     time.Time                    `json:"generatedAt"`
	Summary         backupSummaryJSON            `json:"summary"`
	Schedules       []scheduleSummaryJSON        `json:"schedules"`
	PeriodSummaries map[string]periodSummaryJSON `json:"periodSummaries"`
	Backups         []backupListItemJSON         `json:"backups"`
}

type periodSummaryJSON struct {
	Period          string `json:"period"`
	TotalBackups    int    `json:"totalBackups"`
	Completed       int    `json:"completed"`
	Failed          int    `json:"failed"`
	PartiallyFailed int    `json:"partiallyFailed"`
	AverageDuration string `json:"averageDuration"`
	TotalItems      int    `json:"totalItems"`
}

func (s *Server) handleAPIReport(w http.ResponseWriter, r *http.Request) {
	rpt := report.Generate(s.collector.Backups(), s.collector.Schedules())

	// Sort schedules by name
	sort.Slice(rpt.ScheduleSummaries, func(i, j int) bool {
		return rpt.ScheduleSummaries[i].ScheduleName < rpt.ScheduleSummaries[j].ScheduleName
	})

	schedules := make([]scheduleSummaryJSON, 0, len(rpt.ScheduleSummaries))
	for _, ss := range rpt.ScheduleSummaries {
		failed := ss.TotalBackups - ss.SuccessfulBackups
		schedules = append(schedules, scheduleSummaryJSON{
			ScheduleName:      ss.ScheduleName,
			TotalBackups:      ss.TotalBackups,
			SuccessfulBackups: ss.SuccessfulBackups,
			FailedBackups:     failed,
			LastBackupTime:    ss.LastBackupTime,
			LastBackupStatus:  ss.LastBackupStatus,
			SuccessRate:       ss.SuccessRate,
		})
	}

	// Convert period summaries
	periodSummaries := make(map[string]periodSummaryJSON)
	for periodName, ps := range rpt.PeriodSummaries {
		periodSummaries[periodName] = periodSummaryJSON{
			Period:          ps.Period,
			TotalBackups:    ps.TotalBackups,
			Completed:       ps.Completed,
			Failed:          ps.Failed,
			PartiallyFailed: ps.PartiallyFailed,
			AverageDuration: formatDuration(ps.AverageDuration),
			TotalItems:      ps.TotalItems,
		}
	}

	// Convert backup details to JSON format (already sorted by date in report.Generate)
	backups := make([]backupListItemJSON, 0, len(rpt.Backups))
	for _, b := range rpt.Backups {
		backups = append(backups, backupListItemJSON{
			Name:                b.Name,
			ScheduleName:        b.ScheduleName,
			Status:              b.Status,
			StartTimestamp:      b.StartTime,
			CompletionTimestamp: b.CompletionTime,
			Duration:            formatDuration(b.Duration),
			ItemsBackedUp:       b.ItemsBackedUp,
			Warnings:            b.Warnings,
			Errors:              b.Errors,
		})
	}

	resp := reportResponse{
		GeneratedAt: rpt.GeneratedAt,
		Summary: backupSummaryJSON{
			Total:           rpt.Summary.TotalBackups,
			Completed:       rpt.Summary.Completed,
			Failed:          rpt.Summary.Failed,
			PartiallyFailed: rpt.Summary.PartiallyFailed,
			NotStarted:      rpt.Summary.NotStarted,
			InProgress:      rpt.Summary.InProgress,
			Deleting:        rpt.Summary.Deleting,
			LastSuccessful:  rpt.Summary.LastSuccessful,
			LastFailed:      rpt.Summary.LastFailed,
		},
		Schedules:       schedules,
		PeriodSummaries: periodSummaries,
		Backups:         backups,
	}

	writeJSON(w, http.StatusOK, resp)
}

func formatDuration(d time.Duration) string {
	if d == 0 {
		return ""
	}
	if d < time.Minute {
		return d.Round(time.Second).String()
	}
	minutes := int(d.Minutes())
	seconds := int(d.Seconds()) % 60
	if minutes < 60 {
		return fmt.Sprintf("%dm%ds", minutes, seconds)
	}
	hours := minutes / 60
	minutes = minutes % 60
	return fmt.Sprintf("%dh%dm", hours, minutes)
}

func (s *Server) handleAPIBackups(w http.ResponseWriter, r *http.Request) {
	rpt := report.Generate(s.collector.Backups(), s.collector.Schedules())

	sort.Slice(rpt.Backups, func(i, j int) bool {
		if rpt.Backups[i].StartTime == nil {
			return false
		}
		if rpt.Backups[j].StartTime == nil {
			return true
		}
		return rpt.Backups[i].StartTime.After(*rpt.Backups[j].StartTime)
	})

	items := make([]backupListItemJSON, 0, len(rpt.Backups))
	for _, b := range rpt.Backups {
		items = append(items, backupListItemJSON{
			Name:                b.Name,
			ScheduleName:        b.ScheduleName,
			Status:              b.Status,
			StartTimestamp:      b.StartTime,
			CompletionTimestamp: b.CompletionTime,
			Duration:            formatDuration(b.Duration),
			ItemsBackedUp:       b.ItemsBackedUp,
			Warnings:            b.Warnings,
			Errors:              b.Errors,
		})
	}

	writeJSON(w, http.StatusOK, items)
}

func isTerminalPhase(phase string) bool {
	switch phase {
	case "Completed", "PartiallyFailed", "Failed":
		return true
	}
	return false
}

func (s *Server) handleAPIBackupDetail(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	backup := s.collector.GetBackup(name)
	if backup == nil {
		writeError(w, http.StatusNotFound, "backup not found")
		return
	}

	var duration time.Duration
	if backup.StartTimestamp != nil && backup.CompletionTimestamp != nil {
		duration = backup.CompletionTimestamp.Sub(*backup.StartTimestamp)
	}

	// Fetch volume backups on-demand if kube client is available
	var volumeBackups []volumeBackupJSON
	if s.kubeClient != nil {
		vbs := collector.ListVolumeBackups(r.Context(), s.kubeClient, name, backup.Namespace)
		volumeBackups = make([]volumeBackupJSON, 0, len(vbs))
		for _, vb := range vbs {
			volumeBackups = append(volumeBackups, volumeBackupJSON{
				VolumeName:          vb.VolumeName,
				PodName:             vb.PodName,
				PodNamespace:        vb.PodNamespace,
				NodeName:            vb.NodeName,
				UploaderType:        vb.UploaderType,
				Phase:               vb.Phase,
				StartTimestamp:      vb.StartTimestamp,
				CompletionTimestamp: vb.CompletionTimestamp,
				TotalBytes:          vb.TotalBytes,
				BytesDone:           vb.BytesDone,
				SnapshotID:          vb.SnapshotID,
			})
		}
	}
	if volumeBackups == nil {
		volumeBackups = []volumeBackupJSON{}
	}

	resp := backupDetailJSON{
		Name:                          backup.Name,
		Namespace:                     backup.Namespace,
		Status:                        backup.Phase,
		ScheduleName:                  backup.ScheduleName,
		StartTimestamp:                backup.StartTimestamp,
		CompletionTimestamp:           backup.CompletionTimestamp,
		Expiration:                    backup.Expiration,
		Duration:                      formatDuration(duration),
		TTL:                           backup.TTL,
		StorageLocation:               backup.StorageLocation,
		ItemsBackedUp:                 backup.ItemsBackedUp,
		TotalItems:                    backup.TotalItems,
		Warnings:                      backup.Warnings,
		Errors:                        backup.Errors,
		IncludedNamespaces:            backup.IncludedNamespaces,
		ExcludedNamespaces:            backup.ExcludedNamespaces,
		IncludedResources:             backup.IncludedResources,
		ExcludedResources:             backup.ExcludedResources,
		Labels:                        backup.Labels,
		Annotations:                   backup.Annotations,
		VolumeSnapshotsAttempted:      backup.VolumeSnapshotsAttempted,
		VolumeSnapshotsCompleted:      backup.VolumeSnapshotsCompleted,
		CSIVolumeSnapshotsAttempted:   backup.CSIVolumeSnapshotsAttempted,
		CSIVolumeSnapshotsCompleted:   backup.CSIVolumeSnapshotsCompleted,
		FailureReason:                 backup.FailureReason,
		ValidationErrors:              backup.ValidationErrors,
		IsTerminal:                    isTerminalPhase(backup.Phase),
		HooksAttempted:                backup.HooksAttempted,
		HooksFailed:                   backup.HooksFailed,
		BackupItemOperationsAttempted: backup.BackupItemOperationsAttempted,
		BackupItemOperationsCompleted: backup.BackupItemOperationsCompleted,
		BackupItemOperationsFailed:    backup.BackupItemOperationsFailed,
		FormatVersion:                 backup.FormatVersion,
		VolumeBackups:                 volumeBackups,
	}

	writeJSON(w, http.StatusOK, resp)
}

func (s *Server) handleAPIBackupLogs(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	backup := s.collector.GetBackup(name)
	if backup == nil {
		writeError(w, http.StatusNotFound, "backup not found")
		return
	}

	if !isTerminalPhase(backup.Phase) {
		writeError(w, http.StatusBadRequest, "logs are not available for backups in non-terminal phase")
		return
	}

	if s.kubeClient == nil {
		writeError(w, http.StatusServiceUnavailable, "log retrieval is not available (no Kubernetes client configured)")
		return
	}

	logContent, err := logs.FetchBackupLogs(r.Context(), s.kubeClient, name, backup.Namespace)
	if err != nil {
		log.Printf("ERROR: fetching backup logs for %s: %v", name, err)
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("failed to retrieve logs: %v", err))
		return
	}

	// Logs endpoint returns plain text, override the JSON content type
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, logContent)
}

func (s *Server) handleAPIBackupPDF(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	backup := s.collector.GetBackup(name)
	if backup == nil {
		writeError(w, http.StatusNotFound, "backup not found")
		return
	}

	pdfBytes, err := pdf.GenerateBackupReport(backup)
	if err != nil {
		log.Printf("ERROR: generating PDF for %s: %v", name, err)
		writeError(w, http.StatusInternalServerError, "failed to generate PDF")
		return
	}

	// PDF endpoint returns binary, override the JSON content type
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"backup-%s.pdf\"", name))
	w.Write(pdfBytes)
}

func (s *Server) handleAPIEmailTest(w http.ResponseWriter, r *http.Request) {
	if s.emailSender == nil {
		writeError(w, http.StatusNotImplemented, "email notifications are not enabled")
		return
	}

	rpt := report.Generate(s.collector.Backups(), s.collector.Schedules())
	if err := s.emailSender.Send(rpt); err != nil {
		log.Printf("ERROR: sending test email: %v", err)
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("failed to send test email: %v", err))
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Test email sent successfully"})
}
