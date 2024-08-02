package mailer

import (
	"bytes"
	"embed"
	"html/template"
	"time"

	mail "github.com/wneessen/go-mail"
)

//go:embed "templates"
var templateFS embed.FS

type Mailer struct {
	dialer *mail.Client
	sender string
}

func New(host string, port int, username, password, sender string) (Mailer, error) {
	c, err := mail.NewClient(host,
		mail.WithSMTPAuth(mail.SMTPAuthPlain), mail.WithTLSPortPolicy(mail.TLSMandatory),
		mail.WithUsername(username), mail.WithPassword(password),
	)
	if err != nil {
		return Mailer{}, err
	}

	return Mailer{
		dialer: c,
		sender: sender,
	}, nil
}
func (m Mailer) Send(recipient, templateFile string, data interface{}) error {
	tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	plainBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(plainBody, "plainBody", data)
	if err != nil {
		return err
	}

	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return err
	}

	msg := mail.NewMsg()
	msg.From(m.sender)
	msg.To(recipient)
	msg.Subject(subject.String())
	msg.SetBodyString("text/plain", plainBody.String())
	// err = msg.AddAlternativeHTMLTemplate(tmpl, data)
	// if err != nil {
	// 	return err
	// }
	for i := 0; i < 3; i++ {
		err = m.dialer.DialAndSend(msg)
		if nil == err {
			return nil
		}
		time.Sleep(1 * time.Second)
	}
	return err
}
