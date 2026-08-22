// Package campaign reads the lightweight Campaign planning projection.
package campaign

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/alexsmedile/spectacular/v2/internal/discovery"
	"github.com/alexsmedile/spectacular/v2/internal/domain"
	"go.yaml.in/yaml/v3"
)

type Block struct {
	Ref       string   `yaml:"ref" json:"ref"`
	Title     string   `yaml:"title" json:"title"`
	DependsOn []string `yaml:"after,omitempty" json:"after,omitempty"`
	Missions  []string `yaml:"missions,omitempty" json:"missions,omitempty"`
	State     string   `yaml:"state" json:"state"`
}

type Check struct {
	Path          string   `json:"path"`
	Title         string   `json:"title"`
	StrategicGoal string   `json:"strategic_goal,omitempty"`
	Current       string   `json:"current"`
	CurrentBlock  Block    `json:"current_block"`
	Next          []Block  `json:"next,omitempty"`
	ExitCondition string   `json:"exit_condition,omitempty"`
	Blocks        []Block  `json:"blocks"`
	Order         []string `json:"order"`
	Mermaid       string   `json:"mermaid"`
	EmbeddedMap   string   `json:"embedded_mermaid"`
}

type manifest struct {
	CampaignSchema string  `yaml:"campaign_schema"`
	Title          string  `yaml:"title"`
	Focus          string  `yaml:"focus"`
	Current        string  `yaml:"current"`
	ExitCondition  string  `yaml:"exit_condition"`
	Blocks         []Block `yaml:"blocks"`
}

// Validate reads one Campaign planning document. Campaigns remain free-form
// Markdown, but this command recognizes the compact frontmatter map documented in
// campaigns/README.md so it can validate and render the dependency projection.
func Validate(ws *discovery.Workspace, input string) (Check, error) {
	path, err := campaignPath(ws, input)
	if err != nil {
		return Check{}, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return Check{}, domain.NewRefusal(domain.RefusalInvalidWorkspacePath, input, "read Campaign", err)
	}
	check, err := parse(string(data))
	if err != nil {
		return Check{}, domain.NewRefusal(domain.RefusalInvalidKnownField, input, err.Error(), nil)
	}
	for _, block := range check.Blocks {
		for _, ref := range block.Missions {
			if _, err := ws.Lookup(ref, domain.Mission); err != nil {
				return Check{}, domain.NewRefusal(domain.RefusalTargetNotFound, ref, "Campaign block "+block.Title+" names no Mission", err)
			}
		}
	}
	check.Path = relative(ws.Root, path)
	return check, nil
}

func campaignPath(ws *discovery.Workspace, input string) (string, error) {
	if input == "" || input == "-" {
		return "", domain.NewRefusal(domain.RefusalInvalidWorkspacePath, input, "Campaign path is required", nil)
	}
	path := input
	if !filepath.IsAbs(path) {
		path = filepath.Join(ws.Root, path)
	}
	real, err := filepath.EvalSymlinks(path)
	if err != nil {
		return "", domain.NewRefusal(domain.RefusalInvalidWorkspacePath, input, "Campaign path is missing or unreadable", err)
	}
	root, err := filepath.EvalSymlinks(filepath.Join(ws.MetadataDir, "campaigns"))
	if err != nil {
		return "", domain.NewRefusal(domain.RefusalInvalidWorkspacePath, input, "Campaign directory is missing", err)
	}
	if !within(root, real) || filepath.Ext(real) != ".md" || filepath.Base(real) == "README.md" {
		return "", domain.NewRefusal(domain.RefusalInvalidWorkspacePath, input, "Campaign must be a Markdown file under .spectacular/campaigns", nil)
	}
	return real, nil
}

func parse(text string) (Check, error) {
	frontmatter, err := campaignFrontmatter(text)
	if err != nil {
		return Check{}, err
	}
	var source manifest
	if err := yaml.Unmarshal([]byte(frontmatter), &source); err != nil {
		return Check{}, fmt.Errorf("decode Campaign frontmatter: %w", err)
	}
	if source.CampaignSchema != "spectacular.campaign.v1" {
		return Check{}, fmt.Errorf("campaign_schema must be spectacular.campaign.v1")
	}
	if source.Title == "" || source.Focus == "" || source.Current == "" || source.ExitCondition == "" {
		return Check{}, fmt.Errorf("Campaign frontmatter requires title, focus, current, and exit_condition")
	}
	if len(source.Blocks) == 0 {
		return Check{}, fmt.Errorf("Campaign frontmatter requires at least one block")
	}
	blocks, err := validateBlocks(source.Blocks, source.Current)
	if err != nil {
		return Check{}, err
	}
	order, err := topological(blocks)
	if err != nil {
		return Check{}, err
	}
	current, next := mapPosition(source.Current, blocks)
	graph := mermaid(source.Title, blocks)
	embedded, err := embeddedMermaid(text, graph)
	if err != nil {
		return Check{}, err
	}
	return Check{Title: source.Title, StrategicGoal: source.Focus, Current: source.Current, CurrentBlock: current, Next: next, ExitCondition: source.ExitCondition, Blocks: blocks, Order: order, Mermaid: graph, EmbeddedMap: embedded}, nil
}

