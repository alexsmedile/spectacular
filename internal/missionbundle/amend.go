package missionbundle

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/alexsmedile/spectacular/v2/internal/discovery"
	"github.com/alexsmedile/spectacular/v2/internal/domain"
	"github.com/alexsmedile/spectacular/v2/internal/governance"
	"github.com/alexsmedile/spectacular/v2/internal/workspace"
)

// Amendment is what `contract amend` reports: what changed, or under --dry-run what
// would change. It is the preview and the receipt in one shape, so the owner reads
// the same description before and after.
type Amendment struct {
	Operation  string   `json:"operation"`
	Contract   string   `json:"contract"`
	Path       string   `json:"path"`
	Gap        string   `json:"gap"`
	Resolution string   `json:"resolution"`
	Mission    string   `json:"mission"`
	Source     string   `json:"source"`
	From       string   `json:"from"`
	To         string   `json:"to"`
	Repointed  []string `json:"repointed,omitempty"`
	Log        string   `json:"log"`
	DryRun     bool     `json:"dry_run,omitempty"`
	NoOp       bool     `json:"no_op,omitempty"`
	Changed    []string `json:"changed,omitempty"`
}

// amendableFields are the only frontmatter keys an amendment may reach. Everything
// absent from this set states what was agreed, and changing it is a new Contract
// version rather than an amendment — otherwise "amendment" is "rewrite with extra
// steps" and the only thing separating them is a reason string nothing validates.
var amendableFields = map[string]bool{"gaps": true, "updated": true}

// blockedOnLine matches one Gap entry's blocked_on key, including a folded or
// literal block scalar and its indented continuation lines. The rewrite is textual
// rather than a decode-and-re-emit because re-emitting canonical YAML would reflow
// the whole Contract, and an amendment that reformats prose it did not change is
// not reviewable.
var blockedOnLine = regexp.MustCompile(`(?m)^(\s*)blocked_on:(.*)$`)

// AmendContract closes one Gap on a Contract by rewriting its blocked_on to
// resolution, using the text a Mission froze in its resolves_gaps declaration.
// override carries an owner-supplied resolution for a Gap that predates the
// resolves_gaps declaration. It is the deliberate exception: the mechanism cannot
// reach back before its own existence, and the Gaps that motivated it are all
// older than it. Recorded in the log as owner-supplied so a reader can tell an
// approved-at-activation wording from one typed at a prompt.
func (s Service) AmendContract(contractRef, gapRef, owner, override string, dryRun bool) (Amendment, error) {
	if !dryRun {
		locked, unlock, err := s.beginMutation()
		if err != nil {
			return Amendment{}, err
		}
		defer unlock()
		return locked.amendContract(contractRef, gapRef, owner, override, false)
	}
	return s.amendContract(contractRef, gapRef, owner, override, true)
}

