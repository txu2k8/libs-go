package config

import "libs-go/mail"

// EmailConfig .
var (
	MailToSelf     = "txu@xxx.com"
	MailToDev      = "txu@xxx.com"
	MailToStressQA = "txu@xxx.com;txu1@xxx.com;"
	EmailConfig    = mail.Config{
		ServerHost: "smtp.gmail.com",
		ServerPort: 465,
		FromEmail:  "txu@xxx.com",
		FromPasswd: "Pass@word",
		FromUser:   "StressTest",
		Toers:      MailToSelf,
		CCers:      MailToStressQA,
	}
)
