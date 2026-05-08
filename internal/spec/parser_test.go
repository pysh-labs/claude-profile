package spec

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParse_Minimal(t *testing.T) {
	data, err := os.ReadFile("testdata/minimal.yaml")
	require.NoError(t, err)

	p, err := Parse(data)
	require.NoError(t, err)
	require.Equal(t, "claude-profile.io/v1", p.APIVersion)
	require.Equal(t, "Profile", p.Kind)
	require.Equal(t, "minimal", p.Metadata.Name)
	require.Equal(t, "minimal", p.Statusline.Label)
	require.Equal(t, "green", p.Statusline.Color)
}

func TestParse_Full(t *testing.T) {
	data, err := os.ReadFile("testdata/full.yaml")
	require.NoError(t, err)

	p, err := Parse(data)
	require.NoError(t, err)
	require.Equal(t, "personal", p.Metadata.Name)
	require.Len(t, p.Plugins, 2)
	require.Equal(t, "atlassian@claude-plugins-official", p.Plugins[0])
	require.Equal(t, "github", p.Marketplaces["claude-plugins-official"].Type)
	require.Equal(t, "docker", p.MCPServers["atlassian"].Command)
	require.Equal(t, "dmitri", p.MCPServers["atlassian"].Env["USER_NAME"])
	require.True(t, p.SettingsOverrides["skipDangerousModePermissionPrompt"].(bool))
	require.Contains(t, p.ClaudeMD, "feature/")
}

func TestParse_InvalidColor(t *testing.T) {
	data, err := os.ReadFile("testdata/invalid_color.yaml")
	require.NoError(t, err)

	_, err = Parse(data)
	require.Error(t, err)
	require.Contains(t, err.Error(), "color")
	require.Contains(t, err.Error(), "neon")
}

func TestParse_InvalidPluginFormat(t *testing.T) {
	data := []byte(`
apiVersion: claude-profile.io/v1
kind: Profile
metadata: {name: x}
statusline: {label: x, color: red}
plugins:
  - atlassian
  - superpowers@market
`)
	_, err := Parse(data)
	require.Error(t, err)
	require.Contains(t, err.Error(), "atlassian")
	require.Contains(t, err.Error(), "name@marketplace")
}

func TestParse_MissingName(t *testing.T) {
	data := []byte(`
apiVersion: claude-profile.io/v1
kind: Profile
metadata: {}
statusline: {label: x, color: red}
`)
	_, err := Parse(data)
	require.Error(t, err)
	require.Contains(t, err.Error(), "name")
}
