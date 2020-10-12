package internal

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
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

		out, err := runCobraCmd(NewCmd, "testcase")

		at.Nil(err)
		at.Contains(out, "Done")
	})

	t.Run("custom mod name", func(t *testing.T) {
		defer func() {
			at.Nil(os.Chdir("../"))
			at.Nil(os.RemoveAll("testcase"))
		}()

		setupCmd()
		defer teardownCmd()

		out, err := runCobraCmd(NewCmd, "testcase", "custom")

		at.Nil(err)
		at.Contains(out, "custom")
	})

	t.Run("use --app and fail", func(t *testing.T) {
		defer func() {
			at.Nil(os.Chdir("../"))
			at.Nil(os.RemoveAll("testcase"))
		}()

		setupCmd(errFlag)
		defer teardownCmd()

		out, err := runCobraCmd(NewCmd, "testcase", "--app")

		at.NotNil(err)
		at.Contains(out, "failed to run")
	})

	t.Run("invalid project name", func(t *testing.T) {
		out, err := runCobraCmd(NewCmd, ".")

		at.NotNil(err)
		at.Contains(out, ".")
	})
}
