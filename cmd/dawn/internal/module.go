package internal

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
)

// Module command generates a new dawn module
var Module = &cli.Command{
	Name:      "module",
	Aliases:   []string{"m"},
	Usage:     "Generate a new dawn module",
	UsageText: "dawn module name",
	Action: func(c *cli.Context) error {
		if !c.Args().Present() {
			return exit(c, "Missing module name")
		}
		now := time.Now()

		name := c.Args().First()

		dir, _ := os.Getwd()

		modulePath := dir + "/" + name
		if err := createModule(modulePath, name); err != nil {
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
			_ = os.RemoveAll(modulePath)
		}
	}()

	// create module.go
	if err = createFile(fmt.Sprintf("%s/%s.go", modulePath, name),
		moduleContent(name)); err != nil {
		return
	}

	// create module_test.go
	return createFile(fmt.Sprintf("%s/%s_test.go", modulePath, name),
		moduleTestContent(name))
}

func moduleContent(name string) string {
	temp := `package {{module}}

import (
	"github.com/kiyonlin/dawn"
	"github.com/gofiber/fiber/v2"
)

type {{module}}Module struct {
	dawn.Module
}

// New returns the module
func New() dawn.Moduler {
	return &{{module}}Module{
	}
}

func (m *{{module}}Module) String() string {
	return "{{module}}"
}

func (m *{{module}}Module) Init() dawn.Cleanup {
	// you can implement me or remove me

	// Read config and init module

	return func() {
		// Put cleanup stuff here if any
	}
}

func (m *{{module}}Module) Boot() {
	// you can implement me or remove me
}

func (m *{{module}}Module) RegisterRoutes(router fiber.Router) {
	// implement me or remove me
}`
	return strings.ReplaceAll(temp, "{{module}}", name)
}

func moduleTestContent(name string) string {
	temp := `package {{module}}

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Module_Name(t *testing.T) {
	assert.Equal(t, "{{module}}", New().String())
}

func Test_Module_Init(t *testing.T) {
	m := &{{module}}Module{}

	m.Init()()

	// more assertions
}

func Test_Module_Boot(t *testing.T) {
	m := &{{module}}Module{}

	m.Boot()

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
