package internal

import (
	"fmt"
	"os"
	"time"

	"github.com/urfave/cli/v2"
)

// NewProject command generates a new dawn project
var NewProject = &cli.Command{
	Name:      "new",
	Aliases:   []string{"n"},
	Usage:     "Generate a new dawn project",
	UsageText: "dawn new [options] project [mod name]",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "app",
			Usage: "create an application project",
		},
	},
	Action: func(c *cli.Context) error {
		if !c.Args().Present() {
			return exit(c, "Missing project name")
		}
		now := time.Now()

		projectName := c.Args().First()

		modName := projectName
		if c.Args().Len() > 1 {
			modName = c.Args().Get(1)
		}

		dir, _ := os.Getwd()

		projectPath := dir + "/" + projectName
		if err := createProject(projectPath, modName, c.Bool("app")); err != nil {
			return exit(c, err)
		}

		return success(fmt.Sprintf(newSuccessTemplate, projectPath, modName, projectName, time.Since(now)))
	},
}

func createProject(projectPath, modName string, isApp bool) (err error) {
	if err = os.Mkdir(projectPath, 0750); err != nil {
		return
	}

	defer func() {
		if err != nil {
			_ = os.RemoveAll(projectPath)
		}
	}()

	if err = os.Chdir(projectPath); err != nil {
		return
	}

	// create main.go
	if err = createFile(fmt.Sprintf("%s/main.go", projectPath),
		templateContent(isApp)); err != nil {
		return
	}

	if err = runCmd("go", "mod", "init", modName); err != nil {
		return
	}

	return runCmd("go", "mod", "tidy")
}

func templateContent(isApp bool) string {
	if isApp {
		return newAppTemplate
	}
	return newWebTemplate
}

var (
	newSuccessTemplate = `
Scaffolding project in %s (module %s)

  Done. Now run:

  cd %s
  go run .

✨  Done in %s.
`

	newWebTemplate = `package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/kiyonlin/dawn"
	"github.com/kiyonlin/dawn/fiberx"
)

func main() {
	sloop := dawn.Default()

	router := sloop.Router()
	// GET /  =>  Welcome to dawn 👋
	router.Get("/", func(c *fiber.Ctx) error {
		return fiberx.Message(c, "Welcome to dawn 👋")
	})

	log.Println(sloop.Run(":3000"))
}
`

	newAppTemplate = `package main

import (
	"flag"

	"github.com/kiyonlin/dawn"
	"github.com/kiyonlin/dawn/config"
	"github.com/kiyonlin/dawn/db/redis"
	"github.com/kiyonlin/dawn/db/sql"
	"github.com/kiyonlin/dawn/log"
)

func main() {
	config.Load("./")
	config.LoadEnv()

	log.InitFlags(nil)
	flag.Parse()
	defer log.Flush()

	sloop := dawn.New(dawn.Modulers(
		sql.New(),
		redis.New(),
		// add custom module 
	))

	defer sloop.Cleanup()

	sloop.Setup().Watch()
}
`
)
