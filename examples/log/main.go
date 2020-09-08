package main

import (
	"flag"

	"github.com/kiyonlin/dawn/log"
)

func main() {
	log.InitFlags()
	flag.Parse()
	defer log.Flush()

	log.Errorln("error")
	log.Infoln(0, "info 0")
	log.Infoln(1, "info 1")
	// Won't log if set -v=1
	log.Infoln(2, "info 2")
}
