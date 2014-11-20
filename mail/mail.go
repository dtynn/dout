package mail

import (
	"bytes"
)

var lineSep = "\r\n"

type Mail struct {
	From string
	Tos  []string
	Body string
}

func NewMail(from string, tos []string, body string) *Mail {
	return &Mail{
		From: from,
		Tos:  tos,
		Body: body,
	}
}

func (this *Mail) Bytes() []byte {
	var w mailWriter
	w.writeLine(this.From)
	for _, to := range this.Tos {
		w.writeLine(to)
	}
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
