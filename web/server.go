package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/dtynn/dout/cache/memory"
	"github.com/dtynn/dout/mail"
	"github.com/dtynn/dout/out"
	"github.com/qiniu/log"
)

type server struct {
	d *out.D
}

func New() *server {
	c := memory.New()
	d := out.NewDout(log.Std, c)
	return &server{d}
}

func (this *server) send(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		outputJson(http.StatusBadRequest, w, resp{err.Error(), nil})
		return
	}

	from := r.PostFormValue("from")
	if pieces := strings.Split(from, "@"); len(pieces) != 2 || len(pieces[1]) == 0 {
		outputJson(http.StatusBadRequest, w, resp{fmt.Sprintf("from address <%s> is invalid", from), nil})
		return
	}

	toStr := r.PostFormValue("to")
	if len(toStr) == 0 {
		outputJson(http.StatusBadRequest, w, resp{"no recipients", nil})
		return
	}
	tos := strings.Split(toStr, ",")

	subject := r.PostFormValue("subject")
	body := r.PostFormValue("body")

	m := mail.NewMail(from, subject, tos, body)
	failed, err := this.d.Send(m)

	data := resp{}
	data.Msg = "ok"
	data.Failed = failed
	code := 200

	if err != nil {
		data.Msg = err.Error()
		if out.SomeFailed(err) {
			code = 298
		} else {
			code = http.StatusBadRequest
		}
	}

	outputJson(code, w, data)
	return
}

func (this *server) Run(addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/send", this.send)
	return http.ListenAndServe(addr, mux)
}

type resp struct {
	Msg    string        `json:"msg"`
	Failed []*out.Status `json:"failed"`
}

func outputJson(code int, w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	encoder := json.NewEncoder(w)
	encoder.Encode(data)
}
