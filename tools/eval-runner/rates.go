// Phase 3 of EVALS_V2: the rate tier. Non-strict cases that majority-fail
// no longer red the gate automatically — the verdict consults the case's
// trailing failure rate from past release-gate records. A majority-fail
// consistent with the case's known rate is bad luck, absorbed and loudly
// labeled; one that exceeds it is a regression and stays red. Promises
// ([strict]) never enter this tier: one failed epoch remains the verdict.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func jsonUnmarshal(line string, v any) error { return json.Unmarshal([]byte(line), v) }

// caseRates is one case's trailing history.
type caseRates struct {
	Epochs int // total epochs observed across the window
	Fails  int // FAIL epochs (any layer)
}

// history maps case ID -> trailing rates. Nil means the rate tier is
// inactive (cold start or missing -history dir) and gate semantics are
// unchanged — the report says so either way.
type history map[string]caseRates

// minHistoryRuns gates the rate tier: below this many past runs the
// baseline is too thin to absorb anything (research: ~9+ epochs for a
// meaningful rate; 3 runs x 3 epochs = 9).
const minHistoryRuns = 3

// loadHistory reads every *.jsonl in dir (one file per past run, as
// downloaded from the eval-records artifacts). Malformed lines are
// skipped; a missing dir returns nil (inactive tier), any other read
// error is fatal to the caller — same doctrine as the other ground truth.
func loadHistory(dir string) (history, int, error) {
	if dir == "" {
		return nil, 0, nil
	}
	files, err := filepath.Glob(filepath.Join(dir, "*.jsonl"))
	if err != nil {
		return nil, 0, err
	}
	if len(files) < minHistoryRuns {
		return nil, len(files), nil
	}
	h := history{}
	sort.Strings(files)
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			return nil, 0, fmt.Errorf("rate history: %w", err)
		}
		for _, line := range strings.Split(string(data), "\n") {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			var rec record
			if err := jsonUnmarshal(line, &rec); err != nil || rec.Case == "" {
				continue
			}
			// First attempts only: retries measure the harness, not the case.
			if rec.Attempt > 1 {
				continue
			}
			cr := h[rec.Case]
			cr.Epochs++
			if rec.Verdict == vFail {
				cr.Fails++
			}
			h[rec.Case] = cr
		}
	}
	return h, len(files), nil
}

// binomTail returns P(X >= k) for X ~ Binomial(n, p) — exact, n is small.
func binomTail(n, k int, p float64) float64 {
	total := 0.0
	for i := k; i <= n; i++ {
		total += float64(choose(n, i)) * pow(p, i) * pow(1-p, n-i)
	}
	return total
}

func choose(n, k int) int {
	if k < 0 || k > n {
		return 0
	}
	c := 1
	for i := 0; i < k; i++ {
		c = c * (n - i) / (i + 1)
	}
	return c
}

func pow(x float64, n int) float64 {
	r := 1.0
	for i := 0; i < n; i++ {
		r *= x
	}
	return r
}

// rateAbsorbed decides whether a non-strict majority-fail is consistent
// with the case's trailing rate. Absorbed when the baseline is real
// (enough history) and the current epoch tally is not statistically
// surprising against it (one-sided binomial, alpha 0.05). A case with a
// near-zero baseline can never absorb — any majority-fail there is news.
func rateAbsorbed(h history, id string, failsNow, epochsNow int) (bool, string) {
	cr, ok := h[id]
	if !ok || cr.Epochs < 9 {
		return false, ""
	}
	p := float64(cr.Fails) / float64(cr.Epochs)
	if p <= 0.02 {
		return false, ""
	}
	tail := binomTail(epochsNow, failsNow, p)
	verdictNote := fmt.Sprintf("%d/%d epochs failed vs trailing rate %.0f%% (%d/%d, p=%.3f)",
		failsNow, epochsNow, p*100, cr.Fails, cr.Epochs, tail)
	if tail >= 0.05 {
		return true, verdictNote
	}
	return false, verdictNote
}

// applyRateTier rewrites non-strict majority-fail verdicts whose epoch
// tally is consistent with the trailing baseline. perEpoch carries the
// pre-aggregation results (adjacency per expandEpochs).
func applyRateTier(agg, perEpoch []result, epochs int, h history) {
	if h == nil || epochs < 2 {
		return
	}
	for i := range agg {
		if agg[i].verdict != vFail || agg[i].c.Strict {
			continue
		}
		fails := 0
		for j := i * epochs; j < (i+1)*epochs && j < len(perEpoch); j++ {
			if perEpoch[j].verdict == vFail {
				fails++
			}
		}
		if absorbed, note := rateAbsorbed(h, agg[i].c.ID, fails, epochs); absorbed {
			agg[i].verdict = vPass
			agg[i].reason = "RATE-ABSORBED (majority-fail within trailing rate — tracked, not a regression): " + note + " — " + agg[i].reason
		} else if note != "" {
			agg[i].reason = "rate-tier: REGRESSION (" + note + ") — " + agg[i].reason
		}
	}
}
