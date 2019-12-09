package tmux

import (
	"bufio"
	"bytes"
	"regexp"
	"strconv"
)

var optMatcher = regexp.MustCompile(`^\s*([@\w-][\w-]+)\s+(.*)$`)

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

// GetOptions uses the "show-options" command to get all options within given
// SCOPE.
func (s *Tmux) GetOptions(scope Scope) (Options, error) {
	out, err := s.Exec("show-options", ScopeToFlags(scope))
	if err != nil {
		return nil, err
	}

	return s.parseOptions(out), nil
}

func (s *Tmux) parseOptions(options []byte) Options {
	scanner := bufio.NewScanner(bytes.NewBuffer(options))
	result := Options{}

	for scanner.Scan() {
		match := optMatcher.FindSubmatch(scanner.Bytes())
		if len(match) > 2 {
			key := string(match[1])
			val := string(match[2])

			unquoted, err := strconv.Unquote(val)
			if err == nil {
				result[key] = unquoted
			} else {
				result[key] = val
			}
		}
	}

	return result
}
