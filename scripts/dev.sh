#!/usr/bin/env bash
set -euo pipefail

# Run a command inside a golang toolchain container with the project mounted
# and persistent caches via named volumes. Use this when local Go execution is
# blocked by EDR/sandbox.
#
# Usage:
#   scripts/dev.sh go test ./...
#   scripts/dev.sh go build ./cmd/claude-profile
#   scripts/dev.sh sh -c "go vet ./... && go test ./..."

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
IMAGE="${CLAUDE_PROFILE_DEV_IMAGE:-golang:1.26-alpine}"

exec docker run --rm \
  -v "${REPO_ROOT}:/work" \
  -w /work \
  -v claude-profile-gocache:/root/.cache/go-build \
  -v claude-profile-gomodcache:/go/pkg/mod \
  -e CGO_ENABLED=0 \
  "${IMAGE}" \
  "$@"
