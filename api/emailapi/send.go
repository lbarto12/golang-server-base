package emailapi

import (
	"golang-server-base/api/emailapi/models"

	"gopkg.in/gomail.v2"
)

func SendEmail(options models.EmailOptions) error {
	m := gomail.NewMessage()
	m.SetHeader("From", dialer.Username)
	m.SetHeader("To", options.To)
	m.SetHeader("Subject", options.Subject)
	m.SetBody(options.BodyType, options.Body)

	return dialer.DialAndSend(m)
}

// example use
// log.Fatal(emailapi.SendEmail(models.EmailOptions{
// 	To:       "recipient@host.com",
// 	Subject:  "testing email",
// 	BodyType: "text/plain",
// 	Body:     "did you get this?",
// }))
