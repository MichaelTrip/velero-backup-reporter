## Context

The project is a Go + Vue 3 application with a multi-stage Dockerfile that builds the frontend then the backend. It currently has raw Kubernetes manifests in `deploy/manifests.yaml`. The organization uses a standardized GitHub Actions workflow pattern (seen in `virthorn-scheduler`) for CI/CD with GHCR publishing, semantic versioning, and auto-deployment updates.

## Goals / Non-Goals

**Goals:**
- Mirror the `virthorn-scheduler` CI/CD workflow adapted for a Go + frontend project
- Provide a Helm chart that covers all deployment needs (RBAC, Deployment, Service, configurable values)
- Inject version info into the Go binary at build time

**Non-Goals:**
- Helm chart repository hosting (just the chart in the repo for now)
- Multi-architecture builds (linux/amd64 only, matching existing pattern)
- Ingress or TLS in the Helm chart (users add their own)

## Decisions

### 1. CI/CD workflow structure

Same job structure as `virthorn-scheduler`:
1. `determine-tag` — semantic version from commit messages or manual override
2. `test` — `go vet`, `go test`, plus `npm ci && npm run build` for frontend validation
3. `build-container` — Docker Buildx with GHCR push, GHA cache
4. `create-git-tag` — tag on main
5. `create-github-release` — changelog from commits
6. `update-deployment` — sed the manifest image tag
7. `build-summary` — step summary output

Adaptation for this project: the test job also runs frontend build to catch JS issues early.

### 2. Helm chart structure

Standard Helm chart at `charts/velero-backup-reporter/`:
- `Chart.yaml` — metadata
- `values.yaml` — configurable defaults (image, resources, probes, email config, namespace)
- `templates/` — ServiceAccount, ClusterRole, ClusterRoleBinding, Deployment, Service
- Values for email config passed as env vars to the container

### 3. Version injection

Add `-ldflags "-X main.version=${VERSION}"` to the Go build in the Dockerfile. Accept `VERSION` as a Docker build arg. The CI workflow passes the determined tag as the build arg.

## Risks / Trade-offs

- **Manifest auto-update creates commits on main** — same pattern as virthorn-scheduler, acceptable for this workflow.
- **Helm chart and raw manifests coexist** — raw manifests stay for simple kubectl users, Helm chart for production. Both are valid deployment paths.
