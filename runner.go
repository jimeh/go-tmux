package tmux

type Runner interface {
	Run(string, ...string) ([]byte, error)
}
