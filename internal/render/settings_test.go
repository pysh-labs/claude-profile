package render

import (
	"encoding/json"
	"testing"

	"github.com/dmitriipyshinskii/claude-profile/internal/spec"
	"github.com/stretchr/testify/require"
)

func TestRenderSettings_Minimal(t *testing.T) {
	p := &spec.Profile{
		Statusline: spec.Statusline{Label: "x", Color: "red"},
	}
	data, err := RenderSettings(p)
	require.NoError(t, err)

	var got map[string]any
	require.NoError(t, json.Unmarshal(data, &got))

	sl := got["statusLine"].(map[string]any)
	require.Equal(t, "command", sl["type"])
	require.Equal(t, "claude-profile statusline", sl["command"])
}

func TestRenderSettings_TranslatesMarketplaces(t *testing.T) {
	p := &spec.Profile{
		Statusline: spec.Statusline{Label: "x", Color: "red"},
		Marketplaces: map[string]spec.Marketplace{
			"mkt": {Type: "github", Repo: "anthropics/claude-plugins-official"},
		},
	}
	data, err := RenderSettings(p)
	require.NoError(t, err)

	var got map[string]any
	require.NoError(t, json.Unmarshal(data, &got))

	mkts := got["extraKnownMarketplaces"].(map[string]any)
	mkt := mkts["mkt"].(map[string]any)
	source := mkt["source"].(map[string]any)
	require.Equal(t, "github", source["source"])
	require.Equal(t, "anthropics/claude-plugins-official", source["repo"])
}

func TestRenderSettings_OverridesWinDeepMerge(t *testing.T) {
	p := &spec.Profile{
		Statusline: spec.Statusline{Label: "x", Color: "red"},
		SettingsOverrides: map[string]any{
			"statusLine": map[string]any{"command": "my-custom-cmd"},
			"customKey":  true,
		},
	}
	data, err := RenderSettings(p)
	require.NoError(t, err)

	var got map[string]any
	require.NoError(t, json.Unmarshal(data, &got))

	sl := got["statusLine"].(map[string]any)
	require.Equal(t, "command", sl["type"], "type kept from base")
	require.Equal(t, "my-custom-cmd", sl["command"], "command overridden")
	require.Equal(t, true, got["customKey"])
}
