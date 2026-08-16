package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

// The Skill is the product: `assemble-release` ships SKILL.md and its references
// as the artifact, and for an LLM-executed method the prose is the behavior. So
// these tests are completion evidence for the claims M3 and M4 froze, not a style
// guide.
//
// They therefore pin *concepts*, not sentences. An earlier version asserted whole
// phrases, which locked the wording: rephrasing "The plan supplies meaning" to
// "The plan carries meaning" turned the build red while the shipped behavior was
// unchanged. A concept check still fails when a behavior is dropped, which is the
// only failure worth having.
//
// A named vocabulary token is the exception. `manual-bootstrap` and the FROST lens
// names are strings an agent must emit or match exactly, so rewording them *is* the
// regression and they stay pinned literally.

// concept is one behavior the released Skill must still instruct, identified by the
// keywords that carry it rather than by the sentence that currently phrases it.
type concept struct {
	// behavior states what breaks for a reader of the release if this is missing.
	behavior string
	// keywords must all appear. Keep them to the load-bearing tokens, so an editor
	// can improve a sentence without turning the build red.
	keywords []string
}

// requireConcepts reports every missing concept rather than stopping at the first,
// so one run names the full set of behaviors a Skill rewrite dropped.
func requireConcepts(t *testing.T, claim, content string, concepts []concept) {
	t.Helper()
	for _, want := range concepts {
		var missing []string
		for _, keyword := range want.keywords {
			if !strings.Contains(content, keyword) {
				missing = append(missing, strconv.Quote(keyword))
			}
		}
		if len(missing) == 0 {
			continue
		}
		t.Errorf("released Skill no longer instructs: %s\n"+
			"  missing keywords: %s\n"+
			"  this is completion evidence for claim %q — if the behavior moved or was\n"+
			"  deliberately dropped, update the claim, not just this assertion",
			want.behavior, strings.Join(missing, ", "), claim)
	}
}

// skillText concatenates the named release parts, which is what a reader of the
// shipped artifact sees.
func skillText(t *testing.T, parts ...string) string {
	t.Helper()
	root := filepath.Join("..", "..", "skills", "spectacular")
	var content strings.Builder
	for _, part := range parts {
		data, err := os.ReadFile(filepath.Join(root, part))
		if err != nil {
			t.Fatal(err)
		}
		content.Write(data)
	}
	return content.String()
}

func TestReleasedSkillInstructsTheLaunchPreflight(t *testing.T) {
	content := skillText(t, "SKILL.md", filepath.Join("references", "prepare.md"))
	requireConcepts(t, "lean-launch", content, []concept{
		{
			behavior: "enter the workspace through .spectacular/PROJECT.md before reading anything else",
			keywords: []string{".spectacular/PROJECT.md"},
		},
		{
			behavior: "run the launch check read-only, so a launch never mutates the workspace",
			keywords: []string{"read-only", "preflight"},
		},
	})
}

func TestReleasedSkillInstructsTheOwnerQuestionShape(t *testing.T) {
	content := skillText(t, "SKILL.md", filepath.Join("references", "prepare.md"))
	requireConcepts(t, "dual-layer-questions", content, []concept{
		{
			behavior: "lead an owner question with the plain outcome, then the technical basis",
			keywords: []string{"outcome", "Technical basis"},
		},
		{
			behavior: "state each option as an action paired with its consequence",
			keywords: []string{"action -> consequence"},
		},
		{
			behavior: "name a recommended default so the owner can accept rather than decide",
			keywords: []string{"Recommended default"},
		},
	})
}

func TestReleasedSkillInstructsProgressiveContext(t *testing.T) {
	content := skillText(t, "SKILL.md", filepath.Join("references", "execute.md"))
	requireConcepts(t, "compact-context", content, []concept{
		{
			behavior: "drill down from the Mission card to the current Objective to exact sources",
			keywords: []string{"Mission card -> current Objective -> exact sources"},
		},
		{
			behavior: "separate what the plan is authoritative for from what the tooling is authoritative for",
			keywords: []string{"plan carries meaning", "tooling carries repeatability"},
		},
	})
}

func TestReleasedSkillPinsTheNamedVocabulary(t *testing.T) {
	content := skillText(t,
		"SKILL.md",
		filepath.Join("references", "execute.md"),
		filepath.Join("references", "audit.md"),
	)
	// Unlike the concept checks above, these are exact tokens an agent emits or
	// matches. Rewording one is the regression, so they are pinned literally.
	for _, token := range []string{
		"manual-bootstrap",
		"focused checks",
		"Frozen fit",
		"Truth of proof",
	} {
		if !strings.Contains(content, token) {
			t.Errorf("released Skill omits the exact vocabulary token %q; "+
				"this token is matched or emitted verbatim, so rewording it is a breaking change", token)
		}
	}
}

func TestReleasedSkillInstructsSelfHostingAndActivationAuthority(t *testing.T) {
	content := skillText(t,
		"SKILL.md",
		filepath.Join("references", "execute.md"),
		filepath.Join("references", "audit.md"),
	)
	requireConcepts(t, "governed-self-hosting", content, []concept{
		{
			behavior: "hold the schema and completion criteria frozen while a Mission is active, even when Spectacular develops itself",
			keywords: []string{"active Mission keeps the schema"},
		},
		{
			behavior: "refuse to treat a Decision as authority to activate a Mission",
			keywords: []string{"A Decision is not activation authority"},
		},
	})
}

// The core Skill routes to references instead of absorbing them. These are
// structural assertions, not wording ones: they fail when detail migrates back
// into the always-loaded file and inflates every launch.
func TestCoreSkillRoutesDetailToReferencesInsteadOfAbsorbingIt(t *testing.T) {
	core := skillText(t, "SKILL.md")
	if strings.Contains(core, "Frozen fit**") {
		t.Error("core Skill embeds detailed FROST policy instead of routing to audit.md; " +
			"this reloads audit detail on every launch, which M4 froze as out of the core file")
	}
	if strings.Contains(core, "--event <@Event> --json") {
		t.Error("core Skill still forces full JSON context; " +
			"M3 froze the lean form, where full JSON is opt-in rather than the instructed default")
	}
}

func TestArchiveMetadataIsCanonicalAndReproducible(t *testing.T) {
	entries := []archiveEntry{{name: "spectacular/VERSION", mode: 0o644, data: []byte("2.0.0\n")}}
	first, err := encodeArchive(entries)
	if err != nil {
		t.Fatal(err)
	}
	second, err := encodeArchive(entries)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(first, second) {
		t.Fatal("identical archive inputs produced different bytes")
	}
	gzipReader, err := gzip.NewReader(bytes.NewReader(first))
	if err != nil {
		t.Fatal(err)
	}
	if !gzipReader.ModTime.IsZero() || gzipReader.OS != 255 || gzipReader.Name != "" || gzipReader.Comment != "" {
		t.Fatalf("non-canonical gzip header: %#v", gzipReader.Header)
	}
	tarReader := tar.NewReader(gzipReader)
	header, err := tarReader.Next()
	if err != nil {
		t.Fatal(err)
	}
	if header.Name != entries[0].name || header.Mode != entries[0].mode || header.Uid != 0 || header.Gid != 0 || header.Uname != "root" || header.Gname != "root" || !header.ModTime.Equal(time.Unix(0, 0).UTC()) {
		t.Fatalf("non-canonical tar header: %#v", header)
	}
	data, err := io.ReadAll(tarReader)
	if err != nil || !bytes.Equal(data, entries[0].data) {
		t.Fatalf("archive payload mismatch: %q err=%v", data, err)
	}
}
