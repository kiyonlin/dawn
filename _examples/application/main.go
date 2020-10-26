package main

import (
	"flag"

	"github.com/kiyonlin/dawn"
	"github.com/kiyonlin/dawn/config"
	"github.com/kiyonlin/dawn/db/redis"
	"github.com/kiyonlin/dawn/db/sql"
	"github.com/kiyonlin/dawn/log"
)

func main() {
	config.Load("./_examples/application")
	config.LoadEnv()

	log.InitFlags(nil)
	flag.Parse()
	defer log.Flush()

	sloop := dawn.New().AddModulers(
		sql.New(),
		redis.New(),
		// add custom module
	)

	defer sloop.Cleanup()

	sloop.Setup().Watch()
}
