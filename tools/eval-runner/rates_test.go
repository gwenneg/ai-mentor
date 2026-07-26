package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeHistoryRun(t *testing.T, dir, name string, cases map[string][]string) {
	t.Helper()
	var b strings.Builder
	for id, verdicts := range cases {
		for e, v := range verdicts {
			fmt.Fprintf(&b, `{"case":"%s","group":"%s","epoch":%d,"attempt":1,"verdict":"%s"}`+"\n", id, id[:1], e+1, v)
		}
	}
	if err := os.WriteFile(filepath.Join(dir, name), []byte(b.String()), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestRateTier(t *testing.T) {
	dir := t.TempDir()
	// C05 fails ~1/3 historically; A18 is clean.
	for i := 1; i <= 4; i++ {
		verdicts := []string{vPass, vPass, vPass}
		if i%2 == 0 {
			verdicts = []string{vFail, vPass, vPass}
		}
		writeHistoryRun(t, dir, fmt.Sprintf("run%d.jsonl", i), map[string][]string{
			"C05": verdicts,
			"A18": {vPass, vPass, vPass},
		})
	}
	h, n, err := loadHistory(dir)
	if err != nil || h == nil {
		t.Fatalf("history should load: %v (runs %d)", err, n)
	}
	if cr := h["C05"]; cr.Epochs != 12 || cr.Fails != 2 {
		t.Fatalf("C05 baseline wrong: %+v", cr)
	}

	// A 2/3 majority-fail on C05 (baseline ~17%) is within rate -> absorbed.
	mk := func(id string, strict bool, verdicts ...string) ([]result, []result) {
		var per []result
		for _, v := range verdicts {
			per = append(per, result{c: evalCase{Group: "C", ID: id, Strict: strict}, verdict: v, reason: "r"})
		}
		agg := aggregateEpochs(append([]result(nil), per...), len(verdicts))
		return agg, per
	}
	agg, per := mk("C05", false, vFail, vFail, vPass)
	applyRateTier(agg, per, 3, h)
	if agg[0].verdict != vPass || !strings.Contains(agg[0].reason, "RATE-ABSORBED") {
		t.Errorf("within-rate majority-fail must absorb, got %s (%s)", agg[0].verdict, agg[0].reason)
	}

	// Same tally on a clean-baseline case -> stays red (near-zero baseline).
	agg, per = mk("A18", false, vFail, vFail, vPass)
	applyRateTier(agg, per, 3, h)
	if agg[0].verdict != vFail {
		t.Errorf("clean-baseline majority-fail must stay red, got %s", agg[0].verdict)
	}

	// 3/3 on C05 exceeds the rate -> REGRESSION, stays red.
	agg, per = mk("C05", false, vFail, vFail, vFail)
	applyRateTier(agg, per, 3, h)
	if agg[0].verdict != vFail || !strings.Contains(agg[0].reason, "REGRESSION") {
		t.Errorf("beyond-rate fail must stay red with the regression label, got %s (%s)", agg[0].verdict, agg[0].reason)
	}

	// Strict cases never enter the tier.
	agg, per = mk("C04", true, vFail, vPass, vPass)
	applyRateTier(agg, per, 3, h)
	if agg[0].verdict != vFail {
		t.Errorf("strict cases must never rate-absorb, got %s", agg[0].verdict)
	}
}

func TestRateTierColdStart(t *testing.T) {
	dir := t.TempDir()
	writeHistoryRun(t, dir, "only.jsonl", map[string][]string{"C05": {vFail, vPass, vPass}})
	h, n, err := loadHistory(dir)
	if err != nil {
		t.Fatal(err)
	}
	if h != nil || n != 1 {
		t.Errorf("below minHistoryRuns the tier must be inactive, got %v (%d)", h, n)
	}
	if h2, _, _ := loadHistory(""); h2 != nil {
		t.Error("no -history dir must mean inactive tier")
	}
}
