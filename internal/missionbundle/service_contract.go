package missionbundle

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"

	"github.com/alexsmedile/spectacular/v2/internal/discovery"
	"github.com/alexsmedile/spectacular/v2/internal/domain"
	"github.com/alexsmedile/spectacular/v2/internal/humanlayout"
	"github.com/alexsmedile/spectacular/v2/internal/workspace"
	"go.yaml.in/yaml/v3"
)

func (s Service) CreateContract(ref, title, owner string) (Result, error) {
	locked, unlock, err := s.beginMutation()
	if err != nil {
		return Result{}, err
	}
	defer unlock()
	return locked.createContract(ref, title, owner)
}

func (s Service) createContract(ref, title, owner string) (Result, error) {
	ref = strings.TrimSpace(ref)
	if ref == "" {
		return Result{}, invalid("ref", "contract ref is required (e.g. CC-my-feature)")
	}
	if !strings.HasPrefix(ref, "CC-") {
		ref = "CC-" + ref
	}
	if title == "" {
		title = strings.TrimPrefix(ref, "CC-")
	}
	if owner == "" {
		owner = s.Workspace.Config.Defaults.Operator
		if owner == "" {
			owner = "Alex"
		}
	}
	now := s.now()
	contractID, err := domain.NewID()
	if err != nil {
		return Result{}, err
	}
	slug := humanlayout.Slug(title)
	fileName := ref + "-" + slug + ".md"
	relPath := filepath.ToSlash(filepath.Join(".spectacular", "contracts", fileName))

	doc := &workspace.Document{
		Record: domain.Record{
			Type:    domain.Contract,
			ID:      contractID,
			Title:   stringPtr(title),
			Status:  stringPtr("current"),
			Created: stringPtr(now),
			Updated: stringPtr(now),
		},
		Unknown: map[string]*yaml.Node{},
		Body:    "# " + title + "\n\n" + "## Overview\n\n" + title + " contract specification.\n",
	}
	workspace.SetString(doc, "ref", ref)
	workspace.SetString(doc, "owner", owner)
	workspace.SetString(doc, "contract_version", "1")
	workspace.SetString(doc, "purpose", title)
	workspace.SetString(doc, "outcome", title)
	workspace.SetStrings(doc, "applies_when", []string{"Work on " + title + " begins."})
	workspace.SetStrings(doc, "does_not_apply_when", []string{})
	workspace.SetStrings(doc, "does_not_provide", []string{})
	workspace.SetStrings(doc, "required_behavior", []string{})
	workspace.SetStrings(doc, "command_surface", []string{})
	workspace.SetStrings(doc, "mandatory_validation", []string{})

	paths := map[domain.ID]string{contractID: relPath}
	return s.apply("contract.create:"+contractID.String(), []*workspace.Document{doc}, paths, "contract.create", ref, relPath)
}

func resolveContract(ws *discovery.Workspace, ref string) (Binding, error) {
	entry, err := ws.Lookup(ref, domain.Contract)
	if err != nil {
		return Binding{}, invalidCause("contract.ref", "must be a valid Contract reference (e.g. CC-* or Contract:<UUIDv7>)", err)
	}
	data, err := os.ReadFile(entry.Absolute)
	if err != nil {
		return Binding{}, err
	}
	digest := sha256.Sum256(data)
	canonicalRef := string(domain.Contract) + ":" + entry.Document.Record.ID.String()
	return Binding{Ref: canonicalRef, Fingerprint: "sha256:" + hex.EncodeToString(digest[:])}, nil
}
