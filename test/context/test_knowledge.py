"""Behavioral checks for optional context diagnostics, outside governed validation."""
import importlib.util
import subprocess
import sys
import tempfile
import unittest
from pathlib import Path

SCRIPT = Path(__file__).resolve().parents[2] / 'skills/spectacular/scripts/check-knowledge.py'
try:
    import yaml
except ImportError:
    yaml = None


@unittest.skipIf(yaml is None, 'optional PyYAML unavailable')
class KnowledgeDiagnostics(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        spec = importlib.util.spec_from_file_location('knowledge', SCRIPT)
        cls.module = importlib.util.module_from_spec(spec)
        spec.loader.exec_module(cls.module)
        cls.rules = yaml.safe_load((SCRIPT.parent.parent / 'knowledge-folders.yaml').read_text())

    def setUp(self):
        self.temp = tempfile.TemporaryDirectory()
        self.addCleanup(self.temp.cleanup)
        self.root = Path(self.temp.name)

    def document(self, name, body='', **overrides):
        fields = dict(type='Plan', version='1.0', created='2026-10-07T10:00:00Z', updated='2026-10-07T10:00:00Z')
        fields.update(overrides)
        path = self.root / name
        path.parent.mkdir(parents=True, exist_ok=True)
        path.write_text('---\n' + yaml.safe_dump(fields) + '---\n' + body)
        return path

    def report(self):
        return self.module.inspect(self.root, self.rules)

    def test_raw_aliases_and_governed_claims_are_exempt(self):
        for folder in ('raw', 'sketch', 'scratchpad', 'contracts', 'history'):
            path = self.root / folder / 'capture.md'
            path.parent.mkdir()
            path.write_bytes(b'\xffdraft with no metadata')
        self.document('decisions/choice.md', type='Decision', schema='existing-governed-record', created=None)
        self.assertEqual(self.report()['findings'], [])
        self.assertEqual(self.report()['documents'], [])

    def test_metadata_and_folder_agreement_findings(self):
        self.document('plans/work.md', type='Atlas', version=None, created='unknown')
        issues = [x['issue'] for x in self.report()['findings']]
        self.assertIn('missing version', issues)
        self.assertTrue(any('created must be' in x for x in issues))
        self.assertTrue(any('outside plans agreement' in x for x in issues))

    def test_relative_and_path_qualified_wikilinks(self):
        self.document('atlas/domain.md', type='Entity')
        self.document('INDEX.md', type='Index')
        self.document('plans/work.md', '[Domain](../atlas/domain.md)\n[[atlas/domain|Domain]]\n[[INDEX]]\n```md\n[[missing-example]]\n```\n')
        self.assertEqual(self.report()['findings'], [])

    def test_missing_and_ambiguous_targets(self):
        self.document('atlas/topic.md', type='Atlas')
        self.document('plans/topic.md')
        self.document('plans/work.md', '[[topic]]\n[Absent](absent.md)\n')
        issues = [x['issue'] for x in self.report()['findings']]
        self.assertIn('ambiguous link: topic', issues)
        self.assertIn('missing link: absent.md', issues)

    def test_manual_index_type_and_governed_frontmatter_exclusion(self):
        self.document('INDEX.md')
        self.assertTrue(any('type Index' in x['issue'] for x in self.report()['findings']))


class MissingDependency(unittest.TestCase):
    def test_unavailable_inspection_does_not_claim_success(self):
        result = subprocess.run([sys.executable, '-I', '-S', str(SCRIPT), '.'], capture_output=True, text=True)
        self.assertEqual(result.returncode, 2)
        self.assertIn('inspection unavailable', result.stderr)


if __name__ == '__main__':
    unittest.main()
