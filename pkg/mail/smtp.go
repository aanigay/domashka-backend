package smtpmail

import (
	"bytes"
	"context"
	"crypto/tls"
	"domashka-backend/config"
	"domashka-backend/internal/entity/notifications"
	"fmt"
	"html/template"
	"net/smtp"
	"time"
)

const emailTemplate = `<!DOCTYPE html>
<html lang="ru">
<head>
<meta charset="UTF-8">
<title>{{.Title}}</title>
<style>
body {
	font-family: Arial, sans-serif;
	margin: 20px;
}
.footer {
	margin-top: 20px;
	color: #888;
	font-size: 12px;
}
</style>
</head>
<body>
<h2>Привет, {{.UserName}}!</h2>
<p>{{.Body}}</p>
<div class="footer">
С уважением,<br>
{{.Footer}}
</div>
</body>
</html>`

type SMTPClient struct {
	Host     string
	Port     string
	Username string
	Password string

	MaxRetries int
	RetryDelay time.Duration

	UseSSL bool
}

func New(config *config.SMTPEmailConfig) *SMTPClient {
	return &SMTPClient{
		Host:       config.Host,
		Port:       config.Port,
		Username:   config.Email,
		Password:   config.Password,
		UseSSL:     true,
		MaxRetries: config.MaxRetries,
		RetryDelay: config.RetryDelay,
	}
}

// SendEmail отправляет email через SMTP
func (e *SMTPClient) SendEmail(ctx context.Context, to, subject string, data notifications.EmailData) (int, error) {
	var err error
	attempts := 0
	for attempts = 1; attempts <= e.MaxRetries; attempts++ {
		err = e.send(to, subject, data)
		if err == nil {
			return attempts, nil
		}
		select {
		case <-ctx.Done():
			return attempts, ctx.Err()
		case <-time.After(e.RetryDelay):
		}
	}
	return attempts, fmt.Errorf("достигнут лимит попыток отправки: %w", err)
}

func (e *SMTPClient) send(to, subject string, data notifications.EmailData) error {
	from := e.Username

	tmpl, err := template.New("email").Parse(emailTemplate)
	if err != nil {
		return fmt.Errorf("ошибка парсинга шаблона: %w", err)
	}

	var bodyBuffer bytes.Buffer
	if err := tmpl.Execute(&bodyBuffer, data); err != nil {
		return fmt.Errorf("ошибка рендеринга шаблона: %w", err)
	}

	headers := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n",
		from, to, subject)

	message := []byte(headers + bodyBuffer.String())

	if e.UseSSL {
		return e.sendWithSSL(from, to, message)
	}
	return e.sendWithSTARTTLS(from, to, message)
}

// sendWithSSL отправляет email через SSL (порт 465)
func (e *SMTPClient) sendWithSSL(from, to string, message []byte) error {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         e.Host,
	}

	conn, err := tls.Dial("tcp", e.Host+":"+e.Port, tlsConfig)
	if err != nil {
		return fmt.Errorf("ошибка подключения к SMTP через SSL: %w", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, e.Host)
	if err != nil {
		return fmt.Errorf("ошибка создания SMTP клиента: %w", err)
	}

	auth := smtp.PlainAuth("", e.Username, e.Password, e.Host)
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("ошибка авторизации: %w", err)
	}

	if err := client.Mail(from); err != nil {
		return fmt.Errorf("ошибка MAIL FROM: %w", err)
	}
	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("ошибка RCPT TO: %w", err)
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("ошибка DATA: %w", err)
	}
	_, err = w.Write(message)
	if err != nil {
		return fmt.Errorf("ошибка записи тела письма: %w", err)
	}
	w.Close()

	client.Quit()
	return nil
}

// sendWithSTARTTLS отправляет email через STARTTLS (порт 587)
func (e *SMTPClient) sendWithSTARTTLS(from, to string, message []byte) error {
	client, err := smtp.Dial(e.Host + ":" + e.Port)
	if err != nil {
		return fmt.Errorf("ошибка подключения к SMTP через STARTTLS: %w", err)
	}
	defer client.Close()

	if err = client.StartTLS(&tls.Config{ServerName: e.Host}); err != nil {
		return fmt.Errorf("ошибка при старте TLS: %w", err)
	}
	auth := smtp.PlainAuth("", e.Username, e.Password, e.Host)
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("ошибка авторизации: %w", err)
	}

	if err = client.Mail(from); err != nil {
		return fmt.Errorf("ошибка MAIL FROM: %w", err)
	}
	if err = client.Rcpt(to); err != nil {
		return fmt.Errorf("ошибка RCPT TO: %w", err)
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("ошибка DATA: %w", err)
	}
	_, err = w.Write(message)
	if err != nil {
		return fmt.Errorf("ошибка записи тела письма: %w", err)
	}
	w.Close()

	client.Quit()
	return nil
}
