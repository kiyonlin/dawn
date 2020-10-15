package dawn

import (
	"fmt"
	"os"
	"os/exec"
)

const envDaemonKey = "ENV_DAWN_DAEMON"
const envDaemonVal = "true"

func (s *Sloop) daemon() (err error) {
	return
}

func (s *Sloop) spawn() (cmd *exec.Cmd, err error) {
	cmd = &exec.Cmd{
		Path:        os.Args[0],
		Args:        os.Args,
		Env:         append(os.Environ(), fmt.Sprintf("%s=%s", envDaemonKey, envDaemonVal)),
		SysProcAttr: newSysProcAttr(),
	}

	if s.StdoutLogFile != "" {
		if cmd.Stdout, err = os.OpenFile(s.StdoutLogFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600); err != nil {
			return
		}
	}

	if s.StderrLogFile != "" {
		if cmd.Stderr, err = os.OpenFile(s.StderrLogFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600); err != nil {
			return
		}
	}

	if err = cmd.Start(); err != nil {
		return
	}

	if !isDaemon() {
		os.Exit(0)
	}

	return
}

func isDaemon() bool {
	return os.Getenv(envDaemonKey) == envDaemonVal
}
