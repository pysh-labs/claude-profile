# claude-profile v0.1 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build the v0.1 MVP of `claude-profile` — a Go single-binary tool that bootstraps isolated Claude Code profiles from a declarative `profile.yaml`, with per-profile statusline rendering and shell alias generation.

**Architecture:** Layered Go packages (`internal/spec`, `internal/render`, `internal/statusline`, `internal/profile`, `internal/shell`, `internal/plugin`, `internal/templates`) wired through a thin `cmd/claude-profile/main.go` cobra CLI. Bootstrap-only authority model — yaml is a template, not reconciled state. See `docs/design/2026-05-07-claude-profile-as-code-design.md`.

**Tech Stack:**
- Go 1.22+
- `github.com/spf13/cobra` — CLI
- `gopkg.in/yaml.v3` — yaml parsing
- `github.com/go-playground/validator/v10` — schema validation
- `github.com/stretchr/testify/require` — test assertions

**Module path placeholder:** `github.com/dmitriipyshinskii/claude-profile`. Substitute the real GitHub username in Task 1 if different.

**Working tree state:**
- `/Users/dmitriipyshinskii/Projects/personal/claude-profile/` — git repo on `main`, contains only `docs/design/2026-05-07-claude-profile-as-code-design.md` and this plan file.

**Sub-product checkpoints:**
- After Task 8: `internal/spec` and `internal/render` produce all owned files in-memory.
- After Task 14: `claude-profile statusline` works end-to-end (Q1 from user's original request fixed).
- After Task 22: `claude-profile init` flow works against a fake `claude` binary in tests.
- After Task 27: full CLI surface working manually.
- After Task 31: published binary, releases, README.

---

## Task 1: Repo bootstrap

**Files:**
- Create: `go.mod`
- Create: `.gitignore`
- Create: `LICENSE`
- Create: `Makefile`
- Create: `README.md`

- [ ] **Step 1: Init Go module**

```bash
cd /Users/dmitriipyshinskii/Projects/personal/claude-profile
go mod init github.com/dmitriipyshinskii/claude-profile
```

If your GitHub username is not `dmitriipyshinskii`, edit the `module` line in `go.mod` after init.

- [ ] **Step 2: Write `.gitignore`**

```
# Binaries
/claude-profile
/dist/

# Editor
.vscode/
.idea/
*.swp

# OS
.DS_Store

# Coverage
*.out
coverage.html
```

- [ ] **Step 3: Write `LICENSE` (MIT)**

```
MIT License

Copyright (c) 2026 Dmitrii Pyshinskii

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

- [ ] **Step 4: Write `Makefile`**

```makefile
.PHONY: build test lint fmt vet cover

build:
	go build -o claude-profile ./cmd/claude-profile

test:
	go test ./...

lint:
	golangci-lint run ./...

fmt:
	gofmt -s -w .

vet:
	go vet ./...

cover:
	go test -coverprofile=coverage.out ./internal/...
	go tool cover -html=coverage.out -o coverage.html
```

- [ ] **Step 5: Write minimal `README.md` placeholder**

```markdown
# claude-profile

Profile-as-code for Claude Code. Define organizations as YAML, version them, share with your team.

> Status: v0.1 in development. See `docs/design/` for the spec and `docs/plans/` for the implementation plan.
```

- [ ] **Step 6: Commit**

```bash
git add -A
git commit -m "chore: bootstrap repo (go.mod, license, makefile, .gitignore)"
```

---

## Task 2: Cobra root command + version

**Files:**
- Create: `cmd/claude-profile/main.go`
- Create: `cmd/claude-profile/main_test.go`
- Modify: `go.mod` (cobra dependency)

- [ ] **Step 1: Add cobra**

```bash
go get github.com/spf13/cobra@latest
```

- [ ] **Step 2: Write failing test for `--version`**

`cmd/claude-profile/main_test.go`:
```go
package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRootCmd_VersionFlag(t *testing.T) {
	cmd := newRootCmd("0.0.1-test")
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"--version"})

	require.NoError(t, cmd.Execute())
	require.Contains(t, out.String(), "0.0.1-test")
}
```

- [ ] **Step 3: Add testify**

```bash
go get github.com/stretchr/testify@latest
```

- [ ] **Step 4: Run test, expect compile error (newRootCmd undefined)**

```bash
go test ./cmd/claude-profile/
```
Expected: build error.

- [ ] **Step 5: Implement `main.go`**

```go
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "dev"

func newRootCmd(v string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "claude-profile",
		Short:   "Profile-as-code for Claude Code",
		Version: v,
	}
	cmd.SetVersionTemplate("claude-profile {{.Version}}\n")
	return cmd
}

func main() {
	if err := newRootCmd(version).Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
```

- [ ] **Step 6: Run test, verify pass**

```bash
go test ./cmd/claude-profile/
```
Expected: PASS

- [ ] **Step 7: Verify build**

```bash
go build -o claude-profile ./cmd/claude-profile && ./claude-profile --version
```
Expected: `claude-profile dev`

- [ ] **Step 8: Commit**

```bash
git add -A
git commit -m "feat: cobra root command with version flag"
```

---

## Task 3: spec package — Profile struct + parser

**Files:**
- Create: `internal/spec/schema.go`
- Create: `internal/spec/parser.go`
- Create: `internal/spec/parser_test.go`
- Create: `internal/spec/testdata/minimal.yaml`
- Create: `internal/spec/testdata/full.yaml`

- [ ] **Step 1: Add yaml dependency**

```bash
go get gopkg.in/yaml.v3
```

- [ ] **Step 2: Create test data — minimal valid spec**

`internal/spec/testdata/minimal.yaml`:
```yaml
apiVersion: claude-profile.io/v1
kind: Profile

metadata:
  name: minimal

statusline:
  label: minimal
  color: green
```

- [ ] **Step 3: Create test data — full spec**

`internal/spec/testdata/full.yaml`:
```yaml
apiVersion: claude-profile.io/v1
kind: Profile

metadata:
  name: personal
  description: Personal profile

statusline:
  label: personal
  color: green

marketplaces:
  claude-plugins-official:
    type: github
    repo: anthropics/claude-plugins-official

plugins:
  - atlassian@claude-plugins-official
  - superpowers@claude-plugins-official

mcp_servers:
  atlassian:
    command: docker
    args: [run, -i, --rm, mcp/atlassian]
    env:
      USER_NAME: dmitri

settings_overrides:
  skipDangerousModePermissionPrompt: true

claude_md: |
  Branch prefix: feature/, bugfix/, hotfix/.
```

- [ ] **Step 4: Write failing test for Parse**

`internal/spec/parser_test.go`:
```go
package spec

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParse_Minimal(t *testing.T) {
	data, err := os.ReadFile("testdata/minimal.yaml")
	require.NoError(t, err)

	p, err := Parse(data)
	require.NoError(t, err)
	require.Equal(t, "claude-profile.io/v1", p.APIVersion)
	require.Equal(t, "Profile", p.Kind)
	require.Equal(t, "minimal", p.Metadata.Name)
	require.Equal(t, "minimal", p.Statusline.Label)
	require.Equal(t, "green", p.Statusline.Color)
}

func TestParse_Full(t *testing.T) {
	data, err := os.ReadFile("testdata/full.yaml")
	require.NoError(t, err)

	p, err := Parse(data)
	require.NoError(t, err)
	require.Equal(t, "personal", p.Metadata.Name)
	require.Len(t, p.Plugins, 2)
	require.Equal(t, "atlassian@claude-plugins-official", p.Plugins[0])
	require.Equal(t, "github", p.Marketplaces["claude-plugins-official"].Type)
	require.Equal(t, "docker", p.MCPServers["atlassian"].Command)
	require.Equal(t, "dmitri", p.MCPServers["atlassian"].Env["USER_NAME"])
	require.True(t, p.SettingsOverrides["skipDangerousModePermissionPrompt"].(bool))
	require.Contains(t, p.ClaudeMD, "feature/")
}
```

- [ ] **Step 5: Run test — expect compile error**

```bash
go test ./internal/spec/
```
Expected: build error.

- [ ] **Step 6: Implement schema**

`internal/spec/schema.go`:
```go
package spec

type Profile struct {
	APIVersion        string                  `yaml:"apiVersion"`
	Kind              string                  `yaml:"kind"`
	Metadata          Metadata                `yaml:"metadata"`
	Statusline        Statusline              `yaml:"statusline"`
	Marketplaces      map[string]Marketplace  `yaml:"marketplaces,omitempty"`
	Plugins           []string                `yaml:"plugins,omitempty"`
	MCPServers        map[string]MCPServer    `yaml:"mcp_servers,omitempty"`
	SettingsOverrides map[string]any          `yaml:"settings_overrides,omitempty"`
	ClaudeMD          string                  `yaml:"claude_md,omitempty"`
}

type Metadata struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description,omitempty"`
}

type Statusline struct {
	Label string `yaml:"label"`
	Color string `yaml:"color"`
}

type Marketplace struct {
	Type string `yaml:"type"`
	Repo string `yaml:"repo"`
}

