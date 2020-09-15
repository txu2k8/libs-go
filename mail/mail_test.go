package mail

import (
	"bytes"
	"testing"
)

func TestEmail(t *testing.T) {
	cf := &Config{
		ServerHost: "smtp.gmail.com",
		ServerPort: 465,
		FromEmail:  "txu@xxx.com",
		FromPasswd: "Pass@word",
		Toers:      "txu1@xxx.com",
		CCers:      "",
	}

	subject := "Test"
	content := "ssssssss"
	// Set Body with mail-template.html
	body := new(bytes.Buffer)
	ParseHTML(cf.FromEmail, "TestUser", cf.Toers, content, body)
	SendEmail(cf, subject, body.String())
}