func (s Service) amendContract(contractRef, gapRef, owner, override string, dryRun bool) (Amendment, error) {
	if strings.TrimSpace(owner) == "" {
		return Amendment{}, invalid("by", "an amendment requires the owner who authorizes it")
	}
	entry, err := s.Workspace.Lookup(contractRef, domain.Contract)
	if err != nil {
		return Amendment{}, err
	}
	// The caller may name the Contract by its human ref; Missions bind by exact
	// typed reference. Everything downstream compares against what a Mission stores,
	// so normalize once here rather than at each comparison.
	contractRef = string(domain.Contract) + ":" + entry.Document.Record.ID.String()
	original, err := os.ReadFile(entry.Absolute)
	if err != nil {
		return Amendment{}, invalidCause("contract", "cannot read the Contract", err)
	}
	from := digest(original)

	// The resolution text is never composed here. It comes from the Mission that
	// declared it, frozen in that Mission's activation fingerprint, so the words
	// written are the words the owner approved at the activation gate.
	mission, resolution, err := s.declaringMission(entry.Document.Record.ID.String(), contractRef, gapRef)
	source := "mission-declared"
	if override != "" {
		// An override does not consult a Mission at all: a Gap older than the
		// declaration has no Mission that could have frozen wording for it.
		mission, resolution, err, source = "none", override, nil, "owner-supplied"
	}
	if err != nil {
		return Amendment{}, err
	}

	// A live bound Mission still has the Contract constraining work in flight.
	if live, ref := s.liveBoundMission(contractRef); live {
		return Amendment{}, domain.NewStateRefusal(domain.RefusalUnauthorized, "contract",
			"a Mission bound to this Contract is live", "no live bound Mission", ref,
			"complete or stop "+ref+" before amending the Contract it is bound to", nil)
	}

	gaps, err := ContractGaps(s.Workspace, contractRef)
	if err != nil {
		return Amendment{}, err
	}
	var target *contractGap
	for i := range gaps {
		if gaps[i].Ref == gapRef {
			target = &gaps[i]
			break
		}
	}
	if target == nil {
		return Amendment{}, invalid("gap", "no such Gap on this Contract: "+gapRef)
	}

	result := Amendment{
		Operation: "contract.amend", Contract: contractRef, Path: entry.Path, Gap: gapRef,
		Resolution: resolution, Mission: mission, Source: source, From: from, DryRun: dryRun,
	}

	// An already-closed Gap whose text matches is a no-op rather than a refusal:
	// refusing an identical rewrite is pedantry that gets worked around. Differing
	// text refuses, because two Missions disagreeing about why a Gap closed is a
	// question for the owner.
	if target.Resolution != "" {
		if strings.TrimSpace(target.Resolution) == strings.TrimSpace(resolution) {
			result.NoOp, result.To = true, from
			return result, nil
		}
		return Amendment{}, domain.NewStateRefusal(domain.RefusalInvalidKnownField, "gap",
			"Gap "+gapRef+" already carries a different resolution", target.Resolution, resolution,
			"reconcile the two resolutions with the owner before amending", nil)
	}
	if target.BlockedOn == "" {
		return Amendment{}, invalid("gap", "Gap "+gapRef+" carries neither blocked_on nor resolution")
	}

	amended, err := rewriteGap(string(original), gapRef, resolution)
	if err != nil {
		return Amendment{}, err
	}
	if err := assertOnlyAmendableFieldsChanged(string(original), amended); err != nil {
		return Amendment{}, err
	}
	result.To = digest([]byte(amended))

	logPath := amendmentLogPath(entry.Path)
	logData, err := appendAmendmentLog(s.Workspace.Root, logPath, amendmentEntry{
		At: s.now(), By: owner, Mission: mission, Gap: gapRef, Contract: contractRef,
		Source: source, From: from, To: result.To,
	})
	if err != nil {
		return Amendment{}, err
	}
	result.Log = logPath

	repointed, changes, err := s.repointBoundMissions(contractRef, result.To)
	if err != nil {
		return Amendment{}, err
	}
	result.Repointed = repointed

	changes = append(changes,
		governance.FileChange{Path: entry.Path, Data: []byte(amended), Mode: 0o644},
		governance.FileChange{Path: logPath, Data: logData, Mode: 0o644},
	)
	for _, change := range changes {
		result.Changed = append(result.Changed, filepath.ToSlash(change.Path))
	}
	if dryRun {
		return result, nil
	}

	apply := s.ApplyTransaction
	if apply == nil {
		apply = governance.ApplyTransaction
	}
	if err := apply(s.Workspace.Root, "contract.amend:"+contractRef+":"+gapRef+":"+result.To, changes); err != nil {
		return Amendment{}, err
	}
	return result, nil
}

// declaringMission finds the Mission that froze a resolution for this Gap. An
// amendment with no declaring Mission is refused: the resolution text has to come
// from a boundary an owner approved, not from the command line.
func (s Service) declaringMission(contractID, contractRef, gapRef string) (string, string, error) {
	for _, entry := range s.Workspace.Entries {
		if entry.Document == nil || entry.Document.Record.Type != domain.Mission {
			continue
		}
		bundle, err := Load(s.Workspace, entry.Document.Record.ID.String())
		if err != nil || bundle.Contract.Ref != contractRef {
			continue
		}
		for _, declared := range bundle.ResolvesGaps {
			if declared.Gap == gapRef {
				return bundle.Ref, declared.Resolution, nil
			}
		}
	}
	return "", "", domain.NewStateRefusal(domain.RefusalMissingRequiredField, "gap",
		"no Mission bound to this Contract declares it resolves "+gapRef, "a frozen resolves_gaps entry", "none",
		"declare the Gap in the resolving Mission's resolves_gaps so the owner approves the wording at activation", nil)
}

