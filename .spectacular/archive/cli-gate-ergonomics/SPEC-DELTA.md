# Spec delta — cli-gate-ergonomics (b29)

### MODIFIED
- specs/index.md :: "CLI: `spectacular policy [@<hook> | --principle N]`" -> "Each policy carries a one-sentence `- directive:` (the practice layer, injected verbatim at gates); hook output is tiered — warn rows print directive + principle title, block rows add the full principle line, `--full` restores paragraphs. CLI: `spectacular policy [@<hook> | --principle N] [--full]`"
- specs/index.md :: "(advances lifecycle; renamed from `promote` in v1.19.0, which remains a deprecated alias)" -> "(advances lifecycle; scaffolds SESSION.md on `planned → active` when absent; renamed from `promote` in v1.19.0, which remains a deprecated alias)"
- specs/index.md :: "`--fix` applies mechanical repairs" -> "Text reports repeat all non-pass findings in a `── findings ──` block after the summary count line. `--fix` applies mechanical repairs"
