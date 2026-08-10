// Package index provides deterministic, disposable exact lookup over semantic
// records. It does not persist projections or perform fuzzy discovery.
package index

import (
	"fmt"
	"path"
	"sort"
	"strings"

	"github.com/alexsmedile/spectacular/v2/internal/domain"
)

// Entry is the minimum disposable projection needed for exact M1 lookup.
type Entry struct {
	Path   string
	Record domain.Record
}

// Index keys records independently by stable identity and complete
// workspace-relative path.
type Index struct {
	entries []Entry
	byID    map[domain.ID]Entry
	byPath  map[string]Entry
}

// New discovers all entries first, deterministically rejects duplicate keys,
// and only then validates typed relationships. Input order therefore cannot
// change either successful results or refusal selection.
func New(discovered []Entry) (*Index, error) {
	entries := make([]Entry, len(discovered))
	for index, entry := range discovered {
		entries[index] = cloneEntry(entry)
	}
	sort.Slice(entries, func(left, right int) bool {
		if entries[left].Path != entries[right].Path {
			return entries[left].Path < entries[right].Path
		}
		return recordSortKey(entries[left].Record) < recordSortKey(entries[right].Record)
	})

	result := &Index{
		byID:   make(map[domain.ID]Entry, len(entries)),
		byPath: make(map[string]Entry, len(entries)),
	}
	for _, entry := range entries {
		if err := validateWorkspacePath(entry.Path); err != nil {
			return nil, err
		}
		if err := entry.Record.Validate(); err != nil {
			return nil, fmt.Errorf("index %q: %w", entry.Path, err)
		}
		if existing, found := result.byPath[entry.Path]; found {
			return nil, domain.NewRefusal(
				domain.RefusalDuplicatePath,
				"path",
				fmt.Sprintf(
					"%q is claimed by identities %s and %s",
					entry.Path,
					existing.Record.ID,
					entry.Record.ID,
				),
				nil,
			)
		}
		if existing, found := result.byID[entry.Record.ID]; found {
			return nil, domain.NewRefusal(
				domain.RefusalDuplicateID,
				"id",
				fmt.Sprintf(
					"identity %s is claimed by %q and %q",
					entry.Record.ID,
					existing.Path,
					entry.Path,
				),
				nil,
			)
		}
		result.byPath[entry.Path] = entry
		result.byID[entry.Record.ID] = entry
	}

	// Relationship validation deliberately occurs after complete discovery.
	// A target's filesystem position can never determine whether it is found.
	for _, entry := range entries {
		reference := entry.Record.Source
		if reference == nil {
			continue
		}
		target, found := result.byID[reference.ID]
		if !found {
			return nil, domain.NewRefusal(
				domain.RefusalTargetNotFound,
				"source",
				fmt.Sprintf("%q references absent %s", entry.Path, reference),
				nil,
			)
		}
		if target.Record.Type != reference.Type {
			return nil, domain.NewRefusal(
				domain.RefusalTargetTypeMismatch,
				"source",
				fmt.Sprintf(
					"%q expects %s but identity %s at %q is %s",
					entry.Path,
					reference.Type,
					reference.ID,
					target.Path,
					target.Record.Type,
				),
				nil,
			)
		}
	}

	result.entries = append(result.entries, entries...)
	sort.Slice(result.entries, func(left, right int) bool {
		leftID := result.entries[left].Record.ID.String()
		rightID := result.entries[right].Record.ID.String()
		if leftID == rightID {
			return result.entries[left].Path < result.entries[right].Path
		}
		return leftID < rightID
	})
	return result, nil
}

// LookupID performs exact lookup by already-validated stable identity.
func (index *Index) LookupID(id domain.ID) (Entry, bool) {
	entry, found := index.byID[id]
	return cloneEntry(entry), found
}

// LookupPath performs exact lookup by complete canonical workspace-relative
// path. Basenames, cleaned alternatives, and fuzzy matches are not considered.
func (index *Index) LookupPath(workspacePath string) (Entry, bool) {
	entry, found := index.byPath[workspacePath]
	return cloneEntry(entry), found
}

// Entries returns a stable identity-then-path ordered snapshot.
func (index *Index) Entries() []Entry {
	entries := make([]Entry, len(index.entries))
	for position, entry := range index.entries {
		entries[position] = cloneEntry(entry)
	}
	return entries
}

func validateWorkspacePath(workspacePath string) error {
	if workspacePath == "" || strings.Contains(workspacePath, "\\") || path.IsAbs(workspacePath) {
		return invalidPathRefusal(workspacePath)
	}
	cleaned := path.Clean(workspacePath)
	if cleaned == "." || cleaned != workspacePath || strings.HasPrefix(cleaned, "../") {
		return invalidPathRefusal(workspacePath)
	}
	return nil
}

func invalidPathRefusal(workspacePath string) error {
	return domain.NewRefusal(
		domain.RefusalInvalidWorkspacePath,
		"path",
		fmt.Sprintf("%q is not a canonical workspace-relative path", workspacePath),
		nil,
	)
}

func cloneEntry(entry Entry) Entry {
	entry.Record = cloneRecord(entry.Record)
	return entry
}

func cloneRecord(record domain.Record) domain.Record {
	result := record
	result.Title = cloneString(record.Title)
	result.Description = cloneString(record.Description)
	result.Status = cloneString(record.Status)
	result.CreatedBy = cloneString(record.CreatedBy)
	result.Created = cloneString(record.Created)
	result.Updated = cloneString(record.Updated)
	if record.Source != nil {
		source := *record.Source
		result.Source = &source
	}
	return result
}

func cloneString(value *string) *string {
	if value == nil {
		return nil
	}
	cloned := *value
	return &cloned
}

func recordSortKey(record domain.Record) string {
	source := ""
	if record.Source != nil {
		source = record.Source.String()
	}
	return fmt.Sprintf(
		"%q|%q|%q|%q|%q|%q|%q|%q|%q|%d",
		record.ID.String(),
		record.Type,
		optionalSortKey(record.Title),
		optionalSortKey(record.Description),
		optionalSortKey(record.Status),
		optionalSortKey(record.CreatedBy),
		optionalSortKey(record.Created),
		optionalSortKey(record.Updated),
		source,
		len(source),
	)
}

func optionalSortKey(value *string) string {
	if value == nil {
		return "absent"
	}
	return "present:" + *value
}
