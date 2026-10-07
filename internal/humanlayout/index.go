package humanlayout

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/alexsmedile/spectacular/v2/internal/discovery"
	"github.com/alexsmedile/spectacular/v2/internal/domain"
	"github.com/alexsmedile/spectacular/v2/internal/workspace"
)

type indexRow struct {
	Ref    string `json:"ref"`
	Noun   string `json:"noun"`
	Title  string `json:"title"`
	Path   string `json:"path"`
	Status string `json:"status,omitempty"`
}

// Indexes returns deterministic non-authoritative JSON projections for
// the root and affected collections. The root Mission record is the Mission-bundle index;
// generating another index inside every Mission would duplicate it.
func Indexes(existing []discovery.Entry, pending []*workspace.Document, paths map[domain.ID]string) (map[string][]byte, error) {
	pendingIDs := map[domain.ID]bool{}
	touchedCollections := map[string]bool{}
	rows := make([]indexRow, 0, len(existing)+len(pending))
	for _, doc := range pending {
		pendingIDs[doc.Record.ID] = true
	}
	for _, entry := range existing {
		if pendingIDs[entry.Document.Record.ID] {
			if index := collectionIndex(entry.Path); index != "" {
				touchedCollections[index] = true
			}
			continue
		}
		rows = append(rows, row(entry.Document, entry.Path))
	}
	for _, doc := range pending {
		path := paths[doc.Record.ID]
		if path == "" {
			return nil, fmt.Errorf("missing planned path for %s", doc.Record.ID)
		}
		rows = append(rows, row(doc, path))
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].Ref == rows[j].Ref {
			return rows[i].Path < rows[j].Path
		}
		return rows[i].Ref < rows[j].Ref
	})

	groups := map[string][]indexRow{}
	for path := range touchedCollections {
		groups[path] = []indexRow{}
	}
	for _, item := range rows {
		parts := strings.Split(filepath.ToSlash(item.Path), "/")
		if len(parts) < 3 || parts[0] != ".spectacular" {
			continue
		}
		top := parts[1]
		if top == "PROJECT.md" || top == "PRODUCT.md" || top == "ARCHITECTURE.md" || top == "STACK.md" || top == "ROADMAP.md" || top == "VOCABULARY.md" || top == "ONTOLOGY.md" {
			continue
		}
		if top == "missions" && item.Noun != string(domain.Mission) {
			continue
		}
		if top == "archive" && item.Noun != string(domain.Mission) && item.Noun != string(domain.Proposal) {
			continue
		}
		indexPath := filepath.ToSlash(filepath.Join(".spectacular", top, "index.json"))
		groups[indexPath] = append(groups[indexPath], item)
	}

	out := map[string][]byte{}
	rootIndex, err := generatedIndex(rows)
	if err != nil {
		return nil, err
	}
	out[".spectacular/index.json"] = rootIndex
	for path, items := range groups {
		sort.Slice(items, func(i, j int) bool { return items[i].Ref < items[j].Ref })
		data, err := generatedIndex(items)
		if err != nil {
			return nil, err
		}
		out[path] = data
	}
	return out, nil
}

func collectionIndex(path string) string {
	parts := strings.Split(filepath.ToSlash(path), "/")
	if len(parts) < 3 || parts[0] != ".spectacular" {
		return ""
	}
	switch parts[1] {
	case "PROJECT.md", "PRODUCT.md", "ARCHITECTURE.md", "STACK.md", "ROADMAP.md", "VOCABULARY.md", "ONTOLOGY.md":
		return ""
	default:
		return filepath.ToSlash(filepath.Join(".spectacular", parts[1], "index.json"))
	}
}

func row(doc *workspace.Document, path string) indexRow {
	ref := HumanRef(doc)
	if ref == "" {
		ref = string(doc.Record.Type) + ":" + doc.Record.ID.String()
	} else if !strings.Contains(ref, "/") {
		switch doc.Record.Type {
		case domain.Objective, domain.Run, domain.Review, domain.Evidence, domain.Handoff, domain.Checkpoint:
			if mission, _ := workspace.String(doc, "mission", false); strings.HasPrefix(mission, "M") && !strings.Contains(mission, ":") {
				ref = mission + "/" + ref
			}
		}
	}
	title := ""
	if doc.Record.Title != nil {
		title = *doc.Record.Title
	}
	status := ""
	if doc.Record.Status != nil {
		status = *doc.Record.Status
	}
	return indexRow{Ref: ref, Noun: string(doc.Record.Type), Title: title, Path: filepath.ToSlash(path), Status: status}
}

// Generated inventories describe only the graph that was validated. Manual
// INDEX.md routes durable context and is never a transaction output.
func generatedIndex(rows []indexRow) ([]byte, error) {
	data, err := json.MarshalIndent(struct {
		Type      string     `json:"type"`
		Version   string     `json:"version"`
		Generated bool       `json:"generated"`
		Scope     string     `json:"scope"`
		Entries   []indexRow `json:"entries"`
	}{"Index", "1.0", true, "governed-records", rows}, "", "  ")
	if err != nil {
		return nil, err
	}
	return append(data, '\n'), nil
}
