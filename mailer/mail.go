package mailer

import (
	"bytes"
	"fmt"
	"text/template"
)

type Mail struct {
	// Domain is the url of the sender
	Domain string

	// Templates holds the path to the templates directory
	Templates   string
	Host        string
	Port        string
	Username    string
	Password    string
	Encryption  string
	FromAddress string
	FromName    string
	Jobs        chan Message

	// Results holds the results received after sending a mail
	Results chan Result
	API     string
	APIKey  string
	APIUrl  string
}

type Message struct {
	From        string
	FromName    string
	To          string
	Subject     string
	Template    string
	Attachments []string
	// Data is the data passed to the template when its rendered
	Data interface{}
}

type Result struct {
	Success bool
	Error   error
}

// listens for mail
func (m *Mail) ListenForMail() {
	for {
		// listen for incoming messages on the jobs channel
		msg := <-m.Jobs
		err := m.Send(msg)
		if err != nil {
			// send error to results
			m.Results <- Result{false, err}
		} else {
			m.Results <- Result{true, nil}
		}
	}
}

func (m *Mail) Send(msg Message) error {
	// TODO: using smtp (legacy) or api (like mail gun)?
	return m.SendSMTPMessage(msg)
}

func (m *Mail) SendSMTPMessage(msg Message) error {

	return nil
}

func (m *Mail) buildHTMLMessage(msg Message) (string, error) {
	templateToRender := fmt.Sprintf("%s/%s.html.tmpl", m.Templates, msg.Template)

	t, err := template.New("email-html").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err = t.ExecuteTemplate(&tpl, "body", msg.Data); err != nil {
		return "", err
	}

	formattedMessage := tpl.String()

	return formattedMessage, nil
}

func (m *Mail) buildPlainTextMessage(msg Message) (string, error) {
	templateToRender := fmt.Sprintf("%s/%s.plain.tmpl", m.Templates, msg.Template)

	t, err := template.New("email-plaintext").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err = t.ExecuteTemplate(&tpl, "body", msg.Data); err != nil {
		return "", err
	}

	plainMessage := tpl.String()

	return plainMessage, nil
}