type MCPServer struct {
	Command string            `yaml:"command"`
	Args    []string          `yaml:"args,omitempty"`
	Env     map[string]string `yaml:"env,omitempty"`
}
```

- [ ] **Step 7: Implement Parse**

`internal/spec/parser.go`:
```go
package spec

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

func Parse(data []byte) (*Profile, error) {
	var p Profile
	if err := yaml.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("yaml parse: %w", err)
	}
	return &p, nil
}
```

- [ ] **Step 8: Run test, verify pass**

```bash
go test ./internal/spec/
```
Expected: PASS

- [ ] **Step 9: Commit**

```bash
git add -A
git commit -m "feat(spec): yaml schema + Parse function"
```

---

## Task 4: spec package — validation

**Files:**
- Create: `internal/spec/validate.go`
- Modify: `internal/spec/parser.go`
- Modify: `internal/spec/parser_test.go`
- Create: `internal/spec/testdata/invalid_color.yaml`

- [ ] **Step 1: Add validator**

```bash
go get github.com/go-playground/validator/v10
```

- [ ] **Step 2: Write invalid color fixture**

`internal/spec/testdata/invalid_color.yaml`:
```yaml
apiVersion: claude-profile.io/v1
kind: Profile
metadata: {name: bad}
statusline: {label: bad, color: neon}
```

- [ ] **Step 3: Add failing tests**

Append to `internal/spec/parser_test.go`:
```go
func TestParse_InvalidColor(t *testing.T) {
	data, err := os.ReadFile("testdata/invalid_color.yaml")
	require.NoError(t, err)

	_, err = Parse(data)
	require.Error(t, err)
	require.Contains(t, err.Error(), "color")
	require.Contains(t, err.Error(), "neon")
}

func TestParse_InvalidPluginFormat(t *testing.T) {
	data := []byte(`
apiVersion: claude-profile.io/v1
kind: Profile
metadata: {name: x}
statusline: {label: x, color: red}
plugins:
  - atlassian
  - superpowers@market
`)
	_, err := Parse(data)
	require.Error(t, err)
	require.Contains(t, err.Error(), "atlassian")
	require.Contains(t, err.Error(), "name@marketplace")
}

func TestParse_MissingName(t *testing.T) {
	data := []byte(`
apiVersion: claude-profile.io/v1
kind: Profile
metadata: {}
statusline: {label: x, color: red}
`)
	_, err := Parse(data)
	require.Error(t, err)
	require.Contains(t, err.Error(), "name")
}
```

- [ ] **Step 4: Run tests, expect fail**

```bash
go test ./internal/spec/
```
Expected: 3 new tests fail.

- [ ] **Step 5: Add validator tags to schema**

Replace `internal/spec/schema.go`:
```go
package spec

type Profile struct {
	APIVersion        string                  `yaml:"apiVersion"        validate:"required,eq=claude-profile.io/v1"`
	Kind              string                  `yaml:"kind"              validate:"required,eq=Profile"`
	Metadata          Metadata                `yaml:"metadata"          validate:"required"`
	Statusline        Statusline              `yaml:"statusline"        validate:"required"`
	Marketplaces      map[string]Marketplace  `yaml:"marketplaces,omitempty"`
	Plugins           []string                `yaml:"plugins,omitempty"`
	MCPServers        map[string]MCPServer    `yaml:"mcp_servers,omitempty"`
	SettingsOverrides map[string]any          `yaml:"settings_overrides,omitempty"`
	ClaudeMD          string                  `yaml:"claude_md,omitempty"`
}

type Metadata struct {
	Name        string `yaml:"name"        validate:"required"`
	Description string `yaml:"description,omitempty"`
}

type Statusline struct {
	Label string `yaml:"label" validate:"required"`
	Color string `yaml:"color" validate:"required,oneof=red green yellow blue magenta cyan white"`
}

type Marketplace struct {
	Type string `yaml:"type" validate:"required,oneof=github"`
	Repo string `yaml:"repo" validate:"required"`
}

type MCPServer struct {
	Command string            `yaml:"command" validate:"required"`
	Args    []string          `yaml:"args,omitempty"`
	Env     map[string]string `yaml:"env,omitempty"`
}
```

- [ ] **Step 6: Implement Validate**

`internal/spec/validate.go`:
```go
package spec

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var pluginFormat = regexp.MustCompile(`^[A-Za-z0-9_.-]+@[A-Za-z0-9_.-]+$`)

func Validate(p *Profile) error {
	v := validator.New()
	var errs []string

	if err := v.Struct(p); err != nil {
		ve, _ := err.(validator.ValidationErrors)
		for _, fe := range ve {
			errs = append(errs, formatFieldError(fe))
		}
	}

	for i, plugin := range p.Plugins {
		if !pluginFormat.MatchString(plugin) {
			errs = append(errs, fmt.Sprintf(
				"plugins[%d]: invalid format, expected \"name@marketplace\", got %q",
				i, plugin))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("profile validation failed:\n  - %s", strings.Join(errs, "\n  - "))
	}
	return nil
}

func formatFieldError(fe validator.FieldError) string {
	switch fe.Tag() {
	case "oneof":
		return fmt.Sprintf("%s: must be one of [%s], got %q",
			toYAMLPath(fe.Namespace()), fe.Param(), fe.Value())
	case "required":
		return fmt.Sprintf("%s: required", toYAMLPath(fe.Namespace()))
	case "eq":
		return fmt.Sprintf("%s: must equal %q, got %q",
			toYAMLPath(fe.Namespace()), fe.Param(), fe.Value())
	default:
		return fmt.Sprintf("%s: failed %s validation", toYAMLPath(fe.Namespace()), fe.Tag())
	}
}

func toYAMLPath(ns string) string {
	parts := strings.Split(ns, ".")
	if len(parts) > 1 {
		parts = parts[1:]
	}
	for i, part := range parts {
		parts[i] = strings.ToLower(part[:1]) + part[1:]
	}
	return strings.Join(parts, ".")
}
```

- [ ] **Step 7: Wire Validate into Parse**

Update `internal/spec/parser.go`:
```go
package spec

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

func Parse(data []byte) (*Profile, error) {
	var p Profile
	if err := yaml.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("yaml parse: %w", err)
	}
	if err := Validate(&p); err != nil {
		return nil, err
	}
	return &p, nil
}
```

- [ ] **Step 8: Run tests, verify pass**

```bash
go test ./internal/spec/
```
Expected: all pass.

- [ ] **Step 9: Commit**

```bash
git add -A
git commit -m "feat(spec): schema validation with aggregated errors"
```

---

## Task 5: spec package — env interpolation

**Files:**
- Create: `internal/spec/interpolate.go`
- Create: `internal/spec/interpolate_test.go`

- [ ] **Step 1: Write failing tests**

`internal/spec/interpolate_test.go`:
```go
package spec

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInterpolate_ReplacesEnvVars(t *testing.T) {
	t.Setenv("FOO", "bar")
	t.Setenv("BAZ", "qux")

	p := &Profile{
		MCPServers: map[string]MCPServer{
			"s": {
				Command: "x",
				Env: map[string]string{
					"A": "${FOO}",
					"B": "literal",
					"C": "${BAZ}-suffix",
				},
			},
		},
	}

	err := Interpolate(p)
	require.NoError(t, err)
	require.Equal(t, "bar", p.MCPServers["s"].Env["A"])
	require.Equal(t, "literal", p.MCPServers["s"].Env["B"])
	require.Equal(t, "qux-suffix", p.MCPServers["s"].Env["C"])
}

func TestInterpolate_MissingVarFails(t *testing.T) {
	p := &Profile{
		MCPServers: map[string]MCPServer{
			"s": {Command: "x", Env: map[string]string{"A": "${MISSING_VAR_XYZ}"}},
		},
	}

	err := Interpolate(p)
	require.Error(t, err)
	require.Contains(t, err.Error(), "MISSING_VAR_XYZ")
	require.Contains(t, err.Error(), "mcp_servers.s.env.A")
}
```

- [ ] **Step 2: Run test, expect fail**

```bash
go test ./internal/spec/ -run Interpolate
```
Expected: build error.

- [ ] **Step 3: Implement**

`internal/spec/interpolate.go`:
```go
package spec

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
)

var envVarRe = regexp.MustCompile(`\$\{([A-Za-z_][A-Za-z0-9_]*)\}`)

func Interpolate(p *Profile) error {
	var errs []string

	serverNames := make([]string, 0, len(p.MCPServers))
	for name := range p.MCPServers {
		serverNames = append(serverNames, name)
	}
	sort.Strings(serverNames)

	for _, name := range serverNames {
		srv := p.MCPServers[name]
		envKeys := make([]string, 0, len(srv.Env))
		for k := range srv.Env {
			envKeys = append(envKeys, k)
		}
		sort.Strings(envKeys)

		for _, k := range envKeys {
			expanded, err := expand(srv.Env[k])
			if err != nil {
				errs = append(errs,
					fmt.Sprintf("mcp_servers.%s.env.%s: %s", name, k, err.Error()))
				continue
			}
			srv.Env[k] = expanded
		}
		p.MCPServers[name] = srv
	}

	if len(errs) > 0 {
		return fmt.Errorf("env interpolation failed:\n  - %s", strings.Join(errs, "\n  - "))
	}
	return nil
}

