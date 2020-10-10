package internal

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

func Test_Module_Run(t *testing.T) {
	at := assert.New(t)

	t.Run("success", func(t *testing.T) {
		defer func() {
			at.Nil(os.RemoveAll("testcase"))
		}()

		app := &cli.App{
			Commands: []*cli.Command{Module},
			ExitErrHandler: func(c *cli.Context, err error) {
				at.Contains(err.Error(), "Done")
			}}

		at.NotNil(app.Run([]string{"bin", "module", "testcase"}))
	})

	t.Run("missing module name", func(t *testing.T) {
		app := &cli.App{
			Writer:   &bytes.Buffer{},
			Commands: []*cli.Command{Module},
			ExitErrHandler: func(c *cli.Context, err error) {
				at.Contains(err.Error(), "Missing")
			}}

		at.NotNil(app.Run([]string{"bin", "module"}))
	})

	t.Run("invalid module name", func(t *testing.T) {
		app := &cli.App{
			Writer:   &bytes.Buffer{},
			Commands: []*cli.Command{Module},
			ExitErrHandler: func(c *cli.Context, err error) {
				at.Contains(err.Error(), ".")
			}}

		at.NotNil(app.Run([]string{"bin", "module", "."}))
	})
}

func Test_Module_CreateModule(t *testing.T) {
	t.Parallel()

	at := assert.New(t)

	dir, err := ioutil.TempDir("", "test_create_module")
	at.Nil(err)
	defer func() { _ = os.RemoveAll(dir) }()

	modulePath := fmt.Sprintf("%s%cmodule", dir, os.PathSeparator)

	at.NotNil(createModule(modulePath, "invalid-name/"))
}
