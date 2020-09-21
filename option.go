package dawn

import (
	"github.com/gofiber/fiber/v2"
)

// Option can be applied in server
type Option func(s *Sloop)

// App option sets custom Fiber App to Sloop
func App(app *fiber.App) Option {
	return func(s *Sloop) {
		s.app = app
	}
}

// Modulars option adds several Modulars to server
func Modulars(mods ...Modular) Option {
	return func(s *Sloop) {
		s.mods = append(s.mods, mods...)
	}
}
