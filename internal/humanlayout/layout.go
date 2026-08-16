// Package humanlayout maps durable records to a filesystem that humans can
// navigate without surrendering UUID identity or revision fingerprints.
package humanlayout

import (
	"crypto/sha256"
	"encoding/base32"
	"fmt"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/alexsmedile/spectacular/v2/internal/discovery"
	"github.com/alexsmedile/spectacular/v2/internal/domain"
	"github.com/alexsmedile/spectacular/v2/internal/workspace"
)

// ShortKey is a stable, non-semantic key derived from durable identity. It is
// deliberately not the content fingerprint.
func ShortKey(id domain.ID) string {
	sum := sha256.Sum256([]byte(id.String()))
	return strings.ToLower(base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(sum[:]))[:6]
}

func Slug(text string) string {
	text = strings.ToLower(strings.TrimSpace(text))
	var out []rune
	dash := false
	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			out = append(out, r)
			dash = false
		} else if len(out) > 0 && !dash {
			out = append(out, '-')
			dash = true
		}
	}
	return strings.Trim(string(out), "-")
}

func HumanRef(doc *workspace.Document) string {
	return workspace.RefOrEmpty(doc)
}

// Plan annotates new documents with human_ref and returns canonical paths.
// Ordinals are presentation identities scoped by their owning bundle; UUIDs
// remain the collision-proof authority.
func Plan(existing []discovery.Entry, docs []*workspace.Document) (map[domain.ID]string, error) {
	all := map[domain.ID]*workspace.Document{}
	priorByID := map[domain.ID]*workspace.Document{}
	stableDirectories := map[string]string{}
	used := map[string]bool{}
	for _, entry := range existing {
		all[entry.Document.Record.ID] = entry.Document
		priorByID[entry.Document.Record.ID] = entry.Document
		if ref := HumanRef(entry.Document); ref != "" {
			used[ref] = true
			if rootFile := bundleRootFile(entry.Document.Record.Type); rootFile != "" && filepath.Base(entry.Path) == rootFile {
				stableDirectories[ref] = filepath.ToSlash(filepath.Dir(strings.TrimPrefix(entry.Path, ".spectacular/")))
			}
		}
	}
	for _, doc := range docs {
		all[doc.Record.ID] = doc
	}

	ordered := append([]*workspace.Document(nil), docs...)
	sort.SliceStable(ordered, func(i, j int) bool { return rank(ordered[i].Record.Type) < rank(ordered[j].Record.Type) })
	for _, doc := range ordered {
		if prior := priorByID[doc.Record.ID]; prior != nil {
			if ref := HumanRef(prior); ref != "" {
				workspace.SetString(doc, "human_ref", ref)
				used[ref] = true
				continue
			}
		}
		if HumanRef(doc) != "" {
			used[HumanRef(doc)] = true
			continue
		}
		ref, err := nextRef(doc, all, used)
		if err != nil {
			return nil, err
		}
		workspace.SetString(doc, "human_ref", ref)
		used[ref] = true
	}

	paths := map[domain.ID]string{}
	for _, doc := range docs {
		path, err := pathFor(doc, all, stableDirectories)
		if err != nil {
			return nil, err
		}
		paths[doc.Record.ID] = path
	}
	return paths, nil
}

func bundleRootFile(noun domain.RecordType) string {
	switch noun {
	case domain.Mission:
		return "MISSION.md"
	case domain.Run:
		return "RUN.md"
	default:
		return ""
	}
}

func rank(t domain.RecordType) int {
	switch t {
	case domain.Mission:
		return 0
	case domain.Objective:
		return 1
	case domain.Run:
		return 2
	case domain.Checkpoint:
		return 3
	default:
		return 4
	}
}

func nextRef(doc *workspace.Document, all map[domain.ID]*workspace.Document, used map[string]bool) (string, error) {
	if doc.Record.Type == domain.Contract {
		return "CC-" + ShortKey(doc.Record.ID), nil
	}
	prefix := map[domain.RecordType]string{
		domain.Mission: "M", domain.Proposal: "P",
		domain.Objective: "O", domain.Run: "R", domain.Checkpoint: "C",
		domain.Evidence: "E", domain.Decision: "D", domain.Gap: "G",
		domain.Handoff: "H", domain.Assessment: "A",
	}[doc.Record.Type]
	if prefix == "" {
		return "", fmt.Errorf("no human reference grammar for %s", doc.Record.Type)
	}
	parent, err := parentRef(doc, all)
	if err != nil {
		return "", err
	}
	max := 0
	needle := prefix
	if parent != "" {
		needle = parent + "/" + prefix
	}
	for ref := range used {
		if !strings.HasPrefix(ref, needle) {
			continue
		}
		rest := strings.TrimPrefix(ref, needle)
		digits := strings.TrimLeftFunc(rest, unicode.IsDigit)
		n := len(rest) - len(digits)
		if n == 0 {
			continue
		}
		value, _ := strconv.Atoi(rest[:n])
		if value > max {
			max = value
		}
	}
	base := fmt.Sprintf("%s%d", prefix, max+1)
	if parent != "" {
		base = parent + "/" + base
	}
	if doc.Record.Type == domain.Evidence || doc.Record.Type == domain.Decision || doc.Record.Type == domain.Gap || doc.Record.Type == domain.Handoff || doc.Record.Type == domain.Assessment {
		base += "-" + ShortKey(doc.Record.ID)
	}
	return base, nil
}

