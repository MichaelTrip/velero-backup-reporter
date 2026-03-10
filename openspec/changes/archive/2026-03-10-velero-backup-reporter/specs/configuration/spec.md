## ADDED Requirements

### Requirement: Configuration sources
The system SHALL accept configuration from environment variables, a YAML config file, and CLI flags, with CLI flags taking highest precedence, then environment variables, then config file.

#### Scenario: Precedence order
- **WHEN** the same setting is specified in multiple sources
- **THEN** CLI flags SHALL override environment variables, which SHALL override config file values

#### Scenario: Config file loading
- **WHEN** a config file path is specified via `--config` flag
- **THEN** the system SHALL load configuration from the specified YAML file

### Requirement: Kubernetes configuration
The system SHALL accept configuration for Kubernetes cluster connectivity.

#### Scenario: Kubeconfig path
- **WHEN** `--kubeconfig` flag or `KUBECONFIG` env var is set
- **THEN** the system SHALL use the specified kubeconfig for cluster access

#### Scenario: Namespace configuration
- **WHEN** `--namespace` flag or `NAMESPACE` env var is set
- **THEN** the system SHALL monitor Velero resources in the specified namespace

### Requirement: Web server configuration
The system SHALL accept configuration for the HTTP server.

#### Scenario: Port configuration
- **WHEN** `--port` flag or `PORT` env var is set
- **THEN** the system SHALL serve the web UI on the specified port

### Requirement: Email configuration
The system SHALL accept configuration for email notifications.

#### Scenario: SMTP settings via environment
- **WHEN** SMTP settings are provided via environment variables (`SMTP_HOST`, `SMTP_PORT`, `SMTP_USERNAME`, `SMTP_PASSWORD`, `SMTP_FROM`, `SMTP_TO`, `SMTP_TLS`)
- **THEN** the system SHALL use these values for email configuration

#### Scenario: Email schedule configuration
- **WHEN** `EMAIL_SCHEDULE` env var or `--email-schedule` flag is set with a cron expression
- **THEN** the system SHALL use the specified schedule for sending email reports

### Requirement: Configuration validation
The system SHALL validate configuration at startup.

#### Scenario: Invalid SMTP configuration
- **WHEN** email is enabled but required SMTP fields (host, from, to) are missing
- **THEN** the system SHALL log a warning and disable email notifications

#### Scenario: Invalid port
- **WHEN** an invalid port number is provided
- **THEN** the system SHALL exit with an error message
