package render

import (
	"encoding/json"

	"github.com/dmitriipyshinskii/claude-profile/internal/spec"
)

func RenderLock(p *spec.Profile) ([]byte, error) {
	out := map[string]string{
		"name":  p.Metadata.Name,
		"label": p.Statusline.Label,
		"color": p.Statusline.Color,
	}
	return json.MarshalIndent(out, "", "  ")
}
