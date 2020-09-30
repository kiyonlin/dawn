package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

const version = "v0.0.2"

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
			newProject, module,
		},
	}

	if err := app.Run(args); err != nil {
		fmt.Println(err)
	}
}