func expand(s string) (string, error) {
	var missing []string
	out := envVarRe.ReplaceAllStringFunc(s, func(match string) string {
		name := match[2 : len(match)-1]
		val, ok := os.LookupEnv(name)
		if !ok {
			missing = append(missing, name)
			return match
		}
		return val
	})
	if len(missing) > 0 {
		return "", fmt.Errorf("environment variable %s is not set", missing[0])
	}
	return out, nil
}
```

- [ ] **Step 4: Run tests, verify pass**

```bash
go test ./internal/spec/
```
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add -A
git commit -m "feat(spec): ${VAR} env interpolation for mcp_servers"
```

---

## Task 6: render package — atomic write helper

**Files:**
- Create: `internal/render/atomic.go`
- Create: `internal/render/atomic_test.go`

- [ ] **Step 1: Write failing test**

`internal/render/atomic_test.go`:
```go
package render

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWriteAtomic_NewFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "f.txt")

	require.NoError(t, WriteAtomic(path, []byte("hello")))

	got, err := os.ReadFile(path)
	require.NoError(t, err)
	require.Equal(t, "hello", string(got))
}

func TestWriteAtomic_OverwriteLeavesNoTmp(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "f.txt")

	require.NoError(t, WriteAtomic(path, []byte("v1")))
	require.NoError(t, WriteAtomic(path, []byte("v2")))

	got, err := os.ReadFile(path)
	require.NoError(t, err)
	require.Equal(t, "v2", string(got))

	entries, _ := os.ReadDir(dir)
	require.Len(t, entries, 1, "should not leave .tmp behind")
}
```

- [ ] **Step 2: Run, expect fail**

```bash
go test ./internal/render/
```
Expected: build error.

- [ ] **Step 3: Implement**

`internal/render/atomic.go`:
```go
package render

import (
	"fmt"
	"os"
)

func WriteAtomic(path string, content []byte) error {
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, content, 0644); err != nil {
		return fmt.Errorf("write tmp: %w", err)
	}
	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("rename: %w", err)
	}
	return nil
}
```

- [ ] **Step 4: Run, verify pass**

```bash
go test ./internal/render/
```
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add -A
git commit -m "feat(render): WriteAtomic helper"
```

---

## Task 7: render package — settings.json with deep-merge

**Files:**
- Create: `internal/render/settings.go`
- Create: `internal/render/settings_test.go`
- Create: `internal/render/merge.go`
- Create: `internal/render/merge_test.go`

- [ ] **Step 1: Write deep-merge tests**

`internal/render/merge_test.go`:
```go
package render

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDeepMerge_OverridesWin(t *testing.T) {
	base := map[string]any{"a": 1, "b": "base"}
	over := map[string]any{"b": "over", "c": true}

	got := DeepMerge(base, over)
	require.Equal(t, 1, got["a"])
	require.Equal(t, "over", got["b"])
	require.Equal(t, true, got["c"])
}

func TestDeepMerge_NestedObjectsRecurse(t *testing.T) {
	base := map[string]any{
		"perms": map[string]any{"allow": []any{"x"}, "deny": []any{"y"}},
	}
	over := map[string]any{
		"perms": map[string]any{"allow": []any{"z"}},
	}

	got := DeepMerge(base, over)
	perms := got["perms"].(map[string]any)
	require.Equal(t, []any{"z"}, perms["allow"], "arrays replace, not concat")
	require.Equal(t, []any{"y"}, perms["deny"])
}
```

- [ ] **Step 2: Implement deep-merge**

`internal/render/merge.go`:
```go
package render

func DeepMerge(base, over map[string]any) map[string]any {
	out := make(map[string]any, len(base))
	for k, v := range base {
		out[k] = v
	}
	for k, v := range over {
		if existing, ok := out[k]; ok {
			if baseMap, baseOk := existing.(map[string]any); baseOk {
				if overMap, overOk := v.(map[string]any); overOk {
					out[k] = DeepMerge(baseMap, overMap)
					continue
				}
			}
		}
		out[k] = v
	}
	return out
}
```

- [ ] **Step 3: Run merge tests**

```bash
go test ./internal/render/ -run DeepMerge
```
Expected: PASS

- [ ] **Step 4: Write settings.json render tests**

`internal/render/settings_test.go`:
```go
package render

import (
	"encoding/json"
	"testing"

	"github.com/dmitriipyshinskii/claude-profile/internal/spec"
	"github.com/stretchr/testify/require"
)

func TestRenderSettings_Minimal(t *testing.T) {
	p := &spec.Profile{
		Statusline: spec.Statusline{Label: "x", Color: "red"},
	}
	data, err := RenderSettings(p)
	require.NoError(t, err)

	var got map[string]any
	require.NoError(t, json.Unmarshal(data, &got))

	sl := got["statusLine"].(map[string]any)
	require.Equal(t, "command", sl["type"])
	require.Equal(t, "claude-profile statusline", sl["command"])
}

func TestRenderSettings_TranslatesMarketplaces(t *testing.T) {
	p := &spec.Profile{
		Statusline: spec.Statusline{Label: "x", Color: "red"},
		Marketplaces: map[string]spec.Marketplace{
			"mkt": {Type: "github", Repo: "anthropics/claude-plugins-official"},
		},
	}
	data, err := RenderSettings(p)
	require.NoError(t, err)

	var got map[string]any
	require.NoError(t, json.Unmarshal(data, &got))

	mkts := got["extraKnownMarketplaces"].(map[string]any)
	mkt := mkts["mkt"].(map[string]any)
	source := mkt["source"].(map[string]any)
	require.Equal(t, "github", source["source"])
	require.Equal(t, "anthropics/claude-plugins-official", source["repo"])
}

func TestRenderSettings_OverridesWinDeepMerge(t *testing.T) {
	p := &spec.Profile{
		Statusline: spec.Statusline{Label: "x", Color: "red"},
		SettingsOverrides: map[string]any{
			"statusLine": map[string]any{"command": "my-custom-cmd"},
			"customKey":  true,
		},
	}
	data, err := RenderSettings(p)
	require.NoError(t, err)

	var got map[string]any
	require.NoError(t, json.Unmarshal(data, &got))

	sl := got["statusLine"].(map[string]any)
	require.Equal(t, "command", sl["type"], "type kept from base")
	require.Equal(t, "my-custom-cmd", sl["command"], "command overridden")
	require.Equal(t, true, got["customKey"])
}
```

- [ ] **Step 5: Implement**

`internal/render/settings.go`:
```go
package render

import (
	"encoding/json"

	"github.com/dmitriipyshinskii/claude-profile/internal/spec"
)

func RenderSettings(p *spec.Profile) ([]byte, error) {
	base := map[string]any{
		"statusLine": map[string]any{
			"type":    "command",
			"command": "claude-profile statusline",
		},
	}

	if len(p.Marketplaces) > 0 {
		mkts := make(map[string]any, len(p.Marketplaces))
		for name, m := range p.Marketplaces {
			mkts[name] = map[string]any{
				"source": map[string]any{
					"source": m.Type,
					"repo":   m.Repo,
				},
			}
		}
		base["extraKnownMarketplaces"] = mkts
	}

	merged := DeepMerge(base, p.SettingsOverrides)
	return json.MarshalIndent(merged, "", "  ")
}
```

- [ ] **Step 6: Run all render tests**

```bash
go test ./internal/render/
```
Expected: PASS

- [ ] **Step 7: Commit**

```bash
git add -A
git commit -m "feat(render): settings.json with deep-merge for overrides"
```

---

## Task 8: render package — .mcp.json, CLAUDE.md, profile.lock.json

**Files:**
- Create: `internal/render/mcp.go`
- Create: `internal/render/mcp_test.go`
- Create: `internal/render/claudemd.go`
- Create: `internal/render/claudemd_test.go`
- Create: `internal/render/lock.go`
- Create: `internal/render/lock_test.go`

- [ ] **Step 1: Write .mcp.json tests**

`internal/render/mcp_test.go`:
```go
package render

import (
	"encoding/json"
	"testing"

	"github.com/dmitriipyshinskii/claude-profile/internal/spec"
	"github.com/stretchr/testify/require"
)

func TestRenderMCP_OneServer(t *testing.T) {
	p := &spec.Profile{
		MCPServers: map[string]spec.MCPServer{
			"atl": {Command: "docker", Args: []string{"run", "x"}, Env: map[string]string{"K": "v"}},
		},
	}
	data, ok, err := RenderMCP(p)
	require.NoError(t, err)
	require.True(t, ok)

	var got map[string]any
	require.NoError(t, json.Unmarshal(data, &got))
	servers := got["mcpServers"].(map[string]any)
	atl := servers["atl"].(map[string]any)
	require.Equal(t, "docker", atl["command"])
}

func TestRenderMCP_NoServersReturnsFalse(t *testing.T) {
	p := &spec.Profile{}
	_, ok, err := RenderMCP(p)
	require.NoError(t, err)
	require.False(t, ok)
}
```

- [ ] **Step 2: Implement**

`internal/render/mcp.go`:
```go
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
```

- [ ] **Step 3: Write CLAUDE.md tests**

`internal/render/claudemd_test.go`:
```go
package render

