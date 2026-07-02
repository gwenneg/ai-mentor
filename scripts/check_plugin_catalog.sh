#!/usr/bin/env bash
# Check references/official-plugins.md against the live official marketplace.
#
# Pure GitHub API diffing — no LLM. Exits non-zero on drift so a scheduled
# workflow can open an issue or feed the diff to the maintenance run.
# Requires: curl, jq. Uses GITHUB_TOKEN / GH_TOKEN if set (else unauthenticated,
# subject to rate limits).

set -euo pipefail

REPO="$(cd "$(dirname "$0")/.." && pwd)"
CATALOG="$REPO/skills/mentor/references/official-plugins.md"
API="https://api.github.com/repos/anthropics/claude-plugins-official/contents"

auth_args=()
token="${GITHUB_TOKEN:-${GH_TOKEN:-}}"
[ -n "$token" ] && auth_args=(-H "Authorization: Bearer $token")

fetch_names() { # $1 = directory
  curl -sfL --max-time 30 "${auth_args[@]}" \
    -H "Accept: application/vnd.github+json" "$API/$1" \
    | jq -r '.[] | select(.type == "dir") | .name'
}

live="$( (fetch_names plugins; fetch_names external_plugins) | sort -u)"
documented="$(grep -oE '`[a-z0-9-]+`' "$CATALOG" | tr -d '\`' | sort -u)"

live_count="$(printf '%s\n' "$live" | grep -c .)"
in_both="$(comm -12 <(printf '%s\n' "$live") <(printf '%s\n' "$documented") | grep -c . || true)"
echo "Marketplace plugins: $live_count; documented names found: $in_both"

missing="$(comm -23 <(printf '%s\n' "$live") <(printf '%s\n' "$documented"))"
# Documented multi-word kebab tokens not in the marketplace (may be prose tokens)
removed="$(comm -13 <(printf '%s\n' "$live") <(printf '%s\n' "$documented") \
  | grep -E '^[a-z0-9]+(-[a-z0-9]+)+$' || true)"

drift=0
if [ -n "$missing" ]; then
  drift=1
  echo ""
  echo "NEW plugins not yet in the catalog:"
  printf '  + %s\n' $missing
fi
if [ -n "$removed" ]; then
  echo ""
  echo "Documented names not in the marketplace (verify manually — may be prose tokens):"
  printf '  ? %s\n' $removed
fi

if [ "$drift" -eq 1 ]; then
  echo ""
  echo "Drift detected: run the maintenance skill's catalog sync (step 5)."
  exit 1
fi
echo ""
echo "Plugin catalog: in sync."
