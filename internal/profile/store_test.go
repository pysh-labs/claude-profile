package profile

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDiscover_FindsClaudeDirs(t *testing.T) {
	home := t.TempDir()
	for _, name := range []string{".claude", ".claude-personal", ".claude-work", ".other"} {
		require.NoError(t, os.MkdirAll(filepath.Join(home, name), 0755))
	}

	got, err := Discover(home)
	require.NoError(t, err)

	names := make([]string, 0, len(got))
	for _, p := range got {
		names = append(names, p.Name)
	}
	require.ElementsMatch(t, []string{"default", "personal", "work"}, names)
}

func TestDiscover_ReadsLockFile(t *testing.T) {
	home := t.TempDir()
	dir := filepath.Join(home, ".claude-personal")
	require.NoError(t, os.MkdirAll(dir, 0755))
	require.NoError(t, os.WriteFile(
		filepath.Join(dir, "profile.lock.json"),
		[]byte(`{"name":"personal","label":"personal","color":"green"}`),
		0644,
	))

	got, err := Discover(home)
	require.NoError(t, err)
	require.Len(t, got, 1)
	require.Equal(t, "green", got[0].Color)
}

func TestActive_FromConfigDir(t *testing.T) {
	home := t.TempDir()
	dir := filepath.Join(home, ".claude-work")
	require.NoError(t, os.MkdirAll(dir, 0755))
	t.Setenv("CLAUDE_CONFIG_DIR", dir)

	require.Equal(t, "work", Active(home))
}

func TestActive_DefaultWhenUnset(t *testing.T) {
	home := t.TempDir()
	t.Setenv("CLAUDE_CONFIG_DIR", "")
	require.Equal(t, "default", Active(home))
}
