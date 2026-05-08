# FAQ

## Why does VS Code ignore my profile?

The Claude Code VS Code extension reads from `~/.claude/` regardless of the `CLAUDE_CONFIG_DIR` environment variable. Tracked upstream as [anthropics/claude-code#30538](https://github.com/anthropics/claude-code/issues/30538). `claude-profile` only affects the CLI; until the extension is fixed, profiles are CLI-only.

## Is `init` safe to re-run?

Yes. `init` is idempotent on the files it owns and never touches user data:

| Owned (overwritten) | Not touched |
|---|---|
| `settings.json` | `.credentials.json` |
| `.mcp.json` | `history.jsonl` |
| `CLAUDE.md` | `projects/`, `sessions/`, `shell-snapshots/` |
| `profile.lock.json` | `plugins/*` (managed by `claude` CLI) |

Re-running `init` after editing the yaml applies the new content without losing auth or history.

## How do I delete a profile?

```bash
rm -rf ~/.claude-<name>
```

`claude-profile` does not provide a `delete` subcommand in v0.1 to keep the destructive surface area zero. After removal, regenerate aliases:

```bash
eval "$(claude-profile shell-init zsh)"
```

## How do I remove a plugin?

`claude-profile init` is purely additive — removing a plugin from yaml does not uninstall it. Either:

1. Run `claude /plugin uninstall <name>` inside the profile, or
2. Delete the profile directory and re-init.

## Where do secrets go?

Only environment variable interpolation in `mcp_servers.*.env`:

```yaml
mcp_servers:
  atlassian:
    env:
      TOKEN: ${ATLASSIAN_TOKEN}
```

`${ATLASSIAN_TOKEN}` is expanded at `init` time from your shell environment. The yaml file itself stays free of secrets and is safe to commit.

Keychain / `pass` / 1Password integration is on the v0.2 roadmap.

## What if `claude plugin install` is missing?

Install Claude Code 2.x or newer. Earlier versions did not expose the `plugin install` subcommand. `claude-profile` requires it for the bootstrap flow.

## What does `current` return when I have no `CLAUDE_CONFIG_DIR` set?

`default` — the conventional `~/.claude/` profile. Setting `CLAUDE_CONFIG_DIR=$HOME/.claude-personal` makes `current` return `personal`.

## Does `init` validate my plugins exist before trying to install?

No. Validation is structural (yaml schema, env interpolation). Plugin reachability is checked by `claude plugin install` itself; failures from that subprocess are reported with stderr and counted toward the partial-failure exit code (3).

## How are profile names resolved on disk?

`<name>` → `~/.claude-<name>`. The exception is `default`, which maps to `~/.claude/` (the conventional location).
