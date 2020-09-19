package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/kiyonlin/dawn"
)

func main() {
	server := dawn.Default()

	r := server.Router()
	// GET /  =>  Welcome to dawn 👋
	r.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to dawn 👋")
	})

	log.Println(server.Run(":3000"))
}
