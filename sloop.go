package dawn

import (
	"crypto/tls"
	"fmt"
	"sync"

	"github.com/kiyonlin/dawn/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/kiyonlin/dawn/fiberx"
)

// Sloop denotes Dawn web server
type Sloop struct {
	app      *fiber.App
	wg       sync.WaitGroup
	mods     []Modular
	cleanups []Cleanup
}

// New returns a new blank Sloop.
func New(opts ...Option) *Sloop {
	s := &Sloop{
		app: fiber.New(),
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// Default returns an Sloop instance with the
// - RequestID
// - Logger
// - Recovery
// - Pprof
// middleware already attached in default fiber app.
func Default(cfg ...fiber.Config) *Sloop {
	c := fiber.Config{}
	if len(cfg) > 0 {
		c = cfg[0]
	}
	if c.ErrorHandler == nil {
		c.ErrorHandler = fiberx.ErrHandler
	}
	app := fiber.New(c)
	app.Use(
		requestid.New(),
		fiberx.Logger(),
		recover.New(),
	)

	if config.GetBool("debug") {
		app.Use(pprof.New())
	}

	return &Sloop{
		app: app,
	}
}

// AddModulars appends more Modulars
func (s *Sloop) AddModulars(m ...Modular) {
	s.mods = append(s.mods, m...)
}

// Run runs a web server
func (s *Sloop) Run(addr string) error {
	defer s.cleanup()
	return s.setup().app.Listen(addr)
}

// Run runs a tls web server
func (s *Sloop) RunTls(addr, certFile, keyFile string) error {
	defer s.cleanup()

	// Create tls certificate
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return err
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}

	ln, err := tls.Listen("tcp", addr, tlsConfig)
	if err != nil {
		return err
	}

	return s.setup().app.Listener(ln)
}

// Shutdown gracefully shuts down the server without interrupting any active connections.
func (s *Sloop) Shutdown() error {
	if s.app == nil {
		return fmt.Errorf("shutdown: fiber app is not found")
	}
	return s.app.Shutdown()
}

// Router returns the server router
func (s *Sloop) Router() fiber.Router {
	return s.app
}

func (s *Sloop) setup() *Sloop {
	return s.init().boot().registerRoutes()
}

func (s *Sloop) init() *Sloop {
	for _, mod := range s.mods {
		s.wg.Add(1)
		go mod.Init(&s.wg)
	}
	s.wg.Wait()
	return s
}

func (s *Sloop) boot() *Sloop {
	cleanups := make(chan Cleanup, len(s.mods))

	for _, mod := range s.mods {
		s.wg.Add(1)
		go mod.Boot(&s.wg, cleanups)
	}

	s.wg.Wait()

	close(cleanups)

	for cleanup := range cleanups {
		s.cleanups = append(s.cleanups, cleanup)
	}

	return s
}

func (s *Sloop) registerRoutes() *Sloop {
	for _, mod := range s.mods {
		mod.RegisterRoutes(s.app)
	}
	return s
}

func (s *Sloop) cleanup() {
	for _, fn := range s.cleanups {
		fn()
	}
}
