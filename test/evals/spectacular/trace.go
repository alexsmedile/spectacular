package spectaculareval

import (
	"bufio"
	"encoding/json"
	"path/filepath"
	"regexp"
	"strings"
)

var tracePathPattern = regexp.MustCompile(`[A-Za-z0-9_./-]+\.md`)
var shellSegmentPattern = regexp.MustCompile(`&&|;|\|`)

type TraceCertification struct {
	Valid            bool     `json:"valid"`
	UsageObserved    bool     `json:"usage_observed"`
	SemanticObserved bool     `json:"semantic_observed"`
	ToolCalls        int      `json:"tool_calls"`
	Findings         []string `json:"findings,omitempty"`
}

func CertifyTrace(trace string, requireTools bool) TraceCertification {
	metrics := ParseTraceMetrics(trace)
	result := TraceCertification{UsageObserved: metrics.UsageObserved, SemanticObserved: metrics.SemanticObserved, ToolCalls: metrics.ToolCalls}
	if !metrics.UsageObserved {
		result.Findings = append(result.Findings, "usage telemetry missing")
	}
	if !metrics.SemanticObserved {
		result.Findings = append(result.Findings, "semantic tool telemetry missing or self-reported only")
	}
	if requireTools && metrics.ToolCalls == 0 {
		result.Findings = append(result.Findings, "expected tool telemetry but observed zero tool calls")
	}
	result.Valid = len(result.Findings) == 0
	return result
}

func ParseTraceMetrics(trace string) TraceMetrics {
	var metrics TraceMetrics
	scanner := bufio.NewScanner(strings.NewReader(trace))
	scanner.Buffer(make([]byte, 64*1024), 4*1024*1024)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || !strings.HasPrefix(line, "{") {
			continue
		}
		var event any
		if json.Unmarshal([]byte(line), &event) != nil {
			continue
		}
		metrics.Events++
		collectSemanticEvent(event, &metrics)
		collectNativeTelemetry(event, &metrics)
		walkTrace(event, &metrics)
	}
	return metrics
}

func collectNativeTelemetry(value any, metrics *TraceMetrics) {
	event, ok := value.(map[string]any)
	if !ok {
		return
	}
	eventType, _ := event["type"].(string)
	if eventType == "thread.started" || eventType == "turn.completed" {
		// Codex JSONL is a complete host trace. An observed turn can validly
		// contain zero tool calls; absence is still adapter-authored evidence.
		metrics.SemanticObserved = true
	}
	if item, ok := event["item"].(map[string]any); ok {
		if itemType, _ := item["type"].(string); itemType == "command_execution" {
			if command, _ := item["command"].(string); command != "" {
				metrics.SemanticObserved = true
				metrics.SemanticEvents++
				metrics.ObservedCommands = append(metrics.ObservedCommands, command)
				recordCommandReads(command, metrics)
			}
		}
	}
	message, _ := event["message"].(map[string]any)
	if content, ok := message["content"].([]any); ok {
		for _, rawBlock := range content {
			block, _ := rawBlock.(map[string]any)
			if blockType, _ := block["type"].(string); blockType != "tool_use" {
				continue
			}
			metrics.SemanticObserved = true
			metrics.SemanticEvents++
			name, _ := block["name"].(string)
			input, _ := block["input"].(map[string]any)
			switch name {
			case "Read", "NotebookRead":
				path, _ := input["file_path"].(string)
				if path == "" {
					path, _ = input["notebook_path"].(string)
				}
				recordObservedPath(path, metrics)
			case "Bash":
				if command, _ := input["command"].(string); command != "" {
					metrics.ObservedCommands = append(metrics.ObservedCommands, command)
					recordCommandReads(command, metrics)
				}
			case "Skill":
				if skill, _ := input["skill"].(string); skill != "" {
					metrics.ObservedReferences = append(metrics.ObservedReferences, skill)
				}
			}
		}
	}
	if eventType == "result" {
		if usage, ok := event["usage"].(map[string]any); ok {
			input := number(usage["input_tokens"]) + number(usage["cache_creation_input_tokens"]) + number(usage["cache_read_input_tokens"])
			if input > 0 {
				metrics.InputTokens = input
				metrics.CachedInputTokens = number(usage["cache_read_input_tokens"])
				metrics.OutputTokens = number(usage["output_tokens"])
				metrics.UsageObserved = true
			}
		}
	}
}

