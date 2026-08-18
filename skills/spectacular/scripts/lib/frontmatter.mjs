// frontmatter.mjs — a bounded YAML reader for Spectacular record frontmatter.
//
// Not a general YAML implementation. It handles exactly the shapes Spectacular
// records use: scalars, nested maps, sequences of scalars, sequences of maps,
// and block scalars (`>-`, `|`). Anything it does not understand is left out
// rather than guessed at, because a parser that silently mis-reads a governance
// record is worse than one that admits its limits.
//
// Read-only. Nothing here writes a record or computes a fingerprint.

const INDENT = (line) => line.length - line.trimStart().length;

function parseScalar(raw) {
  const v = raw.trim();
  if (v === '') return '';
  if (v === 'true') return true;
  if (v === 'false') return false;
  if (v === 'null' || v === '~') return null;
  if (/^-?\d+$/.test(v)) return Number(v);
  if ((v.startsWith('"') && v.endsWith('"')) || (v.startsWith("'") && v.endsWith("'"))) {
    return v.slice(1, -1);
  }
  return v;
}

function readBlock(lines, start, parentIndent) {
  const out = [];
  let i = start;
  while (i < lines.length) {
    const line = lines[i];
    if (line.trim() === '') { out.push(''); i++; continue; }
    if (INDENT(line) <= parentIndent) break;
    out.push(line.trim());
    i++;
  }
  // Folded (`>-`) is the shape Spectacular uses; join on spaces.
  return [out.join(' ').replace(/\s+/g, ' ').trim(), i];
}

function parseBlock(lines, start, indent) {
  const map = {};
  let i = start;

  while (i < lines.length) {
    const line = lines[i];
    if (line.trim() === '' || line.trim().startsWith('#')) { i++; continue; }
    const ind = INDENT(line);
    if (ind < indent) break;

    const trimmed = line.trim();
    if (trimmed.startsWith('- ') || trimmed === '-') break;

    const m = trimmed.match(/^([^:]+?):(?:\s+(.*))?$/);
    if (!m) { i++; continue; }

    const key = m[1].trim();
    const rest = m[2] === undefined ? '' : m[2];

    if (rest === '>-' || rest === '>' || rest === '|' || rest === '|-') {
      const [text, next] = readBlock(lines, i + 1, ind);
      map[key] = text;
      i = next;
      continue;
    }

    if (rest !== '') { map[key] = parseScalar(rest); i++; continue; }

    let j = i + 1;
    while (j < lines.length && (lines[j].trim() === '' || lines[j].trim().startsWith('#'))) j++;
    if (j >= lines.length) { map[key] = null; i = j; continue; }

    const childIndent = INDENT(lines[j]);
    if (childIndent <= ind) { map[key] = null; i = j; continue; }

    if (lines[j].trim().startsWith('-')) {
      const [arr, next] = parseSequence(lines, j, childIndent);
      map[key] = arr;
      i = next;
    } else {
      const [obj, next] = parseBlock(lines, j, childIndent);
      map[key] = obj;
      i = next;
    }
  }
  return [map, i];
}

function parseSequence(lines, start, indent) {
  const arr = [];
  let i = start;
  while (i < lines.length) {
    const line = lines[i];
    if (line.trim() === '' || line.trim().startsWith('#')) { i++; continue; }
    const ind = INDENT(line);
    if (ind < indent) break;
    const trimmed = line.trim();
    if (!trimmed.startsWith('-')) break;

    const inline = trimmed.slice(1).trim();
    if (inline === '') { i++; continue; }

    if (/^[^:\s]+(\s[^:]*)?:(\s|$)/.test(inline)) {
      // A map item. Its keys sit at the column where the inline key begins,
      // which is the dash column plus the "- " prefix.
      const keyIndent = ind + (trimmed.length - inline.length);
      const itemLines = [' '.repeat(keyIndent) + inline];
      let j = i + 1;
      while (j < lines.length) {
        if (lines[j].trim() === '') { itemLines.push(lines[j]); j++; continue; }
        const jInd = INDENT(lines[j]);
        if (jInd <= ind) break;
        itemLines.push(lines[j]);
        j++;
      }
      const [obj] = parseBlock(itemLines, 0, keyIndent);
      arr.push(obj);
      i = j;
    } else {
      arr.push(parseScalar(inline));
      i++;
    }
  }
  return [arr, i];
}

/** Split a record into { frontmatter, body }. frontmatter is null when absent. */
export function readRecord(text) {
  if (!text.startsWith('---\n')) return { frontmatter: null, body: text };
  const end = text.indexOf('\n---\n', 4);
  if (end === -1) return { frontmatter: null, body: text };
  const raw = text.slice(4, end);
  const body = text.slice(end + 5);
  const [fm] = parseBlock(raw.split('\n'), 0, 0);
  return { frontmatter: fm, body };
}
