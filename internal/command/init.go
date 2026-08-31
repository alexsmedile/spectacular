package command

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alexsmedile/spectacular/v2/internal/domain"
)

type InitResult struct {
	Root               string   `json:"root"`
	MetadataDir        string   `json:"metadata_dir"`
	CreatedFiles       []string `json:"created_files"`
	SkippedFiles       []string `json:"skipped_files"`
	AlreadyInitialized bool     `json:"already_initialized"`
}

const defaultWorkspaceYAML = `schema_version: spectacular.workspace.v1
record_roots:
  - .
project_anchor: PROJECT.md
`

func defaultProjectAnchor(id domain.ID, title string) string {
	if title == "" {
		title = "Project"
	}
	return fmt.Sprintf(`---
type: Anchor
id: %s
title: %s
current_truth:
  - .spectacular/PROJECT.md
---

# %s

## Purpose
What this project does and why it exists.

## Boundaries
- In scope: Core functionality.
- Out of scope: Non-goals.
`, id.String(), title, title)
}

// InitWorkspace initializes a Spectacular workspace in targetDir.
// It avoids overwriting existing files in .spectacular.
func InitWorkspace(targetDir, projectName string) (InitResult, error) {
	abs, err := filepath.Abs(targetDir)
	if err != nil {
		return InitResult{}, fmt.Errorf("resolve target directory: %w", err)
	}

	info, err := os.Stat(abs)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(abs, 0o755); err != nil {
				return InitResult{}, fmt.Errorf("create target directory: %w", err)
			}
		} else {
			return InitResult{}, fmt.Errorf("inspect target directory: %w", err)
		}
	} else if !info.IsDir() {
		return InitResult{}, fmt.Errorf("target path is not a directory: %s", abs)
	}

	metaDir := filepath.Join(abs, ".spectacular")
	if err := os.MkdirAll(metaDir, 0o755); err != nil {
		return InitResult{}, fmt.Errorf("create metadata directory: %w", err)
	}

	// Standard collection directories
	standardDirs := []string{
		"contracts",
		"proposals",
		"missions",
		"decisions",
		"evidence",
		"gaps",
		"archive/missions",
	}
	for _, dir := range standardDirs {
		dirPath := filepath.Join(metaDir, filepath.FromSlash(dir))
		if err := os.MkdirAll(dirPath, 0o755); err != nil {
			return InitResult{}, fmt.Errorf("create collection directory %s: %w", dir, err)
		}
	}

	var createdFiles []string
	var skippedFiles []string

	// 1. workspace.yaml
	wsMarker := filepath.Join(metaDir, "workspace.yaml")
	if _, err := os.Stat(wsMarker); err == nil {
		skippedFiles = append(skippedFiles, ".spectacular/workspace.yaml")
	} else if os.IsNotExist(err) {
		if err := os.WriteFile(wsMarker, []byte(defaultWorkspaceYAML), 0o644); err != nil {
			return InitResult{}, fmt.Errorf("write workspace.yaml: %w", err)
		}
		createdFiles = append(createdFiles, ".spectacular/workspace.yaml")
	} else {
		return InitResult{}, fmt.Errorf("inspect workspace.yaml: %w", err)
	}

	// 2. PROJECT.md
	projectAnchorPath := filepath.Join(metaDir, "PROJECT.md")
	if _, err := os.Stat(projectAnchorPath); err == nil {
		skippedFiles = append(skippedFiles, ".spectacular/PROJECT.md")
	} else if os.IsNotExist(err) {
		name := projectName
		if strings.TrimSpace(name) == "" {
			name = filepath.Base(abs)
			if name == "" || name == "." || name == "/" {
				name = "Project"
			}
		}
		id, err := domain.NewID()
		if err != nil {
			return InitResult{}, fmt.Errorf("generate anchor ID: %w", err)
		}
		if err := os.WriteFile(projectAnchorPath, []byte(defaultProjectAnchor(id, name)), 0o644); err != nil {
			return InitResult{}, fmt.Errorf("write PROJECT.md: %w", err)
		}
		createdFiles = append(createdFiles, ".spectacular/PROJECT.md")
	} else {
		return InitResult{}, fmt.Errorf("inspect PROJECT.md: %w", err)
	}

	return InitResult{
		Root:               abs,
		MetadataDir:        metaDir,
		CreatedFiles:       createdFiles,
		SkippedFiles:       skippedFiles,
		AlreadyInitialized: len(createdFiles) == 0,
	}, nil
}
