package send_grid

import (
	"blion-auth/internal/env"
	"errors"
	"fmt"
	"github.com/sendgrid/sendgrid-go"

	"blion-auth/internal/logger"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type Model struct {
	Tos         []To
	FromMail    string
	FromName    string
	Subject     string
	HTMLContent string
	Attachments []*mail.Attachment
}

type To struct {
	Name string
	Mail string
}

func (m *Model) AddAttachment(as ...*mail.Attachment) {
	m.Attachments = append(as)
}

func (m *Model) SendMail() error {
	c := env.NewConfiguration()
	var privateKey = c.SendGrid.Key

	client := sendgrid.NewSendClient(privateKey)
	from := mail.NewEmail(m.FromName, m.FromMail)
	per := mail.NewPersonalization()
	for _, v := range m.Tos {
		to := mail.NewEmail(v.Name, v.Mail)
		per.AddTos(to)
	}

	message := mail.NewV3Mail()
	message.SetFrom(from)
	message.AddPersonalizations(per)
	message.Subject = m.Subject

	content := mail.NewContent("text/html", m.HTMLContent)
	message.AddContent(content)
	for _, v := range m.Attachments {
		message.AddAttachment(v)
	}

	r, err := client.Send(message)
	if err != nil {
		logger.Error.Printf("no se pudo enviar el correo: %v", err)
		return err
	}

	if r.StatusCode != 202 {
		logger.Error.Printf("enviando correo: %s\nHeaders Sendgrid: %v", r.Body, r.Headers)
		return errors.New(fmt.Sprintf("error al enviar el correo: %s", r.Body))
	}

	return nil
}
