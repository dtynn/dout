#### dout
=====

attempt to directly send emails, instead of deploying a smtp server.  
download:
```
go get -u github.com/dtynn/dout
```

example: 

```golang
package main

import (
	"strings"

	"github.com/dtynn/dout"
	"github.com/qiniu/log"
)

func main() {
	tos := []string{"a@a.com", "a@b.com", "b@a.com", "c@c.com"}
	lines := []string{
		"From: noreply@example.com",
		"To: a@a.com",
		"To: a@b.com",
		"To: b@a.com",
		"To: c@c.com",
		"Subject: for example",
		"\r\n",
		"This is body",
	}
	body := strings.Join(lines, "\r\n")

	hostname := "example.com"
	from := "noreply@example.com"
	d, err := dout.NewWithMemory(hostname, from, log.Std)
	if err != nil {
		log.Fatal(err)
	}

	d.SendMultiWithDefault(tos, []byte(body))
}
```