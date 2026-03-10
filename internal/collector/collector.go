package collector

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	velerov1api "github.com/vmware-tanzu/velero/pkg/apis/velero/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// BackupInfo holds extracted metadata from a Velero Backup CR.
type BackupInfo struct {
	Name                string
	Namespace           string
	Phase               string
	ScheduleName        string
	StartTimestamp      *time.Time
	CompletionTimestamp *time.Time
	Expiration          *time.Time
	ItemsBackedUp       int
	TotalItems          int
	Warnings            int
	Errors              int

	// Extended fields for detail view
	StorageLocation             string
	TTL                         string
	IncludedNamespaces          []string
	ExcludedNamespaces          []string
	IncludedResources           []string
	ExcludedResources           []string
	Labels                      map[string]string
	Annotations                 map[string]string
	VolumeSnapshotsAttempted    int
	VolumeSnapshotsCompleted    int
	CSIVolumeSnapshotsAttempted int
	CSIVolumeSnapshotsCompleted int
	FailureReason               string
	ValidationErrors            []string

	// Hook and operation status
	HooksAttempted                int
	HooksFailed                   int
	BackupItemOperationsAttempted int
	BackupItemOperationsCompleted int
	BackupItemOperationsFailed    int
	FormatVersion                 string
}

// VolumeBackupInfo holds extracted metadata from a Velero PodVolumeBackup CR.
type VolumeBackupInfo struct {
	VolumeName          string
	PodName             string
	PodNamespace        string
	NodeName            string
	UploaderType        string
	Phase               string
	StartTimestamp      *time.Time
	CompletionTimestamp *time.Time
	TotalBytes          int64
	BytesDone           int64
	SnapshotID          string
}

// ScheduleInfo holds extracted metadata from a Velero Schedule CR.
type ScheduleInfo struct {
	Name           string
	Namespace      string
	Schedule       string
	Paused         bool
	Phase          string
	LastBackupTime *time.Time
}

// Collector periodically fetches Velero backup and schedule data.
type Collector struct {
	client    client.Client
	namespace string
	interval  time.Duration

	mu        sync.RWMutex
	backups   []BackupInfo
	schedules []ScheduleInfo
}

// New creates a new Collector.
func New(c client.Client, namespace string, interval time.Duration) *Collector {
	return &Collector{
		client:    c,
		namespace: namespace,
		interval:  interval,
	}
}

// Backups returns the most recently collected backup data.
func (c *Collector) Backups() []BackupInfo {
	c.mu.RLock()
	defer c.mu.RUnlock()
	result := make([]BackupInfo, len(c.backups))
	copy(result, c.backups)
	return result
}

// GetBackup returns a single backup by name, or nil if not found.
func (c *Collector) GetBackup(name string) *BackupInfo {
	c.mu.RLock()
	defer c.mu.RUnlock()
	for _, b := range c.backups {
		if b.Name == name {
			result := b
			return &result
		}
	}
	return nil
}

// Schedules returns the most recently collected schedule data.
func (c *Collector) Schedules() []ScheduleInfo {
	c.mu.RLock()
	defer c.mu.RUnlock()
	result := make([]ScheduleInfo, len(c.schedules))
	copy(result, c.schedules)
	return result
}

// SetData sets backup and schedule data directly (for testing).
func (c *Collector) SetData(backups []BackupInfo, schedules []ScheduleInfo) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.backups = backups
	c.schedules = schedules
}

// Run starts periodic collection. It blocks until ctx is cancelled.
func (c *Collector) Run(ctx context.Context) {
	// Collect immediately on start
	if err := c.collect(ctx); err != nil {
		log.Printf("ERROR: initial collection failed: %v", err)
	}

	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := c.collect(ctx); err != nil {
				log.Printf("ERROR: collection failed: %v", err)
			}
		}
	}
}

func (c *Collector) collect(ctx context.Context) error {
	backups, err := c.listBackups(ctx)
	if err != nil {
		return fmt.Errorf("listing backups: %w", err)
	}

	schedules, err := c.listSchedules(ctx)
	if err != nil {
		return fmt.Errorf("listing schedules: %w", err)
	}

	// Build schedule name set for association
	scheduleNames := make(map[string]bool)
	for _, s := range schedules {
		scheduleNames[s.Name] = true
	}

	// Associate backups with schedules via the schedule label
	for i := range backups {
		if backups[i].ScheduleName != "" && !scheduleNames[backups[i].ScheduleName] {
			// Schedule reference exists but schedule CR is gone
			backups[i].ScheduleName = backups[i].ScheduleName + " (deleted)"
		}
	}

	c.mu.Lock()
	c.backups = backups
	c.schedules = schedules
	c.mu.Unlock()

	log.Printf("INFO: collected %d backups, %d schedules", len(backups), len(schedules))
	return nil
}

func (c *Collector) listBackups(ctx context.Context) ([]BackupInfo, error) {
	var backupList velerov1api.BackupList
	opts := []client.ListOption{}
	if c.namespace != "" {
		opts = append(opts, client.InNamespace(c.namespace))
	}
	if err := c.client.List(ctx, &backupList, opts...); err != nil {
		return nil, err
	}

	infos := make([]BackupInfo, 0, len(backupList.Items))
	for _, b := range backupList.Items {
		infos = append(infos, extractBackupInfo(b))
	}
	return infos, nil
}

