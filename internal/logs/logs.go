package logs

import (
	"compress/gzip"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"time"

	velerov1api "github.com/vmware-tanzu/velero/pkg/apis/velero/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	pollInterval = 500 * time.Millisecond
	pollTimeout  = 30 * time.Second
	maxLogSize   = 10 * 1024 * 1024 // 10MB
)

// FetchBackupLogs retrieves backup logs via the Velero DownloadRequest mechanism.
func FetchBackupLogs(ctx context.Context, kubeClient client.Client, backupName, namespace string) (string, error) {
	// Create a DownloadRequest CR
	dr := &velerov1api.DownloadRequest{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "backup-reporter-" + backupName + "-",
			Namespace:    namespace,
		},
		Spec: velerov1api.DownloadRequestSpec{
			Target: velerov1api.DownloadTarget{
				Kind: velerov1api.DownloadTargetKindBackupLog,
				Name: backupName,
			},
		},
	}

	if err := kubeClient.Create(ctx, dr); err != nil {
		return "", fmt.Errorf("creating download request: %w", err)
	}

	// Always clean up the DownloadRequest
	defer func() {
		cleanupCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = kubeClient.Delete(cleanupCtx, dr)
	}()

	// Poll until processed
	downloadURL, err := waitForDownloadURL(ctx, kubeClient, dr)
	if err != nil {
		return "", err
	}

	// Fetch and decompress the logs
	logContent, err := fetchAndDecompress(ctx, downloadURL)
	if err != nil {
		return "", fmt.Errorf("fetching logs: %w", err)
	}

	return logContent, nil
}

func waitForDownloadURL(ctx context.Context, kubeClient client.Client, dr *velerov1api.DownloadRequest) (string, error) {
	deadline := time.After(pollTimeout)
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	key := types.NamespacedName{Namespace: dr.Namespace, Name: dr.Name}

	for {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-deadline:
			return "", fmt.Errorf("timed out waiting for download request to be processed")
		case <-ticker.C:
			updated := &velerov1api.DownloadRequest{}
			if err := kubeClient.Get(ctx, key, updated); err != nil {
				return "", fmt.Errorf("getting download request status: %w", err)
			}
			if updated.Status.Phase == velerov1api.DownloadRequestPhaseProcessed && updated.Status.DownloadURL != "" {
				return updated.Status.DownloadURL, nil
			}
		}
	}
}

func fetchAndDecompress(ctx context.Context, url string) (string, error) {
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
		},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	gzReader, err := gzip.NewReader(resp.Body)
	if err != nil {
		return "", fmt.Errorf("creating gzip reader: %w", err)
	}
	defer gzReader.Close()

	limited := io.LimitReader(gzReader, maxLogSize)
	data, err := io.ReadAll(limited)
	if err != nil {
		return "", fmt.Errorf("reading logs: %w", err)
	}

	return string(data), nil
}
