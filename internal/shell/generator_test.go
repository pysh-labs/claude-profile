package shell

import (
	"testing"

	"github.com/dmitriipyshinskii/claude-profile/internal/profile"
	"github.com/stretchr/testify/require"
)

func TestGenerate_Zsh(t *testing.T) {
	profiles := []profile.Profile{
		{Name: "default", Path: "/Users/me/.claude"},
		{Name: "personal", Path: "/Users/me/.claude-personal"},
	}
	got, err := Generate("zsh", profiles)
	require.NoError(t, err)
	require.Contains(t, got, `alias claude-default='CLAUDE_CONFIG_DIR=/Users/me/.claude claude'`)
	require.Contains(t, got, `alias claude-personal='CLAUDE_CONFIG_DIR=/Users/me/.claude-personal claude'`)
}

func TestGenerate_Bash(t *testing.T) {
	profiles := []profile.Profile{{Name: "x", Path: "/p"}}
	got, err := Generate("bash", profiles)
	require.NoError(t, err)
	require.Contains(t, got, `alias claude-x='CLAUDE_CONFIG_DIR=/p claude'`)
}

func TestGenerate_Fish(t *testing.T) {
	profiles := []profile.Profile{{Name: "x", Path: "/p"}}
	got, err := Generate("fish", profiles)
	require.NoError(t, err)
	require.Contains(t, got, `alias claude-x "CLAUDE_CONFIG_DIR=/p claude"`)
}

func TestGenerate_UnknownShell(t *testing.T) {
	_, err := Generate("powershell", nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "powershell")
}
