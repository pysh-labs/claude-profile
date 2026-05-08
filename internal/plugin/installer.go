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

type MarketplaceFailure struct {
	Source string
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

func AddMarketplaces(configDir string, sources []string) ([]MarketplaceFailure, error) {
	if _, err := exec.LookPath("claude"); err != nil {
		return nil, fmt.Errorf("`claude` binary not found in PATH")
	}

	var failures []MarketplaceFailure
	for _, src := range sources {
		if err := addOne(configDir, src); err != nil {
			failures = append(failures, MarketplaceFailure{Source: src, Stderr: err.Error()})
		}
	}
	return failures, nil
}

func addOne(configDir, source string) error {
	return runClaude(configDir, "plugin", "marketplace", "add", source)
}

func one(configDir, plugin string) error {
	return runClaude(configDir, "plugin", "install", plugin)
}

func runClaude(configDir string, args ...string) error {
	ctx, cancel := context.WithTimeout(context.Background(), installTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "claude", args...)
	cmd.Env = append(os.Environ(), "CLAUDE_CONFIG_DIR="+configDir)
	var stderr, stdout bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err == nil {
		return nil
	}
	if ctx.Err() == context.DeadlineExceeded {
		return fmt.Errorf("timed out after %s", installTimeout)
	}
	msg := strings.TrimSpace(stderr.String())
	if msg == "" {
		msg = strings.TrimSpace(stdout.String())
	}
	if msg == "" {
		return err
	}
	return fmt.Errorf("%s", msg)
}
