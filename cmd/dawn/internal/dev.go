package internal

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/fsnotify/fsnotify"

	"github.com/urfave/cli/v2"
)

var Dev = &cli.Command{
	Name:      "dev",
	Usage:     "Rerun the dawn project if watched files change",
	UsageText: "dawn dev [options]",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "root",
			Aliases: []string{"r"},
			Usage:   "root path for watch, all files must be under root",
			Value:   ".",
		},
		&cli.StringFlag{
			Name:    "path",
			Aliases: []string{"p"},
			Usage:   "path for go run",
			Value:   ".",
		},
		&cli.StringSliceFlag{
			Name:    "extensions",
			Aliases: []string{"ext"},
			Usage:   "file extensions to watch",
			Value:   cli.NewStringSlice("go", "tmpl", "tpl", "html"),
		},
		&cli.StringSliceFlag{
			Name:    "exclude_dirs",
			Aliases: []string{"ed"},
			Usage:   "ignore these directories",
			Value:   cli.NewStringSlice("assets", "tmp", "vendor", "node_modules"),
		},
		&cli.StringSliceFlag{
			Name:    "exclude_files",
			Aliases: []string{"ef"},
			Usage:   "ignore these files",
			Value:   cli.NewStringSlice(),
		},
		&cli.DurationFlag{
			Name:    "delay",
			Aliases: []string{"d"},
			Usage:   "delay to trigger rerun",
			Value:   time.Second,
		},
	},
	Action: func(c *cli.Context) error {
		return newEscort(c).run()
	},
}

type escort struct {
	root         string
	path         string
	extensions   []string
	excludeDirs  []string
	excludeFiles []string
	delay        time.Duration

	w *fsnotify.Watcher
}

func newEscort(c *cli.Context) *escort {
	return &escort{
		root:         c.String("root"),
		path:         c.String("path"),
		extensions:   c.StringSlice("extensions"),
		excludeDirs:  c.StringSlice("exclude_dirs"),
		excludeFiles: c.StringSlice("exclude_files"),
		delay:        c.Duration("delay"),
	}
}

func (e *escort) run() (err error) {
	if err = e.init(); err != nil {
		return
	}
	defer func() { _ = e.w.Close() }()

	cmd := execCommand("go", "run", e.path)

	go watchStderrPipe(cmd)
	go watchStdoutPipe(cmd)

	fmt.Println("starting:", cmd.String())
	return nil
}

func (e *escort) init() (err error) {
	if e.w, err = fsnotify.NewWatcher(); err != nil {
		return
	}

	// normalize root

	return
}

func watchStdoutPipe(cmd *exec.Cmd) {

}

func watchStderrPipe(cmd *exec.Cmd) {

}

func (e *escort) hitExtension(ext string) bool {
	for _, e := range e.extensions {
		if ext == e {
			return true
		}
	}

	return false
}

func (e *escort) ignoredDirs(dir string) bool {
	for _, d := range e.excludeDirs {
		if dir == d {
			return true
		}
	}

	return false
}

func (e *escort) ignoredFiles(filename string) bool {
	for _, f := range e.excludeFiles {
		if filename == f {
			return true
		}
	}

	return false
}
