package domain

import (
	"errors"
	"fmt"
)

// RefusalCode is a stable, machine-readable reason that an invalid record or
// workspace state was rejected.
type RefusalCode string

const (
	RefusalInvalidID            RefusalCode = "invalid_id"
	RefusalInvalidType          RefusalCode = "invalid_type"
	RefusalInvalidReference     RefusalCode = "invalid_reference"
	RefusalMissingRequiredField RefusalCode = "missing_required_field"
	RefusalInvalidKnownField    RefusalCode = "invalid_known_field"
	RefusalDuplicateID          RefusalCode = "duplicate_id"
	RefusalDuplicatePath        RefusalCode = "duplicate_path"
	RefusalTargetNotFound       RefusalCode = "target_not_found"
	RefusalTargetTypeMismatch   RefusalCode = "target_type_mismatch"
	RefusalInvalidWorkspacePath RefusalCode = "invalid_workspace_path"
	RefusalInvalidFrontmatter   RefusalCode = "invalid_frontmatter"
	RefusalInvalidUTF8          RefusalCode = "invalid_utf8"
	RefusalPersistence          RefusalCode = "persistence_failed"
)

// Refusal describes a deterministic rejection without granting the caller an
// authorization or lifecycle decision.
type Refusal struct {
	Code   RefusalCode
	Field  string
	Detail string
	Cause  error
}

func (r *Refusal) Error() string {
	location := ""
	if r.Field != "" {
		location = " field " + r.Field
	}
	if r.Detail != "" {
		return fmt.Sprintf("refused %s%s: %s", r.Code, location, r.Detail)
	}
	return fmt.Sprintf("refused %s%s", r.Code, location)
}

func (r *Refusal) Unwrap() error {
	return r.Cause
}

// NewRefusal creates a refusal shared by the Domain, Workspace, and Index
// packages. Detail is diagnostic text, not semantic authority.
func NewRefusal(code RefusalCode, field, detail string, cause error) error {
	return &Refusal{Code: code, Field: field, Detail: detail, Cause: cause}
}

// RefusalHasCode reports whether err contains a refusal with the supplied
// stable code.
func RefusalHasCode(err error, code RefusalCode) bool {
	var refusal *Refusal
	return errors.As(err, &refusal) && refusal.Code == code
}
