#!/usr/bin/env node
// show.mjs — print a record's state and outcome.
//
// Read-only fallback for hosts without the `spectacular` CLI. Parses
// frontmatter properly, unlike the shell tier. It never verifies a fingerprint
// and never writes.
//
// Usage: node show.mjs <ref>
import fs from 'fs';
import { readRecord } from './lib/frontmatter.mjs';
import { findWorkspace, resolveRef } from './lib/workspace.mjs';

const ref = process.argv[2];
if (!ref) { console.error('usage: show.mjs <ref>'); process.exit(2); }

const ws = findWorkspace();
if (!ws) { console.error('no .spectacular/ workspace found'); process.exit(1); }

const file = resolveRef(ws, ref);
if (!file) { console.error(`no record found for ref: ${ref}`); process.exit(1); }

const { frontmatter: f, body } = readRecord(fs.readFileSync(file, 'utf8'));
if (!f) { console.error(`no frontmatter in ${file}`); process.exit(1); }

const line = (k, v) => { if (v !== undefined && v !== null && v !== '') console.log(`${k.padEnd(12)} ${v}`); };

console.log(`${f.ref || ref} — ${f.title || '(untitled)'}\n`);
line('Type:', f.type);
line('Status:', f.status);
line('Owner:', f.owner || f.created_by);
line('Updated:', f.updated);
if (f.resolved_by) line('Resolved by:', f.resolved_by);
if (f.supersedes) line('Supersedes:', f.supersedes);
if (f.activation) {
  line('Activated:', `${f.activation.at || '?'} by ${f.activation.by || '?'}`);
  if (f.activation.fingerprint) line('Fingerprint:', String(f.activation.fingerprint).slice(0, 26) + '… (not verified)');
}
if (f.archive_authorization) line('Archived:', `authorized by ${f.archive_authorization}`);
line('Path:', file.replace(process.cwd() + '/', ''));

if (f.outcome) console.log(`\nOutcome\n  ${f.outcome}`);

if (Array.isArray(f.objectives) && f.objectives.length) {
  console.log('\nObjectives');
  for (const o of f.objectives) {
    // An earlier schema stores objectives as bare `Objective:<id>` references.
    if (typeof o === 'string') { console.log(`  · ${o}`); continue; }
    const mark = o.status === 'implemented' || o.status === 'done' ? '✓' : '·';
    console.log(`  ${mark} ${o.ref || '?'} — ${o.title || o.outcome || ''}`);
    if (o.after) console.log(`      after: ${[].concat(o.after).join(', ')}`);
    if (o.after_interface) console.log(`      after_interface: ${[].concat(o.after_interface).join(', ')}`);
  }
}

if (Array.isArray(f.gaps) && f.gaps.length) {
  console.log('\nGaps');
  for (const g of f.gaps) {
    const state = g.resolution ? 'resolved' : (g.blocked_on ? 'OPEN' : 'unknown');
    console.log(`  ${state === 'OPEN' ? '!' : '·'} ${g.ref} — ${state}`);
    if (g.blocked_on) console.log(`      blocked_on: ${g.blocked_on}`);
  }
}

const heading = body.split('\n').find((l) => l.startsWith('# '));
if (heading) console.log(`\n${heading.replace(/^# /, '')} …`);

console.log('\nread from the file only; no fingerprint or binding was verified');
