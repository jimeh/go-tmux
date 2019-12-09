package tmux

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRunner struct {
	mock.Mock
}

func (s *MockRunner) Run(command string, args ...string) ([]byte, error) {
	called := s.Called(append([]string{command}, args...))

	return called.Get(0).([]byte), called.Error(1)
}

func TestNewTmux(t *testing.T) {
	tmux := New()

	assert.IsType(t, &Tmux{}, tmux)
	assert.IsType(t, &ExecRunner{}, tmux.Runner)
}

func TestTmuxExec(t *testing.T) {
	tests := []struct {
		binPath    string
		socketName string
		socketPath string
		baseArgs   []string
		args       []string
		out        []byte
		error      error
	}{
		{
			args:     []string{"new-session", "-d"},
			baseArgs: []string{"tmux"},
		},
		{
			args:     []string{"new-session", "-d"},
			baseArgs: []string{"tmux"},
		},
		{
			args:     []string{"list-sessions"},
			baseArgs: []string{"tmux"},
			out:      []byte("0: 1 windows (created Fri Dec  6 23:45:19 2019)"),
		},
		{
			binPath:  "/opt/tmux/bin/tmux",
			args:     []string{"list-sessions"},
			baseArgs: []string{"/opt/tmux/bin/tmux"},
		},
		{
			binPath:    "/opt/tmux/bin/tmux",
			socketName: "test-sock",
			args:       []string{"list-sessions"},
			baseArgs:   []string{"/opt/tmux/bin/tmux", "-L", "test-sock"},
		},
		{
			binPath:    "/opt/tmux/bin/tmux",
			socketPath: "/tmp/tmux.sock",
			args:       []string{"list-sessions"},
			baseArgs:   []string{"/opt/tmux/bin/tmux", "-S", "/tmp/tmux.sock"},
		},
		{
			binPath:    "/opt/tmux/bin/tmux",
			socketName: "test-sock",
			socketPath: "/tmp/tmux.sock",
			args:       []string{"list-sessions"},
			baseArgs:   []string{"/opt/tmux/bin/tmux", "-S", "/tmp/tmux.sock"},
		},
		{
			args:     []string{"new-session", "-d"},
			baseArgs: []string{"tmux"},
			error:    errors.New("Something went wrong"),
		},
	}

	for _, tt := range tests {
		runner := new(MockRunner)
		runner.On("Run", append(tt.baseArgs, tt.args...)).
			Return(tt.out, tt.error)

		tmux := Tmux{
			Runner:     runner,
			BinPath:    tt.binPath,
			SocketName: tt.socketName,
			SocketPath: tt.socketPath,
		}

		out, err := tmux.Exec(tt.args...)

		if tt.error == nil {
			assert.NoError(t, err)
		} else {
			assert.Equal(t, tt.error, err)
		}

		if tt.out != nil {
			assert.Equal(t, tt.out, out)
		}

		runner.AssertExpectations(t)
	}
}

func TestTmuxBinary(t *testing.T) {
	tests := []struct {
		binPath    string
		executable string
	}{
		{binPath: "/opt/tmux/bin/tmux", executable: "/opt/tmux/bin/tmux"},
		{executable: "tmux"},
	}

	for _, tt := range tests {
		tmux := &Tmux{}
		tmux.BinPath = tt.binPath

		assert.Equal(t, tt.executable, tmux.Binary())
	}
}

func TestTmuxArgs(t *testing.T) {
	tests := []struct {
		socketName string
		socketPath string
		args       []string
	}{
		{args: []string{}},
		{socketName: "foo", args: []string{"-L", "foo"}},
		{socketPath: "/tmp/bar", args: []string{"-S", "/tmp/bar"}},
		{
			socketName: "foo",
			socketPath: "/tmp/bar",
			args:       []string{"-S", "/tmp/bar"},
		},
	}

	for _, tt := range tests {
		tmux := &Tmux{}
		tmux.SocketName = tt.socketName
		tmux.SocketPath = tt.socketPath

		assert.Equal(t, tt.args, tmux.Args())
	}
}

func TestTmuxGetOptions(t *testing.T) {
	tests := []struct {
		flags string
		scope OptionsScope
		opts  map[string]string
		out   []byte
		error error
	}{
		{
			opts: map[string]string{"hello-world": "FooBar"},
			out:  []byte(`hello-world FooBar`),
		},
		{
			scope: Server,
			flags: "-s",
			opts:  map[string]string{"hello-world": "Foo Bar"},
			out:   []byte(`hello-world "Foo Bar"`),
		},
		{
			scope: GlobalSession,
			flags: "-g",
			opts:  map[string]string{"hello-world": "Foo Bar"},
			out:   []byte(`hello-world   "Foo Bar"`),
		},
		{
			scope: GlobalWindow,
			flags: "-gw",
			opts:  map[string]string{"hello-world": "  Foo Bar   "},
			out:   []byte(`hello-world "  Foo Bar   "`),
		},
		{
			scope: Window,
			flags: "-w",
			opts:  map[string]string{"@foo": "bar"},
			out:   []byte(`@foo bar`),
		},
		{
			opts: map[string]string{
				"@foo":        "bar",
				"@themepack":  "powerline/default/green",
				"status-left": "This Is Left",
			},
			out: []byte(`
  @foo bar
@themepack "powerline/default/green"
status-left This Is Left
`),
		},
	}

	for _, tt := range tests {
		if tt.scope == 0 {
			tt.scope = Session
		}

		runner := new(MockRunner)
		runner.On("Run", append([]string{"tmux", "show-options"}, tt.flags)).
			Return(tt.out, tt.error)

		tmux := Tmux{Runner: runner}

		opts, err := tmux.GetOptions(tt.scope)

		if tt.error == nil {
			assert.NoError(t, err)
			assert.Equal(t, tt.opts, opts)
		} else {
			assert.Equal(t, tt.error, err)
		}

		runner.AssertExpectations(t)
	}
}
