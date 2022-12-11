package gotpb

import (
	"log"

	mail "github.com/xhit/go-simple-mail/v2"
)

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
		log.Fatal(err)
	}
	return smtp
}

func prepareEmail(conf Config, users []User) *mail.Email {
	email := mail.NewMSG()
	email.SetFrom("From noter <apps.brottem@gmail.com>")
	for _, user := range users {
		email.AddTo(user.Email)
	}
	return email
}

func sendEmail(email *mail.Email, conf Config) {
	client := smtpClient(conf)
	err := email.Send(client)
	if err != nil {
		log.Printf("Could not send email because: %v", err)
	}
}
