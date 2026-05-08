package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStatuslineCmd_RendersFromLock(t *testing.T) {
	configDir := t.TempDir()
	require.NoError(t, os.WriteFile(
		filepath.Join(configDir, "profile.lock.json"),
		[]byte(`{"name":"work","label":"work","color":"red"}`),
		0644,
	))
	t.Setenv("CLAUDE_CONFIG_DIR", configDir)

	cmd := newRootCmd("test")
	cmd.SetIn(strings.NewReader(`{"workspace":{"current_dir":"/x"},"model":{"display_name":"M"}}`))
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"statusline"})

	require.NoError(t, cmd.Execute())
	require.Contains(t, out.String(), "[work]")
	require.Contains(t, out.String(), "M")
}
