## Why

The project has no CI/CD pipeline, no Helm chart, and the Dockerfile works but doesn't follow the same patterns as other projects in the organization. To make this production-ready and consistently deployable, it needs a GitHub Actions workflow for automated builds/releases, a Helm chart for Kubernetes deployment, and a Dockerfile aligned with the organization's conventions.

## What Changes

- **GitHub Actions CI/CD pipeline** modeled after the `virthorn-scheduler` workflow:
  - Semantic versioning from commit history (major/minor/patch bumps)
  - Go tests + vet on every push/PR
  - Frontend build integrated into Docker multi-stage
  - Build and push container image to GHCR (`ghcr.io/michaeltrip/velero-backup-reporter`)
  - Git tagging and GitHub Releases on main branch
  - Auto-update deployment manifest with new image tag
  - Build summary in workflow output
  - Manual version override via `workflow_dispatch`
- **Helm chart** for Kubernetes deployment:
  - Replaces raw manifests in `deploy/manifests.yaml`
  - Configurable values for image, resources, probes, email settings, RBAC
  - ServiceAccount, ClusterRole, ClusterRoleBinding, Deployment, Service
- **Dockerfile improvements**:
  - Add build args for version injection (`-ldflags`)
  - Align with organizational patterns

## Capabilities

### New Capabilities
- `ci-cd-pipeline`: Automated build, test, release pipeline via GitHub Actions
- `helm-chart`: Helm chart for Kubernetes deployment

### Modified Capabilities

## Impact

- **`.github/workflows/build-release.yaml`**: New CI/CD workflow
- **`charts/velero-backup-reporter/`**: New Helm chart directory
- **`Dockerfile`**: Updated with version build args
- **`deploy/manifests.yaml`**: Updated image reference for CI auto-update
