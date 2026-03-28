---
applyTo:
  - "**/*_test.go"
  - Makefile
  - Dockerfile
  - ".github/workflows/**"
description: Testing, CI/CD, and build automation for Velero Backup Reporter
---

# Testing & CI/CD Guide

This guide covers testing, building, and deployment automation for Velero Backup Reporter.

## Testing Strategy

### Unit Tests (Backend)

**Test file structure:**
```go
// internal/collector/collector_test.go
package collector

import (
    "context"
    "testing"
)

func TestCollectorRun(t *testing.T) {
    // Arrange: Set up test data and mocks
    mockClient := &MockKubeClient{
        backups: []corev1.Backup{{ /* ... */ }},
    }

    // Act: Call the function being tested
    coll := New(mockClient, "velero", 1*time.Second)
    err := coll.Run(context.Background())

    // Assert: Verify results
    if err != nil {
        t.Fatalf("expected no error, got %v", err)
    }
}
```

**Patterns:**
- Use table-driven tests for multiple scenarios:
  ```go
  tests := []struct {
      name    string
      input   string
      want    string
      wantErr bool
  }{
      {"valid", "backup-1", "BACKUP-1", false},
      {"empty", "", "", true},
  }

  for _, tt := range tests {
      t.Run(tt.name, func(t *testing.T) {
          got, err := transform(tt.input)
          if (err != nil) != tt.wantErr {
              t.Errorf("got error %v, want %v", err, tt.wantErr)
          }
          if got != tt.want {
              t.Errorf("got %q, want %q", got, tt.want)
          }
      })
  }
  ```

- Mock external dependencies (K8s client, SMTP server, filesystem)
- Avoid network calls in unit tests (use mocks or fixtures)

### Integration Tests

**When to use:** Testing multiple components together (collector + config + server)

**Pattern:**
```go
func TestEndToEnd(t *testing.T) {
    // Start a real HTTP server or use httptest
    srv := httptest.NewServer(handler)
    defer srv.Close()

    // Make requests
    res, err := http.Get(srv.URL + "/api/backups")
    if err != nil {
        t.Fatal(err)
    }

    // Verify response
    if res.StatusCode != http.StatusOK {
        t.Errorf("expected 200, got %d", res.StatusCode)
    }
}
```

### Frontend Tests (Future)

Consider adding when UI complexity grows:
- **Unit tests:** Vue Test Utils + Vitest for components/composables
- **E2E tests:** Playwright against running backend at `localhost:8080`

## Build & Release

### Local Build

```bash
# Full build (frontend + backend)
make build

# Backend only
make backend

# Frontend only
make frontend

# Clean artifacts
make clean
```

### Docker Build

```bash
# Build image
docker build -t velero-backup-reporter:latest .

# Run container
docker run \
  -p 8080:8080 \
  -v ~/.kube/config:/kubeconfig:ro \
  -e KUBECONFIG=/kubeconfig \
  velero-backup-reporter:latest

# Tag for registry
docker tag velero-backup-reporter:latest myregistry/velero-backup-reporter:v1.0.0
docker push myregistry/velero-backup-reporter:v1.0.0
```

**Dockerfile features:**
- Multi-stage build (reduces final image size)
- CGO_ENABLED=0 for portable binary (no libc dependency)
- Minimal base image (alpine or scratch)
- Non-root user (security best practice)

### Kubernetes Deployment

**Using Helm chart:**
```bash
helm install velero-reporter charts/velero-backup-reporter/ \
  --namespace velero \
  --set smtp.host=smtp.example.com \
  --set smtp.from=alerts@example.com
```

**Using manifests:**
```bash
kubectl apply -f deploy/manifests.yaml
```

**Verification:**
```bash
kubectl get deployment -n velero
kubectl get pod -n velero
kubectl logs -f -n velero deployment/velero-backup-reporter
```

## Running Tests

### All Tests

```bash
# Run all backend tests
go test ./internal/...

# Verbose output
go test -v ./internal/...

# With coverage
go test -cover ./internal/...
```

### Specific Package

```bash
go test ./internal/collector
go test ./internal/config
go test ./internal/email
```

### Specific Test

```bash
go test -run TestCollectorRun ./internal/collector
```

### Coverage Report

```bash
# Generate coverage profile
go test -coverprofile=coverage.out ./internal/...

# View in browser
go tool cover -html=coverage.out
```

### Continuous Testing (Watch Mode)

Install a file watcher (e.g., `entr`):
```bash
find internal/ -name '*.go' | entr -r go test ./internal/...
```

## Linting & Code Quality

### Format Code

```bash
# Format all Go files
gofmt -w .

# Or use goimports (also fixes imports)
goimports -w .
```

### Lint for Issues

```bash
# Run go vet
go vet ./...

# Run golangci-lint (if installed)
golangci-lint run ./internal/...
```

### Dependencies

```bash
# Tidy unused imports
go mod tidy

# Verify checksums
go mod verify

# List all dependencies
go mod graph
```

## Development Workflow

### Quick Dev Loop

**Terminal 1 (Frontend with hot reload):**
```bash
cd web/frontend
npm run dev
# Opens http://localhost:5173
```

