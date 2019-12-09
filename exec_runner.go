package tmux

import (
	"os/exec"
)

type ExecRunner struct{}

func (r *ExecRunner) Run(command string, args ...string) ([]byte, error) {
	return exec.Command(command, args...).CombinedOutput()
}