func parentRef(doc *workspace.Document, all map[domain.ID]*workspace.Document) (string, error) {
	field := ""
	switch doc.Record.Type {
	case domain.Objective, domain.Run, domain.Evidence, domain.Handoff, domain.Assessment:
		field = "mission"
	case domain.Checkpoint:
		field = "run"
	case domain.Gap:
		field = "scope"
	case domain.Decision:
		if direct, _ := workspace.String(doc, "mission", false); direct != "" {
			field = "mission"
		} else if targets, _ := workspace.Strings(doc, "targets", false); len(targets) > 0 {
			for _, target := range targets {
				if typed, err := domain.ParseReference(target); err == nil && typed.Type == domain.Mission {
					if owner := all[typed.ID]; owner != nil {
						return HumanRef(owner), nil
					}
				}
			}
		}
	}
	if field == "" {
		return "", nil
	}
	value, _ := workspace.String(doc, field, false)
	if value == "" {
		return "", nil
	}
	typed, err := domain.ParseReference(value)
	if err != nil {
		return "", err
	}
	owner := all[typed.ID]
	if owner == nil {
		return "", fmt.Errorf("human layout parent %s is unavailable", value)
	}
	if doc.Record.Type == domain.Checkpoint {
		return HumanRef(owner), nil
	}
	if owner.Record.Type != domain.Mission {
		return "", fmt.Errorf("human layout parent %s is not a Mission", value)
	}
	return HumanRef(owner), nil
}

func Path(doc *workspace.Document, all map[domain.ID]*workspace.Document) (string, error) {
	return pathFor(doc, all, nil)
}

func pathFor(doc *workspace.Document, all map[domain.ID]*workspace.Document, stableDirectories map[string]string) (string, error) {
	ref := HumanRef(doc)
	if ref == "" {
		return "", fmt.Errorf("record %s lacks human_ref", doc.Record.ID)
	}
	title := "record"
	if doc.Record.Title != nil && Slug(*doc.Record.Title) != "" {
		title = Slug(*doc.Record.Title)
	}
	leaf := ref
	if slash := strings.LastIndex(leaf, "/"); slash >= 0 {
		leaf = leaf[slash+1:]
	}
	filename := leaf + ".md"
	switch doc.Record.Type {
	case domain.Mission:
		return filepath.ToSlash(filepath.Join(".spectacular", "missions", ref+"-"+title, "MISSION.md")), nil
	case domain.Objective:
		mission := strings.Split(ref, "/")[0]
		return filepath.ToSlash(filepath.Join(".spectacular", missionDirectory(mission, all, stableDirectories), "objectives", leaf+"-"+title+".md")), nil
	case domain.Run:
		mission := strings.Split(ref, "/")[0]
		return filepath.ToSlash(filepath.Join(".spectacular", missionDirectory(mission, all, stableDirectories), "runs", leaf+"-"+title, "RUN.md")), nil
	case domain.Checkpoint:
		parts := strings.Split(ref, "/")
		mission, run := parts[0], parts[1]
		runDirectory := stableDirectories[mission+"/"+run]
		if runDirectory == "" {
			runDirectory = filepath.Join(missionDirectory(mission, all, stableDirectories), "runs", run+"-"+parentSlug(run, all))
		}
		return filepath.ToSlash(filepath.Join(".spectacular", runDirectory, "checkpoints", leaf+"-"+title+".md")), nil
	case domain.Contract:
		return filepath.ToSlash(filepath.Join(".spectacular", "contracts", ref+"-"+title+".md")), nil
	case domain.Proposal:
		return filepath.ToSlash(filepath.Join(".spectacular", "proposals", ref+"-"+title+".md")), nil
	case domain.Evidence, domain.Decision, domain.Gap, domain.Handoff, domain.Assessment:
		folder := strings.ToLower(string(doc.Record.Type)) + "s"
		if doc.Record.Type == domain.Evidence {
			folder = "evidence"
		}
		parts := strings.Split(ref, "/")
		if len(parts) > 1 && strings.HasPrefix(parts[0], "M") {
			return filepath.ToSlash(filepath.Join(".spectacular", missionDirectory(parts[0], all, stableDirectories), folder, filename)), nil
		}
		return filepath.ToSlash(filepath.Join(".spectacular", folder, filename)), nil
	}
	return "", fmt.Errorf("no canonical human path for %s", doc.Record.Type)
}

func missionDirectory(ref string, all map[domain.ID]*workspace.Document, stableDirectories map[string]string) string {
	if existing := stableDirectories[ref]; existing != "" {
		return existing
	}
	for _, doc := range all {
		if doc.Record.Type == domain.Mission && HumanRef(doc) == ref {
			title := "mission"
			if doc.Record.Title != nil && Slug(*doc.Record.Title) != "" {
				title = Slug(*doc.Record.Title)
			}
			return filepath.Join("missions", ref+"-"+title)
		}
	}
	return filepath.Join("missions", ref+"-mission")
}

func parentSlug(runRef string, all map[domain.ID]*workspace.Document) string {
	for _, doc := range all {
		if doc.Record.Type == domain.Run && strings.HasSuffix(HumanRef(doc), "/"+runRef) {
			if doc.Record.Title != nil && Slug(*doc.Record.Title) != "" {
				return Slug(*doc.Record.Title)
			}
		}
	}
	return "run"
}
