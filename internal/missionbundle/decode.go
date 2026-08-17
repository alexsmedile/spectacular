package missionbundle

import (
	"path/filepath"
	"strings"

	"github.com/alexsmedile/spectacular/v2/internal/discovery"
	"github.com/alexsmedile/spectacular/v2/internal/domain"
	"github.com/alexsmedile/spectacular/v2/internal/workspace"
	"go.yaml.in/yaml/v3"
)

func Load(ws *discovery.Workspace, ref string) (*Bundle, error) {
	entry, err := ws.Lookup(ref, domain.Mission)
	if err != nil {
		return nil, err
	}
	return decode(ws, entry)
}

func decode(ws *discovery.Workspace, entry discovery.Entry) (*Bundle, error) {
	doc := entry.Document
	ref, err := compactRef(doc)
	if err != nil {
		return nil, err
	}
	if ref == "" {
		return decodeLegacy(ws, entry)
	}
	b := &Bundle{
		ID:       doc.Record.ID.String(),
		Ref:      ref,
		Title:    value(doc.Record.Title),
		Status:   value(doc.Record.Status),
		Source:   sourceValue(doc),
		Created:  value(doc.Record.Created),
		Updated:  value(doc.Record.Updated),
		Path:     entry.Path,
		Body:     doc.Body,
		entry:    entry,
		document: doc,
	}
	if b.Owner, err = workspace.String(doc, "owner", true); err != nil {
		return nil, err
	}
	if b.Outcome, err = workspace.String(doc, "outcome", true); err != nil {
		return nil, err
	}
	if b.Review, err = workspace.String(doc, "review", true); err != nil {
		return nil, err
	}
	if err = workspace.DecodeValue(doc, "contract", &b.Contract); err != nil {
		return nil, err
	}
	if err = workspace.DecodeValue(doc, "completion", &b.Completion); err != nil {
		return nil, err
	}
	if err = workspace.DecodeValue(doc, "objectives", &b.Objectives); err != nil {
		return nil, err
	}
	if err = workspace.DecodeValue(doc, "validation", &b.Validation); err != nil {
		return nil, err
	}
	if err = workspace.DecodeValue(doc, "authority", &b.Authority); err != nil {
		return nil, err
	}
	if err = workspace.DecodeValue(doc, "scope", &b.Scope); err != nil {
		return nil, err
	}
	if b.RepairBudget, err = workspace.Int(doc, "repair_budget", true); err != nil {
		return nil, err
	}
	if b.Dependencies, err = workspace.Strings(doc, "dependencies", true); err != nil {
		return nil, err
	}
	if b.Gaps, err = workspace.Strings(doc, "gaps", true); err != nil {
		return nil, err
	}
	if b.Stops, err = workspace.Strings(doc, "stops", true); err != nil {
		return nil, err
	}
	if _, err = decodeOptional(doc, "baseline", &b.Baseline); err != nil {
		return nil, err
	}
	if _, err = decodeOptional(doc, "request", &b.Request); err != nil {
		return nil, err
	}
	if _, err = decodeOptional(doc, "fallbacks", &b.Fallbacks); err != nil {
		return nil, err
	}
	if _, err = decodeOptional(doc, "resolves_gaps", &b.ResolvesGaps); err != nil {
		return nil, err
	}
	if _, err = decodeOptional(doc, "after_mission", &b.AfterMission); err != nil {
		return nil, err
	}
	if _, err = decodeOptional(doc, "run", &b.Run); err != nil {
		return nil, err
	}
	if _, err = decodeOptional(doc, "runs", &b.Runs); err != nil {
		return nil, err
	}
	if _, err = decodeOptional(doc, "activation", &b.Activation); err != nil {
		return nil, err
	}
	if _, err = decodeOptional(doc, "reviews", &b.Reviews); err != nil {
		return nil, err
	}
	if _, err = decodeOptional(doc, "handoffs", &b.Handoffs); err != nil {
		return nil, err
	}
	if _, err = decodeOptional(doc, "completion_record", &b.CompletionRecord); err != nil {
		return nil, err
	}
	if err := resolveObjectives(ws, b); err != nil {
		return nil, err
	}
	if err := resolveRuns(ws, b); err != nil {
		return nil, err
	}
	if err := resolveHandoffs(ws, b); err != nil {
		return nil, err
	}
	if err := resolveReviews(ws, b); err != nil {
		return nil, err
	}
	return b, nil
}

