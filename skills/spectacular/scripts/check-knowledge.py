#!/usr/bin/env python3
"""Optional durable-context diagnostics. Requires Python 3 and PyYAML.

Exit 0: no findings; 1: advisory findings; 2: inspection unavailable.
Never writes files or validates governed identities, bindings, or Contracts.
"""
import argparse
import datetime
import json
import re
import sys
from pathlib import Path
from urllib.parse import unquote, urlsplit

try:
    import yaml
except ImportError:
    print("inspection unavailable: PyYAML is required; inspect the folder agreement manually", file=sys.stderr)
    raise SystemExit(2)


def inspect(root, rules):
    findings, documents = [], []
    exempt = set(rules["unconstrained"]) | set(rules["governed"]["collections"]) | set(rules["historical"])
    required = rules["metadata"]["required"]
    candidates = []
    for path in sorted(root.rglob("*.md")):
        rel = path.relative_to(root)
        if any(part.startswith(".") or part in exempt for part in rel.parts[:-1]):
            continue
        if path.name in rules["navigation"]["legacy"]:
            continue
        if not path.resolve().is_relative_to(root.resolve()):
            findings.append({"path": rel.as_posix(), "issue": "file escapes workspace"})
            continue
        candidates.append(path)
    all_targets = [p for p in root.rglob("*.md") if p.resolve().is_relative_to(root.resolve())]

    def note(path, issue):
        findings.append({"path": path.relative_to(root).as_posix(), "issue": issue})

    for path in candidates:
        try:
            text = path.read_text(encoding="utf-8")
            match = re.match(r"\A---\r?\n(.*?)\r?\n---(?:\r?\n|$)", text, re.S)
            if not match:
                note(path, "missing frontmatter")
                continue
            fields = yaml.safe_load(match[1])
            if not isinstance(fields, dict):
                note(path, "frontmatter must be a mapping")
                continue
            # Existing record claims are validated by their governed interface.
            if any(key in fields for key in ("id", "ref", "human_ref", "schema", "schema_version")):
                continue
            for key in required:
                if key not in fields or fields[key] is None or str(fields[key]).strip() == "":
                    note(path, f"missing {key}")
            if "type" in fields and not isinstance(fields["type"], str):
                note(path, "type must be a non-empty string")
            for key in ("created", "updated"):
                if key in fields:
                    try:
                        stamp = datetime.datetime.fromisoformat(str(fields[key]).replace("Z", "+00:00"))
                        if stamp.tzinfo is None:
                            raise ValueError()
                    except (ValueError, TypeError):
                        note(path, f"{key} must be an ISO timestamp with timezone; do not invent historical dates")
            collection = path.relative_to(root).parts[0]
            agreement = rules["collections"].get(collection)
            if agreement and fields.get("type") not in agreement["types"]:
                note(path, f"type outside {collection} agreement; clarify or extend the agreement")
            if path.name == rules["navigation"]["manual"] and fields.get("type") != "Index":
                note(path, "manual INDEX.md should declare type Index")
            documents.append({"path": path.relative_to(root).as_posix(), "metadata": fields})
            body = text[match.end():]
            # Ignore fenced code: examples are not navigation edges.
            body = re.sub(r"(?ms)^\s*(`{3,}|~{3,}).*?^\s*\1\s*$", "", body)
            targets = [(m.group(1).split("|", 1)[0], True) for m in re.finditer(r"(?<!!)\[\[([^\]]+)\]\]", body)]
            targets += [(m.group(1).strip("<>"), False) for m in re.finditer(r"(?<!!)\[[^\]\n]+\]\(([^\s)]+)\)", body)]
            for target, wiki in targets:
                if urlsplit(target).scheme or target.startswith("//"):
                    continue
                name = unquote(target.split("#", 1)[0])
                if not name:
                    continue
                destination = (root / name.lstrip("/")) if wiki or name.startswith("/") else path.parent / name
                if wiki and not destination.suffix:
                    destination = destination.with_suffix(".md")
                if destination.is_file():
                    continue
                alternatives = [p for p in all_targets if p.stem == Path(name).stem] if wiki and "/" not in name else []
                if len(alternatives) != 1:
                    note(path, f"{'ambiguous' if alternatives else 'missing'} link: {target}")
        except (OSError, UnicodeError, yaml.YAMLError) as exc:
            note(path, f"unreadable metadata: {exc}")
    return {"basis": "durable-context diagnostics; governed validation excluded", "documents": documents, "findings": findings}


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("workspace", type=Path, help="the .spectacular directory")
    args = parser.parse_args()
    if not args.workspace.is_dir():
        parser.error("workspace directory does not exist")
    rules_path = Path(__file__).resolve().parent.parent / "knowledge-folders.yaml"
    try:
        rules = yaml.safe_load(rules_path.read_text())
        report = inspect(args.workspace.resolve(), rules)
    except (OSError, yaml.YAMLError, KeyError) as exc:
        print(f"inspection unavailable: {exc}")
        return 2
    print(json.dumps(report, indent=2, default=str))
    return 1 if report["findings"] else 0


if __name__ == "__main__":
    raise SystemExit(main())
