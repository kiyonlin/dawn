package internal

import (
	"fmt"
	"os"
	"os/exec"
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
	Action: newAction,
}

func newAction(c *cli.Context) error {
	if !c.Args().Present() {
		return exit(c, "Missing project name")
	}
	start := time.Now()

	projectName := c.Args().First()

	modName := projectName
	if c.Args().Len() > 1 {
		modName = c.Args().Get(1)
	}

	dir, _ := os.Getwd()

	projectPath := fmt.Sprintf("%s%c%s", dir, os.PathSeparator, projectName)
	if err := createProject(projectPath, modName, c.Bool("app")); err != nil {
		return exit(c, err)
	}

	return success(fmt.Sprintf(newSuccessTemplate,
		projectPath, modName, projectName, formatLatency(time.Since(start))))
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
	if err = createFile(fmt.Sprintf("%s%cmain.go", projectPath, os.PathSeparator),
		templateContent(isApp)); err != nil {
		return
	}

	if err = runCmd("go", "mod", "init", modName); err != nil {
		return
	}

	if err = runCmd("go", "mod", "tidy"); err != nil {
		return
	}

	if err = createConfigs(projectPath); err != nil {
		return
	}

	// init git
	return initGit(projectPath)
}

func createConfigs(projectPath string) (err error) {
	// create config.toml
	if err = createFile(fmt.Sprintf("%s%cconfig.toml", projectPath, os.PathSeparator), configTemplate); err != nil {
		return
	}

	// create config.example.toml
	if err = createFile(fmt.Sprintf("%s%cconfig.example.toml", projectPath, os.PathSeparator), configTemplate); err != nil {
		return
	}

	return
}

func initGit(projectPath string) (err error) {
	var git string
	if git, err = exec.LookPath("git"); err != nil {
		return nil
	}

	// create .gitignore
	if err = createFile(fmt.Sprintf("%s%c.gitignore", projectPath, os.PathSeparator), gitignoreTemplate); err != nil {
		return
	}

	if err = runCmd(git, "init"); err != nil {
		return
	}

	if err = runCmd(git, "add", "-A"); err != nil {
		return
	}

	if err = runCmd(git, "commit", "-m", "dawn init"); err != nil {
		return
	}

	return
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
  dawn dev

✨  Done in %s.
`

	newWebTemplate = `package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/kiyonlin/dawn"
	"github.com/kiyonlin/dawn/config"
	"github.com/kiyonlin/dawn/db/redis"
	"github.com/kiyonlin/dawn/db/sql"
	"github.com/kiyonlin/dawn/fiberx"
)

func main() {
	config.Load("./")

	sloop := dawn.Default().
		AddModulers(sql.New(), redis.New())

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

	gitignoreTemplate = `# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary, built with ` + "`go test -c`" + `
*.test

# Output of the go coverage tool, specifically when used with LiteIDE
*.out

# Dependency directories (remove the comment below to include it)
# vendor/

# IDEs
.idea

# Application
config.toml
`

	configTemplate = `Debug = true

[Sql]
Default = "testing"

[Sql.Connections]
[Sql.Connections.Testing]
Driver = "sqlite"
Database = "file:dawn?mode=memory&cache=shared&_fk=1"
Prefix = "dawn_"
Log = true
# Uncomment to use other sql connections
#[Sql.Connections.Mysql]
#Driver = "mysql"
#Username = "username"
#Password = "password"
#Host = "127.0.0.1"
#Port = "3306"
#Database = "database"
#Location = "Asia/Shanghai"
#Charset = "utf8mb4"
#ParseTime = true
#Prefix = "dawn_"
#Log = false
#MaxIdleConns = 10
#MaxOpenConns = 100
#ConnMaxLifetime = "5m"
#
#[Sql.Connections.Postgres]
#Driver = "postgres"
#Host = "127.0.0.1"
#Port = "5432"
#Database = "database"
#Username = "username"
#Password = "password"
#Sslmode = "disable"
#TimeZone = "Asia/Shanghai"
#Prefix = "dawn_"
#Log = false
#MaxIdleConns = 10
#MaxOpenConns = 100
#ConnMaxLifetime = "5m"

[Redis]
Default = "default"

[Redis.Connections]
[Redis.Connections.default]
Network = "tcp"
Addr = "127.0.0.1:6379"
Username = ""
Password = ""
DB = 0
MaxRetries = 5
DialTimeout = "5s"
ReadTimeout = "5s"
WriteTimeout = "5s"
PoolSize = 1024
MinIdleConns = 10
MaxConnAge = "1m"
PoolTimeout = "1m"
IdleTimeout = "1m"
IdleCheckFrequency = "1m"
`
)
