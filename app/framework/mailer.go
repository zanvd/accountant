package framework

import (
	"bytes"
	"net/smtp"
)

var baseMailTmpls = []string{
	"templates/mail/base.gohtml",
}

type MailerConsumer interface {
	GetMailTemplates() map[string]string
}

type Mail struct {
	Body    []byte
	From    string
	Subject string
	To      []string
}

func (m *Mail) RenderBody(r *Routes, rd *RequestData, tb *TemplateBuilder) (err error) {
	b := new(bytes.Buffer)
	if err = tb.Render(r, rd, b); err != nil {
		return
	}
	m.Body = b.Bytes()
	return
}

type Mailer struct {
	Auth        smtp.Auth
	DefaultFrom string
	Host        string
	Port        string
}

func NewMailer(config *Config) *Mailer {
	return &Mailer{
		Auth:        smtp.PlainAuth("", config.Mail.Username, config.Mail.Password, config.Mail.Host),
		DefaultFrom: config.Mail.DefaultSender,
		Host:        config.Mail.Host,
		Port:        config.Mail.Port,
	}
}

func (m *Mailer) Send(mail Mail) error {
	head := []byte(
		"Subject: " + mail.Subject + "\n" + "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n",
	)
	return smtp.SendMail(
		m.Host+":"+m.Port, m.Auth, mail.From, mail.To, bytes.Join([][]byte{head, mail.Body}, []byte("\n")),
	)
}
