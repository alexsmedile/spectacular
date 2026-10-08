// Repository utility: go run ./scripts/rebuild-workspace-index.go [--write].
// Uses governed discovery; it does not index soft context or certify bindings.
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alexsmedile/spectacular/v2/internal/discovery"
	"github.com/alexsmedile/spectacular/v2/internal/humanlayout"
)

func main() {
	write := flag.Bool("write", false, "replace only the generated root index.json; default prints JSON")
	flag.Parse()
	if err := rebuild(*write); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func rebuild(write bool) error {
	ws, err := discovery.Open(".")
	if err != nil {
		return err
	}
	indexes, err := humanlayout.Indexes(ws.Entries, nil, nil)
	if err != nil {
		return err
	}
	data := indexes[".spectacular/index.json"]
	if !write {
		_, err = os.Stdout.Write(data)
		return err
	}
	path := filepath.Join(ws.MetadataDir, "index.json")
	if info, statErr := os.Lstat(path); statErr == nil && !info.Mode().IsRegular() {
		return fmt.Errorf("refuse non-regular inventory: %s", path)
	}
	old, err := os.ReadFile(path)
	if err == nil {
		var claim struct {
			Generated bool   `json:"generated"`
			Scope     string `json:"scope"`
		}
		if json.Unmarshal(old, &claim) != nil || !claim.Generated || claim.Scope != "governed-records" {
			return fmt.Errorf("refuse authored or unrecognized inventory: %s", path)
		}
		if bytes.Equal(old, data) {
			return nil
		}
	} else if !os.IsNotExist(err) {
		return err
	}
	temp, err := os.CreateTemp(ws.MetadataDir, ".index-*.json")
	if err != nil {
		return err
	}
	defer os.Remove(temp.Name())
	if _, err = temp.Write(data); err != nil {
		temp.Close()
		return err
	}
	if err = temp.Chmod(0644); err != nil {
		temp.Close()
		return err
	}
	if err = temp.Close(); err != nil {
		return err
	}
	return os.Rename(temp.Name(), path)
}
