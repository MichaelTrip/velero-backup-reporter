package collector

import (
	"testing"
	"time"

	velerov1api "github.com/vmware-tanzu/velero/pkg/apis/velero/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestExtractBackupInfo_Basic(t *testing.T) {
	now := time.Now()
	start := metav1.NewTime(now.Add(-10 * time.Minute))
	completion := metav1.NewTime(now)

	backup := velerov1api.Backup{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-backup",
			Namespace: "velero",
			Labels: map[string]string{
				"velero.io/schedule-name": "daily-schedule",
			},
		},
		Status: velerov1api.BackupStatus{
			Phase:               velerov1api.BackupPhaseCompleted,
			StartTimestamp:      &start,
			CompletionTimestamp: &completion,
			Warnings:            2,
			Errors:              0,
			Progress: &velerov1api.BackupProgress{
				TotalItems:    100,
				ItemsBackedUp: 100,
			},
		},
	}

	info := extractBackupInfo(backup)

	if info.Name != "test-backup" {
		t.Errorf("expected name 'test-backup', got '%s'", info.Name)
	}
	if info.Namespace != "velero" {
		t.Errorf("expected namespace 'velero', got '%s'", info.Namespace)
	}
	if info.Phase != "Completed" {
		t.Errorf("expected phase 'Completed', got '%s'", info.Phase)
	}
	if info.ScheduleName != "daily-schedule" {
		t.Errorf("expected schedule 'daily-schedule', got '%s'", info.ScheduleName)
	}
	if info.Warnings != 2 {
		t.Errorf("expected 2 warnings, got %d", info.Warnings)
	}
	if info.Errors != 0 {
		t.Errorf("expected 0 errors, got %d", info.Errors)
	}
	if info.TotalItems != 100 {
		t.Errorf("expected 100 total items, got %d", info.TotalItems)
	}
	if info.ItemsBackedUp != 100 {
		t.Errorf("expected 100 items backed up, got %d", info.ItemsBackedUp)
	}
	if info.StartTimestamp == nil {
		t.Fatal("expected start timestamp to be set")
	}
	if info.CompletionTimestamp == nil {
		t.Fatal("expected completion timestamp to be set")
	}
}

func TestExtractBackupInfo_ExtendedFields(t *testing.T) {
	ttl := metav1.Duration{Duration: 720 * time.Hour}

	backup := velerov1api.Backup{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "full-backup",
			Namespace: "velero",
			Labels: map[string]string{
				"velero.io/schedule-name": "weekly",
				"env":                    "prod",
			},
			Annotations: map[string]string{
				"note": "weekly full backup",
			},
		},
		Spec: velerov1api.BackupSpec{
			StorageLocation:    "default",
			TTL:                ttl,
			IncludedNamespaces: []string{"app-ns", "db-ns"},
			ExcludedNamespaces: []string{"kube-system"},
			IncludedResources:  []string{"deployments", "services"},
			ExcludedResources:  []string{"secrets"},
		},
		Status: velerov1api.BackupStatus{
			Phase:                         velerov1api.BackupPhaseFailed,
			FailureReason:                 "storage location unavailable",
			ValidationErrors:              []string{"invalid label selector", "unknown resource"},
			VolumeSnapshotsAttempted:       5,
			VolumeSnapshotsCompleted:       3,
			CSIVolumeSnapshotsAttempted:    2,
			CSIVolumeSnapshotsCompleted:    1,
			FormatVersion:                  "1.1.0",
			BackupItemOperationsAttempted:  3,
			BackupItemOperationsCompleted:  2,
			BackupItemOperationsFailed:     1,
			HookStatus: &velerov1api.HookStatus{
				HooksAttempted: 5,
				HooksFailed:    1,
			},
		},
	}

	info := extractBackupInfo(backup)

	if info.StorageLocation != "default" {
		t.Errorf("expected storage location 'default', got '%s'", info.StorageLocation)
	}
	if info.TTL != "720h0m0s" {
		t.Errorf("expected TTL '720h0m0s', got '%s'", info.TTL)
	}
	if len(info.IncludedNamespaces) != 2 {
		t.Errorf("expected 2 included namespaces, got %d", len(info.IncludedNamespaces))
	}
	if len(info.ExcludedNamespaces) != 1 {
		t.Errorf("expected 1 excluded namespace, got %d", len(info.ExcludedNamespaces))
	}
	if len(info.IncludedResources) != 2 {
		t.Errorf("expected 2 included resources, got %d", len(info.IncludedResources))
	}
	if len(info.ExcludedResources) != 1 {
		t.Errorf("expected 1 excluded resource, got %d", len(info.ExcludedResources))
	}
	if info.Labels["env"] != "prod" {
		t.Errorf("expected label env=prod, got '%s'", info.Labels["env"])
	}
	if info.Annotations["note"] != "weekly full backup" {
		t.Errorf("expected annotation note='weekly full backup', got '%s'", info.Annotations["note"])
	}
	if info.VolumeSnapshotsAttempted != 5 {
		t.Errorf("expected 5 volume snapshots attempted, got %d", info.VolumeSnapshotsAttempted)
	}
	if info.VolumeSnapshotsCompleted != 3 {
		t.Errorf("expected 3 volume snapshots completed, got %d", info.VolumeSnapshotsCompleted)
	}
	if info.CSIVolumeSnapshotsAttempted != 2 {
		t.Errorf("expected 2 CSI snapshots attempted, got %d", info.CSIVolumeSnapshotsAttempted)
	}
	if info.CSIVolumeSnapshotsCompleted != 1 {
		t.Errorf("expected 1 CSI snapshot completed, got %d", info.CSIVolumeSnapshotsCompleted)
	}
	if info.FailureReason != "storage location unavailable" {
		t.Errorf("expected failure reason, got '%s'", info.FailureReason)
	}
	if len(info.ValidationErrors) != 2 {
		t.Errorf("expected 2 validation errors, got %d", len(info.ValidationErrors))
	}
	if info.FormatVersion != "1.1.0" {
		t.Errorf("expected format version '1.1.0', got '%s'", info.FormatVersion)
	}
	if info.HooksAttempted != 5 {
		t.Errorf("expected 5 hooks attempted, got %d", info.HooksAttempted)
	}
	if info.HooksFailed != 1 {
		t.Errorf("expected 1 hook failed, got %d", info.HooksFailed)
	}
	if info.BackupItemOperationsAttempted != 3 {
		t.Errorf("expected 3 backup item operations attempted, got %d", info.BackupItemOperationsAttempted)
	}
	if info.BackupItemOperationsCompleted != 2 {
		t.Errorf("expected 2 backup item operations completed, got %d", info.BackupItemOperationsCompleted)
	}
	if info.BackupItemOperationsFailed != 1 {
		t.Errorf("expected 1 backup item operation failed, got %d", info.BackupItemOperationsFailed)
	}
}

