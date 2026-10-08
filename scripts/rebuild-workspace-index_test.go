package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/alexsmedile/spectacular/v2/internal/command"
)

func TestRebuildPreservesManualNavigationAndRefusesAuthoredInventory(t *testing.T) {
	root := t.TempDir()
	if _, err := command.InitWorkspace(root, "Inventory test"); err != nil {
		t.Fatal(err)
	}
	t.Chdir(root)
	manual := filepath.Join(root, ".spectacular", "INDEX.md")
	manualData := []byte("# My reading routes\n")
	if err := os.WriteFile(manual, manualData, 0644); err != nil {
		t.Fatal(err)
	}
	if err := rebuild(true); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, ".spectacular", "index.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var inventory struct {
		Scope   string `json:"scope"`
		Entries []struct {
			Path string `json:"path"`
		} `json:"entries"`
	}
	if err := json.Unmarshal(data, &inventory); err != nil {
		t.Fatal(err)
	}
	if inventory.Scope != "governed-records" || len(inventory.Entries) != 1 || inventory.Entries[0].Path != ".spectacular/PROJECT.md" {
		t.Fatalf("unexpected inventory: %s", data)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := rebuild(true); err != nil {
		t.Fatal(err)
	}
	again, err := os.Stat(path)
	if err != nil || !info.ModTime().Equal(again.ModTime()) {
		t.Fatal("unchanged inventory was rewritten", err)
	}
	authored := []byte("{\"purpose\":\"authored data\"}\n")
	if err := os.WriteFile(path, authored, 0644); err != nil {
		t.Fatal(err)
	}
	if err := rebuild(true); err == nil {
		t.Fatal("accepted replacement of authored JSON")
	}
	current, err := os.ReadFile(path)
	if err != nil || !bytes.Equal(current, authored) {
		t.Fatal("authored inventory changed", err)
	}
	current, err = os.ReadFile(manual)
	if err != nil || !bytes.Equal(current, manualData) {
		t.Fatal("manual navigation changed", err)
	}
}
