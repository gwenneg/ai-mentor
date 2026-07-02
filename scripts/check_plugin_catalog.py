#!/usr/bin/env python3
"""Check references/official-plugins.md against the live official marketplace.

Compares the plugin names documented in the catalog reference file with the
directories in anthropics/claude-plugins-official (plugins/ and
external_plugins/). Pure GitHub API diffing — no LLM. Exits non-zero on drift
so a scheduled workflow can open an issue or feed the diff to the maintenance
run.

Requires: GITHUB_TOKEN or gh auth (uses unauthenticated API as fallback,
subject to rate limits).
"""

import json
import os
import re
import sys
import urllib.request
from pathlib import Path

REPO = Path(__file__).resolve().parent.parent
CATALOG = REPO / "skills" / "mentor" / "references" / "official-plugins.md"
API = "https://api.github.com/repos/anthropics/claude-plugins-official/contents/{}"


def fetch_names(directory: str) -> set[str]:
    req = urllib.request.Request(API.format(directory))
    token = os.environ.get("GITHUB_TOKEN") or os.environ.get("GH_TOKEN")
    if token:
        req.add_header("Authorization", f"Bearer {token}")
    req.add_header("Accept", "application/vnd.github+json")
    with urllib.request.urlopen(req, timeout=30) as resp:
        entries = json.load(resp)
    return {e["name"] for e in entries if e["type"] == "dir"}


def main() -> int:
    live = fetch_names("plugins") | fetch_names("external_plugins")
    documented = set(re.findall(r"`([a-z0-9-]+)`", CATALOG.read_text()))

    missing = sorted(live - documented)   # in marketplace, not in catalog
    removed = sorted(p for p in documented - live
                     # catalog also backticks non-plugin tokens (commands, files)
                     if re.fullmatch(r"[a-z0-9]+(-[a-z0-9]+)+", p))

    print(f"Marketplace plugins: {len(live)}; documented names found: {len(documented & live)}")
    drift = False
    if missing:
        drift = True
        print("\nNEW plugins not yet in the catalog:")
        for p in missing:
            print(f"  + {p}")
    if removed:
        print("\nDocumented names not in the marketplace (verify manually — may be prose tokens):")
        for p in removed:
            print(f"  ? {p}")

    if drift:
        print("\nDrift detected: run the maintenance skill's catalog sync (step 5).")
        return 1
    print("\nPlugin catalog: in sync.")
    return 0


if __name__ == "__main__":
    sys.exit(main())
