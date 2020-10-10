package main

import (
	"fmt"
	"os"

	"github.com/kiyonlin/dawn/cmd/dawn/internal"
	"github.com/urfave/cli/v2"
)

const version = "v0.0.3"

func init() {
	cli.AppHelpTemplate = AppHelpTemplate
	cli.CommandHelpTemplate = CommandHelpTemplate
	cli.SubcommandHelpTemplate = SubcommandHelpTemplate
}

func main() {
	run(os.Args)
}

func run(args []string) {
	app := &cli.App{
		Version: version,
		Commands: []*cli.Command{
			internal.NewProject, internal.Module, internal.Dev,
		},
	}

	if err := app.Run(args); err != nil {
		fmt.Println(err)
	}
}
