package main

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

var (
	needError bool
	errFlag   = struct{}{}
)

func fakeExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	if needError {
		cmd.Env = append(cmd.Env, "GO_WANT_HELPER_NEED_ERR=1")
	}
	return cmd
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	args := os.Args
	for len(args) > 0 {
		if args[0] == "--" {
			args = args[1:]
			break
		}
		args = args[1:]
	}

	if len(args) == 0 {
		_, _ = fmt.Fprintf(os.Stderr, "No command")
		os.Exit(2)
	}

	if os.Getenv("GO_WANT_HELPER_NEED_ERR") == "1" {
		_, _ = fmt.Fprintf(os.Stderr, "fake error")
		os.Exit(1)
	}

	os.Exit(0)
}

func setupCmd(flag ...struct{}) {
	execCommand = fakeExecCommand
	if len(flag) > 0 {
		needError = true
	}
}

func teardownCmd() {
	execCommand = exec.Command
	needError = false
}
