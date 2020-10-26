package daemon

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/kiyonlin/dawn/config"
)

const envDaemon = "DAWN_DAEMON=1"
const envDaemonWorker = "DAWN_DAEMON_WORKER=1"

var stdoutLogFile *os.File
var stderrLogFile *os.File

func Run() {
	if isWorker() {
		log.Println("I'm a worker")
		return
	}

	// Panic if the initial spawned daemon process has error
	if _, err := spawn(true); err != nil {
		panic(fmt.Sprintf("dawn: failed to run in daemon mode: %s", err))
	}

	setupLogFiles()
	defer teardownLogFiles()

	_, _ = stdoutLogFile.WriteString("xxxxxx")

	var (
		cmd    *exec.Cmd
		err    error
		count  int
		max    = config.GetInt("daemon.tries", 10)
		logger = log.New(stderrLogFile, "", log.LstdFlags)
	)

	for {
		if count++; count > max {
			break
		}

		if cmd, err = spawn(false); err != nil {
			continue
		}

		// reset count
		count = 0

		err = cmd.Wait()

		logger.Printf("%s(pid:%d) exist with err: %v", cmd.Args[0], cmd.Process.Pid, err)
	}

	logger.Printf("Already attempted %d times", max)
}

func spawn(skip bool) (cmd *exec.Cmd, err error) {
	if inDaemon() && skip {
		log.Println("skip in daemon")
		return
	}

	cmd = &exec.Cmd{
		Path:        os.Args[0],
		Args:        os.Args,
		Env:         parseEnv(),
		SysProcAttr: newSysProcAttr(),
	}

	if inDaemon() {
		if cmd.Stdout = ioutil.Discard; stdoutLogFile != nil {
			cmd.Stdout = stdoutLogFile
		}

		if cmd.Stderr = ioutil.Discard; stderrLogFile != nil {
			cmd.Stderr = stderrLogFile
		}
	}

	if err = cmd.Start(); err != nil {
		return
	}

	log.Println("master process pid", cmd.Process.Pid)

	// Exit main process
	if !inDaemon() {
		log.Println("exit main process")
		os.Exit(0)
	}

	return
}

func parseEnv() []string {
	env := os.Environ()
	if !inDaemon() {
		env = append(env, envDaemon)
	} else if !isWorker() {
		env = append(env, envDaemonWorker)
	}

	return env
}

func inDaemon() bool {
	_, ok := os.LookupEnv(envDaemon)
	return ok
}

func isWorker() bool {
	_, ok := os.LookupEnv(envDaemonWorker)
	return ok
}

func setupLogFiles() {
	var err error
	if f := config.GetString("daemon.stdoutLogFile"); f != "" {
		if stdoutLogFile, err = os.OpenFile(f, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600); err != nil {
			panic(fmt.Sprintf("dawn: failed to open stdout log file %s: %s", f, err))
		}
	}

	if f := config.GetString("daemon.stderrLogFile"); f != "" {
		if stderrLogFile, err = os.OpenFile(f, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600); err != nil {
			panic(fmt.Sprintf("dawn: failed to open stderr log file %s: %s", f, err))
		}
	}
}

func teardownLogFiles() {
	if stdoutLogFile != nil {
		_ = stdoutLogFile.Close()
	}

	if stderrLogFile != nil {
		_ = stderrLogFile.Close()
	}
}
