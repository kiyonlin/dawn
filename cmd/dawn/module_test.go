package main

import (
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
			Commands: []*cli.Command{module},
			ExitErrHandler: func(c *cli.Context, err error) {
				at.Contains(err.Error(), "Done")
			}}

		at.NotNil(app.Run([]string{"bin", "module", "testcase"}))
	})

	t.Run("missing module name", func(t *testing.T) {
		app := &cli.App{
			Commands: []*cli.Command{module},
			ExitErrHandler: func(c *cli.Context, err error) {
				at.Contains(err.Error(), "Missing")
			}}

		at.NotNil(app.Run([]string{"bin", "module"}))
	})

	t.Run("invalid module name", func(t *testing.T) {
		app := &cli.App{
			Commands: []*cli.Command{module},
			ExitErrHandler: func(c *cli.Context, err error) {
				at.Contains(err.Error(), ".")
			}}

		at.NotNil(app.Run([]string{"bin", "module", "."}))
	})
}
