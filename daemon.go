package dawn

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/kiyonlin/dawn/config"
)

const envDaemonKey = "ENV_DAWN_DAEMON"

var count = daemonCount()

func (s *Sloop) daemon() {
	var (
		cmd   *exec.Cmd
		err   error
		tries = 1
	)

	cmd, err = s.spawn()

	for {
		tries++
		if cmd, err = s.spawn(); err != nil {
			continue
		}

		if cmd == nil {
			break
		}

		err = cmd.Wait()
	}
}

func (s *Sloop) spawn() (cmd *exec.Cmd, err error) {
	count++

	if count <= daemonCount() {
		return nil, nil
	}

	cmd = &exec.Cmd{
		Path:        os.Args[0],
		Args:        os.Args,
		Env:         append(os.Environ(), fmt.Sprintf("%s=%d", envDaemonKey, count)),
		SysProcAttr: newSysProcAttr(),
	}

	if stdoutLogFile := config.GetString("daemon.stdoutLogFile"); stdoutLogFile != "" {
		if cmd.Stdout, err = os.OpenFile(stdoutLogFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600); err != nil {
			return
		}
	}

	if stderrLogFile := config.GetString("daemon.stderrLogFile"); stderrLogFile != "" {
		if cmd.Stderr, err = os.OpenFile(stderrLogFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600); err != nil {
			return
		}
	}

	if err = cmd.Start(); err != nil {
		return
	}

	if daemonCount() == 0 {
		os.Exit(0)
	}

	return
}

func daemonCount() int {
	c, _ := strconv.Atoi(os.Getenv(envDaemonKey))
	return c
}
