package main

import (
	"os"
	"testing"

	"github.com/urfave/cli/v2"
)

func Test_Main_Run(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		run([]string{"bin", "help"})
	})

	t.Run("panic", func(t *testing.T) {
		cli.OsExiter = func(code int) {}
		defer func() { cli.OsExiter = os.Exit }()
		run([]string{"bin", "non"})
	})
}
