package utils

import (
	"context"
	"fmt"
	"net/smtp"
	"strings"

	"github.com/TmzFranck/books-api-golang/internal/jobs"
	"github.com/spf13/viper"
)

type Mail struct {
	Sender  string
	To      []string
	Subject string
	Body    string
}
type SMTPConfig struct {
	SMTPHost string
	SMTPPort uint
	Username string
	Password string
	From     string
	FromName string
}

func (s *SMTPConfig) Validate() error {
	if s.SMTPHost == "" {
		return fmt.Errorf("SMTP host is required")
	}
	if s.SMTPPort == 0 {
		return fmt.Errorf("SMTP port is required")
	}
	if s.Username == "" {
		return fmt.Errorf("SMTP username is required")
	}
	if s.Password == "" {
		return fmt.Errorf("SMTP password is required")
	}
	if s.From == "" {
		return fmt.Errorf("SMTP from is required")
	}
	if s.FromName == "" {
		return fmt.Errorf("SMTP form name is required")
	}
	return nil
}

func (m *Mail) Validate() error {
	if m.Sender == "" {
		return fmt.Errorf("sender cannot be empty")
	}
	if len(m.To) == 0 {
		return fmt.Errorf("to cannot be empty")
	}
	if m.Subject == "" {
		return fmt.Errorf("subject cannot be empty")
	}
	if m.Body == "" {
		return fmt.Errorf("body cannot be empty")
	}
	return nil
}

func sendViaSMTP(ctx context.Context, config *SMTPConfig, mail *Mail) error {
	auth := smtp.PlainAuth("", config.Username, config.Password, config.SMTPHost)

	addr := fmt.Sprintf("%s:%d", config.SMTPHost, config.SMTPPort)

	headers := buildEmailHeaders(config, mail)

	message := headers + "\r\n\r\n" + mail.Body

	done := make(chan error, 1)

	go func() {
		done <- smtp.SendMail(addr, auth, config.From, mail.To, []byte(message))
	}()

	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return fmt.Errorf("email send cancelled: %w", ctx.Err())
	}

}

func buildEmailHeaders(config *SMTPConfig, mail *Mail) string {
	from := config.From
	if config.FromName != "" {
		from = fmt.Sprintf("%s <%s>", config.FromName, config.From)
	}
	return fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/html; charset=UTF-8",
		from,
		strings.Join(mail.To, ", "),
		mail.Subject,
	)
}

func getPayloadString(payload map[string]any, key string) string {
	val, ok := payload[key].(string)
	if !ok {
		return ""
	}
	return val
}

func getPayloadStrings(payload map[string]any, key string) []string {
	val, ok := payload[key].([]any)
	if !ok {
		if strSlice, ok := payload[key].([]string); ok {
			return strSlice
		}
		return []string{}
	}

	result := make([]string, 0, len(val))
	for _, v := range val {
		if s, ok := v.(string); ok {
			result = append(result, s)
		}
	}
	return result
}

func SendMail(worker *jobs.WokerPool, mail *Mail) error {
	if err := mail.Validate(); err != nil {
		return fmt.Errorf("invalid mail: %w", err)
	}

	payload := map[string]any{
		"sender":  mail.Sender,
		"to":      mail.To,
		"subject": mail.Subject,
		"body":    mail.Body,
	}
	job := &jobs.Job{
		Type:    "SendMail",
		Payload: payload,
	}

	return worker.Submit(job)
}

func Send(ctx context.Context, job *jobs.Job) error {
	mail := &Mail{
		Sender:  getPayloadString(job.Payload, "sender"),
		To:      getPayloadStrings(job.Payload, "to"),
		Subject: getPayloadString(job.Payload, "subject"),
		Body:    getPayloadString(job.Payload, "body"),
	}

	if err := mail.Validate(); err != nil {
		return fmt.Errorf("mail validation failed: %w", err)
	}

	config := &SMTPConfig{
		SMTPHost: viper.GetString("SMTP.server"),
		SMTPPort: viper.GetUint("SMTP.port"),
		Username: viper.GetString("SMTP.username"),
		Password: viper.GetString("SMTP.password"),
		From:     viper.GetString("SMTP.from"),
		FromName: viper.GetString("SMTP.fromName"),
	}

	if err := config.Validate(); err != nil {
		return fmt.Errorf("SMTP config invalid: %w", err)
	}

	return sendViaSMTP(ctx, config, mail)

}