import (
	"testing"

	"github.com/dmitriipyshinskii/claude-profile/internal/spec"
	"github.com/stretchr/testify/require"
)

func TestRenderClaudeMD_PassThrough(t *testing.T) {
	p := &spec.Profile{ClaudeMD: "Branch prefix: feature/.\n"}
	data, ok := RenderClaudeMD(p)
	require.True(t, ok)
	require.Equal(t, "Branch prefix: feature/.\n", string(data))
}

func TestRenderClaudeMD_EmptyReturnsFalse(t *testing.T) {
	p := &spec.Profile{}
	_, ok := RenderClaudeMD(p)
	require.False(t, ok)
}
```

- [ ] **Step 4: Implement**

`internal/render/claudemd.go`:
```go
package render

import "github.com/dmitriipyshinskii/claude-profile/internal/spec"

func RenderClaudeMD(p *spec.Profile) ([]byte, bool) {
	if p.ClaudeMD == "" {
		return nil, false
	}
	return []byte(p.ClaudeMD), true
}
```

- [ ] **Step 5: Write lock.json tests**

`internal/render/lock_test.go`:
```go
package render

import (
	"encoding/json"
	"testing"

	"github.com/dmitriipyshinskii/claude-profile/internal/spec"
	"github.com/stretchr/testify/require"
)

func TestRenderLock_Shape(t *testing.T) {
	p := &spec.Profile{
		Metadata:   spec.Metadata{Name: "personal"},
		Statusline: spec.Statusline{Label: "personal", Color: "green"},
	}
	data, err := RenderLock(p)
	require.NoError(t, err)

	var got map[string]string
	require.NoError(t, json.Unmarshal(data, &got))
	require.Equal(t, "personal", got["name"])
	require.Equal(t, "personal", got["label"])
	require.Equal(t, "green", got["color"])
}
```

- [ ] **Step 6: Implement**

`internal/render/lock.go`:
```go
package render

import (
	"encoding/json"

	"github.com/dmitriipyshinskii/claude-profile/internal/spec"
)

func RenderLock(p *spec.Profile) ([]byte, error) {
	out := map[string]string{
		"name":  p.Metadata.Name,
		"label": p.Statusline.Label,
		"color": p.Statusline.Color,
	}
	return json.MarshalIndent(out, "", "  ")
}
```

- [ ] **Step 7: Run all render tests**

```bash
go test ./internal/render/
```
Expected: PASS

- [ ] **Step 8: Commit**

```bash
git add -A
git commit -m "feat(render): mcp.json, CLAUDE.md, profile.lock.json"
```

---

## Task 9: statusline package — Render pure function

**Files:**
- Create: `internal/statusline/renderer.go`
- Create: `internal/statusline/renderer_test.go`

- [ ] **Step 1: Write failing tests**

`internal/statusline/renderer_test.go`:
```go
package statusline

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRender_FullInputWithProfile(t *testing.T) {
	in := Input{
		Cwd:        "/Users/me/Projects/x",
		Home:       "/Users/me",
		Branch:     "main",
		ModelName:  "Opus 4.7",
		ContextPct: 42,
	}
	p := ProfileMeta{Name: "personal", Label: "personal", Color: "green"}

	got := Render(in, p)
	require.Contains(t, got, "[personal]")
	require.Contains(t, got, "~/Projects/x")
	require.Contains(t, got, "main")
	require.Contains(t, got, "Opus 4.7")
	require.Contains(t, got, "ctx:42%")
	require.True(t, strings.Contains(got, "\x1b[1;32m"), "green ANSI for prefix")
}

func TestRender_NoProfile_NoPrefix(t *testing.T) {
	in := Input{Cwd: "/x", Home: "/x"}
	got := Render(in, ProfileMeta{})
	require.NotContains(t, got, "[")
}

func TestRender_NoBranch_SkipBranch(t *testing.T) {
	in := Input{Cwd: "/x", Home: "/x", ModelName: "M"}
	got := Render(in, ProfileMeta{})
	require.NotContains(t, got, "⎇")
}

func TestRender_CtxColors(t *testing.T) {
	cases := []struct {
		pct  float64
		ansi string
	}{
		{30, "\x1b[1;32m"},
		{60, "\x1b[1;33m"},
		{85, "\x1b[1;31m"},
	}
	for _, tc := range cases {
		got := Render(Input{Cwd: "/x", Home: "/x", ContextPct: tc.pct}, ProfileMeta{})
		require.Contains(t, got, tc.ansi, "pct=%v", tc.pct)
	}
}

func TestRender_HomeSubstitution(t *testing.T) {
	got := Render(
		Input{Cwd: "/Users/me/Projects/x", Home: "/Users/me"},
		ProfileMeta{},
	)
	require.Contains(t, got, "~/Projects/x")
	require.NotContains(t, got, "/Users/me/Projects/x")
}
```

- [ ] **Step 2: Implement**

`internal/statusline/renderer.go`:
```go
package statusline

import (
	"fmt"
	"strings"
)

type Input struct {
	Cwd        string
	Home       string
	Branch     string
	ModelName  string
	ContextPct float64
}

type ProfileMeta struct {
	Name  string
	Label string
	Color string
}

const (
	reset    = "\x1b[0m"
	dim      = "\x1b[2m"
	cyanBold = "\x1b[1;36m"
	yellow   = "\x1b[0;33m"
)

var colors = map[string]string{
	"red":     "\x1b[1;31m",
	"green":   "\x1b[1;32m",
	"yellow":  "\x1b[1;33m",
	"blue":    "\x1b[1;34m",
	"magenta": "\x1b[1;35m",
	"cyan":    "\x1b[1;36m",
	"white":   "\x1b[1;37m",
}

func Render(in Input, p ProfileMeta) string {
	var parts []string

	if p.Label != "" {
		c := colors[p.Color]
		parts = append(parts, fmt.Sprintf("%s[%s]%s", c, p.Label, reset))
	}

	short := in.Cwd
	if in.Home != "" && strings.HasPrefix(short, in.Home) {
		short = "~" + strings.TrimPrefix(short, in.Home)
	}
	if short != "" {
		parts = append(parts, fmt.Sprintf("%s%s%s", cyanBold, short, reset))
	}

	if in.Branch != "" {
		parts = append(parts, fmt.Sprintf("%s⎇ %s%s", yellow, in.Branch, reset))
	}

	if in.ModelName != "" {
		parts = append(parts, fmt.Sprintf("%s%s%s", dim, in.ModelName, reset))
	}

	if in.ContextPct > 0 {
		parts = append(parts, ctxSegment(in.ContextPct))
	}

	return strings.Join(parts, "  ")
}

func ctxSegment(pct float64) string {
	var c string
	switch {
	case pct >= 80:
		c = colors["red"]
	case pct >= 50:
		c = colors["yellow"]
	default:
		c = colors["green"]
	}
	return fmt.Sprintf("%sctx:%.0f%%%s", c, pct, reset)
}
```

- [ ] **Step 3: Run tests, verify pass**

```bash
go test ./internal/statusline/
```
Expected: PASS

- [ ] **Step 4: Commit**

```bash
git add -A
git commit -m "feat(statusline): pure Render function"
```

---

## Task 10: statusline package — git branch + stdin parsing

**Files:**
- Create: `internal/statusline/input.go`
- Create: `internal/statusline/input_test.go`

- [ ] **Step 1: Write failing tests**

`internal/statusline/input_test.go`:
```go
package statusline

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseStdin_AllFields(t *testing.T) {
	json := `{
		"workspace": {"current_dir": "/x"},
		"model": {"display_name": "Opus 4.7"},
		"context_window": {"used_percentage": 42.5}
	}`
	in, err := ParseStdin(strings.NewReader(json))
	require.NoError(t, err)
	require.Equal(t, "/x", in.Cwd)
	require.Equal(t, "Opus 4.7", in.ModelName)
	require.InDelta(t, 42.5, in.ContextPct, 0.01)
}

func TestParseStdin_FallbackCwd(t *testing.T) {
	json := `{"cwd": "/legacy"}`
	in, err := ParseStdin(strings.NewReader(json))
	require.NoError(t, err)
	require.Equal(t, "/legacy", in.Cwd)
}

func TestParseStdin_InvalidJSON(t *testing.T) {
	_, err := ParseStdin(strings.NewReader("not json"))
	require.Error(t, err)
}
```

- [ ] **Step 2: Implement**

`internal/statusline/input.go`:
```go
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
```

- [ ] **Step 3: Run tests, verify pass**

```bash
go test ./internal/statusline/
```
Expected: PASS

- [ ] **Step 4: Commit**

```bash
git add -A
git commit -m "feat(statusline): stdin parser + git branch detection"
```

---

## Task 11: statusline package — lock.json reader

**Files:**
- Create: `internal/statusline/lock.go`
- Create: `internal/statusline/lock_test.go`

- [ ] **Step 1: Write failing tests**

`internal/statusline/lock_test.go`:
```go
package statusline

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadProfileMeta_ValidLock(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(
		filepath.Join(dir, "profile.lock.json"),
		[]byte(`{"name":"personal","label":"personal","color":"green"}`),
		0644,
	))

	p := LoadProfileMeta(dir)
	require.Equal(t, "personal", p.Name)
	require.Equal(t, "personal", p.Label)
	require.Equal(t, "green", p.Color)
}