func (s Service) liveBoundMission(contractRef string) (bool, string) {
	for _, entry := range s.Workspace.Entries {
		if entry.Document == nil || entry.Document.Record.Type != domain.Mission {
			continue
		}
		bundle, err := Load(s.Workspace, entry.Document.Record.ID.String())
		if err != nil || bundle.Contract.Ref != contractRef {
			continue
		}
		if bundle.Status == "active" {
			return true, bundle.Ref
		}
	}
	return false, ""
}

// repointBoundMissions updates only contract.fingerprint on every bound Mission.
// Re-pointing is not re-activation: the frozen envelope is untouched, so a
// completed Mission's plan and any review bound to its activation stay valid.
func (s Service) repointBoundMissions(contractRef, fingerprint string) ([]string, []governance.FileChange, error) {
	var refs []string
	var changes []governance.FileChange
	for _, entry := range s.Workspace.Entries {
		if entry.Document == nil || entry.Document.Record.Type != domain.Mission {
			continue
		}
		bundle, err := Load(s.Workspace, entry.Document.Record.ID.String())
		if err != nil || bundle.Contract.Ref != contractRef || bundle.Contract.Fingerprint == fingerprint {
			continue
		}
		data, err := os.ReadFile(entry.Absolute)
		if err != nil {
			return nil, nil, invalidCause("mission", "cannot read a bound Mission", err)
		}
		updated := strings.Replace(string(data), bundle.Contract.Fingerprint, fingerprint, 1)
		if updated == string(data) {
			continue
		}
		refs = append(refs, bundle.Ref)
		changes = append(changes, governance.FileChange{Path: entry.Path, Data: []byte(updated), Mode: 0o644})
	}
	return refs, changes, nil
}

// rewriteGap replaces one Gap's blocked_on with resolution, in place, leaving every
// other byte of the Contract untouched.
func rewriteGap(contract, gapRef, resolution string) (string, error) {
	lines := strings.Split(contract, "\n")
	start := -1
	for i, line := range lines {
		if strings.TrimSpace(line) == "- ref: "+gapRef || strings.TrimSpace(line) == "ref: "+gapRef {
			start = i
			break
		}
	}
	if start < 0 {
		return "", invalid("gap", "cannot locate Gap "+gapRef+" in the Contract text")
	}
	for i := start + 1; i < len(lines); i++ {
		trimmed := strings.TrimSpace(lines[i])
		if strings.HasPrefix(trimmed, "- ") || (trimmed != "" && !strings.HasPrefix(lines[i], " ")) {
			break // next Gap entry, or the end of the block
		}
		match := blockedOnLine.FindStringSubmatch(lines[i])
		if match == nil {
			continue
		}
		indent := match[1]
		end := i + 1
		for end < len(lines) {
			next := lines[end]
			if strings.TrimSpace(next) == "" || len(next)-len(strings.TrimLeft(next, " ")) > len(indent) {
				end++
				continue
			}
			break
		}
		replacement := []string{indent + "resolution: >-"}
		for _, wrapped := range foldText(resolution, indent+"    ", 96) {
			replacement = append(replacement, wrapped)
		}
		rewritten := append([]string{}, lines[:i]...)
		rewritten = append(rewritten, replacement...)
		rewritten = append(rewritten, lines[end:]...)
		return strings.Join(rewritten, "\n"), nil
	}
	return "", invalid("gap", "Gap "+gapRef+" carries no blocked_on to rewrite")
}

