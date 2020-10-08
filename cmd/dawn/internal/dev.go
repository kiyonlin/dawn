package internal

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
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
			Name:    "target",
			Aliases: []string{"t"},
			Usage:   "target path for go build",
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
	target       string
	binPath      string
	extensions   []string
	excludeDirs  []string
	excludeFiles []string
	delay        time.Duration

	w     *fsnotify.Watcher
	bin   *exec.Cmd
	hitCh chan struct{}
}

func newEscort(c *cli.Context) *escort {
	return &escort{
		root:         c.String("root"),
		target:       c.String("target"),
		extensions:   c.StringSlice("extensions"),
		excludeDirs:  c.StringSlice("exclude_dirs"),
		excludeFiles: c.StringSlice("exclude_files"),
		delay:        c.Duration("delay"),
		hitCh:        make(chan struct{}, 1),
	}
}

func (e *escort) run() (err error) {
	if err = e.init(); err != nil {
		return
	}

	defer func() {
		_ = e.w.Close()
		_ = os.Remove(e.binPath)
	}()

	go e.runBin()
	go e.watchingBin()
	go e.watchingFiles()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	<-c

	return nil
}

func (e *escort) init() (err error) {
	if e.w, err = fsnotify.NewWatcher(); err != nil {
		return
	}

	// normalize root
	if e.root, err = filepath.Abs(e.root); err != nil {
		return
	}

	// create bin target
	var f *os.File
	if f, err = ioutil.TempFile("", ""); err != nil {
		return
	}
	defer func() {
		if e := f.Close(); e != nil {
			err = e
		}
	}()

	e.binPath = f.Name()

	return
}

func (e *escort) watchingFiles() {
	// walk root and add all dirs
	e.walkForWatcher(e.root)

	var (
		info os.FileInfo
		err  error
	)

	for {
		select {
		case event, ok := <-e.w.Events:
			if ok {
				log.Println("debug========", event)
				p, op := event.Name, event.Op

				if isRemoved(op) {
					e.tryRemoveWatch(p)
					continue
				}

				if info, err = os.Stat(p); err != nil {
					log.Printf("failed to get info of %s: %s\n", p, err)
					continue
				}

				base := filepath.Base(p)

				if info.IsDir() && isCreated(op) {
					log.Println("add", p, "to watch")
					e.walkForWatcher(p)
					continue
				}

				if e.ignoredFiles(base) {
					continue
				}

				if e.hitExtension(filepath.Ext(base)) {
					e.hitCh <- struct{}{}
				}
			}
		case err, ok := <-e.w.Errors:
			if ok {
				log.Printf("watcher error: %s\n", err)
			}
		}
	}
}

func (e *escort) watchingBin() {
	var timer *time.Timer
	for {
		select {
		case <-e.hitCh:
			// reset timer
			if timer != nil && !timer.Stop() {
				select {
				case <-timer.C:
				default:
				}
			}
			timer = time.AfterFunc(e.delay, e.runBin)
		}
	}
}

func (e *escort) runBin() {
	if e.bin != nil {
		if err := e.bin.Process.Kill(); err != nil {
			log.Printf("failed to kill old bin (pid %d): %s\n", e.bin.Process.Pid, err)
		}
		_, _ = e.bin.Process.Wait()
		e.bin = nil
	}

	// build target
	compile := execCommand("go", "build", "-o", e.binPath, e.target)
	if out, err := compile.CombinedOutput(); err != nil {
		log.Printf("failed to compile %s: %s\n", e.target, out)
		return
	}

	e.bin = execCommand(e.binPath)

	e.watchPipe()

	if err := e.bin.Start(); err != nil {
		log.Printf("failed to start bin: %s\n", err)
	}

	log.Println("pid", e.bin.Process.Pid)
}

func (e *escort) watchPipe() {
	log.Println("starting watch pipe", e.bin.String())
	// TODO
}

func (e *escort) walkForWatcher(root string) {
	if err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info != nil && !info.IsDir() {
			return nil
		}

		base := filepath.Base(path)

		// exclude hidden directories like .git, .idea, etc.
		if isHiddenDirectory(base) {
			return filepath.SkipDir
		}

		if e.ignoredDirs(base) {
			return filepath.SkipDir
		}

		log.Println("check", path)

		return e.w.Add(path)
	}); err != nil {
		log.Printf("failed to walk root %s\n", e.root)
	}
}

func (e *escort) tryRemoveWatch(p string) {
	log.Println("try to remove", p, "from watch")
	if err := e.w.Remove(p); err != nil && !strings.Contains(err.Error(), "non-existent") {
		log.Printf("failed to remove %s from watch: %s\n", p, err)
	}
}

func (e *escort) hitExtension(ext string) bool {
	if ext == "" {
		return false
	}
	// remove '.'
	ext = ext[1:]
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

func isHiddenDirectory(d string) bool {
	return len(d) > 1 && d[0] == '.'
}

func isRemoved(op fsnotify.Op) bool {
	return op&fsnotify.Remove != 0
}

func isCreated(op fsnotify.Op) bool {
	return op&fsnotify.Create != 0
}
