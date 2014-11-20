package main

import (
	"github.com/dtynn/dout/web"
	"github.com/qiniu/log"
)

func main() {
	srv := web.New()
	addr := ":10025"
	log.Infof("Dout Listen on %s", addr)
	err := srv.Run(addr)
	if err != nil {
		log.Fatal("Err: ", err)
	}
}