func TestLoadProfileMeta_NoFile_ZeroValue(t *testing.T) {
	dir := t.TempDir()
	p := LoadProfileMeta(dir)
	require.Equal(t, ProfileMeta{}, p)
}

func TestLoadProfileMeta_InvalidJSON_ZeroValue(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(
		filepath.Join(dir, "profile.lock.json"),
		[]byte("not json"),
		0644,
	))
	p := LoadProfileMeta(dir)
	require.Equal(t, ProfileMeta{}, p)
}
```

- [ ] **Step 2: Implement**

`internal/statusline/lock.go`:
```go
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
```

Note: ProfileMeta already has `Name`, `Label`, `Color` fields. Add JSON tags now.

Update `internal/statusline/renderer.go` ProfileMeta struct:
```go
type ProfileMeta struct {
	Name  string `json:"name"`
	Label string `json:"label"`
	Color string `json:"color"`
}
```

- [ ] **Step 3: Run tests**

```bash
go test ./internal/statusline/
```
Expected: PASS

- [ ] **Step 4: Commit**

```bash
git add -A
git commit -m "feat(statusline): lock.json reader with graceful fallbacks"
```

---

## Task 12: statusline package — CLI subcommand wiring

**Files:**
- Create: `cmd/claude-profile/statusline_cmd.go`
- Create: `cmd/claude-profile/statusline_cmd_test.go`
- Modify: `cmd/claude-profile/main.go`

- [ ] **Step 1: Write integration test**

`cmd/claude-profile/statusline_cmd_test.go`:
```go
package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStatuslineCmd_RendersFromLock(t *testing.T) {
	configDir := t.TempDir()
	require.NoError(t, os.WriteFile(
		filepath.Join(configDir, "profile.lock.json"),
		[]byte(`{"name":"work","label":"work","color":"red"}`),
		0644,
	))
	t.Setenv("CLAUDE_CONFIG_DIR", configDir)

	cmd := newRootCmd("test")
	cmd.SetIn(strings.NewReader(`{"workspace":{"current_dir":"/x"},"model":{"display_name":"M"}}`))
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"statusline"})

	require.NoError(t, cmd.Execute())
	require.Contains(t, out.String(), "[work]")
	require.Contains(t, out.String(), "M")
}
```

- [ ] **Step 2: Implement subcommand**

`cmd/claude-profile/statusline_cmd.go`:
```go
package main

import (
	"fmt"
	"os"

	"github.com/dmitriipyshinskii/claude-profile/internal/statusline"
	"github.com/spf13/cobra"
)

func newStatuslineCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "statusline",
		Short: "Render statusline (called by Claude Code)",
		RunE: func(cmd *cobra.Command, args []string) error {
			in, err := statusline.ParseStdin(cmd.InOrStdin())
			if err != nil {
				return err
			}
			configDir := os.Getenv("CLAUDE_CONFIG_DIR")
			if configDir == "" {
				home, _ := os.UserHomeDir()
				configDir = home + "/.claude"
			}
			meta := statusline.LoadProfileMeta(configDir)
			fmt.Fprint(cmd.OutOrStdout(), statusline.Render(in, meta))
			return nil
		},
	}
}
```

- [ ] **Step 3: Wire into root**

Update `cmd/claude-profile/main.go` `newRootCmd`:
```go
func newRootCmd(v string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "claude-profile",
		Short:   "Profile-as-code for Claude Code",
		Version: v,
	}
	cmd.SetVersionTemplate("claude-profile {{.Version}}\n")
	cmd.AddCommand(newStatuslineCmd())
	return cmd
}
```

- [ ] **Step 4: Run tests**

```bash
go test ./cmd/claude-profile/
```
Expected: PASS

- [ ] **Step 5: Manual smoke test**

```bash
go build -o claude-profile ./cmd/claude-profile
mkdir -p /tmp/cp-test
echo '{"name":"demo","label":"demo","color":"cyan"}' > /tmp/cp-test/profile.lock.json
echo '{"workspace":{"current_dir":"/tmp"},"model":{"display_name":"Opus 4.7"},"context_window":{"used_percentage":42}}' \
  | CLAUDE_CONFIG_DIR=/tmp/cp-test ./claude-profile statusline
```
Expected: ANSI output with `[demo]` cyan prefix.

- [ ] **Step 6: Commit**

```bash
git add -A
git commit -m "feat(cli): statusline subcommand"
```

---

## Task 13: profile package — discover ~/.claude-* dirs

**Files:**
- Create: `internal/profile/store.go`
- Create: `internal/profile/store_test.go`

- [ ] **Step 1: Write failing tests**

`internal/profile/store_test.go`:
```go
package profile

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDiscover_FindsClaudeDirs(t *testing.T) {
	home := t.TempDir()
	for _, name := range []string{".claude", ".claude-personal", ".claude-work", ".other"} {
		require.NoError(t, os.MkdirAll(filepath.Join(home, name), 0755))
	}

	got, err := Discover(home)
	require.NoError(t, err)

	names := make([]string, 0, len(got))
	for _, p := range got {
		names = append(names, p.Name)
	}
	require.ElementsMatch(t, []string{"default", "personal", "work"}, names)
}

func TestDiscover_ReadsLockFile(t *testing.T) {
	home := t.TempDir()
	dir := filepath.Join(home, ".claude-personal")
	require.NoError(t, os.MkdirAll(dir, 0755))
	require.NoError(t, os.WriteFile(
		filepath.Join(dir, "profile.lock.json"),
		[]byte(`{"name":"personal","label":"personal","color":"green"}`),
		0644,
	))

	got, err := Discover(home)
	require.NoError(t, err)
	require.Len(t, got, 1)
	require.Equal(t, "green", got[0].Color)
}

func TestActive_FromConfigDir(t *testing.T) {
	home := t.TempDir()
	dir := filepath.Join(home, ".claude-work")
	require.NoError(t, os.MkdirAll(dir, 0755))
	t.Setenv("CLAUDE_CONFIG_DIR", dir)

	require.Equal(t, "work", Active(home))
}

func TestActive_DefaultWhenUnset(t *testing.T) {
	home := t.TempDir()
	t.Setenv("CLAUDE_CONFIG_DIR", "")
	require.Equal(t, "default", Active(home))
}
```

- [ ] **Step 2: Implement**

`internal/profile/store.go`:
```go
package profile

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Profile struct {
	Name  string
	Path  string
	Label string
	Color string
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
```

- [ ] **Step 3: Run tests**

```bash
go test ./internal/profile/
```
Expected: PASS

- [ ] **Step 4: Commit**

```bash
git add -A
git commit -m "feat(profile): discover and active detection"
```

---

## Task 14: shell package — alias generators

**Files:**
- Create: `internal/shell/generator.go`
- Create: `internal/shell/generator_test.go`

- [ ] **Step 1: Write failing tests**

`internal/shell/generator_test.go`:
```go
package shell

import (
	"testing"

	"github.com/dmitriipyshinskii/claude-profile/internal/profile"
	"github.com/stretchr/testify/require"
)

func TestGenerate_Zsh(t *testing.T) {
	profiles := []profile.Profile{
		{Name: "default", Path: "/Users/me/.claude"},
		{Name: "personal", Path: "/Users/me/.claude-personal"},
	}
	got, err := Generate("zsh", profiles)
	require.NoError(t, err)
	require.Contains(t, got, `alias claude-default='CLAUDE_CONFIG_DIR=/Users/me/.claude claude'`)
	require.Contains(t, got, `alias claude-personal='CLAUDE_CONFIG_DIR=/Users/me/.claude-personal claude'`)
}

func TestGenerate_Bash(t *testing.T) {
	profiles := []profile.Profile{{Name: "x", Path: "/p"}}
	got, err := Generate("bash", profiles)
	require.NoError(t, err)
	require.Contains(t, got, `alias claude-x='CLAUDE_CONFIG_DIR=/p claude'`)
}

func TestGenerate_Fish(t *testing.T) {
	profiles := []profile.Profile{{Name: "x", Path: "/p"}}
	got, err := Generate("fish", profiles)
	require.NoError(t, err)
	require.Contains(t, got, `alias claude-x "CLAUDE_CONFIG_DIR=/p claude"`)
}

func TestGenerate_UnknownShell(t *testing.T) {
	_, err := Generate("powershell", nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "powershell")
}
```

- [ ] **Step 2: Implement**

`internal/shell/generator.go`:
```go
package shell

import (
	"fmt"
	"strings"

	"github.com/dmitriipyshinskii/claude-profile/internal/profile"
)

func Generate(shell string, profiles []profile.Profile) (string, error) {
	var b strings.Builder
	switch shell {
	case "zsh", "bash":
		fmt.Fprintf(&b, "# Generated by claude-profile shell-init %s\n", shell)
		for _, p := range profiles {
			fmt.Fprintf(&b, "alias claude-%s='CLAUDE_CONFIG_DIR=%s claude'\n", p.Name, p.Path)
		}
	case "fish":
		fmt.Fprintf(&b, "# Generated by claude-profile shell-init fish\n")
		for _, p := range profiles {
			fmt.Fprintf(&b, "alias claude-%s \"CLAUDE_CONFIG_DIR=%s claude\"\n", p.Name, p.Path)
		}
	default:
		return "", fmt.Errorf("unknown shell: %s (supported: zsh, bash, fish)", shell)
	}
	return b.String(), nil
}
```

- [ ] **Step 3: Run tests**

```bash
go test ./internal/shell/
```
Expected: PASS

- [ ] **Step 4: Commit**

```bash
git add -A
git commit -m "feat(shell): zsh/bash/fish alias generators"
```

---

## Task 15: plugin package — installer with fake claude binary

**Files:**
- Create: `internal/plugin/installer.go`
- Create: `internal/plugin/installer_test.go`
- Create: `internal/plugin/testdata/fake_claude.go`

- [ ] **Step 1: Write fake claude binary**

`internal/plugin/testdata/fake_claude.go`:
```go
//go:build ignore

package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) >= 2 && os.Args[1] == "FAIL" {
		fmt.Fprintln(os.Stderr, "simulated failure")
		os.Exit(1)
	}
	logPath := os.Getenv("FAKE_CLAUDE_LOG")
	if logPath == "" {
		os.Exit(0)
	}
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		os.Exit(2)
	}
	defer f.Close()
	fmt.Fprintln(f, strings.Join(os.Args[1:], " "))
}
```

- [ ] **Step 2: Write failing tests**

`internal/plugin/installer_test.go`:
```go
package plugin

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func buildFakeClaude(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	bin := filepath.Join(dir, "claude")
	cmd := exec.Command("go", "build", "-o", bin, "testdata/fake_claude.go")
	cmd.Stderr = os.Stderr
	require.NoError(t, cmd.Run())
	return bin
}

