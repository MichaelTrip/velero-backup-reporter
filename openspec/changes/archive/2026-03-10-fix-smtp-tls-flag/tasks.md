## 1. Fix SMTP TLS Handling

- [x] 1.1 In `internal/email/email.go`, replace `smtp.SendMail()` with manual SMTP client that conditionally calls `StartTLS()` based on `s.cfg.TLS`

## 2. Gate Test Email Behind Flag

- [x] 2.1 Add `TestEnabled bool` field to `EmailConfig` in `internal/config/config.go` and load from `email-test-enabled`
- [x] 2.2 Add `--email-test-enabled` CLI flag (default `false`) in `cmd/velero-backup-reporter/main.go`
- [x] 2.3 Only pass `WithEmailSender` to the server when both `email-enabled` and `email-test-enabled` are true

## 3. Build and Verify

- [x] 3.1 Build the Go binary and run tests
