package mail

import (
	"bytes"
	"crypto/tls"
	"html/template"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/go-gomail/gomail"
	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("test")

// Config .
type Config struct {
	ServerHost string // ServerHost, eg: smtp.gmail.com
	ServerPort int    // ServerPort eg: 465
	FromEmail  string // FromEmail Sender email address
	FromPasswd string // FromPasswd Sender email pwd
	FromUser   string // FromEmail --> Name
	Toers      string // Toers Send to? if more than one, split by ";"
	CCers      string // CCers Copy to? if more than one, split by ";"
}

// ParseHTML .
func ParseHTML(fromEmail, fromUser, toUsers, content string, body *bytes.Buffer) error {
	_, file, _, _ := runtime.Caller(0)
	t, err := template.ParseFiles(path.Join(path.Dir(file), "mail-template-default.html"))
	if err != nil {
		return err
	}

	t.Execute(body, struct {
		FromUser  string
		FromEmail string
		ToUsers   string
		TimeDate  string
		Content   string
	}{
		FromUser:  fromUser,
		FromEmail: fromEmail,
		ToUsers:   toUsers,
		TimeDate:  time.Now().Format("2006/01/02 15:04:05"),
		Content:   content, //strings.ReplaceAll(content, "\\n", "%0d%0a"),
	})
	return nil
}

// SendEmail body/html
func SendEmail(cf *Config, subject, content string, attachments ...string) {
	logger.Infof("> Send Email To:%v; Cc:%v", cf.Toers, cf.CCers)
	if len(cf.Toers) == 0 {
		logger.Fatal("Email Send \"To\" no one")
	}

	msg := gomail.NewMessage()

	// SetHeader To
	toers := []string{}
	for _, tmpT := range strings.Split(cf.Toers, ";") {
		toers = append(toers, strings.TrimSpace(tmpT))
	}
	msg.SetHeader("To", toers...)

	// SetHeader Cc
	if len(cf.CCers) != 0 {
		ccers := []string{}
		for _, tmpC := range strings.Split(cf.CCers, ";") {
			ccers = append(ccers, strings.TrimSpace(tmpC))
		}
		msg.SetHeader("Cc", ccers...)
	}

	// SetHeader From, Subject
	msg.SetAddressHeader("From", cf.FromEmail, cf.FromUser)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", content)

	// Set Attachments
	for _, att := range attachments {
		msg.Attach(att)
	}

	d := gomail.NewPlainDialer(cf.ServerHost, cf.ServerPort, cf.FromEmail, cf.FromPasswd)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	err := d.DialAndSend(msg)
	if err != nil {
		logger.Fatal(err)
	}
}
