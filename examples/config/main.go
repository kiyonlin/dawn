package main

import (
	"log"

	"github.com/kiyonlin/dawn/config"
)

func init() {
	config.Load("./examples/config")
	config.LoadEnv("dawn")
}

func main() {
	// output: bar
	log.Println(config.GetString("foo"))

	// output: baz
	log.Println(config.GetString("bar", "baz"))

	// DAWN_FROM_ENV=hello go run ./examples/config
	// output: hello
	log.Println(config.GetString("from.env"))
}
