package email



import (
	"blion-auth/internal/env"
	"blion-auth/internal/logger"
	"blion-auth/internal/template"
	"crypto/tls"
	"gopkg.in/gomail.v2"

)

func Send(tos []string, parameters map[string]string, pathTemplate, subject string) error {
	e := env.NewConfiguration()
	body, err := template.GenerateTemplateMail(parameters, pathTemplate)
	if err != nil {
		logger.Error.Printf("couldn't generate body in notification email")
		return err
	}
	myeMail := &Email{
		From:    e.Smtp.Email,
		To:      tos,
		Subject: subject,
		Body:    body,
	}
	err = myeMail.email()
	if err != nil {
		logger.Error.Printf("couldn't send email NotificationMail: %V", err)
		return err
	}
	return nil
}

func (m *Email) email() error {
	c := env.NewConfiguration()
	e := gomail.NewMessage()
	e.SetHeader("From", m.From)
	e.SetHeader("To", m.To...)
	e.SetHeader("Cc", m.CC...)
	e.SetHeader("Subject", m.Subject)
	e.SetBody("text/html", m.Body)
	if len(m.Attach) > 0 {
		//m.Attach(e.Attach)
	}

	for _, v := range m.Attachments {
		e.Attach(v)
	}

	mp := c.Smtp.Port
	d := gomail.NewDialer(c.Smtp.Host, mp, c.Smtp.Email, c.Smtp.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	err := d.DialAndSend(e)
	if err != nil {
		logger.Error.Printf("couldn't emil to: %s, subject: %s, %v", m.To, m.Subject, err)
		return err
	}

	return nil
}

