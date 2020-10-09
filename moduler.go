package dawn

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// Cleanup is a function does cleanup works
type Cleanup func()

type Moduler interface {
	// Stringer indicates Module's name
	fmt.Stringer

	// Init does initialization works and should return
	// a cleanup function.
	Init() Cleanup

	// Boot boots the module.
	Boot()

	// RegisterRoutes add routes to fiber router
	RegisterRoutes(fiber.Router)
}

// Module is an empty struct implements Moduler interface
// and can be embedded into custom struct as a Moduler
type Module struct{}

func (Module) String() string              { return "anonymous" }
func (Module) Init() Cleanup               { return nil }
func (Module) Boot()                       {}
func (Module) RegisterRoutes(fiber.Router) {}
