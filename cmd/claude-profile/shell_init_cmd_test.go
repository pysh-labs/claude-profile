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
	personal := filepath.Join(home, ".claude-personal")
	require.NoError(t, os.MkdirAll(personal, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(personal, "profile.lock.json"), []byte(`{}`), 0644))

	cmd := newRootCmd("test")
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"shell-init", "zsh"})
	require.NoError(t, cmd.Execute())

	require.Contains(t, out.String(), "alias claude-default=")
	require.Contains(t, out.String(), "alias claude-personal=")
}

func TestShellInitCmd_HidesUnmanaged(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	require.NoError(t, os.MkdirAll(filepath.Join(home, ".claude"), 0755))
	require.NoError(t, os.MkdirAll(filepath.Join(home, ".claude-flow"), 0755))
	managed := filepath.Join(home, ".claude-work")
	require.NoError(t, os.MkdirAll(managed, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(managed, "profile.lock.json"), []byte(`{}`), 0644))

	cmd := newRootCmd("test")
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"shell-init", "zsh"})
	require.NoError(t, cmd.Execute())

	require.Contains(t, out.String(), "alias claude-default=")
	require.Contains(t, out.String(), "alias claude-work=")
	require.NotContains(t, out.String(), "alias claude-flow=")
}

func TestShellInitCmd_AllFlagShowsUnmanaged(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	require.NoError(t, os.MkdirAll(filepath.Join(home, ".claude"), 0755))
	require.NoError(t, os.MkdirAll(filepath.Join(home, ".claude-flow"), 0755))

	cmd := newRootCmd("test")
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"shell-init", "zsh", "--all"})
	require.NoError(t, cmd.Execute())

	require.Contains(t, out.String(), "alias claude-flow=")
}
