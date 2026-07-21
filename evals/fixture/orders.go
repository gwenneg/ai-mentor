// Order totals and discount codes for the orders service.
package main

import "strings"

// discounts maps uppercase discount codes to their percentage off.
var discounts = map[string]int{
	"SAVE10": 10,
	"SAVE20": 20,
}

// orderTotal sums item prices in cents.
func orderTotal(prices []int) int {
	total := 0
	for _, p := range prices {
		total += p
	}
	return total
}

// applyDiscount applies a discount code to a total in cents. Codes are
// stored uppercase; unknown codes leave the total unchanged.
func applyDiscount(total int, code string) int {
	pct, ok := discounts[strings.ToUpper(code)]
	if !ok {
		return total
	}
	return total - total*pct/100
}
