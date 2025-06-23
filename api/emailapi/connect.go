package emailapi

import (
	"golang-server-base/api/emailapi/models"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

var dialer *gomail.Dialer

func EnvGetOptions() models.DialerOptions {
	smtpHost, ok := os.LookupEnv("SMTP_HOST")
	if !ok {
		panic("SMTP_HOST environment variable not set")
	}

	smtpPort, ok := os.LookupEnv("SMTP_PORT")
	if !ok {
		panic("SMTP_PORT environment variable not set")
	}

	smtpSender, ok := os.LookupEnv("SMTP_SENDER")
	if !ok {
		panic("SMTP_SENDER environment variable not set")
	}

	smtpPassword, ok := os.LookupEnv("SMTP_PASSWORD")
	if !ok {
		panic("SMTP_PASSWORD environment variable not set")
	}

	return models.DialerOptions{
		SMTPServer: smtpHost,
		Port:       smtpPort,
		Sender:     smtpSender,
		Password:   smtpPassword,
	}
}

func Init(options models.DialerOptions) error {
	port, err := strconv.ParseInt(options.Port, 10, 32)
	if err != nil {
		return err
	}
	dialer = gomail.NewDialer(options.SMTPServer, int(port), options.Sender, options.Password)
	return nil
}
