package main

import (
	"fmt"
	"gopkg.in/gomail.v2"
)

type SmtpConf struct {
	Host string
	Port int
	User string
	Pass string
}

type EmailNotifier struct {
	from string
	to   string
	conf SmtpConf
}

func NewEmailNotifier(from string, to string, conf SmtpConf) *EmailNotifier {
	return &EmailNotifier{from: from, to: to, conf: conf}
}

const notificationTemplate = `Hey webhook admin,

i failed to execute one of your hooks.

command: %s
error:   %s

You can find the command output below.

---
%s
`

func (en *EmailNotifier) Notify(n *Notification) {
	msg := gomail.NewMessage()
	msg.SetHeader("From", en.from)
	msg.SetHeader("To", en.to)
	msg.SetHeader("Subject", fmt.Sprintf("[webhook] Failed to execute hook: %s", n.Command))
	msg.SetBody("text/plain", fmt.Sprintf(notificationTemplate, n.Command, n.Error, n.Output))
	d := gomail.NewDialer(en.conf.Host, en.conf.Port, en.conf.User, en.conf.Pass)

	if err := d.DialAndSend(msg); err != nil {
		fmt.Println("failed to send notification email:", err)
	}
}
