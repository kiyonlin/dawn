package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/urfave/cli/v2"
)

var newProject = &cli.Command{
	Name:      "new",
	Aliases:   []string{"n"},
	Usage:     "Generate a new dawn project",
	UsageText: "new project [mod name] [options]",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "template",
			Aliases: []string{"t"},
			Usage:   "project templates: default, custom, app",
			Value:   "default",
		},
	},
	Action: func(c *cli.Context) error {
		if c.Args().Len() == 0 {
			return exit(c, "Missing project name")
		}
		now := time.Now()

		projectName := c.Args().First()

		modName := projectName
		if c.Args().Len() > 1 {
			modName = c.Args().Get(1)
		}

		dir, err := os.Getwd()
		if err != nil {
			return exit(c, err)
		}

		projectPath := dir + "/" + projectName
		if err = createProject(projectPath, modName, c.String("template")); err != nil {
			return exit(c, err)
		}

		return success(fmt.Sprintf(newSuccessTemplate, projectPath, modName, projectName, time.Since(now)))
	},
}

func createProject(projectPath, modName, template string) (err error) {
	if err = os.Mkdir(projectPath, 0750); err != nil {
		return
	}

	defer func() {
		if err != nil {
			if e := os.RemoveAll(projectPath); e != nil {
				err = e
			}
		}
	}()

	if err = os.Chdir(projectPath); err != nil {
		return
	}

	// create main.go
	var mainFile *os.File
	if mainFile, err = os.Create(fmt.Sprintf("%s/main.go", projectPath)); err != nil {
		return
	}
	defer func() {
		if e := mainFile.Close(); e != nil {
			err = e
		}
	}()

	if _, err = mainFile.WriteString(newQuickstartTemplate); err != nil {
		return
	}

	var (
		rc  io.ReadCloser
		buf = make([]byte, 1024)
		n   int
	)

	cmdInit := exec.Command("go", "mod", "init", modName)
	if rc, err = cmdInit.StderrPipe(); err != nil {
		return
	}
	if err = cmdInit.Start(); err != nil {
		return
	}
	for {
		if n, err = rc.Read(buf); err != nil {
			if err == io.EOF {
				break
			}
		}
		_, _ = os.Stdout.Write(buf[:n])
	}

	if err = cmdInit.Wait(); err != nil {
		return
	}

	cmdTidy := exec.Command("go", "mod", "tidy")

	if rc, err = cmdTidy.StderrPipe(); err != nil {
		return
	}
	if err = cmdTidy.Start(); err != nil {
		return
	}
	for {
		if n, err = rc.Read(buf); err != nil {
			if err == io.EOF {
				break
			}
		}
		_, _ = os.Stdout.Write(buf[:n])
	}

	return cmdTidy.Wait()
}

var (
	newSuccessTemplate = `
Scaffolding project in %s (module %s)

  Done. Now run:

  cd %s
  go run .

✨  Done in %s.
`

	newQuickstartTemplate = `
package main

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
)
