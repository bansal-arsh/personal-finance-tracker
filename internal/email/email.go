package email

import (
	"errors"
	"net/mail"

	gomail "gopkg.in/mail.v2"
)

var ErrNoSender = errors.New("No sender email address provided")
var ErrNoGmailAppPassword = errors.New("No gmail app password for sender email address provided")
var ErrNoRecevier = errors.New("No receiver email address provided")
var ErrNoTextBody = errors.New("No plain text body provided")

// fields are not meant to be interacted with directly
// must use their constructors to create these
// still exporting types for use by handlers and other functions
type Email struct {
	receiver mail.Address
	subject  string
	htmlBody string // can be empty
	textBody string
}

type GmailDialer struct {
	sender      mail.Address
	appPassword string
}

const GMAIL_SMTP_URL = "smtp.gmail.com"
const GMAIL_SMTP_PORT = 587

func NewEmail(recevier, subject, htmlBody, textBody string) (*Email, error) {
	if recevier == "" {
		return nil, ErrNoRecevier
	}

	parsedEmailAddress, err := mail.ParseAddress(recevier)
	if err != nil {
		return nil, err
	}

	if textBody == "" {
		return nil, ErrNoTextBody
	}

	return &Email{*parsedEmailAddress, subject, htmlBody, textBody}, nil
}

func NewGmailDialer(sender, appPassword string) (*GmailDialer, error) {
	if sender == "" {
		return nil, ErrNoSender
	}

	parsedEmailAddress, err := mail.ParseAddress(sender)
	if err != nil {
		return nil, err
	}

	if appPassword == "" {
		return nil, ErrNoGmailAppPassword
	}

	return &GmailDialer{*parsedEmailAddress, appPassword}, nil
}

// function not unit tested because all validations are done in the constructors
// smtp email errors can't be unit tested
func (d *GmailDialer) Send(e *Email) error {
	message := gomail.NewMessage()
	message.SetHeader("From", d.sender.Address)
	message.SetHeader("To", e.receiver.Address)
	message.SetHeader("Subject", e.subject)

	if e.htmlBody == "" {
		message.SetBody("text/plain", e.textBody)
	} else {
		message.AddAlternative("text/plain", e.textBody)
		message.SetBody("text/html", e.htmlBody)
	}

	dialer := gomail.NewDialer(GMAIL_SMTP_URL, GMAIL_SMTP_PORT, d.sender.Address, d.appPassword)

	if err := dialer.DialAndSend(message); err != nil {
		return err
	}
	return nil
}
