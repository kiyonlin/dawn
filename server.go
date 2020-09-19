package dawn

import (
	"crypto/tls"
	"fmt"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/kiyonlin/dawn/config"
	"github.com/kiyonlin/dawn/fiberx"
)

// Server denotes Dawn web server
type Server struct {
	app      *fiber.App
	wg       sync.WaitGroup
	mods     []Modular
	cleanups []Cleanup
}

// New returns a new blank Server.
func New(opts ...Option) *Server {
	s := &Server{
		app: fiber.New(),
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// Default returns an Server instance with the
// - RequestID
// - Logger
// - Recovery
// - Pprof
// middleware already attached in default fiber app.
func Default() *Server {
	app := fiber.New(fiber.Config{
		ErrorHandler: fiberx.ErrHandler,
	})
	app.Use(
		requestid.New(),
		fiberx.Logger(),
		recover.New(),
	)

	if config.GetBool("debug") {
		app.Use(pprof.New())
	}

	return &Server{
		app: app,
	}
}

// AddModulars appends more Modulars
func (s *Server) AddModulars(m ...Modular) {
	s.mods = append(s.mods, m...)
}

// Run runs a web server
func (s *Server) Run(addr string) error {
	return s.setup().app.Listen(addr)
}

// Run runs a tls web server
func (s *Server) RunTls(addr, certFile, keyFile string) error {
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
func (s *Server) Shutdown() error {
	if s.app == nil {
		return fmt.Errorf("shutdown: fiber app is not found")
	}
	return s.app.Shutdown()
}

// Router returns the server router
func (s *Server) Router() fiber.Router {
	return s.app
}

func (s *Server) setup() *Server {
	return s.init().boot().registerRoutes()
}

func (s *Server) init() *Server {
	for _, mod := range s.mods {
		s.wg.Add(1)
		go mod.Init(&s.wg)
	}
	s.wg.Wait()
	return s
}

func (s *Server) boot() *Server {
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

func (s *Server) registerRoutes() *Server {
	for _, mod := range s.mods {
		mod.RegisterRoutes(s.app)
	}
	return s
}
