package main

import (
	"bytes"
	"crypto/tls"
	"gopkg.in/gomail.v2"
	"html/template"
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

func (m *Mail) SendSMTPMessage(msg Message) error {
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

	message := gomail.NewMessage()

	message.SetHeader("From", msg.From)
	message.SetHeader("To", msg.To)
	message.SetHeader("Subject", msg.Subject)
	message.SetBody("text/plain", plainMessage)
	message.AddAlternative("text/html", formattedMessage)

	if len(msg.Attachments) > 0 {
		for _, v := range msg.Attachments {
			message.Attach(v)
		}
	}

	server := gomail.NewDialer(
		m.Host,
		m.Port,
		m.Username,
		m.Password,
	)

	m.setEncryption(server)

	if err = server.DialAndSend(message); err != nil {
		return err
	}

	return nil

}

func (m *Mail) buildHTMLMessage(msg Message) (string, error) {
	// TODO: Getting no such file or directory should be fixed
	templateToRender := "./templates/mail.html.gohtml"

	t, err := template.New("email-html").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err = t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", err
	}

	formattedMessage := tpl.String()

	return formattedMessage, nil
}

func (m *Mail) buildPlainTextMessage(msg Message) (string, error) {
	// TODO: Getting no such file or directory should be fixed
	templateToRender := "./templates/mail.plain.gohtml"

	t, err := template.New("email-plain").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err = t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", err
	}

	plainMessage := tpl.String()

	return plainMessage, nil
}

func (m *Mail) setEncryption(dialer *gomail.Dialer) {
	switch m.Encryption {
	case "SSL":
		dialer.SSL = true
	case "TLS":
		dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	case "None", "":
		dialer.SSL = false
		dialer.TLSConfig = nil
	default:
		dialer.SSL = true
	}
}
