#!/usr/bin/env bash
# Deterministic structural audit for the ai-mentor catalog.
#
# Checks goal files, approach files, cross-references, and SKILL.md
# consistency. Exits non-zero if any issue is found. No network, no LLM —
# safe as a PR gate. Portable: bash + POSIX grep/sed/awk only.

set -u

REPO="$(cd "$(dirname "$0")/.." && pwd)"
SKILL_DIR="$REPO/skills/mentor"
GOALS="$SKILL_DIR/goals"
APPROACHES="$SKILL_DIR/approaches"
SKILL_MD="$SKILL_DIR/SKILL.md"

ISSUES=0
issue() { # $1 = file, $2 = message
  printf '  - %s: %s\n' "${1#"$REPO"/}" "$2"
  ISSUES=$((ISSUES + 1))
}

# first_line <file> <fixed-string>  -> line number of first occurrence, or empty
first_line() { grep -n -F -- "$2" "$1" | head -1 | cut -d: -f1; }

check_order() { # $1 = file, remaining args = sections in required order
  local f="$1"; shift
  local pos=0 ln s
  for s in "$@"; do
    ln="$(first_line "$f" "$s")"
    if [ -z "$ln" ]; then
      issue "$f" "missing section '$s'"
    elif [ "$ln" -lt "$pos" ]; then
      issue "$f" "section '$s' out of order"
    else
      pos="$ln"
    fi
  done
}

check_date_line() { # $1 = file
  sed -n '2p' "$1" | grep -qE '^\*Last (reviewed|verified): [0-9]{4}-[0-9]{2}-[0-9]{2}\*' \
    || issue "$1" "line 2 must be '*Last reviewed: YYYY-MM-DD*'"
}

lc() { printf '%s' "$1" | tr '[:upper:]' '[:lower:]'; }

# ---------- goal files ----------
check_goal() {
  local f="$1"
  check_date_line "$f"
  check_order "$f" "## When You're Here" "## Quick Decision Guide" "**Hidden gem:**" "## Approaches (Ranked)"

  # Hidden gem must name a ranked approach (containment, case-insensitive)
  local gem_count gem headers h matched=0
  gem_count="$(grep -c -F '**Hidden gem:**' "$f")"
  if [ "$gem_count" -ne 1 ]; then
    issue "$f" "expected exactly one Hidden gem line, found $gem_count"
  else
    gem="$(grep -F '**Hidden gem:**' "$f" | sed -e 's/.*\*\*Hidden gem:\*\* //' -e 's/ —.*//' -e 's/[[:space:]]*$//')"
    headers="$(grep -E '^### [0-9]+\. ' "$f" | sed -e 's/^### [0-9]*\. //' -e 's/ —.*//' -e 's/[[:space:]]*$//')"
    while IFS= read -r h; do
      case "$(lc "$h")" in *"$(lc "$gem")"*) matched=1 ;; esac
      case "$(lc "$gem")" in *"$(lc "$h")"*) matched=1 ;; esac
    done <<EOF
$headers
EOF
    [ "$matched" -eq 1 ] || issue "$f" "Hidden gem '$gem' does not match any ranked approach"
  fi

  # Numbering sequential, at least 3 entries
  local nums expected n_entries
  nums="$(grep -E '^### [0-9]+\.' "$f" | sed -E 's/^### ([0-9]+)\..*/\1/' | tr '\n' ' ')"
  n_entries="$(printf '%s' "$nums" | wc -w | tr -d ' ')"
  expected="$(seq 1 "$n_entries" | tr '\n' ' ')"
  [ "$nums" = "$expected" ] || issue "$f" "approach numbering not sequential: $nums"
  if [ "$n_entries" -lt 3 ]; then
    issue "$f" "$n_entries approach entries (expected at least 3)"
  fi

  # Per-entry required fields in order (entries = blocks between '---' after the ranked header)
  awk -v file="${f#"$REPO"/}" '
    BEGIN {
      nf = 6
      fields[1] = "**Level:**";           fields[2] = "**Try it now:**"
      fields[3] = "**Why this works:**";  fields[4] = "**Pros:**"
      fields[5] = "**Cons:**";            fields[6] = "**Deeper:** See `approaches/"
    }
    /^## Approaches \(Ranked\)$/ { ranked = 1; next }
    !ranked { next }
    /^---$/ { flush(); next }
    {
      if ($0 ~ /^### [0-9]+\./) { header = $0; nlines = 0 }
      if (header != "") { nlines++; lines[nlines] = $0 }
    }
    function flush(   i, j, pos, found) {
      if (header == "") return
      pos = 0
      for (i = 1; i <= nf; i++) {
        found = 0
        for (j = 1; j <= nlines; j++) if (index(lines[j], fields[i])) { found = j; break }
        if (!found)          printf "  - %s: entry %c%s%c missing field %c%s%c\n", file, 39, header, 39, 39, fields[i], 39
        else if (found < pos) printf "  - %s: entry %c%s%c field %c%s%c out of order\n", file, 39, header, 39, 39, fields[i], 39
        else pos = found
      }
      header = ""; nlines = 0
    }
    END { flush() }
  ' "$f" | {
    local out; out="$(cat)"
    if [ -n "$out" ]; then printf '%s\n' "$out"; return 1; fi
  } || ISSUES=$((ISSUES + 1))

  # No trailing '---' after the last entry
  [ "$(grep -vE '^[[:space:]]*$' "$f" | tail -1)" = "---" ] \
    && issue "$f" "trailing '---' after the last approach entry"

  # All approach references resolve
  local ref
  for ref in $(grep -oE 'approaches/[a-z0-9-]+\.md' "$f" | sort -u); do
    [ -f "$SKILL_DIR/$ref" ] || issue "$f" "broken reference $ref"
  done
}

