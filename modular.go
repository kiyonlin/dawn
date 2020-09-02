package dawn

import (
	"fmt"
	"sync"

	"github.com/gofiber/fiber"
)

// Cleanup is a function does cleanup works
type Cleanup func()

type Modular interface {
	// Stringer indicates Module's name
	fmt.Stringer

	// Init does initialization works.
	// Should call `WaitGroup.Done()` once done.
	Init(*sync.WaitGroup)

	// Boot boots the module
	// and pass a cleanup function
	// Should call `WaitGroup.Done()` once done.
	Boot(*sync.WaitGroup, chan<- Cleanup)

	// RegisterRoutes add routes to fiber router
	RegisterRoutes(fiber.Router)
}

// Module is an empty struct implements Modular interface
// and can be embedded into custom struct as a Modular
type Module struct{}

func (Module) String() string          { return "anonymous" }
func (Module) Init(wg *sync.WaitGroup) { wg.Done() }
func (Module) Boot(wg *sync.WaitGroup, cleanup chan<- Cleanup) {
	defer wg.Done()
	cleanup <- func() {}
}
func (Module) RegisterRoutes(fiber.Router) {}
