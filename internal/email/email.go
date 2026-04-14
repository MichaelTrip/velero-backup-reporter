package email

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"strings"
	"time"

	"github.com/michael/velero-backup-reporter/internal/config"
	"github.com/michael/velero-backup-reporter/internal/report"
)

// Sender handles sending backup report emails.
type Sender struct {
	cfg           config.SMTPConfig
	detailsWindow time.Duration
	template      *template.Template
}

// NewSender creates a new email Sender.
func NewSender(cfg config.SMTPConfig, emailCfg config.EmailConfig) (*Sender, error) {
	funcMap := template.FuncMap{
		"formatTime": func(t *time.Time) string {
			if t == nil {
				return "-"
			}
			return t.Format("2006-01-02 15:04:05 UTC")
		},
		"formatTimeVal": func(t time.Time) string {
			if t.IsZero() {
				return "-"
			}
			return t.Format("2006-01-02 15:04:05 UTC")
		},
		"formatDuration": func(d time.Duration) string {
			if d == 0 {
				return "-"
			}
			return d.Round(time.Second).String()
		},
		"formatRate": func(rate float64) string {
			return fmt.Sprintf("%.1f%%", rate)
		},
		"statusColor": func(status string) string {
			switch status {
			case "Completed":
				return "#22c55e"
			case "Failed":
				return "#ef4444"
			case "Missed":
				return "#ef4444"
			case "PartiallyFailed":
				return "#f59e0b"
			case "InProgress":
				return "#3b82f6"
			default:
				switch {
				case strings.Contains(status, "PartiallyFailed"):
					return "#f59e0b"
				case strings.Contains(status, "Failed"):
					return "#ef4444"
				case status == "New" || status == "Queued" || status == "ReadyToStart" || status == "WaitingForPluginOperations" || status == "Finalizing":
					return "#3b82f6"
				default:
					return "#6b7280"
				}
			}
		},
	}

	tmpl, err := template.New("email").Funcs(funcMap).Parse(emailTemplate)
	if err != nil {
		return nil, fmt.Errorf("parsing email template: %w", err)
	}

	return &Sender{
		cfg:           cfg,
		detailsWindow: emailDetailsWindowOrDefault(emailCfg.DetailsWindow),
		template:      tmpl,
	}, nil
}

// Send sends a backup report email.
func (s *Sender) Send(rpt report.BackupReport) error {
	rpt.Backups = filterBackupDetailsWithinWindow(rpt.Backups, time.Now(), s.detailsWindow)

	var body bytes.Buffer
	if err := s.template.Execute(&body, rpt); err != nil {
		return fmt.Errorf("rendering email template: %w", err)
	}

	subject := fmt.Sprintf("Velero Backup Report - %s", rpt.GeneratedAt.Format("2006-01-02"))

	msg := buildMessage(s.cfg.From, s.cfg.To, subject, body.String())

	addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)

	c, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("connecting to SMTP server: %w", err)
	}
	defer c.Close()

	if err := c.Hello("localhost"); err != nil {
		return fmt.Errorf("SMTP EHLO: %w", err)
	}

	if s.cfg.TLS {
		if err := c.StartTLS(&tls.Config{ServerName: s.cfg.Host}); err != nil {
			return fmt.Errorf("SMTP STARTTLS: %w", err)
		}
	}

	if s.cfg.Username != "" {
		auth := smtp.PlainAuth("", s.cfg.Username, s.cfg.Password, s.cfg.Host)
		if err := c.Auth(auth); err != nil {
			return fmt.Errorf("SMTP auth: %w", err)
		}
	}

	if err := c.Mail(s.cfg.From); err != nil {
		return fmt.Errorf("SMTP MAIL FROM: %w", err)
	}
	for _, rcpt := range s.cfg.To {
		if err := c.Rcpt(rcpt); err != nil {
			return fmt.Errorf("SMTP RCPT TO %s: %w", rcpt, err)
		}
	}

	w, err := c.Data()
	if err != nil {
		return fmt.Errorf("SMTP DATA: %w", err)
	}
	if _, err := w.Write([]byte(msg)); err != nil {
		return fmt.Errorf("writing email body: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("closing email body: %w", err)
	}

	if err := c.Quit(); err != nil {
		return fmt.Errorf("SMTP QUIT: %w", err)
	}

	log.Printf("INFO: backup report email sent to %s", strings.Join(s.cfg.To, ", "))
	return nil
}

func emailDetailsWindowOrDefault(d time.Duration) time.Duration {
	if d <= 0 {
		return 24 * time.Hour
	}
	return d
}

