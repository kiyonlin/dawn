package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/kiyonlin/dawn"
	"github.com/kiyonlin/dawn/fiberx"
)

func main() {
	app := fiber.New(fiber.Config{
		Prefork: true,
	})

	// GET /  =>  I'm in prefork mode 🚀
	app.Get("/", func(c *fiber.Ctx) error {
		return fiberx.Message(c, "I'm in prefork mode 🚀")
	})

	server := dawn.New(dawn.App(app))

	log.Println(server.Run(":3000"))
}
