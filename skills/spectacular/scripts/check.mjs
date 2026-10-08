#!/usr/bin/env node
// check.mjs — structural validation of a record, or the whole workspace.
//
// Read-only fallback for hosts without the `spectacular` CLI. It checks shape:
// extracted fields exist and referenced records resolve. YAML parsing is
// unverified: the bundled bounded reader is not a conformant YAML parser. It
// does NOT verify fingerprints, bindings, claim drift, or authority — those
// need the CLI, and this says so rather than implying coverage it lacks.
//
// Usage: node check.mjs [<ref>]      no ref checks every record
import fs from 'fs';
import path from 'path';
import { readRecord } from './lib/frontmatter.mjs';
import { findWorkspace, resolveRef } from './lib/workspace.mjs';

const REQUIRED = { Mission: ['type','id','title','status'], Proposal: ['type','id','ref','title','status'],
  Decision: ['type','id','ref','title'], Contract: ['type','id','ref','title'], default: ['type','id'] };

const ws = findWorkspace();
if (!ws) { console.error('no .spectacular/ workspace found'); process.exit(1); }

function walk(dir, out = []) {
  for (const e of fs.readdirSync(dir, { withFileTypes: true })) {
    const p = path.join(dir, e.name);
    if (e.isDirectory() && !['raw', 'sketch', 'scratchpad', 'atlas', 'plans', 'campaigns', 'retrospectives', '.engine', '.cache', 'history'].includes(e.name)) walk(p, out);
    else if (e.isFile() && e.name.endsWith('.md') && e.name.toLowerCase() !== 'index.md') out.push(p);
  }
  return out;
}

function checkOne(file) {
  const problems = [];
  const notices = [];
  const text = fs.readFileSync(file, 'utf8');
  const { frontmatter: f } = readRecord(text);
  if (!f) return { file, skipped: 'no frontmatter', problems, notices };

  const claimed = ['id', 'ref', 'human_ref', 'schema', 'schema_version'].some(key => key in f);
  const parts = path.relative(ws, file).split(path.sep);
  const archivedScenario = parts[0] === 'archive' && parts[1] === 'scenarios' && ['Scenario', 'Reference'].includes(f.type);
  const strict = parts.some(part => ['contracts', 'proposals', 'missions', 'evidence', 'gaps', 'handoffs', 'assessments', 'reviews'].includes(part)) || (parts.includes('archive') && !archivedScenario);
  const canonical = /^(CC-|[A-Z]+[0-9]+([-\.]|$))/.test(path.basename(file));
  if (!claimed && !strict && path.basename(file) !== 'PROJECT.md' && (!canonical || f.governance === 'context')) {
    return { file, skipped: 'soft context; use check-knowledge.py for advisory metadata checks', problems, notices };
  }

  const required = REQUIRED[f.type] || REQUIRED.default;
  for (const key of required) {
    if (f[key] !== undefined) continue;
    // `human_ref:` is the legacy spelling of `ref:`. The CLI reports this as a
    // notice and keeps the record valid; match that rather than failing it.
    if (key === 'ref' && f.human_ref !== undefined) {
      notices.push('uses legacy `human_ref:`; new records declare `ref:`');
      continue;
    }
    problems.push(`missing \`${key}:\``);
  }

  if (f.status !== undefined && typeof f.status !== 'string') problems.push('`status:` is not a scalar');

  // Referenced records must resolve.
  if (f.resolved_by && !resolveRef(ws, String(f.resolved_by))) {
    problems.push(`\`resolved_by: ${f.resolved_by}\` does not resolve`);
  }
  for (const g of (Array.isArray(f.gaps) ? f.gaps : [])) {
    if (!g.ref) problems.push('a gap entry has no `ref:`');
    if (!g.resolution && !g.blocked_on) problems.push(`gap \`${g.ref}\` has neither resolution nor blocked_on`);
  }
  for (const o of (Array.isArray(f.objectives) ? f.objectives : [])) {
    // Two legal shapes: an inline map with `ref:`, or a bare `Objective:<id>`
    // reference from an earlier schema. Neither is a defect.
    if (typeof o === 'string') {
      if (!/^Objective:/.test(o)) problems.push(`objective entry is not a reference: ${o}`);
      continue;
    }
    if (!o.ref) problems.push('an objective has no `ref:`');
  }
  return { file, problems, notices };
}

const ref = process.argv[2];
const targets = ref ? [resolveRef(ws, ref)] : walk(ws);
if (ref && !targets[0]) { console.error(`no record found for ref: ${ref}`); process.exit(1); }

let checked = 0, skipped = 0, failing = 0, noticed = 0;
for (const file of targets) {
  const r = checkOne(file);
  if (r.skipped) { skipped++; continue; }
  checked++;
  if (r.notices.length) noticed++;
  if (r.problems.length || (ref && r.notices.length)) {
    if (r.problems.length) failing++;
    console.log(`\n${file.replace(process.cwd() + '/', '')}`);
    for (const p of r.problems) console.log(`  - ${p}`);
    for (const n of r.notices) console.log(`  notice: ${n}`);
  }
}

console.log(`\n${checked} record(s) inspected, ${failing} with problems, ${noticed} with notices, ${skipped} outside governed inspection`);
console.log('best-effort inspection only — YAML parsing, fingerprints, bindings, claim drift, and authority were NOT verified');
console.log('run `spectacular mission check <ref>` for the real gate');
// Exit 2 distinguishes incomplete inspection from a successful validation gate.
process.exit(failing ? 1 : 2);