**Terminal 2 (Backend):**
```bash
go run ./cmd/velero-backup-reporter/ --kubeconfig ~/.kube/config --namespace velero
# Serves on http://localhost:8080
```

**Terminal 3 (Run tests on save):**
```bash
find internal/ -name '*.go' | entr -r go test ./internal/... -v
```

### Pre-commit Checks

Before committing, run:
```bash
go fmt ./...
go vet ./...
go test ./internal/...
```

Or use a `pre-commit` hook:
```bash
#!/bin/bash
set -e
go fmt ./...
go vet ./...
go test ./internal/...
echo "✓ All checks passed"
```

## CI/CD Pipeline (Recommended Setup)

### GitHub Actions Workflow Example

Create `.github/workflows/test.yml`:

```yaml
name: Tests

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.25.7'

      - name: Set up Node
        uses: actions/setup-node@v3
        with:
          node-version: '18'

      - name: Lint
        run: |
          go vet ./...
          go fmt ./... && git diff --exit-code

      - name: Test
        run: go test -v -cover ./internal/...

      - name: Build
        run: make build

      - name: Build Docker image
        run: docker build -t velero-backup-reporter:test .
```

### Build Pipeline

```yaml
name: Release

on:
  push:
    tags:
      - 'v*'

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.25.7'

      - name: Build
        run: make build

      - name: Create GitHub Release
        uses: softprops/action-gh-release@v1
        with:
          files: velero-backup-reporter
```

## Test Coverage Goals

- **Target:** ≥ 70% total coverage
- **Critical paths:** Collector, config validation, email sending (≥ 80%)
- **Lower priority:** Server handlers, JSON serialization (≥ 50%)

**Check coverage:**
```bash
go test -coverprofile=coverage.out ./internal/...
go tool cover -func=coverage.out | grep total
```

## Mocking & Test Utilities

### Mock K8s Client

```go
type MockKubeClient struct {
    backups   []veleroapi.Backup
    schedules []veleroapi.Schedule
    err       error
}

func (m *MockKubeClient) ListBackups(ctx context.Context, ns string) ([]veleroapi.Backup, error) {
    if m.err != nil {
        return nil, m.err
    }
    return m.backups, nil
}

func (m *MockKubeClient) WatchBackups(ctx context.Context, ns string, handler func(*veleroapi.Backup)) error {
    for _, b := range m.backups {
        handler(&b)
    }
    return m.err
}
```

### Test Server (httptest)

```go
import "net/http/httptest"

func TestAPIHandler(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        fmt.Fprint(w, `{"status":"ok"}`)
    }))
    defer server.Close()

    res, _ := http.Get(server.URL)
    // Assert res
}
```

### Test Data Fixtures

```go
// testdata/backups.json
[
  {
    "name": "backup-1",
    "status": "Completed",
    "timestamp": "2026-03-28T10:00:00Z"
  }
]

// In test:
data, _ := os.ReadFile("testdata/backups.json")
var backups []Backup
json.Unmarshal(data, &backups)
```

## Troubleshooting Tests

| Issue | Solution |
|-------|----------|
| `kubeconfig not found` in tests | Use mocks or set `KUBECONFIG` env var to test kubeconfig |
| Tests timeout | Increase timeout: `go test -timeout 30s ./internal/...` |
| Race conditions | Run with race detector: `go test -race ./internal/...` |
| Flaky tests (intermittent failures) | Check for time dependencies; mock time if needed |
| Import cycles | Use dependency injection to break cycles |

## Common Test Commands

```bash
# All tests with output
go test -v ./internal/...

# Stop on first failure
go test -failfast ./internal/...

# Run with race detector
go test -race ./internal/...

# Show code coverage percentage
go test -cover ./internal/...

# Generate HTML coverage report
go test -coverprofile=coverage.out ./internal/...
go tool cover -html=coverage.out -o coverage.html

# Benchmark
go test -bench=. ./internal/...

# Run specific test by name
go test -run ^TestCollector ./internal/collector
```

## Release Checklist

Before tagging a release (`v1.0.0`):

- [ ] All tests pass: `go test ./internal/...`
- [ ] Code formatted: `gofmt -w .`
- [ ] No linting issues: `go vet ./...`
- [ ] Frontend built: `make frontend`
- [ ] Backend built successfully: `make build`
- [ ] Docker image builds: `docker build .`
- [ ] Helm chart valid: `helm lint charts/velero-backup-reporter/`
- [ ] Update version in `cmd/main.go` (if applicable)
- [ ] Update [README.md](../README.md) changelog section
- [ ] Tag commit: `git tag v1.0.0 && git push origin v1.0.0`
- [ ] Push Docker image to registry
- [ ] Create GitHub release with binary artifacts

## Useful Commands Summary

```bash
# Test
go test -v -cover ./internal/...

# Lint
go fmt ./... && go vet ./...

# Build
make build                # Full build
make frontend             # Frontend only
make backend              # Backend only
docker build -t velero-backup-reporter .

# Run
./velero-backup-reporter --kubeconfig ~/.kube/config --namespace velero
kubectl apply -f deploy/manifests.yaml

# Deploy
helm install velero-reporter charts/velero-backup-reporter/ -n velero
```