func decodeLegacy(ws *discovery.Workspace, entry discovery.Entry) (*Bundle, error) {
	doc := entry.Document
	ref, _, err := workspace.Ref(doc)
	if err != nil {
		return nil, err
	}
	if ref == "" {
		ref = "Mission:" + doc.Record.ID.String()
	}
	b := &Bundle{
		ID:       doc.Record.ID.String(),
		Ref:      ref,
		Title:    value(doc.Record.Title),
		Status:   value(doc.Record.Status),
		Source:   sourceValue(doc),
		Created:  value(doc.Record.Created),
		Updated:  value(doc.Record.Updated),
		Path:     entry.Path,
		Body:     doc.Body,
		Legacy:   true,
		entry:    entry,
		document: doc,
		Validation: Validation{
			Schema: "legacy-v2",
			Mode:   "read-only",
		},
	}
	b.Outcome, _ = workspace.String(doc, "outcome", false)
	if refs, listErr := workspace.Strings(doc, "objectives", false); listErr == nil {
		for _, objectiveRef := range refs {
			objective := Objective{Ref: objectiveRef}
			if target, lookupErr := ws.Lookup(objectiveRef, domain.Objective); lookupErr == nil {
				objective.ID = target.Document.Record.ID.String()
				objective.Outcome = value(target.Document.Record.Title)
				objective.Status = value(target.Document.Record.Status)
				objective.File = target.Path
			}
			b.Objectives = append(b.Objectives, objective)
		}
	}
	return b, nil
}

// resolveHandoffs loads each Handoff a Mission points at so any reader sees the
// record, not just the pointer. It reuses the schema decoder rather than
// restating the fields, so a field added to a Handoff reaches readers without a
// second place to update.
func resolveHandoffs(ws *discovery.Workspace, b *Bundle) error {
	if len(b.Handoffs) == 0 {
		return nil
	}
	base := filepath.Dir(b.entry.Absolute)
	current := workingTree(ws.Root)
	for i := range b.Handoffs {
		pointer := &b.Handoffs[i]
		path, err := containedFile(base, pointer.File)
		if err != nil {
			return err
		}
		doc, err := workspace.ReadFile(path)
		if err != nil {
			return err
		}
		if doc.Record.Type != domain.Handoff || doc.Record.ID.String() != pointer.ID {
			return domain.NewRefusal(domain.RefusalTargetTypeMismatch, "handoffs.file", "Handoff identity or type does not match its pointer", nil)
		}
		resolved, err := decodeHandoff(doc, pointer.File)
		if err != nil {
			return err
		}
		if resolved.Ref != pointer.Ref {
			return domain.NewRefusal(domain.RefusalInvalidReference, "handoffs.file", "Handoff ref does not match its pointer", nil)
		}
		resolved.TreeCurrent = resolved.Reviewed.Tree == current
		pointer.Document = resolved
	}
	return nil
}

func resolveReviews(ws *discovery.Workspace, b *Bundle) error {
	_ = ws
	base := filepath.Dir(b.entry.Absolute)
	for i := range b.Reviews {
		pointer := &b.Reviews[i]
		path, err := containedFile(base, pointer.File)
		if err != nil {
			return err
		}
		doc, err := workspace.ReadFile(path)
		if err != nil {
			return err
		}
		if doc.Record.Type != domain.Review || doc.Record.ID.String() != pointer.ID {
			return domain.NewRefusal(domain.RefusalTargetTypeMismatch, "reviews.file", "Review identity or type does not match its pointer", nil)
		}
		resolved := &Review{
			ID:       doc.Record.ID.String(),
			Title:    value(doc.Record.Title),
			Status:   value(doc.Record.Status),
			Source:   sourceValue(doc),
			Created:  value(doc.Record.Created),
			Path:     pointer.File,
			Body:     doc.Body,
			document: doc,
		}
		if resolved.Ref, err = workspace.String(doc, "ref", true); err != nil {
			return err
		}
		if resolved.Mission, err = workspace.String(doc, "mission", true); err != nil {
			return err
		}
		if err = workspace.DecodeValue(doc, "reviewed", &resolved.Reviewed); err != nil {
			return err
		}
		if err = workspace.DecodeValue(doc, "reviewer", &resolved.Reviewer); err != nil {
			return err
		}
		if err = workspace.DecodeValue(doc, "claims", &resolved.Claims); err != nil {
			return err
		}
		if resolved.Findings, err = workspace.Strings(doc, "findings", true); err != nil {
			return err
		}
		if resolved.Limitations, err = workspace.Strings(doc, "limitations", true); err != nil {
			return err
		}
		if resolved.Ref != pointer.Ref {
			return domain.NewRefusal(domain.RefusalInvalidReference, "reviews.file", "Review ref does not match its pointer", nil)
		}
		pointer.Document = resolved
	}
	return nil
}

