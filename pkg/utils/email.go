package utils

import "gopkg.in/gomail.v2"

func SendEmail(host string, port int, username, password, from, to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(host, port, username, password)

	return d.DialAndSend(m)
}