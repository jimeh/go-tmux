package tmux

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptionsScopeFlags(t *testing.T) {
	tests := []struct {
		scope OptionsScope
		flags string
	}{
		{0, ""},
		{Server, "-s"},
		{GlobalSession, "-g"},
		{Session, ""},
		{GlobalWindow, "-gw"},
		{Window, "-w"},
		{38404, ""},
		{934, ""},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.flags, OptionsScopeFlags(tt.scope))
	}
}
