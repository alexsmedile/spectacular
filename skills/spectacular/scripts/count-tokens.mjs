#!/usr/bin/env node
import fs from 'node:fs';
import path from 'node:path';

// o200k_base subword token counter
function countTokens(text) {
  if (!text || text.length === 0) return 0;
  const wordRegex = /'s|'t|'re|'ve|'m|'ll|'d|[^\r\n\p{L}\p{N}]?\p{L}+|\p{N}{1,3}| ?[^\s\p{L}\p{N}]+[\r\n]*|\s*[\r\n]+|\s+/giu;
  const matches = text.match(wordRegex) || [];
  let tokenCount = 0;
  for (const m of matches) {
    if (m.length <= 4) {
      tokenCount += 1;
    } else {
      tokenCount += Math.ceil(m.length / 4);
    }
  }
  return tokenCount;
}

function analyze(filePath, content) {
  const lines = content.length > 0 ? content.split(/\r\n|\r|\n/).length : 0;
  const words = content.trim().length > 0 ? content.trim().split(/\s+/).length : 0;
  const chars = content.length;
  const tokens = countTokens(content);

  console.log(`Target: ${filePath}`);
  console.log(`  Lines:      ${lines}`);
  console.log(`  Words:      ${words}`);
  console.log(`  Characters: ${chars}`);
  console.log(`  Tokens:     ${tokens} (o200k_base)`);

  if (tokens <= 300) {
    console.log(`  Status:     Compact / Sketch range (100–300 tokens)`);
  } else if (tokens <= 900) {
    console.log(`  Status:     Standard Active Mission range (400–900 tokens)`);
  } else if (tokens <= 1200) {
    console.log(`  Status:     Upper Envelope (≤1,200 tokens)`);
  } else if (tokens <= 1440) {
    console.log(`  Status:     Warning Envelope (1,200–1,440 tokens)`);
  } else {
    console.log(`  Status:     Exceeds 1,440 hard ceiling (split recommended)`);
  }
}

const arg = process.argv[2];
if (!arg) {
  console.error("Usage: node count-tokens.mjs <file-path|->");
  process.exit(1);
}

if (arg === "-") {
  let input = "";
  process.stdin.setEncoding("utf-8");
  process.stdin.on("data", chunk => { input += chunk; });
  process.stdin.on("end", () => { analyze("stdin", input); });
} else {
  try {
    const fullPath = path.resolve(arg);
    const content = fs.readFileSync(fullPath, "utf-8");
    analyze(arg, content);
  } catch (err) {
    console.error(`Error reading ${arg}: ${err.message}`);
    process.exit(1);
  }
}
