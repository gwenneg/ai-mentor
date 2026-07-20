package main

import "testing"

func TestOrderTotal(t *testing.T) {
	if got := orderTotal([]int{250, 100}); got != 350 {
		t.Fatalf("orderTotal = %d, want 350", got)
	}
}

func TestApplyDiscount(t *testing.T) {
	if got := applyDiscount(1000, "SAVE10"); got != 900 {
		t.Fatalf("applyDiscount = %d, want 900", got)
	}
	if got := applyDiscount(1000, "nope"); got != 1000 {
		t.Fatalf("unknown code must not discount, got %d", got)
	}
}
