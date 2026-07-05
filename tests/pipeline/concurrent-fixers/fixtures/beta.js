// P4 fixture B — independent closed bug (wrong comparison operator).
function isAdult(age) {
  // bug: should be >= 18, not > 18 (18-year-olds are adults)
  return age > 18;
}

if (require.main === module) {
  const assert = require("assert");
  assert.strictEqual(isAdult(18), true, `got ${isAdult(18)}`);
  assert.strictEqual(isAdult(17), false);
  console.log("ok");
}

module.exports = { isAdult };