func resolveObjectives(ws *discovery.Workspace, b *Bundle) error {
	_ = ws
	base := filepath.Dir(b.entry.Absolute)
	for i := range b.Objectives {
		item := &b.Objectives[i]
		if item.File == "" {
			continue
		}
		path, err := containedFile(base, item.File)
		if err != nil {
			return err
		}
		doc, err := workspace.ReadFile(path)
		if err != nil {
			return err
		}
		if doc.Record.Type != domain.Objective || doc.Record.ID.String() != item.ID {
			return domain.NewRefusal(domain.RefusalTargetTypeMismatch, "objectives.file", "promoted Objective identity or type does not match its pointer", nil)
		}
		var resolved Objective
		if err := decodeRecordFields(doc, &resolved); err != nil {
			return err
		}
		resolved.File = item.File
		resolved.Source = sourceValue(doc)
		resolved.Body = doc.Body
		resolved.document = doc
		if resolved.Ref != item.Ref || resolved.ID != item.ID {
			return domain.NewRefusal(domain.RefusalInvalidReference, "objectives.file", "promoted Objective ref or identity does not match its pointer", nil)
		}
		*item = resolved
	}
	return nil
}

func resolveRuns(ws *discovery.Workspace, b *Bundle) error {
	_ = ws
	if b.Run != nil && b.Run.File != "" {
		return resolveRunFile(b, b.Run)
	}
	for i := range b.Runs {
		if b.Runs[i].File != "" {
			if err := resolveRunFile(b, &b.Runs[i]); err != nil {
				return err
			}
		}
	}
	return nil
}

func resolveRunFile(b *Bundle, item *Run) error {
	path, err := containedFile(filepath.Dir(b.entry.Absolute), item.File)
	if err != nil {
		return err
	}
	doc, err := workspace.ReadFile(path)
	if err != nil {
		return err
	}
	if doc.Record.Type != domain.Run || doc.Record.ID.String() != item.ID {
		return domain.NewRefusal(domain.RefusalTargetTypeMismatch, "runs.file", "promoted Run identity or type does not match its pointer", nil)
	}
	var resolved Run
	if err := decodeRecordFields(doc, &resolved); err != nil {
		return err
	}
	resolved.File = item.File
	resolved.Title = value(doc.Record.Title)
	resolved.Source = sourceValue(doc)
	resolved.Body = doc.Body
	resolved.document = doc
	if resolved.Ref != item.Ref || resolved.ID != item.ID {
		return domain.NewRefusal(domain.RefusalInvalidReference, "runs.file", "promoted Run ref or identity does not match its pointer", nil)
	}
	*item = resolved
	return nil
}

func decodeRecordFields(doc *workspace.Document, output any) error {
	fields := map[string]any{}
	for name, node := range doc.Unknown {
		var value any
		if err := node.Decode(&value); err != nil {
			return domain.NewRefusal(domain.RefusalInvalidKnownField, name, "cannot decode promoted record", err)
		}
		fields[name] = value
	}
	fields["id"] = doc.Record.ID.String()
	fields["status"] = value(doc.Record.Status)
	data, err := yaml.Marshal(fields)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, output)
}

func decodeOptional(doc *workspace.Document, name string, output any) (bool, error) {
	node := doc.Unknown[name]
	if node == nil {
		return false, nil
	}
	if err := node.Decode(output); err != nil {
		return true, domain.NewRefusal(domain.RefusalInvalidKnownField, name, "cannot decode structured value", err)
	}
	return true, nil
}

func compactRef(doc *workspace.Document) (string, error) {
	return workspace.String(doc, "ref", false)
}

func containedFile(base, relative string) (string, error) {
	if relative == "" || filepath.IsAbs(relative) || filepath.Clean(relative) != relative || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return "", domain.NewRefusal(domain.RefusalPathEscape, "file", "bundle pointer must be a canonical relative path", nil)
	}
	path := filepath.Join(base, relative)
	rel, err := filepath.Rel(base, path)
	if err != nil || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", domain.NewRefusal(domain.RefusalPathEscape, "file", "bundle pointer escapes its Mission", err)
	}
	return path, nil
}

func value(input *string) string {
	if input == nil {
		return ""
	}
	return *input
}

func sourceValue(doc *workspace.Document) string {
	if doc == nil || doc.Record.Source == nil {
		return ""
	}
	return doc.Record.Source.String()
}

func scopedRef(raw string) (string, string, error) {
	mission, child, ok := strings.Cut(raw, "/")
	if !ok || mission == "" || child == "" || strings.Contains(child, "/") {
		return "", "", domain.NewRefusal(domain.RefusalInvalidReference, "ref", "expected <mission-ref>/<local-ref>", nil)
	}
	return mission, child, nil
}
