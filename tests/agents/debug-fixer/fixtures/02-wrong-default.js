// Pristine fixture 02 — wrong default parameter value.
//
// Bug: paginate() should default to 20 per page, but defaults to 0,
// which makes the first page empty. Closed, single-site, mechanical.

function paginate(items, page = 1, perPage = 0) {
  const start = (page - 1) * perPage;
  return items.slice(start, start + perPage);
}

function _demo() {
  const items = Array.from({ length: 50 }, (_, i) => i);
  const first = paginate(items, 1);
  if (first.length !== 20) throw new Error(`expected 20, got ${first.length}`);
  if (first[0] !== 0 || first[19] !== 19) throw new Error("wrong slice");
  console.log("ok");
}

if (require.main === module) _demo();

module.exports = { paginate };
