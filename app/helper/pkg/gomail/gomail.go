package gomail

import (
	"github.com/go-gomail/gomail"
	"mybedv2/app/system/model/mail"
	"strconv"
)

type Config struct {
	MailTo  []string
	Cc      string
	Subject string
	Config  mail.BindForm
}

func SendMailV1(conf Config) error {
	port, _ := strconv.Atoi(conf.Config.EmailPort)
	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(conf.Config.EmailUser, conf.Config.EmailName))
	m.SetHeader("To", conf.MailTo...)
	//m.SetAddressHeader("Cc", conf.Cc, conf.Config.EmailUser)
	m.SetHeader("Subject", conf.Subject)
	m.SetBody("text/html", conf.Config.EmailTemplate)
	d := gomail.NewDialer(conf.Config.EmailHost, port, conf.Config.EmailUser, conf.Config.EmailPwd)
	err := d.DialAndSend(m)
	return err
}
