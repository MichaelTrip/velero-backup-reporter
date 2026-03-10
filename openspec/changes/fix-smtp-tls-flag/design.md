## Context

Two issues:

1. Go's `smtp.SendMail()` always attempts STARTTLS if the server advertises it. There is no way to opt out. The `SMTPConfig.TLS` bool is loaded but never read by the sender.

2. The "Send Test Email" button currently appears whenever `--email-enabled` is set. There's no way to hide it in production while still having scheduled email delivery.

## Goals / Non-Goals

**Goals:**
- `--smtp-tls=false` skips STARTTLS, allowing plain SMTP connections
- `--smtp-tls=true` (default) performs STARTTLS with proper cert validation
- `--email-test-enabled` flag gates the test endpoint and UI button (default `false`)

**Non-Goals:**
- Adding `InsecureSkipVerify` option (separate concern)
- Supporting implicit TLS on port 465 (SMTPS)
- Auth/RBAC on the test endpoint

## Decisions

### 1. Replace `smtp.SendMail()` with manual SMTP client

Use `smtp.Dial()` and manually step through the SMTP conversation:

1. `smtp.Dial(addr)` — connect
2. `client.Hello()` — EHLO
3. If `cfg.TLS` is true: `client.StartTLS(&tls.Config{ServerName: host})`
4. If auth configured: `client.Auth(auth)`
5. `client.Mail(from)` / `client.Rcpt(to)` / `client.Data()` — send
6. `client.Quit()`

### 2. Gate test email behind `--email-test-enabled`

Add `TestEnabled bool` to `EmailConfig`. In `main.go`, only pass `WithEmailSender` to the server when **both** `--email-enabled` and `--email-test-enabled` are true. This means:

- `--email-enabled` alone: scheduled emails work, no test button, no test endpoint
- `--email-enabled --email-test-enabled`: scheduled emails + test button + test endpoint
- Neither: no email functionality at all

The dashboard API's `emailEnabled` field reflects the test gate, not the general email flag, since its only consumer is the test button.

## Risks / Trade-offs

- **Plain SMTP is unencrypted** — credentials and email content in cleartext when TLS disabled. Acceptable for dev/testing and air-gapped environments.
- **Extra flag complexity** — one more flag to document. But keeps the test action opt-in, which is the safer default for production.
