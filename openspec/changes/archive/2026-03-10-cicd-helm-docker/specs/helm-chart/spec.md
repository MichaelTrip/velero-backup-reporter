## ADDED Requirements

### Requirement: Helm chart for Kubernetes deployment
The project SHALL include a Helm chart at `charts/velero-backup-reporter/` that deploys the application to Kubernetes with all required resources.

#### Scenario: Helm install creates all resources
- **WHEN** a user runs `helm install` with the chart
- **THEN** a ServiceAccount, ClusterRole, ClusterRoleBinding, Deployment, and Service are created in the target namespace

### Requirement: Configurable values
The Helm chart SHALL expose configurable values for image repository/tag, resource limits, probe settings, replica count, namespace, and email configuration.

#### Scenario: Custom image tag
- **WHEN** a user sets `image.tag` in values
- **THEN** the Deployment uses that image tag

#### Scenario: Email configuration via values
- **WHEN** a user sets email-related values (smtp host, port, credentials, schedule)
- **THEN** the corresponding environment variables are set on the container

### Requirement: RBAC resources
The Helm chart SHALL create ClusterRole and ClusterRoleBinding with the minimum permissions needed for Velero backup monitoring (get/list/watch on backups, schedules, podvolumebackups; create/get/delete on downloadrequests).

#### Scenario: RBAC permissions match application needs
- **WHEN** the chart is installed
- **THEN** the ClusterRole grants exactly the permissions defined in the current `deploy/manifests.yaml`
