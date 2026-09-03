package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/alexsmedile/spectacular/v2/internal/charter/tokenizer"
)

type MissionStats struct {
	Ref         string
	Path        string
	PlanTokens  int
	TotalFiles  int
	TotalTokens int
}

func main() {
	missions := make(map[string]*MissionStats)

	err := filepath.Walk(".spectacular/missions", func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || filepath.Ext(path) != ".md" || filepath.Base(path) == "index.md" {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		tokens, err := tokenizer.Count(string(data))
		if err != nil {
			return nil
		}

		rel, _ := filepath.Rel(".spectacular/missions", path)
		parts := strings.Split(rel, string(filepath.Separator))
		dirName := parts[0]
		ref := strings.Split(dirName, "-")[0]

		stats, ok := missions[ref]
		if !ok {
			stats = &MissionStats{Ref: ref}
			missions[ref] = stats
		}
		stats.TotalFiles++
		stats.TotalTokens += tokens

		if filepath.Base(path) == dirName+".md" {
			stats.PlanTokens = tokens
			stats.Path = path
		}
		return nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Walk error: %v\n", err)
		os.Exit(1)
	}

	var list []*MissionStats
	for _, s := range missions {
		list = append(list, s)
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Ref < list[j].Ref
	})

	fmt.Printf("%-6s | %-12s | %-12s | %-10s | %s\n", "Ref", "Plan Tokens", "Bundle Tokens", "Files", "Status vs Active Budget (400-900)")
	fmt.Println(strings.Repeat("-", 80))
	for _, s := range list {
		status := "Pass"
		if s.PlanTokens > 900 && s.PlanTokens <= 1200 {
			status = "Warn (Upper Envelope)"
		} else if s.PlanTokens > 1200 {
			status = "Exceeds 1200 Cap"
		} else if s.PlanTokens < 400 {
			status = "Under (Lean / Sketch)"
		}
		fmt.Printf("%-6s | %-12d | %-12d | %-10d | %s\n", s.Ref, s.PlanTokens, s.TotalTokens, s.TotalFiles, status)
	}
}
