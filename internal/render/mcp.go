package render

import (
	"encoding/json"

	"github.com/dmitriipyshinskii/claude-profile/internal/spec"
)

func RenderMCP(p *spec.Profile) ([]byte, bool, error) {
	if len(p.MCPServers) == 0 {
		return nil, false, nil
	}
	servers := make(map[string]any, len(p.MCPServers))
	for name, s := range p.MCPServers {
		entry := map[string]any{"command": s.Command}
		if len(s.Args) > 0 {
			entry["args"] = s.Args
		}
		if len(s.Env) > 0 {
			entry["env"] = s.Env
		}
		servers[name] = entry
	}
	out := map[string]any{"mcpServers": servers}
	data, err := json.MarshalIndent(out, "", "  ")
	return data, true, err
}
