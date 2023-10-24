package utils

import (
	"bytes"
	"crypto/tls"
	"github.com/example/golang-test/config"
	"github.com/example/golang-test/models"
	"github.com/k3a/html2text"
	"gopkg.in/gomail.v2"
	"log"
)

type EmailData struct {
	URL       string
	FirstName string
	Subject   string
}

// ? Email template parser

func SendEmail(user *models.DBResponse, data *EmailData, templateName string) error {
	config, err := config.LoadConfig1(".")

	if err != nil {
		log.Fatal("could not load config", err)
	}

	// Sender data.
	from := config.EmailFrom
	smtpPass := config.SMTPPass
	smtpUser := config.SMTPUser
	to := user.Email
	smtpHost := config.SMTPHost
	smtpPort := config.SMTPPort

	var body bytes.Buffer

	//template, err := ParseTemplateDir("templates")
	//if err != nil {
	//	log.Fatal("Could not parse template", err)
	//}
	//
	//template = template.Lookup(templateName)
	//template.Execute(&body, &data)
	//fmt.Println(template.Name())

	m := gomail.NewMessage()

	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", data.Subject)
	m.SetBody("text/html", body.String())
	m.AddAlternative("text/plain", html2text.HTML2Text(body.String()))

	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Send Email
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
