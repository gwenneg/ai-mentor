#!/usr/bin/env bash
# Deterministic structural audit for the ai-mentor catalog.
#
# Implemented in Go (see scripts/audit/); this wrapper keeps the stable
# entry point used by CI, the maintenance allowlist, and REVIEW.md.
# No network, no LLM — safe as a PR gate. Requires a Go toolchain.
set -euo pipefail

REPO="$(cd "$(dirname "$0")/.." && pwd)"
exec go -C "$REPO/scripts/audit" run . "$REPO"
