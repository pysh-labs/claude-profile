package statusline

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
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
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	cmd := exec.CommandContext(ctx, "git", "-C", dir, "-c", "core.fsmonitor=false", "symbolic-ref", "--short", "HEAD")
	cmd.Stderr = io.Discard
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0", "GIT_OPTIONAL_LOCKS=0")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}
