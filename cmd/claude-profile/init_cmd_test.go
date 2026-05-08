package main

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func buildFakeClaudeForCLI(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	bin := filepath.Join(dir, "claude")
	cmd := exec.Command("go", "build", "-o", bin, "../../internal/plugin/testdata/fake_claude.go")
	cmd.Stderr = os.Stderr
	require.NoError(t, cmd.Run())
	return bin
}

func TestInitCmd_FromFile(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	bin := buildFakeClaudeForCLI(t)
	t.Setenv("PATH", filepath.Dir(bin)+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Setenv("FAKE_CLAUDE_LOG", filepath.Join(t.TempDir(), "calls.log"))

	specPath := filepath.Join(t.TempDir(), "p.yaml")
	require.NoError(t, os.WriteFile(specPath, []byte(`apiVersion: claude-profile.io/v1
kind: Profile
metadata: {name: t}
statusline: {label: t, color: red}
plugins: [a@m1]
marketplaces:
  m1: {type: github, repo: o/r}
claude_md: "rules"
`), 0644))

	cmd := newRootCmd("test")
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{"init", "t", "-f", specPath})

	require.NoError(t, cmd.Execute())

	target := filepath.Join(home, ".claude-t")
	require.DirExists(t, target)

	settingsRaw, err := os.ReadFile(filepath.Join(target, "settings.json"))
	require.NoError(t, err)
	var settings map[string]any
	require.NoError(t, json.Unmarshal(settingsRaw, &settings))
	sl := settings["statusLine"].(map[string]any)
	require.Equal(t, "claude-profile statusline", sl["command"])

	lockRaw, err := os.ReadFile(filepath.Join(target, "profile.lock.json"))
	require.NoError(t, err)
	require.Contains(t, string(lockRaw), `"label": "t"`)

	mdRaw, err := os.ReadFile(filepath.Join(target, "CLAUDE.md"))
	require.NoError(t, err)
	require.Equal(t, "rules", string(mdRaw))
}

func TestInitCmd_FromTemplate(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	bin := buildFakeClaudeForCLI(t)
	t.Setenv("PATH", filepath.Dir(bin)+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Setenv("FAKE_CLAUDE_LOG", filepath.Join(t.TempDir(), "calls.log"))

	cmd := newRootCmd("test")
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{"init", "p", "-t", "personal"})

	require.NoError(t, cmd.Execute())

	require.DirExists(t, filepath.Join(home, ".claude-p"))
}
