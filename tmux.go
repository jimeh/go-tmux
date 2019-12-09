package tmux

import (
	"bufio"
	"bytes"
	"regexp"
)

var optMatcher = regexp.MustCompile(`^\s*([@\w-][\w-]+)\s+(.*)$`)
var quote = []byte(`"`)

// Tmux enables easily running tmux commands.
type Tmux struct {
	BinPath    string
	SocketName string
	SocketPath string
	Runner     Runner
}

// New returns a Tmux objects with a Runner capable of executing shell commands.
func New() *Tmux {
	return &Tmux{Runner: &ExecRunner{}}
}

// Exec runs the given tmux command.
func (s *Tmux) Exec(args ...string) ([]byte, error) {
	args = append(s.Args(), args...)

	return s.Runner.Run(s.Binary(), args...)
}

func (s *Tmux) Binary() string {
	if s.BinPath != "" {
		return s.BinPath
	} else {
		return "tmux"
	}
}

func (s *Tmux) Args() []string {
	args := []string{}

	if s.SocketPath != "" {
		args = append(args, "-S", s.SocketPath)
	} else if s.SocketName != "" {
		args = append(args, "-L", s.SocketName)
	}

	return args
}

func (s *Tmux) GetOptions(scope OptionsScope) (map[string]string, error) {
	out, err := s.Exec("show-options", OptionsScopeFlags(scope))
	if err != nil {
		return nil, err
	}

	return s.parseOptions(out), nil
}

func (s *Tmux) parseOptions(options []byte) map[string]string {
	scanner := bufio.NewScanner(bytes.NewBuffer(options))
	result := map[string]string{}

	for scanner.Scan() {
		match := optMatcher.FindSubmatch(scanner.Bytes())
		if len(match) > 2 {
			result[string(match[1])] = string(s.unwrap(match[2], quote))
		}
	}

	return result
}

func (s *Tmux) unwrap(input, wrap []byte) []byte {
	if bytes.HasPrefix(input, wrap) && bytes.HasSuffix(input, wrap) {
		return bytes.TrimSuffix(bytes.TrimPrefix(input, wrap), wrap)
	}

	return input
}
