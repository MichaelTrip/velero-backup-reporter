## ADDED Requirements

### Requirement: Test email endpoint
The server SHALL expose a `POST /api/v1/email/test` endpoint that generates a backup report and sends it via the configured SMTP sender immediately.

#### Scenario: Successful test email
- **WHEN** a POST request is made to `/api/v1/email/test` and email is configured
- **THEN** the server generates a backup report, sends it via SMTP, and returns HTTP 200 with `{"message": "Test email sent successfully"}`

#### Scenario: Email not configured
- **WHEN** a POST request is made to `/api/v1/email/test` and email is not enabled
- **THEN** the server returns HTTP 501 with `{"error": "email notifications are not enabled"}`

#### Scenario: SMTP delivery failure
- **WHEN** a POST request is made to `/api/v1/email/test` and the SMTP send fails
- **THEN** the server returns HTTP 500 with `{"error": "..."}` containing the failure reason

### Requirement: Email status in dashboard API
The dashboard API response SHALL include an `emailEnabled` boolean field indicating whether email notifications are configured.

#### Scenario: Email enabled
- **WHEN** the dashboard API is called and email is configured
- **THEN** the response includes `"emailEnabled": true`

#### Scenario: Email disabled
- **WHEN** the dashboard API is called and email is not configured
- **THEN** the response includes `"emailEnabled": false`

### Requirement: Test email UI button
The web UI dashboard SHALL display a "Send Test Email" button when email is enabled, allowing operators to trigger an immediate test email.

#### Scenario: Button visible when email enabled
- **WHEN** the dashboard loads and the API returns `emailEnabled: true`
- **THEN** a "Send Test Email" button is displayed

#### Scenario: Button hidden when email disabled
- **WHEN** the dashboard loads and the API returns `emailEnabled: false`
- **THEN** no test email button is displayed

#### Scenario: Button triggers email and shows feedback
- **WHEN** the operator clicks "Send Test Email"
- **THEN** the button shows a loading state, a POST is sent to `/api/v1/email/test`, and a success or error toast message is displayed