# ---------- approach files ----------
check_approach() {
  local f="$1"
  check_date_line "$f"
  check_order "$f" \
    "## What It Is" "## Why It Works" "## When to Use It" "## When NOT to Use It" \
    "## How It Works" "### Basic (Beginner)" \
    "### Composing with Other Approaches (Intermediate)" "### Advanced Patterns" \
    "## Common Pitfalls" "## Real-World Example" "## Sources"

  local n srcs
  n="$(wc -l < "$f" | tr -d ' ')"
  if [ "$n" -lt 60 ]; then
    issue "$f" "$n lines (expected at least 60)"
  fi
  srcs="$(awk '/^## Sources$/{s=1;next} s' "$f" | grep -cE '^- \[[^]]+\]\(https?://')"
  if [ "$srcs" -lt 1 ]; then
    issue "$f" "$srcs Sources entries (expected at least 1)"
  fi
}

# ---------- cross-references ----------
check_orphans() {
  local a base
  for a in "$APPROACHES"/*.md; do
    base="approaches/$(basename "$a")"
    grep -rqF -- "$base" "$GOALS"/ || issue "$a" "orphan: not referenced by any goal file"
  done
}

# ---------- SKILL.md consistency ----------
check_skill_md() {
  local goal_names table_names missing stale count n_goals
  n_goals="$(ls "$GOALS"/*.md | wc -l | tr -d ' ')"
  goal_names="$(ls "$GOALS" | sort)"
  table_names="$(grep -oE '^\| `[a-z0-9-]+\.md` \|' "$SKILL_MD" | sed -e 's/^| `//' -e 's/` |$//' | sort)"

  missing="$(comm -23 <(printf '%s\n' "$goal_names") <(printf '%s\n' "$table_names"))"
  stale="$(comm -13 <(printf '%s\n' "$goal_names") <(printf '%s\n' "$table_names"))"
  local x
  for x in $missing; do issue "$SKILL_MD" "goal file $x missing from the Phase 1 classification table"; done
  for x in $stale;   do issue "$SKILL_MD" "Phase 1 table row $x has no matching goal file"; done

  for count in $(grep -oE '[0-9]+ goal categories' "$SKILL_MD" | grep -oE '^[0-9]+'); do
    [ "$count" -eq "$n_goals" ] \
      || issue "$SKILL_MD" "prose says '$count goal categories' but there are $n_goals"
  done
}

# ---------- main ----------
n_goals="$(ls "$GOALS"/*.md 2>/dev/null | wc -l | tr -d ' ')"
n_apprs="$(ls "$APPROACHES"/*.md 2>/dev/null | wc -l | tr -d ' ')"
if [ "$n_goals" -eq 0 ] || [ "$n_apprs" -eq 0 ]; then
  echo "FATAL: goal or approach directory empty/missing" >&2
  exit 2
fi

for f in "$GOALS"/*.md;      do check_goal "$f";     done
for f in "$APPROACHES"/*.md; do check_approach "$f"; done
check_orphans
check_skill_md

echo "Audited $n_goals goals, $n_apprs approaches."
if [ "$ISSUES" -gt 0 ]; then
  echo ""
  echo "$ISSUES issue(s) found (listed above)."
  exit 1
fi
echo "Structural audit: PASS"
