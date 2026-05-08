package statusline

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

type stdinPayload struct {
	Workspace struct {
		CurrentDir string `json:"current_dir"`
	} `json:"workspace"`
	Cwd   string `json:"cwd"`
	Model struct {
		DisplayName string `json:"display_name"`
	} `json:"model"`
	ContextWindow struct {
		UsedPercentage float64 `json:"used_percentage"`
	} `json:"context_window"`
}

func ParseStdin(r io.Reader) (Input, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return Input{}, fmt.Errorf("read stdin: %w", err)
	}
	var p stdinPayload
	if err := json.Unmarshal(data, &p); err != nil {
		return Input{}, fmt.Errorf("parse json: %w", err)
	}
	cwd := p.Workspace.CurrentDir
	if cwd == "" {
		cwd = p.Cwd
	}
	return Input{
		Cwd:        cwd,
		Home:       os.Getenv("HOME"),
		Branch:     gitBranch(cwd),
		ModelName:  p.Model.DisplayName,
		ContextPct: p.ContextWindow.UsedPercentage,
	}, nil
}

func gitBranch(dir string) string {
	if dir == "" {
		return ""
	}
	cmd := exec.Command("git", "-C", dir, "-c", "core.fsmonitor=false", "symbolic-ref", "--short", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}
