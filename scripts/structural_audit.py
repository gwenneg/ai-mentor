#!/usr/bin/env python3
"""Deterministic structural audit for the ai-mentor catalog.

Checks goal files, approach files, cross-references, and SKILL.md consistency
against the structure defined in .claude/skills/ai-mentor-update/templates/.
Exits non-zero if any issue is found. No network, no LLM — safe as a PR gate.
"""

import re
import sys
from pathlib import Path

REPO = Path(__file__).resolve().parent.parent
SKILL_DIR = REPO / "skills" / "mentor"
GOALS = SKILL_DIR / "goals"
APPROACHES = SKILL_DIR / "approaches"
SKILL_MD = SKILL_DIR / "SKILL.md"

APPROACH_SECTIONS = [
    "## What It Is",
    "## Why It Works",
    "## When to Use It",
    "## When NOT to Use It",
    "## How It Works",
    "### Basic (Beginner)",
    "### Composing with Other Approaches (Intermediate)",
    "### Advanced Patterns",
    "## Common Pitfalls",
    "## Real-World Example",
    "## Sources",
]
GOAL_SECTIONS = [
    "## When You're Here",
    "## Quick Decision Guide",
    "**Hidden gem:**",
    "## Approaches (Ranked)",
]
ENTRY_FIELDS = [
    "**Level:**",
    "**Try it now:**",
    "**Why this works:**",
    "**Pros:**",
    "**Cons:**",
    "**Deeper:** See `approaches/",
]
APPROACH_LINE_RANGE = (60, 110)
GOAL_ENTRY_RANGE = (3, 7)
DATE_RE = re.compile(r"^\*Last (?:reviewed|verified): \d{4}-\d{2}-\d{2}\*", re.M)

issues: list[str] = []


def issue(path: Path, msg: str) -> None:
    issues.append(f"{path.relative_to(REPO)}: {msg}")


def check_order(path: Path, text: str, sections: list[str]) -> None:
    pos = -1
    for s in sections:
        found = text.find(s)
        if found == -1:
            issue(path, f"missing section {s!r}")
        elif found < pos:
            issue(path, f"section {s!r} out of order")
        else:
            pos = found


def check_date_line(path: Path, text: str) -> None:
    lines = text.splitlines()
    if len(lines) < 2 or not DATE_RE.match(lines[1]):
        issue(path, "line 2 must be '*Last reviewed: YYYY-MM-DD*'")


def check_goal(path: Path) -> None:
    text = path.read_text()
    check_date_line(path, text)
    check_order(path, text, GOAL_SECTIONS)

    gems = re.findall(r"\*\*Hidden gem:\*\* ([^—\n]+)", text)
    headers = re.findall(r"^### \d+\. ([^—\n]+)", text, re.M)
    if len(gems) != 1:
        issue(path, f"expected exactly one Hidden gem line, found {len(gems)}")
    elif not any(
        gems[0].strip().lower() in h.lower() or h.strip().lower() in gems[0].strip().lower()
        for h in headers
    ):
        issue(path, f"Hidden gem {gems[0].strip()!r} does not match any ranked approach")

    nums = [int(n) for n in re.findall(r"^### (\d+)\.", text, re.M)]
    if nums != list(range(1, len(nums) + 1)):
        issue(path, f"approach numbering not sequential: {nums}")
    if not GOAL_ENTRY_RANGE[0] <= len(nums) <= GOAL_ENTRY_RANGE[1]:
        issue(path, f"{len(nums)} approach entries (expected {GOAL_ENTRY_RANGE[0]}-{GOAL_ENTRY_RANGE[1]})")

    ranked = text.split("## Approaches (Ranked)", 1)[-1]
    entries = re.split(r"^---$", ranked, flags=re.M)
    entries = [e for e in entries if re.search(r"^### \d+\.", e, re.M)]
    for e in entries:
        header = re.search(r"^### \d+\. .*", e, re.M).group(0)
        pos = -1
        for f in ENTRY_FIELDS:
            found = e.find(f)
            if found == -1:
                issue(path, f"entry {header!r}: missing field {f!r}")
            elif found < pos:
                issue(path, f"entry {header!r}: field {f!r} out of order")
            else:
                pos = found
    if ranked.rstrip().endswith("---"):
        issue(path, "trailing '---' after the last approach entry")

    for ref in re.findall(r"approaches/[a-z0-9-]+\.md", text):
        if not (SKILL_DIR / ref).exists():
            issue(path, f"broken reference {ref}")


def check_approach(path: Path) -> None:
    text = path.read_text()
    check_date_line(path, text)
    check_order(path, text, APPROACH_SECTIONS)
    n = len(text.splitlines())
    lo, hi = APPROACH_LINE_RANGE
    if not lo <= n <= hi:
        issue(path, f"{n} lines (expected {lo}-{hi})")
    sources = re.findall(r"^- \[[^\]]+\]\(https?://[^)]+\)", text.split("## Sources", 1)[-1], re.M)
    if not 1 <= len(sources) <= 3:
        issue(path, f"{len(sources)} Sources entries (expected 1-3)")


def check_cross_references(goal_files: list[Path], approach_files: list[Path]) -> None:
    all_goal_text = "".join(p.read_text() for p in goal_files)
    for a in approach_files:
        if f"approaches/{a.name}" not in all_goal_text:
            issue(a, "orphan: not referenced by any goal file")


def check_skill_md(goal_files: list[Path]) -> None:
    text = SKILL_MD.read_text()
    goal_names = {p.name for p in goal_files}
    table_names = set(re.findall(r"^\| `([a-z0-9-]+\.md)` \|", text, re.M))
    for missing in sorted(goal_names - table_names):
        issue(SKILL_MD, f"goal file {missing} missing from the Phase 1 classification table")
    for stale in sorted(table_names - goal_names):
        issue(SKILL_MD, f"Phase 1 table row {stale} has no matching goal file")
    for count in re.findall(r"(\d+) goal categories", text):
        if int(count) != len(goal_names):
            issue(SKILL_MD, f"prose says '{count} goal categories' but there are {len(goal_names)}")


def main() -> int:
    goal_files = sorted(GOALS.glob("*.md"))
    approach_files = sorted(APPROACHES.glob("*.md"))
    if not goal_files or not approach_files:
        print("FATAL: goal or approach directory empty/missing", file=sys.stderr)
        return 2
    for p in goal_files:
        check_goal(p)
    for p in approach_files:
        check_approach(p)
    check_cross_references(goal_files, approach_files)
    check_skill_md(goal_files)

    print(f"Audited {len(goal_files)} goals, {len(approach_files)} approaches.")
    if issues:
        print(f"\n{len(issues)} issue(s):")
        for i in issues:
            print(f"  - {i}")
        return 1
    print("Structural audit: PASS")
    return 0


if __name__ == "__main__":
    sys.exit(main())
