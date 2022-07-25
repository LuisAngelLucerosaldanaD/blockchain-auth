package send_grid

import (
	"testing"
)

func TestModel_SendMail(t *testing.T) {
	to := To{
		Name: "Mr. jonhjar guerrero",
		Mail: "guerrerojv@gmail.com",
	}
	var too []To
	too = append(too, to)
	m := Model{
		Tos:         too,
		FromMail:    "no-reply@e-capture.co",
		FromName:    "ecapture",
		Subject:     "test-sendgrid",
		HTMLContent: "<h1>Hola amigos sendgrid</h1><p>Esta es una prueba<p>",
		Attachments: nil,
	}
	err := m.SendMail()
	if err != nil {
		t.Fatalf("no se envío correo SendMail: %v", err)
	}
}
