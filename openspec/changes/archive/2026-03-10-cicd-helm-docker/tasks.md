## 1. Dockerfile Improvements

- [x] 1.1 Add `VERSION` build arg and inject via `-ldflags` into Go build; add `version` variable to `main.go`

## 2. GitHub Actions CI/CD Pipeline

- [x] 2.1 Create `.github/workflows/build-release.yaml` with the full pipeline: determine-tag, test (Go vet + test + frontend build check), build-container (GHCR push), create-git-tag, create-github-release, update-deployment, build-summary

## 3. Helm Chart

- [x] 3.1 Create `charts/velero-backup-reporter/Chart.yaml` and `values.yaml` with configurable defaults
- [x] 3.2 Create Helm templates: `_helpers.tpl`, `serviceaccount.yaml`, `clusterrole.yaml`, `clusterrolebinding.yaml`, `deployment.yaml`, `service.yaml`

## 4. Verify

- [x] 4.1 Validate the Helm chart templates render correctly (helm template dry-run), verify Go builds, run tests
