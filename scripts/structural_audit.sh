#!/usr/bin/env bash
# Deterministic structural audit for the ai-mentor catalog.
#
# Checks goal files, approach files, cross-references, and SKILL.md
# consistency. Exits non-zero if any issue is found. No network, no LLM —
# safe as a PR gate. Portable: bash + POSIX grep/sed/awk only.

set -u

REPO="$(cd "$(dirname "$0")/.." && pwd)"
SKILL_DIR="$REPO/skills/mentor"
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
  sed -n '2p' "$1" | grep -qE '^\*Last verified: [0-9]{4}-[0-9]{2}-[0-9]{2}\*' \
    || issue "$1" "line 2 must be '*Last verified: YYYY-MM-DD*'"
}

lc() { printf '%s' "$1" | tr '[:upper:]' '[:lower:]'; }

# ---------- routing table ----------
ROUTING="$SKILL_DIR/routing.md"

check_routing() {
  if [ ! -f "$ROUTING" ]; then
    issue "$ROUTING" "missing routing table"
    return
  fi
  check_date_line "$ROUTING"

  local out
  out="$(
    awk -v file="${ROUTING#"$REPO"/}" '
      /^## / {
        if (section != "" && section != "extraction-notes") flush()
        section = substr($0, 4); n = 0; gem = 0; next
      }
      section == "" || section == "extraction-notes" { next }
      /^\*\*Hidden gem:\*\*/ { gem = 1; gemline = $0 }
      /^\| [0-9]+ \|/ {
        n++
        split($0, cells, "|")
        num = cells[2]; gsub(/ /, "", num)
        if (num != n) { printf "section %s: row numbering not sequential at row %d\n", section, n }
      }
      END { if (section != "" && section != "extraction-notes") flush() }
      function flush() {
        if (n < 3) printf "section %s: only %d rows (expected at least 3)\n", section, n
        if (!gem)  printf "section %s: missing Hidden gem line\n", section
      }
    ' "$ROUTING"
  )"
  if [ -n "$out" ]; then
    printf '%s\n' "$out" | while IFS= read -r line; do
      printf '  - %s: %s\n' "${ROUTING#"$REPO"/}" "$line"
    done
    ISSUES=$((ISSUES + $(printf '%s\n' "$out" | wc -l)))
  fi

  # Every approach link resolves
  local ref
  for ref in $(grep -oE 'approaches/[a-z0-9-]+\.md' "$ROUTING" | sort -u); do
    [ -f "$SKILL_DIR/$ref" ] || issue "$ROUTING" "broken reference $ref"
  done

  # Every level valid
  grep -E '^\| [0-9]+ \|' "$ROUTING" | awk -F'|' '{gsub(/ /,"",$4); print $4}' | sort -u | while read -r lvl; do
    case "$lvl" in
      (Beginner|Intermediate|Advanced) ;;
      (*) printf '  - %s: invalid level %s\n' "${ROUTING#"$REPO"/}" "$lvl" ;;
    esac
  done | { local o; o="$(cat)"; if [ -n "$o" ]; then printf '%s\n' "$o"; ISSUES=$((ISSUES + $(printf '%s\n' "$o" | wc -l))); fi; }

  # Hidden gem names a ranked approach in its own section (containment, case-insensitive)
  local gemout
  gemout="$(python3 - "$ROUTING" <<'PYEOF'
import re, sys
text = open(sys.argv[1]).read()
sections = re.split(r"^## ", text, flags=re.M)[1:]
for sec in sections:
    name = sec.splitlines()[0].strip()
    if name == "extraction-notes":
        continue
    gem = re.search(r"^\*\*Hidden gem:\*\* ([^—\n]+)", sec, re.M)
    rows = re.findall(r"^\| \d+ \| \[([^\]]+)\]", sec, re.M)
    if gem:
        g = gem.group(1).strip().lower()
        if not any(g in r.lower() or r.lower() in g for r in rows):
            print(f"section {name}: Hidden gem '{gem.group(1).strip()}' does not match any ranked row")
PYEOF
)"
  if [ -n "$gemout" ]; then
    printf '%s\n' "$gemout" | while IFS= read -r line; do
      printf '  - %s: %s\n' "${ROUTING#"$REPO"/}" "$line"
    done
    ISSUES=$((ISSUES + $(printf '%s\n' "$gemout" | wc -l)))
  fi
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
    grep -qF -- "$base" "$ROUTING" || issue "$a" "orphan: not referenced by the routing table"
  done
}

# ---------- processed-changelog ledger ----------
LEDGER="$SKILL_DIR/references/processed-changelogs.md"