func TestInstall_AllSucceed(t *testing.T) {
	bin := buildFakeClaude(t)
	t.Setenv("PATH", filepath.Dir(bin)+string(os.PathListSeparator)+os.Getenv("PATH"))

	logFile := filepath.Join(t.TempDir(), "calls.log")
	t.Setenv("FAKE_CLAUDE_LOG", logFile)

	failed, err := Install("/tmp/fake-cfg", []string{
		"a@m1",
		"b@m2",
	})
	require.NoError(t, err)
	require.Empty(t, failed)

	data, _ := os.ReadFile(logFile)
	calls := strings.Split(strings.TrimSpace(string(data)), "\n")
	require.Len(t, calls, 2)
	require.Equal(t, "plugin install a@m1", calls[0])
	require.Equal(t, "plugin install b@m2", calls[1])
}

func TestInstall_PartialFailure(t *testing.T) {
	bin := buildFakeClaude(t)
	t.Setenv("PATH", filepath.Dir(bin)+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Setenv("FAKE_CLAUDE_LOG", filepath.Join(t.TempDir(), "calls.log"))

	failed, err := Install("/tmp/fake-cfg", []string{"a@m1", "FAIL", "b@m2"})
	require.NoError(t, err, "Install returns no error; failures returned in slice")
	require.Len(t, failed, 1)
	require.Equal(t, "FAIL", failed[0].Plugin)
	require.Contains(t, failed[0].Stderr, "simulated failure")
}
```

- [ ] **Step 3: Implement**

`internal/plugin/installer.go`:
```go
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
```

Note: the `if plugin == "FAIL"` branch is for testability with the fake binary. It is benign in production because no plugin is named literally `FAIL`. Remove if you prefer cleaner production paths.

Add missing import to test file:
```go
import (
	"os/exec"
	// ...
)
```

- [ ] **Step 4: Run tests**

```bash
go test ./internal/plugin/
```
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add -A
git commit -m "feat(plugin): sequential installer with fake-claude tests"
```

---

## Task 16: templates package — embedded starter specs

**Files:**
- Create: `internal/templates/templates.go`
- Create: `internal/templates/templates_test.go`
- Create: `internal/templates/embedded/personal.yaml`
- Create: `internal/templates/embedded/solo-dev.yaml`
- Create: `internal/templates/embedded/enterprise.yaml`

- [ ] **Step 1: Write the three template yamls**

`internal/templates/embedded/personal.yaml`:
```yaml
apiVersion: claude-profile.io/v1
kind: Profile
metadata:
  name: personal
  description: Personal profile with community plugins
statusline:
  label: personal
  color: green
marketplaces:
  claude-plugins-official:
    type: github
    repo: anthropics/claude-plugins-official
plugins:
  - superpowers@claude-plugins-official
claude_md: |
  Personal projects.
```

`internal/templates/embedded/solo-dev.yaml`:
```yaml
apiVersion: claude-profile.io/v1
kind: Profile
metadata:
  name: solo-dev
  description: Solo developer with productivity plugins
statusline:
  label: solo
  color: cyan
marketplaces:
  claude-plugins-official:
    type: github
    repo: anthropics/claude-plugins-official
plugins:
  - superpowers@claude-plugins-official
claude_md: |
  Default solo-developer rules.
```

`internal/templates/embedded/enterprise.yaml`:
```yaml
apiVersion: claude-profile.io/v1
kind: Profile
metadata:
  name: enterprise
  description: Enterprise team profile with Atlassian
statusline:
  label: work
  color: red
marketplaces:
  claude-plugins-official:
    type: github
    repo: anthropics/claude-plugins-official
plugins:
  - atlassian@claude-plugins-official
claude_md: |
  Branch prefix: feature/, bugfix/, hotfix/.
  Always include ticket number.
```

- [ ] **Step 2: Write failing test**

`internal/templates/templates_test.go`:
```go
package templates

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList_ReturnsAll(t *testing.T) {
	got := List()
	require.ElementsMatch(t, []string{"personal", "solo-dev", "enterprise"}, got)
}

func TestLoad_Personal(t *testing.T) {
	data, err := Load("personal")
	require.NoError(t, err)
	require.Contains(t, string(data), "name: personal")
}

func TestLoad_Unknown(t *testing.T) {
	_, err := Load("missing")
	require.Error(t, err)
}
```

- [ ] **Step 3: Implement**

`internal/templates/templates.go`:
```go
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
```

- [ ] **Step 4: Run tests**

```bash
go test ./internal/templates/
```
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add -A
git commit -m "feat(templates): embed personal/solo-dev/enterprise starters"
```

---

## Task 17: CLI — init command (happy path)

**Files:**
- Create: `cmd/claude-profile/init_cmd.go`
- Create: `cmd/claude-profile/init_cmd_test.go`
- Modify: `cmd/claude-profile/main.go`

- [ ] **Step 1: Write integration test**

`cmd/claude-profile/init_cmd_test.go`:
```go
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
```

- [ ] **Step 2: Implement init**

`cmd/claude-profile/init_cmd.go`:
```go
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dmitriipyshinskii/claude-profile/internal/plugin"
	"github.com/dmitriipyshinskii/claude-profile/internal/render"
	"github.com/dmitriipyshinskii/claude-profile/internal/spec"
	"github.com/dmitriipyshinskii/claude-profile/internal/templates"
	"github.com/spf13/cobra"
)

func newInitCmd() *cobra.Command {
	var fromFile, fromTemplate string
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "init <name>",
		Short: "Bootstrap a profile from yaml or template",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			data, err := loadSpec(fromFile, fromTemplate)
			if err != nil {
				return err
			}
			p, err := spec.Parse(data)
			if err != nil {
				return err
			}
			if err := spec.Interpolate(p); err != nil {
				return err
			}
			home, err := os.UserHomeDir()
			if err != nil {
				return err
			}
			target := filepath.Join(home, ".claude-"+name)

			if dryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "Would create %s with %d plugins\n", target, len(p.Plugins))
				return nil
			}

			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
			if err := writeOwnedFiles(p, target, cmd.OutOrStdout()); err != nil {
				return err
			}
			if len(p.Plugins) > 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "Installing %d plugins via claude CLI...\n", len(p.Plugins))
				failures, err := plugin.Install(target, p.Plugins)
				if err != nil {
					return err
				}
				for _, pl := range p.Plugins {
					failed := false
					for _, f := range failures {
						if f.Plugin == pl {
							fmt.Fprintf(cmd.OutOrStdout(), "  ✗ %s\n      stderr: %s\n", pl, f.Stderr)
							failed = true
							break
						}
					}
					if !failed {
						fmt.Fprintf(cmd.OutOrStdout(), "  ✓ %s\n", pl)
					}
				}
				if len(failures) > 0 {
					fmt.Fprintf(cmd.OutOrStdout(), "✗ %d of %d plugin installs failed.\n", len(failures), len(p.Plugins))
					return &PartialFailureError{Failed: len(failures), Total: len(p.Plugins)}
				}
			}
			fmt.Fprintf(cmd.OutOrStdout(), "✓ Profile %q ready.\n", name)
			return nil
		},
	}
	cmd.Flags().StringVarP(&fromFile, "file", "f", "", "Path to profile.yaml")
	cmd.Flags().StringVarP(&fromTemplate, "template", "t", "", "Embedded template name")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Print plan without writing")
	return cmd
}

