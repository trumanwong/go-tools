package mail

import (
	"crypto/tls"
	"github.com/jordan-wright/email"
	"net/smtp"
)

type Smtp struct {
	From     string
	Addr     string
	Identity string
	Username string
	Password string
	Host     string
}

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

type SendMailRequest struct {
	To      []string
	Subject []byte
	Text    []byte
	Html    []byte
	Tls     bool
}

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
