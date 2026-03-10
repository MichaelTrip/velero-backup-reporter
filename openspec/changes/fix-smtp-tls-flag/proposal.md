## Why

Two issues with the email system:

1. **TLS flag ignored**: The `--smtp-tls` flag is accepted but never used. The email sender calls `smtp.SendMail()` which unconditionally attempts STARTTLS and fails with `x509` errors when connecting to SMTP servers without valid TLS.

2. **Test button always visible**: The "Send Test Email" button shows on the dashboard whenever `--email-enabled` is set. This exposes an operational action to any user with dashboard access. The button should only appear when explicitly opted in via a dedicated flag, keeping it hidden by default even when email is enabled.

## What Changes

- Replace `smtp.SendMail()` with manual SMTP dial + conditional STARTTLS using `net/smtp` directly
- When `TLS` is `false`, skip the STARTTLS step entirely so plain SMTP connections work
- When `TLS` is `true` (default), perform STARTTLS as before with proper certificate validation
- Add a new `--email-test-enabled` CLI flag / `EMAIL_TEST_ENABLED` env var (default `false`)
- Only expose the `POST /api/v1/email/test` endpoint and the dashboard button when this flag is `true`

## Capabilities

### New Capabilities
- `email-test-gate`: The test email endpoint and UI button SHALL be gated behind an explicit opt-in flag

### Modified Capabilities
- `email-test`: The email sender SHALL respect the TLS configuration flag

## Impact

- **`internal/email/email.go`**: Replace `smtp.SendMail()` with manual SMTP client
- **`internal/config/config.go`**: Add `TestEnabled` field to `EmailConfig`
- **`cmd/velero-backup-reporter/main.go`**: Add `--email-test-enabled` flag, pass to server conditionally
- **`internal/server/server.go`**: Only register test endpoint and set `emailEnabled` when test is opted in
- **`web/frontend/src/views/DashboardView.vue`**: No change needed (already conditional on `emailEnabled`)
