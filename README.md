# claude-profile

> Profile-as-code for Claude Code. Define organizations as YAML, version them, share with your team.

## Why

Claude Code natively has three configuration layers — user, project, project-local. There is no **organization layer**. If you work across multiple orgs on one machine, each with its own plugins, MCP servers, branch rules, and statusline conventions, you end up either logging in and out, manually editing `~/.claude/`, or running shell aliases without any reproducibility.

`claude-profile` fills that gap. A single declarative `profile.yaml` describes everything an org needs; `claude-profile new` bootstraps an isolated `~/.claude-<name>/` from it. The yaml is small enough to commit to your team repo — onboarding becomes one command.

## Quickstart

```bash
# Install (from source for now)
go install github.com/dmitriipyshinskii/claude-profile/cmd/claude-profile@latest

# Bootstrap a profile from an embedded template
claude-profile new personal -t personal
# → output ends with the exact `eval` and `>> ~/.zshrc` lines you need

# Activate the alias for this shell
eval "$(claude-profile init zsh)"

# Persist for new shells
echo 'eval "$(claude-profile init zsh)"' >> ~/.zshrc

# Launch Claude Code in the new profile
claude-personal
```

## Example `profile.yaml`

```yaml
apiVersion: claude-profile.io/v1
kind: Profile

metadata:
  name: personal

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
  Branch prefix: feature/, bugfix/, hotfix/.
```

`claude-profile new personal -f profile.yaml` produces `~/.claude-personal/` with `settings.json`, `CLAUDE.md`, `profile.lock.json`, and the plugins installed.

## Comparison

| | claude-profile | [claude-code-profiles](https://github.com/quinnjr/claude-code-profiles) | [claude-account-switcher](https://github.com/ukogan/claude-account-switcher) |
|---|---|---|---|
| Declarative `profile.yaml` | ✅ | ❌ | ❌ |
| Shareable spec (commit to repo) | ✅ | ❌ | ❌ |
| Templated bootstrap | ✅ | ❌ | ❌ |
| Per-profile statusline with visual identifier | ✅ | ❌ | ❌ |
| Composable CLAUDE.md | v0.2 | ❌ | ❌ |
| Full isolation per profile | ✅ | ✅ | partial (symlinked shared config) |
| Single binary | ✅ | shell scripts | shell scripts |
| `doctor` health check | v0.2 | ❌ | ❌ |

## CLI

```
claude-profile init <shell>                                       # activate shell integration (zsh|bash|fish)
claude-profile new  <name> [-f profile.yaml | -t <template>] [--dry-run]
claude-profile list  [--all] [--json]
claude-profile current
claude-profile which <name>
claude-profile templates
claude-profile statusline                                         # internal, called by Claude Code
```

## Status

v0.1 — initial release. Scope:

- Bootstrap profiles from yaml or embedded templates
- Per-profile statusline with named-color prefix
- Shell alias generation for zsh, bash, fish
- Plugin install via `claude` CLI subprocess
- MCP server config with `${VAR}` env interpolation
- Settings overrides with deep-merge

Out of scope for v0.1: drift reconciliation, plugin uninstall, secret managers, custom statusline templates, profile inheritance, VS Code integration. See [Roadmap](#roadmap).

## Install

### `go install`

```bash
go install github.com/dmitriipyshinskii/claude-profile/cmd/claude-profile@latest
```

### From release tarball

Grab the platform-appropriate archive from [Releases](https://github.com/dmitriipyshinskii/claude-profile/releases) and put `claude-profile` on your `$PATH`.

### Build locally

```bash
git clone https://github.com/dmitriipyshinskii/claude-profile
cd claude-profile
make build
```

## Development

If your local environment blocks execution of newly-built binaries (EDR, Gatekeeper), `scripts/dev.sh` runs Go inside a `golang:1.26-alpine` container with persistent caches:

```bash
scripts/dev.sh go test ./...
scripts/dev.sh go build -o /work/bin/claude-profile ./cmd/claude-profile
```

## Documentation

- [`docs/profile-yaml-reference.md`](docs/profile-yaml-reference.md) — full schema
- [`docs/faq.md`](docs/faq.md) — VS Code limitation, idempotency, plugin removal, secrets
- [`examples/`](examples/) — working profile specs
- [`docs/design/`](docs/design/) — original design doc

## Roadmap

- `doctor` health check (plugins installed, MCP responsive, credentials valid)
- `claude_md_include` — base + overlay rendering
- `apply` — re-render owned files without re-running plugin install
- `extends:` — profile inheritance
- Plugin pinning: `name@marketplace@version`
- Custom statusline templates
- Secrets from macOS Keychain / `pass` / 1Password
- VS Code integration helper for [anthropics/claude-code#30538](https://github.com/anthropics/claude-code/issues/30538)
- Generated zsh/bash completion

## License

MIT — see [`LICENSE`](LICENSE).
