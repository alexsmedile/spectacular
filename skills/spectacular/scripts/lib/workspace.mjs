// workspace.mjs — locate a workspace and resolve refs to record paths.
// Read-only. Mirrors where.sh so the two tiers agree on resolution.
import fs from 'fs';
import path from 'path';

export function findWorkspace(start = process.cwd()) {
  let dir = path.resolve(start);
  for (;;) {
    const candidate = path.join(dir, '.spectacular');
    if (fs.existsSync(candidate) && fs.statSync(candidate).isDirectory()) return candidate;
    const parent = path.dirname(dir);
    if (parent === dir) return null;
    dir = parent;
  }
}

function firstMatch(dir, predicate) {
  if (!fs.existsSync(dir)) return null;
  for (const name of fs.readdirSync(dir).sort()) {
    if (predicate(name)) return path.join(dir, name);
  }
  return null;
}

const startsWithRef = (ref) => (name) => name === ref || name.startsWith(ref + '-');

export function resolveRef(ws, ref) {
  if (ref.includes('/')) {
    const [missionRef, childRef] = ref.split('/');
    for (const base of [path.join(ws, 'missions'), path.join(ws, 'archive', 'missions')]) {
      const bundle = firstMatch(base, startsWithRef(missionRef));
      if (!bundle) continue;
      for (const sub of ['runs', 'objectives', 'evidence', 'reviews', 'handoffs', 'decisions', 'gaps', 'assessments']) {
        const subdir = path.join(bundle, sub);
        const hit = firstMatch(subdir, startsWithRef(childRef));
        if (!hit) continue;
        if (fs.statSync(hit).isDirectory()) {
          const inner = path.join(hit, path.basename(hit) + '.md');
          if (fs.existsSync(inner)) return inner;
        } else return hit;
      }
      // Checkpoints nest inside a Run bundle.
      const runs = path.join(bundle, 'runs');
      if (fs.existsSync(runs)) {
        for (const run of fs.readdirSync(runs)) {
          const hit = firstMatch(path.join(runs, run, 'checkpoints'), startsWithRef(childRef));
          if (hit) return hit;
        }
      }
    }
    return null;
  }

  for (const base of [path.join(ws, 'missions'), path.join(ws, 'archive', 'missions')]) {
    const bundle = firstMatch(base, startsWithRef(ref));
    if (bundle && fs.statSync(bundle).isDirectory()) {
      const inner = path.join(bundle, path.basename(bundle) + '.md');
      if (fs.existsSync(inner)) return inner;
    }
  }
  for (const base of ['proposals', 'decisions', 'contracts', path.join('archive', 'proposals')]) {
    const hit = firstMatch(path.join(ws, base), (n) => n.endsWith('.md') && startsWithRef(ref)(n.replace(/\.md$/, '')));
    if (hit) return hit;
  }
  return null;
}

export function listMissions(ws) {
  const out = [];
  for (const [base, archived] of [[path.join(ws, 'missions'), false], [path.join(ws, 'archive', 'missions'), true]]) {
    if (!fs.existsSync(base)) continue;
    for (const name of fs.readdirSync(base).sort()) {
      const record = path.join(base, name, name + '.md');
      if (fs.existsSync(record)) out.push({ ref: name.split('-')[0], dir: name, path: record, archived });
    }
  }
  return out;
}
