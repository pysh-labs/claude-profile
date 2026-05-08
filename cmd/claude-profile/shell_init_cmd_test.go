package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestShellInitCmd_Zsh(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	require.NoError(t, os.MkdirAll(filepath.Join(home, ".claude"), 0755))
	require.NoError(t, os.MkdirAll(filepath.Join(home, ".claude-personal"), 0755))

	cmd := newRootCmd("test")
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"shell-init", "zsh"})
	require.NoError(t, cmd.Execute())

	require.Contains(t, out.String(), "alias claude-default=")
	require.Contains(t, out.String(), "alias claude-personal=")
}
