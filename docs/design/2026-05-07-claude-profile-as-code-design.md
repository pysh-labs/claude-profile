# claude-profile: profile-as-code for Claude Code

**Status:** design — pending implementation
**Date:** 2026-05-07
**Owner:** Dmitrii Pyshinskii

## Why

Claude Code natively supports three configuration layers — user, project, project-local — and exposes the `CLAUDE_CONFIG_DIR` environment variable to redirect all state into an isolated directory. The community uses this to run multiple "accounts" on one machine.

What is missing is an **organization layer**. Working across multiple organizations on a single laptop, each with its own plugins, MCP servers, branch-naming rules, statusline conventions, and access policies, requires a per-org configuration set that is reproducible, shareable, and visually distinguishable at runtime.

Existing tools either:

- Manage profiles imperatively without a declarative spec ([`quinnjr/claude-code-profiles`](https://github.com/quinnjr/claude-code-profiles))
- Symlink shared config across profiles, sacrificing isolation ([`ukogan/claude-account-switcher`](https://github.com/ukogan/claude-account-switcher))
- Solve a different problem (parallel agents — [`smtg-ai/claude-squad`](https://github.com/smtg-ai/claude-squad))

`claude-profile` fills the gap: a declarative `profile.yaml`, a templated bootstrap, a per-profile statusline with visual identifier, and reproducibility from spec.

Positioning: **profile-as-code for Claude Code**. Analogies: `direnv` for environment, `terraform` for declared state, `nix-darwin` for reproducible workstations.

## Non-goals (v0.1)

- Reconciliation / drift detection (yaml is bootstrap-only)
- Auth bootstrap (each profile logs in once via `claude` CLI)
- VS Code extension support (known upstream issue [anthropics/claude-code#30538](https://github.com/anthropics/claude-code/issues/30538) — documented, not solved here)
- Plugin uninstall, version pinning, secret management beyond env interpolation
- Custom statusline templates
- `extends:` profile inheritance

These belong in v0.2+ roadmap.

## Assumptions

- The `claude plugin install <name>@<marketplace>` subcommand exists and is stable in Claude Code 2.x. The first implementation task is to verify this against the installed `claude` version; if the subcommand differs, the spec must be revised before further work.
- `CLAUDE_CONFIG_DIR` redirects all profile state including plugin install state, so `claude plugin install` invoked with that env var operates on the target profile. Verified against user's current setup (`~/.claude-personal/plugins/installed_plugins.json` is populated when `CLAUDE_CONFIG_DIR=~/.claude-personal`).

## Architecture

### Authority model

`profile.yaml` is a **bootstrap template**. `init` reads it, creates the profile directory, and writes a fixed set of "owned" files. After `init`, on-disk state can drift from spec — this is intentional and matches user expectations for v0.1.

### Stack

- Language: Go (single static binary, cross-platform via `goreleaser`)
- CLI: `spf13/cobra`
- YAML: `gopkg.in/yaml.v3`
- Validation: `go-playground/validator/v10`
- Tests: standard library + `stretchr/testify`

### Package layout

```
claude-profile/
├── cmd/claude-profile/main.go
├── internal/
│   ├── spec/         parse + validate profile.yaml
│   ├── render/       spec → on-disk files (settings.json, .mcp.json, CLAUDE.md, profile.lock.json)
│   ├── plugin/       drive `claude plugin install` via subprocess
│   ├── statusline/   stdin JSON → ANSI output
│   ├── shell/        generate aliases for zsh/bash/fish
│   ├── profile/      discover ~/.claude-* dirs, detect active
│   └── templates/    embedded starter specs (go:embed)
├── docs/profile-yaml-reference.md
├── docs/faq.md
├── examples/personal.yaml
├── examples/enterprise.yaml
├── examples/solo-dev.yaml
├── .github/workflows/{lint,test,release}.yml
├── .goreleaser.yaml
├── go.mod
├── Makefile
├── LICENSE                   # MIT
└── README.md
```

**Decoupling rules:**

- `spec` knows nothing about the filesystem — pure parsing + validation
- `render` accepts `spec.Profile`, writes files, knows file paths but no business logic
- `plugin` accepts a list of `name@marketplace` strings, does not depend on `spec`
- `statusline` is standalone — works without any other package being initialised

## Profile schema

```yaml
apiVersion: claude-profile.io/v1
kind: Profile

metadata:
  name: personal
  description: Personal projects + community plugins

statusline:
  label: personal                    # short prefix in statusline
  color: green                       # red|green|yellow|blue|magenta|cyan|white

marketplaces:                        # written to settings.json extraKnownMarketplaces
  claude-plugins-official:
    type: github
    repo: anthropics/claude-plugins-official

plugins:                             # list of plugin@marketplace
  - atlassian@claude-plugins-official
  - superpowers@claude-plugins-official
  - clangd-lsp@claude-plugins-official

mcp_servers:                         # optional, written to .mcp.json
  atlassian:
    command: docker
    args: [run, -i, --rm, mcp/atlassian]
    env:
      ATLASSIAN_USERNAME: ${ATLASSIAN_USERNAME}    # interpolated at init-time

settings_overrides:                  # deep-merged into settings.json after our keys
  skipDangerousModePermissionPrompt: true
  permissions:
    allow: ["Bash(az *:*)"]

claude_md: |
  ALWAYS use feature/, bugfix/, hotfix/ branch prefixes.
  ALWAYS include Jira ticket in branch name.
```

### Schema decisions

- `apiVersion` + `kind`: k8s-style versioning. Cheap now, pays off when v2 schema arrives.
- `plugins` as `name@marketplace` strings: matches `installed_plugins.json` format, easy to copy. Object form with `version` pinning is v0.2.
- `statusline.color` is a named palette of seven colors. Raw ANSI codes deferred.
- `mcp_servers.*.env` supports `${VAR}` interpolation at init-time. Missing var → fail fast with clear error. Secrets stay in environment, yaml is safe to commit.
- `settings_overrides` deep-merges on top of our generated keys. Merge semantics: objects deep-merge recursively; arrays and scalars replace; conflict on objects → overrides win. Documented in `docs/profile-yaml-reference.md`.
- `claude_md` is inline string in v0.1. Includes / overlays come in v0.2.
- The render layer translates the simplified `marketplaces` spec form into Claude Code's actual `extraKnownMarketplaces` shape (`{source: {source: "github", repo: "..."}}`). Users do not need to know the underlying shape.

## CLI surface

```
claude-profile init <name> [-f profile.yaml | -t <template>] [--dry-run] [-v]
claude-profile list  [--json]
claude-profile current
claude-profile which <name>
claude-profile templates
claude-profile shell-init zsh|bash|fish
claude-profile statusline                   # internal, called by Claude Code
claude-profile version
claude-profile help [<command>]
```

### Sample sessions

**Init from spec:**

```
$ claude-profile init personal -f ./profiles/personal.yaml
✓ Created ~/.claude-personal/
✓ Wrote settings.json
✓ Wrote .mcp.json (1 server)
✓ Wrote CLAUDE.md
✓ Wrote profile.lock.json
✓ Installing 3 plugins via claude CLI...
  ✓ atlassian@claude-plugins-official
  ✓ superpowers@claude-plugins-official
  ✓ clangd-lsp@claude-plugins-official
✓ Profile 'personal' ready.

Next:
  • eval "$(claude-profile shell-init zsh)" >> ~/.zshrc
  • Launch:  claude-personal
```

**Init from embedded template:**

```
$ claude-profile init work -t enterprise
```

**List:**

```
$ claude-profile list
NAME       PATH                  PLUGINS  ACTIVE
default    ~/.claude             1
personal   ~/.claude-personal    3        ●
work       ~/.claude-work        5
```

**Shell aliases:**

```
$ claude-profile shell-init zsh
alias claude-default='CLAUDE_CONFIG_DIR=$HOME/.claude claude'
alias claude-personal='CLAUDE_CONFIG_DIR=$HOME/.claude-personal claude'
alias claude-work='CLAUDE_CONFIG_DIR=$HOME/.claude-work claude'
```

### UX rules

- `init` is **idempotent**: re-running on an existing profile re-renders owned files and skips already-installed plugins. User data (credentials, history, sessions) is never touched.
- `init` is **purely additive**: removing a plugin from yaml does not uninstall it. Documented limitation.
- `--dry-run` prints the plan with full env interpolation but writes nothing.
- `--json` is supported on `list` only; other commands are human-oriented.
- `current` reads `CLAUDE_CONFIG_DIR`; falls back to `default` (`~/.claude`).
- No `delete` subcommand in v0.1 — `rm -rf ~/.claude-<name>` is documented. Avoids accidental destruction surface.

## Statusline rendering

```
Claude Code  ──stdin JSON──▶  claude-profile statusline  ──stdout ANSI──▶  Claude Code UI
                                       │
                                       └── reads $CLAUDE_CONFIG_DIR/profile.lock.json
```

`profile.lock.json` is written at `init` and contains the metadata needed by the renderer:

```json
{ "name": "personal", "label": "personal", "color": "green" }
```

If the lock file is absent (profile created externally), the renderer skips the prefix and degrades gracefully.

### Default render layout (v0.1, fixed)

```
[personal]  ~/Projects/digital-patology  ⎇ master  Opus 4.7  ctx:42%
└─green─┘   └─cyan bold─────────────┘    └yellow─┘ └─dim─┘    └─by%─┘
```

- Profile prefix: bold + color from spec
- `short_cwd`: bold cyan, `~` substitution for `$HOME`
- Git branch: yellow, only when inside a git repo
- Model: dim
- Context %: green <50, yellow <80, red ≥80

Custom layouts deferred to v0.2.

### Renderer contract

```go
type Input struct {
    Cwd        string
    ModelName  string
    ContextPct float64
}

type ProfileMeta struct {
    Name  string
    Label string
    Color string
}

func Render(in Input, p ProfileMeta) string
```

Pure function, no I/O. Testable in isolation. I/O (stdin read, lock.json read, stdout write) lives in `cmd/claude-profile/statusline.go`.

### Settings.json wiring

After `init`, the profile's `settings.json` contains:

```json
{
  "statusLine": { "type": "command", "command": "claude-profile statusline" },
  "enabledPlugins": { ... },
  "extraKnownMarketplaces": { ... }
}
```

A single binary on `$PATH` works in any profile. No per-profile bash scripts, no absolute paths to other profiles' files.

### Edge cases

| Case | Behavior |
|---|---|
| `profile.lock.json` missing | Skip prefix, render rest |
| Color in lock invalid | Render prefix without color |
| stdin not valid JSON | Render only prefix (best-effort) |
| Not a git repo | Skip branch |
| No `model.display_name` in stdin | Skip model |

## Init flow

```
1. Parse + validate profile.yaml
   ├─ yaml syntax
   ├─ schema (validator tags + custom rules: color enum, plugin@marketplace format)
   ├─ env interpolation for mcp_servers.*.env
   └─ pre-flight: `claude` is on PATH

2. Resolve target dir
   target = $HOME/.claude-<name>

3. mkdir -p target

4. Render owned files (atomic write: tmp + rename)
   ├─ settings.json
   │    statusLine block
   │    extraKnownMarketplaces from spec.marketplaces
   │    settings_overrides deep-merged on top
   ├─ .mcp.json (when mcp_servers present)
   ├─ CLAUDE.md
   └─ profile.lock.json

5. Plugin install loop  (sequential, not parallel)
   for plugin in spec.plugins:
     CLAUDE_CONFIG_DIR=$target claude plugin install <plugin>
     ├─ idempotent at claude CLI level
     ├─ stderr captured
     └─ failure → record, continue, return non-zero exit

6. Print summary + next steps
```

### File ownership

| File | Status |
|---|---|
| `settings.json` | owned, overwritten |
| `.mcp.json` | owned, overwritten |
| `CLAUDE.md` | owned, overwritten |
| `profile.lock.json` | owned, overwritten |
| `.credentials.json` | never touched |
| `history.jsonl` | never touched |
| `projects/`, `sessions/`, `shell-snapshots/`, `paste-cache/` | never touched |
| `plugins/*` (incl. `installed_plugins.json`) | managed by `claude` CLI |

Marketplaces go into `settings.json` `extraKnownMarketplaces` rather than `plugins/known_marketplaces.json` (internal `claude` state).

### Atomic writes

```go
func writeAtomic(path string, content []byte) error {
    tmp := path + ".tmp"
    if err := os.WriteFile(tmp, content, 0644); err != nil {
        return err
    }
    return os.Rename(tmp, path)
}
```

A crashed mid-init leaves the previous file intact.

## Error handling

| Class | When | Reaction | Exit code |
|---|---|---|---|
| Config validation | yaml parse, schema, env interpolation | Aggregate all errors, print list, exit before any write | 2 |
| Pre-flight | `claude` not on PATH, `$HOME` unwritable | Exit before any write | 1 |
| File write | Cannot write owned file | Exit. Atomic write keeps old file intact | 1 |
| Plugin install | Specific plugin fails | Continue with remaining plugins, log failure | 3 |
| Subprocess timeout | `claude plugin install` hangs > 60s | Kill, count as failure | 3 |

### Sample validation error

```
✗ profile.yaml validation failed:

  • statusline.color: must be one of [red green yellow blue magenta
    cyan white], got "neon"

  • plugins[1]: invalid format, expected "name@marketplace",
    got "atlassian"

  • mcp_servers.atlassian.env.ATLASSIAN_USERNAME: environment variable
    is not set (used in mcp_servers config)

Hint: see docs/profile-yaml-reference.md
```

All errors aggregated in one pass — no fix-one-find-next loop.

### Sample partial plugin failure

```
✓ Created ~/.claude-personal/
✓ Wrote settings.json
✓ Wrote CLAUDE.md
✓ Wrote profile.lock.json
✓ Installing 3 plugins via claude CLI...
  ✓ atlassian@claude-plugins-official
  ✗ superpowers@claude-plugins-official
      stderr: marketplace 'claude-plugins-official' returned 404 for
      plugin 'superpowers'
  ✓ clangd-lsp@claude-plugins-official
✗ 1 of 3 plugin installs failed.
  Profile structure is valid; re-run `init` to retry failed plugins.

Exit code: 3
```

### Idempotency guarantees

Re-running `init` on an existing profile:

- Owned files: rewritten byte-for-byte if spec unchanged
- Plugins: `claude plugin install` is no-op when already installed
- User data (credentials, history, sessions): untouched

Safe to re-run after editing yaml to apply changes.

### Verbose mode

`-v` / `--verbose`: prints subprocess stdout/stderr from plugin installs. Default: only ✓/✗ summary plus failure stderr on error.

## Testing

| Layer | Scope | Mechanism |
|---|---|---|
| Unit per package | Pure logic in each `internal/*` | `go test ./internal/...` |
| Integration: init flow | Full `init` with fake `claude` on PATH | Temp `$HOME`, fake binary in `$PATH` |
| E2E | Real `claude plugin install` | Manual smoke before release |

**Per-package coverage:**

```
internal/spec/        valid/invalid yaml parse
                      validation rules: color enum, plugin format, required fields
                      env interpolation: success / missing var

internal/render/      settings.json with/without mcp_servers, settings_overrides
                      deep-merge precedence (overrides > defaults)
                      CLAUDE.md content
                      profile.lock.json shape

internal/statusline/  full input → expected ANSI
                      each color → correct escape
                      missing fields → graceful skip
                      ProfileMeta zero value → no prefix

internal/shell/       zsh/bash/fish output snapshot
                      multiple profiles enumerated

internal/profile/     discover ~/.claude-*, parse profile.lock.json
                      active detection from CLAUDE_CONFIG_DIR

internal/plugin/      subprocess wiring (via fake binary)
```

**Fake `claude` pattern (integration tests):**

The test compiles a small helper binary, places it on a temp `$PATH`, and asserts our subprocess calls it with the expected arguments.

```go
func TestInit_callsClaudePluginInstall(t *testing.T) {
    fakeClaude := buildFakeClaude(t)
    t.Setenv("PATH", filepath.Dir(fakeClaude)+":"+os.Getenv("PATH"))
    home := t.TempDir()
    t.Setenv("HOME", home)

    err := initProfile("personal", "testdata/spec.yaml")
    require.NoError(t, err)

    calls := readFakeClaudeLog(t)
    require.Equal(t, 3, len(calls))
    require.Contains(t, calls[0], "plugin install atlassian@claude-plugins-official")
}
```

**CI gate:** ≥75% coverage on `internal/*`. CLI command code (`cmd/`) is not counted — it is a thin wrapper.

## Distribution

**Channels:**

1. `brew install <user>/tap/claude-profile` — primary on macOS/Linux
2. `curl -fsSL https://raw.githubusercontent.com/<user>/claude-profile/main/install.sh | sh` — CI installations
3. `go install github.com/<user>/claude-profile/cmd/claude-profile@latest` — Go developers
4. GitHub Releases tarball — manual

**Platforms (`goreleaser`):**

| OS | Arch | Status |
|---|---|---|
| macOS | arm64, amd64 | first-class |
| Linux | arm64, amd64 | first-class |
| Windows | amd64 | best-effort, with VS Code limitation noted |

**Versioning:**

- Semver from start
- v0.1.0 = MVP described in this spec
- v0.x.y while features stabilise
- v1.0.0 freezes the yaml schema (deprecation cycle required from then on)

## Roadmap (post-v0.1)

- `doctor` — health check (plugins installed, MCP responsive, credentials valid)
- `claude_md_include` — base + overlay rendering, removes manual CLAUDE.md duplication
- `apply` — re-render owned files without re-running plugin install
- `extends:` — profile inheritance
- Plugin pinning: `name@marketplace@version`
- Custom statusline templates
- Secrets from macOS Keychain / `pass` / 1Password
- VS Code integration helper for [#30538](https://github.com/anthropics/claude-code/issues/30538)
- Generated zsh/bash completion

## README outline (positioning)

```
# claude-profile
> Profile-as-code for Claude Code. Define organizations as YAML, version them, share with your team.

## Why                       — the gap (no organization layer)
## Quickstart                — three commands, working profile
## Example profile.yaml      — the hook
## Status                    — v0.1 scope, what's in/out
## Comparison                — vs claude-code-profiles, claude-account-switcher
## Install                   — brew, curl, go install
## FAQ                       — statusline, secrets, VS Code
## Roadmap                   — v0.2: doctor, includes, reconcile
```

The comparison table is critical — explains exactly why another tool exists and what gap it fills.