func TestExtractBackupInfo_NilHookStatus(t *testing.T) {
	backup := velerov1api.Backup{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "no-hooks-backup",
			Namespace: "velero",
		},
		Status: velerov1api.BackupStatus{
			Phase:      velerov1api.BackupPhaseCompleted,
			HookStatus: nil,
		},
	}

	info := extractBackupInfo(backup)

	if info.HooksAttempted != 0 {
		t.Errorf("expected 0 hooks attempted, got %d", info.HooksAttempted)
	}
	if info.HooksFailed != 0 {
		t.Errorf("expected 0 hooks failed, got %d", info.HooksFailed)
	}
}

func TestExtractBackupInfo_NoSchedule(t *testing.T) {
	backup := velerov1api.Backup{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "adhoc-backup",
			Namespace: "velero",
		},
		Status: velerov1api.BackupStatus{
			Phase: velerov1api.BackupPhaseFailed,
		},
	}

	info := extractBackupInfo(backup)

	if info.ScheduleName != "" {
		t.Errorf("expected empty schedule name, got '%s'", info.ScheduleName)
	}
	if info.Phase != "Failed" {
		t.Errorf("expected phase 'Failed', got '%s'", info.Phase)
	}
}

func TestExtractBackupInfo_NoProgress(t *testing.T) {
	backup := velerov1api.Backup{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "in-progress-backup",
			Namespace: "velero",
		},
		Status: velerov1api.BackupStatus{
			Phase:    velerov1api.BackupPhaseInProgress,
			Progress: nil,
		},
	}

	info := extractBackupInfo(backup)

	if info.TotalItems != 0 {
		t.Errorf("expected 0 total items, got %d", info.TotalItems)
	}
	if info.ItemsBackedUp != 0 {
		t.Errorf("expected 0 items backed up, got %d", info.ItemsBackedUp)
	}
}

func TestExtractBackupInfo_NilTimestamps(t *testing.T) {
	backup := velerov1api.Backup{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "new-backup",
			Namespace: "velero",
		},
		Status: velerov1api.BackupStatus{
			Phase: velerov1api.BackupPhaseNew,
		},
	}

	info := extractBackupInfo(backup)

	if info.StartTimestamp != nil {
		t.Error("expected nil start timestamp")
	}
	if info.CompletionTimestamp != nil {
		t.Error("expected nil completion timestamp")
	}
	if info.Expiration != nil {
		t.Error("expected nil expiration")
	}
}
