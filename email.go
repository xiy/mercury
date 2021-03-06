package main

import (
	"log"
	"net/mail"

	"github.com/domodwyer/mailyak"
)

// Email represents a single email message
type Email struct {
	From     string   `json:"from"`
	To       string   `json:"to"`
	Cc       []string `json:"cc"`
	Bcc      []string `json:"bcc"`
	Subject  string   `json:"subject"`
	TextBody string   `json:"text_body"`
	HTMLBody string   `json:"html_body"`
}

// Send sends the instance of Email using the given instance of mailyak.MailYak.
// It expects the instance of mailyak.MailYak to have been set up previously
// with a valid hostname and implementer of smtp.Auth, such as smtp.PlainAuth.
func (e *Email) Send(yak *mailyak.MailYak) error {
	from, err := mail.ParseAddress(e.From)
	if err != nil {
		log.Fatal(err)
	}

	to, err := mail.ParseAddress(e.To)
	if err != nil {
		log.Fatal(err)
	}

	// Decompress email text/html content
	textBody, htmlBody := decompressBody(e)

	yak.To(to.Address)
	yak.From(from.Address)
	yak.Subject(e.Subject)
	yak.HTML().Set(htmlBody)
	yak.Plain().Set(textBody)

	if err := yak.Send(); err != nil {
		return err
	}

	logger.Printf("info: Sent email to %s", e.To)

	return nil
}

func decompressBody(e *Email) (string, string) {
	return Decompress(e.TextBody), Decompress(e.HTMLBody)
}
