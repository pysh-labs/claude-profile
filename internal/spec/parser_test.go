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
