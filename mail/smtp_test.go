package mail

import (
	"os"
	"testing"
)

func TestSendMail(t *testing.T) {
	from := os.Getenv("FROM")
	password := os.Getenv("PASSWORD")
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	address := host + ":" + port
	smtp := NewSmtp(from, address, "", from, password, host)
	err := smtp.SendMail(&SendMailRequest{
		To:      []string{os.Getenv("TO")},
		Subject: []byte("Test subject"),
		Text:    []byte("Test body"),
		Html:    nil,
		Tls:     true,
	})
	if err != nil {
		t.Error(err)
	}
}