// foldText wraps prose into a YAML block scalar body. Long single-line values are
// where a bare colon silently breaks the document, so the emitted form is always a
// block scalar rather than a plain scalar.
func foldText(text, indent string, width int) []string {
	words := strings.Fields(text)
	var lines []string
	current := indent
	for _, word := range words {
		if current != indent && len(current)+1+len(word) > width {
			lines = append(lines, current)
			current = indent + word
			continue
		}
		if current == indent {
			current += word
			continue
		}
		current += " " + word
	}
	if strings.TrimSpace(current) != "" {
		lines = append(lines, current)
	}
	return lines
}

// assertOnlyAmendableFieldsChanged is the guard that keeps an amendment from
// becoming a rewrite. It compares top-level frontmatter keys before and after and
// refuses if any key outside amendableFields moved.
func assertOnlyAmendableFieldsChanged(before, after string) error {
	changed, err := changedTopLevelKeys(before, after)
	if err != nil {
		return err
	}
	for _, key := range changed {
		if !amendableFields[key] {
			return domain.NewStateRefusal(domain.RefusalUnauthorized, key,
				"an amendment may not change a field that states what was agreed", "gaps or editorial fields", key,
				"route a semantic change through a new contract_version rather than an amendment", nil)
		}
	}
	return nil
}

func changedTopLevelKeys(before, after string) ([]string, error) {
	first, err := topLevelBlocks(before)
	if err != nil {
		return nil, err
	}
	second, err := topLevelBlocks(after)
	if err != nil {
		return nil, err
	}
	seen := map[string]bool{}
	var changed []string
	for key, value := range first {
		if second[key] != value {
			changed, seen[key] = append(changed, key), true
		}
	}
	for key := range second {
		if _, ok := first[key]; !ok && !seen[key] {
			changed = append(changed, key)
		}
	}
	return changed, nil
}

// topLevelBlocks splits frontmatter into top-level key to raw-body pairs. A textual
// split rather than a decode: the comparison must notice a reflowed block scalar,
// which a decode would normalize away.
func topLevelBlocks(record string) (map[string]string, error) {
	_, frontmatter, _, err := splitRecord(record)
	if err != nil {
		return nil, err
	}
	blocks := map[string]string{}
	key := ""
	var builder strings.Builder
	for _, line := range strings.Split(frontmatter, "\n") {
		if line != "" && !strings.HasPrefix(line, " ") && !strings.HasPrefix(line, "#") && strings.Contains(line, ":") {
			if key != "" {
				blocks[key] = builder.String()
			}
			key = strings.TrimSpace(strings.SplitN(line, ":", 2)[0])
			builder.Reset()
		}
		builder.WriteString(line + "\n")
	}
	if key != "" {
		blocks[key] = builder.String()
	}
	return blocks, nil
}

func splitRecord(record string) (string, string, string, error) {
	if !strings.HasPrefix(record, "---\n") {
		return "", "", "", invalid("contract", "record does not open with YAML frontmatter")
	}
	rest := record[4:]
	end := strings.Index(rest, "\n---\n")
	if end < 0 {
		return "", "", "", invalid("contract", "record frontmatter is unterminated")
	}
	return "---\n", rest[:end], rest[end+5:], nil
}

func digest(data []byte) string {
	sum := sha256.Sum256(data)
	return "sha256:" + hex.EncodeToString(sum[:])
}

type amendmentEntry struct {
	At, By, Mission, Gap, Contract, Source, From, To string
}

// amendmentLogPath names the companion record beside the Contract. The log lives
// outside the Contract so an entry can record the fingerprint its own amendment
// produced, which an entry inside the fingerprinted file cannot.
func amendmentLogPath(contractPath string) string {
	return strings.TrimSuffix(filepath.ToSlash(contractPath), ".md") + ".amendments.md"
}

