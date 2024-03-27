package mail

import (
	"crypto/tls"
	"github.com/jordan-wright/email"
	"net/smtp"
)

// Smtp is a struct that holds the necessary information for sending an email via SMTP.
// It includes the sender's email address (From), the SMTP server address (Addr),
// the identity string for authentication (Identity), the username and password for authentication,
// and the host name for the SMTP server (Host).
type Smtp struct {
	From     string
	Addr     string
	Identity string
	Username string
	Password string
	Host     string
}

// NewSmtp is a function that creates a new instance of the Smtp struct.
// It takes the sender's email address, the SMTP server address, the identity string for authentication,
// the username and password for authentication, and the host name for the SMTP server as arguments.
// The function returns a pointer to the newly created Smtp instance.
func NewSmtp(from, addr, identity, username, password, host string) *Smtp {
	return &Smtp{
		From:     from,
		Addr:     addr,
		Identity: identity,
		Username: username,
		Password: password,
		Host:     host,
	}
}

// SendMailRequest is a struct that holds the necessary information for a mail sending request.
// It includes the recipient email addresses (To), the subject of the email (Subject),
// the plain text content of the email (Text), the HTML content of the email (Html),
// and a boolean indicating whether to use TLS for the connection (Tls).
type SendMailRequest struct {
	To      []string
	Subject []byte
	Text    []byte
	Html    []byte
	Tls     bool
}

// SendMail is a method on the Smtp struct.
// It takes a pointer to a SendMailRequest struct as an argument.
// The method creates a new email with the information from the SendMailRequest,
// and sends the email via the SMTP server specified in the Smtp struct.
// If the Tls field in the SendMailRequest is true, the method uses a TLS connection.
// The method returns an error if the email sending fails.
func (s *Smtp) SendMail(req *SendMailRequest) error {
	e := email.NewEmail()
	e.From = s.From
	e.To = req.To
	e.Subject = string(req.Subject)
	e.Text = req.Text
	e.HTML = req.Html
	if req.Tls {
		return e.SendWithTLS(s.Addr, smtp.PlainAuth(s.Identity, s.Username, s.Password, s.Host), &tls.Config{
			ServerName: s.Host,
		})
	}
	return e.Send(s.Addr, smtp.PlainAuth(s.Identity, s.Username, s.Password, s.Host))
}
