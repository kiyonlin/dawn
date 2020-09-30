package internal

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"

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
		rc  io.ReadCloser
		buf = make([]byte, 1024)
		n   int
	)

	if rc, err = cmd.StderrPipe(); err != nil {
		return
	}

	if err = cmd.Start(); err != nil {
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

	if err = cmd.Wait(); err != nil {
		err = errors.New(fmt.Sprintf("failed to run %s", cmd.String()))
	}

	return
}