func filterBackupDetailsWithinWindow(backups []report.BackupDetail, now time.Time, window time.Duration) []report.BackupDetail {
	if len(backups) == 0 {
		return backups
	}

	cutoff := now.Add(-window)
	filtered := make([]report.BackupDetail, 0, len(backups))

	for _, b := range backups {
		runTime := b.StartTime
		if runTime == nil {
			runTime = b.CompletionTime
		}

		if runTime == nil {
			if isNotStartedStatus(b.Status) {
				filtered = append(filtered, b)
			}
			continue
		}

		if !runTime.Before(cutoff) {
			filtered = append(filtered, b)
		}
	}

	return filtered
}

func isNotStartedStatus(status string) bool {
	switch status {
	case "New", "Queued", "ReadyToStart", "FailedValidation", "Missed":
		return true
	default:
		return false
	}
}

func buildMessage(from string, to []string, subject, htmlBody string) string {
	var b strings.Builder
	b.WriteString("From: " + from + "\r\n")
	b.WriteString("To: " + strings.Join(to, ", ") + "\r\n")
	b.WriteString("Subject: " + subject + "\r\n")
	b.WriteString("MIME-Version: 1.0\r\n")
	b.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n")
	b.WriteString("\r\n")
	b.WriteString(htmlBody)
	return b.String()
}

