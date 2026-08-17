package missionbundle

import (
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/alexsmedile/spectacular/v2/internal/discovery"
	"github.com/alexsmedile/spectacular/v2/internal/domain"
	"github.com/alexsmedile/spectacular/v2/internal/workspace"
)

// handoffPattern is the compact ref a Handoff carries inside a Mission bundle.
// A Handoff ref ends in a short key derived from its identity, because a Handoff
// is a leaf record like Evidence rather than an ordinal like a Review.
var handoffPattern = regexp.MustCompile(`^H[1-9][0-9]*-[a-z2-7]{6}$`)

// HandoffBinding is the commit and tree a Handoff was sent against. It is not
// Reviewed: a Review also binds the Mission's activation fingerprint, because a
// review judges a Mission's frozen claims. A Handoff delegates work against a
// state of the repository and has no verdict on the Mission, so carrying an
// activation fingerprint would write an empty field into every Handoff and
// imply a binding the record does not have.
type HandoffBinding struct {
	Commit string `yaml:"commit" json:"commit"`
	Tree   string `yaml:"tree" json:"tree"`
}

// HandoffPointer is the Mission's record that a Handoff exists, mirroring
// ReviewPointer. The Mission owns the list; the Handoff owns its own content.
type HandoffPointer struct {
	Ref      string   `yaml:"ref" json:"ref"`
	ID       string   `yaml:"id" json:"id"`
	File     string   `yaml:"file" json:"file"`
	Document *Handoff `yaml:"-" json:"document,omitempty"`
}

// Sender is the identity a Handoff is filed under and that identity's relation
// to the receiver. Accountability for a delegated task stays with the sender, so
// a Handoff that does not say who sent it records nothing that can be answered
// for.
type Sender struct {
	Actor              string `yaml:"actor" json:"actor"`
	RelationToReceiver string `yaml:"relation_to_receiver" json:"relation_to_receiver"`
}

// Handoff is a delegation stated as a record: what the sender verified, what the
// sender is only assuming, and what the receiver owes back.
//
// Asserted and Assumed are the point of the record. They are required to be
// present and are never scored — no validator here or elsewhere judges whether a
// sender actually verified something filed under Asserted. That is deliberate and
// is stated at the field rather than left for a reader to infer from the absence
// of a check: the distinction is a claim the sender signs, and a mechanism that
// scored it would be asserting a fact it cannot know. The receiver re-verifies
// what arrives under Assumed before acting on it.
type Handoff struct {
	ID       string         `json:"id"`
	Ref      string         `json:"ref"`
	Title    string         `json:"title"`
	Status   string         `json:"status,omitempty"`
	Source   string         `json:"source,omitempty"`
	Mission  string         `json:"mission"`
	Created  string         `json:"created,omitempty"`
	Reviewed HandoffBinding `json:"reviewed"`
	Sender   Sender         `json:"sender"`
	Task     string         `json:"task"`
	Asserted []string       `json:"asserted"`
	Assumed  []string       `json:"assumed"`
	Stops    []string       `json:"stops"`
	Returns  []string       `json:"returns"`
	// Supersedes names the Handoff this one corrects, empty when it corrects
	// nothing. A Handoff is never edited; a correction is a new Handoff. The
	// superseded record survives as what its sender believed at the time.
	Supersedes string `json:"supersedes,omitempty"`
	Path       string `json:"path"`
	// TreeCurrent reports whether the tree this Handoff bound is still the tree on
	// disk. A receiver acts on the state the sender bound, so a stale binding is
	// the difference between a usable instruction and a pointer into a state that
	// no longer exists. It is derived on read and never stored.
	TreeCurrent bool   `yaml:"-" json:"tree_current"`
	Body        string `json:"-"`

	document *workspace.Document
}

// decodeHandoff reads a Handoff document into its schema, refusing any field
// that is required but absent. The refusal names the field so a sender correcting
// a rejected Handoff is told which one, rather than that the record is invalid.
func decodeHandoff(doc *workspace.Document, path string) (*Handoff, error) {
	if doc.Record.Type != domain.Handoff {
		return nil, invalid("handoffs.file", "handoff pointer must resolve to a Handoff record")
	}
	h := &Handoff{
		ID:      doc.Record.ID.String(),
		Ref:     workspace.RefOrEmpty(doc),
		Title:   value(doc.Record.Title),
		Status:  value(doc.Record.Status),
		Source:  sourceValue(doc),
		Created: value(doc.Record.Created),
		Path:    path,
		Body:    doc.Body,

		document: doc,
	}
	h.Mission, _ = workspace.String(doc, "mission", false)
	h.Task, _ = workspace.String(doc, "task", false)
	h.Supersedes, _ = workspace.String(doc, "supersedes", false)
	if err := workspace.DecodeValue(doc, "reviewed", &h.Reviewed); err != nil {
		return nil, err
	}
	if err := workspace.DecodeValue(doc, "sender", &h.Sender); err != nil {
		return nil, err
	}
	// These lists are decoded as required, which distinguishes absent from empty.
	// An absent asserted: refuses; an empty one is legal. A sender who verified
	// nothing states that, and a sender who assumed nothing states that too.
	// Absence is silence, and the asserted/assumed split only means something when
	// both lists are deliberately filed.
	for _, field := range []struct {
		name   string
		target *[]string
	}{
		{"asserted", &h.Asserted},
		{"assumed", &h.Assumed},
		{"stops", &h.Stops},
		{"returns", &h.Returns},
	} {
		values, err := workspace.Strings(doc, field.name, true)
		if err != nil {
			return nil, err
		}
		*field.target = values
	}
	return h, nil
}

