## Why

There is no way to verify SMTP configuration works without waiting for the cron schedule to fire. When setting up email notifications, operators need immediate feedback that their SMTP host, credentials, and recipient addresses are correctly configured. A manual trigger endpoint eliminates the guesswork and shortens the feedback loop from hours to seconds.

## What Changes

- Add a `POST /api/v1/email/test` API endpoint that generates a backup report and sends it immediately via the configured SMTP sender
- Pass the email `Sender` (and `Collector`) into the `Server` so the endpoint has access to them
- Return a clear success/error JSON response indicating whether the email was sent
- Add a "Send Test Email" button in the web UI dashboard so operators can trigger it without curl

## Capabilities

### New Capabilities
- `email-test`: On-demand email send endpoint for SMTP configuration verification

### Modified Capabilities

## Impact

- **`internal/server/server.go`**: Add `emailSender` field to `Server` struct, new `WithEmailSender` option, new `POST` handler
- **`cmd/velero-backup-reporter/main.go`**: Pass sender into server via the new option (when email is enabled)
- **`web/frontend/src/views/DashboardView.vue`**: Add "Send Test Email" button with loading/success/error feedback
- **API surface**: New `POST /api/v1/email/test` endpoint (additive, non-breaking)
