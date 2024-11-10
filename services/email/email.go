package email

import (
	"bytes"
	"errors"
	"html/template"
	"path/filepath"

	"gopkg.in/gomail.v2"
)

type Email struct {
	From    string         `json:"string"`
	To      string         `json:"string"`
	CC      []string       `json:"cc"`
	Subject string         `json:"subject"`
	Body    string         `json:"body"`
	Dialer  *gomail.Dialer `json:"dialer"`
}

type IEmail interface {
	Send(gomail.Message) error
	CreateMessage(from string, to string, cc []string, subject string, body string) (*gomail.Message, error)
}

const (
	WELCOME_EMAIL = "Welcome to Sync"
)

type WelcomeEmail struct {
	FullName string
}

type TemplateConfig struct {
	TemplateDir  string
	TemplateFile string
}

// SetupMailer creates a new & setup new dialer
func SetupMailer(host string, port int, username string, password string) *Email {
	e := Email{}
	e.Dialer = gomail.NewDialer(host, port, username, password)
	return &e
}

// Send message to SMTP server
func (e *Email) Send(msg gomail.Message) error {
	if e.From == "" {
		return errors.New("from can't be empty")
	}
	if e.To == "" {
		return errors.New("to can't be empty")
	}
	if e.Subject == "" {
		return errors.New("subject can't be empty")
	}
	return e.Dialer.DialAndSend(&msg)
}

// CreateMessage creates a message with all the headers
func (e *Email) CreateMessage(from string, to string, cc []string, subject string, templateConfig TemplateConfig, data any) (*gomail.Message, error) {
	tPath := filepath.Join(".", templateConfig.TemplateDir, templateConfig.TemplateFile)
	tmp, err := template.New(templateConfig.TemplateFile).ParseFiles(tPath)
	if err != nil {
		return nil, err
	}
	var body bytes.Buffer
	if err := tmp.Execute(&body, data); err != nil {
		return nil, err
	}
	m := gomail.NewMessage()
	if len(cc) != 0 {
		m.SetHeader("Cc", cc...)
	}
	m.SetHeader("To", to)
	m.SetHeader("From", from)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body.String())
	return m, nil
}
