package internal

import (
	"context"
	"io"
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

	ctx       context.Context
	terminate context.CancelFunc

	w             *fsnotify.Watcher
	watcherEvents chan fsnotify.Event
	watcherErrors chan error
	sig           chan os.Signal

	bin        *exec.Cmd
	stdoutPipe io.ReadCloser
	stderrPipe io.ReadCloser
	hitCh      chan struct{}
	hitFunc    func()
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
		sig:          make(chan os.Signal, 1),
	}
}

func (e *escort) run() (err error) {
	if err = e.init(); err != nil {
		return
	}

	log.Println("Welcome to dawn dev 👋")

	defer func() {
		_ = e.w.Close()
		_ = os.Remove(e.binPath)
	}()

	go e.runBin()
	go e.watchingBin()
	go e.watchingFiles()

	signal.Notify(e.sig, syscall.SIGTERM, syscall.SIGINT)
	<-e.sig

	e.terminate()

	log.Println("See you next time 👋")

	return nil
}

func (e *escort) init() (err error) {
	if e.w, err = fsnotify.NewWatcher(); err != nil {
		return
	}

	e.watcherEvents = e.w.Events
	e.watcherErrors = e.w.Errors

	e.ctx, e.terminate = context.WithCancel(context.Background())

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

	e.hitFunc = e.runBin

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
		case <-e.ctx.Done():
			return
		case event := <-e.watcherEvents:
			p, op := event.Name, event.Op

			// ignore chmod
			if isChmoded(op) {
				continue
			}

			if isRemoved(op) {
				e.tryRemoveWatch(p)
				continue
			}

			if info, err = os.Stat(p); err != nil {
				log.Printf("Failed to get info of %s: %s\n", p, err)
				continue
			}

			base := filepath.Base(p)

			if info.IsDir() && isCreated(op) {
				log.Println("Add", p, "to watch")
				e.walkForWatcher(p)
				continue
			}

			if e.ignoredFiles(base) {
				continue
			}

			if e.hitExtension(filepath.Ext(base)) {
				e.hitCh <- struct{}{}
			}
		case err := <-e.watcherErrors:
			log.Printf("Watcher error: %v\n", err)
		}
	}
}

func (e *escort) watchingBin() {
	var timer *time.Timer
	for range e.hitCh {
		// reset timer
		if timer != nil && !timer.Stop() {
			select {
			case <-timer.C:
			default:
			}
		}
		timer = time.AfterFunc(e.delay, e.hitFunc)
	}
}

func (e *escort) runBin() {
	if e.bin != nil {
		e.cleanOldBin()
		log.Println("Restarting...")
	}

	// build target
	compile := execCommand("go", "build", "-o", e.binPath, e.target)
	if out, err := compile.CombinedOutput(); err != nil {
		log.Printf("Failed to compile %s: %s\n", e.target, out)
		return
	}

	log.Println("Compile done!")

	e.bin = execCommand(e.binPath)

	e.watchingPipes()

	if err := e.bin.Start(); err != nil {
		log.Printf("Failed to start bin: %s\n", err)
		e.bin = nil
		return
	}

	log.Println("New pid is", e.bin.Process.Pid)
}

func (e *escort) cleanOldBin() {
	defer func() {
		if e.stdoutPipe != nil {
			_ = e.stdoutPipe.Close()
		}
		if e.stderrPipe != nil {
			_ = e.stderrPipe.Close()
		}
	}()

	pid := e.bin.Process.Pid
	log.Println("Killing old pid", pid)
	if err := e.bin.Process.Kill(); err != nil {
		log.Printf("Failed to kill old pid %d: %s\n", pid, err)
	}

	_, _ = e.bin.Process.Wait()

	e.bin = nil
}

func (e *escort) watchingPipes() {
	var err error
	if e.stdoutPipe, err = e.bin.StdoutPipe(); err != nil {
		log.Printf("Failed to get stdout pipe: %s", err)
	} else {
		go func() { _, _ = io.Copy(os.Stdout, e.stdoutPipe) }()
	}

	if e.stderrPipe, err = e.bin.StderrPipe(); err != nil {
		log.Printf("Failed to get stderr pipe: %s", err)
	} else {
		go func() { _, _ = io.Copy(os.Stderr, e.stderrPipe) }()
	}
}

func (e *escort) walkForWatcher(root string) {
	if err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info != nil && !info.IsDir() {
			return nil
		}

		base := filepath.Base(path)

		if e.ignoredDirs(base) {
			return filepath.SkipDir
		}

		return e.w.Add(path)
	}); err != nil {
		log.Printf("Failed to walk root %s: %s\n", e.root, err)
	}
}

func (e *escort) tryRemoveWatch(p string) {
	if err := e.w.Remove(p); err != nil && !strings.Contains(err.Error(), "non-existent") {
		log.Printf("Failed to remove %s from watch: %s\n", p, err)
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
	// exclude hidden directories like .git, .idea, etc.
	if len(dir) > 1 && dir[0] == '.' {
		return true
	}

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

func isRemoved(op fsnotify.Op) bool {
	return op&fsnotify.Remove != 0
}

func isCreated(op fsnotify.Op) bool {
	return op&fsnotify.Create != 0
}

func isChmoded(op fsnotify.Op) bool {
	return op&fsnotify.Chmod != 0
}
