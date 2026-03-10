## Context

The email subsystem (`internal/email/`) already has a fully functional `Sender` that generates HTML reports and delivers them via SMTP. A `Scheduler` wraps the sender with cron-based delivery. Currently the `Sender` is only accessible within the scheduler — the HTTP server has no reference to it. The `Server` struct holds a `collector` and an optional `kubeClient`, but nothing email-related.

## Goals / Non-Goals

**Goals:**
- Allow operators to verify SMTP configuration by sending a real email on demand
- Surface the trigger in the web UI so no CLI/curl knowledge is required
- Return actionable error messages when SMTP is misconfigured

**Non-Goals:**
- Customising the test email content (it sends the same report the scheduler would)
- Adding email configuration via the UI (config stays in flags/env/YAML)
- Rate limiting or authentication on the test endpoint

## Decisions

### 1. Pass `Sender` into `Server` via functional option

Add a `WithEmailSender(*email.Sender)` option alongside the existing `WithKubeClient`. This keeps the server's constructor signature stable and makes the sender optional — when email is disabled, it's simply not passed in. The endpoint returns 404 or 501 when no sender is configured.

**Alternative considered**: Passing the entire `Scheduler` to the server. Rejected because the server only needs `Sender.Send()` and the collector it already has. No reason to couple to the scheduler.

### 2. Endpoint design: `POST /api/v1/email/test`

- `POST` because it triggers a side effect (sending an email)
- Returns `{"message": "..."}` on success (200) or `{"error": "..."}` on failure (500/501)
- Generates a fresh `report.Generate()` using the server's existing collector, same as the scheduler does

### 3. UI button on the Dashboard

A "Send Test Email" button in the dashboard header area, next to the title. Uses PrimeVue `Button` with loading state. Shows a PrimeVue `Toast` on success/error for non-intrusive feedback.

**Alternative considered**: Separate settings/email page. Rejected as over-engineering for a single button — the dashboard is where operators check system health.

## Risks / Trade-offs

- **No rate limiting** → An operator could spam the endpoint. Acceptable for an internal tool; the SMTP server itself will rate-limit if needed.
- **Endpoint available without auth** → Same as all other endpoints in this app. Auth is out of scope.
- **Button hidden when email disabled** → The endpoint returns 501 and the UI conditionally shows the button. Requires a way for the frontend to know if email is enabled — solved by adding an `emailEnabled` field to the dashboard API response.
