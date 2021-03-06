package mail

import (
	"bytes"
	"fmt"
)

var lineSep = "\r\n"

type Mail struct {
	From    string
	Tos     []string
	Subject string
	Body    string
}

func NewMail(from, subject string, tos []string, body string) *Mail {
	return &Mail{
		From:    from,
		Tos:     tos,
		Subject: subject,
		Body:    body,
	}
}

func (this *Mail) Bytes() []byte {
	var w mailWriter
	w.writeLine(fmt.Sprintf("From: %s", this.From))
	for _, to := range this.Tos {
		w.writeLine(fmt.Sprintf("To: %s", to))
	}
	w.writeLine(fmt.Sprintf("Subject: %s", this.Subject))
	w.writeLine("MIME-Version: 1.0")
	w.writeLine("Content-Type: text/html; charset=UTF-8")
	w.writeLine("")
	w.writeLine(this.Body)
	return w.Bytes()
}

type mailWriter struct {
	bytes.Buffer
}

func (this *mailWriter) writeLine(line string) {
	this.WriteString(line + lineSep)
}