func appendAmendmentLog(root, path string, entry amendmentEntry) ([]byte, error) {
	existing, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(path)))
	if err != nil && !os.IsNotExist(err) {
		return nil, invalidCause("amendments", "cannot read the amendment log", err)
	}
	body := string(existing)
	if body == "" {
		body = "<!-- append-only amendment log; each entry records one amendment to the Contract beside it -->\n" +
			"# Amendments\n"
	}
	body += fmt.Sprintf("\n- at: %q\n  by: %s\n  mission: %s\n  gap: %s\n  contract: %s\n  source: %s\n  fields: [gaps]\n  from: %s\n  to: %s\n",
		entry.At, entry.By, entry.Mission, entry.Gap, entry.Contract, entry.Source, entry.From, entry.To)
	return []byte(body), nil
}

// assertDeclaredGapsClosed refuses completion while any Gap the Mission declared it
// resolves still reads blocked_on. The refusal names the Gap and the command that
// closes it, so a reader learns the fix rather than only the problem — the failure
// the original stale_fingerprint refusal made, promising an amend path that did not
// exist.
func (s Service) assertDeclaredGapsClosed(b *Bundle) error {
	if len(b.ResolvesGaps) == 0 {
		return nil
	}
	gaps, err := ContractGaps(s.Workspace, b.Contract.Ref)
	if err != nil {
		return err
	}
	state := make(map[string]contractGap, len(gaps))
	for _, gap := range gaps {
		state[gap.Ref] = gap
	}
	for _, declared := range b.ResolvesGaps {
		gap, known := state[declared.Gap]
		if !known {
			return invalid("resolves_gaps.gap", "declared Gap is no longer on the bound Contract: "+declared.Gap)
		}
		if gap.Resolution != "" {
			continue
		}
		return domain.NewStateRefusal(domain.RefusalInvalidKnownField, "resolves_gaps.gap",
			"Gap "+declared.Gap+" is still open and this Mission declared it would be closed",
			"resolution", "blocked_on",
			"close it with: spectacular contract amend "+b.Contract.Ref+" --gap "+declared.Gap+" --by "+b.Owner, nil)
	}
	return nil
}

// ContractVersion reads a Contract's declared version. The field existed on every
// Contract and was read by nothing, which is why one Contract reaching version 2
// meant nothing mechanically. Versioning is where semantic change goes, so the field
// has to be validated and reported or the rule routes change into a field no reader
// can trust.
func ContractVersion(ws *discovery.Workspace, ref string) (int, error) {
	entry, err := ws.Lookup(ref, domain.Contract)
	if err != nil {
		return 0, err
	}
	// Existing Contracts quote the value, so both `contract_version: 2` and
	// `contract_version: "2"` are accepted. Rewriting the records to unquote it
	// would be a Contract edit for a cosmetic reason.
	version, err := workspace.Int(entry.Document, "contract_version", true)
	if err != nil {
		text, textErr := workspace.String(entry.Document, "contract_version", true)
		if textErr != nil {
			return 0, invalidCause("contract_version", "a Contract must declare contract_version", err)
		}
		version, err = strconv.Atoi(strings.TrimSpace(text))
		if err != nil {
			return 0, invalidCause("contract_version", "contract_version must be an integer", err)
		}
	}
	if version < 1 {
		return 0, invalid("contract_version", "contract_version must be a positive integer")
	}
	return version, nil
}

// validateContractVersion reports the bound Contract's version and never refuses
// over it. A Mission bound to an earlier version is simply outdated: it ran against
// that version, which is a true fact about it and not a problem to solve. No
// migration and no superseded copies — the history of what changed is the commit
// history.
//
// A malformed version refuses, because a version nobody can read is worse than
// none. An absent one does not: requiring it would refuse every Mission bound to a
// Contract authored before the field was read, which is the same
// no-legal-correction failure this Mission exists to remove.
func validateContractVersion(ws *discovery.Workspace, b *Bundle) error {
	version, err := ContractVersion(ws, b.Contract.Ref)
	if err != nil {
		if domain.RefusalHasCode(err, domain.RefusalMissingRequiredField) {
			return nil
		}
		var refusal *domain.Refusal
		if errors.As(err, &refusal) && refusal.Field == "contract_version" &&
			strings.Contains(refusal.Actual, "required property is absent") {
			return nil
		}
		return err
	}
	b.contractVersion = version
	return nil
}
