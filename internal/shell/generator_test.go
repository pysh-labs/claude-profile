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

func TestGenerate_PosixIncludesUseFunction(t *testing.T) {
	for _, sh := range []string{"zsh", "bash"} {
		got, err := Generate(sh, nil)
		require.NoError(t, err)
		require.Contains(t, got, "claude-profile() {", "shell %s", sh)
		require.Contains(t, got, `if [ "$1" = "use" ]`, "shell %s", sh)
		require.Contains(t, got, "command claude-profile which", "shell %s", sh)
		require.Contains(t, got, "export CLAUDE_CONFIG_DIR=", "shell %s", sh)
	}
}

func TestGenerate_FishIncludesUseFunction(t *testing.T) {
	got, err := Generate("fish", nil)
	require.NoError(t, err)
	require.Contains(t, got, "function claude-profile")
	require.Contains(t, got, `test "$argv[1]" = "use"`)
	require.Contains(t, got, "set -gx CLAUDE_CONFIG_DIR")
}
