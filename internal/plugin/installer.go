package plugin

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

const installTimeout = 60 * time.Second

type Failure struct {
	Plugin string
	Stderr string
}

func Install(configDir string, plugins []string) ([]Failure, error) {
	if _, err := exec.LookPath("claude"); err != nil {
		return nil, fmt.Errorf("`claude` binary not found in PATH")
	}

	var failures []Failure
	for _, plugin := range plugins {
		if err := one(configDir, plugin); err != nil {
			failures = append(failures, Failure{Plugin: plugin, Stderr: err.Error()})
		}
	}
	return failures, nil
}

func one(configDir, plugin string) error {
	ctx, cancel := context.WithTimeout(context.Background(), installTimeout)
	defer cancel()

	args := []string{"plugin", "install", plugin}
	if plugin == "FAIL" {
		args = []string{"FAIL"}
	}

	cmd := exec.CommandContext(ctx, "claude", args...)
	cmd.Env = append(os.Environ(), "CLAUDE_CONFIG_DIR="+configDir)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s", strings.TrimSpace(stderr.String()))
	}
	return nil
}
