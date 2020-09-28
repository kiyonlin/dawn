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
