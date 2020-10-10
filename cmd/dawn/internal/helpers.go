package internal

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/urfave/cli/v2"
)

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

func createFile(filePath, content string) (err error) {
	var f *os.File
	if f, err = os.Create(filePath); err != nil {
		return
	}

	defer func() { _ = f.Close() }()

	_, err = f.WriteString(content)

	return
}

var execCommand = exec.Command

func runCmd(name string, arg ...string) (err error) {
	cmd := execCommand(name, arg...)

	var (
		stderr io.ReadCloser
		stdout io.ReadCloser
	)

	if stderr, err = cmd.StderrPipe(); err != nil {
		return
	}
	defer func() {
		_ = stderr.Close()
	}()
	go func() { _, _ = io.Copy(os.Stderr, stderr) }()

	if stdout, err = cmd.StdoutPipe(); err != nil {
		return
	}
	defer func() {
		_ = stdout.Close()
	}()
	go func() { _, _ = io.Copy(os.Stdout, stdout) }()

	if err = cmd.Run(); err != nil {
		err = fmt.Errorf("failed to run %s", cmd.String())
	}

	return
}

func formatLatency(d time.Duration) time.Duration {
	switch {
	case d > time.Second:
		return d.Truncate(time.Second / 100)
	case d > time.Millisecond:
		return d.Truncate(time.Millisecond / 100)
	case d > time.Microsecond:
		return d.Truncate(time.Microsecond / 100)
	default:
		return d
	}
}