func extractBackupInfo(b velerov1api.Backup) BackupInfo {
	info := BackupInfo{
		Name:      b.Name,
		Namespace: b.Namespace,
		Phase:     string(b.Status.Phase),
		Warnings:  b.Status.Warnings,
		Errors:    b.Status.Errors,
	}

	// Extract schedule name from label
	if labels := b.GetLabels(); labels != nil {
		if sched, ok := labels["velero.io/schedule-name"]; ok {
			info.ScheduleName = sched
		}
	}

	if b.Status.StartTimestamp != nil {
		t := b.Status.StartTimestamp.Time
		info.StartTimestamp = &t
	}
	if b.Status.CompletionTimestamp != nil {
		t := b.Status.CompletionTimestamp.Time
		info.CompletionTimestamp = &t
	}
	if b.Status.Expiration != nil {
		t := b.Status.Expiration.Time
		info.Expiration = &t
	}
	if b.Status.Progress != nil {
		info.ItemsBackedUp = b.Status.Progress.ItemsBackedUp
		info.TotalItems = b.Status.Progress.TotalItems
	}

	// Extended fields
	info.StorageLocation = b.Spec.StorageLocation
	if b.Spec.TTL.Duration > 0 {
		info.TTL = b.Spec.TTL.Duration.String()
	}
	info.IncludedNamespaces = b.Spec.IncludedNamespaces
	info.ExcludedNamespaces = b.Spec.ExcludedNamespaces
	info.IncludedResources = b.Spec.IncludedResources
	info.ExcludedResources = b.Spec.ExcludedResources
	info.Labels = b.GetLabels()
	info.Annotations = b.GetAnnotations()
	info.VolumeSnapshotsAttempted = b.Status.VolumeSnapshotsAttempted
	info.VolumeSnapshotsCompleted = b.Status.VolumeSnapshotsCompleted
	info.CSIVolumeSnapshotsAttempted = b.Status.CSIVolumeSnapshotsAttempted
	info.CSIVolumeSnapshotsCompleted = b.Status.CSIVolumeSnapshotsCompleted
	info.FailureReason = b.Status.FailureReason
	info.ValidationErrors = b.Status.ValidationErrors

	// Hook and operation status
	if b.Status.HookStatus != nil {
		info.HooksAttempted = b.Status.HookStatus.HooksAttempted
		info.HooksFailed = b.Status.HookStatus.HooksFailed
	}
	info.BackupItemOperationsAttempted = b.Status.BackupItemOperationsAttempted
	info.BackupItemOperationsCompleted = b.Status.BackupItemOperationsCompleted
	info.BackupItemOperationsFailed = b.Status.BackupItemOperationsFailed
	info.FormatVersion = b.Status.FormatVersion

	return info
}

// ListVolumeBackups queries PodVolumeBackup CRs associated with the given backup name.
func ListVolumeBackups(ctx context.Context, c client.Client, backupName, namespace string) []VolumeBackupInfo {
	var pvbList velerov1api.PodVolumeBackupList
	opts := []client.ListOption{
		client.MatchingLabels{"velero.io/backup-name": backupName},
	}
	if namespace != "" {
		opts = append(opts, client.InNamespace(namespace))
	}

	if err := c.List(ctx, &pvbList, opts...); err != nil {
		log.Printf("ERROR: listing PodVolumeBackups for backup %s: %v", backupName, err)
		return nil
	}

	infos := make([]VolumeBackupInfo, 0, len(pvbList.Items))
	for _, pvb := range pvbList.Items {
		info := VolumeBackupInfo{
			VolumeName:   pvb.Spec.Volume,
			PodName:      pvb.Spec.Pod.Name,
			PodNamespace: pvb.Spec.Pod.Namespace,
			NodeName:     pvb.Spec.Node,
			UploaderType: pvb.Spec.UploaderType,
			Phase:        string(pvb.Status.Phase),
			TotalBytes:   pvb.Status.Progress.TotalBytes,
			BytesDone:    pvb.Status.Progress.BytesDone,
			SnapshotID:   pvb.Status.SnapshotID,
		}
		if pvb.Status.StartTimestamp != nil {
			t := pvb.Status.StartTimestamp.Time
			info.StartTimestamp = &t
		}
		if pvb.Status.CompletionTimestamp != nil {
			t := pvb.Status.CompletionTimestamp.Time
			info.CompletionTimestamp = &t
		}
		infos = append(infos, info)
	}
	return infos
}

func (c *Collector) listSchedules(ctx context.Context) ([]ScheduleInfo, error) {
	var scheduleList velerov1api.ScheduleList
	opts := []client.ListOption{}
	if c.namespace != "" {
		opts = append(opts, client.InNamespace(c.namespace))
	}
	if err := c.client.List(ctx, &scheduleList, opts...); err != nil {
		return nil, err
	}

	infos := make([]ScheduleInfo, 0, len(scheduleList.Items))
	for _, s := range scheduleList.Items {
		info := ScheduleInfo{
			Name:      s.Name,
			Namespace: s.Namespace,
			Schedule:  s.Spec.Schedule,
			Paused:    s.Spec.Paused,
			Phase:     string(s.Status.Phase),
		}
		if s.Status.LastBackup != nil {
			t := s.Status.LastBackup.Time
			info.LastBackupTime = &t
		}
		infos = append(infos, info)
	}
	return infos, nil
}