// validateHandoffs checks every Handoff a Mission points at. It runs on every
// validation rather than only at completion, because a Handoff is read by a
// receiver while the Mission is still live — a Handoff that only became checkable
// at completion would be checked after the work it delegated was already done.
func validateHandoffs(ws *discovery.Workspace, b *Bundle) error {
	if len(b.Handoffs) == 0 {
		return nil
	}
	base := filepath.Dir(b.entry.Absolute)
	seenRefs, seenIDs := map[string]bool{}, map[string]bool{}
	byRef := map[string]*Handoff{}
	for i := range b.Handoffs {
		pointer := &b.Handoffs[i]
		if _, err := domain.ParseID(pointer.ID); err != nil {
			return invalidCause("handoffs.id", "must be canonical UUIDv7", err)
		}
		if !handoffPattern.MatchString(pointer.Ref) || seenRefs[pointer.Ref] || seenIDs[pointer.ID] {
			return invalid("handoffs.ref", "handoff refs and identities must be unique H<number>-<key> values")
		}
		seenRefs[pointer.Ref], seenIDs[pointer.ID] = true, true

		path, err := containedFile(base, pointer.File)
		if err != nil {
			return err
		}
		doc, err := workspace.ReadFile(path)
		if err != nil {
			return err
		}
		handoff, err := decodeHandoff(doc, path)
		if err != nil {
			return err
		}
		if handoff.ID != pointer.ID {
			return invalid("handoffs.file", "handoff pointer must resolve to the same Handoff identity")
		}
		if err := validateHandoffContent(handoff, b); err != nil {
			return err
		}
		pointer.Document = handoff
		byRef[pointer.Ref] = handoff
	}
	// supersedes is resolved after every Handoff is known, so a correction may be
	// recorded in any order relative to the record it corrects.
	for ref, handoff := range byRef {
		if handoff.Supersedes == "" {
			continue
		}
		if handoff.Supersedes == ref {
			return invalid("handoff.supersedes", "a Handoff cannot supersede itself")
		}
		if byRef[handoff.Supersedes] == nil {
			return invalid("handoff.supersedes", "supersedes must name a Handoff recorded on this Mission")
		}
	}
	return nil
}

// HandoffDraft is the input a sender writes. It is deliberately not the Handoff
// record: the sender states the delegation, and the command derives identity,
// ref, and placement so two senders filing the same delegation cannot disagree
// about where it lives.
type HandoffDraft struct {
	Type     string         `yaml:"type"`
	Title    string         `yaml:"title"`
	Reviewed HandoffBinding `yaml:"reviewed"`
	Sender   Sender         `yaml:"sender"`
	Task     string         `yaml:"task"`
	// Asserted and Assumed are pointers so an absent list is distinguishable
	// from an empty one at the command boundary, the same distinction the
	// schema draws on disk. A sender who omits asserted: is told to state it.
	Asserted   *[]string `yaml:"asserted"`
	Assumed    *[]string `yaml:"assumed"`
	Stops      []string  `yaml:"stops"`
	Returns    []string  `yaml:"returns"`
	Supersedes string    `yaml:"supersedes"`
}

// validateHandoffContent enforces the schema of one Handoff against the Mission
// that carries it.
func validateHandoffContent(h *Handoff, b *Bundle) error {
	// The Mission reference is checked against the Mission that points at this
	// Handoff. A Handoff naming a different Mission is a record filed in the wrong
	// bundle, which no later reader would catch.
	if h.Mission == "" {
		return invalid("handoff.mission", "a Handoff must name the Mission it belongs to")
	}
	if h.Mission != b.Ref && h.Mission != b.ID {
		if typed, err := domain.ParseReference(h.Mission); err != nil || typed.ID.String() != b.ID {
			return invalid("handoff.mission", "handoff must name the Mission that carries it")
		}
	}
	if h.Title == "" {
		return invalid("handoff.title", "a Handoff must have a title")
	}
	if h.Task == "" {
		return invalid("handoff.task", "a Handoff must state the task in the receiver's terms")
	}
	if h.Sender.Actor == "" || h.Sender.RelationToReceiver == "" {
		return invalid("handoff.sender", "a Handoff must name its sender and that sender's relation to the receiver")
	}
	if len(h.Stops) == 0 {
		return invalid("handoff.stops", "a Handoff must state at least one stop")
	}
	if len(h.Returns) == 0 {
		return invalid("handoff.returns", "a Handoff must state what the receiver returns")
	}
	// asserted is what the sender claims to have verified. An empty asserted list
	// is a Handoff that verified nothing, which is a real thing to send but not a
	// thing to send silently, so it is stated rather than defaulted.
	if !commitPattern.MatchString(h.Reviewed.Commit) || !commitPattern.MatchString(h.Reviewed.Tree) {
		return invalid("handoff.reviewed", "a Handoff must bind an exact commit and tree")
	}
	return nil
}

// workingTree reports the tree a Handoff would be compared against: HEAD's tree
// when the working tree is clean, and empty when it is dirty or unreadable.
//
// git write-tree would give an exact answer for a dirty tree, but it writes
// objects into the repository, and a read command must not mutate the thing it
// is reporting on. An uncommitted change means the sender's bound state is not
// what is on disk, which is exactly what "moved" reports, so the cheaper
// non-mutating check answers the question the reader is asking.
func workingTree(root string) string {
	status := exec.Command("git", "status", "--porcelain")
	status.Dir = root
	dirty, err := status.Output()
	if err != nil || len(strings.TrimSpace(string(dirty))) > 0 {
		return ""
	}
	head := exec.Command("git", "rev-parse", "HEAD^{tree}")
	head.Dir = root
	out, err := head.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}
