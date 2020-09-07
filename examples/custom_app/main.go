package main

import (
	"log"

	"github.com/gofiber/fiber"
	"github.com/kiyonlin/dawn"
)

func main() {
	app := fiber.New(&fiber.Settings{
		Prefork: true,
	})

	// GET /  =>  I'm in prefork mode 🚀
	app.Get("/", func(c *fiber.Ctx) {
		c.SendString("I'm in prefork mode 🚀")
	})

	server := dawn.New(dawn.App(app))

	log.Println(server.Run(":3000"))
}
