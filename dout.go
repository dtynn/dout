package dout

import (
	"crypto/tls"
	"errors"
	"fmt"
	"strings"

	"github.com/dtynn/dout/cache/memory"
	"github.com/dtynn/dout/mx"
	"github.com/dtynn/dout/smtp"
)

var (
	errInvalidEmail = errors.New("invalid email address")
)

var cacheExpire int64 = 60 * 10

type D struct {
	l        Logger
	hostname string
	from     string
	cache    Cache
}

type sendConfig struct {
	domain string
	from   string
	tos    []string
	msg    []byte
	tlsCfg *tls.Config
}

func New(hostname, from string, l Logger, cache Cache) (*D, error) {
	if pieces := strings.Split(from, "@"); len(pieces) != 2 {
		return nil, errInvalidEmail
	}

	return &D{
		l:        l,
		hostname: hostname,
		from:     from,
		cache:    cache,
	}, nil
}

func NewWithMemory(hostname, from string, l Logger) (*D, error) {
	cache := memory.New()
	return New(hostname, from, l, cache)
}

func (this *D) getMX(name string) (string, error) {
	if val, err := this.cache.Get(name); err == nil {
		return val.(string), nil
	}
	rec, err := mx.ChoseMx(name)
	if err != nil {
		return "", err
	}
	if err := this.cache.Setex(name, rec.Host, cacheExpire); err != nil {
		this.l.Warnf("set cache for %s failed: %s", name, err)
	}
	return rec.Host, nil
}

func (this *D) clearCache(name string) error {
	return this.cache.Del(name)
}

func (this *D) send(cfg *sendConfig) error {
	host, err := this.getMX(cfg.domain)
	if err != nil {
		this.l.Warnf("failed to resolve %s", cfg.domain)
		return err
	}
	addr := fmt.Sprintf("%s:%d", host, smtp.DefaultPort)
	err = smtp.SendEmail(addr, this.hostname, cfg.from, cfg.tos, cfg.msg, cfg.tlsCfg)
	if err != nil {
		this.clearCache(cfg.domain)
	}
	return err
}

func (this *D) SendOne(from, to string, msg []byte) error {
	pieces := strings.Split(to, "@")
	if len(pieces) != 2 {
		this.l.Warnf("Invalid email: %s", to)
		return errInvalidEmail
	}

	domain := pieces[1]

	cfg := &sendConfig{
		domain: domain,
		from:   from,
		tos:    []string{to},
		msg:    msg,
		tlsCfg: nil,
	}
	return this.send(cfg)
}

func (this *D) SendOneWithDefault(to string, msg []byte) error {
	return this.SendOne(this.from, to, msg)
}

func (this *D) SendMulti(from string, tos []string, msg []byte) {
	tosMap := map[string][]string{}
	for _, to := range tos {
		pieces := strings.Split(to, "@")
		if len(pieces) != 2 {
			this.l.Warnf("Abandoned email address: %s", to)
			continue
		}

		domain := pieces[1]
		slice, has := tosMap[domain]
		if has {
			tosMap[domain] = append(slice, to)
		} else {
			tosMap[domain] = []string{to}
		}
	}

	for domain, domainTos := range tosMap {
		cfg := &sendConfig{
			domain: domain,
			from:   from,
			tos:    domainTos,
			msg:    msg,
		}
		go func(cfg *sendConfig) {
			err := this.send(cfg)
			if err != nil {
				this.l.Errorf("failed to send: ", err)
			}
		}(cfg)
	}
	return
}

func (this *D) SendMultiWithDefault(tos []string, msg []byte) {
	this.SendMulti(this.from, tos, msg)
}
