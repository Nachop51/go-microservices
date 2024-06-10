package main

import (
	"bytes"
	"html/template"
	"time"

	"github.com/vanng822/go-premailer/premailer"
	mail "github.com/xhit/go-simple-mail/v2"
)

type Mail struct {
	Domain      string
	Host        string
	Port        int
	Username    string
	Password    string
	Encryption  string
	FromAddress string
	FromName    string
}

type Message struct {
	From        string
	FromName    string
	To          string
	Subject     string
	Attachments []string
	Data        any
	DataMap     map[string]any
}

func (m *Mail) SendSTMPMessage(msg *Message) error {
	if msg.From == "" {
		msg.From = m.FromAddress
	}

	if msg.FromName == "" {
		msg.FromName = m.FromName
	}

	data := map[string]any{
		"message": msg.Data,
	}

	msg.DataMap = data

	formattedMessage, err := m.buildHTMLMessage(msg)

	if err != nil {
		return err
	}

	plainMessage, err := m.buildPlainTextMessage(msg)

	if err != nil {
		return err
	}

	server := mail.NewSMTPClient()

	server.Host = m.Host
	server.Port = m.Port
	server.Username = m.Username
	server.Password = m.Password
	server.Encryption = m.getEncryption(m.Encryption)
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	client, err := server.Connect()

	if err != nil {
		return err
	}

	email := mail.NewMSG()

	email.SetFrom(msg.From)
	email.AddTo(msg.To)
	email.SetSubject(msg.Subject)
	email.SetBody(mail.TextPlain, plainMessage)
	email.AddAlternative(mail.TextHTML, formattedMessage)

	for _, attachment := range msg.Attachments {
		email.AddAttachment(attachment)
	}

	err = email.Send(client)

	if err != nil {
		return err
	}

	return nil
}

func (m *Mail) getEncryption(s string) mail.Encryption {
	switch s {
	case "ssl":
		return mail.EncryptionSSLTLS
	case "tls":
		return mail.EncryptionSTARTTLS
	default:
		return mail.EncryptionNone
	}
}

func (m *Mail) buildPlainTextMessage(msg *Message) (string, error) {
	tmpl := "./templates/mail-plain.html"

	t, err := template.New("mail-plain").ParseFiles(tmpl)

	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer

	err = t.ExecuteTemplate(&tpl, "body", msg.DataMap)

	if err != nil {
		return "", err
	}

	formatted := tpl.String()

	return formatted, nil
}

func (m *Mail) buildHTMLMessage(msg *Message) (string, error) {
	tmpl := "./templates/mail.html"

	t, err := template.New("mail-html").ParseFiles(tmpl)

	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer

	err = t.ExecuteTemplate(&tpl, "body", msg.DataMap)

	if err != nil {
		return "", err
	}

	formatted := tpl.String()

	formatted, err = m.inlineCSS(formatted)

	return formatted, nil
}

func (m *Mail) inlineCSS(html string) (string, error) {
	options := premailer.Options{
		RemoveClasses:     false,
		CssToAttributes:   false,
		KeepBangImportant: true,
	}

	prem, err := premailer.NewPremailerFromString(html, &options)

	if err != nil {
		return "", err
	}

	html, err = prem.Transform()

	if err != nil {
		return "", err
	}

	return html, nil
}
