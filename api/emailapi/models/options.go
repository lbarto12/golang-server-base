package models

type DialerOptions struct {
	SMTPServer string
	Port       string
	Sender     string
	Password   string
}

// TODO: implement attachments etc.
type EmailOptions struct {
	To       string
	Subject  string
	BodyType string
	Body     string
}
