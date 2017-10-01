package utils

import (
	"strings"
	"net/smtp"
)

func SendMail(srvHost, userEmail, password, to, subject, body string) error {
	pos := strings.Index(srvHost, `:`)
	if pos < 0 {
		pos = len(srvHost)
	}
	auth := smtp.PlainAuth("", userEmail, password, srvHost[:pos])
	content_type := "Content-Type: text/html; charset=UTF-8"
	msg := make([]byte, 0, len(to)+len(userEmail)+len(subject)+len(body)+128)
	msg = append(msg, "To: "...)
	msg = append(msg, to...)
	msg = append(msg, "\r\nFrom: "...)
	msg = append(msg, userEmail...)
	msg = append(msg, "\r\nSubject: "...)
	msg = append(msg, subject...)
	msg = append(msg, "\r\n"...)
	msg = append(msg, content_type...)
	msg = append(msg, "\r\n\r\n"...)
	msg = append(msg, body...)
	msg = append(msg, "\r\n"...)

	send_to := strings.Split(to, ";")
	err := smtp.SendMail(srvHost, auth, userEmail, send_to, msg)
	return err
}
