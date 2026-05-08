package statusline

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func LoadProfileMeta(configDir string) ProfileMeta {
	path := filepath.Join(configDir, "profile.lock.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return ProfileMeta{}
	}
	var m ProfileMeta
	if err := json.Unmarshal(data, &m); err != nil {
		return ProfileMeta{}
	}
	return m
}
