// check-architecture verifies production imports for every platform, without
// compiling or executing the product. Test fixtures may compose real adapters.
package main

import (
	"encoding/json"
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type boundary struct {
	Internal []string  `json:"internal"`
	External *[]string `json:"external,omitempty"`
}
type policy struct {
	Module   string              `json:"module"`
	Packages map[string]boundary `json:"packages"`
}

func main() {
	root := "."
	if len(os.Args) > 1 {
		root = os.Args[1]
	}
	if err := check(root); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("Architecture boundaries: PASS")
}

func check(root string) error {
	data, err := os.ReadFile(filepath.Join(root, "scripts", "architecture-boundaries.json"))
	if err != nil {
		return err
	}
	var p policy
	decoder := json.NewDecoder(strings.NewReader(string(data)))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&p); err != nil {
		return err
	}
	if p.Module == "" || len(p.Packages) == 0 {
		return fmt.Errorf("architecture policy requires module and packages")
	}
	edges := map[string][]string{}
	seen := map[string]bool{}
	err = filepath.WalkDir(filepath.Join(root, "internal"), func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		rel, err := filepath.Rel(filepath.Join(root, "internal"), filepath.Dir(path))
		if err != nil {
			return err
		}
		name := filepath.ToSlash(rel)
		rule, ok := p.Packages[name]
		if !ok {
			return fmt.Errorf("%s: package %s has no architecture policy", path, name)
		}
		seen[name] = true
		file, err := parser.ParseFile(token.NewFileSet(), path, nil, parser.ImportsOnly)
		if err != nil {
			return err
		}
		for _, spec := range file.Imports {
			imp, err := strconv.Unquote(spec.Path.Value)
			if err != nil {
				return err
			}
			internalPrefix := p.Module + "/internal/"
			if strings.HasPrefix(imp, internalPrefix) {
				target := strings.TrimPrefix(imp, internalPrefix)
				if _, ok := p.Packages[target]; !ok {
					return fmt.Errorf("%s: import %s has no architecture policy", path, target)
				}
				if !contains(rule.Internal, target) {
					return fmt.Errorf("%s: forbidden dependency %s -> %s", path, name, target)
				}
				edges[name] = append(edges[name], target)
			} else if strings.HasPrefix(imp, p.Module+"/") {
				return fmt.Errorf("%s: internal package %s must not import product composition %s", path, name, imp)
			} else if rule.External != nil && !contains(*rule.External, imp) {
				return fmt.Errorf("%s: forbidden external dependency %s -> %s", path, name, imp)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	names := make([]string, 0, len(p.Packages))
	for name := range p.Packages {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		if !seen[name] {
			return fmt.Errorf("architecture policy names missing production package %s", name)
		}
	}
	return acyclic(edges)
}
func contains(values []string, value string) bool {
	for _, candidate := range values {
		if candidate == value {
			return true
		}
	}
	return false
}
func acyclic(edges map[string][]string) error {
	state := map[string]int{}
	var visit func(string) error
	visit = func(name string) error {
		if state[name] == 1 {
			return fmt.Errorf("architecture import cycle at %s", name)
		}
		if state[name] == 2 {
			return nil
		}
		state[name] = 1
		for _, target := range edges[name] {
			if err := visit(target); err != nil {
				return err
			}
		}
		state[name] = 2
		return nil
	}
	names := make([]string, 0, len(edges))
	for name := range edges {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		if err := visit(name); err != nil {
			return err
		}
	}
	return nil
}
