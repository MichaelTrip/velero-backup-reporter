package pdf

import (
	"testing"
	"time"

	"github.com/michael/velero-backup-reporter/internal/collector"
)

func TestGenerateBackupReport(t *testing.T) {
	now := time.Now()
	start := now.Add(-10 * time.Minute)
	completion := now

	backup := &collector.BackupInfo{
		Name:                        "test-backup-20240101",
		Namespace:                   "velero",
		Phase:                       "Completed",
		ScheduleName:                "daily",
		StartTimestamp:              &start,
		CompletionTimestamp:         &completion,
		StorageLocation:             "default",
		TTL:                         "720h0m0s",
		IncludedNamespaces:          []string{"app-ns"},
		ExcludedNamespaces:          []string{"kube-system"},
		IncludedResources:           []string{"deployments"},
		ExcludedResources:           []string{"secrets"},
		Labels:                      map[string]string{"env": "prod"},
		Annotations:                 map[string]string{"note": "test"},
		ItemsBackedUp:               100,
		TotalItems:                  100,
		Warnings:                    1,
		Errors:                      0,
		VolumeSnapshotsAttempted:    3,
		VolumeSnapshotsCompleted:    3,
		CSIVolumeSnapshotsAttempted: 1,
		CSIVolumeSnapshotsCompleted: 1,
	}

	data, err := GenerateBackupReport(backup)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("expected non-empty PDF output")
	}
	// Verify it's a valid PDF (starts with %PDF)
	if string(data[:4]) != "%PDF" {
		t.Errorf("expected PDF header, got %q", string(data[:4]))
	}
}

func TestGenerateBackupReport_FailedBackup(t *testing.T) {
	backup := &collector.BackupInfo{
		Name:             "failed-backup",
		Namespace:        "velero",
		Phase:            "Failed",
		FailureReason:    "storage location not found",
		ValidationErrors: []string{"invalid selector", "unknown resource"},
		Warnings:         0,
		Errors:           5,
	}

	data, err := GenerateBackupReport(backup)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("expected non-empty PDF output")
	}
}

func TestGenerateBackupReport_MinimalBackup(t *testing.T) {
	backup := &collector.BackupInfo{
		Name:      "minimal",
		Namespace: "velero",
		Phase:     "New",
	}

	data, err := GenerateBackupReport(backup)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("expected non-empty PDF output")
	}
}
