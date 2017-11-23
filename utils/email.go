package utils

import (
	"strings"
	"net/smtp"
	"encoding/base64"
)

type MailBody struct {
	content []byte
}

func NewMailBody(ln int) (mb *MailBody) {
	mb = &MailBody{make([]byte, 0, ln)}
	return
}
func (this *MailBody) AddHeader(key, val string) {
	this.content = append(this.content, key...)
	this.content = append(this.content, `: `...)
	this.content = append(this.content, val...)
	this.content = append(this.content, "\r\n"...)
}
func (this *MailBody) AddBody(msg []byte) {
	this.content = append(this.content, "\r\n"...)
	this.content = append(this.content, msg...)
	this.content = append(this.content, "\r\n"...)
}
func (this *MailBody) Bytes() (bs []byte) {
	bs = this.content
	return
}

func SendMail(smtpSrv, user, password, from, to, subject, body string) (err error) {
	host := smtpSrv
	pos := strings.Index(smtpSrv, `:`)
	if pos >= 0 {
		host = host[:pos]
	}
	auth := smtp.PlainAuth("", user, password, host)
	content_type := "Content-Type: text/html; charset=UTF-8"
	msg := NewMailBody(len(to) + len(user) + len(subject) + len(body) + 128)
	msg.AddHeader(`From`, from)
	msg.AddHeader(`To`, to)
	msg.AddHeader(`Subject`, subject)
	msg.AddHeader("MIME-Version", "1.0")
	msg.AddHeader(`Content-Type`, content_type)
	msg.AddHeader("Content-Transfer-Encoding", "base64")
	enc := base64.StdEncoding
	buf := make([]byte, enc.EncodedLen(len(body)))
	enc.Encode(buf, []byte(body))
	msg.AddBody(buf)

	send_to := strings.Split(to, ";")
	err = smtp.SendMail(smtpSrv, auth, from, send_to, msg.Bytes())
	return
}