func recordCommandReads(command string, metrics *TraceMetrics) {
	for _, segment := range shellSegmentPattern.Split(command, -1) {
		lower := strings.ToLower(segment)
		if !strings.Contains(lower, "sed ") && !strings.Contains(lower, "cat ") && !strings.Contains(lower, "head ") && !strings.Contains(lower, "tail ") && !strings.Contains(lower, "bat ") {
			continue
		}
		for _, path := range tracePathPattern.FindAllString(segment, -1) {
			recordObservedPath(path, metrics)
		}
	}
}

func recordObservedPath(path string, metrics *TraceMetrics) {
	path = strings.Trim(path, "'\"`.,:()[]")
	if path == "" {
		return
	}
	metrics.ObservedFiles = append(metrics.ObservedFiles, path)
	normalized := filepath.ToSlash(path)
	base := filepath.Base(normalized)
	if strings.Contains(normalized, "/references/") || primaryReferenceNames[base] || base == "SKILL.md" {
		metrics.ObservedReferences = append(metrics.ObservedReferences, path)
	}
}

func number(value any) int {
	if result, ok := jsonNumber(value); ok {
		return result
	}
	return 0
}

func collectSemanticEvent(value any, metrics *TraceMetrics) {
	event, ok := value.(map[string]any)
	if !ok {
		return
	}
	eventType, _ := event["type"].(string)
	path, _ := event["path"].(string)
	command, _ := event["command"].(string)
	switch eventType {
	case "spectacular.eval.observations":
		metrics.SemanticEvents++
		metrics.SemanticObserved = true
		metrics.ObservedFiles = append(metrics.ObservedFiles, stringList(event["files_read"])...)
		metrics.ObservedReferences = append(metrics.ObservedReferences, stringList(event["references_loaded"])...)
		metrics.ObservedCommands = append(metrics.ObservedCommands, stringList(event["commands_run"])...)
	case "spectacular.eval.file_read":
		if path != "" {
			metrics.SemanticEvents++
			metrics.SemanticObserved = true
			recordObservedPath(path, metrics)
		}
	case "spectacular.eval.reference_loaded":
		if path != "" {
			metrics.SemanticEvents++
			metrics.SemanticObserved = true
			metrics.ObservedReferences = append(metrics.ObservedReferences, path)
		}
	case "spectacular.eval.command":
		if command != "" {
			metrics.SemanticEvents++
			metrics.SemanticObserved = true
			metrics.ObservedCommands = append(metrics.ObservedCommands, command)
		}
	}
}

func stringList(value any) []string {
	items, ok := value.([]any)
	if !ok {
		return nil
	}
	result := make([]string, 0, len(items))
	for _, item := range items {
		if text, ok := item.(string); ok {
			result = append(result, text)
		}
	}
	return result
}

func walkTrace(value any, metrics *TraceMetrics) {
	switch typed := value.(type) {
	case map[string]any:
		for key, item := range typed {
			switch key {
			case "input_tokens":
				if number, ok := jsonNumber(item); ok && number > metrics.InputTokens {
					metrics.InputTokens = number
					metrics.UsageObserved = true
				}
			case "cached_input_tokens":
				if number, ok := jsonNumber(item); ok && number > metrics.CachedInputTokens {
					metrics.CachedInputTokens = number
				}
			case "output_tokens":
				if number, ok := jsonNumber(item); ok && number > metrics.OutputTokens {
					metrics.OutputTokens = number
				}
			case "type":
				if text, ok := item.(string); ok && isToolEvent(text) {
					metrics.ToolCalls++
				}
			}
			walkTrace(item, metrics)
		}
	case []any:
		for _, item := range typed {
			walkTrace(item, metrics)
		}
	}
}

func jsonNumber(value any) (int, bool) {
	number, ok := value.(float64)
	return int(number), ok
}

func isToolEvent(eventType string) bool {
	lower := strings.ToLower(eventType)
	return strings.Contains(lower, "tool_call") || strings.Contains(lower, "tool_use") || strings.Contains(lower, "command_execution") || strings.Contains(lower, "mcp_call")
}
