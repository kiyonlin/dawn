package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

const version = "v0.0.1"

func init() {
	cli.AppHelpTemplate = AppHelpTemplate
	cli.CommandHelpTemplate = CommandHelpTemplate
	cli.SubcommandHelpTemplate = SubcommandHelpTemplate

	cli.VersionFlag = &cli.BoolFlag{
		Name: "version", Aliases: []string{"v"},
		Usage: "print dawn version",
	}

	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("dawn %s(latest %s)\n", currentVersion(), latestVersion())
	}
}

func main() {
	app := &cli.App{
		Version: version,
		Commands: []*cli.Command{
			newProject,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func latestVersion() string {
	return "v0.0.1"
}

func currentVersion() string {
	return "v0.0.1"
}

func exit(c *cli.Context, message interface{}) error {
	if err := c.App.Run([]string{"dawn", "help", c.Command.Name}); err != nil {
		return cli.Exit(err, 1)
	}
	fmt.Println()

	msg := fmt.Sprintf("%s %s: %v", c.App.Name, c.Command.Name, message)
	return cli.Exit(msg, 1)
}

func success(message interface{}) error {
	return cli.Exit(message, 0)
}
