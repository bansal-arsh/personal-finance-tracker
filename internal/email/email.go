package email

import (
	"errors"
	"net/mail"

	gomail "gopkg.in/mail.v2"
)

var ErrNoSender = errors.New("No sender email address provided")
var ErrNoRecevier = errors.New("No receiver email address provided")
var ErrNoTextBody = errors.New("No plain text body provided")

type email struct {
	receiver mail.Address
	subject  string
	htmlBody string
	textBody string
}

type GmailDialer struct {
	sender   mail.Address
	password string
}

const GMAIL_SMTP_URL = "smtp.gmail.com"
const GMAIL_SMTP_PORT = 587

func NewEmail(recevier, subject, htmlBody, textBody string) (*email, error) {
	emailAddress, err := mail.ParseAddress(recevier)
	if err != nil {
		return nil, err
	}

	return &email{*emailAddress, subject, htmlBody, textBody}, nil
}

func NewGmailDialer(sender, password string) (*GmailDialer, error) {
	emailAddress, err := mail.ParseAddress(sender)
	if err != nil {
		return nil, err
	}

	return &GmailDialer{*emailAddress, password}, nil
}

func (d *GmailDialer) Send(e *email) error {
	message := gomail.NewMessage()

	if d.sender.Address == "" {
		return ErrNoSender
	}
	message.SetHeader("From", d.sender.Address)

	if e.receiver.Address == "" {
		return ErrNoRecevier
	}
	message.SetHeader("To", e.receiver.Address)
	message.SetHeader("Subject", e.subject)

	if e.textBody == "" {
		return ErrNoTextBody
	}

	if e.htmlBody == "" {
		message.SetBody("text/plain", e.textBody)
	} else {
		message.AddAlternative("text/plain", e.textBody)
		message.SetBody("text/html", e.htmlBody)
	}

	dialer := gomail.NewDialer(GMAIL_SMTP_URL, GMAIL_SMTP_PORT, d.sender.Address, d.password)

	if err := dialer.DialAndSend(message); err != nil {
		return err
	}
	return nil
}
