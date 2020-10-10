package main

import (
	"fmt"
	"io"
	"os"

	"github.com/kiyonlin/dawn/cmd/dawn/internal"
	"github.com/urfave/cli/v2"
)

const version = "v0.0.5"

func init() {
	cli.AppHelpTemplate = AppHelpTemplate
	cli.CommandHelpTemplate = CommandHelpTemplate
	cli.SubcommandHelpTemplate = SubcommandHelpTemplate
}

func main() {
	run(os.Args, os.Stdout, os.Stderr)
}

func run(args []string, w io.Writer, ew io.Writer) {
	app := &cli.App{
		Version: version,
		Commands: []*cli.Command{
			internal.NewProject, internal.Module, internal.Dev,
		},
		Writer:    w,
		ErrWriter: ew,
	}

	if err := app.Run(args); err != nil {
		fmt.Println(err)
	}
}
