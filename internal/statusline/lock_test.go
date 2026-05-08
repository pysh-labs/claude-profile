package statusline

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadProfileMeta_ValidLock(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(
		filepath.Join(dir, "profile.lock.json"),
		[]byte(`{"name":"personal","label":"personal","color":"green"}`),
		0644,
	))

	p := LoadProfileMeta(dir)
	require.Equal(t, "personal", p.Label)
	require.Equal(t, "green", p.Color)
}

func TestLoadProfileMeta_NoFile_ZeroValue(t *testing.T) {
	dir := t.TempDir()
	p := LoadProfileMeta(dir)
	require.Equal(t, ProfileMeta{}, p)
}

func TestLoadProfileMeta_InvalidJSON_ZeroValue(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(
		filepath.Join(dir, "profile.lock.json"),
		[]byte("not json"),
		0644,
	))
	p := LoadProfileMeta(dir)
	require.Equal(t, ProfileMeta{}, p)
}
