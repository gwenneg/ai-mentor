"use strict";

// Order helpers for the fixture service. Deliberately small and boring:
// eval cases need real paths and commands to ground against, not real logic.

function orderTotal(items) {
  return items.reduce((sum, item) => sum + item.price * item.quantity, 0);
}

function applyDiscount(total, code) {
  if (code === "SAVE10") return total * 0.9;
  if (code === "SAVE20") return total * 0.8;
  return total;
}

module.exports = { orderTotal, applyDiscount };
