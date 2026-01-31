package commons

import (
	"crypto/tls"

	"gopkg.in/gomail.v2"

	"github.com/dekaiju/go-skeleton/config"
)

func SendMail(subject string, body string, mailTo ...string) error {
	emailConfig := config.MailConfig
	if mailTo == nil {
		mailTo = config.AlertEmailList
	}
	m := gomail.NewMessage()
	m.SetHeader("From", "<"+emailConfig.User+">")
	m.SetHeader("To", mailTo...)    // Send to multiple users
	m.SetHeader("Subject", subject) // Set email subject
	m.SetBody("text/html", body)    // Set email body
	d := gomail.NewDialer(
		emailConfig.Host,
		emailConfig.Port,
		emailConfig.User,
		emailConfig.Password,
	)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	err := d.DialAndSend(m)
	return err
}
