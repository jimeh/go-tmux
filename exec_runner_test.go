package tmux

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecRunnerRun(t *testing.T) {
	tests := []struct {
		command string
		args    []string
	}{
		{command: "pwd"},
		{command: "hostname"},
		{command: "uname"},
	}

	for _, tt := range tests {
		runner := &ExecRunner{}

		expected, err := exec.Command(tt.command, tt.args...).CombinedOutput()
		assert.NoError(t, err)
		actual, err := runner.Run(tt.command, tt.args...)
		assert.NoError(t, err)

		assert.Equal(t, expected, actual)
	}
}