const emailTemplate = `<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"></head>
<body style="margin: 0; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background: #eaf2ff; padding: 20px 12px; color: #0f172a;">
<div style="max-width: 700px; margin: 0 auto; background: #ffffff; border-radius: 12px; overflow: hidden; box-shadow: 0 8px 24px rgba(15,23,42,0.12); border: 1px solid #dbeafe;">

<div style="background: #0f172a; background-image: linear-gradient(120deg, #0f172a 0%, #1d4ed8 55%, #0ea5e9 100%); color: #fff; padding: 22px 24px;">
	<h1 style="margin: 0; font-size: 20px; letter-spacing: 0.2px;">Velero Backup Report</h1>
	<p style="margin: 6px 0 0; color: #dbeafe; font-size: 14px;">Generated at {{formatTimeVal .GeneratedAt}}</p>
</div>
<div style="height: 6px; background: #f59e0b;"></div>

<div style="padding: 24px;">
	<h2 style="margin: 0 0 12px; font-size: 16px; color: #1d4ed8;">Summary</h2>
	<table style="width: 100%; border-collapse: collapse; margin-bottom: 24px; background: #f8fbff; border: 1px solid #dbeafe;">
        <tr>
			<td style="padding: 8px 12px; background: #e0e7ff; border: 1px solid #dbeafe;"><strong>Total Backups</strong></td>
			<td style="padding: 8px 12px; border: 1px solid #dbeafe;">{{.Summary.TotalBackups}}</td>
			<td style="padding: 8px 12px; background: #dcfce7; border: 1px solid #dbeafe;"><strong>Completed</strong></td>
			<td style="padding: 8px 12px; border: 1px solid #dbeafe; color: #166534; font-weight: 600;">{{.Summary.Completed}}</td>
        </tr>
        <tr>
			<td style="padding: 8px 12px; background: #fee2e2; border: 1px solid #dbeafe;"><strong>Failed</strong></td>
			<td style="padding: 8px 12px; border: 1px solid #dbeafe; color: #991b1b; font-weight: 600;">{{.Summary.Failed}}</td>
			<td style="padding: 8px 12px; background: #fef3c7; border: 1px solid #dbeafe;"><strong>Partially Failed</strong></td>
			<td style="padding: 8px 12px; border: 1px solid #dbeafe; color: #92400e; font-weight: 600;">{{.Summary.PartiallyFailed}}</td>
        </tr>
		<tr>
			<td style="padding: 8px 12px; background: #fee2e2; border: 1px solid #dbeafe;"><strong>Missing / Not Started</strong></td>
			<td style="padding: 8px 12px; border: 1px solid #dbeafe; color: #991b1b; font-weight: 600;">{{.Summary.NotStarted}}</td>
			<td style="padding: 8px 12px; background: #f3f4f6; border: 1px solid #dbeafe;"><strong></strong></td>
			<td style="padding: 8px 12px; border: 1px solid #dbeafe;"></td>
		</tr>
        <tr>
			<td style="padding: 8px 12px; background: #ecfeff; border: 1px solid #dbeafe;"><strong>Last Successful</strong></td>
			<td style="padding: 8px 12px; border: 1px solid #dbeafe;" colspan="3">{{formatTime .Summary.LastSuccessful}}</td>
        </tr>
        <tr>
			<td style="padding: 8px 12px; background: #fef2f2; border: 1px solid #dbeafe;"><strong>Last Failed</strong></td>
			<td style="padding: 8px 12px; border: 1px solid #dbeafe;" colspan="3">{{formatTime .Summary.LastFailed}}</td>
        </tr>
    </table>

    {{if .ScheduleSummaries}}
	<h2 style="margin: 0 0 12px; font-size: 16px; color: #0f766e;">Schedules</h2>
	<table style="width: 100%; border-collapse: collapse; margin-bottom: 24px; background: #f0fdfa; border: 1px solid #99f6e4;">
		<tr style="background: #14b8a6; color: #ffffff;">
			<th style="padding: 8px 12px; text-align: left; border: 1px solid #0d9488;">Schedule</th>
			<th style="padding: 8px 12px; text-align: left; border: 1px solid #0d9488;">Last Status</th>
			<th style="padding: 8px 12px; text-align: left; border: 1px solid #0d9488;">Total</th>
			<th style="padding: 8px 12px; text-align: left; border: 1px solid #0d9488;">Success Rate</th>
        </tr>
        {{range .ScheduleSummaries}}
        <tr>
			<td style="padding: 8px 12px; border: 1px solid #99f6e4;">{{.ScheduleName}}</td>
			<td style="padding: 8px 12px; border: 1px solid #99f6e4; color: {{statusColor .LastBackupStatus}}; font-weight: 600;">{{if .LastBackupStatus}}{{.LastBackupStatus}}{{else}}-{{end}}</td>
			<td style="padding: 8px 12px; border: 1px solid #99f6e4;">{{.TotalBackups}}</td>
			<td style="padding: 8px 12px; border: 1px solid #99f6e4;">{{formatRate .SuccessRate}}</td>
        </tr>
        {{end}}
    </table>
    {{end}}

    {{if .Backups}}
	<h2 style="margin: 0 0 12px; font-size: 16px; color: #7c2d12;">Backup Details (Last 24 Hours)</h2>
	{{range .Backups}}
	<table style="width: 100%; border-collapse: collapse; margin-bottom: 14px; border: 1px solid #fdba74; border-radius: 10px; overflow: hidden; background: #fff7ed;">
		<tr>
			<td style="padding: 12px; border: 1px solid #fed7aa; vertical-align: top;">
				<div style="font-size: 14px; font-weight: 600; line-height: 1.4; word-break: break-word;">{{.Name}}</div>
				{{if .ScheduleName}}
				<div style="margin-top: 4px; font-size: 12px; color: #6b7280;">Schedule: {{.ScheduleName}}</div>
				{{end}}
			</td>
			<td style="padding: 12px; border: 1px solid #fed7aa; vertical-align: top; text-align: right; white-space: nowrap; color: {{statusColor .Status}}; font-size: 13px; font-weight: 700; background: #fff1e6;">
				{{.Status}}
			</td>
		</tr>
		<tr>
			<td style="padding: 8px 12px; border: 1px solid #fed7aa; background: #ffedd5; font-size: 12px; font-weight: 600; color: #7c2d12; width: 35%;">Started</td>
			<td style="padding: 8px 12px; border: 1px solid #fed7aa; font-size: 13px; line-height: 1.4;">{{formatTime .StartTime}}</td>
		</tr>
		<tr>
			<td style="padding: 8px 12px; border: 1px solid #fed7aa; background: #ffedd5; font-size: 12px; font-weight: 600; color: #7c2d12;">Duration</td>
			<td style="padding: 8px 12px; border: 1px solid #fed7aa; font-size: 13px;">{{formatDuration .Duration}}</td>
		</tr>
		<tr>
			<td style="padding: 8px 12px; border: 1px solid #fed7aa; background: #ffedd5; font-size: 12px; font-weight: 600; color: #7c2d12;">Items Backed Up</td>
			<td style="padding: 8px 12px; border: 1px solid #fed7aa; font-size: 13px;">{{.ItemsBackedUp}} / {{.TotalItems}}</td>
		</tr>
		<tr>
			<td style="padding: 8px 12px; border: 1px solid #fed7aa; background: #ffedd5; font-size: 12px; font-weight: 600; color: #7c2d12;">Warnings / Errors</td>
			<td style="padding: 8px 12px; border: 1px solid #fed7aa; font-size: 13px;">{{.Warnings}} / {{.Errors}}</td>
		</tr>
		<tr>
			<td style="padding: 8px 12px; border: 1px solid #fed7aa; background: #ffedd5; font-size: 12px; font-weight: 600; color: #7c2d12;">Failure Reason</td>
			<td style="padding: 8px 12px; border: 1px solid #fed7aa; font-size: 13px;">{{if .FailureReason}}{{.FailureReason}}{{else}}-{{end}}</td>
		</tr>
		<tr>
			<td style="padding: 8px 12px; border: 1px solid #fed7aa; background: #ffedd5; font-size: 12px; font-weight: 600; color: #7c2d12;">Validation Errors</td>
			<td style="padding: 8px 12px; border: 1px solid #fed7aa; font-size: 13px;">{{if .ValidationErrors}}{{range .ValidationErrors}}{{.}}<br>{{end}}{{else}}-{{end}}</td>
		</tr>
	</table>
	{{end}}
    {{end}}
</div>

</div>
</body>
</html>`
