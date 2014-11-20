package out

import (
	"fmt"
	"strings"

	"github.com/dtynn/dout/cache"
	"github.com/dtynn/dout/mail"
	"github.com/dtynn/dout/mx"
	"github.com/dtynn/dout/smtp"
)

var cacheExpire int64 = 60 * 10

type D struct {
	l     Logger
	cache cache.Cache
}

func NewDout(l Logger, c cache.Cache) *D {
	return &D{
		l:     l,
		cache: c,
	}
}

func (this *D) Send(m *mail.Mail) ([]*Status, error) {
	pieces := strings.Split(m.From, "@")
	if len(pieces) != 2 || len(pieces[1]) == 0 {
		return nil, InvalidEmail(m.From)
	}

	local := pieces[1]

	if len(m.Tos) == 0 {
		return nil, noRecipients
	}

	msg := m.Bytes()
	stats := []*Status{}
	for _, to := range m.Tos {
		// check email
		pieces := strings.Split(to, "@")
		if len(pieces) != 2 || len(pieces[1]) == 0 {
			stats = append(stats, NewStatus(StatusInvalidToAddress, to, "invalid email"))
			continue
		}

		// dns mx
		name := pieces[1]
		host, err := this.cache.Get(name)
		if err == cache.CacheNotFound {
			mxRec, err := mx.ChoseMx(name)
			if err != nil {
				stats = append(stats, NewStatus(StatusDnsErr, to, err.Error()))
				continue
			}
			host = mxRec.Host
		}

		// send
		addr := fmt.Sprintf("%s:%d", host, smtp.DefaultPort)
		err = smtp.SendEmail(addr, local, m.From, []string{to}, msg, nil)
		if err != nil {
			stats = append(stats, NewStatus(StatusSmtpErr, to, err.Error()))
			continue
		}
		this.cache.Setex(name, host, cacheExpire)
	}

	if len(stats) == 0 {
		return stats, nil
	}
	return stats, someFailed
}
