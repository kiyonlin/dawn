package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
)

var module = &cli.Command{
	Name:      "module",
	Aliases:   []string{"m"},
	Usage:     "Generate a new dawn module",
	UsageText: "module name",
	Action: func(c *cli.Context) error {
		if !c.Args().Present() {
			return exit(c, "Missing module name")
		}
		now := time.Now()

		name := c.Args().First()

		dir, err := os.Getwd()
		if err != nil {
			return exit(c, err)
		}

		modulePath := dir + "/" + name
		if err = createModule(modulePath, name); err != nil {
			return exit(c, err)
		}

		return success(fmt.Sprintf(moduleCreatedTemplate, modulePath, name, time.Since(now)))
	},
}

func createModule(modulePath, name string) (err error) {
	if err = os.Mkdir(modulePath, 0750); err != nil {
		return
	}

	defer func() {
		if err != nil {
			if e := os.RemoveAll(modulePath); e != nil {
				err = e
			}
		}
	}()

	if err = os.Chdir(modulePath); err != nil {
		return
	}

	// create module.go
	var moduleFile *os.File
	if moduleFile, err = os.Create(fmt.Sprintf("%s/%s.go", modulePath, name)); err != nil {
		return
	}
	defer func() {
		if e := moduleFile.Close(); e != nil {
			err = e
		}
	}()

	if _, err = moduleFile.WriteString(moduleContent(name)); err != nil {
		return
	}

	// create module_test.go
	var moduleTestFile *os.File
	if moduleTestFile, err = os.Create(fmt.Sprintf("%s/%s_test.go", modulePath, name)); err != nil {
		return
	}
	defer func() {
		if e := moduleTestFile.Close(); e != nil {
			err = e
		}
	}()

	if _, err = moduleTestFile.WriteString(moduleTestContent(name)); err != nil {
		return
	}

	return
}

func moduleContent(name string) string {
	temp := `package {{module}}

import (
	"sync"

	"github.com/kiyonlin/dawn"
	"github.com/gofiber/fiber/v2"
)

type {{module}}Module struct {
	dawn.Module
}

// New returns the module
func New() dawn.Modular {
	return &{{module}}Module{
	}
}

func (m *{{module}}Module) String() string {
	return "{{module}}"
}

func (m *{{module}}Module) Init(wg *sync.WaitGroup) {
	// you can implement me or remove me
	defer wg.Done()
	
	// Read config and init module
}

func (m *{{module}}Module) Boot(wg *sync.WaitGroup, cleanup chan<- dawn.Cleanup) {
	// you can implement me or remove me
	defer wg.Done()

	// Do some tasks here

	cleanup <- func() {
		// Put cleanup stuff here if any
	}
}

func (m *{{module}}Module) RegisterRoutes(router fiber.Router) {
	// implement me or remove me
}`
	return strings.ReplaceAll(temp, "{{module}}", name)
}

func moduleTestContent(name string) string {
	temp := `package {{module}}

import (
	"sync"
	"testing"

	"github.com/kiyonlin/dawn"
	"github.com/stretchr/testify/assert"
)

func Test_Module_Name(t *testing.T) {
	assert.Equal(t, "{{module}}", New().String())
}

func Test_Module_Init(t *testing.T) {
	m := &{{module}}Module{}

	var wg sync.WaitGroup
	wg.Add(1)
	m.Init(&wg)

	wg.Wait()

	// more assertions
}

func Test_Module_Boot(t *testing.T) {
	m := &{{module}}Module{}

	var wg sync.WaitGroup
	wg.Add(1)
	cleanup := make(chan dawn.Cleanup, 1)

	m.Boot(&wg, cleanup)

	wg.Wait()

	assert.Len(t, cleanup, 1)
	// more assertions
}`
	return strings.ReplaceAll(temp, "{{module}}", name)
}

var (
	moduleCreatedTemplate = `
Scaffolding module in %s

  Done. Now run:

  cd %s
  go test . -cover

🎊  Done in %s.
`
)