package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/michael/velero-backup-reporter/internal/collector"
)

// mockCollector creates a collector with pre-populated data for testing.
func mockCollector() *collector.Collector {
	c := collector.New(nil, "velero", 5*time.Minute)
	c.SetData(
		[]collector.BackupInfo{
			{
				Name:                "daily-20240101",
				Namespace:           "velero",
				Phase:               "Completed",
				ScheduleName:        "daily",
				StartTimestamp:      timePtr(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
				CompletionTimestamp: timePtr(time.Date(2024, 1, 1, 0, 5, 0, 0, time.UTC)),
				ItemsBackedUp:      100,
				TotalItems:         100,
				Warnings:           0,
				Errors:             0,
				StorageLocation:    "default",
				TTL:                "720h0m0s",
				HooksAttempted:     3,
				HooksFailed:        0,
				BackupItemOperationsAttempted: 2,
				BackupItemOperationsCompleted: 2,
				BackupItemOperationsFailed:    0,
				FormatVersion:      "1.1.0",
			},
			{
				Name:                "daily-20240102",
				Namespace:           "velero",
				Phase:               "Failed",
				ScheduleName:        "daily",
				StartTimestamp:      timePtr(time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)),
				CompletionTimestamp: timePtr(time.Date(2024, 1, 2, 0, 3, 0, 0, time.UTC)),
				Errors:             3,
				FailureReason:      "something went wrong",
			},
		},
		[]collector.ScheduleInfo{
			{Name: "daily", Schedule: "0 0 * * *", Phase: "Enabled"},
		},
	)
	return c
}

func timePtr(t time.Time) *time.Time {
	return &t
}

func TestHealthz(t *testing.T) {
	c := mockCollector()
	srv, err := New(c)
	if err != nil {
		t.Fatalf("creating server: %v", err)
	}

	req := httptest.NewRequest("GET", "/healthz", nil)
	w := httptest.NewRecorder()
	srv.Handler().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var body map[string]string
	json.NewDecoder(w.Body).Decode(&body)
	if body["status"] != "ok" {
		t.Errorf("expected status 'ok', got '%s'", body["status"])
	}
}

func TestAPIDashboard(t *testing.T) {
	c := mockCollector()
	srv, err := New(c)
	if err != nil {
		t.Fatalf("creating server: %v", err)
	}

	req := httptest.NewRequest("GET", "/api/v1/dashboard", nil)
	w := httptest.NewRecorder()
	srv.Handler().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	ct := w.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", ct)
	}

	var resp dashboardResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decoding response: %v", err)
	}

	if resp.Summary.Total != 2 {
		t.Errorf("expected 2 total backups, got %d", resp.Summary.Total)
	}
	if resp.Summary.Completed != 1 {
		t.Errorf("expected 1 completed, got %d", resp.Summary.Completed)
	}
	if resp.Summary.Failed != 1 {
		t.Errorf("expected 1 failed, got %d", resp.Summary.Failed)
	}
	if len(resp.Schedules) != 1 {
		t.Errorf("expected 1 schedule, got %d", len(resp.Schedules))
	}
	if resp.Schedules[0].ScheduleName != "daily" {
		t.Errorf("expected schedule name 'daily', got '%s'", resp.Schedules[0].ScheduleName)
	}
}

func TestAPIBackups(t *testing.T) {
	c := mockCollector()
	srv, err := New(c)
	if err != nil {
		t.Fatalf("creating server: %v", err)
	}

	req := httptest.NewRequest("GET", "/api/v1/backups", nil)
	w := httptest.NewRecorder()
	srv.Handler().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	ct := w.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", ct)
	}

	var items []backupListItemJSON
	if err := json.NewDecoder(w.Body).Decode(&items); err != nil {
		t.Fatalf("decoding response: %v", err)
	}

	if len(items) != 2 {
		t.Fatalf("expected 2 backups, got %d", len(items))
	}

	// Should be sorted by start time descending (most recent first)
	if items[0].Name != "daily-20240102" {
		t.Errorf("expected first backup to be daily-20240102, got %s", items[0].Name)
	}
}

