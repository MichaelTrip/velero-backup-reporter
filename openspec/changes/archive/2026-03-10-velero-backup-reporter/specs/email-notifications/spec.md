## ADDED Requirements

### Requirement: Send backup reports via email
The system SHALL send backup reports as HTML emails via SMTP when email notifications are enabled.

#### Scenario: Email delivery enabled
- **WHEN** SMTP configuration is provided and email is enabled
- **THEN** the system SHALL send backup report emails to the configured recipients

#### Scenario: Email delivery disabled
- **WHEN** no SMTP configuration is provided or email is explicitly disabled
- **THEN** the system SHALL NOT attempt to send emails and SHALL operate in web-only mode

### Requirement: Configure SMTP settings
The system SHALL support SMTP configuration for email delivery.

#### Scenario: SMTP configuration
- **WHEN** configuring email notifications
- **THEN** the system SHALL accept: SMTP host, port, username, password, from address, TLS enabled flag, and recipient list

### Requirement: Email schedule
The system SHALL send reports on a configurable schedule.

#### Scenario: Default email schedule
- **WHEN** email is enabled but no schedule is configured
- **THEN** the system SHALL send a daily report at 08:00 in the configured timezone

#### Scenario: Custom email schedule
- **WHEN** a cron expression is provided for the email schedule
- **THEN** the system SHALL send reports according to the cron schedule

### Requirement: Email content
The system SHALL send well-formatted HTML email reports.

#### Scenario: Email report content
- **WHEN** an email report is sent
- **THEN** the email SHALL contain the same summary and detail information as the web UI dashboard, formatted as an HTML email

### Requirement: Email failure handling
The system SHALL handle email delivery failures gracefully.

#### Scenario: SMTP connection failure
- **WHEN** the system cannot connect to the SMTP server
- **THEN** the system SHALL log the error and continue operating the web UI normally

#### Scenario: Email send failure
- **WHEN** an email fails to send
- **THEN** the system SHALL log the error with details and retry on the next scheduled interval
