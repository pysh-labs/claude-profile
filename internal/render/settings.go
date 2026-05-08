package render

import (
	"encoding/json"

	"github.com/dmitriipyshinskii/claude-profile/internal/spec"
)

func RenderSettings(p *spec.Profile) ([]byte, error) {
	base := map[string]any{
		"statusLine": map[string]any{
			"type":    "command",
			"command": "claude-profile statusline",
		},
	}

	if len(p.Marketplaces) > 0 {
		mkts := make(map[string]any, len(p.Marketplaces))
		for name, m := range p.Marketplaces {
			mkts[name] = map[string]any{
				"source": map[string]any{
					"source": m.Type,
					"repo":   m.Repo,
				},
			}
		}
		base["extraKnownMarketplaces"] = mkts
	}

	merged := DeepMerge(base, p.SettingsOverrides)
	return json.MarshalIndent(merged, "", "  ")
}