func TestAPIBackupDetail_Found(t *testing.T) {
	c := mockCollector()
	srv, err := New(c)
	if err != nil {
		t.Fatalf("creating server: %v", err)
	}

	req := httptest.NewRequest("GET", "/api/v1/backups/daily-20240101", nil)
	w := httptest.NewRecorder()
	srv.Handler().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var resp backupDetailJSON
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decoding response: %v", err)
	}

	if resp.Name != "daily-20240101" {
		t.Errorf("expected name 'daily-20240101', got '%s'", resp.Name)
	}
	if resp.Status != "Completed" {
		t.Errorf("expected status 'Completed', got '%s'", resp.Status)
	}
	if resp.Duration != "5m0s" {
		t.Errorf("expected duration '5m0s', got '%s'", resp.Duration)
	}
	if !resp.IsTerminal {
		t.Error("expected isTerminal to be true for Completed backup")
	}
	if resp.HooksAttempted != 3 {
		t.Errorf("expected 3 hooks attempted, got %d", resp.HooksAttempted)
	}
	if resp.HooksFailed != 0 {
		t.Errorf("expected 0 hooks failed, got %d", resp.HooksFailed)
	}
	if resp.BackupItemOperationsAttempted != 2 {
		t.Errorf("expected 2 backup item operations attempted, got %d", resp.BackupItemOperationsAttempted)
	}
	if resp.FormatVersion != "1.1.0" {
		t.Errorf("expected format version '1.1.0', got '%s'", resp.FormatVersion)
	}
	if resp.VolumeBackups == nil {
		t.Error("expected volumeBackups to be non-nil (empty array)")
	}
	if len(resp.VolumeBackups) != 0 {
		t.Errorf("expected 0 volume backups (no kube client), got %d", len(resp.VolumeBackups))
	}
}

func TestAPIBackupDetail_NotFound(t *testing.T) {
	c := mockCollector()
	srv, err := New(c)
	if err != nil {
		t.Fatalf("creating server: %v", err)
	}

	req := httptest.NewRequest("GET", "/api/v1/backups/nonexistent", nil)
	w := httptest.NewRecorder()
	srv.Handler().ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decoding error response: %v", err)
	}

	if resp["error"] != "backup not found" {
		t.Errorf("expected error 'backup not found', got '%s'", resp["error"])
	}
}

func TestAPIBackupLogs_NotFound(t *testing.T) {
	c := mockCollector()
	srv, err := New(c)
	if err != nil {
		t.Fatalf("creating server: %v", err)
	}

	req := httptest.NewRequest("GET", "/api/v1/backups/nonexistent/logs", nil)
	w := httptest.NewRecorder()
	srv.Handler().ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestAPIBackupLogs_NoKubeClient(t *testing.T) {
	c := mockCollector()
	srv, err := New(c) // no WithKubeClient
	if err != nil {
		t.Fatalf("creating server: %v", err)
	}

	req := httptest.NewRequest("GET", "/api/v1/backups/daily-20240101/logs", nil)
	w := httptest.NewRecorder()
	srv.Handler().ServeHTTP(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503, got %d", w.Code)
	}
}

func TestAPIBackupPDF_NotFound(t *testing.T) {
	c := mockCollector()
	srv, err := New(c)
	if err != nil {
		t.Fatalf("creating server: %v", err)
	}

	req := httptest.NewRequest("GET", "/api/v1/backups/nonexistent/pdf", nil)
	w := httptest.NewRecorder()
	srv.Handler().ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestAPIBackupPDF_Found(t *testing.T) {
	c := mockCollector()
	srv, err := New(c)
	if err != nil {
		t.Fatalf("creating server: %v", err)
	}

	req := httptest.NewRequest("GET", "/api/v1/backups/daily-20240101/pdf", nil)
	w := httptest.NewRecorder()
	srv.Handler().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	ct := w.Header().Get("Content-Type")
	if ct != "application/pdf" {
		t.Errorf("expected Content-Type application/pdf, got %s", ct)
	}

	if w.Body.Len() == 0 {
		t.Error("expected non-empty PDF body")
	}
}

func TestAPIDashboard_EmptyData(t *testing.T) {
	c := collector.New(nil, "velero", 5*time.Minute)
	c.SetData(nil, nil)

	srv, err := New(c)
	if err != nil {
		t.Fatalf("creating server: %v", err)
	}

	req := httptest.NewRequest("GET", "/api/v1/dashboard", nil)
	w := httptest.NewRecorder()
	srv.Handler().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var resp dashboardResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decoding response: %v", err)
	}

	if resp.Summary.Total != 0 {
		t.Errorf("expected 0 total, got %d", resp.Summary.Total)
	}
}

func TestAPIBackups_EmptyData(t *testing.T) {
	c := collector.New(nil, "velero", 5*time.Minute)
	c.SetData(nil, nil)

	srv, err := New(c)
	if err != nil {
		t.Fatalf("creating server: %v", err)
	}

	req := httptest.NewRequest("GET", "/api/v1/backups", nil)
	w := httptest.NewRecorder()
	srv.Handler().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var items []backupListItemJSON
	if err := json.NewDecoder(w.Body).Decode(&items); err != nil {
		t.Fatalf("decoding response: %v", err)
	}

	if len(items) != 0 {
		t.Errorf("expected 0 backups, got %d", len(items))
	}
}
