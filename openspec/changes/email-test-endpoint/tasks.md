## 1. Server Wiring

- [x] 1.1 Add `WithEmailSender(*email.Sender)` option and `emailSender` field to the `Server` struct in `internal/server/server.go`
- [x] 1.2 Pass `sender` into `server.New()` via `WithEmailSender` in `cmd/velero-backup-reporter/main.go` when email is enabled

## 2. API Endpoint

- [x] 2.1 Add `POST /api/v1/email/test` route and handler in `server.go` that generates a report and calls `sender.Send()`; return 200 on success, 501 if no sender, 500 on SMTP error
- [x] 2.2 Add `emailEnabled` boolean field to `dashboardResponse` and set it based on whether `emailSender` is non-nil

## 3. Frontend

- [x] 3.1 Add "Send Test Email" button to `DashboardView.vue` (visible only when `emailEnabled` is true), with loading state and success/error toast feedback

## 4. Build and Verify

- [x] 4.1 Build the frontend and Go binary, run tests
