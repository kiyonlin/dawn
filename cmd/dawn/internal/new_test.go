package internal

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

func Test_New_Run(t *testing.T) {
	at := assert.New(t)

	t.Run("new project", func(t *testing.T) {
		defer func() {
			at.Nil(os.Chdir("../"))
			at.Nil(os.RemoveAll("testcase"))
		}()

		setupCmd()
		defer teardownCmd()

		app := &cli.App{
			Commands: []*cli.Command{NewProject},
			ExitErrHandler: func(c *cli.Context, err error) {
				at.Contains(err.Error(), "Done")
			}}

		at.NotNil(app.Run([]string{"bin", "new", "testcase"}))
	})

	t.Run("custom mod name", func(t *testing.T) {
		defer func() {
			at.Nil(os.Chdir("../"))
			at.Nil(os.RemoveAll("testcase"))
		}()

		setupCmd()
		defer teardownCmd()

		app := &cli.App{
			Commands: []*cli.Command{NewProject},
			ExitErrHandler: func(c *cli.Context, err error) {
				at.Contains(err.Error(), "custom")
			}}

		at.NotNil(app.Run([]string{"bin", "new", "testcase", "custom"}))
	})

	t.Run("use --app", func(t *testing.T) {
		defer func() {
			at.Nil(os.Chdir("../"))
			at.Nil(os.RemoveAll("testcase"))
		}()

		setupCmd(errFlag)
		defer teardownCmd()

		app := &cli.App{
			Commands: []*cli.Command{NewProject},
			ExitErrHandler: func(c *cli.Context, err error) {
				at.Contains(err.Error(), "failed to run")
			}}

		at.NotNil(app.Run([]string{"bin", "new", "--app", "testcase"}))
	})

	t.Run("missing project name", func(t *testing.T) {
		app := &cli.App{
			Commands: []*cli.Command{NewProject},
			ExitErrHandler: func(c *cli.Context, err error) {
				at.Contains(err.Error(), "Missing")
			}}

		at.NotNil(app.Run([]string{"bin", "new"}))
	})

	t.Run("invalid project name", func(t *testing.T) {
		app := &cli.App{
			Commands: []*cli.Command{NewProject},
			ExitErrHandler: func(c *cli.Context, err error) {
				at.Contains(err.Error(), ".")
			}}

		at.NotNil(app.Run([]string{"bin", "new", "."}))
	})
}