check_ledger() {
  if [ ! -f "$LEDGER" ]; then
    issue "$LEDGER" "missing processed-changelog ledger"
    return
  fi
  sed -n '2p' "$LEDGER" | grep -qE '^\*Updated: [0-9]{4}-[0-9]{2}-[0-9]{2}\*' \
    || issue "$LEDGER" "line 2 must be '*Updated: YYYY-MM-DD*'"

  local out
  out="$(
    grep -E '^\| *\[' "$LEDGER" | while IFS='|' read -r _ week date outcome _; do
      slug="$(printf '%s' "$week" | sed -E 's/^ *\[([^]]+)\].*/\1/')"
      d="$(printf '%s' "$date" | sed 's/^ *//; s/ *$//')"
      o="$(printf '%s' "$outcome" | sed 's/^ *//; s/ *$//')"
      printf '%s' "$slug" | grep -qE '^[0-9]{4}-w[0-9]{2}$' \
        || echo "row '$slug' is not a week slug like 2026-w26"
      printf '%s' "$d" | grep -qE '^[0-9]{4}-[0-9]{2}-[0-9]{2}$' \
        || echo "row '$slug' has invalid processed date '$d'"
      [ -n "$o" ] || echo "row '$slug' has an empty outcome"
    done
  )"
  if [ -n "$out" ]; then
    printf '%s\n' "$out" | while IFS= read -r line; do
      printf '  - %s: %s\n' "${LEDGER#"$REPO"/}" "$line"
    done
    ISSUES=$((ISSUES + $(printf '%s\n' "$out" | wc -l)))
  fi

  # duplicate week slugs
  local dups
  dups="$(grep -E '^\| *\[' "$LEDGER" | sed -E 's/^\| *\[([^]]+)\].*/\1/' | sort | uniq -d)"
  local d
  for d in $dups; do issue "$LEDGER" "duplicate ledger row for '$d'"; done
}

# ---------- adoption signals ----------
SIGNALS="$SKILL_DIR/references/adoption-signals.md"

check_signals() {
  if [ ! -f "$SIGNALS" ]; then
    issue "$SIGNALS" "missing adoption-signals table"
    return
  fi
  sed -n '2p' "$SIGNALS" | grep -qE '^\*Last reviewed: [0-9]{4}-[0-9]{2}-[0-9]{2}\*' \
    || issue "$SIGNALS" "line 2 must be '*Last reviewed: YYYY-MM-DD*'"

  local approach_names signal_names missing stale x
  approach_names="$(ls "$APPROACHES" | sed 's/\.md$//' | sort)"
  signal_names="$(grep -E '^\| [a-z0-9-]+ \|' "$SIGNALS" | sed -E 's/^\| ([a-z0-9-]+) \|.*/\1/' | sort)"

  missing="$(comm -23 <(printf '%s\n' "$approach_names") <(printf '%s\n' "$signal_names"))"
  stale="$(comm -13 <(printf '%s\n' "$approach_names") <(printf '%s\n' "$signal_names"))"
  for x in $missing; do issue "$SIGNALS" "approach '$x' has no adoption-signals row"; done
  for x in $stale;   do issue "$SIGNALS" "row '$x' has no matching approach file"; done

  local dups d
  dups="$(printf '%s\n' "$signal_names" | uniq -d)"
  for d in $dups; do issue "$SIGNALS" "duplicate signals row for '$d'"; done
}

# ---------- SKILL.md consistency ----------
check_skill_md() {
  local goal_names table_names missing stale count n_goals
  goal_names="$(grep -E '^## ' "$ROUTING" | sed 's/^## //' | grep -v '^extraction-notes$' | sort)"
  n_goals="$(printf '%s\n' "$goal_names" | grep -c .)"
  table_names="$(grep -oE '^\| `[a-z0-9-]+` \|' "$SKILL_MD" | sed -e 's/^| `//' -e 's/` |$//' | sort)"

  missing="$(comm -23 <(printf '%s\n' "$goal_names") <(printf '%s\n' "$table_names"))"
  stale="$(comm -13 <(printf '%s\n' "$goal_names") <(printf '%s\n' "$table_names"))"
  local x
  for x in $missing; do issue "$SKILL_MD" "routing goal $x missing from the Phase 1 classification table"; done
  for x in $stale;   do issue "$SKILL_MD" "Phase 1 table row $x has no matching routing section"; done

  for count in $(grep -oE '[0-9]+ goal categories' "$SKILL_MD" | grep -oE '^[0-9]+'); do
    [ "$count" -eq "$n_goals" ] \
      || issue "$SKILL_MD" "prose says '$count goal categories' but there are $n_goals"
  done
}

# ---------- main ----------
n_apprs="$(ls "$APPROACHES"/*.md 2>/dev/null | wc -l | tr -d ' ')"
if [ "$n_apprs" -eq 0 ]; then
  echo "FATAL: approach directory empty/missing" >&2
  exit 2
fi

check_routing
for f in "$APPROACHES"/*.md; do check_approach "$f"; done
check_orphans
check_ledger
check_signals
check_skill_md

n_weeks="$(grep -cE '^\| *\[' "$LEDGER" 2>/dev/null || echo 0)"
n_goals="$(grep -E '^## ' "$ROUTING" 2>/dev/null | grep -vc '^## extraction-notes$' || echo 0)"
echo "Audited $n_goals routing goals, $n_apprs approaches, $n_weeks processed changelogs."
if [ "$ISSUES" -gt 0 ]; then
  echo ""
  echo "$ISSUES issue(s) found (listed above)."
  exit 1
fi
echo "Structural audit: PASS"
