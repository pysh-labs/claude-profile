package templates

import (
	"embed"
	"fmt"
	"sort"
	"strings"
)

//go:embed embedded/*.yaml
var fs embed.FS

func List() []string {
	entries, _ := fs.ReadDir("embedded")
	out := make([]string, 0, len(entries))
	for _, e := range entries {
		out = append(out, strings.TrimSuffix(e.Name(), ".yaml"))
	}
	sort.Strings(out)
	return out
}

func Load(name string) ([]byte, error) {
	path := fmt.Sprintf("embedded/%s.yaml", name)
	data, err := fs.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("template %q not found", name)
	}
	return data, nil
}
