package logs

import (
	"bytes"
	"compress/gzip"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchAndDecompress(t *testing.T) {
	// Create gzip-compressed test content
	var buf bytes.Buffer
	gzWriter := gzip.NewWriter(&buf)
	logContent := "time=\"2024-01-01T00:00:00Z\" level=info msg=\"backup started\"\n"
	gzWriter.Write([]byte(logContent))
	gzWriter.Close()

	// Serve it from a test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/gzip")
		w.Write(buf.Bytes())
	}))
	defer server.Close()

	result, err := fetchAndDecompress(context.Background(), server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != logContent {
		t.Errorf("expected %q, got %q", logContent, result)
	}
}

func TestFetchAndDecompress_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	_, err := fetchAndDecompress(context.Background(), server.URL)
	if err == nil {
		t.Fatal("expected error for 404 response")
	}
}

func TestFetchAndDecompress_ContextCancelled(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Never respond
		<-r.Context().Done()
	}))
	defer server.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, err := fetchAndDecompress(ctx, server.URL)
	if err == nil {
		t.Fatal("expected error for cancelled context")
	}
}
