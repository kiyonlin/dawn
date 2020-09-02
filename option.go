package dawn

import (
	"github.com/gofiber/fiber"
)

// Option can be applied in server
type Option func(s *Server)

// App option sets custom Fiber App to Server
func App(app *fiber.App) Option {
	return func(s *Server) {
		s.app = app
	}
}

// Modulars option adds several Modulars to server
func Modulars(mods ...Modular) Option {
	return func(s *Server) {
		s.mods = append(s.mods, mods...)
	}
}