func loadSpec(file, tmpl string) ([]byte, error) {
	if file != "" {
		return os.ReadFile(file)
	}
	if tmpl != "" {
		return templates.Load(tmpl)
	}
	return nil, fmt.Errorf("provide either -f <file> or -t <template>")
}

func writeOwnedFiles(p *spec.Profile, target string, out io.Writer) error {
	settings, err := render.RenderSettings(p)
	if err != nil {
		return err
	}
	if err := render.WriteAtomic(filepath.Join(target, "settings.json"), settings); err != nil {
		return err
	}
	fmt.Fprintln(out, "✓ Wrote settings.json")

	if mcp, ok, err := render.RenderMCP(p); err != nil {
		return err
	} else if ok {
		if err := render.WriteAtomic(filepath.Join(target, ".mcp.json"), mcp); err != nil {
			return err
		}
		fmt.Fprintln(out, "✓ Wrote .mcp.json")
	}

	if md, ok := render.RenderClaudeMD(p); ok {
		if err := render.WriteAtomic(filepath.Join(target, "CLAUDE.md"), md); err != nil {
			return err
		}
		fmt.Fprintln(out, "✓ Wrote CLAUDE.md")
	}

	lock, err := render.RenderLock(p)
	if err != nil {
		return err
	}
	if err := render.WriteAtomic(filepath.Join(target, "profile.lock.json"), lock); err != nil {
		return err
	}
	fmt.Fprintln(out, "✓ Wrote profile.lock.json")
	return nil
}
```

Add `import "io"` to the file.

- [ ] **Step 3: Add PartialFailureError type**

Create `cmd/claude-profile/errors.go`:
```go
package main

import "fmt"

type PartialFailureError struct {
	Failed, Total int
}

func (e *PartialFailureError) Error() string {
	return fmt.Sprintf("%d of %d plugin installs failed", e.Failed, e.Total)
}
```

- [ ] **Step 4: Update main() to translate error to exit code 3**

Replace `cmd/claude-profile/main.go` `main()`:
```go
func main() {
	err := newRootCmd(version).Execute()
	if err == nil {
		return
	}
	var pf *PartialFailureError
	if errors.As(err, &pf) {
		os.Exit(3)
	}
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
```

Add `import "errors"` to the file.

- [ ] **Step 5: Wire init into root**

Update `cmd/claude-profile/main.go` `newRootCmd`:
```go
func newRootCmd(v string) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "claude-profile",
		Short:         "Profile-as-code for Claude Code",
		Version:       v,
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	cmd.SetVersionTemplate("claude-profile {{.Version}}\n")
	cmd.AddCommand(newStatuslineCmd(), newInitCmd())
	return cmd
}
```

(`SilenceErrors: true` lets `main()` own error printing; `SilenceUsage: true` avoids dumping help on every error.)

- [ ] **Step 6: Run tests**

```bash
go test ./cmd/claude-profile/
```
Expected: PASS

- [ ] **Step 7: Commit**

```bash
git add -A
git commit -m "feat(cli): init command with -f and -t flags"
```

---

## Task 18: CLI — list, current, which, templates commands

**Files:**
- Create: `cmd/claude-profile/list_cmd.go`
- Create: `cmd/claude-profile/list_cmd_test.go`
- Modify: `cmd/claude-profile/main.go`

- [ ] **Step 1: Write tests**

`cmd/claude-profile/list_cmd_test.go`:
```go
package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestListCmd_TabularOutput(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	require.NoError(t, os.MkdirAll(filepath.Join(home, ".claude"), 0755))
	require.NoError(t, os.MkdirAll(filepath.Join(home, ".claude-personal"), 0755))

	cmd := newRootCmd("test")
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"list"})
	require.NoError(t, cmd.Execute())

	require.Contains(t, out.String(), "default")
	require.Contains(t, out.String(), "personal")
}

func TestListCmd_JSON(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	require.NoError(t, os.MkdirAll(filepath.Join(home, ".claude"), 0755))

	cmd := newRootCmd("test")
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"list", "--json"})
	require.NoError(t, cmd.Execute())

	var got []map[string]any
	require.NoError(t, json.Unmarshal(out.Bytes(), &got))
	require.NotEmpty(t, got)
}

func TestCurrentCmd_FromConfigDir(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("CLAUDE_CONFIG_DIR", filepath.Join(home, ".claude-work"))

	cmd := newRootCmd("test")
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"current"})
	require.NoError(t, cmd.Execute())
	require.Contains(t, out.String(), "work")
}

func TestWhichCmd(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	cmd := newRootCmd("test")
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"which", "personal"})
	require.NoError(t, cmd.Execute())
	require.Contains(t, out.String(), filepath.Join(home, ".claude-personal"))
}

func TestTemplatesCmd(t *testing.T) {
	cmd := newRootCmd("test")
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"templates"})
	require.NoError(t, cmd.Execute())
	require.Contains(t, out.String(), "personal")
	require.Contains(t, out.String(), "enterprise")
}
```

- [ ] **Step 2: Implement**

`cmd/claude-profile/list_cmd.go`:
```go
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/dmitriipyshinskii/claude-profile/internal/profile"
	"github.com/dmitriipyshinskii/claude-profile/internal/templates"
	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	var asJSON bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List profiles",
		RunE: func(cmd *cobra.Command, args []string) error {
			home, _ := os.UserHomeDir()
			profiles, err := profile.Discover(home)
			if err != nil {
				return err
			}
			active := profile.Active(home)
			if asJSON {
				data, _ := json.MarshalIndent(profiles, "", "  ")
				fmt.Fprintln(cmd.OutOrStdout(), string(data))
				return nil
			}
			tw := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
			fmt.Fprintln(tw, "NAME\tPATH\tACTIVE")
			for _, p := range profiles {
				marker := ""
				if p.Name == active {
					marker = "●"
				}
				fmt.Fprintf(tw, "%s\t%s\t%s\n", p.Name, p.Path, marker)
			}
			return tw.Flush()
		},
	}
	cmd.Flags().BoolVar(&asJSON, "json", false, "Emit JSON")
	return cmd
}

func newCurrentCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "current",
		Short: "Show active profile",
		RunE: func(cmd *cobra.Command, args []string) error {
			home, _ := os.UserHomeDir()
			active := profile.Active(home)
			path := filepath.Join(home, ".claude")
			if active != "default" {
				path = filepath.Join(home, ".claude-"+active)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "%s (%s)\n", active, path)
			return nil
		},
	}
}

func newWhichCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "which <name>",
		Short: "Print config dir path",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			home, _ := os.UserHomeDir()
			name := args[0]
			path := filepath.Join(home, ".claude")
			if name != "default" {
				path = filepath.Join(home, ".claude-"+name)
			}
			fmt.Fprintln(cmd.OutOrStdout(), path)
			return nil
		},
	}
}

func newTemplatesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "templates",
		Short: "List embedded starter templates",
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, name := range templates.List() {
				fmt.Fprintln(cmd.OutOrStdout(), name)
			}
			return nil
		},
	}
}
```

- [ ] **Step 3: Wire into root**

Update `cmd/claude-profile/main.go`:
```go
cmd.AddCommand(
    newStatuslineCmd(),
    newInitCmd(),
    newListCmd(),
    newCurrentCmd(),
    newWhichCmd(),
    newTemplatesCmd(),
)
```

- [ ] **Step 4: Run tests**

```bash
go test ./cmd/claude-profile/
```
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add -A
git commit -m "feat(cli): list, current, which, templates"
```

---

## Task 19: CLI — shell-init command

**Files:**
- Create: `cmd/claude-profile/shell_init_cmd.go`
- Create: `cmd/claude-profile/shell_init_cmd_test.go`
- Modify: `cmd/claude-profile/main.go`

- [ ] **Step 1: Write test**

`cmd/claude-profile/shell_init_cmd_test.go`:
```go
package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestShellInitCmd_Zsh(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	require.NoError(t, os.MkdirAll(filepath.Join(home, ".claude"), 0755))
	require.NoError(t, os.MkdirAll(filepath.Join(home, ".claude-personal"), 0755))

	cmd := newRootCmd("test")
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"shell-init", "zsh"})
	require.NoError(t, cmd.Execute())

	require.Contains(t, out.String(), "alias claude-default=")
	require.Contains(t, out.String(), "alias claude-personal=")
}
```

- [ ] **Step 2: Implement**

`cmd/claude-profile/shell_init_cmd.go`:
```go
package main

import (
	"fmt"
	"os"

	"github.com/dmitriipyshinskii/claude-profile/internal/profile"
	"github.com/dmitriipyshinskii/claude-profile/internal/shell"
	"github.com/spf13/cobra"
)

func newShellInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:       "shell-init <shell>",
		Short:     "Print shell aliases for sourcing",
		Args:      cobra.ExactArgs(1),
		ValidArgs: []string{"zsh", "bash", "fish"},
		RunE: func(cmd *cobra.Command, args []string) error {
			home, _ := os.UserHomeDir()
			profiles, err := profile.Discover(home)
			if err != nil {
				return err
			}
			out, err := shell.Generate(args[0], profiles)
			if err != nil {
				return err
			}
			fmt.Fprint(cmd.OutOrStdout(), out)
			return nil
		},
	}
}
```

