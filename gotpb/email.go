package gotpb

import (
	"log"

	mail "github.com/xhit/go-simple-mail/v2"
)

// Email represent a simple interface with methods the methds
// required by this application to send emails
type Email interface {
	AddTo(...string) *mail.Email
	Send(client *mail.SMTPClient) error
	SetBody(contenttype mail.ContentType, body string) *mail.Email
	SetFrom(from string) *mail.Email
	SetSubject(subject string) *mail.Email
}

func smtpClient(conf Config) *mail.SMTPClient {
	client := mail.NewSMTPClient()
	client.Host = conf.EmailClientConfig.Host
	client.Port = conf.EmailClientConfig.Port
	client.Username = conf.EmailClientConfig.Username
	client.Password = conf.EmailClientConfig.Password
	client.Encryption = mail.EncryptionSTARTTLS
	client.KeepAlive = false

	smtp, err := client.Connect()
	if err != nil {
		log.Printf("ERROR: %v\n", err)
	}
	return smtp
}

func prepareEmail(email Email, users []User) {
	email.SetFrom("From noter <apps.brottem@gmail.com>")
	for _, user := range users {
		email.AddTo(user.Email)
	}
}

func sendEmail(email Email, conf Config) {
	client := smtpClient(conf)
	err := email.Send(client)
	panicOnErr(err)
}
