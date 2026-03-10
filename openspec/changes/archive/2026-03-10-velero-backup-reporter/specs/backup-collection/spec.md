## ADDED Requirements

### Requirement: Connect to Kubernetes cluster
The system SHALL connect to a Kubernetes cluster using in-cluster service account credentials when running as a pod, or a kubeconfig file when running externally.

#### Scenario: In-cluster authentication
- **WHEN** the application starts inside a Kubernetes pod without a kubeconfig path configured
- **THEN** the system uses the mounted service account token to authenticate with the Kubernetes API

#### Scenario: Out-of-cluster authentication
- **WHEN** the application starts with a kubeconfig path configured
- **THEN** the system uses the specified kubeconfig to authenticate with the Kubernetes API

#### Scenario: Authentication failure
- **WHEN** the system cannot authenticate with the Kubernetes API
- **THEN** the system SHALL log an error and exit with a non-zero status code

### Requirement: List Velero Backups
The system SHALL list all Velero Backup custom resources from the configured namespace(s).

#### Scenario: List backups from default Velero namespace
- **WHEN** no specific namespace is configured
- **THEN** the system SHALL list Backup CRs from the `velero` namespace

#### Scenario: List backups from configured namespace
- **WHEN** a namespace is configured
- **THEN** the system SHALL list Backup CRs only from that namespace

#### Scenario: List backups from all namespaces
- **WHEN** the namespace configuration is set to all namespaces
- **THEN** the system SHALL list Backup CRs from all namespaces

### Requirement: List Velero Schedules
The system SHALL list all Velero Schedule custom resources to correlate backups with their schedules.

#### Scenario: Schedules retrieved
- **WHEN** the system collects backup data
- **THEN** it SHALL also retrieve all Schedule CRs from the same namespace(s) and associate backups with their owning schedules

### Requirement: Periodic data collection
The system SHALL periodically collect backup data at a configurable interval.

#### Scenario: Default collection interval
- **WHEN** no collection interval is configured
- **THEN** the system SHALL collect data every 5 minutes

#### Scenario: Custom collection interval
- **WHEN** a collection interval is configured (e.g., 10m)
- **THEN** the system SHALL collect data at the specified interval

### Requirement: Extract backup metadata
The system SHALL extract relevant metadata from each Backup CR.

#### Scenario: Backup metadata extraction
- **WHEN** a Backup CR is read
- **THEN** the system SHALL extract: name, namespace, phase (status), start timestamp, completion timestamp, expiration, items backed up, total items, warnings count, errors count, and associated schedule name