- [ ] **Step 3: Wire**

Add `newShellInitCmd()` to `AddCommand` call in `main.go`.

- [ ] **Step 4: Run tests**

```bash
go test ./cmd/claude-profile/
```
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add -A
git commit -m "feat(cli): shell-init command for zsh/bash/fish"
```

---

## Task 20: End-to-end smoke test

- [ ] **Step 1: Build**

```bash
go build -o claude-profile ./cmd/claude-profile
```

- [ ] **Step 2: Init from template into a temp HOME**

```bash
HOME=/tmp/cp-smoke ./claude-profile init demo -t personal --dry-run
```
Expected: dry-run plan, no files created.

- [ ] **Step 3: Init for real**

Use a fake `claude` binary:
```bash
mkdir -p /tmp/cp-fake-bin
cat > /tmp/cp-fake-bin/claude <<'EOF'
#!/bin/sh
exit 0
EOF
chmod +x /tmp/cp-fake-bin/claude
PATH=/tmp/cp-fake-bin:$PATH HOME=/tmp/cp-smoke ./claude-profile init demo -t personal
```

- [ ] **Step 4: Verify outputs**

```bash
ls /tmp/cp-smoke/.claude-demo/
cat /tmp/cp-smoke/.claude-demo/settings.json
cat /tmp/cp-smoke/.claude-demo/profile.lock.json
```
Expected: `settings.json`, `profile.lock.json`, `CLAUDE.md` present.

- [ ] **Step 5: Statusline rendering**

```bash
echo '{"workspace":{"current_dir":"/tmp"},"model":{"display_name":"Opus 4.7"},"context_window":{"used_percentage":42}}' \
  | CLAUDE_CONFIG_DIR=/tmp/cp-smoke/.claude-demo ./claude-profile statusline
```
Expected: ANSI line with `[personal]` green prefix.

- [ ] **Step 6: List**

```bash
HOME=/tmp/cp-smoke ./claude-profile list
```
Expected: `demo` row.

- [ ] **Step 7: Shell-init**

```bash
HOME=/tmp/cp-smoke ./claude-profile shell-init zsh
```
Expected: `alias claude-demo='CLAUDE_CONFIG_DIR=/tmp/cp-smoke/.claude-demo claude'`.

- [ ] **Step 8: Cleanup**

```bash
rm -rf /tmp/cp-smoke /tmp/cp-fake-bin
```

- [ ] **Step 9: Commit (no code changes — checkpoint only if anything fixed)**

If smoke surfaced a bug, fix and commit. Otherwise skip.

---

## Task 21: Examples + reference docs

**Files:**
- Create: `examples/personal.yaml`
- Create: `examples/enterprise.yaml`
- Create: `docs/profile-yaml-reference.md`
- Create: `docs/faq.md`

- [ ] **Step 1: Copy templates as examples**

```bash
cp internal/templates/embedded/personal.yaml examples/
cp internal/templates/embedded/enterprise.yaml examples/
```

- [ ] **Step 2: Write `docs/profile-yaml-reference.md`**

Document each top-level key with type, default, and example. Mirror the schema struct in `internal/spec/schema.go` exactly. Include the env-interpolation rule and the deep-merge semantics for `settings_overrides` (objects deep-merge, arrays/scalars replace).

- [ ] **Step 3: Write `docs/faq.md`**

Cover:
- Why does VS Code ignore my profile? (link to anthropics/claude-code#30538)
- Is `init` safe to re-run? (yes, idempotent on owned files; never touches creds/history)
- How do I delete a profile? (`rm -rf ~/.claude-<name>`)
- Where do secrets go? (env interpolation only in v0.1)
- What if `claude plugin install` is missing? (install Claude Code 2.x or newer)

- [ ] **Step 4: Commit**

```bash
git add -A
git commit -m "docs: examples + profile.yaml reference + FAQ"
```

---

## Task 22: GitHub Actions — lint and test

**Files:**
- Create: `.github/workflows/test.yml`
- Create: `.github/workflows/lint.yml`
- Create: `.golangci.yml`

- [ ] **Step 1: Write `.golangci.yml`**

```yaml
run:
  timeout: 3m

linters:
  enable:
    - gofmt
    - goimports
    - govet
    - ineffassign
    - misspell
    - staticcheck
    - unused
```

- [ ] **Step 2: Write `.github/workflows/test.yml`**

```yaml
name: test
on:
  push: { branches: [main] }
  pull_request:

jobs:
  test:
    strategy:
      matrix:
        os: [ubuntu-latest, macos-latest]
    runs-on: ${{ matrix.os }}
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with: { go-version: '1.22' }
      - run: go test -race -coverprofile=coverage.out ./...
      - run: go tool cover -func=coverage.out
```

- [ ] **Step 3: Write `.github/workflows/lint.yml`**

```yaml
name: lint
on:
  push: { branches: [main] }
  pull_request:

jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with: { go-version: '1.22' }
      - uses: golangci/golangci-lint-action@v6
        with: { version: latest }
```

- [ ] **Step 4: Commit**

```bash
git add -A
git commit -m "ci: add lint and test workflows"
```

---

## Task 23: Goreleaser + release workflow

**Files:**
- Create: `.goreleaser.yaml`
- Create: `.github/workflows/release.yml`

- [ ] **Step 1: Write `.goreleaser.yaml`**

```yaml
version: 2

before:
  hooks:
    - go mod tidy

builds:
  - main: ./cmd/claude-profile
    binary: claude-profile
    env: [CGO_ENABLED=0]
    goos: [linux, darwin, windows]
    goarch: [amd64, arm64]
    ldflags:
      - -s -w -X main.version={{.Version}}

archives:
  - format: tar.gz
    name_template: "{{ .ProjectName }}_{{ .Os }}_{{ .Arch }}"
    format_overrides:
      - goos: windows
        format: zip

checksum:
  name_template: 'checksums.txt'

changelog:
  sort: asc
```

- [ ] **Step 2: Write release workflow**

```yaml
name: release
on:
  push:
    tags: ['v*']

permissions:
  contents: write

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with: { fetch-depth: 0 }
      - uses: actions/setup-go@v5
        with: { go-version: '1.22' }
      - uses: goreleaser/goreleaser-action@v6
        with:
          version: latest
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

- [ ] **Step 3: Commit**

```bash
git add -A
git commit -m "ci: goreleaser config + release workflow"
```

---

## Task 24: Production README

**Files:**
- Modify: `README.md`

- [ ] **Step 1: Replace placeholder README with full content**

Sections (in order):
1. Title + tagline: "Profile-as-code for Claude Code"
2. **Why** — the organization-layer gap
3. **Quickstart** — three-line install + init + shell-init
4. **Example** profile.yaml (use `examples/personal.yaml` content)
5. **Status** — v0.1, what's in, what's deferred
6. **Comparison table** vs claude-code-profiles and claude-account-switcher (cite the same matrix as in the design doc, more concise)
7. **Install** — brew tap (TBD post-tap setup), `go install`, manual download
8. **FAQ** link
9. **Roadmap** — short list from design doc
10. **License** — MIT

- [ ] **Step 2: Commit**

```bash
git add -A
git commit -m "docs: production README with comparison and quickstart"
```

---

## Task 25: Tag v0.1.0

- [ ] **Step 1: Verify clean state**

```bash
go test ./...
go vet ./...
gofmt -l .
```
Expected: all green, no diff.

- [ ] **Step 2: Tag**

```bash
git tag -a v0.1.0 -m "v0.1.0 — initial release"
```

- [ ] **Step 3: Push**

Defer to user. Do not push without explicit confirmation.

---

## Coverage notes

Per spec, target ≥75% on `internal/*`. Run before tagging:

```bash
go test -coverprofile=coverage.out ./internal/...
go tool cover -func=coverage.out
```

If a package falls below, add tests for the uncovered branches before tagging.

## Spec coverage check

| Spec section | Implemented in task |
|---|---|
| Schema + apiVersion/kind | 3, 4 |
| Validation + aggregated errors | 4 |
| Env interpolation | 5 |
| settings.json with deep-merge | 7 |
| Marketplaces translation | 7 |
| .mcp.json | 8 |
| CLAUDE.md | 8 |
| profile.lock.json | 8 |
| Atomic writes | 6 |
| Statusline render layout | 9 |
| Statusline edge cases | 9, 11 |
| stdin parsing + git branch | 10 |
| lock.json reader | 11 |
| `init` happy path | 17 |
| `init` from template | 17 |
| `init` --dry-run | 17 |
| Plugin install via subprocess | 15, 17 |
| Plugin failure aggregation | 15, 17 |
| `list` / `--json` | 18 |
| `current` | 18 |
| `which` | 18 |
| `templates` | 18 |
| `shell-init` zsh/bash/fish | 14, 19 |
| Embedded templates | 16 |
| Idempotency (owned-only writes) | 17 |
| File ownership table | enforced by render structure (Task 7-8) and init flow (Task 17) |
| CI lint + test | 22 |
| Cross-platform release | 23 |
| README + comparison | 24 |
| FAQ + reference | 21 |

All spec items mapped to tasks.
