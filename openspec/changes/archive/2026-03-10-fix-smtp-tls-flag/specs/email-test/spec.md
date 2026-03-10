## MODIFIED Requirements

### Requirement: SMTP TLS configuration
The email sender SHALL respect the `smtp-tls` configuration flag when establishing SMTP connections.

#### Scenario: TLS disabled allows plain SMTP
- **WHEN** `smtp-tls` is set to `false`
- **THEN** the sender connects via plain SMTP without attempting STARTTLS, and the email is delivered successfully

#### Scenario: TLS enabled performs STARTTLS
- **WHEN** `smtp-tls` is set to `true` (default)
- **THEN** the sender performs STARTTLS with certificate validation before sending the email

## ADDED Requirements

### Requirement: Test email opt-in gate
The test email endpoint and UI button SHALL only be available when the `email-test-enabled` flag is explicitly set to `true`.

#### Scenario: Test disabled by default
- **WHEN** the application starts with `--email-enabled` but without `--email-test-enabled`
- **THEN** the `POST /api/v1/email/test` endpoint returns 501 and the dashboard does not show the test button

#### Scenario: Test enabled explicitly
- **WHEN** the application starts with `--email-enabled --email-test-enabled`
- **THEN** the `POST /api/v1/email/test` endpoint is functional and the dashboard shows the "Send Test Email" button

#### Scenario: Test flag without email enabled
- **WHEN** the application starts with `--email-test-enabled` but without `--email-enabled`
- **THEN** the test endpoint returns 501 and no button is shown (email must be enabled first)
