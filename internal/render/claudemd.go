package render

import "github.com/dmitriipyshinskii/claude-profile/internal/spec"

func RenderClaudeMD(p *spec.Profile) ([]byte, bool) {
	if p.ClaudeMD == "" {
		return nil, false
	}
	return []byte(p.ClaudeMD), true
}