func mapPosition(current string, blocks []Block) (Block, []Block) {
	var position Block
	next := []Block{}
	for _, block := range blocks {
		if block.Ref == current {
			position = block
		}
		for _, dependency := range block.DependsOn {
			if dependency == current {
				next = append(next, block)
				break
			}
		}
	}
	sort.Slice(next, func(i, j int) bool { return next[i].Ref < next[j].Ref })
	return position, next
}

func embeddedMermaid(text, graph string) (string, error) {
	const start = "<!-- spectacular:campaign-mermaid:start -->"
	const end = "<!-- spectacular:campaign-mermaid:end -->"
	from := strings.Index(text, start)
	to := strings.Index(text, end)
	if from < 0 && to < 0 {
		return "absent", nil
	}
	if from < 0 || to < 0 || to < from {
		return "", fmt.Errorf("embedded Campaign Mermaid must use matching start and end markers")
	}
	actual := strings.TrimSpace(text[from+len(start) : to])
	expected := strings.TrimSpace("```mermaid\n" + graph + "```")
	if actual != expected {
		return "", fmt.Errorf("embedded Campaign Mermaid is stale; re-render it from campaign check output")
	}
	return "current", nil
}

func campaignFrontmatter(text string) (string, error) {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	if !strings.HasPrefix(text, "---\n") {
		return "", fmt.Errorf("requires Campaign frontmatter")
	}
	frontmatter, _, found := strings.Cut(strings.TrimPrefix(text, "---\n"), "\n---\n")
	if !found {
		return "", fmt.Errorf("Campaign frontmatter must close with ---")
	}
	return frontmatter, nil
}

func validateBlocks(blocks []Block, current string) ([]Block, error) {
	seen := map[string]bool{}
	for i, block := range blocks {
		if block.Ref == "" || block.Ref != fmt.Sprintf("B%d", i+1) || block.Title == "" {
			return nil, fmt.Errorf("blocks must use ordered refs B1, B2, ... and titles")
		}
		if seen[block.Ref] {
			return nil, fmt.Errorf("blocks must have unique refs")
		}
		if block.State != "planned" && block.State != "active" && block.State != "complete" && block.State != "paused" {
			return nil, fmt.Errorf("Block %q state must be planned, active, complete, or paused", block.Ref)
		}
		seen[block.Ref] = true
	}
	if !seen[current] {
		return nil, fmt.Errorf("current names unknown Block %q", current)
	}
	for _, block := range blocks {
		for _, dependency := range block.DependsOn {
			if !seen[dependency] {
				return nil, fmt.Errorf("Block %q depends on unknown Block %q", block.Ref, dependency)
			}
		}
	}
	return blocks, nil
}

func topological(blocks []Block) ([]string, error) {
	byName := map[string]Block{}
	degrees := map[string]int{}
	children := map[string][]string{}
	for _, block := range blocks {
		byName[block.Ref] = block
		degrees[block.Ref] = len(block.DependsOn)
		for _, dep := range block.DependsOn {
			children[dep] = append(children[dep], block.Ref)
		}
	}
	ready := []string{}
	for _, block := range blocks {
		if degrees[block.Ref] == 0 {
			ready = append(ready, block.Ref)
		}
	}
	order := make([]string, 0, len(blocks))
	for len(ready) > 0 {
		sort.Strings(ready)
		next := ready[0]
		ready = ready[1:]
		order = append(order, next)
		for _, child := range children[next] {
			degrees[child]--
			if degrees[child] == 0 {
				ready = append(ready, child)
			}
		}
	}
	if len(order) != len(blocks) {
		return nil, fmt.Errorf("Blocks dependency graph contains a cycle")
	}
	return order, nil
}

func mermaid(title string, blocks []Block) string {
	var out strings.Builder
	out.WriteString("flowchart LR\n")
	for i, block := range blocks {
		fmt.Fprintf(&out, "  B%d[\"%s\\n%s\"]\n", i+1, escape(block.Title), escape(block.State))
	}
	positions := map[string]int{}
	for i, block := range blocks {
		positions[block.Ref] = i + 1
	}
	for i, block := range blocks {
		for _, dep := range block.DependsOn {
			fmt.Fprintf(&out, "  B%d --> B%d\n", positions[dep], i+1)
		}
	}
	_ = title
	return out.String()
}

func escape(value string) string {
	return strings.ReplaceAll(strings.ReplaceAll(value, "\"", "'"), "\n", " ")
}
func within(root, path string) bool {
	relative, err := filepath.Rel(root, path)
	return err == nil && relative != ".." && !strings.HasPrefix(relative, ".."+string(filepath.Separator))
}
func relative(root, path string) string {
	value, err := filepath.Rel(root, path)
	if err != nil {
		return path
	}
	return filepath.ToSlash(value)
}
