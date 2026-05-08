package profile

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Profile struct {
	Name  string `json:"name"`
	Path  string `json:"path"`
	Label string `json:"label,omitempty"`
	Color string `json:"color,omitempty"`
}

func Discover(home string) ([]Profile, error) {
	entries, err := os.ReadDir(home)
	if err != nil {
		return nil, fmt.Errorf("read home: %w", err)
	}
	var out []Profile
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		name := e.Name()
		if name == ".claude" {
			out = append(out, load(home, name, "default"))
			continue
		}
		if strings.HasPrefix(name, ".claude-") {
			out = append(out, load(home, name, strings.TrimPrefix(name, ".claude-")))
		}
	}
	return out, nil
}

func load(home, dirName, profileName string) Profile {
	path := filepath.Join(home, dirName)
	p := Profile{Name: profileName, Path: path}
	data, err := os.ReadFile(filepath.Join(path, "profile.lock.json"))
	if err == nil {
		var lock struct {
			Label string `json:"label"`
			Color string `json:"color"`
		}
		if err := json.Unmarshal(data, &lock); err == nil {
			p.Label = lock.Label
			p.Color = lock.Color
		}
	}
	return p
}

func Active(home string) string {
	dir := os.Getenv("CLAUDE_CONFIG_DIR")
	if dir == "" {
		return "default"
	}
	base := filepath.Base(dir)
	if base == ".claude" {
		return "default"
	}
	return strings.TrimPrefix(base, ".claude-")
}
