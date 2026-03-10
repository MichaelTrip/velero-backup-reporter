## ADDED Requirements

### Requirement: Automated semantic versioning
The CI pipeline SHALL determine version tags automatically from commit messages using semantic versioning conventions (BREAKING CHANGE = major, feat: = minor, default = patch).

#### Scenario: Patch bump on regular commit
- **WHEN** a commit is pushed to main with no special prefix
- **THEN** the pipeline increments the patch version (e.g., v1.0.0 → v1.0.1)

#### Scenario: Minor bump on feature commit
- **WHEN** a commit message starts with "feat:" or "feature:"
- **THEN** the pipeline increments the minor version (e.g., v1.0.1 → v1.1.0)

#### Scenario: Manual version override
- **WHEN** the workflow is triggered via workflow_dispatch with a version-override input
- **THEN** the pipeline uses the exact version provided

### Requirement: Go test and vet on every push
The pipeline SHALL run `go vet ./...` and `go test ./...` on every push and pull request.

#### Scenario: Tests run on PR
- **WHEN** a pull request is opened against main
- **THEN** the pipeline runs Go vet and tests, and the build fails if any test fails

### Requirement: Container image build and publish
The pipeline SHALL build a Docker image and push it to GHCR on every push (not on PRs).

#### Scenario: Image pushed to GHCR on main
- **WHEN** a commit is pushed to main
- **THEN** the pipeline builds and pushes the image to `ghcr.io/michaeltrip/velero-backup-reporter` with the semantic version tag and `latest`

#### Scenario: Feature branch image tag
- **WHEN** a commit is pushed to a feature branch
- **THEN** the image is tagged with `{branch-name}-{short-sha}`

### Requirement: GitHub Release creation
The pipeline SHALL create a GitHub Release with a changelog when pushing to main.

#### Scenario: Release created on main push
- **WHEN** a commit is pushed to main and the image is built
- **THEN** a GitHub Release is created with the version tag and a changelog generated from commit messages

### Requirement: Deployment manifest auto-update
The pipeline SHALL update the deployment manifest image tag on main branch pushes.

#### Scenario: Manifest updated after build
- **WHEN** a container image is built on main
- **THEN** the pipeline updates `deploy/manifests.yaml` with the new image reference and commits the change
