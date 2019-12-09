package tmux

// OptionsScope represents one of the five scopes that Tmux holds options
// within.
type OptionsScope int

const (
	Server OptionsScope = iota + 1
	GlobalSession
	Session
	GlobalWindow
	Window
)

// OptionsScopeFlags converts a given OptionsScope to the command line flags
// needed to restrict "set-option" and "show-options" commands to the scope in
// question.
func OptionsScopeFlags(scope OptionsScope) string {
	switch scope {
	case 0, Session:
		return ""
	case Server:
		return "-s"
	case GlobalSession:
		return "-g"
	case GlobalWindow:
		return "-gw"
	case Window:
		return "-w"
	default:
		return ""
	}
}
