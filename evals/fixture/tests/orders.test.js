"use strict";

const { test } = require("node:test");
const assert = require("node:assert");
const { orderTotal, applyDiscount } = require("../src/orders.js");

test("orderTotal sums price times quantity", () => {
  assert.strictEqual(orderTotal([{ price: 5, quantity: 2 }]), 10);
});

test("applyDiscount handles SAVE10", () => {
  assert.strictEqual(applyDiscount(100, "SAVE10"), 90);
});
