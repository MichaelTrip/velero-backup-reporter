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
	cfg      config.SMTPConfig
	template *template.Template
}

// NewSender creates a new email Sender.
func NewSender(cfg config.SMTPConfig) (*Sender, error) {
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
			case "PartiallyFailed":
				return "#f59e0b"
			case "InProgress":
				return "#3b82f6"
			default:
				return "#6b7280"
			}
		},
	}

	tmpl, err := template.New("email").Funcs(funcMap).Parse(emailTemplate)
	if err != nil {
		return nil, fmt.Errorf("parsing email template: %w", err)
	}

	return &Sender{
		cfg:      cfg,
		template: tmpl,
	}, nil
}

// Send sends a backup report email.
func (s *Sender) Send(rpt report.BackupReport) error {
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
<body style="font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background: #f5f7fa; padding: 20px;">
<div style="max-width: 700px; margin: 0 auto; background: #fff; border-radius: 8px; overflow: hidden; box-shadow: 0 1px 3px rgba(0,0,0,0.1);">

<div style="background: #1a2332; color: #fff; padding: 20px 24px;">
    <h1 style="margin: 0; font-size: 20px;">Velero Backup Report</h1>
    <p style="margin: 4px 0 0; color: #a0b4cc; font-size: 14px;">Generated at {{formatTimeVal .GeneratedAt}}</p>
</div>

<div style="padding: 24px;">
    <h2 style="margin: 0 0 16px; font-size: 16px; color: #333;">Summary</h2>
    <table style="width: 100%; border-collapse: collapse; margin-bottom: 24px;">
        <tr>
            <td style="padding: 8px 12px; background: #f8f9fb; border: 1px solid #eee;"><strong>Total Backups</strong></td>
            <td style="padding: 8px 12px; border: 1px solid #eee;">{{.Summary.TotalBackups}}</td>
            <td style="padding: 8px 12px; background: #f8f9fb; border: 1px solid #eee;"><strong>Completed</strong></td>
            <td style="padding: 8px 12px; border: 1px solid #eee; color: #166534;">{{.Summary.Completed}}</td>
        </tr>
        <tr>
            <td style="padding: 8px 12px; background: #f8f9fb; border: 1px solid #eee;"><strong>Failed</strong></td>
            <td style="padding: 8px 12px; border: 1px solid #eee; color: #991b1b;">{{.Summary.Failed}}</td>
            <td style="padding: 8px 12px; background: #f8f9fb; border: 1px solid #eee;"><strong>Partially Failed</strong></td>
            <td style="padding: 8px 12px; border: 1px solid #eee; color: #92400e;">{{.Summary.PartiallyFailed}}</td>
        </tr>
        <tr>
            <td style="padding: 8px 12px; background: #f8f9fb; border: 1px solid #eee;"><strong>Last Successful</strong></td>
            <td style="padding: 8px 12px; border: 1px solid #eee;" colspan="3">{{formatTime .Summary.LastSuccessful}}</td>
        </tr>
        <tr>
            <td style="padding: 8px 12px; background: #f8f9fb; border: 1px solid #eee;"><strong>Last Failed</strong></td>
            <td style="padding: 8px 12px; border: 1px solid #eee;" colspan="3">{{formatTime .Summary.LastFailed}}</td>
        </tr>
    </table>

    {{if .ScheduleSummaries}}
    <h2 style="margin: 0 0 16px; font-size: 16px; color: #333;">Schedules</h2>
    <table style="width: 100%; border-collapse: collapse; margin-bottom: 24px;">
        <tr style="background: #f0f2f5;">
            <th style="padding: 8px 12px; text-align: left; border: 1px solid #eee;">Schedule</th>
            <th style="padding: 8px 12px; text-align: left; border: 1px solid #eee;">Last Status</th>
            <th style="padding: 8px 12px; text-align: left; border: 1px solid #eee;">Total</th>
            <th style="padding: 8px 12px; text-align: left; border: 1px solid #eee;">Success Rate</th>
        </tr>
        {{range .ScheduleSummaries}}
        <tr>
            <td style="padding: 8px 12px; border: 1px solid #eee;">{{.ScheduleName}}</td>
            <td style="padding: 8px 12px; border: 1px solid #eee; color: {{statusColor .LastBackupStatus}};">{{if .LastBackupStatus}}{{.LastBackupStatus}}{{else}}-{{end}}</td>
            <td style="padding: 8px 12px; border: 1px solid #eee;">{{.TotalBackups}}</td>
            <td style="padding: 8px 12px; border: 1px solid #eee;">{{formatRate .SuccessRate}}</td>
        </tr>
        {{end}}
    </table>
    {{end}}

    {{if .Backups}}
    <h2 style="margin: 0 0 16px; font-size: 16px; color: #333;">Backup Details</h2>
    <table style="width: 100%; border-collapse: collapse;">
        <tr style="background: #f0f2f5;">
            <th style="padding: 8px 12px; text-align: left; border: 1px solid #eee;">Name</th>
            <th style="padding: 8px 12px; text-align: left; border: 1px solid #eee;">Status</th>
            <th style="padding: 8px 12px; text-align: left; border: 1px solid #eee;">Start</th>
            <th style="padding: 8px 12px; text-align: left; border: 1px solid #eee;">Duration</th>
            <th style="padding: 8px 12px; text-align: left; border: 1px solid #eee;">Items</th>
        </tr>
        {{range .Backups}}
        <tr>
            <td style="padding: 8px 12px; border: 1px solid #eee;">{{.Name}}</td>
            <td style="padding: 8px 12px; border: 1px solid #eee; color: {{statusColor .Status}};">{{.Status}}</td>
            <td style="padding: 8px 12px; border: 1px solid #eee;">{{formatTime .StartTime}}</td>
            <td style="padding: 8px 12px; border: 1px solid #eee;">{{formatDuration .Duration}}</td>
            <td style="padding: 8px 12px; border: 1px solid #eee;">{{.ItemsBackedUp}}/{{.TotalItems}}</td>
        </tr>
        {{end}}
    </table>
    {{end}}
</div>

</div>
</body>
</html>`
